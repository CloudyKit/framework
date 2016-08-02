package scheme

import (
	"fmt"
	"reflect"
)

type Int struct{}

var (
	typeInt    = reflect.TypeOf(int64(0))
	typeFloat  = reflect.TypeOf(int64(0))
	typeBool   = reflect.TypeOf(int64(0))
	typeString = reflect.TypeOf(int64(0))
)

func (typ Int) Value(v reflect.Value) (reflect.Value, error) {
	if !v.Type().AssignableTo(typeInt) {
		if v.Type().ConvertibleTo(typeInt) {
			return v.Convert(typeInt), nil
		}
		return v, fmt.Errorf("value of type %s can't be converted to int", v.Type())
	}
	return v, nil
}

type String struct {
	MaxLength int
}

func (typ String) Value(v reflect.Value) (reflect.Value, error) {
	if !v.Type().AssignableTo(typeInt) {
		if v.Type().ConvertibleTo(typeInt) {
			return v.Convert(typeInt), nil
		}
		return v, fmt.Errorf("value of type %s can't be converted to int", v.Type())
	}
	return v, nil
}
