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
	"github.com/CloudyKit/framework/events"
	"github.com/CloudyKit/framework/request"
	"reflect"
	"regexp"
	"sync"
)

type (
	Mapper struct {
		Prefix    string
		name      string
		typ       reflect.Type
		zeroValue reflect.Value

		pool   *sync.Pool
		app    *App
		Global *container.IoC

		*ctlGen
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

func (app *App) BindContext(contexts ...Controller) {
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

		name := structTyp.String()

		// creates a new di for this controller
		newDi := app.IoC.Fork()

		// creates a new cascade url generator
		myGen := new(ctlGen)
		// injects parent url generator
		newDi.Inject(myGen)
		myGen.urlGen = app.urlGen
		myGen.id = name + "."

		newDi.MapValue(common.URLGenType, myGen)

		emitter := app.emitter.(*events.Emitter)
		newDi.MapValue(events.EmitterType, func(c *container.IoC) interface{} {
			return emitter.Inherit()
		})

		controller.Mx(&Mapper{
			name:      name,
			app:       app,
			typ:       ptrTyp,
			ctlGen:    myGen,
			emitter:   app.emitter,
			Global:    newDi,
			zeroValue: zero,
			pool: &sync.Pool{
				New: func() interface{} {
					return reflect.New(structTyp).Interface()
				},
			},
		})
	}
}

func (handler *controllerHandler) Handle(c *request.Context) {
	ii := handler.pool.Get()

	// get's or allocates a new context
	ctx := reflect.ValueOf(ii)
	c.IoC.InjectValue(ctx.Elem())

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

	mx.app.urlGen[mx.typ.Elem().String()+"."+action] = acRegex.ReplaceAllStringFunc(mx.app.Prefix+mx.Prefix+path, func(st string) string {
		if st[1] == '*' {
			return "%v"
		}
		return "/%v"
	})

	mx.app.AddHandlerContextName(mx.Global, mx.name, method, mx.Prefix+path, &controllerHandler{
		pool:      mx.pool,
		isPtr:     isPtr,
		zeroValue: mx.zeroValue,
		funcValue: methodByName.Func,
	}, mx.reslice(filters...)...)
}
