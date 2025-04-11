package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("REDIGO")
	fmt.Println("Listening on port 6379")

	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("TCP Error: ", err)
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Connection Error: ", err)
	}

	defer conn.Close()

	for {
		input := make([]byte, 512)
		_, err := conn.Read(input)
		if err != nil {
			fmt.Println("Reading Error: ", err)
		}

		conn.Write([]byte("+OK\r\n"))
	}
}
