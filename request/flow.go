// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package request

import (
	"github.com/CloudyKit/framework/container"
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
func Advance(c *Context, name string, w http.ResponseWriter, r *http.Request, p router.Parameter, v *container.IoC, h []Handler) {
	c.Name = name
	c.Response = w
	c.Request = r
	c.Parameters = p
	c.IoC = v
	c.handlers = h

	//maps the request context into the scoped variables
	v.Map(c)

	c.Advance()
}
