package change

import "reflect"

type Operation interface{}

type Set struct {
	Field string
	Value reflect.Value
}
