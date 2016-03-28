package request

type Filters struct {
	filters []func(ContextChain)
}

func (f *Filters) AddFilter(filters ...func(ContextChain)) {
	f.filters = append(f.filters, filters...)
}
func (f *Filters) MakeFilters(filters ...func(ContextChain)) []func(ContextChain) {
	newFilter := make([]func(ContextChain), 0, len(f.filters)+len(filters))
	newFilter = append(newFilter, f.filters...)
	newFilter = append(newFilter, filters...)
	return newFilter
}

func NewContextChain(r *Context, handler Handler, filters []func(ContextChain)) ContextChain {
	return ContextChain{Context: r, Handler: handler, filters: filters}
}

type ContextChain struct {
	*Context
	filters []func(ContextChain)
	Handler Handler
}

type Handler interface {
	Handle(*Context)
}

func (c ContextChain) Next() {
	if len(c.filters) > 0 {
		f := c.filters[0]
		c.filters = c.filters[1:]
		f(c)
		return
	}
	c.Handler.Handle(c.Context)
}
