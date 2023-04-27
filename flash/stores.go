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

package flash

import (
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/session"
)

const defaultKey = "###flash-variables###"

type Session struct {
	Key string
}

func (sess Session) getKey() string {
	if sess.Key == "" {
		return defaultKey
	}
	return sess.Key
}

func (sess Session) Read(r *request.Context) (map[string]interface{}, error) {
	sessContext := session.GetSessionManager(r.Registry)
	if ii, has := sessContext.Lookup(sess.getKey()); has {
		sessContext.Unset(sess.getKey())
		return ii.(map[string]interface{}), nil
	}
	return nil, nil
}

func (sess Session) Save(r *request.Context, val map[string]interface{}) error {
	sessContext := session.GetSessionManager(r.Registry)
	sessContext.Set(sess.getKey(), val)
	return nil
}
