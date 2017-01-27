// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package session

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/request"
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

func LoadSession(cdi *container.IoC) *Session {
	return cdi.LoadType(SessionType).(*Session)
}

func (component *Component) Handle(ctx *request.Context) {
	s := _sessionPool.Get().(*Session)
	s.data = make(sessionData)
	ctx.IoC.MapValue(SessionType, s)

	if readedcookie, _ := ctx.Request.Cookie(component.CookieOptions.Name); readedcookie == nil {
		s._id = component.Manager.Generator.Generate("", component.CookieOptions.Name)
	} else {
		s._id = component.Manager.Generator.Generate(readedcookie.Value, component.CookieOptions.Name)
		if s._id != readedcookie.Value {
			component.Manager.Remove(ctx.IoC, readedcookie.Value)
		}
		err := component.Manager.Open(ctx.IoC, readedcookie.Value, &s.data) //todo: use this error message here can be helpful
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

	err := component.Manager.Save(ctx.IoC, s._id, s.data)
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

	app.Get(a.IoC).AddMiddleHandlers(component)
}
