package stub

import (
	"github.com/JIAKUNHUANG/krpc/client"
)

type DoubleRequest struct {
	Num float64 `json:"num"`
}

type DoubleResponse struct {
	Num float64 `json:"num"`
}

type Proxy struct {
	client *client.Client
}

func NewProxy() *Proxy {
	p := &Proxy{}
	p.client = client.NewClient()
	return p
}

func (p *Proxy) RegisterProxy(addr string) error {
	err := p.client.RegisterClient(addr)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Double(clientReq *DoubleRequest) (*DoubleResponse, error) {
	req := client.Request{
		Method: "Double",
		Params: clientReq,
	}

	rsp, err := p.client.Call(req)
	if err != nil {
		return nil, err
	}

	clientRsp := (rsp.Result).(*DoubleResponse)
	return clientRsp, nil

}
