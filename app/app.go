package app

import (
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
	"github.com/CloudyKit/router"

	"github.com/CloudyKit/framework/events"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
)

var AppType = reflect.TypeOf((*App)(nil))

func Get(c *scope.Variables) *App {
	return c.GetByType(AppType).(*App)
}

var Default = New()

func New() *App {
	_app := &App{Variables: scope.New(), Router: router.New(), urlGen: make(urlGen), emitter: events.NewEmitter()}

	// provide application urlGen as URLer
	_app.Variables.MapType(common.URLGenType, _app.urlGen)
	// provide the Router
	_app.Variables.Map(_app.Router)
	// provide the app
	_app.Variables.MapType(AppType, _app)
	_app.Variables.MapType(events.EmitterType, _app.emitter)
	return _app
}

type filterHandlers struct {
	filters []request.Handler
}

func (f *filterHandlers) ResetMiddleHandlers() {
	f.filters = nil
}

// AddMiddleHandler adds filters to the request chain
func (f *filterHandlers) AddMiddleHandlers(filters ...request.Handler) {
	f.filters = append(f.filters, filters...)
}

//func (f *filterHandlers) AddMiddleHandlersFunc(filters ...request.HandlerFunc) {
//	nlen := len(filters) + len(f.filters)
//
//	nfilters := make([]request.Handler, nlen)
//
//	copy(nfilters, f.filters)
//
//	for i, j := len(f.filters), 0; i < nlen; i, j = i + 1, j + 1 {
//		nfilters[i] = filters[j]
//	}
//
//	f.filters = nfilters
//}

func (f *filterHandlers) reslice(filters ...request.Handler) []request.Handler {
	newFilter := make([]request.Handler, 0, len(f.filters)+len(filters))
	newFilter = append(newFilter, f.filters...)
	newFilter = append(newFilter, filters...)
	return newFilter
}

type emitter interface {
	Subscribe(groups string, handler interface{}) *events.Emitter
	Emit(groupName, key string, context interface{}) (canceled bool, err error)
}

// App app holds your top level data for you application
// Router, Emitter, Scope
type App struct {
	emitter

	Variables *scope.Variables // App Variables dependency injection context
	Router    *router.Router   // Router
	Prefix    string           // Prefix prefix for path added in this app
	urlGen    urlGen
	filterHandlers
}

// Component represents a application component, a component need to implement
// a bootstrap method which is responsible to setup the component with the app,
// ex: register a type Providers, or add middleware handler
type Component interface {
	Bootstrap(app *App)
}

// Root returns the root app
func (app *App) Root() *App {
	return Get(app.Variables)
}

// Snapshot causes a sub app to be created and inserted in the scope
// calling app.Root will return the created sub app
func (app *App) Snapshot() *App {
	_app := *app

	_app.Variables = app.Variables.Inherit()
	_app.Variables.MapType(AppType, _app)

	return &_app
}

// ComponentFunc func implementing Component interface
type ComponentFunc func(*App)

func (component ComponentFunc) Bootstrap(a *App) {
	component(a)
}

// Bootstrap bootstrap a list of components, a sub scope will be created, and a copy of the
// original app is used, in such form that modifing the app.Prefix will not reflect outside this
// call.
func (app App) Bootstrap(b ...Component) {
	c := app.Variables.Inherit()
	defer c.EndForce() // require 0 references at this point

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

// End same as app.Variables.End() invoke this func before exiting the app to cleanup
func (app *App) End() {
	app.Variables.End()
}

// AddHandlerFunc register a func handler, see: request.Handler
func (add *App) AddHandlerFunc(method, path string, fn request.HandlerFunc, filters ...request.Handler) {
	add.AddHandler(method, path, fn, filters...)
}

// AddHandlerFunc register a handler, see: request.Handler
func (app *App) AddHandler(method, path string, handler request.Handler, filters ...request.Handler) {
	app.AddHandlerName("", method, path, handler, filters...)
}

// AddHandlerFunc register a named handler, see: request.Handler
func (app *App) AddHandlerName(name, method, path string, handler request.Handler, filters ...request.Handler) {
	app.AddHandlerContextName(app.Variables, name, method, path, handler, filters...)
}

// AddHandlerContextName accepts a context, a name identifier, http method|methods, pattern path, handler and filters
// ex: one handler app.AddHandlerContextName(myContext,"mySectionIdentifier","GET", "/public",fileServer,checkAuth)
//     multiples handles app.AddHandlerContextName(myContext,"mySectionIdentifier","GET|POST|SEARCH", "/products",productHandler,checkAuth)
func (app *App) AddHandlerContextName(variables *scope.Variables, name, method, path string, handler request.Handler, filters ...request.Handler) {

	filters = append(app.reslice(filters...), handler)

	if variables == nil {
		variables = app.Variables
	}

	for _, method := range strings.Split(method, "|") {
		app.Router.AddRoute(method, app.Prefix+path, func(rw http.ResponseWriter, r *http.Request, v router.Parameter) {

			c := newRequestContext()
			defer requestRecover(c)

			request.Advance(c, name, rw, r, v, variables.Inherit(), filters)
		})
	}
}

// requestRecover finalizes and cleanup request allocated scope variables
func requestRecover(c *request.Context) {

	variables := c.Variables
	request_Context_pool.Put(c)

	// we call scope EndForce, this require that all children scopes Ended in this call if not
	// panic is raised
	variables.EndForce()
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

// RunServer runs the server with the specified host
// Calling this func will emit a "app.run" event in the app
func (app *App) RunServer(host string) error {
	app.Emit("app.run", host, app)
	return http.ListenAndServe(host, app.Router)
}

// RunServerTLS runs the server in tls mode
// Calling this func will emit a "app.run.tls" event in the app
func (app *App) RunServerTLS(host, certfile, keyfile string) error {
	app.Emit("app.run.tls", host, app)
	return http.ListenAndServeTLS(app.host(host), certfile, keyfile, app.Router)
}

var request_Context_pool = sync.Pool{
	New: func() interface{} {
		return new(request.Context)
	},
}

func newRequestContext() *request.Context {
	return request_Context_pool.Get().(*request.Context)
}
