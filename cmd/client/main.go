package main

import (
	"log"

	"github.com/JIAKUNHUANG/krpc/test/stub"
)

func main() {
	CallDouble()
}

func CallDouble() {
	// 注册客户端
	p := stub.NewProxy()

	err := p.RegisterProxy()
	if err != nil {
		log.Fatal(err)
	}

	req1 := stub.Teacher{
		Name: "JIAKUNHUANG",
		Sex:  true,
		StudentData: stub.Student{
			Name: "JIAKUNHUANG",
			Sex:  true,
		},
	}
	log.Println("request1:", req1)
	rsp1, err := p.SexExchange(req1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response1:", rsp1)

	req2 := stub.NumRequest{
		Num: 1.0,
	}
	log.Println("request2:", req2)
	rsp2, err := p.Double(req2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response2:", rsp2)
}
