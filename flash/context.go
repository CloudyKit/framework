package flash

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
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
	writeData map[string]interface{}
	Data      map[string]interface{}
	store     Store
	context   *request.Context
}

func (c *Flasher) CountMessages() int {
	return len(c.Data)
}

func (c *Flasher) Reflash(keys ...string) {
	for i := 0; i < len(keys); i++ {
		if val, has := c.Data[keys[i]]; has {
			c.writeData[keys[i]] = val
		}
	}

}

func (c *Flasher) Get(key string) interface{} {
	return c.Data[key]
}

func (c *Flasher) Contains(key string) (isset bool) {
	_, isset = c.Data[key]
	return
}

func (c *Flasher) Lookup(key string) (val interface{}, has bool) {
	val, has = c.Data[key]
	return
}

func (c *Flasher) Set(key string, val interface{}) {
	if c.writeData == nil {
		c.writeData = make(map[string]interface{})
	}
	c.writeData[key] = val
}

type Boot struct {
	Store
}

var FlasherType = reflect.TypeOf((*Flasher)(nil))

func GetFlasher(cdi *cdi.Global) *Flasher {
	return cdi.Val4Type(FlasherType).(*Flasher)
}

type flasher Flasher

func (f *flasher) finalize() {
	if f.writeData != nil {
		err := f.store.Save(f.context, f.writeData)
		if err != nil {
			panic(err)
		}
	}
}

func (f *flasher) Provide(cdi *cdi.Global) interface{} {
	return (*Flasher)(f)
}

func (p *Boot) Bootstrap(a *app.App) {

	a.Root().AddFilter(func(c *request.Context, f request.Flow) {
		readData, err := p.Read(c)
		if err != nil {
			panic(err)
		}
		cc := &flasher{Data: readData, store: p.Store, context: c}
		c.Global.MapType(FlasherType, cc)
		f.Continue()
		defer cc.finalize()
	})

}
