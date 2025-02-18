package main

import (
	"log"

	"github.com/JIAKUNHUANG/krpc/test/stub"
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

func main() {
	CallDouble()
}

func CallDouble() {
	// 注册客户端
	p := stub.NewProxy()

	err := p.RegisterProxy("127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	req := &stub.DoubleRequest{
		Num: 1.0,
	}
	log.Println("request:", req.Num)

	rsp, err := p.Double(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response:", rsp.Num)
}
