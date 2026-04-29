# GHTTP

GHTTP is an educational HTTP server built entirely from scratch in Go using raw TCP sockets (`net.Listen`). It was created for learning purposes to understand the inner workings of the HTTP protocol, connection management, and routing.

## 🚀 Features

- **Raw TCP Sockets**: No `net/http` used! Everything is parsed and written directly from/to TCP connections.
- **Radix Tree Routing**: Efficient routing with support for dynamic path parameters (e.g., `/users/:id`).
- **HTTP Keep-Alive**: Reuses TCP connections for sequential requests to maximize throughput.
- **Chunked Transfer Streaming**: Automatically uses `Transfer-Encoding: chunked` for true streaming when no `Content-Length` is provided.
- **Query Parameter Parsing**: Extracts URL query strings (e.g., `?q=search`) seamlessly.
- **Middleware Support**: Includes custom `Logger` and `Recovery` (Panic catching) middleware.
- **Security & Timeouts**: Implements socket-level read/write deadlines and provides `context.Context` for request cancellation.

## 🛠️ Project Structure

- `main.go` - Entry point. Sets up the router, registers middlewares and routes, and starts the server.
- `server.go` - Handles the underlying TCP listener and spawns goroutines for incoming connections.
- `handler.go` - The core HTTP parser and response writer. Handles Keep-Alive loops, chunked encoding, and HTTP formatting.
- `router.go` - A custom Radix Tree (Trie) implementation for performant, parameterized route matching.
- `middleware.go` - Request interceptors for logging and panic recovery (returning 500 status codes).
- `routes.go` - Example endpoint definitions.
- `types.go` - Core struct definitions (`Request`, `Response`, `Handler`).

## 🏃 How to Run

1. Clone the repository and navigate to the project directory:
   ```bash
   cd ghttp
   ```

2. Run the server:
   ```bash
   go run .
   ```

3. Test it out using `curl` or your browser:

   **Basic Request**
   ```bash
   curl -v http://localhost:8080/
   ```

   **Dynamic Path Parameters**
   ```bash
   curl -v http://localhost:8080/users/123
   ```

   **Query Parameters**
   ```bash
   curl -v "http://localhost:8080/users/123?q=search"
   ```

   **Panic Recovery Test**
   ```bash
   curl -v http://localhost:8080/panic
   ```

## 📚 Learning Goals

This project demonstrates how Go's powerful standard library features like `bufio.Reader`, `io.Reader`, and `net.Conn` can be used to construct a fully functional web server without relying on the heavily optimized `net/http` package.

---
*Created for learning and exploration.*
