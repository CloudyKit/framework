package database

import (
	"github.com/CloudyKit/framework/validation"
	"reflect"
)

var modelTYPE = reflect.TypeOf((*Model)(nil)).Elem()

func validModelSlice(v reflect.Type) bool {
	kind := v.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		v = v.Elem()
		if v.Implements(modelTYPE) {
			return true
		}
		if reflect.PtrTo(v).Implements(modelTYPE) {
			return true
		}
	}
	return false
}

func canBeNil(r reflect.Kind) bool {
	return r > reflect.Array && r < reflect.String || r == reflect.UnsafePointer
}

func getRefModel(v reflect.Value) (reflect.Value, bool) {
	t := v.Type()
	ok := t.Implements(modelTYPE)
	if !ok {
		kind := t.Kind()
		if kind == reflect.Struct && reflect.PtrTo(t).Implements(modelTYPE) {
			if v.CanAddr() {
				v = v.Addr()
				ok = true
			}
		} else if kind == reflect.Ptr && t.Elem().Implements(modelTYPE) {
			v = v.Elem()
			ok = true
		}
	}
	return v, ok
}

func runFieldTesters(v *validation.Validator, fieldName string, val reflect.Value, testers ...validation.Tester) *validation.Validator {
	if v == nil {
		v = new(validation.Validator)
	}
	v.Name = fieldName
	v.Value = val
	for i := 0; i < len(testers); i++ {
		testers[i](v)
	}
	return v
}
