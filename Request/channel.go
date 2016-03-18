package Request

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
