package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/cdi"
	"reflect"
)

type provider interface {
	Provide(c *cdi.DI) interface{}
}

func init() {
	var defaultGlobal = Globals{}
	app.Default.Global.Map(defaultGlobal)
}

type valueProvider struct {
	v interface{}
}

func (v valueProvider) Provide(c *cdi.DI) interface{} {
	return v.v
}

type contextProvider struct {
	typeof reflect.Type
}

func (v contextProvider) Provide(c *cdi.DI) interface{} {
	return c.Val4Type(v.typeof)
}

type Globals map[string]provider

func GlobalInjectName(ci *cdi.DI, name string, typ interface{}) error {
	typeof := reflect.TypeOf(typ)
	if typeof.Kind() == reflect.Ptr && typeof.Elem().Kind() == reflect.Interface {
		typeof = typeof.Elem()
	}
	return globalNameProvider(ci, name, contextProvider{typeof})
}
func GlobalName(ci *cdi.DI, name string, v interface{}) error {
	return globalNameProvider(ci, name, valueProvider{v})
}

func globalNameProvider(ci *cdi.DI, name string, v provider) error {
	globals := ci.Get(Globals(nil)).(Globals)
	globals[name] = v
	return nil
}
