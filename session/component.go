package session

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/scope"
	"log"
	"net/http"
	"reflect"
)

type Component struct {
	CookieOptions *CookieOptions
	Manager       *Manager
}

var (
	SessionType = reflect.TypeOf((*Session)(nil))
)

func GetSession(cdi *scope.Variables) *Session {
	return cdi.GetByType(SessionType).(*Session)
}

func (component *Component) Handle(ctx *request.Context) {
	s := _sessionPool.Get().(*Session)
	s.data = make(sessionData)
	ctx.Variables.MapType(SessionType, s)

	if readedcookie, _ := ctx.Request.Cookie(component.CookieOptions.Name); readedcookie == nil {
		s._id = component.Manager.Generator.Generate("", component.CookieOptions.Name)
	} else {
		s._id = component.Manager.Generator.Generate(readedcookie.Value, component.CookieOptions.Name)
		if s._id != readedcookie.Value {
			component.Manager.Remove(ctx.Variables, readedcookie.Value)
		}
		err := component.Manager.Open(ctx.Variables, readedcookie.Value, &s.data) //todo: use this error message here can be helpful
		if err != nil {
			log.Println("Session read err:", err.Error())
		}
	}

	// resets the cookie
	http.SetCookie(ctx.Response, &http.Cookie{
		Name:     component.CookieOptions.Name,
		Value:    s._id,
		Path:     component.CookieOptions.Path,
		Domain:   component.CookieOptions.Domain,
		Secure:   component.CookieOptions.Secure,
		HttpOnly: component.CookieOptions.HttpOnly,
		MaxAge:   component.CookieOptions.MaxAge,
		Expires:  component.CookieOptions.Expires,
	})

	ctx.Advance()

	for sesskey, sessvalue := range s.data {
		of := reflect.ValueOf(sessvalue)
		switch of.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.UnsafePointer, reflect.Slice, reflect.Ptr:
			if of.IsNil() {
				delete(s.data, sesskey)
			}
		}
	}

	err := component.Manager.Save(ctx.Variables, s._id, s.data)
	_sessionPool.Put(s)

	if err != nil {
		log.Println("Session write err:", err.Error())
	}
}

func (component *Component) Bootstrap(a *app.App) {

	if component.CookieOptions == nil {
		component.CookieOptions = &CookieOptions{
			Name: "__gsid",
			Path: "/",
		}
	} else {
		if component.CookieOptions.Path == "" {
			component.CookieOptions.Path = "/"
		}
	}

	app.Get(a.Variables).AddMiddleHandlers(component)
}
