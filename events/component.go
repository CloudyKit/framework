//go:generate cdi -f=component.go -o=component.go
//cdi: *Emitter
package events

import "github.com/CloudyKit/framework/scope"

///cdi:generated
var EmitterType = scope.TypeOfElem((**Emitter)(nil))

func GetEmitter(c *scope.Variables) *Emitter {
	v, _ := c.GetByType(EmitterType).(*Emitter)
	return v
}

///cdi:generated
