// main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	// Start the server
	router := NewRouter()
	// use middlewares
	router.Use(Logger)
	router.Use(Recovery)

	// setup paths (routes) to the actual endpoiunt logic
	setupRoutes(router)
	if err := Run(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func Run(addr string, router *Router) error {
	fmt.Println("ghttp listening on http://localhost:8080")
	return StartServer(":8080", router)
}

func setupRoutes(r *Router) {
	r.Handle("GET", "/", homeHandler)
	r.Handle("GET", "/health", healthHandler)
	r.Handle("POST", "/echo", echoHandler)
	r.Handle("GET", "/panic", panicTest)
}
