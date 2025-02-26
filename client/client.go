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

func ConnectServiceFinding(serviceFindingAddr string, serviceName string) (serviceAddr string, err error) {
	addr, err := net.ResolveTCPAddr("tcp", serviceFindingAddr)
	if err != nil {
		return "", err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return "", err
	}

	req := Request{
		Method: "searchService",
		Params: serviceName,
	}
	var rsp Response

	reqBuf, _ := json.Marshal(req)
	reqBufLen := len(reqBuf)
	reqBufHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(reqBufHeader, uint32(reqBufLen))

	reqBuf = append(reqBufHeader, reqBuf...)
	conn.Write(reqBuf)

	rspBufHeader := make([]byte, 4)
	conn.Read(rspBufHeader)
	rspBufLen := binary.BigEndian.Uint32(rspBufHeader)
	rspBuf := make([]byte, rspBufLen)
	conn.Read(rspBuf)
	json.Unmarshal(rspBuf, &rsp)

	if rsp.Result == nil {
		rsp.Error = "no service found"
	}
	if rsp.Error != "" {
		return "", fmt.Errorf(rsp.Error)
	}
	return rsp.Result.(string), nil
}

func (c *Client) ConnectService(targetAddr string) error {

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
	json.Unmarshal(rspBuf, &rsp)

	return rsp, nil

}
