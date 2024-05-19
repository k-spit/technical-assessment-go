package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Client connected from %v\n", c.RemoteAddr())
	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Println(input.Text()) // Echo the received data
	}
	c.Close()
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("TCP server listening on port 8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn) // Handle each connection concurrently
	}
}
