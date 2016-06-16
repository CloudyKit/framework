package session

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/session/store/file"
	"sync"
	"time"
)

var (
	DefaultManager = New(time.Hour, time.Hour*2, file.New("./resources/sessions"), GobSerializer{}, RandGenerator{})

	DefaultCookieOptions = &CookieOptions{
		Name: "__gsid",
	}

	_sessionPool = sync.Pool{
		New: func() interface{} {
			return &Session{}
		},
	}
)

func init() {
	gob.Register(SessionData(nil))
	app.Default.Bootstrap(&Boot{Manager: DefaultManager, CookieOptions: DefaultCookieOptions})
}

func New(gcEvery time.Duration, duration time.Duration, store Store, serializer Serializer, generator IdGenerator) *Manager {
	manager := &Manager{
		Generator:  generator,
		gcEvery:    gcEvery,
		Duration:   duration,
		Store:      store,
		Serializer: serializer,
	}

	manager.donechan = make(chan *mJob)
	manager.workchan = make(chan *mJob)

	manager.Store.GC(manager.Global, time.Now().Add(-manager.Duration))

	go manager.work()
	return manager
}

type SessionData map[string]interface{}

type Session struct {
	ID   string
	data SessionData
}

func (c *Session) Contains(key string) (contains bool) {
	_, contains = c.data[key]
	return
}

func (c *Session) Get(name string) (value interface{}) {
	value, _ = c.data[name]
	return
}

func (c *Session) Lookup(name string) (val interface{}, has bool) {
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
