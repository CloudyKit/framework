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

package flash

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/ensure"
	"github.com/CloudyKit/framework/request"
	"reflect"
)

func init() {
	gob.Register((map[string]interface{})(nil))
	app.Default.Bootstrap(&Component{Session{defaultKey}})
}

type Store interface {
	Read(*request.Context) (map[string]interface{}, error)
	Save(*request.Context, map[string]interface{}) error
}

type Flasher struct {
	readed    bool
	writeData map[string]interface{}
	Data      map[string]interface{}
	store     Store
	context   *request.Context
}

func (c *Flasher) initWriter() {
	if c.writeData == nil {
		c.writeData = make(map[string]interface{})
	}
}
func (c *Flasher) initReader() {
	if c.readed == false {
		var err error
		c.Data, err = c.store.Read(c.context)
		ensure.NilErr(err)
		c.readed = true
	}
}

func (c *Flasher) CountMessages() int {
	return len(c.Data)
}

func (c *Flasher) Get(key string) interface{} {
	c.initReader()
	return c.Data[key]
}

func (c *Flasher) Contains(key string) (isset bool) {
	c.initReader()
	_, isset = c.Data[key]
	return
}

func (c *Flasher) Lookup(key string) (val interface{}, has bool) {
	c.initReader()
	val, has = c.Data[key]
	return
}

func (c *Flasher) Set(key string, val interface{}) {
	c.initWriter()
	c.writeData[key] = val
}

func (c *Flasher) Reflash(keys ...string) {
	c.initWriter()
	for _, key := range keys {
		if val, has := c.Data[key]; has {
			c.writeData[key] = val
		}
	}
}

type Component struct {
	Store
}

var FlasherType = reflect.TypeOf((*Flasher)(nil))

func GetFlasher(cdi *container.Registry) *Flasher {
	return cdi.LoadType(FlasherType).(*Flasher)
}

type flasher Flasher

func (f *flasher) dispose() {
	if len(f.writeData) > 0 {
		err := f.store.Save(f.context, f.writeData)
		ensure.NilErr(err)
	}
}

func (f *flasher) Provide(cdi *container.Registry) interface{} {
	return (*Flasher)(f)
}

func (component *Component) Handle(ctx *request.Context) {

	// allocates the flasher|flasherProvider
	flasher := &flasher{store: component.Store, context: ctx}

	// maps flasher in the request scope
	ctx.Registry.MapProvider(FlasherType, flasher)

	// advance with the request
	ctx.Next()

	// finalize the request
	flasher.dispose()
}

func (component *Component) Bootstrap(a *app.Kernel) {
	a.Root().BindFilterHandlers(component)
}
