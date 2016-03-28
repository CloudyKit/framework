package request

type Filters struct {
	filters []func(Channel)
}

func (f *Filters) AddFilter(filters ...func(Channel)) {
	f.filters = append(f.filters, filters...)
}
func (f *Filters) MakeFilters(filters ...func(Channel)) []func(Channel) {
	newFilter := make([]func(Channel), 0, len(f.filters) + len(filters))
	newFilter = append(newFilter, f.filters...)
	newFilter = append(newFilter, filters...)
	return newFilter
}

type Channel struct {
	*Context
	Filters []func(Channel)
	Handler Handler
}

type Handler interface {
	Handle(*Context)
}

func (c Channel) Next() {
	if len(c.Filters) > 0 {
		f := c.Filters[0]
		c.Filters = c.Filters[1:]
		f(c)
		return
	}
	c.Handler.Handle(c.Context)
}
