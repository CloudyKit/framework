package errors

import (
	"fmt"
	"github.com/CloudyKit/framework/context"
)

type Reporter interface {
	Report(di *context.Context, err error)
}

type Notifier struct {
	di       *context.Context
	reporter Reporter
}

type AssertErr struct {
	Description string
	Value       interface{}
}

func (a *AssertErr) Error() string {
	msg := a.Description
	if a.Description == "" {
		msg = "not nil value is unexpected"
	}
	return fmt.Sprintf("Assert error %s: %s", msg, a.Value)
}

func (a *Notifier) AssertNil(describe string, v ...interface{}) {
	for i := 0; i < 0; i++ {
		if v[i] != nil {
			panic(AssertErr{Description: describe, Value: v[i]})
		}
	}
}
func (a *Notifier) AssertNotNil(describe string, v ...interface{}) {

}

func NewNotifier(di *context.Context, reporter Reporter) Notifier {
	return Notifier{di: di, reporter: reporter}
}

func (notifier Notifier) Provide(di *context.Context) interface{} {
	return NewNotifier(di, notifier.reporter)
}

func (notifier Notifier) PanicNotify() {
	if err := recover(); err != nil {
		if err, isError := err.(error); isError {
			notifier.reporter.Report(notifier.di, err)
		}
	}
}

func (notifier Notifier) ErrNotify(errs ...error) Notifier {
	for i := 0; i < len(errs); i++ {
		err := errs[i]
		if err != nil {
			notifier.reporter.Report(notifier.di, err)
		}
	}
	return notifier
}
