package request

import (
	"encoding/json"
	"github.com/monoculum/formam"
)

// BindGetForm decodes the request url values into target
func (c *Context) BindGetForm(target interface{}) error {
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	return formam.Decode(c.Request.Form, target)
}

// BindJSON decodes request post data into target
func (c *Context) BindForm(target interface{}) error {
	if c.Request.PostForm == nil {
		c.Request.ParseForm()
	}
	return formam.Decode(c.Request.PostForm, target)
}

// BindJSON decodes request body as json into the target
func (c *Context) BindJSON(target interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(target)
}
