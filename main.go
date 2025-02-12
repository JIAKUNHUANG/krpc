package main

import (
	"fmt"
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
	server.ServerRegister(addr)
	fmt.Println("Server is running...")
}
