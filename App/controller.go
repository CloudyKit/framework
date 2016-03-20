package App

import (
	"github.com/CloudyKit/framework/Request"
	"reflect"
	"regexp"
	"sync"
)

type Mapper struct {
	name string
	typ  reflect.Type
	pool *sync.Pool
	app  *Application
}

type controllerHandler struct {
	pool      *sync.Pool
	isPtr     bool
	funcValue reflect.Value
}

func (c *controllerHandler) Handle(rcxt *Request.Context) {
	ii := c.pool.Get()
	defer c.pool.Put(ii)
	rcxt.Context.Inject(ii)

	var arguments = [1]reflect.Value{reflect.ValueOf(ii)}
	if c.isPtr == false {
		arguments[0] = arguments[0].Elem()
	}
	c.funcValue.Call(arguments[0:])
}

var acRegex = regexp.MustCompile("[:*][^/]+")

func (muxmap Mapper) AddHandler(method, path, action string, filters ...func(Request.Channel)) {
	methodByname, isPtr := muxmap.typ.MethodByName(action)
	if !isPtr {
		methodByname, _ = muxmap.typ.Elem().MethodByName(action)
		if methodByname.Type == nil {
			panic("Inválid action " + action + " not found in controller " + muxmap.typ.String())
		}
	}

	muxmap.app.Gen[muxmap.typ.Elem().String()+"."+action] = acRegex.ReplaceAllLiteralString(path, "%v")
	muxmap.app.AddHandlerName(muxmap.name, method, path, &controllerHandler{
		pool:      muxmap.pool,
		isPtr:     isPtr,
		funcValue: methodByname.Func,
	}, filters...)
}
