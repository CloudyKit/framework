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
			return &Session{
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

var _ = view.AvailableKey(view.DefaultManager, "session", (*Session)(nil))

type sessionData map[string]interface{}

type Session struct {
	id   string
	data sessionData
}

func (c *Session) Id() string {
	return c.id
}

func (c *Session) IsSet(key string) (isset bool) {
	_, isset = c.data[key]
	return
}
func (c *Session) Get(name string) (value interface{}) {
	value, _ = c.GetValue(name)
	return
}

func (c *Session) GetValue(name string) (val interface{}, has bool) {
	val, has = c.data[name]
	return
}

// Set sets a key in the session
func (c *Session) Set(name string, val interface{}) {
	c.data[name] = val
}

// Unset deletes a key in the session map
func (c *Session) Unset(keys ...string) {
	for i := 0; i < len(keys); i++ {
		delete(c.data, keys[i])
	}
}