package main

import (
	"fmt"
	"strings"
	"time"
)

type Middleware func(Handler) Handler

// logs every request
func Logger(next Handler) Handler {
	return func(r Request) Response {
		start := time.Now()
		res := next(r)
		duration := time.Since(start)
		fmt.Printf(
			"[LOG] %s %s %d (%v)\n",
			r.Method, r.Path, res.StatusCode, duration,
		)
		return res
	}
}

// Recovery to catch all the panics
func Recovery(next Handler) Handler {
	return func(r Request) (res Response) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] Recovered: %v\n", err)
				res = Response{
					StatusCode: 500,
					Body:       strings.NewReader("Internal Server Error\n"),
				}
			}
		}()
		return next(r)
	}
}
