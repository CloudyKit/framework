package common

import (
	"bytes"
	"fmt"
	"github.com/CloudyKit/framework/scope"
	"reflect"
)

var (
	NamedType = reflect.TypeOf((*Named)(nil)).Elem()
	URLerType = reflect.TypeOf((*URLer)(nil)).Elem()
)

type Named interface {
	Name() string
}

type URLer interface {
	URL(resource string, v ...interface{}) string
}

func GetURLer(cdi *scope.Variables) URLer {
	urler, _ := cdi.GetByType(URLerType).(URLer)
	return urler
}

func GenURL(cdi *scope.Variables, resource string, v ...interface{}) string {

	if cdi == nil {
		if len(v) == 0 {
			return resource
		}
		return fmt.Sprintf(resource, v...)
	}

	urLer := GetURLer(cdi)
	if urLer == nil {
		if len(v) == 0 {
			return resource
		}
		return fmt.Sprintf(resource, v...)
	}

	return urLer.URL(resource, v...)
}

type BaseURL func(...interface{}) string

func NewBaseURL(url string) BaseURL {
	return func(v ...interface{}) string {
		narguments := len(v)
		if narguments > 0 {
			buf := bytes.NewBuffer(nil)
			buf.WriteString(url)
			fmt.Fprintf(buf, "?%s=%s", v[0], v[1])
			for i := 2; i+1 < narguments; i += 2 {
				fmt.Fprintf(buf, "&%s=%s", v[i], v[i+1])
			}
			return buf.String()
		}
		return url
	}
}

func (fn BaseURL) New(urlPath string) BaseURL {
	return NewBaseURL(fn() + urlPath)
}

func (fn BaseURL) String() string {
	return fn()
}

// GenQS generates a url + query string
// ex: GenQS(nil,"http://google.com/")("q","cats") => Generates http://google.com/?q=cats
//     or use with GenQS("app.ProductController.ActionHandler","urlParam")("page",5)
func GenQS(global *scope.Variables, resource string, parameters ...interface{}) BaseURL {
	return NewBaseURL(GenURL(global, resource, parameters...))
}
