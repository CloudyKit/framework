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
	Registry struct {
		parent     *Registry
		references int64
		values     map[reflect.Type]interface{}
	}

	Disposer interface {
		Dispose()
	}

	Provider interface {
		Provide(c *Registry) interface{}
	}

	Initializer interface {
		Initialize(c *Registry, value reflect.Value)
	}

	ProviderFunc func(c *Registry) interface{}

	InitializerFunc func(c *Registry, value reflect.Value)
)

var (
	registryPool = sync.Pool{
		New: func() interface{} {
			cc := new(Registry)
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

func (r *Registry) Container() *Registry {
	return r
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
func New() (r *Registry) {
	r = registryPool.Get().(*Registry)
	r.references = 0
	return
}

func (r *Registry) Load(dst interface{}) {
	r.LoadValue(reflect.ValueOf(dst))
}

func (r *Registry) LoadValue(dst reflect.Value) {
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
		typ := dst.Type()

		if providedValue, isSet := r.resolveType2Value(typ, dst); providedValue != nil || isSet {
			if !isSet {
				dst.Set(reflect.ValueOf(providedValue))
			}
		} else if __type == typ {
			dst.Set(reflect.ValueOf(r))
		} else if _, ok := injectables[typ]; ok {
			r.InjectValue(dst)
		}
	}
	return
}

// Inject walks the target looking the for exported fields that types match injectable types in Registry
func (r *Registry) Inject(target interface{}) {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	r.InjectValue(value)
}

var __type = reflect.TypeOf((*Registry)(nil))

// InjectValue walks the struct value looking to injectable fields
func (r *Registry) InjectValue(value reflect.Value) {
	if value.Kind() != reflect.Struct {
		panic("Invalid value passed to inject, required kind is struct get " + value.Kind().String())
	}
	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldTyp := field.Type()

		if providedValue, wasSet := r.resolveType2Value(fieldTyp, field); providedValue != nil || wasSet {
			if !wasSet {
				field.Set(reflect.ValueOf(providedValue))
			}
		} else if __type == fieldTyp {
			field.Set(reflect.ValueOf(r))
		} else if _, ok := injectables[fieldTyp]; ok {
			r.InjectValue(field)
		}
	}
	return
}

func (r *Registry) MapProvider(typ reflect.Type, provider Provider) {
	r.values[typ] = provider
}

func (r *Registry) WithTypeAndProviderFunc(typ reflect.Type, provider ProviderFunc) {
	r.values[typ] = provider
}

func (r *Registry) MapInitializer(typ reflect.Type, initializer Initializer) {
	r.values[typ] = initializer
}

func (r *Registry) MapInitializerFunc(typ reflect.Type, initializer InitializerFunc) {
	r.values[typ] = initializer
}

// WithTypeAndValue sets a provider for the type of typ with value of val
func (r *Registry) WithTypeAndValue(typOf reflect.Type, value interface{}) {
	valOf := reflect.ValueOf(value)
	if typOf.Kind() == reflect.Ptr && typOf.Elem().Kind() == reflect.Interface {
		typOf = typOf.Elem()
		if valOf.Kind() == reflect.Ptr {
			valOf = valOf.Elem()
			value = valOf.Interface()
		}
	}
	r.values[typOf] = value
}

// WithValues puts the list of values into the current context
func (r *Registry) WithValues(values ...interface{}) {
	for i := 0; i < len(values); i++ {
		r.WithTypeAndValue(reflect.TypeOf(values[i]), values[i])
	}
}

// Fork creates a new context using current value's repository as provider for the new context
func (r *Registry) Fork() (child *Registry) {
	if atomic.LoadInt64(&r.references) > -1 {
		atomic.AddInt64(&r.references, 1)
		child = New()
		child.parent = r
	} else {
		panic(errors.New("invoking child in a context already recycled"))
	}
	return
}

// resolveType search's for value of type typ, walking the context tree from the current to the top parent looking for the value with type typ
func (r *Registry) resolveType(typ reflect.Type) (val interface{}) {
	for {
		val = r.values[typ]
		if val != nil || r.parent == nil {
			return
		}
		r = r.parent
	}
	return
}

// resolveType2Value returns a value for the specified type typ
func (r *Registry) resolveType2Value(typ reflect.Type, valOf reflect.Value) (val interface{}, ok bool) {
	val = r.resolveType(typ)
	switch provider := val.(type) {
	case ProviderFunc:
		val = provider(r)
	case InitializerFunc:
		provider(r, valOf)
		ok = true
	case Provider:
		val = provider.Provide(r)
	case Initializer:
		provider.Initialize(r, valOf)
		ok = true
	}
	return
}

// LoadType returns a value for the specified type typ
func (r *Registry) LoadType(typ reflect.Type) (val interface{}) {
	val = r.resolveType(typ)
	switch provider := val.(type) {
	case ProviderFunc:
		val = provider(r)
	case InitializerFunc:
		valOf := reflect.New(typ).Elem()
		provider(r, valOf)
		val = valOf.Interface()
	case Provider:
		val = provider.Provide(r)
	case Initializer:
		valOf := reflect.New(typ).Elem()
		provider.Initialize(r, valOf)
		val = valOf.Interface()
	}
	return
}

// Dispose call end when the request is not need any more, this will cause all finalizers to run,
// this will my cause parent scopes to end also, this will happen case the parent scopes already ended but
// is waiting all children to end
func (r *Registry) Dispose() int64 {

	// check if this is the last active reference
	atomic.AddInt64(&r.references, -1)

	if r.references == -1 {
		r.finalize()
	} else if r.references < -1 {
		panic(fmt.Errorf("Inválid reference counting expected value is -1 got %v", r.references))
	}
	return r.references
}

var err = errors.New("scope.Registry.EndForce: requested that at this point all references to this context are previous cleared")

// MustDispose works same as End, but require that all children to be terminated when called
func (r *Registry) MustDispose() {
	if r.Dispose() > -1 {
		panic(err)
	}
}

func (r *Registry) recycle() {
	if r.parent != nil {
		r.parent.Dispose()
		r.parent = nil
	}
	registryPool.Put(r)
}

// finalize walks all values in the current context and invokes finalizers
// decrease reference counter into the parent
// and recycle the private data
func (r *Registry) finalize() {
	// invokes parent Done method
	defer r.recycle()
	//runs recycle here
	for _typ, _val := range r.values {
		delete(r.values, _typ)
		if _finalizer, isFinalizer := _val.(Disposer); isFinalizer {
			_finalizer.Dispose()
		}
	}
}
