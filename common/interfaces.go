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

package common

import (
	"bytes"
	"fmt"
	"github.com/CloudyKit/framework/container"
	"reflect"
)

var (
	URLGenType = reflect.TypeOf((*URLGen)(nil)).Elem()
)

// URLGen url generator
type URLGen interface {
	URL(resource string, v ...interface{}) string // URL generates an URL
}

// Gets an URL generator from the scope
func GetURLGen(cdi *container.IoC) URLGen {
	urlGen, _ := cdi.LoadType(URLGenType).(URLGen)
	return urlGen
}

// GenURL generates an URL with the URLGen available in the scope
func GenURL(cdi *container.IoC, resource string, v ...interface{}) string {

	if cdi == nil {
		if len(v) == 0 {
			return resource
		}
		return fmt.Sprintf(resource, v...)
	}

	urLer := GetURLGen(cdi)
	if urLer == nil {
		if len(v) == 0 {
			return resource
		}
		return fmt.Sprintf(resource, v...)
	}

	return urLer.URL(resource, v...)
}

// BaseURL holds an base url, invoking this func will return the base url with query string,
// ex: NewBaseURL("/search")("q", "my search input","page",5) will result in /search?q=my search input&page=5
type BaseURL func(...interface{}) string

// NewBaseURL creates a new BaseURL, see type BaseURL func(...interface{}) string
func NewBaseURL(url string) BaseURL {
	return func(v ...interface{}) string {
		numOfArgs := len(v)
		if numOfArgs > 0 {
			buf := bytes.NewBuffer(nil)
			buf.WriteString(url)
			fmt.Fprintf(buf, "?%s=%s", v[0], v[1])
			for i := 2; i+1 < numOfArgs; i += 2 {
				fmt.Fprintf(buf, "&%s=%s", v[i], v[i+1])
			}
			return buf.String()
		}
		return url
	}
}

// New will append urlPath to the current base url
// ex: NewBaseURL("/api").New("/users") is will return an equivalent BaseURL{"/api/users"}
func (fn BaseURL) New(urlPath string) BaseURL {
	return NewBaseURL(fn() + urlPath)
}

func (fn BaseURL) String() string {
	return fn()
}

// GenQS generates a url + query string
// ex: GenQS(nil,"http://google.com/")("q","cats") => Generates http://google.com/?q=cats
//     or use with GenQS("app.ProductController.ActionHandler","urlParam")("page",5)
func GenQS(global *container.IoC, resource string, parameters ...interface{}) BaseURL {
	return NewBaseURL(GenURL(global, resource, parameters...))
}
