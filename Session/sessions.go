package Session

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/App"
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"github.com/CloudyKit/framework/Session/Store/File"
	"net/http"
	"sync"
	"time"
)

var (
	DefaultManager = New(time.Hour, time.Hour*2, File.Store{"./sessions"}, GobSerializer{}, RandGenerator{})

	DefaultCookieOptions = &CookieOptions{
		Name: "__gsid",
	}

	finalizersPool = sync.Pool{
		New: func() interface{} {
			return new(sessionFinalizer)
		},
	}
)

func init() {

	gob.Register(Context{})
	gob.Register(FlashContext{})
	//setup the default's
	SetupSessionProvider(App.Default.Context, DefaultManager, DefaultCookieOptions)
}

func New(gcEvery time.Duration, duration time.Duration, store Store, serializer Serializer, generator IdGenerator) *Manager {
	return &Manager{
		gcEvery:    gcEvery,
		Duration:   duration,
		Store:      store,
		Serializer: serializer,
		Generator:  generator,
	}
}

const flashWritingKey = "__wflash#session"
const flashReadingKey = "__rflash#session"

type Context map[string]interface{}
type FlashContext map[string]interface{}

// Flash adds a flash variable to the session
func (s Context) Flash(k string, v interface{}) {
	s.wflashes(true)[k] = v
}

// HasFlash check if the k flash is set
func (s Context) HasFlash(k string) (has bool) {
	_, has = s.Flashes()[k]
	return
}

// ReadFlash reads a flash variable from the session
func (s Context) ReadFlash(k string) interface{} {
	return s.Flashes()[k]
}

// ReadWritingFlash reads a flash variable which is marked to be written
func (s Context) ReadWritingFlash(k string) interface{} {
	return s.wflashes(false)[k]
}

func (s Context) wflashes(set bool) FlashContext {
	flashes, _ := s[flashWritingKey].(FlashContext)
	if flashes == nil && set {
		flashes = make(FlashContext)
		s[flashWritingKey] = flashes
	}
	return flashes
}

// Flashes returns all flash variables stored from last request

func (s Context) Flashes() FlashContext {
	flashes, _ := s[flashReadingKey].(FlashContext)
	return flashes
}

// Keep keep's the specified flash variables to the next request
func (s Context) Keep(k ...string) {
	flashes := s.Flashes()
	for i := 0; i < len(k); i++ {
		if v, ok := flashes[k[i]]; ok {
			s.Flash(k[i], v)
		}
	}
}

// ReFlash will keep all stored flash variables to the next request
func (s Context) Reflash() {
	for k, v := range s.Flashes() {
		s.Flash(k, v)
	}
}

type sessionFinalizer struct {
	Id      string
	Manager *Manager
	Session Context
}

func (sessFinalizerMem *sessionFinalizer) Finalize() {
	sessFinalizer := *sessFinalizerMem
	finalizersPool.Put(sessFinalizerMem)

	flashVariables := sessFinalizer.Session[flashWritingKey]
	if flashVariables != nil {
		sessFinalizer.Session[flashReadingKey] = flashVariables
		delete(sessFinalizer.Session, flashWritingKey)
	} else {
		delete(sessFinalizer.Session, flashReadingKey)
		delete(sessFinalizer.Session, flashWritingKey)
	}

	sessFinalizer.Manager.Save(sessFinalizer.Id, sessFinalizer.Session)
	go sessFinalizer.Manager.GcCheckAndRun()
}

// SetupSessionProvider this func will create an session loader into the context, the recommended
// way is to pass the app context
func SetupSessionProvider(c Di.Context, sm *Manager, so *CookieOptions) {

	if sm == nil {
		sm = DefaultManager
	}
	if so == nil {
		so = DefaultCookieOptions
	}

	c.Set((Context)(nil), func(c Di.Context) interface{} {
		rctx := c.Get((*Request.Context)(nil)).(*Request.Context)

		sess, _ := finalizersPool.Get().(*sessionFinalizer)
		sess.Manager = sm
		sess.Session = make(Context)

		rCookie, _ := rctx.Rq.Cookie(so.Name)
		if rCookie == nil {
			// generate new session
			sess.Id = sm.Generator.Generate(so.Name)
		} else {
			sess.Id = rCookie.Value
			sm.Open(rCookie.Value, &sess.Session)
		}

		http.SetCookie(rctx.Rw, &http.Cookie{
			Name:     so.Name,
			Value:    sess.Id,
			Path:     so.Path,
			Domain:   so.Domain,
			Secure:   so.Secure,
			HttpOnly: so.HttpOnly,
			MaxAge:   so.MaxAge,
			Expires:  so.Expires,
		})

		c.Put(sess)
		c.Put(sess.Session)
		return sess.Session
	})
}
