package common

import (
	"bytes"
	"fmt"
	"github.com/CloudyKit/framework/scope"
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
func GetURLGen(cdi *scope.Variables) URLGen {
	urlGen, _ := cdi.GetByType(URLGenType).(URLGen)
	return urlGen
}

// GenURL generates an URL with the URLGen available in the scope
func GenURL(cdi *scope.Variables, resource string, v ...interface{}) string {

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
func GenQS(global *scope.Variables, resource string, parameters ...interface{}) BaseURL {
	return NewBaseURL(GenURL(global, resource, parameters...))
}
