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
	"github.com/CloudyKit/framework/concurrent"
	"github.com/CloudyKit/framework/container"
	"time"
)

type CookieOptions struct {
	Name   string
	Path   string
	Domain string

	MaxAge int

	Expires time.Time

	Secure   bool
	HttpOnly bool
}

type mJob struct {
	name  string
	ok    chan struct{}
	start chan bool
	end   chan struct{}
}

type Manager struct {
	Global     *container.Registry
	Generator  IdGenerator
	Store      Store
	Serializer Serializer
	Duration   time.Duration
	gcEvery    time.Duration
	kMX        *concurrent.KeyLocker
}

func (manager *Manager) gcgoroutine() {
	for n := range time.NewTicker(manager.gcEvery).C {
		manager.Store.GC(manager.Global, n.Add(-manager.Duration))
	}
}

// Open load stored session and un serialize the stored data into dst
func (manager *Manager) Open(ctx *container.Registry, sessionName string, dst interface{}) error {
	defer manager.kMX.Lock(sessionName).Unlock()
	reader, err := manager.Store.Reader(ctx, sessionName, time.Now().Add(-manager.Duration))
	if err == nil && reader != nil {
		err = manager.Serializer.Unserialize(dst, reader)
		reader.Close()
	} else if reader != nil {
		reader.Close()
	}
	return err
}

// Save save the session
func (manager *Manager) Save(ctx *container.Registry, sessionName string, session interface{}) error {
	defer manager.kMX.Lock(sessionName).Unlock()
	writer, err := manager.Store.Writer(ctx, sessionName)
	if err == nil && writer != nil {
		err = manager.Serializer.Serialize(session, writer)
		writer.Close()
	} else if writer != nil {
		writer.Close()
	}
	return err
}

// Remove remove the session
func (manager *Manager) Remove(ctx *container.Registry, sessionName string) error {
	defer manager.kMX.Lock(sessionName).Unlock()
	return manager.Store.Remove(ctx, sessionName)
}
