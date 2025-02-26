package stub

import (
	"encoding/json"

	"github.com/JIAKUNHUANG/krpc/client"
	"github.com/JIAKUNHUANG/krpc/server"
)

type Teacher struct {
	Name        string  `json:"name"`
	Sex         bool    `json:"sex"`
	StudentData Student `json:"studentData"`
}

type Student struct {
	Name string `json:"name"`
	Sex  bool   `json:"sex"`
}

type NumRequest struct {
	Num float64 `json:"num"`
}

type NumResponse struct {
	Num float64 `json:"num"`
}

type Proxy struct {
	client *client.Client
}

func RegisterTestService(s *server.Service) {
	s.ServiceName = "ServiceTest"
	s.AddMethod("SexExchange", SexExchangeFunc)
	s.AddMethod("Double", DoubleFunc)
	addr := server.GetAddr(s.Config.ServiceAddr.Ip, s.Config.ServiceAddr.Port)

	err := s.RegisterService(addr)
	if err != nil {
		panic(err)
	}
}

func NewProxy() *Proxy {
	p := &Proxy{}
	p.client = client.NewClient()
	return p
}

func (p *Proxy) RegisterProxy() error {
	serviceAddr, err := client.ConnectServiceFinding("127.0.0.1:8000", "ServiceTest")
	if err != nil {
		return err
	}

	err = p.client.ConnectService(serviceAddr)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) SexExchange(clientReq Teacher) (clientRsp Teacher, err error) {
	req := client.Request{
		Method: "SexExchange",
		Params: clientReq,
	}

	rsp, err := p.client.Call(req)
	if err != nil {
		return clientRsp, err
	}

	rspResultBuf, _ := json.Marshal(rsp.Result)
	err = json.Unmarshal(rspResultBuf, &clientRsp)
	if err != nil {
		return clientRsp, err
	}

	return clientRsp, nil

}

func (p *Proxy) Double(clientReq NumRequest) (clientRsp NumResponse, err error) {
	req := client.Request{
		Method: "Double",
		Params: clientReq,
	}

	rsp, err := p.client.Call(req)
	if err != nil {
		return clientRsp, err
	}

	rspResultBuf, _ := json.Marshal(rsp.Result)
	err = json.Unmarshal(rspResultBuf, &clientRsp)
	if err != nil {
		return clientRsp, err
	}

	return clientRsp, nil

}
