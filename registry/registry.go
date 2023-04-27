package registry

import (
	"github.com/CloudyKit/framework/container"
	"reflect"
)

func Get[Type any](c *container.Registry) Type {
	var t Type
	c.Load(&t)
	return t
}

type ContainerAware interface {
	Container() *container.Registry
}

func Set[Type any](c ContainerAware, val Type) (err error) {
	c.Container().WithValues(val)
	return
}

func Provider[Type any](c *container.Registry, builder func(registry *container.Registry) Type) (err error) {
	c.WithTypeAndProviderFunc(reflect.TypeOf((*Type)(nil)).Elem(), func(c *container.Registry) interface{} {
		return builder(c)
	})
	return
}
