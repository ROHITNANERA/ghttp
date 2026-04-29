package main

import "strings"

type node struct {
	path     string
	part     string
	children []*node
	handler  Handler
	isParam  bool
}

// holds all routes
type Router struct {
	roots       map[string]*node
	middlewares []Middleware
}

// create new router
func NewRouter() *Router {
	return &Router{
		roots: make(map[string]*node),
	}
}

func parsePath(path string) []string {
	var parts []string
	for _, p := range strings.Split(path, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func (n *node) insert(path string, parts []string, height int, handler Handler) {
	if len(parts) == height {
		n.path = path
		n.handler = handler
		return
	}
	part := parts[height]
	var child *node
	for _, c := range n.children {
		if c.part == part || c.isParam {
			child = c
			break
		}
	}
	if child == nil {
		child = &node{part: part, isParam: part[0] == ':'}
		n.children = append(n.children, child)
	}
	child.insert(path, parts, height+1, handler)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height {
		if n.handler != nil {
			return n
		}
		return nil
	}
	part := parts[height]
	for _, child := range n.children {
		if child.part == part || child.isParam {
			result := child.search(parts, height+1)
			if result != nil {
				return result
			}
		}
	}
	return nil
}

// add route
func (r *Router) Handle(method, path string, handler Handler) {
	wrapped := handler
	// wrap the handler inside middlewares
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		wrapped = r.middlewares[i](wrapped)
	}
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	parts := parsePath(path)
	r.roots[method].insert(path, parts, 0, wrapped)
}

// route request to the handler
func (r *Router) Route(req *Request) (Handler, bool) {
	root, ok := r.roots[req.Method]
	if !ok {
		return nil, false
	}
	parts := parsePath(req.Path)
	n := root.search(parts, 0)
	if n != nil {
		nParts := parsePath(n.path)
		req.PathParams = make(map[string]string)
		for i, p := range nParts {
			if p[0] == ':' {
				req.PathParams[p[1:]] = parts[i]
			}
		}
		return n.handler, true
	}
	return nil, false
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}
