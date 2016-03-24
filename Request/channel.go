package Request

type Filter func(Channel)
type Filters struct {
	filters []Filter
}

func (f *Filters) AddFilter(filters ...Filter) {
	f.filters = append(f.filters, filters...)
}
func (f *Filters) MakeFilters(filters ...Filter) []Filter {
	newFilter := make([]Filter, 0, len(f.filters)+len(filters))
	newFilter = append(newFilter, f.filters...)
	newFilter = append(newFilter, filters...)
	return newFilter
}

type Channel struct {
	*Context
	Filters []Filter
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
