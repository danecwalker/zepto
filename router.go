package zepto

import (
	"net/http"
)

type IRouter interface {
	Handle(method, path string, handler HandlerFunc)
	Find(method, path string) (HandlerFunc, []Param, bool)
}

type Router struct {
	router     IRouter
	middleware []MiddlewareFunc
}

type Group struct {
	*Router
	// Prefix is the prefix for the group.
	Prefix string
}

func New(r IRouter) *Router {
	return &Router{
		router: r,
	}
}

func NewDefault() *Router {
	return &Router{
		router: NewTrie(),
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware...)
}

// Group creates a new route group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	return &Group{
		Router: r,
		Prefix: prefix,
	}
}

// Group Methods
func (g *Group) GET(path string, handler HandlerFunc) {
	path = g.Prefix + path
	g.Router.GET(path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	path = g.Prefix + path
	g.Router.POST(path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc) {
	path = g.Prefix + path
	g.Router.PUT(path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc) {
	path = g.Prefix + path
	g.Router.DELETE(path, handler)
}

// GET registers a handler for GET requests.
func (r *Router) GET(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodGet, path, r.applyMiddleware(handler))
}

// POST registers a handler for POST requests.
func (r *Router) POST(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodPost, path, r.applyMiddleware(handler))
}

// PUT registers a handler for PUT requests.
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodPut, path, r.applyMiddleware(handler))
}

// DELETE registers a handler for DELETE requests.
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodDelete, path, r.applyMiddleware(handler))
}

// applyMiddleware applies all middleware to the handler
func (r *Router) applyMiddleware(handler HandlerFunc) HandlerFunc {
	if len(r.middleware) == 0 {
		return handler
	}
	return Chain(r.middleware...)(handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, params, ok := r.router.Find(req.Method, req.URL.Path)
	if !ok {
		http.NotFound(w, req)
		return
	}

	c := NewContext(req, w, params)
	if err := handler(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
