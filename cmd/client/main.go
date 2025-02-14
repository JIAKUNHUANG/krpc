package main

import (
	"fmt"
	"log"
	"net"

	"github.com/JIAKUNHUANG/krpc/client"
)

func main() {

	rsp, err := Double(2)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(rsp)


	rsp, err = Add(2)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(rsp)

}

func Double(input float64) (output float64, err error) {

	c := client.CreateClient()
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {

		return 0, err
	}
	err = c.RegisterClient(addr)
	if err != nil {
		return 0, err
	}

	req := client.Request{
		Method: "Double",
		Params: []interface{}{input},
	}

	rsp, err := c.Call(req)
	if err != nil {
		return 0, err
	}
	return rsp.Result[0].(float64), nil
}


func Add(input float64) (output float64, err error) {

	c := client.CreateClient()
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {

		return 0, err
	}
	err = c.RegisterClient(addr)
	if err != nil {
		return 0, err
	}

	req := client.Request{
		Method: "Add",
		Params: []interface{}{input},
	}

	rsp, err := c.Call(req)
	if err != nil {
		return 0, err
	}
	return rsp.Result[0].(float64), nil
}