package events

import "github.com/CloudyKit/framework/cdi"

var sub = &Emitter{}

func Subscribe(global *cdi.Global, groupName string, handler interface{}) *Emitter {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Subscribe(groupName, handler)
		}
	}
	return sub.Subscribe(groupName, handler)
}

func NewEmitter() *Emitter {
	return sub.Inherit()
}

func Emit(global *cdi.Global, groupName, key string, c interface{}) (bool, error) {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Emit(groupName, key, c)
		}
	}
	return sub.Emit(groupName, key, c)
}

func Reset(global *cdi.Global, groupName string) bool {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Reset(groupName)
		}
	}
	return sub.Reset(groupName)
}
