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
func (session *Manager) GcInvoke(now time.Time) {
	session.Store.Gc(now.Add(-session.Duration))
}

// GcCheckAndRun checks the last time the garbage collector ran and if necessary
// runs the gc again and update the gcLastCall
func (session *Manager) GcCheckAndRun() bool {
	now := time.Now()
	session.lock.Lock()
	invokeGc := session.gcLastCall.Add(session.gcEvery).Before(now)
	if invokeGc {
		session.gcLastCall = now
		session.lock.Unlock()
		session.GcInvoke(now)
	} else {
		session.lock.Unlock()
	}
	return invokeGc
}

// Open opens a stored session and unserialize into dst
func (session *Manager) Open(sessionName string, dst interface{}) {
	reader := session.Store.Reader(sessionName)
	defer reader.Close()
	session.Serializer.Unserialize(dst, reader)
}

// Save saves the session
func (session *Manager) Save(sessionName string, src interface{}) {
	writer := session.Store.Writer(sessionName)
	defer writer.Close()
	session.Serializer.Serialize(src, writer)
}
