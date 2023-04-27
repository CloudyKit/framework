package dynamic

import (
	"fmt"
	"reflect"
	"strings"
)

func PropertyGetDefault[Type any](val any, property string, def Type) (Type, bool) {
	value, found := PropertyGet[Type](val, property)
	if !found {
		value = def
	}
	return value, found
}
func PropertyGet[Type any](val any, property string) (value Type, found bool) {
	valOf := reflect.ValueOf(val)
	if valOf.Kind() == reflect.Ptr || valOf.Kind() == reflect.Interface {
		valOf = valOf.Elem()
	}
	if valOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected value kind (Struct) got (%s)", valOf.Kind().String()))
	}

	index := strings.Index(property, ".")
	if index == -1 {
		index = len(property)
	}
	for len(property) > 0 {
		curProperty := property[0:index]
		valOf = PropertyGetReflect(valOf, curProperty)
		if !valOf.IsValid() {
			found = false
			return
		}
		if curProperty == property {
			break
		}
		property = property[index+1:]
		index = strings.Index(property, ".")
		if index == -1 {
			index = len(property)
		}
	}
	found = true
	value = valOf.Interface().(Type)
	return
}

func PropertyGetReflect(valOf reflect.Value, property string) reflect.Value {
	if valOf.Kind() == reflect.Ptr || valOf.Kind() == reflect.Interface {
		valOf = valOf.Elem()
	}
	if valOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected value kind (Struct) got (%s)", valOf.Kind().String()))
	}
	fields := getFields(valOf.Type())
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		if field.Name == property {
			return valOf.Field(i)
		}
	}
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		if field.Anonymous {
			valOf = valOf.Field(i)
			kind := valOf.Kind()
			if kind == reflect.Struct {
				foundValue := PropertyGetReflect(valOf, property)
				if foundValue.IsValid() {
					return foundValue
				}
			} else if kind == reflect.Ptr || kind == reflect.Interface {
				valOf = valOf.Elem()
				if valOf.Kind() == reflect.Struct {
					foundValue := PropertyGetReflect(valOf, property)
					if foundValue.IsValid() {
						return foundValue
					}
				}
			}
		}
	}
	return reflect.Value{}
}
