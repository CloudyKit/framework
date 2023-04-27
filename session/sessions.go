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

package session

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/concurrent"
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
	gob.Register(sessionData(nil))
	app.Default.Bootstrap(&Bundle{Manager: DefaultManager, CookieOptions: DefaultCookieOptions})
}

// New returns a new session
func New(gcEvery time.Duration, duration time.Duration, store Store, serializer Serializer, generator IdGenerator) *Manager {

	manager := &Manager{
		Generator:  generator,
		gcEvery:    gcEvery,
		Duration:   duration,
		Store:      store,
		Serializer: serializer,
		kMX:        concurrent.NewKeyLocker(),
	}

	//collect expired sessions
	manager.Store.GC(manager.Global, time.Now().Add(-manager.Duration))

	//starts the garbage collect goroutine
	go manager.gcgoroutine()

	return manager
}

type sessionData map[string]interface{}

type Session struct {
	_id  string
	data sessionData
}

func (c *Session) ID() string {
	return c._id
}

func (c *Session) Contains(key string) (contains bool) {
	_, contains = c.data[key]
	return
}

func (c *Session) Get(name string) (value interface{}) {
	value, _ = c.data[name]
	return
}

func (c *Session) Lookup(name string) (data interface{}, found bool) {
	data, found = c.data[name]
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
