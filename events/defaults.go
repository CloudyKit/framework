package events

import "github.com/CloudyKit/framework/scope"

var sub = &Emitter{}

func Subscribe(global *scope.Variables, groupName string, handler interface{}) *Emitter {
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

func Emit(global *scope.Variables, groupName, key string, c interface{}) (bool, error) {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Emit(groupName, key, c)
		}
	}
	return sub.Emit(groupName, key, c)
}

func Reset(global *scope.Variables, groupName string) bool {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Reset(groupName)
		}
	}
	return sub.Reset(groupName)
}
