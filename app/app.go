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

var AppType = reflect.TypeOf((*App)(nil))

func Get(c *cdi.DI) *App {
	return c.Val4Type(AppType).(*App)
}

var Default = New()

func New() *App {

	newApp := &App{Global: cdi.New(), Router: router.New(), urlGen: make(urlGen)}

	// provide application urlGen as URLer
	newApp.Global.MapType(common.URLerType, newApp.urlGen)

	// provide the Router
	newApp.Global.Map(newApp.Router)
	// provide the app
	newApp.Global.Map(newApp)
	return newApp
}

type filterManager struct {
	modified bool
	filters  []func(*request.Context, request.Flow)
}

func (f *filterManager) AddFilter(filters ...func(*request.Context, request.Flow)) {
	f.filters = append(f.filters, filters...)
	f.modified = true
}

func (f *filterManager) reslice(filters ...func(*request.Context, request.Flow)) []func(*request.Context, request.Flow) {
	if f.modified {
		newFilter := make([]func(*request.Context, request.Flow), 0, len(f.filters)+len(filters))
		newFilter = append(newFilter, f.filters...)
		newFilter = append(newFilter, filters...)
		f.modified = false
		return newFilter
	}
	return f.filters[0:len(f.filters)]
}

type App struct {
	Global *cdi.DI        // App Global dependency injection context
	Router *router.Router // Router
	Prefix string         // Prefix prefix for path added in this app
	urlGen urlGen
	filterManager
}

type Bootstrapper interface {
	Bootstrap(app *App)
}

func (app App) Bootstrap(b ...Bootstrapper) {
	c := app.Global.Child()
	defer c.Done()

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
}

func (app *App) Done() {
	app.Global.Done()
}

type funcHandler func(*request.Context)

func (fn funcHandler) Handle(c *request.Context) {
	fn(c)
}

func (add *App) AddHandlerFunc(method, path string, fn funcHandler, filters ...func(*request.Context, request.Flow)) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *App) AddHandler(method, path string, handler request.Handler, filters ...func(*request.Context, request.Flow)) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *App) AddHandlerName(name, method, path string, handler request.Handler, filters ...func(*request.Context, request.Flow)) {
	app.AddHandlerContextName(app.Global, name, method, path, handler, filters...)
}

// AddHandlerContextName accepts a context, a name identifier, http method|methods, pattern path, handler and filters
// ex: one handler app.AddHandlerContextName(myContext,"mySectionIdentifier","GET", "/public",fileServer,checkAuth)
//     multiples handles app.AddHandlerContextName(myContext,"mySectionIdentifier","GET|POST|SEARCH", "/products",productHandler,checkAuth)
func (app *App) AddHandlerContextName(context *cdi.DI, name, method, path string, handler request.Handler, filters ...func(*request.Context, request.Flow)) {

	filters = app.reslice(filters...)

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

func (app *App) host(host string) (servein string) {
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

func (app *App) RunServer(host string) error {
	return http.ListenAndServe(app.host(host), app.Router)
}

func (app *App) RunServerTls(host, certfile, keyfile string) error {
	return http.ListenAndServeTLS(app.host(host), certfile, keyfile, app.Router)
}
