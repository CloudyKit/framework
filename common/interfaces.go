package common

import (
	"github.com/CloudyKit/framework/cdi"
	"reflect"
)

var (
	NamedType = reflect.TypeOf((*Named)(nil)).Elem()
	URLerType = reflect.TypeOf((*URLer)(nil)).Elem()
)

type Named interface {
	Name() string
}

type URLer interface {
	URL(resource string, v ...interface{}) string
}

func GetURLer(cdi *cdi.DI) URLer {
	return cdi.Val4Type(URLerType).(URLer)
}

func GenURL(cdi *cdi.DI, resource string, v ...interface{}) string {
	return GetURLer(cdi).URL(resource, v...)
}
