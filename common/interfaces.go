package common

import (
	"fmt"
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

func GetURLer(cdi *cdi.Global) URLer {
	urler, _ := cdi.Val4Type(URLerType).(URLer)
	return urler
}

func GenURL(cdi *cdi.Global, resource string, v ...interface{}) string {
	urLer := GetURLer(cdi)
	if urLer == nil {
		return fmt.Sprintf(resource, v...)
	}
	return urLer.URL(resource, v...)
}
