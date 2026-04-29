package main

import (
	"context"
	"io"
)

// type for incoming request
type Request struct {
	Method      string
	Path        string
	QueryParams map[string]string
	PathParams  map[string]string
	Headers     map[string]string
	Body        []byte
	Context     context.Context
}

// type for response to send back
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       io.Reader
}

// handler function type
type Handler func(Request) Response
