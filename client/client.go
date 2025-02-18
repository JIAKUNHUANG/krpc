package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Method string
	Params interface{}
}

type Response struct {
	Result interface{}
	Error  string
}

type Client struct {
	Conn *net.TCPConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) RegisterClient(targetAddr string) error {
	localAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Client) Call(req Request) (rsp Response, err error) {
	reqBuf, _ := json.Marshal(req)
	fmt.Println(req.Params)
	reqBufLen := len(reqBuf)
	reqBufHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(reqBufHeader, uint32(reqBufLen))

	reqBuf = append(reqBufHeader, reqBuf...)

	c.Conn.Write(reqBuf)

	rspBufHeader := make([]byte, 4)
	c.Conn.Read(rspBufHeader)
	rspBufLen := binary.BigEndian.Uint32(rspBufHeader)
	rspBuf := make([]byte, rspBufLen)
	c.Conn.Read(rspBuf)
	fmt.Println(string(rspBuf))
	json.Unmarshal(rspBuf, &rsp)
	return rsp, nil

}
