package App

import (
	"github.com/CloudyKit/framework/Common"
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"github.com/CloudyKit/framework/Router"

	"net/http"
)

var Default = New()

func New() *Application {
	newApp := &Application{Di: Di.New(), Router: Router.New(), urlGen: make(urlGen), Filters: new(Request.Filters)}
	newApp.Di.Set((*Common.URLer)(nil), newApp.urlGen) // injects
	newApp.Di.Put(newApp.Filters)                      // make this available in the di
	return newApp
}

type Application struct {
	Di     *Di.Context
	Router *Router.Router

	urlGen urlGen
	*Request.Filters
}

type Plugin interface {
	Init(*Di.Context)
}

func LoadPlugins(di *Di.Context, plugins ...Plugin) {
	for i := 0; i < len(plugins); i++ {
		di.Inject(plugins[i])
		plugins[i].Init(di)
	}
}

func (app *Application) AddPlugin(plugins ...Plugin) {
	LoadPlugins(app.Di, plugins...)
}

func (app *Application) Done() {
	app.Di.Done()
}

type FuncHandler func(*Request.Context)

func (fn FuncHandler) Handle(c *Request.Context) {
	fn(c)
}

func (add *Application) AddFunc(method, path string, fn FuncHandler, filters ...Request.Filter) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *Application) AddHandler(method, path string, handler Request.Handler, filters ...Request.Filter) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *Application) AddHandlerName(name, method, path string, handler Request.Handler, filters ...Request.Filter) {
	app.AddHandlerContextName(app.Di, name, method, path, handler, filters...)
}

func (app *Application) AddHandlerContextName(context *Di.Context, name, method, path string, handler Request.Handler, filters ...Request.Filter) {
	filters = app.MakeFilters(filters...)
	app.Router.AddRoute(method, path, func(rw http.ResponseWriter, r *http.Request, v Router.Values) {
		cc := Request.New(Request.Context{Id: name, Rw: rw, Request: r, Rv: v, Di: context.Child()})
		defer cc.Done() // call finalizers
		cc.Di.Put(cc)   // self inject
		(Request.Channel{
			Filters: filters,
			Handler: handler,
			Context: cc,
		}).Next()
	})
}
