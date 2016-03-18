package Controller

import (
	"github.com/CloudyKit/framework/App"
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"

	"fmt"
	"github.com/CloudyKit/framework/Log"
)

var _ = Di.Walkable(Context{})

type Context struct {
	*Request.Context

	Application *App.AppContext
	*Log.LogContext
}

func (cc *Context) GenURL(dst string, v ...interface{}) string {

	if dst, ok := cc.Application.Gen[dst]; ok {
		return fmt.Sprintf(dst, v...)
	}

	if dst, ok := cc.Application.Gen[cc.Name+"."+dst]; ok {
		return fmt.Sprintf(dst, v...)
	}

	return fmt.Sprintf(dst, v...)
}
