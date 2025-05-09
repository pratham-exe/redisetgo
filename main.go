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

	aof_method := resp.Create_aof("redis.aof")
	defer resp.Aof_close(aof_method)

	resp.Read_aof(aof_method, func(ci resp.Client_input) {
		command := strings.ToUpper(ci.Array[0].Bulk)
		arguments := ci.Array[1:]

		command_output, ok := resp.Command_store[command]
		if !ok {
			fmt.Println("IDK this command: ", command)
			return
		}

		command_output(arguments)
	})

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Connection Error: ", err)
		return
	}

	defer conn.Close()

	for {
		input := resp.Resp_input_buffer(conn)
		client_input, err := resp.Read_buffer(input)
		if err != nil {
			fmt.Println("Read buffer error train: ", err)
			return
		}

		if client_input.Tipe != "array" {
			fmt.Println("I want array type")
			break
		}

		if len(client_input.Array) == 0 {
			fmt.Println("I want more array length")
			break
		}

		command := strings.ToUpper(client_input.Array[0].Bulk)
		arguments := client_input.Array[1:]

		output := resp.Resp_output_buffer(conn)

		command_output, ok := resp.Command_store[command]
		if !ok {
			fmt.Println("IDK this command: ", command)
			resp.Write_buffer(output, resp.Client_input{Tipe: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			resp.Write_aof(aof_method, client_input)
		}

		result := command_output(arguments)
		resp.Write_buffer(output, result)
	}
}
