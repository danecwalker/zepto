package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danecwalker/zepto"
)

// Logger middleware
func Logger() zepto.MiddlewareFunc {
	return func(next zepto.HandlerFunc) zepto.HandlerFunc {
		return func(c *zepto.Context) error {
			start := time.Now()
			err := next(c)
			fmt.Printf("[%s] %s %s %v\n", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, time.Since(start))
			return err
		}
	}
}

// Auth middleware
func Auth() zepto.MiddlewareFunc {
	return func(next zepto.HandlerFunc) zepto.HandlerFunc {
		return func(c *zepto.Context) error {
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			return next(c)
		}
	}
}

func main() {
	r := zepto.NewDefault()

	// Add global middleware
	r.Use(Logger())

	// Public route
	r.GET("/hello/:id", func(c *zepto.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, " + c.Param("id") + "!"})
	})

	// Protected route group
	protected := r.Group("/api")
	protected.Use(Auth())
	protected.GET("/users", func(c *zepto.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Protected route accessed"})
	})

	http.ListenAndServe(":3000", r)
}
