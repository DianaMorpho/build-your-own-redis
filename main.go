package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Hello, World!")
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
