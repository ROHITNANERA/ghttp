package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func handleConnection(conn net.Conn, router *Router) {
	defer conn.Close()

	fmt.Printf("new connection from %s\n", conn.RemoteAddr())

	req, err := parseRequest(conn)
	if err != nil {
		// write response here
		writeResponse(
			conn,
			Response{
				StatusCode: 400,
				Body:       "Bad Request",
			},
		)

	}
	// route the request
	handler, found := router.Route(req)
	if !found {
		writeResponse(
			conn,
			Response{
				StatusCode: 404,
				Body:       "Bad Request",
			},
		)
		return
	}
	// execute the handler
	resp := handler(req)
	writeResponse(conn, resp)
}

func parseRequest(conn net.Conn) (Request, error) {
	reader := bufio.NewReader(conn)
	req := Request{
		Headers: make(map[string]string),
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		return req, err
	}
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return req, fmt.Errorf("invalid request line")
	}
	req.Method = parts[0]
	req.Path = parts[1]

	// headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break

		}
		line = strings.TrimSuffix(line, "\r\n")
		if line == "" {
			break
		}
		if colon := strings.Index(line, ":"); colon > 0 {
			key := strings.TrimSpace(line[:colon])
			value := strings.TrimSpace(line[colon+1:])
			req.Headers[key] = value
		}
	}

	// body if content-length
	if cl, ok := req.Headers["Content-Length"]; ok {
		if length, err := strconv.Atoi(cl); err == nil && length > 0 {
			body := make([]byte, length)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				return req, err
			}
			req.Body = body
		}

	}
	return req, nil
}

// write Response to connection
func writeResponse(conn net.Conn, res Response) {
	statusTxt := statusText(res.StatusCode)
	headers := res.Headers
	if headers == nil {
		headers = make(map[string]string)
	}

	headers["Content-type"] = "text/plain: charset=utf-8"

	body := res.Body
	headers["Content-Length"] = strconv.Itoa(len(body))

	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", res.StatusCode, statusTxt)

	for k, v := range headers {
		fmt.Fprintf(conn, "%s: %s\r\n", k, v)
	}
	// empty line
	fmt.Fprintf(conn, "\r\n")

	if body != "" {
		conn.Write([]byte(body))
	}
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	default:
		return "OK"
	}
}
