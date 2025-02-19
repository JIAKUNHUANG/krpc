package main

import (
	"github.com/JIAKUNHUANG/krpc/server"
	"github.com/JIAKUNHUANG/krpc/test/stub"
)

func main() {
	// 创建服务端结构体
	s := server.CreateService()

	// 注册服务
	stub.RegisterTestService(s)
	defer s.Listener.Close()

	// 注册服务
	go s.Service()

	select {}
}
