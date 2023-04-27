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
	"github.com/CloudyKit/jet/v6"
	"reflect"
)

var DefaultSet = jet.NewSet(jet.NewOSFileSystemLoader("./resources/views"))

func init() {
	app.Default.Bootstrap(Component{DefaultSet})
}

type Component struct {
	Set *jet.Set
}

var RendererType = reflect.TypeOf((*Renderer)(nil))

func GetRenderer(cdi *container.Registry) *Renderer {
	c, _ := cdi.LoadType(RendererType).(*Renderer)
	return c
}

func Render(global *container.Registry, viewName string, c interface{}) {
	GetRenderer(global).Render(viewName, c)
}

var JetSetType = reflect.TypeOf((*jet.Set)(nil))
var globalType = reflect.TypeOf(Globals(nil))

func GetJetSet(cdi *container.Registry) *jet.Set {
	c, _ := cdi.LoadType(JetSetType).(*jet.Set)
	return c
}

func (component Component) Bootstrap(a *app.Kernel) {

	a.Registry.WithTypeAndValue(JetSetType, component.Set)

	a.Registry.WithTypeAndProviderFunc(RendererType, func(cdi *container.Registry) interface{} {
		cc := &Renderer{
			set:      component.Set,
			rcontext: request.GetContext(cdi),
		}
		for key, value := range cdi.LoadType(globalType).(Globals) {
			cc.With(key, value.Provide(cdi))
		}
		cdi.WithTypeAndValue(RendererType, cc)
		return cc
	})
}

type Renderer struct {
	set      *jet.Set
	scope    jet.VarMap
	rcontext *request.Context
	global   Globals
}

func (renderer *Renderer) JetSet() *jet.Set {
	return renderer.set
}

func (renderer *Renderer) render(templateName string, context interface{}) error {
	t, err := renderer.set.GetTemplate(templateName)
	if err != nil {
		return err
	}
	return t.Execute(renderer.rcontext.Response, renderer.scope, context)
}

func (renderer *Renderer) Render(templateName string, context interface{}) {
	if err := renderer.render(templateName, context); err != nil {
		panic(err)
	}
}

func (renderer *Renderer) Execute(t *jet.Template, context interface{}) error {
	return t.Execute(renderer.rcontext.Response, renderer.scope, context)
}

func (renderer *Renderer) WithValue(name string, v reflect.Value) *Renderer {
	if renderer.scope == nil {
		renderer.scope = make(jet.VarMap)
	}
	renderer.scope[name] = v
	return renderer
}

func (renderer *Renderer) With(name string, v interface{}) *Renderer {
	if renderer.scope == nil {
		renderer.scope = make(jet.VarMap)
	}
	renderer.scope.Set(name, v)
	return renderer
}
