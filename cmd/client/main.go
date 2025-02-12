package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/JIAKUNHUANG/krpc/client"
	"github.com/JIAKUNHUANG/krpc/server"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
		return
	}
	conn, err := client.ClientRegister(addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	testMessage := server.Message{
		Number: 2,
		Type:   "test",
		Data:   "test",
	}
	testMessageBytes, _ := json.Marshal(testMessage)
	testMessageBytesLength := len(testMessageBytes)
	testMesssgeBytesHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(testMesssgeBytesHeader, uint32(testMessageBytesLength))
	testMessageBytes = append(testMesssgeBytesHeader, testMessageBytes...)
	conn.Write(testMessageBytes)

	for {

		testMessageBytesHeader := make([]byte, 4)
		conn.Read(testMessageBytesHeader)
		testMessageBytesLength := binary.BigEndian.Uint32(testMessageBytesHeader)
		fmt.Println(testMessageBytesLength)
		testMessageBuf := make([]byte, testMessageBytesLength)
		conn.Read(testMessageBuf)
		fmt.Println(testMessageBuf)
		testMessage:= server.Message{}
		json.Unmarshal(testMessageBuf, &testMessage)
		fmt.Println(testMessage.Data, testMessage.Type, testMessage.Number)
	}
}
