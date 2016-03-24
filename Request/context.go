package Request

import (
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Router"
	"github.com/CloudyKit/framework/Validator"

	"encoding/json"
	"net/http"
	"sync"
)

type Context struct {
	*http.Request // Request data passed by the router

	Id string // The id associated with this route

	Di *Di.Context // Dependency injection context

	Rw http.ResponseWriter // Response Writer passed by the router
	Rv Router.Values       // Route Variables passed by the router
}

func (cc *Context) ValidateRoute(c func(Validator.At)) Validator.Result {
	return Validator.Run(Validator.NewRouterValueProvider(cc.Rv), c)
}

func (cc *Context) ValidateGet(c func(Validator.At)) Validator.Result {
	return Validator.Run(Validator.NewRequestValueProvider(cc.Request), c)
}

func (cc *Context) ValidatePost(c func(Validator.At)) Validator.Result {
	cc.ParseForm()
	return Validator.Run(Validator.NewURLValueProvider(cc.PostForm), c)
}

func (cc *Context) DecodeJson(target interface{}) error {
	return json.NewDecoder(cc.Body).Decode(target)
}

func (cc *Context) EncodeJson(from interface{}) error {
	return json.NewEncoder(cc.Rw).Encode(from)
}

func (cc *Context) WriteString(txt string) (int, error) {
	return cc.Rw.Write([]byte(txt))
}

func (cc *Context) Redirect(urlStr string) {
	cc.RedirectCode(urlStr, http.StatusFound)
}

func (cc *Context) RedirectCode(urlStr string, code int) {
	http.Redirect(cc.Rw, cc.Request, urlStr, code)
}

var _New = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func New(c Context) (cc *Context) {
	cc = _New.Get().(*Context)
	*cc = c
	return
}

func (cc *Context) Done() {
	// recycle cc
	_New.Put(cc)
	// recycle depedency injection table
	cc.Di.Done()
}
