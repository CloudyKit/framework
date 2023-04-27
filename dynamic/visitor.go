package dynamic

import (
	"fmt"
	"reflect"
	"sync"
)

type VisitorFunc[Type any] func(value Type, field reflect.StructField) (ret Type)

var (
	structFieldsCache = map[reflect.Type][]reflect.StructField{}
	mutex             = sync.RWMutex{}
)

func getFields(typOf reflect.Type) []reflect.StructField {
	mutex.RLock()
	if fields, found := structFieldsCache[typOf]; found {
		mutex.RUnlock()
		return fields
	}
	mutex.RUnlock()
	mutex.Lock()
	fields, found := structFieldsCache[typOf]
	if !found {
		numField := typOf.NumField()
		for i := 0; i < numField; i++ {
			fields = append(fields, typOf.Field(i))
		}
	}
	mutex.Unlock()
	return fields
}

func findStructVal(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr && val.Type().Elem().Kind() == reflect.Struct {
		return val.Elem()
	}
	return val
}

func MapVisitor[Type any](v any, handler VisitorFunc[Type]) {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Map {
		panic(fmt.Errorf("expected value kind (Slice) got (%s)", valueOf.Kind().String()))
	}
	lenOf := valueOf.Len()
	mapKeys := valueOf.MapKeys()
	for i := 0; i < lenOf; i++ {
		//todo handle
		elem := findStructVal(mapKeys[i])
		if elem.Kind() == reflect.Struct && elem.CanInterface() {
			StructVisitor(elem.Interface(), handler)
		}

		index := valueOf.MapIndex(mapKeys[i])
		elem = findStructVal(index)
		if elem.Kind() == reflect.Struct && elem.CanInterface() {
			StructVisitor(elem.Interface(), handler)
		}
	}
}

func SliceVisitor[Type any](v any, handler VisitorFunc[Type]) {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Slice {
		panic(fmt.Errorf("expected value kind (Slice) got (%s)", valueOf.Kind().String()))
	}
	lenOf := valueOf.Len()
	for i := 0; i < lenOf; i++ {
		elem := findStructVal(valueOf.Index(i))
		if elem.Kind() == reflect.Struct && elem.CanInterface() {
			StructVisitor(elem.Interface(), handler)
		}
	}
}

func StructVisitor[Type any](v any, handler VisitorFunc[Type]) {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() == reflect.Ptr || valueOf.Kind() == reflect.Interface {
		valueOf = valueOf.Elem()
	}
	reflectValueStructVisit(valueOf, handler)
}
