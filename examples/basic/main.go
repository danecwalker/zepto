package main

import (
	"fmt"
	"net/http"

	"github.com/danecwalker/zepto"
)

func main() {
	r := zepto.NewDefault()

	r.GET("/hello/:id", func(c *zepto.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	fmt.Println(r)

	http.ListenAndServe(":3000", r)
}
