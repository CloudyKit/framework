package request

import (
	"encoding/json"
	"github.com/monoculum/formam"
)

// BindGetForm decodes the request url values into target
func (ctx *Context) BindGetForm(target interface{}) error {
	if ctx.Request.Form == nil {
		ctx.Request.ParseForm()
	}
	return formam.Decode(ctx.Request.Form, target)
}

// BindJSON decodes request post data into target
func (ctx *Context) BindForm(target interface{}) error {
	if ctx.Request.PostForm == nil {
		ctx.Request.ParseForm()
	}
	return formam.Decode(ctx.Request.PostForm, target)
}

// BindJSON decodes request body as json into the target
func (ctx *Context) BindJSON(target interface{}) error {
	return json.NewDecoder(ctx.Request.Body).Decode(target)
}

// todo: add a generic bind func which will decode values conforming with
// the request content-type or a query string contentType containing the mime type.
// func (ctx *Context) Bind(target interface{}) error {
//	return ctx.BindForm(target)
// }
