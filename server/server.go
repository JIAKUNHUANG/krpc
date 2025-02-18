package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type Request struct {
	Method string
	Params []interface{}
}

type Response struct {
	Result []interface{}
	Error  string
}

type Service struct {
	MethodList  map[string]interface{}
	Listener *net.TCPListener
}

func CreateService() *Service {
	server := &Service{
		MethodList: make(map[string]interface{}),
	}
	return server
}

func (s *Service) AddMethod(name string, Method interface{}) {
	s.MethodList[name] = Method

}

func (s *Service) RegisterService(addr *net.TCPAddr) (err error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	s.Listener = listener

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

		// 处理包
		backBuf, err := s.HandlerBuf(buf[:n])
		if err != nil {
			break
		}

		// 发包
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
	fmt.Println(req.Method)
	if serviceFunc := s.MethodList[req.Method]; serviceFunc != nil {
		fmt.Println(serviceFunc)
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
	resultsInterfaces := make([]interface{}, len(results))
	for i, result := range results {
		resultsInterfaces[i] = result.Interface()
	}

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
