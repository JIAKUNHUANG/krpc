package main

import (
	"github.com/JIAKUNHUANG/krpc/server"
	"github.com/JIAKUNHUANG/krpc/test/stub"
)

func main() {
	// 创建服务端结构体
	s := server.CreateService()

	s.GetConfig("./config.json")
	// 注册服务
	stub.RegisterTestService(s)
	defer s.Listener.Close()

	if s.Config.ServiceFindingAdrr.Execute {
		err := s.ServiceFinding()
		if err != nil {
			panic(err)
		}
		defer s.ServiceFindingConn.Close()
	}

	// 注册服务
	go s.Service()

	select {}
}
