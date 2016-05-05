package flash

import (
	"errors"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/session"
)

const defaultKey = "###flash-variables###"

type Session struct {
	Key string
}

func (sess Session) getKey() string {
	if sess.Key == "" {
		return defaultKey
	}
	return sess.Key
}

func (sess Session) getSession(r *request.Context) (*session.Session, error) {
	sessContext := r.Global.Get((*session.Session)(nil)).(*session.Session)
	if sessContext == nil {
		return nil, errors.New("Session is not available in the context")
	}
	return sessContext, nil
}

func (sess Session) Read(r *request.Context) (map[string]interface{}, error) {
	sessContext, err := sess.getSession(r)
	if err != nil {
		return nil, err
	}
	if ii, has := sessContext.Lookup(sess.getKey()); has {
		sessContext.Unset(sess.getKey())
		return ii.(map[string]interface{}), nil
	}
	return nil, nil

}

func (sess Session) Save(r *request.Context, val map[string]interface{}) error {
	sessContext, err := sess.getSession(r)
	if err != nil {
		return err
	}
	sessContext.Set(sess.getKey(), val)
	return nil
}
