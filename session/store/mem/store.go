package mem

import (
	"github.com/CloudyKit/framework/context"
	"io"
	"time"
	"bytes"
)

type Session struct {
	*bytes.Buffer
	lastUpdate time.Time
}

func (sess *Session) Close() error {
	return nil
}

type Store struct {
	sessions map[string]*Session
}

func New() Store {
	return Store{sessions:make(map[string]*Session)}
}

func (store Store) Reader(_ *context.Context, name string) (io.ReadCloser, error) {
	if reader, ok := store.sessions[name]; ok && reader != nil {
		return reader, nil
	}
	return nil, nil
}

func (store Store) Writer(_ *context.Context, name string) (writer io.WriteCloser, err error) {
	sess, ok := store.sessions[name]
	if sess == nil || ok == false {
		sess = &Session{Buffer:bytes.NewBuffer(nil), lastUpdate:time.Now()}
		store.sessions[name] = sess
	}else {
		sess.Reset()
		sess.lastUpdate = time.Now()
	}
	writer = sess
	return
}

func (store Store) Remove(_ *context.Context, name string) error {
	delete(store.sessions, name)
	return nil
}

func (store Store) Gc(_ *context.Context, before time.Time) {
	for id, sess := range store.sessions {
		if sess.lastUpdate.Before(before) {
			delete(store.sessions, id)
		}
	}
}
