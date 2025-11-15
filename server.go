package main

import (
	"log"
	"net"
)

func StartServer(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		// handle connections conccurrently
		go handleConnection(conn)
	}
}
