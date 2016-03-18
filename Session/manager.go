package Session

import (
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
	Store      Store
	Serializer Serializer
	Generator  IdGenerator
	Duration   time.Duration
	gcEvery    time.Duration
	gcLastCall time.Time
	lock       sync.Mutex
}

// GcInvoke invokes garbage collection, will not update the timer
// this method should be called only when explicit necessary otherwise you should call GcCheckAndRun
// to only run gc periodically
func (session *Manager) GcInvoke(now time.Time) error {
	return session.Store.Gc(now.Add(-session.Duration))
}

// GcCheckAndRun checks the last time the garbage collector ran and if necessary
// runs the gc again and update the gcLastCall
func (session *Manager) GcCheckAndRun() (bool, error) {
	now := time.Now()
	var err error
	session.lock.Lock()
	invokeGc := session.gcLastCall.Add(session.gcEvery).Before(now)
	if invokeGc {
		session.gcLastCall = now
		session.lock.Unlock()
		session.GcInvoke(now)
	} else {
		session.lock.Unlock()
	}
	return invokeGc, err
}

// Open opens a stored session and unserialize into dst
func (session *Manager) Open(sessionName string, dst interface{}) error {
	reader, err := session.Store.Reader(sessionName)
	defer reader.Close()
	if err != nil || reader == nil {
		return err
	}
	return session.Serializer.Unserialize(dst, reader)
}

// Save saves the session
func (session *Manager) Save(sessionName string, src interface{}) error {
	writer, err := session.Store.Writer(sessionName)
	defer writer.Close()
	if err != nil || writer == nil {
		return err
	}
	return session.Serializer.Serialize(src, writer)
}
