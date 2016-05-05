package request

import (
	"errors"
)

type Filters struct {
	filters []func(*Context, Flow)
}

func (f *Filters) AddFilter(filters ...func(*Context, Flow)) {
	f.filters = append(f.filters, filters...)
}

func (f *Filters) MakeFilters(filters ...func(*Context, Flow)) []func(*Context, Flow) {
	newFilter := make([]func(*Context, Flow), 0, len(f.filters)+len(filters))
	newFilter = append(newFilter, f.filters...)
	newFilter = append(newFilter, filters...)
	return newFilter
}

func NewContextChain(r *Context, handler Handler, filters []func(*Context, Flow)) Flow {
	return Flow{context: r, handler: handler, filters: filters}
}

type Flow struct {
	context *Context
	filters []func(*Context, Flow)
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
