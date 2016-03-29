package app

import (
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
	"reflect"
	"regexp"
	"sync"
)

type (
	Mapper struct {
		name string
		typ  reflect.Type
		pool *sync.Pool
		app  *Application
		Di   *context.Context
		*request.Filters
	}

	invokeController struct {
		pool      *sync.Pool
		isPtr     bool
		funcValue reflect.Value
	}

	Controller interface {
		Mux(Mapper)
	}
)

func (app *Application) AddController(controllers ...Controller) {
	for i := 0; i < len(controllers); i++ {
		controller := controllers[i]

		ptrTyp := reflect.TypeOf(controller)
		structTyp := ptrTyp
		if ptrTyp.Kind() == reflect.Ptr {
			structTyp = ptrTyp.Elem()
		} else {
			ptrTyp = reflect.PtrTo(ptrTyp)
		}

		name := structTyp.String()

		// creates a new di for this controller
		newDi := app.Di.Child()
		// creates a new cascade url generator
		myGen := new(ctlGen)
		// injects parent url generator
		newDi.Inject(myGen)

		myGen.urlGen = app.urlGen
		myGen.id = name + "."

		newDi.MapType((*common.URLer)(nil), myGen)

		newFilter := new(request.Filters)
		newDi.Map(newFilter)

		controller.Mux(Mapper{
			name:    name,
			app:     app,
			typ:     ptrTyp,
			Di:      newDi,
			Filters: newFilter,
			pool: &sync.Pool{
				New: func() interface{} {
					return reflect.New(structTyp).Interface()
				},
			},
		})
	}
}

func (c *invokeController) Handle(rDi *request.Context) {

	ii := c.pool.Get()
	defer c.pool.Put(ii)
	rDi.Di.Inject(ii)

	if ii, isInitializer := ii.(interface {
		Init()
	}); isInitializer {
		ii.Init()
	}

	var arguments = [1]reflect.Value{reflect.ValueOf(ii)}
	if c.isPtr == false {
		arguments[0] = arguments[0].Elem()
	}

	c.funcValue.Call(arguments[0:])

	if ii, isFinalizer := ii.(interface {
		Finalize()
	}); isFinalizer {
		ii.Finalize()
	}
}

var acRegex = regexp.MustCompile("[:*][^/]+")

func (muxmap *Mapper) AddHandler(method, path, action string, filters ...func(request.ContextChain)) {
	methodByname, isPtr := muxmap.typ.MethodByName(action)
	if !isPtr {
		methodByname, _ = muxmap.typ.Elem().MethodByName(action)
		if methodByname.Type == nil {
			panic("Inválid action " + action + " not found in controller " + muxmap.typ.String())
		}
	}
	muxmap.app.urlGen[muxmap.typ.Elem().String()+"."+action] = acRegex.ReplaceAllLiteralString(path, "%v")
	muxmap.app.AddHandlerContextName(muxmap.Di, muxmap.name, method, path, &invokeController{
		pool:      muxmap.pool,
		isPtr:     isPtr,
		funcValue: methodByname.Func,
	}, muxmap.MakeFilters(filters...)...)
}

func (muxmap *Mapper) AddPlugin(plugins ...Plugin) {
	LoadPlugins(muxmap.Di, plugins...)
}
