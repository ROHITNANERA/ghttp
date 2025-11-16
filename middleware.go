package main

import (
	"fmt"
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
	return func(r Request) Response {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[PANIC] Recovered: %v\n", r)
			}
		}()
		return next(r)
	}
}
