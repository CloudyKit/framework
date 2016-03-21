package Request

import (
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Router"

	"net/http"
	"sync"
)

type Context struct {
	Di.Context        // Depedency injection base
	Id         string // The id associated with this route

	Rq *http.Request       // Request data passed by the router
	Rw http.ResponseWriter // Response Writer passed by the router
	Rv Router.Values       // Route Variables passed by the router
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

	di := cc.Context // copy

	// reset the values
	cc.Rq = nil
	cc.Rw = nil
	cc.Rv = Router.Values{}
	cc.Context = Di.Context{}

	// recycle cc
	_New.Put(cc)

	// recycle depedency injection table
	di.Done()
}
