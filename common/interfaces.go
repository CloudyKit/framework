package common

import (
	"bytes"
	"fmt"
	"github.com/CloudyKit/framework/cdi"
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

func GetURLer(cdi *cdi.Global) URLer {
	urler, _ := cdi.GetByType(URLerType).(URLer)
	return urler
}

func GenURL(cdi *cdi.Global, resource string, v ...interface{}) string {
	urLer := GetURLer(cdi)
	if urLer == nil {
		return fmt.Sprintf(resource, v...)
	}
	return urLer.URL(resource, v...)
}

// GenQS generates a url + query string
// ex: GenQS("http://google.com/","q","cats") => Generates http://google.com/?q=cats
//     or use with GenQS(GenURL("app.ProductController.ActionHandler",""),"page",5)
func GenQS(url string, v ...interface{}) string {
	var _bytesBuffer [2083]byte
	buf := bytes.NewBuffer(_bytesBuffer[:])
	buf.WriteString(url)
	buf.WriteString("?")
	for i, v := range v {
		if i%2 == 0 {
			buf.WriteString("=")
			fmt.Fprint(buf, v)
			buf.WriteString("&")
		} else {
			fmt.Fprint(buf, v)
		}
	}
	return buf.String()
}
