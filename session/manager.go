package session

import (
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/insync"
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
	Global     *cdi.Global
	Generator  IdGenerator
	Store      Store
	Serializer Serializer
	Duration   time.Duration
	gcEvery    time.Duration
	kMX        insync.KMutex
}

func (manager *Manager) gcgoroutine() {
	for n := range time.NewTicker(manager.gcEvery).C {
		manager.Store.GC(manager.Global, n.Add(-manager.Duration))
	}
}

//Open load stored session and un serialize the stored data into dst
func (manager *Manager) Open(ctx *cdi.Global, sessionName string, dst interface{}) error {
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

//Save save the session
func (manager *Manager) Save(ctx *cdi.Global, sessionName string, session interface{}) error {
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

//Remove remove the session
func (manager *Manager) Remove(ctx *cdi.Global, sessionName string) error {
	defer manager.kMX.Lock(sessionName).Unlock()
	return manager.Store.Remove(ctx, sessionName)
}
