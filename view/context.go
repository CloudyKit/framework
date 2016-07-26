package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
	"github.com/CloudyKit/jet"
	"reflect"
)

var DefaultSet = jet.NewHTMLSet("./resources/views")

func init() {
	app.Default.Bootstrap(Component{DefaultSet})
}

type Component struct {
	Set *jet.Set
}

var RendererType = reflect.TypeOf((*Renderer)(nil))

func GetRenderer(cdi *scope.Variables) *Renderer {
	c, _ := cdi.GetByType(RendererType).(*Renderer)
	return c
}

func Render(global *scope.Variables, viewName string, c interface{}) {
	GetRenderer(global).Render(viewName, c)
}

var JetSetType = reflect.TypeOf((*jet.Set)(nil))

func GetJetSet(cdi *scope.Variables) *jet.Set {
	c, _ := cdi.GetByType(JetSetType).(*jet.Set)
	return c
}

func (p Component) Bootstrap(a *app.App) {

	a.Variables.MapType(JetSetType, p.Set)

	a.Variables.MapType(RendererType, func(cdi *scope.Variables) interface{} {
		cc := &Renderer{
			set:      p.Set,
			rcontext: request.GetContext(cdi),
		}
		for key, value := range cdi.GetByPtr(Globals(nil)).(Globals) {
			cc.With(key, value.Provide(cdi))
		}
		cdi.MapType(RendererType, cc)
		return cc
	})
}

type Renderer struct {
	set      *jet.Set
	scope    jet.VarMap
	rcontext *request.Context
	global   Globals
}

func (s *Renderer) JetSet() *jet.Set {
	return s.set
}

func (c *Renderer) render(templateName string, context interface{}) error {
	t, err := c.set.GetTemplate(templateName)
	if err != nil {
		return err
	}
	return t.Execute(c.rcontext.Response, c.scope, context)
}

func (c *Renderer) Render(templateName string, context interface{}) {
	if err := c.render(templateName, context); err != nil {
		panic(err)
	}
}

func (c *Renderer) Execute(t *jet.Template, context interface{}) error {
	return t.Execute(c.rcontext.Response, c.scope, context)
}

func (c *Renderer) WithValue(name string, v reflect.Value) *Renderer {
	if c.scope == nil {
		c.scope = make(jet.VarMap)
	}
	c.scope[name] = v
	return c
}

func (c *Renderer) With(name string, v interface{}) *Renderer {
	if c.scope == nil {
		c.scope = make(jet.VarMap)
	}
	c.scope.Set(name, v)
	return c
}
