// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package events

import "github.com/CloudyKit/framework/container"

var sub = &Emitter{}

func Subscribe(global *container.IoC, groupName string, handler interface{}) *Emitter {
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

func Emit(global *container.IoC, groupName, key string, c interface{}) (bool, error) {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Emit(groupName, key, c)
		}
	}
	return sub.Emit(groupName, key, c)
}

func Reset(global *container.IoC, groupName string) bool {
	if global != nil {
		if sub := GetEmitter(global); sub != nil {
			return sub.Reset(groupName)
		}
	}
	return sub.Reset(groupName)
}
