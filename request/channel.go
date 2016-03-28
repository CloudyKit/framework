package request

type Filters struct {
	filters []func(FContext)
}

func (f *Filters) AddFilter(filters ...func(FContext)) {
	f.filters = append(f.filters, filters...)
}
func (f *Filters) MakeFilters(filters ...func(FContext)) []func(FContext) {
	newFilter := make([]func(FContext), 0, len(f.filters) + len(filters))
	newFilter = append(newFilter, f.filters...)
	newFilter = append(newFilter, filters...)
	return newFilter
}

func NewFContext(r *Context, handler Handler, filters []func(FContext)) FContext {
	return FContext{Context:r, Handler:handler, filters:filters}
}

type FContext struct {
	*Context
	filters []func(FContext)
	Handler Handler
}

type Handler interface {
	Handle(*Context)
}

func (c FContext) Next() {
	if len(c.filters) > 0 {
		f := c.filters[0]
		c.filters = c.filters[1:]
		f(c)
		return
	}
	c.Handler.Handle(c.Context)
}
