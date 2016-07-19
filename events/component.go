//go:generate cdi -f=component.go -o=component.go
//cdi: *Emitter
package events

import "github.com/CloudyKit/framework/cdi"

///cdi:generated
var EmitterType = cdi.TypeOfElem((**Emitter)(nil))

func GetEmitter(c *cdi.Global) *Emitter {
	v, _ := c.GetByType(EmitterType).(*Emitter)
	return v
}

///cdi:generated
