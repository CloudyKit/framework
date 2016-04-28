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
	"reflect"
	"strings"
)

var Default = New()

func New() *Application {

	newApp := &Application{Context: context.New(), router: router.New(), urlGen: make(urlGen), Filters: new(request.Filters)}

	// setups default err reporter
	newApp.Notifier = errors.NewNotifier(newApp.Context, reporters.LogReporter{})
	// provide application urlGen as URLer
	newApp.Context.MapType((*common.URLer)(nil), newApp.urlGen)
	// provide Filters plugins added in the application can setup filters
	newApp.Context.Map(newApp.Filters)
	// provide the Router
	newApp.Context.Map(newApp.router)
	// provide the app
	newApp.Context.Map(newApp)
	// provide error catcher
	newApp.Context.Map(newApp.Notifier)
	return newApp
}

type Application struct {
	Context *context.Context // Di dependency injection context
	router  *router.Router   // Router

	urlGen urlGen //
	*request.Filters
	Notifier errors.Notifier
}

type Bootstrapper interface {
	Bootstrap(app *Application)
}

type Plugin interface {
	PluginInit(*context.Context)
}

func LoadPlugins(di *context.Context, plugins ...Plugin) {
	for i := 0; i < len(plugins); i++ {
		di.Inject(plugins[i])
		plugins[i].PluginInit(di)
	}
}

func (app *Application) AddPlugin(plugins ...Plugin) {
	LoadPlugins(app.Context, plugins...)
}

func (app Application) Bootstrap(b ...Bootstrapper) {
	c := app.Context.Child()
	for i := 0; i < len(b); i++ {
		bv := reflect.ValueOf(b[i])
		if bv.Kind() == reflect.Ptr {
			bv = bv.Elem()
			if bv.Kind() == reflect.Struct {
				c.InjectStructValue(bv)
			}
		}
		b[i].Bootstrap(&app)
	}
	c.Done()
}

func (app *Application) Done() {
	app.Context.Done()
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
	app.AddHandlerContextName(app.Context, name, method, path, handler, filters...)
}

// AddHandlerContextName accepts a context, a name identifier, http method|methods, pattern path, handler and filters
// ex: one handler app.AddHandlerContextName(myContext,"mySectionIdentifier","GET", "/public",fileServer,checkAuth)
//     multiples handles app.AddHandlerContextName(myContext,"mySectionIdentifier","GET|POST|SEARCH", "/products",productHandler,checkAuth)
func (app *Application) AddHandlerContextName(context *context.Context, name, method, path string, handler request.Handler, filters ...func(request.ContextChain)) {
	filters = app.MakeFilters(filters...)
	if context == nil {
		context = app.Context
	}
	for _, method := range strings.Split(method, "|") {
		app.router.AddRoute(method, path, func(rw http.ResponseWriter, r *http.Request, v router.Parameter) {
			cc := request.New(request.Context{Name: name, Response: rw, Request: r, Parameters: v, Context: context.Child()})
			defer cc.Context.Done() // call finalizers
			cc.Context.Map(cc)      // self inject
			cc.Notifier = cc.Context.Get(cc.Notifier).(errors.Notifier)
			request.NewContextChain(cc, handler, filters).Next()
		})
	}
}

// RunServer runs the application in the default HTTP server
// optional arguments port,certfile,keyfile
func (app *Application) RunServer(args ...string) *Application {
	var host string
	if len(args) > 0 {
		host = os.Getenv(args[0])
		if host == "" {
			host = args[0]
		}
	} else {
		host = os.Getenv("PORT")
	}
	if len(args) < 2 {
		app.Notifier.ErrNotify(http.ListenAndServe(host, app.router))
	} else if len(args) == 3 {
		app.Notifier.ErrNotify(http.ListenAndServeTLS(host, args[1], args[2], app.router))
	} else {
		app.Notifier.ErrNotify(fmt.Errorf("Inválid number of arguments on App.RunServer"))
	}
	return app
}
