package main

import (
	"log"
	"net"

	"github.com/JIAKUNHUANG/krpc/server"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
		return
	}
	s := server.CreateServer()
	s.AddServer("Double", func(input float64) (output float64) {
		output = input * 2
		return
	})
	s.AddServer("Add", func(input float64) (output float64) {
		output = input +1
		return
	})
	err = s.RegisterServer(addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer s.Listener.Close()
	go s.Server()
	select {}
}
