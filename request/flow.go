package request

import (
	"github.com/CloudyKit/framework/scope"
	"github.com/CloudyKit/router"
	"net/http"
)

// HandlerFunc func implementing Handler interface
type HandlerFunc func(*Context)

func (fn HandlerFunc) Handle(c *Context) {
	fn(c)
}

// Handler is responsible to handle the request or part of the request, ex: a middleware handler would
// process some data put the data into the scope.Variables and invoke Advance which will invoke the next
// handler, the last handler is responsible for the main logic of the request.
// calling Advance in the last handler will panic.
type Handler interface {
	Handle(*Context)
}

// Advance entry point
func Advance(c *Context, name string, w http.ResponseWriter, r *http.Request, p router.Parameter, v *scope.Variables, h []Handler) {
	c.Name = name
	c.Response = w
	c.Request = r
	c.Parameters = p
	c.Variables = v
	c.handlers = h

	//maps the request context into the scoped variables
	v.Map(c)

	c.Advance()
}
