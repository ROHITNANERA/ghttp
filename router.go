package main

// holds all routes
type Router struct {
	routes map[string]Handler
}

// create new router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

// add route
func (r *Router) Handle(method, path string, handler Handler) {
	key := method + " " + path
	r.routes[key] = handler
}

// route request to the handler
func (r *Router) Route(req Request) (Handler, bool) {
	key := req.Method + " " + req.Path
	handler, exists := r.routes[key]
	return handler, exists
}
