package view

import (
	"fmt"
	"github.com/CloudyKit/framework/common"
	"reflect"
)

func defaultName(typ reflect.Type) string {
	var plural string
restart:
	switch typ.Kind() {
	case reflect.Ptr:
		typ = typ.Elem()
		goto restart
	case reflect.Array, reflect.Slice:
		typ.Elem()
		plural = "List"
		goto restart
	}

	return typ.Name() + plural
}

func AvailableKey(m *Manager, name string, _typ interface{}) error {

	typ := reflect.TypeOf(_typ)

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	if name == "" {
		if _typ, isNamed := _typ.(common.Named); isNamed {
			name = _typ.Name()
		} else {
			name = defaultName(typ)
		}
	}

	for i := 0; i < len(m.injectables); i++ {
		autosetvar := m.injectables[i]
		if autosetvar.name == name {
			if autosetvar.typ == typ {
				return fmt.Errorf("Variable %s is already setted with %s", name, typ)
			}
			panic(fmt.Errorf("Variable %s was already setted with a diferent typ %s - new typ %s", name, autosetvar.typ, typ))
		}
	}
	m.injectables = append(m.injectables, autoset{name: name, typ: typ})
	return nil
}

func Available(m *Manager, typ interface{}) error {
	return AvailableKey(m, "", typ)
}
