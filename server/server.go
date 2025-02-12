package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Number int    `json:"number"`
	Type   string `json:"type"`
	Data   string `json:"data"`
}

func ServerRegister(addr *net.TCPAddr) (err error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		defer conn.Close()

		for {
			header := make([]byte, 4)
			_, err := receiveBag(conn, header)
			if err != nil {
				break
			}
			bufLen := binary.BigEndian.Uint32(header)
			buf := make([]byte, bufLen)
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			backBuf, err := HandlerBuf(buf[:n])
			if err != nil {
				break
			}
			conn.Write(backBuf)
		}
	}

}

func HandlerBuf(buf []byte) (backBuf []byte, err error) {
	fmt.Println(buf)
	msg := Message{}
	json.Unmarshal(buf, &msg)
	backmsg, err := ServerLogic(msg)
	if err != nil {
		return nil, err
	}
	backmsgBuf, _ := json.Marshal(backmsg)
	backBufLen := uint32(len(backmsgBuf))
	fmt.Println(backBufLen)
	backBufHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(backBufHeader, backBufLen)
	backBuf = append(backBufHeader, backmsgBuf...)
	fmt.Println(backmsgBuf)
	fmt.Println(backBuf)
	return backBuf, nil
}

func ServerLogic(input Message) (output Message, err error) {
	fmt.Println(input.Data, input.Type, input.Number)
	//input.Number+=1
	output = input
	return output, nil
}

func receiveBag(conn net.Conn, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := conn.Read(buf[total:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}
