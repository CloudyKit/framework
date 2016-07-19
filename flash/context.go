package flash

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/assert"
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/request"
	"reflect"
)

func init() {
	gob.Register((map[string]interface{})(nil))
	app.Default.Bootstrap(&Boot{Session{defaultKey}})
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

type Boot struct {
	Store
}

var FlasherType = reflect.TypeOf((*Flasher)(nil))

func GetFlasher(cdi *cdi.Global) *Flasher {
	return cdi.GetByType(FlasherType).(*Flasher)
}

type flasher Flasher

func (f *flasher) finalize() {
	if len(f.writeData) > 0 {
		err := f.store.Save(f.context, f.writeData)
		assert.NilErr(err)
	}
}

func (f *flasher) Provide(cdi *cdi.Global) interface{} {
	return (*Flasher)(f)
}

func (p *Boot) Bootstrap(a *app.App) {
	a.Root().AddFilter(func(c *request.Context, f request.Flow) {
		cc := &flasher{store: p.Store, context: c}
		defer cc.finalize()
		c.Global.MapType(FlasherType, cc)
		f.Continue()
	})
}
