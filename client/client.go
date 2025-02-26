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
		Method: "ServiceFindingMethod",
		Params: FindingRequest{
			ReqType: "finding", 
			ServiceName: serviceName,
		},
	}
	var rsp Response
	var findingResponse FindingResponse

	reqBuf, _ := json.Marshal(req)
	reqBufLen := len(reqBuf)
	reqBufHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(reqBufHeader, uint32(reqBufLen))

	reqBuf = append(reqBufHeader, reqBuf...)
	conn.Write(reqBuf)

	rspBufHeader := make([]byte, 4)
	_,err=conn.Read(rspBufHeader)
	if err != nil {
		return "", err
	}
	rspBufLen := binary.BigEndian.Uint32(rspBufHeader)
	rspBuf := make([]byte, rspBufLen)
	conn.Read(rspBuf)
	json.Unmarshal(rspBuf, &rsp)

	findingResponseByte,_:=json.Marshal(rsp.Result)
	json.Unmarshal(findingResponseByte,&findingResponse)


	if findingResponse.Status != "ok" {
		return "",fmt.Errorf(findingResponse.ErrMsg)
	}
	return findingResponse.Addr, nil
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
