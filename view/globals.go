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
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/container"
	"reflect"
)

type provider interface {
	Provide(c *container.IoC) interface{}
}

func init() {
	var defaultGlobal = Globals{}
	app.Default.IoC.Map(defaultGlobal)
	GlobalInjectName(app.Default.IoC, "link", common.URLGenType)
}

type valueProvider struct {
	v interface{}
}

func (v valueProvider) Provide(c *container.IoC) interface{} {
	return v.v
}

type contextProvider struct {
	typeof reflect.Type
}

func (v contextProvider) Provide(c *container.IoC) interface{} {
	return c.LoadType(v.typeof)
}

type Globals map[string]provider

func GlobalInjectName(ci *container.IoC, name string, typ reflect.Type) error {
	return globalNameProvider(ci, name, contextProvider{typ})
}

func GlobalName(ci *container.IoC, name string, v interface{}) error {
	return globalNameProvider(ci, name, valueProvider{v})
}

func globalNameProvider(ci *container.IoC, name string, v provider) error {
	globals := ci.LoadType(globalType).(Globals)
	globals[name] = v
	return nil
}
