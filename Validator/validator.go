package Validator

import (
	"net/http"
	"net/url"
	"reflect"
)

type Tester func(c *Context)
type Provider func(i string) reflect.Value

type Error struct {
	Field, Description string
}

type Result []Error

func (err Result) Accepted() bool {
	return len(err) == 0
}

func (err Result) Rejected() bool {
	return len(err) > 0
}

type Context struct {
	Name     string
	Value    reflect.Value
	target   reflect.Value
	provider Provider
	errors   Result
}

func (cc *Context) Field(name string) reflect.Value {
	if cc.provider != nil {
		return cc.provider(name)
	}
	return cc.target.FieldByName(name)
}

func (cc *Context) Done() Result {
	cc.errors = nil
	return cc.errors
}

func (cc *Context) Err(msg string) {
	cc.errors = append(cc.errors, Error{Field: cc.Name, Description: msg})
}

func (cc *Context) At(fieldName string, vs ...Tester) *Context {
	numValidators := len(vs)
	cc.Value = cc.Field(fieldName)
	cc.Name = fieldName
	for i := 0; i < numValidators; i++ {
		vs[i](cc)
	}
	return cc
}

func New(target interface{}) *Context {
	switch target := target.(type) {
	case *http.Request:
		return &Context{provider: func(key string) reflect.Value {
			return reflect.ValueOf(target.FormValue(key))
		}}
	case url.Values:
		return &Context{provider: func(key string) reflect.Value {
			return reflect.ValueOf(target.Get(key))
		}}
	case Provider:
		return &Context{provider: target}
	default:
		return &Context{target: reflect.Indirect(reflect.ValueOf(target))}
	}
}

type At func(fieldName string, vs ...Tester) *Context

func Run(target interface{}, aa func(At)) Result {
	cc := New(target)
	aa(cc.At)
	return cc.errors
}
