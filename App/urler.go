package app

import (
	"fmt"
	"github.com/CloudyKit/framework/common"
)

type urlGen map[string]string
type ctlGen struct {
	urlGen urlGen
	id     string
	Parent common.URLer
}

func (urler *ctlGen) URL(dst string, v ...interface{}) string {

	if dst, ok := urler.urlGen[urler.id+dst]; ok {
		return fmt.Sprintf(dst, v...)
	}

	if dst, ok := urler.urlGen[dst]; ok {
		return fmt.Sprintf(dst, v...)
	}

	if urler.Parent != nil {
		return urler.Parent.URL(dst, v...)
	}

	return fmt.Sprintf(dst, v...)
}

func (urler urlGen) URL(dst string, v ...interface{}) string {
	if dst, ok := urler[dst]; ok {
		return fmt.Sprintf(dst, v...)
	}
	return fmt.Sprintf(dst, v...)
}
