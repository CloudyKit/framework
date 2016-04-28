package validator

import (
	"github.com/CloudyKit/router"
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

func NewRouterValueProvider(vl router.Parameter) Provider {
	return func(name string) reflect.Value {
		if vl.IndexOf(name) == -1 {
			return reflect.Value{}
		}
		return reflect.ValueOf(vl.ByName(name))
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
	prefix      string
	Name        string
	Value       reflect.Value
	target      reflect.Value
	provider    Provider
	errors      Result
	aterror     bool
	stoponerror bool
	stopped     bool
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

	if cc.stoponerror {
		cc.stopped = true
	}
	cc.aterror = true

	cc.errors = append(cc.errors, Error{Field: cc.prefix + cc.Name, Description: msg})
}

func (cc *Context) at(fieldName string, vs ...Tester) *Context {
	if !cc.stopped {
		numValidators := len(vs)
		cc.Value = cc.Field(fieldName)
		cc.Name = fieldName
		cc.aterror = false
		for i := 0; i < numValidators; i++ {
			vs[i](cc)
			if cc.aterror {
				return cc
			}
		}
	}
	return cc
}

func newContext(target interface{}) *Context {
	if target, isProvider := target.(Provider); isProvider {
		return &Context{provider: target}
	}
	return &Context{target: reflect.Indirect(reflect.ValueOf(target))}
}

type At func(fieldName string, vs ...Tester) *Context

func Run(target interface{}, aa func(At)) Result {
	cc := newContext(target)
	aa(cc.at)
	return cc.errors
}

func RunStop(target interface{}, aa func(At)) Result {
	cc := newContext(target)
	cc.stoponerror = true
	aa(cc.at)
	return cc.errors
}
