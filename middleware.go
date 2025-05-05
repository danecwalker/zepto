package zepto

// MiddlewareFunc defines the function signature for middleware
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// Chain creates a middleware chain
func Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
