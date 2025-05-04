package main

import (
	"fmt"
	"net"
	"redisetgo/resp"
	"strings"
)

func main() {
	fmt.Println("REDISETGO")
	fmt.Println("I am listening on 6379")

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

		if client_input.Tipe != "Array" {
			fmt.Println("I want array type")
			break
		}

		if len(client_input.Array) == 0 {
			fmt.Println("I want more array length")
		}

		input_len := len(client_input.Array) / 2
		command := strings.ToUpper(client_input.Array[input_len].Bulk)
		arguments := client_input.Array[input_len+1:]

		output := resp.Resp_output_buffer(conn)

		command_output, ok := resp.Command_store[command]
		if !ok {
			fmt.Println("IDK this command: ", command)
			resp.Write_buffer(output, resp.Client_input{Tipe: "string", Str: ""})
			continue
		}

		result := command_output(arguments)
		resp.Write_buffer(output, result)
	}
}
