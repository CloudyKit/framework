package Request

import (
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Router"

	"encoding/json"
	"net/http"
	"sync"
)

type Context struct {
	Di.Context

	Rq *http.Request
	Rw http.ResponseWriter

	Ps Router.Values
	Id string
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

func (cc *Context) ReceiveJson(target interface{}) error {
	return json.NewDecoder(cc.Rq.Body).Decode(target)
}

func (cc *Context) SendJson(from interface{}) error {
	return json.NewEncoder(cc.Rw).Encode(from)
}

func (cc *Context) Done() {
	defer _New.Put(cc)
	cc.Rq = nil
	cc.Rw = nil
	cc.Ps.Values = nil
	cc.Context.Done()
	cc.Context = Di.Context{}
}

func (cc *Context) SendText(txt string) (int, error) {
	return cc.Rw.Write([]byte(txt))
}
