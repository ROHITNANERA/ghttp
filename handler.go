package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("new connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("read error: ", err)
	}

	parts := strings.Fields(line)

	if len(parts) < 3 {
		return
	}
	method, path, _ := parts[0], parts[1], parts[2]
	fmt.Printf("REQUEST: %s -> %s\n", method, path)

	// HTTP Response
	response := "HTTP/1.1 200 OK\r\n"          // Status line
	response += "Content-Type: text/plain\r\n" // Header 1
	response += "Content-Length: 13\r\n"       // Header 2
	response += "\r\n"                         // ← **EMPTY LINE** (CRITICAL)
	response += "Hello, World!"                // Body

	conn.Write([]byte(response))

}
