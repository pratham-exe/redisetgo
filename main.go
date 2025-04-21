package main

import (
	"fmt"
	"net"
	"redigo/resp"
)

func main() {
	fmt.Println("REDIGO")
	fmt.Println("Listening on port 6379")

	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("TCP Error: ", err)
		return
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Connection Error: ", err)
		return
	}

	defer conn.Close()

	for {
		input := resp.Resp_input_buffer(conn)
		client_input := resp.Read_buffer(input)

		fmt.Println(client_input)

		output := resp.Resp_output_buffer(conn)
		resp.Write_buffer(output, resp.Client_input{Tipe: "string", Str: "Hello"})
	}
}
