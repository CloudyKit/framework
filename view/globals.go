package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/context"
)

type provider interface {
	Provide(c *context.Context) interface{}
}

func init() {
	var defaultGlobal = Globals{}
	app.Default.Context.Map(defaultGlobal)
}

type valueProvider struct {
	v interface{}
}

func (v valueProvider) Provide(c *context.Context) interface{} {
	return v.v
}

type contextProvider struct {
	v interface{}
}

func (v contextProvider) Provide(c *context.Context) interface{} {
	return c.Get(v.v)
}

type Globals map[string]provider

func GlobalInjectName(ci *context.Context, name string, typ interface{}) error {
	return globalNameProvider(ci, name, contextProvider{typ})
}
func GlobalName(ci *context.Context, name string, v interface{}) error {
	return globalNameProvider(ci, name, valueProvider{v})
}

func globalNameProvider(ci *context.Context, name string, v provider) error {
	globals := ci.Get(Globals(nil)).(Globals)
	globals[name] = v
	return nil
}
