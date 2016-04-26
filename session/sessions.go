package session

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/session/store/file"
	"sync"
	"time"
)

var (
	DefaultManager = New(time.Hour, time.Hour*2, file.New("./sessions"), GobSerializer{}, RandGenerator{})

	DefaultCookieOptions = &CookieOptions{
		Name: "__gsid",
	}

	contextPool = sync.Pool{
		New: func() interface{} {
			return &Session{}
		},
	}
)

func init() {
	gob.Register(sessionData(nil))
	app.Default.AddPlugin(&Plugin{Manager: DefaultManager, CookieOptions: DefaultCookieOptions})
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

type sessionData map[string]interface{}

type Session struct {
	id   string
	Data sessionData
}

func (c *Session) Id() string {
	return c.id
}

func (c *Session) Contains(key string) (contains bool) {
	_, contains = c.Data[key]
	return
}
func (c *Session) Get(name string) (value interface{}) {
	value, _ = c.Data[name]
	return
}

func (c *Session) Lookup(name string) (val interface{}, has bool) {
	val, has = c.Data[name]
	return
}

// Set sets a key in the session
func (c *Session) Set(name string, val interface{}) {
	c.Data[name] = val
}

// Unset deletes a key in the session map
func (c *Session) Unset(keys ...string) {
	for i := 0; i < len(keys); i++ {
		delete(c.Data, keys[i])
	}
}
