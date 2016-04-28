package flash

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
)

func init() {
	gob.Register((map[string]interface{})(nil))
	app.Default.AddPlugin(NewPlugin(Session{defaultKey}))
}

type Store interface {
	Read(*request.Context) (map[string]interface{}, error)
	Save(*request.Context, map[string]interface{}) error
}

type Flasher struct {
	writeData map[string]interface{}
	Data      map[string]interface{}
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

func NewPlugin(store Store) app.Plugin {
	plugin := new(flashPlugin)
	plugin.Store = store
	return plugin
}

type flashPlugin struct {
	Store
	Filters *request.Filters
}

func (plugin *flashPlugin) PluginInit(di *context.Context) {
	store := plugin.Store
	di.Inject(plugin)

	if store != nil {
		plugin.Store = store
	}

	plugin.Filters.AddFilter(func(c request.ContextChain) {
		readData, err := plugin.Read(c.Request)
		c.Request.Notifier.ErrNotify(err)
		cc := &Flasher{Data: readData}
		di.Map(cc)
		c.Next()
		if cc.writeData != nil {
			c.Request.Notifier.ErrNotify(plugin.Save(c.Request, cc.writeData))
		}
	})

}
