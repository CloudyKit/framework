package dynamic

import (
	"fmt"
	"reflect"
)

type Context struct {
	Property  string
	FieldData reflect.StructField
}

func VisitReflectValue[Type any](value reflect.Value, handler VisitorFunc[Type]) {

	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		VisitReflectValue(value.Elem(), handler)
	case reflect.Struct:
		reflectValueStructVisit(value, handler)
	}
}

func reflectValueStructVisit[Type any](valueOf reflect.Value, handler VisitorFunc[Type]) {
	if valueOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected value kind (Struct) got (%s)", valueOf.Kind().String()))
	}
	fields := getFields(valueOf.Type())
	numFields := len(fields)
	inTyp := reflect.TypeOf(handler).In(0)
	isPtr := inTyp.Kind() == reflect.Ptr

	for i := 0; i < numFields; i++ {
		fieldValue := valueOf.Field(i)

		if fieldValue.CanInterface() {

			fieldKind := fieldValue.Kind()

			if isPtr && fieldKind != reflect.Ptr {
				fieldValue = fieldValue.Addr()
			} else if !isPtr && (fieldKind == reflect.Ptr || fieldKind == reflect.Interface) {
				fieldValue = fieldValue.Elem()
			}
			if fieldValue.Type().AssignableTo(inTyp) {
				if val, ok := fieldValue.Interface().(Type); ok {
					newValueOf := reflect.ValueOf(handler(val, fields[i]))
					if fieldValue.CanSet() {
						fieldValue.Set(newValueOf)
					}
				}
			} else {
				VisitReflectValue(fieldValue, handler)
			}
		}
	}
}
