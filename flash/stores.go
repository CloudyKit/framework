package flash

import (
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

func (sess Session) Read(r *request.Context) (map[string]interface{}, error) {
	sessContext := session.GetSession(r.Global)
	if ii, has := sessContext.Lookup(sess.getKey()); has {
		sessContext.Unset(sess.getKey())
		return ii.(map[string]interface{}), nil
	}
	return nil, nil
}

func (sess Session) Save(r *request.Context, val map[string]interface{}) error {
	sessContext := session.GetSession(r.Global)
	sessContext.Set(sess.getKey(), val)
	return nil
}
