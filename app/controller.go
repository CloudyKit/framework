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

package app

import (
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/event"
	"github.com/CloudyKit/framework/request"
	"reflect"
	"regexp"
	"sync"
)

type (
	Mapper struct {
		Prefix    string
		Name      string
		typ       reflect.Type
		zeroValue reflect.Value

		pool     *sync.Pool
		app      *Kernel
		Registry *container.Registry

		*ControllerURLGen
		emitter
		filterHandlers
	}

	controllerHandler struct {
		pool      *sync.Pool
		isPtr     bool
		funcValue reflect.Value
		zeroValue reflect.Value
	}

	Controller interface {
		Mx(*Mapper)
	}
)

func (kernel *Kernel) AddControllers(contexts ...Controller) {
	for i := 0; i < len(contexts); i++ {
		controller := contexts[i]

		ptrTyp := reflect.TypeOf(controller)
		structTyp := ptrTyp
		zero := reflect.ValueOf(controller)
		if ptrTyp.Kind() == reflect.Ptr {
			structTyp = ptrTyp.Elem()
			zero = zero.Elem()
		} else {
			ptrTyp = reflect.PtrTo(ptrTyp)
		}

		// creates a new di for this controller
		registry := kernel.Registry.Fork()

		// creates a new cascade url generator
		myURLGen := new(ControllerURLGen)
		// injects parent url generator
		registry.Inject(myURLGen)
		myURLGen.urlGen = kernel.URLGen

		registry.WithTypeAndValue(common.URLGenType, myURLGen)

		emitter := kernel.emitter.(*event.Dispatcher)
		registry.WithTypeAndProviderFunc(event.EmitterType, func(c *container.Registry) interface{} {
			return emitter.Inherit()
		})

		mapper := &Mapper{
			Name:             structTyp.String(),
			app:              kernel,
			typ:              ptrTyp,
			ControllerURLGen: myURLGen,
			emitter:          kernel.emitter,
			Registry:         registry,
			zeroValue:        zero,
			pool: &sync.Pool{
				New: func() interface{} {
					return reflect.New(structTyp).Interface()
				},
			},
		}

		controller.Mx(mapper)
		myURLGen.id = mapper.Name + "."
	}
}

func (handler *controllerHandler) Handle(c *request.Context) {
	ii := handler.pool.Get()

	// gets or allocates a new context
	ctx := reflect.ValueOf(ii)
	c.Registry.InjectValue(ctx.Elem())

	var arguments = [1]reflect.Value{ctx}
	if handler.isPtr == false {
		arguments[0] = arguments[0].Elem()
	}

	handler.funcValue.Call(arguments[0:])

	ctx.Elem().Set(handler.zeroValue)
	handler.pool.Put(ii)
}

var acRegex = regexp.MustCompile("/[:*][^/]+")

func (mx *Mapper) BindAction(method, path, action string, filters ...request.Handler) {
	methodByName, isPtr := mx.typ.MethodByName(action)
	if !isPtr {
		methodByName, _ = mx.typ.Elem().MethodByName(action)
		if methodByName.Type == nil {
			panic("Inválid action " + action + " not found in controller " + mx.typ.String())
		}
	}

	mx.app.URLGen[mx.Name+"."+action] = acRegex.ReplaceAllStringFunc(mx.app.Prefix+mx.Prefix+path, func(st string) string {
		if st[1] == '*' {
			return "%v"
		}
		return "/%v"
	})

	mx.app.AddHandlerContextName(mx.Registry, mx.Name, method, mx.Prefix+path, &controllerHandler{
		pool:      mx.pool,
		isPtr:     isPtr,
		zeroValue: mx.zeroValue,
		funcValue: methodByName.Func,
	}, mx.reSlice(filters...)...)
}
