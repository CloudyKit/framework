package app

import (
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/errors"
	"github.com/CloudyKit/framework/errors/reporters"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/router"

	"fmt"
	"net/http"
	"os"
)

var Default = New()

func New() *Application {

	newApp := &Application{Di: context.New(), router: router.New(), urlGen: make(urlGen), Filters: new(request.Filters)}

	// setups default err reporter
	newApp.Notifier = errors.NewCatcher(newApp.Di, reporters.LogReporter{})
	// provide application urlGen as URLer
	newApp.Di.MapType((*common.URLer)(nil), newApp.urlGen)
	// provide Filters plugins added in the application can setup filters
	newApp.Di.Map(newApp.Filters)
	// provide the Router
	newApp.Di.Map(newApp.router)
	// provide the app
	newApp.Di.Map(newApp)
	// provide error catcher
	newApp.Di.Map(newApp.Notifier)
	return newApp
}

type Application struct {
	Di     *context.Context // Di dependency injection context
	router *router.Router   // Router

	urlGen urlGen //
	*request.Filters
	Notifier errors.Notifier
}

type Plugin interface {
	Init(*context.Context)
}

func LoadPlugins(di *context.Context, plugins ...Plugin) {
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

func (add *Application) AddFunc(method, path string, fn FuncHandler, filters ...func(request.ContextChain)) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *Application) AddHandler(method, path string, handler request.Handler, filters ...func(request.ContextChain)) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *Application) AddHandlerName(name, method, path string, handler request.Handler, filters ...func(request.ContextChain)) {
	app.AddHandlerContextName(app.Di, name, method, path, handler, filters...)
}

func (app *Application) AddHandlerContextName(context *context.Context, name, method, path string, handler request.Handler, filters ...func(request.ContextChain)) {
	filters = app.MakeFilters(filters...)
	app.router.AddRoute(method, path, func(rw http.ResponseWriter, r *http.Request, v router.Parameter) {
		cc := request.New(request.Context{Name: name, Response: rw, Request: r, Parameters: v, Di: context.Child()})
		defer cc.Done() // call finalizers
		cc.Di.Map(cc)   // self inject
		request.NewContextChain(cc, handler, filters).Next()
	})
}

func (app *Application) RunServer(args ...string) *Application {
	var host string
	if len(args) > 0 {
		host = args[0]
	} else {
		host = os.Getenv("PORT")
	}
	if len(args) < 2 {
		app.Notifier.NotifyIfNotNil(http.ListenAndServe(host, app.router))
	} else if len(args) == 3 {
		app.Notifier.NotifyIfNotNil(http.ListenAndServeTLS(host, args[1], args[2], app.router))
	} else {
		app.Notifier.NotifyIfNotNil(fmt.Errorf("Inválid number of arguments on App.RunServer"))
	}
	return app
}
