package zepto

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	// Request is the HTTP request.
	Request *http.Request
	// Response is the HTTP response.
	Response http.ResponseWriter
	// Params are the URL parameters.
	params []Param
	// StatusCode is the HTTP status code.
	StatusCode int
	// Body is the request body.
	Body []byte
}

// NewContext creates a new context.
func NewContext(req *http.Request, w http.ResponseWriter, params []Param) *Context {
	return &Context{
		Request:    req,
		Response:   w,
		params:     params,
		StatusCode: http.StatusOK,
		Body:       nil,
	}
}

func (c *Context) Param(key string) string {
	for _, param := range c.params {
		if param.Key == key {
			return param.Value
		}
	}
	return ""
}

// JSON writes a JSON response.
func (c *Context) JSON(status int, v interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	if err := json.NewEncoder(c.Response).Encode(v); err != nil {
		return err
	}

	return nil
}

// Text writes a text response.
func (c *Context) Text(status int, v string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	if _, err := c.Response.Write([]byte(v)); err != nil {
		return err
	}

	return nil
}
