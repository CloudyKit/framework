package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/jet"
	"reflect"
)

var DefaultSet = jet.NewHTMLSet("./views")

func init() {
	app.Default.Bootstrap(BootJet{DefaultSet})
}

type BootJet struct {
	Set *jet.Set
}

var JetContextType = reflect.TypeOf((*JetContext)(nil))

func GetJetContext(cdi *cdi.DI) *JetContext {
	return cdi.Val4Type(JetContextType).(*JetContext)
}

func (p BootJet) Bootstrap(a *app.App) {
	a.Global.MapType(JetContextType, func(cdi *cdi.DI) interface{} {
		cc := &JetContext{
			set:      p.Set,
			rcontext: request.Get(cdi),
		}
		for key, value := range cdi.Get(Globals(nil)).(Globals) {
			cc.With(key, value.Provide(cdi))
		}
		cdi.MapType(JetContextType, cc)
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
	if err != nil {
		panic(err)
	}
	err = t.Execute(c.rcontext.Response, c.scope, context)
	if err != nil {
		panic(err)
	}
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
