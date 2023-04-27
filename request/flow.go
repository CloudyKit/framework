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
// process some data put the data into the scope.Registry and invoke DispatchNext which will invoke the next
// handler, the last handler is responsible for the main logic of the request.
// calling DispatchNext in the last handler will panic.
type Handler interface {
	Handle(*Context)
}

// DispatchNext entry point
func DispatchNext(context *Context, name string, writer http.ResponseWriter, request *http.Request, parameter router.Parameter, registry *container.Registry, handlers []Handler) error {
	context.Name = name
	context.Response = writer
	context.Request = request
	context.Parameters = parameter
	context.Registry = registry
	context.handlers = handlers
	if context.Request.Body != nil {
		context.body = context.Request.Body
		context.Request.Body = context.GetBodyReader()
	}

	//maps the request context into the scoped variables
	registry.WithValues(context)

	return context.Next()
}
