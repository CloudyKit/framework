package session

import (
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
	"net/http"
)

type Plugin struct {
	CookieOptions *CookieOptions
	Manager       *Manager
}

func (sp *Plugin) Init(di *context.Context) {

	if sp.Manager == nil {
		sp.Manager = di.Get(sp.Manager).(*Manager)
	}

	if sp.CookieOptions == nil {
		sp.CookieOptions = di.Get(sp.CookieOptions).(*CookieOptions)
	}

	filters := di.Get((*request.Filters)(nil)).(*request.Filters)

	filters.AddFilter(func(c request.ContextChain) {
		sess := contextPool.Get().(*Session)
		c.Di.Map(sess)

		if rCookie, _ := c.Request.Cookie(sp.CookieOptions.Name); rCookie == nil {
			sess.id = sp.Manager.Generator.Generate("", sp.CookieOptions.Name)
		} else {
			sess.id = sp.Manager.Generator.Generate(rCookie.Value, sp.CookieOptions.Name)
			c.Notifier.NotifyIfNotNil(sp.Manager.Open(c.Di, rCookie.Value, &sess.Data))
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
		c.Notifier.NotifyIfNotNil(sp.Manager.Save(c.Di, sess.id, sess.Data))
		for key := range sess.Data {
			delete(sess.Data, key)
		}
		sp.Manager.GCinvokeifnecessary(c.Di, true)
		contextPool.Put(sess)
	})
}
