package client

import (
	"net"
)

func ClientRegister(addr *net.TCPAddr) (*net.TCPConn,error) {
	localAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return nil, err
	}
	defer listener.Close()
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return conn,nil
}
