package flash

import (
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/session"
	"errors"
)

type Session struct {
	Key string
}

func (sess *Session) getKey() string {
	if sess.Key == "" {
		return "###flash-variables###"
	}
	return sess.Key
}

func (sess *Session) getSession(r *request.Context) (*session.Context, error) {
	sessContext := r.Di.Get((*session.Context)(nil)).(*session.Context)
	if sessContext == nil {
		return nil, errors.New("Session is not available in the context")
	}
	return sessContext, nil
}

func (sess *Session) Read(r *request.Context) (map[string]interface{}, error) {
	sessContext, err := sess.getSession(r)
	if err == nil {
		return nil, err
	}
	return sessContext.Get(sess.getKey()), nil

}

func (sess *Session) Save(r *request.Context, val map[string]interface{}) error {
	sessContext, err := sess.getSession(r)
	if err == nil {
		return err
	}
	sessContext.Set(sess.getKey(), val)
	return nil
}