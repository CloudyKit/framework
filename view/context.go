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

func (viewPlugin viewPlugin) PluginInit(di *context.Context) {
	di.MapType((*JetContext)(nil), func(c *context.Context) interface{} {
		cc := &JetContext{
			set:      viewPlugin.set,
			rcontext: c.Get((*request.Context)(nil)).(*request.Context),
		}
		for key, value := range c.Get(Globals(nil)).(Globals) {
			cc.With(key, value.Provide(c))
		}
		c.Map(cc)
		return cc
	})
}

type JetContext struct {
	set      *jet.Set
	scope    jet.VarMap
	rcontext *request.Context
	global   Globals
}

func (c *JetContext) Render(templateName string, context interface{}) *JetContext {
	t, err := c.set.LoadTemplate(templateName, "")
	c.rcontext.Notifier.AssertNil("unexpected error loading template", err)
	err = t.Execute(c.rcontext.Response, c.scope, context)
	c.rcontext.Notifier.AssertNil("unexpected error executing template", err)
	return c
}

func (c *JetContext) WithValue(name string, v reflect.Value) *JetContext {
	if c.scope == nil {
		c.scope = make(jet.VarMap)
	}
	c.scope[name] = v
	return c
}

func (c *JetContext) With(name string, v interface{}) *JetContext {
	if c.scope == nil {
		c.scope = make(jet.VarMap)
	}
	c.scope.Set(name, v)
	return c
}
