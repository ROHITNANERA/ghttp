package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func handleConnection(conn net.Conn, router *Router) {
	defer conn.Close()
	fmt.Printf("new connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		// Set a read deadline for the next request
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		req, err := parseRequest(reader)
		if err != nil {
			if err == io.EOF {
				break // Client closed connection
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break // Read timeout
			}
			writeResponse(conn, Response{StatusCode: 400, Body: strings.NewReader("Bad Request")})
			break
		}

		// Clear read deadline during request processing
		conn.SetReadDeadline(time.Time{})
		
		req.Context = context.Background()

		handler, found := router.Route(&req)
		if !found {
			writeResponse(conn, Response{StatusCode: 404, Body: strings.NewReader("Not Found")})
		} else {
			resp := handler(req)
			writeResponse(conn, resp)
		}

		// Check for connection close
		if strings.ToLower(req.Headers["Connection"]) == "close" {
			break
		}
	}
}

func parseRequest(reader *bufio.Reader) (Request, error) {
	req := Request{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
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
	
	// Parse path and query params
	u, err := url.Parse(parts[1])
	if err != nil {
		req.Path = parts[1]
	} else {
		req.Path = u.Path
		for k, v := range u.Query() {
			if len(v) > 0 {
				req.QueryParams[k] = v[0]
			}
		}
	}

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
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	defer conn.SetWriteDeadline(time.Time{})

	statusTxt := statusText(res.StatusCode)
	headers := res.Headers
	if headers == nil {
		headers = make(map[string]string)
	}

	if _, ok := headers["Content-type"]; !ok {
		headers["Content-type"] = "text/plain; charset=utf-8"
	}

	var hasContentLength bool
	if _, ok := headers["Content-Length"]; ok {
		hasContentLength = true
	}

	if !hasContentLength && res.Body != nil {
		headers["Transfer-Encoding"] = "chunked"
	} else if res.Body == nil {
		headers["Content-Length"] = "0"
	}

	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", res.StatusCode, statusTxt)

	for k, v := range headers {
		fmt.Fprintf(conn, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(conn, "\r\n")

	if res.Body != nil {
		if !hasContentLength {
			// Write chunked
			buf := make([]byte, 4096)
			for {
				n, err := res.Body.Read(buf)
				if n > 0 {
					fmt.Fprintf(conn, "%x\r\n", n)
					conn.Write(buf[:n])
					fmt.Fprintf(conn, "\r\n")
				}
				if err == io.EOF {
					fmt.Fprintf(conn, "0\r\n\r\n")
					break
				}
				if err != nil {
					break
				}
			}
		} else {
			io.Copy(conn, res.Body)
		}
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
	case 500:
		return "Internal Server Error"
	default:
		return "OK"
	}
}
