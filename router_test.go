//go:build bench
// +build bench

// Package benchmarks contains performance tests for the radix-trie HTTP router.
package zepto

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// dummyHandler is a no-op HTTP handler used for benchmarking.
func dummyHandler(c *Context) error { return nil }

// benchGET sets up a router with the given number of GET routes and benchmarks a request to the last route.
func benchGET(b *testing.B, routeCount int) {
	r := NewDefault()
	for i := 0; i < routeCount; i++ {
		path := fmt.Sprintf("/route%d", i)
		r.GET(path, dummyHandler)
	}
	// Test the last route for worst-case lookup
	testPath := fmt.Sprintf("/route%d", routeCount-1)
	req := httptest.NewRequest(http.MethodGet, testPath, nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterGET runs sub-benchmarks with 10, 100, and 1000 routes.
func BenchmarkRouterGET(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%dRoutes", count), func(b *testing.B) {
			benchGET(b, count)
		})
	}
}

// benchGET sets up a router with the given number of GET routes and benchmarks a request to the last route.
func benchGETdynamic(b *testing.B, routeCount int) {
	r := NewDefault()
	for i := 0; i < routeCount; i++ {
		path := fmt.Sprintf("/route%d/:id", i)
		r.GET(path, dummyHandler)
	}
	// Test the last route for worst-case lookup
	testPath := fmt.Sprintf("/route%d/%d", routeCount-1, (routeCount-1)*2)
	req := httptest.NewRequest(http.MethodGet, testPath, nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}

// BenchmarkRouterGET runs sub-benchmarks with 10, 100, and 1000 routes.
func BenchmarkRouterGETDynamic(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%dRoutes", count), func(b *testing.B) {
			benchGETdynamic(b, count)
		})
	}
}
