package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
	"reflect"
)

type provider interface {
	Provide(c *scope.Variables) interface{}
}

func init() {
	var defaultGlobal = Globals{}
	app.Default.Variables.Map(defaultGlobal)
	GlobalInjectName(app.Default.Variables, "link", common.URLerType)
}

type valueProvider struct {
	v interface{}
}

func (v valueProvider) Provide(c *scope.Variables) interface{} {
	return v.v
}

type contextProvider struct {
	typeof reflect.Type
}

func (v contextProvider) Provide(c *scope.Variables) interface{} {
	return c.GetByType(v.typeof)
}

type Globals map[string]provider

func GlobalInjectName(ci *scope.Variables, name string, typ reflect.Type) error {
	return globalNameProvider(ci, name, contextProvider{typ})
}
func GlobalName(ci *scope.Variables, name string, v interface{}) error {
	return globalNameProvider(ci, name, valueProvider{v})
}

func globalNameProvider(ci *scope.Variables, name string, v provider) error {
	globals := ci.GetByPtr(Globals(nil)).(Globals)
	globals[name] = v
	return nil
}

type setFilter_Item struct {
	Name string
	Val  interface{}
}
type setFilter struct {
	filters []setFilter_Item
}

func (s *setFilter) Set(name string, val interface{}) *setFilter {
	s.filters = append(s.filters, setFilter_Item{name, val})
	return s
}

func (s *setFilter) Build() request.Filter {
	return func(c *request.Context, f request.Flow) {
		v := GetRenderer(c.Variables)
		for i := 0; i < len(s.filters); i++ {
			v.With(s.filters[i].Name, s.filters[i].Val)
		}
		f.Continue()
	}
}

func NewSetterFilter(name string, val interface{}) *setFilter {
	return (&setFilter{}).Set(name, val)
}
