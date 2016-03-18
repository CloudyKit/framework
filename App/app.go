package App

import (
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Router"
)

var Default = New()

func New() *AppContext {
	newApp := &AppContext{Context: Di.New(), Router: Router.New(), Gen: make(map[string]string)}
	newApp.Put(newApp)
	return newApp
}

type AppContext struct {
	Di.Context
	Router *Router.Router
	Gen    map[string]string
}
