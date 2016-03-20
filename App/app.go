package App

import (
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"github.com/CloudyKit/framework/Router"

	"fmt"
)

var (
	Default = New()
	_       = Di.Walkable(Controller{})
)

func New() *Application {
	newApp := &Application{Context: Di.New(), Router: Router.New(), Gen: make(map[string]string)}
	newApp.Put(newApp)
	return newApp
}

type Application struct {
	Di.Context
	Router *Router.Router
	Gen    map[string]string
}

type Controller struct {
	Application *Application
	*Request.Context
}

func (cc *Controller) GenURL(dst string, v ...interface{}) string {

	if dst, ok := cc.Application.Gen[dst]; ok {
		return fmt.Sprintf(dst, v...)
	}

	if dst, ok := cc.Application.Gen[cc.Context.Id+"."+dst]; ok {
		return fmt.Sprintf(dst, v...)
	}

	return fmt.Sprintf(dst, v...)
}
