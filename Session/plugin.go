package session

import (
	"github.com/CloudyKit/framework/di"
	"github.com/CloudyKit/framework/request"
	"net/http"
)

type Plugin struct {
	CookieOptions *CookieOptions
	Manager       *Manager
}

func (sp *Plugin) Init(di *di.Context) {
	println("Session plugin is initializing")
	if sp.Manager == nil {
		sp.Manager = di.Get(sp.Manager).(*Manager)
	}

	if sp.CookieOptions == nil {
		sp.CookieOptions = di.Get(sp.CookieOptions).(*CookieOptions)
	}

	filters := di.Get((*request.Filters)(nil)).(*request.Filters)

	filters.AddFilter(func(c request.Channel) {
		sess := contextPool.Get().(*Context)
		c.Di.Map(sess)

		if rCookie, _ := c.Request.Cookie(sp.CookieOptions.Name); rCookie == nil {
			sess.id = sp.Manager.Generator.Generate("", sp.CookieOptions.Name)
		} else {
			sess.id = sp.Manager.Generator.Generate(rCookie.Value, sp.CookieOptions.Name)
			sp.Manager.Open(rCookie.Value, &sess.data)
		}
		// sets the cookie
		http.SetCookie(c.Response, &http.Cookie{
			Name:     sp.CookieOptions.Name,
			Value:    sess.id,
			Path:     sp.CookieOptions.Path,
			Domain:   sp.CookieOptions.Domain,
			Secure:   sp.CookieOptions.Secure,
			HttpOnly: sp.CookieOptions.HttpOnly,
			MaxAge:   sp.CookieOptions.MaxAge,
			Expires:  sp.CookieOptions.Expires,
		})

		c.Next()
		finalize(sp.Manager, sess)
		contextPool.Put(sess)
	})
}

func finalize(m *Manager, sess *Context) {
	m.Save(sess.id, sess.data)
	for key := range sess.data {
		delete(sess.data, key)
	}
	go m.GCinvokeifnecessary()
}