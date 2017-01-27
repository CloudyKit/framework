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

package request

import (
	"encoding/json"
)

// BindGetForm decodes the request url values into target
func (ctx *Context) BindGetForm(target interface{}) error {
	if ctx.Request.Form == nil {
		ctx.Request.ParseForm()
	}
	return formamDecoder(ctx.Request.Form, target)
}

// BindForm decodes request post data into target
func (ctx *Context) BindForm(target interface{}) error {
	if ctx.Request.PostForm == nil {
		ctx.Request.ParseForm()
	}
	return formamDecoder(ctx.Request.PostForm, target)
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
