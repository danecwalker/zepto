package zepto

import (
	"net/http"
)

type IRouter interface {
	Handle(method, path string, handler HandlerFunc)
	Find(method, path string) (HandlerFunc, []Param, bool)
}

type Router struct {
	router IRouter
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

// Group creates a new route group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	return &Group{
		Router: r,
		Prefix: prefix,
	}
}

// CleanPath cleans the path by removing trailing slashes.
func cleanPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}

// Group Methods
func (g *Group) GET(path string, handler HandlerFunc) {
	path = cleanPath(g.Prefix + path)
	g.Router.GET(path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	path = cleanPath(g.Prefix + path)
	g.Router.POST(path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc) {
	path = cleanPath(g.Prefix + path)
	g.Router.PUT(path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc) {
	path = cleanPath(g.Prefix + path)
	g.Router.DELETE(path, handler)
}

// GET registers a handler for GET requests.
func (r *Router) GET(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodGet, path, handler)
}

// POST registers a handler for POST requests.
func (r *Router) POST(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodPost, path, handler)
}

// PUT registers a handler for PUT requests.
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodPut, path, handler)
}

// DELETE registers a handler for DELETE requests.
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.router.Handle(http.MethodDelete, path, handler)
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
