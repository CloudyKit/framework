package App

import (
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"github.com/CloudyKit/framework/Router"
	"github.com/CloudyKit/framework/Validator"

	"encoding/json"
	"fmt"
)

var (
	Default = New()
)

func init() {
	Di.Walkable(Controller{})
}

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

func (cc *Controller) ValidateRoute(c func(Validator.At)) Validator.Result {
	return Validator.Run(Validator.NewRouterValueProvider(cc.Rv), c)
}

func (cc *Controller) Validate(c func(Validator.At)) Validator.Result {
	return Validator.Run(Validator.NewRequestValueProvider(cc.Rq), c)
}

func (cc *Controller) ValidatePost(c func(Validator.At)) Validator.Result {
	return Validator.Run(Validator.NewURLValueProvider(cc.Rq.PostForm), c)
}

func (cc *Controller) DecodeJson(target interface{}) error {
	return json.NewDecoder(cc.Rq.Body).Decode(target)
}

func (cc *Controller) EncodeJson(from interface{}) error {
	return json.NewEncoder(cc.Rw).Encode(from)
}

func (cc *Controller) WriteString(txt string) (int, error) {
	return cc.Rw.Write([]byte(txt))
}
