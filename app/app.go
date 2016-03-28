package app

import (

	"github.com/CloudyKit/framework/errors/reporters"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/errors"
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/di"
	"github.com/CloudyKit/router"

	"net/http"
)

var Default = New()

func New() *Application {

	newApp := &Application{Di: di.New(), Router: router.New(), urlGen: make(urlGen), Filters: new(request.Filters)}

	newApp.Error = reporters.LogReporter{}
	// provide application urlGen as URLer
	newApp.Di.Set((*common.URLer)(nil), newApp.urlGen)
	// provide Filters plugins added in the application can setup filters
	newApp.Di.Map(newApp.Filters)
	// provide the Router
	newApp.Di.Map(newApp.Router)
	// provide the app
	newApp.Di.Map(newApp)
	// provide error catcher
	newApp.Di.Map(newApp.Error)

	return newApp
}

type Application struct {
	Di     *di.Context
	Router *router.Router

	urlGen urlGen
	*request.Filters
	Error  errors.Catcher
}

type Plugin interface {
	Init(*di.Context)
}

func LoadPlugins(di *di.Context, plugins ...Plugin) {
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

type FuncHandler func(*request.Context)

func (fn FuncHandler) Handle(c *request.Context) {
	fn(c)
}

func (add *Application) AddFunc(method, path string, fn FuncHandler, filters ...func(request.Channel)) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *Application) AddHandler(method, path string, handler request.Handler, filters ...func(request.Channel)) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *Application) AddHandlerName(name, method, path string, handler request.Handler, filters ...func(request.Channel)) {
	app.AddHandlerContextName(app.Di, name, method, path, handler, filters...)
}

func (app *Application) AddHandlerContextName(context *di.Context, name, method, path string, handler request.Handler, filters ...func(request.Channel)) {
	filters = app.MakeFilters(filters...)
	app.Router.AddRoute(method, path, func(rw http.ResponseWriter, r *http.Request, v router.Parameter) {
		cc := request.New(request.Context{Name: name, Response: rw, Request: r, Parameters: v, Di: context.Child()})
		defer cc.Done() // call finalizers
		cc.Di.Map(cc)   // self inject
		(request.Channel{
			Filters: filters,
			Handler: handler,
			Context: cc,
		}).Next()
	})
}
