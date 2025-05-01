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
	Params map[string]string
	// StatusCode is the HTTP status code.
	StatusCode int
	// Body is the request body.
	Body []byte
}

// NewContext creates a new context.
func NewContext(req *http.Request, w http.ResponseWriter, params []Param) *Context {
	return &Context{
		Request:  req,
		Response: w,
	}
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
