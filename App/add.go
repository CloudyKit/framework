package App

import (
	"github.com/CloudyKit/framework/Request"
	"github.com/CloudyKit/framework/Router"

	"net/http"
	"reflect"
	"sync"
)

type FuncHandler func(*Request.Context)

func (fn FuncHandler) Handle(c *Request.Context) {
	fn(c)
}

func (add *Application) AddFunc(method, path string, fn FuncHandler, filters ...func(Request.Channel)) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *Application) AddHandler(method, path string, handler Request.Handler, filters ...func(Request.Channel)) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *Application) AddHandlerName(name, method, path string, handler Request.Handler, filters ...func(Request.Channel)) {
	app.Router.AddRoute(method, path, func(rw http.ResponseWriter, r *http.Request, v Router.Values) {
		cc := Request.New(Request.Context{Id: name, Rw: rw, Rq: r, Ps: v, Context: app.Child()})
		defer cc.Done()
		cc.Put(cc) // self inject
		(Request.Channel{
			Filters: filters,
			Handler: handler,
			Context: cc,
		}).Next()
	})
}

type MuxHandler interface {
	Mux(Mapper)
}

func (app *Application) AddController(controllers ...MuxHandler) {
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

		controller.Mux(Mapper{
			name: name,
			app:  app,
			typ:  ptrTyp,
			pool: &sync.Pool{
				New: func() interface{} {
					return reflect.New(structTyp).Interface()
				},
			},
		})
	}
}
