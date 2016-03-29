package reporters

import (
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
	"log"
)

type LogReporter struct{}

func (logReporter LogReporter) Report(di *context.Context, err error) {
	c, _ := di.Get((*request.Context)(nil)).(*request.Context)
	if c != nil {
		log.Println(c.Name, err.Error())
	} else {
		log.Println(err.Error())
	}
}
