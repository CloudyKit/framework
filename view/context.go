// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package view

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/request"
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

func GetRenderer(cdi *container.IoC) *Renderer {
	c, _ := cdi.LoadType(RendererType).(*Renderer)
	return c
}

func Render(global *container.IoC, viewName string, c interface{}) {
	GetRenderer(global).Render(viewName, c)
}

var JetSetType = reflect.TypeOf((*jet.Set)(nil))
var globalType = reflect.TypeOf(Globals(nil))

func GetJetSet(cdi *container.IoC) *jet.Set {
	c, _ := cdi.LoadType(JetSetType).(*jet.Set)
	return c
}

func (p Component) Bootstrap(a *app.App) {

	a.IoC.MapValue(JetSetType, p.Set)

	a.IoC.MapProviderFunc(RendererType, func(cdi *container.IoC) interface{} {
		cc := &Renderer{
			set:      p.Set,
			rcontext: request.GetContext(cdi),
		}
		for key, value := range cdi.LoadType(globalType).(Globals) {
			cc.With(key, value.Provide(cdi))
		}
		cdi.MapValue(RendererType, cc)
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
