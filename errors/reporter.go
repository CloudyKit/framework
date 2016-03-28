package errors

import (
	"github.com/CloudyKit/framework/di"
)

type Reporter interface {
	Report(di *di.Context, err error)
}

type Catcher struct {
	Reporter
}

func (catcher Catcher) CatchPanic(di *di.Context) {
	if err := recover(); err != nil {
		if err, isError := err.(error); isError {
			catcher.Report(di, err)
		}
	}
}

func (catcher Catcher) ReportIfNotNil(di *di.Context, err error) Catcher {
	if err != nil {
		catcher.Report(di, err)
	}
	return catcher
}
