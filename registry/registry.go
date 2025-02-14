package registry

import (
	"net"
)

func Register() (err error) {
	//var ServiceList []dataTypes.Service
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	go func() {

	}()
	return nil
}
