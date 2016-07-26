package session

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
	"log"
	"net/http"
	"reflect"
)

type Boot struct {
	CookieOptions *CookieOptions
	Manager       *Manager
}

var (
	SessionType = reflect.TypeOf((*Session)(nil))
)

func GetSession(cdi *scope.Variables) *Session {
	return cdi.GetByType(SessionType).(*Session)
}

func (sp *Boot) Bootstrap(a *app.App) {

	if sp.CookieOptions == nil {
		sp.CookieOptions = &CookieOptions{
			Name: "__gsid",
			Path: "/",
		}
	} else {
		if sp.CookieOptions.Path == "" {
			sp.CookieOptions.Path = "/"
		}
	}

	app.Get(a.Variables).AddFilter(func(c *request.Context, f request.Flow) {
		s := _sessionPool.Get().(*Session)
		s.data = make(sessionData)
		c.Variables.MapType(SessionType, s)
		if readedcookie, _ := c.Request.Cookie(sp.CookieOptions.Name); readedcookie == nil {
			s._id = sp.Manager.Generator.Generate("", sp.CookieOptions.Name)
		} else {
			s._id = sp.Manager.Generator.Generate(readedcookie.Value, sp.CookieOptions.Name)
			if s._id != readedcookie.Value {
				sp.Manager.Remove(c.Variables, readedcookie.Value)
			}
			err := sp.Manager.Open(c.Variables, readedcookie.Value, &s.data) //todo: use this error message here can be helpful
			if err != nil {
				log.Println("Session read err:", err.Error())
			}
		}

		// resets the cookie
		http.SetCookie(c.Response, &http.Cookie{
			Name:     sp.CookieOptions.Name,
			Value:    s._id,
			Path:     sp.CookieOptions.Path,
			Domain:   sp.CookieOptions.Domain,
			Secure:   sp.CookieOptions.Secure,
			HttpOnly: sp.CookieOptions.HttpOnly,
			MaxAge:   sp.CookieOptions.MaxAge,
			Expires:  sp.CookieOptions.Expires,
		})

		f.Continue()

		for sesskey, sessvalue := range s.data {
			of := reflect.ValueOf(sessvalue)
			switch of.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.UnsafePointer, reflect.Slice, reflect.Ptr:
				if of.IsNil() {
					delete(s.data, sesskey)
				}
			}
		}

		err := sp.Manager.Save(c.Variables, s._id, s.data)
		_sessionPool.Put(s)

		if err != nil {
			log.Println("Session write err:", err.Error())
		}
	})
}
