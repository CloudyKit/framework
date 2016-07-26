package app

import (
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/events"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
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
		Global *scope.Variables

		*ctlGen
		emitter
		filterManager
	}

	contextHandler struct {
		pool      *sync.Pool
		isPtr     bool
		funcValue reflect.Value
		zeroValue reflect.Value
	}

	Context interface {
		Mx(*Mapper)
	}
)

func (app *App) BindContext(contexts ...Context) {
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
		newDi := app.Variables.Inherit()

		// creates a new cascade url generator
		myGen := new(ctlGen)
		// injects parent url generator
		newDi.Inject(myGen)
		myGen.urlGen = app.urlGen
		myGen.id = name + "."

		newDi.MapType(common.URLerType, myGen)

		emitter := app.emitter.(*events.Emitter)
		newDi.MapType(events.EmitterType, func(c *scope.Variables) interface{} {
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

func (handler *contextHandler) Handle(c *request.Context) {
	ii := handler.pool.Get()

	// get's or allocates a new context
	ctx := reflect.ValueOf(ii)
	c.Variables.InjectInStructValue(ctx.Elem())

	var arguments = [1]reflect.Value{ctx}
	if handler.isPtr == false {
		arguments[0] = arguments[0].Elem()
	}

	handler.funcValue.Call(arguments[0:])

	ctx.Elem().Set(handler.zeroValue)
	handler.pool.Put(ii)
}

var acRegex = regexp.MustCompile("/[:*][^/]+")

func (mx *Mapper) BindAction(method, path, action string, filters ...request.Filter) {
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

	mx.app.AddHandlerContextName(mx.Global, mx.name, method, mx.Prefix+path, &contextHandler{
		pool:      mx.pool,
		isPtr:     isPtr,
		zeroValue: mx.zeroValue,
		funcValue: methodByName.Func,
	}, mx.reslice(filters...)...)
}
