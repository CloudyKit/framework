package Validator

import (
	"github.com/CloudyKit/framework/Router"
	"net/http"
	"net/url"
	"reflect"
)

type Tester func(c *Context)
type Provider func(i string) reflect.Value

func NewURLValueProvider(vl url.Values) Provider {
	return func(name string) reflect.Value {
		return reflect.ValueOf(vl.Get(name))
	}
}

func NewRequestValueProvider(vl *http.Request) Provider {
	return func(name string) reflect.Value {
		return reflect.ValueOf(vl.FormValue(name))
	}
}

func NewRouterValueProvider(vl Router.Values) Provider {
	return func(name string) reflect.Value {
		if vl.Index(name) == -1 {
			return reflect.Value{}
		}
		return reflect.ValueOf(vl.Get(name))
	}
}

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
	if target, isProvider := target.(Provider); isProvider {
		return &Context{provider: target}
	}
	return &Context{target: reflect.Indirect(reflect.ValueOf(target))}
}

type At func(fieldName string, vs ...Tester) *Context

func Run(target interface{}, aa func(At)) Result {
	cc := New(target)
	aa(cc.At)
	return cc.errors
}
