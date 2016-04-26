package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/jet"
	"reflect"
)

var DefaultSet = jet.NewHTMLSet("./views")

func init() {
	app.Default.AddPlugin(viewPlugin{DefaultSet})
}

type viewPlugin struct {
	set *jet.Set
}

func (viewPlugin viewPlugin) Init(di *context.Context) {
	di.MapType((*Context)(nil), func(c *context.Context) interface{} {
		cc := &Context{set: viewPlugin.set, rcontext: c.Get((*request.Context)(nil)).(*request.Context)}
		for key, value := range c.Get(Globals(nil)).(Globals) {
			cc.With(key, value.Provide(c))
		}
		c.Map(cc) // remap
		return cc
	})
}

type Context struct {
	set      *jet.Set
	scope    jet.VarMap
	rcontext *request.Context
	global   Globals
}

func (c *Context) Render(templateName string, context interface{}) error {
	t, err := c.set.LoadTemplate(templateName, "")
	if err == nil {
		err = t.Execute(c.rcontext.Response, c.scope, context)
	}
	return err
}

func (c *Context) WithValue(name string, v reflect.Value) *Context {
	if c.scope == nil {
		c.scope = make(jet.VarMap)
	}
	c.scope[name] = v
	return c
}

func (c *Context) With(name string, v interface{}) *Context {
	if c.scope == nil {
		c.scope = make(jet.VarMap)
	}
	c.scope.Set(name, v)
	return c
}
