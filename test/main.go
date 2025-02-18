package main

import (
	"fmt"
	"log"
	"net"

	"github.com/JIAKUNHUANG/krpc/server"
	"github.com/JIAKUNHUANG/krpc/test/stub"
)

func main() {
	// 定义服务端口
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
		return
	}

	//创建服务端结构体
	s := server.CreateService()

	// 添加方法
	s.AddMethod("Double", func(input stub.DoubleRequest) (output stub.DoubleResponse) {
		output.Num = input.Num * 2
		fmt.Println(input.Num, output.Num)
		return
	})

	// 监听服务端口
	err = s.RegisterService(addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer s.Listener.Close()

	// 注册服务
	go s.Service()

	select {}
}
