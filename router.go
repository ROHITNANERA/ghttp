package main

// holds all routes
type Router struct {
	routes      map[string]Handler
	middlewares []Middleware
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
	wrapped := handler
	// wrap the handler inside middlewares
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		wrapped = r.middlewares[i](wrapped)
	}
	r.routes[key] = wrapped
}

// route request to the handler
func (r *Router) Route(req Request) (Handler, bool) {
	key := req.Method + " " + req.Path
	handler, exists := r.routes[key]
	return handler, exists
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}
