package session

import (
	"github.com/CloudyKit/framework/context"
	"sync"
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

type Manager struct {
	Generator  IdGenerator
	Store      Store
	Serializer Serializer
	Duration   time.Duration
	gcEvery    time.Duration
	gcLastCall time.Time
	lock       sync.Mutex
}

// GCinvoke invokes garbage collection, will not update the timer
// this method should be called only when explicit necessary otherwise you should call GCinvokeifnecessary
// to only run gc periodically
func (manager *Manager) GCinvoke(c *context.Context, now time.Time, goroutine bool) {
	if goroutine {
		go manager.gcinvoke(c.Child(), now, goroutine)
	} else {
		manager.gcinvoke(c, now, goroutine)
	}
}
func (manager *Manager) gcinvoke(c *context.Context, now time.Time, goroutine bool) {
	if goroutine {
		defer c.Done()
	}
	manager.Store.Gc(c, now.Add(-manager.Duration))
}

// GCinvokeifnecessary checks the last time the garbage collector ran and if necessary
// runs the gc again and update the gcLastCall
func (manager *Manager) GCinvokeifnecessary(c *context.Context, goroutine bool) bool {
	now := time.Now()
	manager.lock.Lock()
	invokeGc := manager.gcLastCall.Add(manager.gcEvery).Before(now)
	if invokeGc {
		manager.gcLastCall = now
		manager.lock.Unlock()
		manager.GCinvoke(c, now, goroutine)
	} else {
		manager.lock.Unlock()
	}
	return invokeGc
}

// Open opens a stored session and unserialize into dst
func (manager *Manager) Open(ctx *context.Context, sessionName string, dst interface{}) error {
	reader, err := manager.Store.Reader(ctx, sessionName)
	if err == nil && reader != nil {
		err = manager.Serializer.Unserialize(dst, reader)
		reader.Close()
	} else if reader != nil {
		reader.Close()
	}
	return err
}

// Save saves the session
func (manager *Manager) Save(ctx *context.Context, sessionName string, src interface{}) error {
	writer, err := manager.Store.Writer(ctx, sessionName)
	if err == nil && writer != nil {
		err = manager.Serializer.Serialize(src, writer)
		writer.Close()
	} else if writer != nil {
		writer.Close()
	}
	return err
}

func (manager *Manager) Unregister(ctx *context.Context, sessionName string) error {
	return manager.Store.Remove(ctx, sessionName)
}
