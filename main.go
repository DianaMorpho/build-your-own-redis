package main

import (
	"fmt"
	"net"
	"strings"
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
	writer := NewWriter(conn)

	defer conn.Close()
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		// the value should always be an array
		if value.typ != "array" {
			fmt.Println("Invalid request")
			writer.Write(Value{typ: "string", str: "Invalid request"})
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request")
			writer.Write(Value{typ: "string", str: "Invalid request"})
			continue
		}
		fmt.Println(value)

		// the first element of the array should be the command
		cmd := strings.ToUpper(value.array[0].bulk)
		// the rest of the elements should be the arguments
		args := value.array[1:]
		// see if we have a handler for the command
		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Unknown command")
			writer.Write(Value{typ: "string", str: "Unknown command"})
			continue
		}

		// call the handler
		res := handler(args)
		fmt.Println("Response:", res)
		writer.Write(res)
	}
}
