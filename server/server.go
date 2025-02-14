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

type Server struct {
	Service  map[string]interface{}
	Listener *net.TCPListener
}

func CreateServer() *Server {
	server := &Server{
		Service: make(map[string]interface{}),
	}
	return server
}

func (s *Server) AddServer(name string, function interface{}) {
	s.Service[name] = function

}

func (s *Server) RegisterServer(addr *net.TCPAddr) (err error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	s.Listener = listener

	return nil
}

func (s *Server) Server() {
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			continue
		}

		go s.HandleServerConnection(conn)
	}
}

func (s *Server) HandleServerConnection(conn *net.TCPConn) {
	defer conn.Close()
	for {
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
		backBuf, err := s.HandlerBuf(buf[:n])
		if err != nil {
			break
		}
		backBufLen := uint32(len(backBuf))
		backBufHeader := make([]byte, 4)
		binary.BigEndian.PutUint32(backBufHeader, backBufLen)
		backBuf = append(backBufHeader, backBuf...)
		conn.Write(backBuf)
	}
}

func (s *Server) HandlerBuf(buf []byte) (backBuf []byte, err error) {
	var req Request
	if err := json.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	fmt.Println(req.Method)
	if serviceFunc := s.Service[req.Method]; serviceFunc != nil {
		fmt.Println(serviceFunc)
		return s.CallServiceMethod(serviceFunc, req.Params)

	} else {
		return nil, fmt.Errorf("method %s not found", req.Method)
	}

}

func (s *Server) CallServiceMethod(serviceFunc interface{}, params []interface{}) ([]byte, error) {
	funcValue := reflect.ValueOf(serviceFunc)
	fmt.Println(funcValue)
	// 指针解引用
	if funcValue.Kind() == reflect.Ptr {
		funcValue = funcValue.Elem()
	}
	fmt.Println(funcValue)

	funcParams := make([]reflect.Value, len(params))

	for i, param := range params {
		funcParams[i] = reflect.ValueOf(param)
	}

	results := funcValue.Call(funcParams)

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
