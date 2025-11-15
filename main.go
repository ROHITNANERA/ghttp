// main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	// Start the server
	if err := Run(); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func Run() error {
	fmt.Println("ghttp listening on http://localhost:8080")
	return StartServer(":8080")
}
