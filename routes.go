package main

import (
	"fmt"
	"strings"
)

func homeHandler(req Request) Response {
	return Response{
		StatusCode: 200,
		Body:       strings.NewReader("Welcome to GHTTP -- HTTP server from raw TCP in Go!"),
	}
}

func healthHandler(req Request) Response {
	return Response{
		StatusCode: 200,
		Body:       strings.NewReader("OK\n"),
	}
}
func panicTest(req Request) Response {
	panic("Boom")
}

func echoHandler(req Request) Response {
	body := "Echo: "
	if len(req.Body) > 0 {
		body += string(req.Body)
	} else {
		body += "(No payload provided)"
	}
	body += "\n"
	return Response{
		StatusCode: 200,
		Body:       strings.NewReader(body),
	}
}

func userProfileHandler(req Request) Response {
	id := req.PathParams["id"]
	body := fmt.Sprintf("User Profile for ID: %s\n", id)
	return Response{
		StatusCode: 200,
		Body:       strings.NewReader(body),
	}
}
