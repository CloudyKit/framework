package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/jet"
	"reflect"
)

var DefaultSet = jet.NewHTMLSet("./resources/views")

func init() {
	app.Default.Bootstrap(JetComponent{DefaultSet})
}

type JetComponent struct {
	Set *jet.Set
}

var JetContextType = reflect.TypeOf((*JetContext)(nil))

func GetJetContext(cdi *cdi.Global) *JetContext {
	c, _ := cdi.GetByType(JetContextType).(*JetContext)
	return c
}

func Render(global *cdi.Global, viewName string, c interface{}) {
	GetJetContext(global).Render(viewName, c)
}

var JetSetType = reflect.TypeOf((*jet.Set)(nil))

func GetJetSet(cdi *cdi.Global) *jet.Set {
	c, _ := cdi.GetByType(JetSetType).(*jet.Set)
	return c
}

func (p JetComponent) Bootstrap(a *app.App) {

	a.Global.MapType(JetSetType, p.Set)

	a.Global.MapType(JetContextType, func(cdi *cdi.Global) interface{} {
		cc := &JetContext{
			set:      p.Set,
			rcontext: request.GetContext(cdi),
		}
		for key, value := range cdi.GetByPtr(Globals(nil)).(Globals) {
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

func (s *JetContext) JetSet() *jet.Set {
	return s.set
}

func (c *JetContext) render(templateName string, context interface{}) error {
	t, err := c.set.GetTemplate(templateName)
	if err != nil {
		return err
	}
	return t.Execute(c.rcontext.Response, c.scope, context)
}

func (c *JetContext) Render(templateName string, context interface{}) {
	if err := c.render(templateName, context); err != nil {
		panic(err)
	}
}

func (c *JetContext) Execute(t *jet.Template, context interface{}) error {
	return t.Execute(c.rcontext.Response, c.scope, context)
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
