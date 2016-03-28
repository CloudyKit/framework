package session

import (
	"github.com/CloudyKit/framework/session/store/file"

	"encoding/gob"
	"github.com/CloudyKit/framework/view"
	"sync"
	"time"
	"github.com/CloudyKit/framework/app"
)

var (
	DefaultManager = New(time.Hour, time.Hour * 2, file.Store{"./sessions"}, GobSerializer{}, RandGenerator{})

	DefaultCookieOptions = &CookieOptions{
		Name: "__gsid",
	}

	contextPool = sync.Pool{
		New: func() interface{} {
			return &Context{
				data:make(sessionData),
			}
		},
	}
)

func init() {
	gob.Register(sessionData(nil))
	app.Default.AddPlugin(&Plugin{Manager:DefaultManager, CookieOptions:DefaultCookieOptions})
}

func New(gcEvery time.Duration, duration time.Duration, store Store, serializer Serializer, generator IdGenerator) *Manager {
	return &Manager{
		Generator:  generator,
		gcEvery:    gcEvery,
		Duration:   duration,
		Store:      store,
		Serializer: serializer,
	}
}

var _ = view.AvailableKey(view.DefaultManager, "Session", (*Context)(nil))

type sessionData map[string]interface{}

type Context struct {
	id   string
	data sessionData
}

func (c *Context) HasSession(key string) (isset bool) {
	_, isset = c.data[key]
	return
}
func (c *Context) Get(name string) (value interface{}) {
	value, _ = c.GetValue(name)
	return
}

func (c *Context) GetValue(name string) (val interface{}, has bool) {
	val, has = c.data[name]
	return
}

func (c *Context) Set(name string, val interface{}) {
	c.data[name] = val
}