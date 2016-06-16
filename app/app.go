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
	"sync"
)

var AppType = reflect.TypeOf((*App)(nil))

func Get(c *cdi.Global) *App {
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
	newApp.Global.MapType(AppType, newApp)
	return newApp
}

type filterManager struct {
	filters []request.Filter
}

// AddFilter adds filters to the request chain
func (f *filterManager) AddFilter(filters ...request.Filter) {
	f.filters = append(f.filters, filters...)
}

func (f *filterManager) reslice(filters ...request.Filter) []request.Filter {
	if len(filters) > 0 {
		newFilter := make([]request.Filter, 0, len(f.filters)+len(filters))
		newFilter = append(newFilter, f.filters...)
		newFilter = append(newFilter, filters...)
		return newFilter
	}
	return f.filters[0:len(f.filters)]
}

type App struct {
	Global *cdi.Global    // App Global dependency injection context
	Router *router.Router // Router
	Prefix string         // Prefix prefix for path added in this app
	urlGen urlGen
	filterManager
}

// Component an component
type Component interface {
	Bootstrap(app *App)
}

// Root returns the root app
func (app *App) Root() *App {
	return Get(app.Global)
}

type ComponentFunc func(*App)

func (component ComponentFunc) Bootstrap(a *App) {
	component(a)
}

// Bootstrap bootstrap a list of components, Bootstrap will created a child CDI context used
func (app App) Bootstrap(b ...Component) {
	c := app.Global.Child()
	defer c.Done4C() // require 0 references at this point

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

// Done invoke *(cdi.DI).Done
func (app *App) Done() {
	app.Global.Done()
}

type funcHandler func(*request.Context)

func (fn funcHandler) Handle(c *request.Context) {
	fn(c)
}

func (add *App) AddHandlerFunc(method, path string, fn funcHandler, filters ...request.Filter) {
	add.AddHandler(method, path, fn, filters...)
}

func (app *App) AddHandler(method, path string, handler request.Handler, filters ...request.Filter) {
	app.AddHandlerName("", method, path, handler, filters...)
}

func (app *App) AddHandlerName(name, method, path string, handler request.Handler, filters ...request.Filter) {
	app.AddHandlerContextName(app.Global, name, method, path, handler, filters...)
}

// AddHandlerContextName accepts a context, a name identifier, http method|methods, pattern path, handler and filters
// ex: one handler app.AddHandlerContextName(myContext,"mySectionIdentifier","GET", "/public",fileServer,checkAuth)
//     multiples handles app.AddHandlerContextName(myContext,"mySectionIdentifier","GET|POST|SEARCH", "/products",productHandler,checkAuth)
func (app *App) AddHandlerContextName(context *cdi.Global, name, method, path string, handler request.Handler, filters ...request.Filter) {

	filters = app.reslice(filters...)

	if context == nil {
		context = app.Global
	}

	for _, method := range strings.Split(method, "|") {
		app.Router.AddRoute(method, app.Prefix+path, func(rw http.ResponseWriter, r *http.Request, v router.Parameter) {

			c := newContext(request.Context{Name: name, Response: rw, Request: r, Parameters: v, Global: context.Child()})
			defer func() {
				global := c.Global
				contextPool.Put(c)
				global.Done4C() // at this point all finalizers need to be called
			}() // call finalizers
			c.Global.Map(c)
			request.NewRequestFlow(c, handler, filters).Continue()
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
	return http.ListenAndServe(host, app.Router)
}

func (app *App) RunServerTls(host, certfile, keyfile string) error {
	return http.ListenAndServeTLS(app.host(host), certfile, keyfile, app.Router)
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return new(request.Context)
	},
}

// New make a new request context,
func newContext(c request.Context) (cc *request.Context) {
	cc = contextPool.Get().(*request.Context)
	*cc = c
	return
}
