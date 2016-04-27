package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/context"
	"reflect"
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
	typeof reflect.Type
}

func (v contextProvider) Provide(c *context.Context) interface{} {
	return c.Val4Type(v.typeof)
}

type Globals map[string]provider

func GlobalInjectName(ci *context.Context, name string, typ interface{}) error {
	typeof := reflect.TypeOf(typ)
	if typeof.Kind() == reflect.Ptr && typeof.Elem().Kind() == reflect.Interface {
		typeof = typeof.Elem()
	}
	return globalNameProvider(ci, name, contextProvider{typeof})
}
func GlobalName(ci *context.Context, name string, v interface{}) error {
	return globalNameProvider(ci, name, valueProvider{v})
}

func globalNameProvider(ci *context.Context, name string, v provider) error {
	globals := ci.Get(Globals(nil)).(Globals)
	globals[name] = v
	return nil
}
