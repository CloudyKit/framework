package flash

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/assert"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
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
		assert.NilErr(err)
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

func GetFlasher(cdi *scope.Variables) *Flasher {
	return cdi.GetByType(FlasherType).(*Flasher)
}

type flasher Flasher

func (f *flasher) finalize() {
	if len(f.writeData) > 0 {
		err := f.store.Save(f.context, f.writeData)
		assert.NilErr(err)
	}
}

func (f *flasher) Provide(cdi *scope.Variables) interface{} {
	return (*Flasher)(f)
}

func (component *Component) Handle(ctx *request.Context) {

	// allocates the flasher|flasherProvider
	flasher := &flasher{store: component.Store, context: ctx}

	// maps flasher in the request scope
	ctx.Variables.MapType(FlasherType, flasher)

	// advance with the request
	ctx.Advance()

	// finalize the request
	flasher.finalize()
}

func (component *Component) Bootstrap(a *app.App) {
	a.Root().AddMiddleHandlers(component)
}
