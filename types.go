package main

// type for incoming request
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

// type for response to send back
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

// handler function type
type Handler func(Request) Response
