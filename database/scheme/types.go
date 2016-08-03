package scheme

import (
	"fmt"
	"reflect"
	"strconv"
)

type (
	Int   struct{}
	Uint  struct{}
	Float struct{}
	Bool  struct{}

	Date     struct{}
	DateTime struct{}
	Time     struct{}
	Duration struct{}

	String struct {
		MaxLength int
	}
)

func (typ Int) Value(v reflect.Value) (reflect.Value, error) {
	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		return v, nil
	} else if kind >= reflect.Uint && kind <= reflect.Uintptr {
		return reflect.ValueOf(int64(v.Uint())), nil
	} else if kind == reflect.String {
		nint, err := strconv.ParseInt(v.String(), 10, 64)
		return reflect.ValueOf(nint), err
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		return reflect.ValueOf(int64(v.Float())), nil
	} else if kind == reflect.Bool {
		if v.Bool() {
			return reflect.ValueOf(int64(1)), nil
		}
		return reflect.ValueOf(int64(0)), nil
	}
	return v, fmt.Errorf("Scheme type: value of type %s can't be converted to int", v.Type())
}

func (typ Uint) Value(v reflect.Value) (reflect.Value, error) {
	kind := v.Kind()
	if kind >= reflect.Uint && kind <= reflect.Uintptr {
		return v, nil
	} else if kind >= reflect.Int && kind <= reflect.Int64 {
		_int := v.Int()
		v = reflect.ValueOf(uint64(_int))
		if _int < 0 {
			return v, fmt.Errorf("Scheme type: can't convert negative value into uint")
		}
		return v, nil
	} else if kind == reflect.String {
		nint, err := strconv.ParseUint(v.String(), 10, 64)
		return reflect.ValueOf(nint), err
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		_float := v.Float()
		v = reflect.ValueOf(uint64(_float))
		if _float < 0 {
			return v, fmt.Errorf("Scheme type: can't convert negative value into uint")
		}
		return v, nil
	} else if kind == reflect.Bool {
		if v.Bool() {
			return reflect.ValueOf(uint64(1)), nil
		}
		return reflect.ValueOf(uint64(0)), nil
	}
	return v, fmt.Errorf("Scheme type: value of type %s can't be converted to uint", v.Type())
}

func (typ Float) Value(v reflect.Value) (reflect.Value, error) {
	kind := v.Kind()
	if kind == reflect.Float32 || kind == reflect.Float64 {
		return v, nil
	} else if kind >= reflect.Uint && kind <= reflect.Uintptr {
		return reflect.ValueOf(float64(v.Uint())), nil
	} else if kind >= reflect.Int && kind <= reflect.Int64 {
		return reflect.ValueOf(float64(v.Int())), nil
	} else if kind == reflect.String {
		nint, err := strconv.ParseFloat(v.String(), 64)
		return reflect.ValueOf(nint), err
	} else if kind == reflect.Bool {
		if v.Bool() {
			return reflect.ValueOf(float64(1)), nil
		}
		return reflect.ValueOf(float64(0)), nil
	}
	return v, fmt.Errorf("Scheme type: value of type %s can't be converted to float", v.Type())
}

func (typ Bool) Value(v reflect.Value) (reflect.Value, error) {
	kind := v.Kind()
	if kind == reflect.Bool {
		return v, nil
	} else if kind >= reflect.Uint && kind <= reflect.Uintptr {
		return reflect.ValueOf(v.Uint() > 0), nil
	} else if kind >= reflect.Int && kind <= reflect.Int64 {
		return reflect.ValueOf(v.Int() > 0), nil
	} else if kind == reflect.String {
		nint, err := strconv.ParseBool(v.String())
		return reflect.ValueOf(nint), err
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		return reflect.ValueOf(v.Float() > 0), nil
	}
	return v, fmt.Errorf("Scheme type: value of type %s can't be converted to bool", v.Type())
}

func (typ String) Value(v reflect.Value) (reflect.Value, error) {
	if v.Kind() != reflect.String {
		v = reflect.ValueOf(fmt.Sprint(v.Interface()))
	}
	return v, nil
}
