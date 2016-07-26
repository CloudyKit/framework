package request

import (
	"github.com/CloudyKit/framework/scope"
	"github.com/CloudyKit/router"

	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

var ContextType = reflect.TypeOf((*Context)(nil))

//GetContext get's a Context from the Global context
func GetContext(cdi *scope.Variables) *Context {
	return cdi.GetByType(ContextType).(*Context)
}

//Context holds context information about the incoming request
type Context struct {
	Name       string              //The name associated with the route
	Variables  *scope.Variables    //Dependency injection context
	Request    *http.Request       //Request data passed by the router
	Response   http.ResponseWriter //Response Writer passed by the router
	Parameters router.Parameter    //Route Variables passed by the router
}

//WriteString writes the string txt into the the response
func (cc *Context) WriteString(txt string) (int, error) {
	return cc.Response.Write([]byte(txt))
}

//Printf
func (cc *Context) Printf(format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(cc.Response, format, v...)
}

//Redirect redirects the request to the specified urlStr and send a http StatusFound code
func (c *Context) Redirect(urlStr string) {
	c.RedirectStatus(urlStr, http.StatusFound)
}

//RedirectStatus redirects the request to the specified urlStr and send the the status code specified by httpStatus
func (c *Context) RedirectStatus(urlStr string, httpStatus int) {
	http.Redirect(c.Response, c.Request, urlStr, httpStatus)
}

//ParamByName returns a parameter from the url route, ParamByName is shortcut for Context.Parameters.ByName method
func (cc *Context) ParamByName(name string) string {
	return cc.Parameters.ByName(name)
}

//FormByName  returns a form value from the request, FormByName is shortcut for Context.Request.Form.Get method
func (cc *Context) FormByName(name string) string {
	if cc.Request.PostForm == nil {
		cc.Request.ParseForm()
	}
	return cc.Request.PostForm.Get(name)
}

//URLFormByName  returns a form value from the request, FormByName is shortcut for Context.Request.Form.Get method
func (cc *Context) URLFormByName(name string) string {
	if cc.Request.Form == nil {
		cc.Request.ParseForm()
	}
	return cc.Request.Form.Get(name)
}

//CookieByName returns a cookie value from the request
func (cc *Context) CookieByName(name string) (value string) {
	if cookie, _ := cc.Request.Cookie(name); cookie != nil {
		value, _ = url.QueryUnescape(cookie.Value)
	}
	return
}
