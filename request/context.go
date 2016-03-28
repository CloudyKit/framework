package request

import (
	"github.com/CloudyKit/framework/di"
	"github.com/CloudyKit/framework/errors"
	"github.com/CloudyKit/framework/validator"
	"github.com/CloudyKit/router"

	"encoding/json"
	"net/http"
	"net/url"
	"sync"
)

type Context struct {
	Name       string              // The name associated with the route
	Di         *di.Context         // Dependency injection context
	Request    *http.Request       // Request data passed by the router
	Response   http.ResponseWriter // Response Writer passed by the router
	Parameters router.Parameter    // Route Variables passed by the router
	Error      errors.Catcher      // Error errors reporters
}

func (cc *Context) ValidateRoute(c func(validator.At)) validator.Result {
	return validator.Run(validator.NewRouterValueProvider(cc.Parameters), c)
}

func (cc *Context) ValidateGet(c func(validator.At)) validator.Result {
	return validator.Run(validator.NewRequestValueProvider(cc.Request), c)
}

func (cc *Context) ValidatePost(c func(validator.At)) validator.Result {
	cc.Request.ParseForm()
	return validator.Run(validator.NewURLValueProvider(cc.Request.PostForm), c)
}

func (cc *Context) JsonReadto(target interface{}) error {
	return json.NewDecoder(cc.Request.Body).Decode(target)
}

func (cc *Context) JsonWritefrom(from interface{}) error {
	return json.NewEncoder(cc.Response).Encode(from)
}

func (cc *Context) WriteString(txt string) (int, error) {
	return cc.Response.Write([]byte(txt))
}

func (cc *Context) Redirect(urlStr string) {
	cc.RedirectCode(urlStr, http.StatusFound)
}

func (cc *Context) RedirectCode(urlStr string, code int) {
	http.Redirect(cc.Response, cc.Request, urlStr, code)
}

func (cc *Context) Get(name string) string {
	return cc.Request.Form.Get(name)
}

func (cc *Context) Post(name string) string {
	cc.Request.ParseForm()
	return cc.Request.PostForm.Get(name)
}

func (cc *Context) Cookie(name string) (value string) {
	if cookie, _ := cc.Request.Cookie(name); cookie != nil {
		var err error
		value, err = url.QueryUnescape(cookie.Value)
		cc.Error.ReportIfNotNil(cc.Di, err)
	}
	return
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func New(c Context) (cc *Context) {
	cc = contextPool.Get().(*Context)
	*cc = c
	return
}

func (cc *Context) Done() {
	contextPool.Put(cc)
	cc.Di.Done()
}
