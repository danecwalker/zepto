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
