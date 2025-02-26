package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Request struct {
	Method string
	Params interface{}
}

type Response struct {
	Result interface{}
	Error  string
}

type Service struct {
	MethodList         map[string]interface{}
	Listener           *net.TCPListener
	Config             Config
	ServiceName        string
	ServiceFindingConn *net.TCPConn
}

type Config struct {
	ServiceAddr        NetField `json:"serviceAddr"`
	ServiceFindingAdrr NetField `json:"serviceFindingAddr"`
}

type NetField struct {
	Ip      string `json:"ip"`
	Port    int    `json:"port"`
	Execute bool   `json:"execute"`
}

type FindingRequest struct {
	ReqType     string `json:"reqType"`
	Addr        string `json:"addr"`
	ServiceName string `json:"serviceName"`
}

type FindingResponse struct {
	ServiceName string `json:"serviceName"`
	ErrMsg      string `json:"errMsg"`
	Status      string `json:"status"`
	Addr        string `json:"addr"`
}

func CreateService() *Service {
	server := &Service{
		MethodList: make(map[string]interface{}),
	}
	return server
}

func (s *Service) GetConfig(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	json.Unmarshal(byteValue, &s.Config)
	return nil
}

func (s *Service) AddMethod(name string, Method interface{}) {
	s.MethodList[name] = Method

}

func (s *Service) RegisterService(serviceAddr string) (err error) {
	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	ip, port := SplitAddr(listener.Addr().String())
	s.Config.ServiceAddr.Ip = ip
	s.Config.ServiceAddr.Port = port

	s.Listener = listener

	return nil
}

func (s *Service) ServiceFinding() error {
	serviceFindingAddr := GetAddr(s.Config.ServiceFindingAdrr.Ip, s.Config.ServiceFindingAdrr.Port)
	localAddr := GetAddr(s.Config.ServiceAddr.Ip, s.Config.ServiceAddr.Port)

	addr, err := net.ResolveTCPAddr("tcp", serviceFindingAddr)
	if err != nil {
		return err
	}

	s.ServiceFindingConn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	req := FindingRequest{
		ReqType:     "connect",
		Addr:        localAddr,
		ServiceName: s.ServiceName,
	}
	var rsp FindingResponse

	reqBuf, _ := json.Marshal(req)
	reqBufLen := len(reqBuf)
	reqBufHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(reqBufHeader, uint32(reqBufLen))

	reqBuf = append(reqBufHeader, reqBuf...)
	s.ServiceFindingConn.Write(reqBuf)

	rspBufHeader := make([]byte, 4)
	_, err = s.ServiceFindingConn.Read(rspBufHeader)
	if err != nil {
		return err
	}
	rspBufLen := binary.BigEndian.Uint32(rspBufHeader)

	rspBuf := make([]byte, rspBufLen)
	s.ServiceFindingConn.Read(rspBuf)
	json.Unmarshal(rspBuf, &rsp)

	if rsp.Status != "ok" {
		return fmt.Errorf(rsp.ErrMsg)
	}

	return nil

}

func (s *Service) Service() {
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			continue
		}

		go s.HandleServerConnection(conn)
	}
}

func (s *Service) HandleServerConnection(conn *net.TCPConn) {
	defer conn.Close()
	for {
		// 收包
		header := make([]byte, 4)
		_, err := receiveBag(conn, header)
		if err != nil {
			break
		}
		bufLen := binary.BigEndian.Uint32(header)
		buf := make([]byte, bufLen)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		log.Println("recieve:", string(buf))

		// 处理包
		backBuf, err := s.HandlerBuf(buf[:n])
		if err != nil {
			log.Println(err)
			break
		}

		// 发包
		log.Println("send:", string(backBuf))
		backBufLen := uint32(len(backBuf))
		backBufHeader := make([]byte, 4)
		binary.BigEndian.PutUint32(backBufHeader, backBufLen)
		backBuf = append(backBufHeader, backBuf...)
		conn.Write(backBuf)
	}
}

func (s *Service) HandlerBuf(buf []byte) (backBuf []byte, err error) {
	var req Request
	if err := json.Unmarshal(buf, &req); err != nil {
		return nil, err
	}

	// 调用方法处理请求
	if serviceFunc := s.MethodList[req.Method]; serviceFunc != nil {
		return s.CallServiceMethod(serviceFunc, req.Params)

	} else {
		return nil, fmt.Errorf("method %s not found", req.Method)
	}

}

func (s *Service) CallServiceMethod(serviceFunc interface{}, params interface{}) ([]byte, error) {
	funcValue := reflect.ValueOf(serviceFunc)

	// 确保 serviceFunc 是一个函数
	if funcValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("serviceFunc must be a function")
	}

	// 获取函数的参数类型
	funcType := funcValue.Type()
	if funcType.NumIn() != 1 {
		return nil, fmt.Errorf("serviceFunc must have exactly one input parameter")
	}

	paramType := funcType.In(0)

	// 创建指定类型的实例
	paramValue := reflect.New(paramType).Interface()

	// 将 params 转换为 JSON 并解析到 paramValue
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %v", err)
	}

	if err := json.Unmarshal(paramsBytes, paramValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal params: %v", err)
	}

	// 调用方法
	funcParams := []reflect.Value{reflect.ValueOf(paramValue).Elem()}
	results := funcValue.Call(funcParams)

	// 返回值处理
	resultsInterfaces := results[0].Interface()

	return json.Marshal(Response{Result: resultsInterfaces, Error: ""})

}

func receiveBag(conn net.Conn, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := conn.Read(buf[total:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func GetAddr(ip string, port int) (addr string) {
	addr = ip + ":" + strconv.Itoa(port)
	return addr
}

func SplitAddr(addr string) (ip string, port int) {
	addrArr := strings.Split(addr, ":")
	ip = addrArr[0]
	port, _ = strconv.Atoi(addrArr[1])
	return ip, port
}
