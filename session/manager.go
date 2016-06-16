package session

import (
	"github.com/CloudyKit/framework/cdi"
	"log"
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

	workchan chan *mJob
	donechan chan *mJob
}

// work manages concurrent read,write,remove,cleanup of sessions
// see *Manager.pushWork, *mWork.done
func (manager *Manager) work() {
	// pending operators in file
	var pending = make(map[string]struct {
		pending int
		end     chan struct{}
	})
	var gcTimer = time.NewTicker(manager.gcEvery)
	for {
		select {
		case n := <-gcTimer.C:
			// it's cleaning time
			go func(manager *Manager) {
				defer func() {
					// recovers from errors inside the garbage collection
					if err := recover(); err != nil {
						log.Println(err)
					}
				}()
				// invokes store GC
				manager.Store.GC(manager.Global, n.Add(-manager.Duration))
			}(manager)

		case w := <-manager.workchan:
			// receives a work
			// restore pending works state
			worksState := pending[w.name]
			// update counter of pending works
			worksState.pending++
			// check if this the first pending work
			if worksState.pending == 1 {
				// store the end channel in the state
				worksState.end = w.end
				// update state
				pending[w.name] = worksState
				// send continue signal with false, means continue directly and not wait ok signal
				// since this is the first job in the queue
				w.start <- false
			} else {
				// there is already running jobs in the queue
				// store the end channel of the last pushed job,
				// the end channel of the stored job is the ok channel of
				// the next job
				w.ok = worksState.end
				// store the end channel of the job being pushed
				worksState.end = w.end
				// update the state
				pending[w.name] = worksState
				// send continue signal with true, means continue and wait for ok signal
				w.start <- true
			}
		case done := <-manager.donechan:
			// receives a done signal with the address of the finished job
			work := pending[done.name]
			// decrement the counter
			work.pending--
			if work.pending == 0 {
				// pending 0 all jobs executed, execute cleanup
				delete(pending, done.name)
			} else {
				// update the state
				pending[done.name] = work
				// send ok signal will continue pending job
				done.end <- struct{}{}
			}
			_jobPool.Put(done)
		}
	}
}

var _jobPool = sync.Pool{New: func() interface{} {
	return &mJob{start: make(chan bool), end: make(chan struct{})}
}}

// pushWork push a work on sessionName
func (manager *Manager) pushWork(name string) *mJob {
	// creates the job
	work := _jobPool.Get().(*mJob)
	work.name = name
	// push the job into the manager queue
	manager.workchan <- work
	// await the start signal
	if <-work.start {
		// await an ok signal, wait my turn to execute my job
		<-work.ok
	}
	return work
}

// done notify the job is done
func (w *mJob) done(m *Manager) {
	m.donechan <- w
}

// Open opens a stored session and unserialize into dst
func (manager *Manager) Open(ctx *cdi.Global, sessionName string, dst interface{}) error {
	defer manager.pushWork(sessionName).done(manager)
	reader, err := manager.Store.Reader(ctx, sessionName, time.Now().Add(-manager.Duration))
	if err == nil && reader != nil {
		err = manager.Serializer.Unserialize(dst, reader)
		reader.Close()
	} else if reader != nil {
		reader.Close()
	}
	return err
}

// Save saves the session
func (manager *Manager) Save(ctx *cdi.Global, sessionName string, session interface{}) error {
	defer manager.pushWork(sessionName).done(manager)
	writer, err := manager.Store.Writer(ctx, sessionName)
	if err == nil && writer != nil {
		err = manager.Serializer.Serialize(session, writer)
		writer.Close()
	} else if writer != nil {
		writer.Close()
	}
	return err
}

// Sessions remove will remove the session and it's data
func (manager *Manager) Remove(ctx *cdi.Global, sessionName string) error {
	defer manager.pushWork(sessionName).done(manager)
	return manager.Store.Remove(ctx, sessionName)
}
