package request

import (
	"errors"
)

func NewRequestFlow(r *Context, handler Handler, filters []Filter) Flow {
	return Flow{context: r, handler: handler, filters: filters}
}

type Filter func(*Context, Flow)

type Flow struct {
	context *Context
	filters []Filter
	handler Handler
}

type funcHandler func(*Context)

func (fn funcHandler) Handle(c *Context) {
	fn(c)
}

type Handler interface {
	Handle(*Context)
}

func (c *Flow) SetHandler(handler Handler) {
	if c == nil {
		panic(errors.New("Setting nil handler"))
	}
	c.handler = handler
}

func (c *Flow) SetHandlerFunc(handler funcHandler) {
	c.SetHandler(handler)
}

func (c *Flow) Handler() Handler {
	return c.handler
}

func (flow Flow) Continue() {
	if len(flow.filters) > 0 {
		f := flow.filters[0]
		flow.filters = flow.filters[1:]
		f(flow.context, flow)
		return
	}
	flow.handler.Handle(flow.context)
}
