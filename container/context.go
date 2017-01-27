// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	IoC struct {
		parent     *IoC
		references int64
		values     map[reflect.Type]interface{}
	}

	Disposer interface {
		Dispose()
	}

	Provider interface {
		Provide(c *IoC) interface{}
	}

	Initializer interface {
		Initialize(c *IoC, value reflect.Value)
	}

	ProviderFunc func(c *IoC) interface{}

	InitializerFunc func(c *IoC, value reflect.Value)
)

var (
	dipool = sync.Pool{
		New: func() interface{} {
			cc := new(IoC)
			cc.values = make(map[reflect.Type]interface{})
			return cc
		},
	}
	injectables = map[reflect.Type]struct{}{}
)

func TypeOfElem(i interface{}) reflect.Type {
	return reflect.TypeOf(i).Elem()
}

func TypeOf(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}

// Injectable marks the type as field to be injectable while searching for fields to inject
func Injectable(v ...interface{}) uint {
	for i := 0; i < len(v); i++ {
		typ := reflect.TypeOf(v[i])
	TRY:
		switch typ.Kind() {
		case reflect.Struct:
		case reflect.Ptr:
			typ = typ.Elem()
			goto TRY
		default:
			continue
		}

		injectables[typ] = struct{}{}
	}
	return 0
}

// New creates a new instance of context object
func New() (cc *IoC) {
	cc = dipool.Get().(*IoC)
	cc.references = 0
	return
}

func (c *IoC) Load(into interface{}) {
	c.LoadValue(reflect.ValueOf(into))
}

func (c *IoC) LoadValue(into reflect.Value) {
	if into.Kind() == reflect.Ptr {
		into = into.Elem()
		typ := into.Type()

		if provided_value, isSet := c.resolveType2Value(typ, into); provided_value != nil || isSet {
			if !isSet {
				into.Set(reflect.ValueOf(provided_value))
			}
		} else if __type == typ {
			into.Set(reflect.ValueOf(c))
		} else if _, ok := injectables[typ]; ok {
			c.InjectValue(into)
		}
	}
	return
}

// Inject walks the target looking the for exported fields that types match injectable types in Global
func (c *IoC) Inject(target interface{}) {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	c.InjectValue(value)
}

var __type = reflect.TypeOf((*IoC)(nil))

// InjectValue walks the struct value looking to injectable fields
func (c *IoC) InjectValue(value reflect.Value) {
	if value.Kind() != reflect.Struct {
		panic("Invalid value passed to inject, required kind is struct get " + value.Kind().String())
	}
	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldTyp := field.Type()

		if provided_value, wasSetted := c.resolveType2Value(fieldTyp, field); provided_value != nil || wasSetted {
			if !wasSetted {
				field.Set(reflect.ValueOf(provided_value))
			}
		} else if __type == fieldTyp {
			field.Set(reflect.ValueOf(c))
		} else if _, ok := injectables[fieldTyp]; ok {
			c.InjectValue(field)
		}
	}
	return
}

func (c *IoC) MapProvider(typ reflect.Type, provider Provider) {
	c.values[typ] = provider
}

func (c *IoC) MapProviderFunc(typ reflect.Type, provider ProviderFunc) {
	c.values[typ] = provider
}

func (c *IoC) MapInitializer(typ reflect.Type, initializer Initializer) {
	c.values[typ] = initializer
}

func (c *IoC) MapInitializerFunc(typ reflect.Type, initializer InitializerFunc) {
	c.values[typ] = initializer
}

// MapValue sets a provider for the type of typ with value of val
func (c *IoC) MapValue(typOf reflect.Type, val interface{}) {
	c.values[typOf] = val
}

// Map puts the list of values into the current context
func (c *IoC) Map(value ...interface{}) {
	for i := 0; i < len(value); i++ {
		vof := value[i]
		v := reflect.ValueOf(vof)
		c.values[v.Type()] = vof
	}
}

// Fork creates a new context using current values repository as provider for the new context
func (c *IoC) Fork() (child *IoC) {
	if atomic.LoadInt64(&c.references) > -1 {
		c.references = atomic.AddInt64(&c.references, 1)
		child = New()
		child.parent = c
	} else {
		panic(errors.New("Invoking Child in a context already recycled"))
	}
	return
}

// resolveType search's for value of type typ, walking the context tree from the current to the top parent looking for the value with type typ
func (_context *IoC) resolveType(typ reflect.Type) (val interface{}) {
	for {
		val = _context.values[typ]
		if val != nil || _context.parent == nil {
			return
		}
		_context = _context.parent
	}
	return
}

// resolveType2Value returns a value for the specified type typ
func (c *IoC) resolveType2Value(typ reflect.Type, valOf reflect.Value) (val interface{}, ok bool) {
	val = c.resolveType(typ)
	switch provider := val.(type) {
	case ProviderFunc:
		val = provider(c)
	case InitializerFunc:
		provider(c, valOf)
		ok = true
	case Provider:
		val = provider.Provide(c)
	case Initializer:
		provider.Initialize(c, valOf)
		ok = true
	}
	return
}

// LoadType returns a value for the specified type typ
func (c *IoC) LoadType(typ reflect.Type) (val interface{}) {
	val = c.resolveType(typ)
	switch provider := val.(type) {
	case ProviderFunc:
		val = provider(c)
	case InitializerFunc:
		valOf := reflect.New(typ).Elem()
		provider(c, valOf)
		val = valOf.Interface()
	case Provider:
		val = provider.Provide(c)
	case Initializer:
		valOf := reflect.New(typ).Elem()
		provider.Initialize(c, valOf)
		val = valOf.Interface()
	}
	return
}

// Dispose call end when the request is not need any more, this will cause all finalizers to run,
// this will my cause parent scopes to end also, this will happen case the parent scopes already ended but
// is waiting all children to end
func (c *IoC) Dispose() int64 {

	// check if this is the last active reference
	c.references = atomic.AddInt64(&c.references, -1)

	if c.references == -1 {
		c.finalize()
	} else if c.references < -1 {
		panic(fmt.Errorf("Inválid reference counting expected value is -1 got %v", c.references))
	}
	return c.references
}

var err = errors.New("scope.Variables.EndForce: requested that at this point all references to this context are previous cleared")

// MustDispose works same as End, but require that all children to be terminated when called
func (c *IoC) MustDispose() {
	if c.Dispose() > -1 {
		panic(err)
	}
}

func (c *IoC) recycle() {
	if c.parent != nil {
		c.parent.Dispose()
		c.parent = nil
	}
	dipool.Put(c)
}

// finalize walks all values in the current context and invokes finalizers
// decrease reference counter into the parent
// and recycle the private data
func (c *IoC) finalize() {
	// invokes parent Done method
	defer c.recycle()
	//runs recycle here
	for _typ, _val := range c.values {
		// not delete the keys
		delete(c.values, _typ)
		if _finalizer, isFinalizer := _val.(Disposer); isFinalizer {
			_finalizer.Dispose()
		}
	}
}
