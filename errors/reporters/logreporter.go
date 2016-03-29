package reporters

import (
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
	"log"
)

type LogReporter struct{}

func (logReporter LogReporter) Report(di *context.Context, err error) {
	c := di.Get((*request.Context)(nil)).(*request.Context)
	log.Println(c.Name, err.Error())
}
