package app

import (
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/router"

	"net/http"
	"os"
	"reflect"
	"strings"
)

var Default = New()

func New() *Application {

	newApp := &Application{Global: cdi.New(), Router: router.New(), urlGen: make(urlGen), Filters: new(request.Filters)}

	// provide application urlGen as URLer
	newApp.Global.MapType((*common.URLer)(nil), newApp.urlGen)
	// provide Filters plugins added in the application can setup filters
	newApp.Global.Map(newApp.Filters)
	// provide the Router
	newApp.Global.Map(newApp.Router)
	// provide the app
	newApp.Global.Map(newApp)
	return newApp
}

type Application struct {
	Global *cdi.DI        // Di dependency injection context
	Router *router.Router // Router

	urlGen urlGen //
	*request.Filters
	Prefix string
}

type Bootstrapper interface {
	Bootstrap(app *Application)
}

type Plugin interface {
	PluginInit(*cdi.DI)
}

func LoadPlugins(di *cdi.DI, plugins ...Plugin) {
	for i := 0; i < len(plugins); i++ {
		di.Inject(plugins[i])
		plugins[i].PluginInit(di)
	}
}

func (app *Application) AddPlugin(plugins ...Plugin) {
	LoadPlugins(app.Global, plugins...)
}

func (app Application) Bootstrap(b ...Bootstrapper) {
	c := app.Global.Child()
	for i := 0; i < len(b); i++ {
		bv := reflect.ValueOf(b[i])
		if bv.Kind() == reflect.Ptr {
			bv = bv.Elem()
			if bv.Kind() == reflect.Struct {
				c.InjectInStructValue(bv)
			}
		}
		b[i].Bootstrap(&app)
	}
	c.Done()
}

func (app *Application) Done() {
	app.Global.Done()
}

type funcHandler func(*request.Context)

func (fn funcHandler) Handle(c *request.Context) {
	fn(c)
}

func (add *Application) AddHandlerFunc(method, path string, fn funcHandler, filters ...func(*request.Context, request.Flow)) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *Application) AddHandler(method, path string, handler request.Handler, filters ...func(*request.Context, request.Flow)) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *Application) AddHandlerName(name, method, path string, handler request.Handler, filters ...func(*request.Context, request.Flow)) {
	app.AddHandlerContextName(app.Global, name, method, path, handler, filters...)
}

// AddHandlerContextName accepts a context, a name identifier, http method|methods, pattern path, handler and filters
// ex: one handler app.AddHandlerContextName(myContext,"mySectionIdentifier","GET", "/public",fileServer,checkAuth)
//     multiples handles app.AddHandlerContextName(myContext,"mySectionIdentifier","GET|POST|SEARCH", "/products",productHandler,checkAuth)
func (app *Application) AddHandlerContextName(context *cdi.DI, name, method, path string, handler request.Handler, filters ...func(*request.Context, request.Flow)) {

	filters = app.MakeFilters(filters...)

	if context == nil {
		context = app.Global
	}

	for _, method := range strings.Split(method, "|") {
		app.Router.AddRoute(method, app.Prefix+path, func(rw http.ResponseWriter, r *http.Request, v router.Parameter) {
			cc := request.New(request.Context{Name: name, Response: rw, Request: r, Parameters: v, Global: context.Child()})
			defer cc.Global.Done() // call finalizers
			cc.Global.Map(cc)      // self inject
			request.NewContextChain(cc, handler, filters).Continue()
		})
	}
}

func (app *Application) host(host string) (servein string) {
	// if host is empty set host apphost
	if host == "" {
		host = "apphost"
	}
	// check if host is an env variable containing a host string
	servein = os.Getenv(host)
	// if host is not an env variable than is a host string
	if servein == "" {
		servein = host
	}
	return
}

func (app *Application) RunServer(host string) error {
	return http.ListenAndServe(app.host(host), app.Router)
}

func (app *Application) RunServerTls(host, certfile, keyfile string) error {
	return http.ListenAndServeTLS(app.host(host), certfile, keyfile, app.Router)
}
