package reporters

import (
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/di"
	"log"
)

type LogReporter struct{}

func (logReporter LogReporter ) Report(di *di.Context, err error) {
	c := di.Get((*request.Context)(nil)).(*request.Context)
	log.Println(c.Name, err.Error())
}
