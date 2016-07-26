package validation

import (
	"github.com/CloudyKit/router"
	"net/url"
	"reflect"
)

type Tester func(c *Validator)
type Provider func(i string) reflect.Value

func NewURLValueProvider(vl url.Values) Provider {
	return func(name string) reflect.Value {
		return reflect.ValueOf(vl.Get(name))
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

func (result Result) Good() bool {
	return len(result) == 0
}

func (result Result) Bad() bool {
	return len(result) > 0
}

func (result Result) Lookup(fieldName string) (err *Error, has bool) {
	for i := 0; i < len(result); i++ {
		err = &result[i]
		if err.Field == fieldName {
			return err, true
		}
	}
	return nil, false
}

func (result Result) Get(fieldName string) (err *Error) {
	err, _ = result.Lookup(fieldName)
	return
}

type Validator struct {
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

func (cc *Validator) Field(name string) reflect.Value {
	if cc.provider != nil {
		return cc.provider(name)
	}
	return cc.target.FieldByName(name)
}

func (cc *Validator) Done() Result {
	return cc.errors
}

func (cc *Validator) Err(msg string) {

	if cc.stoponerror {
		cc.stopped = true
	}
	cc.aterror = true

	cc.errors = append(cc.errors, Error{Field: cc.prefix + cc.Name, Description: msg})
}

func (cc *Validator) at(fieldName string, vs ...Tester) *Validator {
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

func newContext(target interface{}) *Validator {
	if target, isProvider := target.(Provider); isProvider {
		return &Validator{provider: target}
	}
	return &Validator{target: reflect.Indirect(reflect.ValueOf(target))}
}

type At func(fieldName string, vs ...Tester) *Validator

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
