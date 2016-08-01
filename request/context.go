package request

import (
	"github.com/CloudyKit/framework/scope"
	"github.com/CloudyKit/router"

	"errors"
	"fmt"
	"github.com/CloudyKit/framework/common"
	"net/http"
	"net/url"
	"reflect"
)

var ContextType = reflect.TypeOf((*Context)(nil))

// GetContext get's a Context from the Global context
func GetContext(cdi *scope.Variables) *Context {
	return cdi.GetByType(ContextType).(*Context)
}

// Context holds context information about the incoming request
type Context struct {
	Name      string           // The name associated with the route
	Variables *scope.Variables // Dependency injection context
	Request   *http.Request    // Request data passed by the router
	Gen       *common.URLGen

	handlers []Handler

	Response   http.ResponseWriter // Response Writer passed by the router
	Parameters router.Parameter    // Route Variables passed by the router
}

// Advance will continue with the request flow
func (ctx *Context) Advance() error {

	if len(ctx.handlers) == 0 {
		return errors.New("request.Context: no available handlers to advance")
	}

	// todo: with this behavior we can allow retry, a func can advance multiple times
	// handlers := ctx.handlers
	// ctx.handlers = ctx.handlers[1:]
	// handlers[0].Handle(ctx)
	// ctx.handlers = handlers

	handler := ctx.handlers[0]
	ctx.handlers = ctx.handlers[1:]
	handler.Handle(ctx)
	return nil
}

// WriteString writes the string txt into the the response
func (ctx *Context) WriteString(txt string) (int, error) {
	return ctx.Response.Write([]byte(txt))
}

// Printf prints a formatted text to response writer
func (ctx *Context) Printf(format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(ctx.Response, format, v...)
}

// Redirect redirects the request to the specified urlStr and send a http StatusFound code
func (ctx *Context) Redirect(urlStr string) {
	ctx.RedirectStatus(urlStr, http.StatusFound)
}

// RedirectStatus redirects the request to the specified urlStr and send the the status code specified by httpStatus
func (ctx *Context) RedirectStatus(urlStr string, httpStatus int) {
	http.Redirect(ctx.Response, ctx.Request, urlStr, httpStatus)
}

// ParamByName returns a parameter from the url route, ParamByName is shortcut for Context.Parameters.ByName method
func (ctx *Context) ParamByName(name string) string {
	return ctx.Parameters.ByName(name)
}

// FormByName  returns a form value from the request, FormByName is shortcut for Context.Request.Form.Get method
func (ctx *Context) FormByName(name string) string {
	if ctx.Request.PostForm == nil {
		ctx.Request.ParseForm()
	}
	return ctx.Request.PostForm.Get(name)
}

// URLFormByName  returns a form value from the request, FormByName is shortcut for Context.Request.Form.Get method
func (ctx *Context) URLFormByName(name string) string {
	if ctx.Request.Form == nil {
		ctx.Request.ParseForm()
	}
	return ctx.Request.Form.Get(name)
}

// CookieByName returns a cookie value from the request
func (ctx *Context) CookieByName(name string) (value string) {
	if cookie, _ := ctx.Request.Cookie(name); cookie != nil {
		value, _ = url.QueryUnescape(cookie.Value)
	}
	return
}
