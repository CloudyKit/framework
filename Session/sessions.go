package Session

import (
	"encoding/gob"
	"github.com/CloudyKit/framework/App"
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"github.com/CloudyKit/framework/Session/Store/File"
	"net/http"
	"time"
)

var (
	DefaultManager       = New(time.Hour, time.Hour*2, File.Store{"./sessions"}, GobSerializer{}, RandGenerator{})
	DefaultCookieOptions = &CookieOptions{
		Name: "__gsid",
	}
)

func init() {

	gob.Register(SessionContext{})
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

type SessionContext map[string]interface{}
type FlashContext map[string]interface{}

// Flash adds a flash variable to the session
func (s SessionContext) Flash(k string, v interface{}) {
	s.wflashes(true)[k] = v
}

// HasFlash check if the k flash is set
func (s SessionContext) HasFlash(k string) (has bool) {
	_, has = s.Flashes()[k]
	return
}

// ReadFlash reads a flash variable from the session
func (s SessionContext) ReadFlash(k string) interface{} {
	return s.Flashes()[k]
}

// ReadWritingFlash reads a flash variable which is marked to be written
func (s SessionContext) ReadWritingFlash(k string) interface{} {
	return s.wflashes(false)[k]
}

func (s SessionContext) wflashes(set bool) FlashContext {
	flashes, _ := s[flashWritingKey].(FlashContext)
	if flashes == nil && set {
		flashes = make(FlashContext)
		s[flashWritingKey] = flashes
	}
	return flashes
}

// Flashes returns all flash variables stored from last request

func (s SessionContext) Flashes() FlashContext {
	flashes, _ := s[flashReadingKey].(FlashContext)
	return flashes
}

// Keep keep's the specified flash variables to the next request
func (s SessionContext) Keep(k ...string) {
	flashes := s.Flashes()
	for i := 0; i < len(k); i++ {
		if v, ok := flashes[k[i]]; ok {
			s.Flash(k[i], v)
		}
	}
}

// ReFlash will keep all stored flash variables to the next request
func (s SessionContext) ReFlash() {
	for k, v := range s.Flashes() {
		s.Flash(k, v)
	}
}

type sessionFinalizer struct {
	Id      string
	Manager *Manager
	Session SessionContext
}

func (sessionFinalizer *sessionFinalizer) Done() {

	flashVariables := sessionFinalizer.Session[flashWritingKey]

	if flashVariables != nil {
		sessionFinalizer.Session[flashReadingKey] = flashVariables
		delete(sessionFinalizer.Session, flashWritingKey)
	} else {
		delete(sessionFinalizer.Session, flashReadingKey)
		delete(sessionFinalizer.Session, flashWritingKey)
	}

	err := sessionFinalizer.Manager.Save(sessionFinalizer.Id, sessionFinalizer.Session)
	if err != nil {
		println(err.Error())
	}
	go sessionFinalizer.Manager.GcCheckAndRun()
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

	c.Set((SessionContext)(nil), func(c Di.Context) interface{} {
		rctx := c.Get((*Request.Context)(nil)).(*Request.Context)
		sess := &sessionFinalizer{Manager: sm, Session: make(SessionContext)}

		rCookie, _ := rctx.R.Cookie(so.Name)
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
