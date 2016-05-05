package session

import (
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/request"
	"net/http"
	"reflect"
)

type Plugin struct {
	CookieOptions *CookieOptions
	Manager       *Manager
}

var SessionType = reflect.TypeOf((*Session)(nil))

func Get(cdi *cdi.DI) *Session {
	return cdi.Val4Type(SessionType).(*Session)
}

func (sp *Plugin) PluginInit(di *cdi.DI) {

	if sp.Manager == nil {
		sp.Manager = di.Get(sp.Manager).(*Manager)
	}

	if sp.CookieOptions == nil {
		sp.CookieOptions = di.Get(sp.CookieOptions).(*CookieOptions)
	}

	filters := di.Get((*request.Filters)(nil)).(*request.Filters)

	filters.AddFilter(func(c *request.Context, f request.Flow) {

		s := contextPool.Get().(*Session)
		s.Data = make(sessionData)

		c.Global.Map(s)

		if readedcookie, _ := c.Request.Cookie(sp.CookieOptions.Name); readedcookie == nil {
			s.id = sp.Manager.Generator.Generate("", sp.CookieOptions.Name)
		} else {
			s.id = sp.Manager.Generator.Generate(readedcookie.Value, sp.CookieOptions.Name)
			err := sp.Manager.Open(c.Global, readedcookie.Value, &s.Data)
			if err != nil {
				s.done(err)
			}
		}

		// resets the cookie
		http.SetCookie(c.Response, &http.Cookie{
			Name:     sp.CookieOptions.Name,
			Value:    s.id,
			Path:     sp.CookieOptions.Path,
			Domain:   sp.CookieOptions.Domain,
			Secure:   sp.CookieOptions.Secure,
			HttpOnly: sp.CookieOptions.HttpOnly,
			MaxAge:   sp.CookieOptions.MaxAge,
			Expires:  sp.CookieOptions.Expires,
		})

		f.Continue()

		err := sp.Manager.Save(c.Global, s.id, s.Data)
		s.done(err)

		sp.Manager.GCinvokeifnecessary(c.Global, true)
	})
}
