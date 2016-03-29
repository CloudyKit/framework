package errors

import (
	"github.com/CloudyKit/framework/context"
)

type Reporter interface {
	Report(di *context.Context, err error)
}

type Notifier struct {
	di       *context.Context
	reporter Reporter
}

func NewNotifier(di *context.Context, reporter Reporter) Notifier {
	return Notifier{di: di, reporter: reporter}
}

func (notifier Notifier) Provide(di *context.Context) interface{} {
	return NewNotifier(di, notifier.reporter)
}

func (notifier Notifier) NotifyPanic() {
	if err := recover(); err != nil {
		if err, isError := err.(error); isError {
			notifier.reporter.Report(notifier.di, err)
		}
	}
}

func (notifier Notifier) NotifyIfNotNil(errs ...error) Notifier {
	for i := 0; i < len(errs); i++ {
		err := errs[i]
		if err != nil {
			notifier.reporter.Report(notifier.di, err)
		}
	}
	return notifier
}
