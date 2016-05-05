package app

import (
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/request"
	"reflect"
	"regexp"
	"sync"
)

type (
	Mapper struct {
		name      string
		typ       reflect.Type
		zeroValue reflect.Value

		pool    *sync.Pool
		app     *App
		Context *cdi.DI

		filterManager
	}

	contextHandler struct {
		pool      *sync.Pool
		isPtr     bool
		funcValue reflect.Value
		zeroValue reflect.Value
	}

	appContext interface {
		Mx(*Mapper)
	}
)

func (app *App) AddController(controllers ...appContext) {
	for i := 0; i < len(controllers); i++ {
		controller := controllers[i]

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
		newDi := app.Global.Child()
		// creates a new cascade url generator
		myGen := new(ctlGen)
		// injects parent url generator
		newDi.Inject(myGen)

		myGen.urlGen = app.urlGen
		myGen.id = name + "."

		newDi.MapType(common.URLerType, myGen)

		controller.Mx(&Mapper{
			name:          name,
			app:           app,
			typ:           ptrTyp,
			Context:       newDi,
			filterManager: filterManager{filters: app.reslice()},
			zeroValue:     zero,
			pool: &sync.Pool{
				New: func() interface{} {
					return reflect.New(structTyp).Interface()
				},
			},
		})
	}
}
func (c *contextHandler) Handle(rDi *request.Context) {
	ii := c.pool.Get()
	// get's or allocates a new context
	ctx := reflect.ValueOf(ii)
	rDi.Global.InjectInStructValue(ctx.Elem())

	var arguments = [1]reflect.Value{ctx}
	if c.isPtr == false {
		arguments[0] = arguments[0].Elem()
	}

	c.funcValue.Call(arguments[0:])

	ctx.Elem().Set(c.zeroValue)
	c.pool.Put(ii)
}

var acRegex = regexp.MustCompile("/[:*][^/]+")

func (muxmap *Mapper) AddHandler(method, path, action string, filters ...func(*request.Context, request.Flow)) {
	methodByname, isPtr := muxmap.typ.MethodByName(action)
	if !isPtr {
		methodByname, _ = muxmap.typ.Elem().MethodByName(action)
		if methodByname.Type == nil {
			panic("Inválid action " + action + " not found in controller " + muxmap.typ.String())
		}
	}

	muxmap.app.urlGen[muxmap.typ.Elem().String()+"."+action] = acRegex.ReplaceAllStringFunc(path, func(st string) string {
		if st[1] == '*' {
			return "%v"
		}
		return "/%v"
	})

	muxmap.app.AddHandlerContextName(muxmap.Context, muxmap.name, method, path, &contextHandler{
		pool:      muxmap.pool,
		isPtr:     isPtr,
		zeroValue: muxmap.zeroValue,
		funcValue: methodByname.Func,
	}, muxmap.reslice(filters...)...)
}
