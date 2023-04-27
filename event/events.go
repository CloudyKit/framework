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

package event

import (
	"errors"
	"fmt"
	"github.com/CloudyKit/framework/container"
	"reflect"
	"strings"
	"sync"
)

var eventPayloadType = reflect.TypeOf((*Payload)(nil)).Elem()

var (
	_ = Payload(&Event{})
)

type Event struct {
	eventName   string
	err         error
	canceled    bool
	unsubscribe bool
	registry    *container.Registry
}

type Payload interface {
	init(registry *container.Registry, eventName string)
	error() error
	unsubscribed() bool
	WasCanceled() bool
	Registry() *container.Registry
	EventName() string
	Cancel()
	CancelWithError(err error)
	CancelWithErrorf(format string, v ...interface{})
}

func (e *Event) init(registry *container.Registry, eventName string) {
	*e = Event{registry: registry, eventName: eventName}
}

func (e *Event) WasCanceled() bool {
	return e.canceled
}

func (e *Event) Registry() *container.Registry {
	return e.registry
}

func (e *Event) EventName() string {
	return e.eventName
}

func (e *Event) UnSubscribe() {
	e.unsubscribe = true
}

func (e *Event) unsubscribed() bool {
	return e.unsubscribe
}

func (e *Event) error() error {
	return e.err
}

func (e *Event) Cancel() {
	e.canceled = true
}

func (e *Event) CancelWithError(err error) {
	e.err = err
	e.canceled = true
}

func (e *Event) CancelWithErrorf(format string, v ...interface{}) {
	e.CancelWithError(fmt.Errorf(format, v...))
}

type subscriptionGroups struct {
	mutex        sync.RWMutex
	name         string
	handlers     []interface{}
	topHandler   int64
	runningEmits int64
}

func (group *subscriptionGroups) clearNilHandlers() {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	for i := 0; i < len(group.handlers); i++ {
		if group.handlers[i] == nil {
			group.handlers = append(group.handlers[0:i], group.handlers[i+1:len(group.handlers)-1]...)
		}
	}
}

type Dispatcher struct {
	parent        *Dispatcher
	mx            sync.RWMutex
	subscriptions []subscriptionGroups
}

func (dispatcher *Dispatcher) Inherit() *Dispatcher {
	return &Dispatcher{parent: dispatcher}
}

var err = errors.New("unexpected handler signature: func(*Event,ContextType) is expected")

func validateHandler(h interface{}) error {
	t := reflect.TypeOf(h)

	if t.Kind() != reflect.Func || t.NumOut() != 0 {
		return err
	}

	numIn := t.NumIn()
	if numIn > 1 || !t.In(0).AssignableTo(eventPayloadType) || t.IsVariadic() {
		return err
	}

	return nil
}

func (dispatcher *Dispatcher) Reset(groupName string) bool {
	dispatcher.assert()

	subsgroup, ok := dispatcher.group(groupName)
	if ok {
		subsgroup.mutex.Lock()
		subsgroup.topHandler = -1
		subsgroup.handlers = nil
		subsgroup.mutex.Unlock()
	}
	return ok
}

func (dispatcher *Dispatcher) subscribe(groupName string, handler interface{}) {
	err := validateHandler(handler)
	if err != nil {
		panic(err)
	}

	dispatcher.mx.Lock()
	for i := 0; i < len(dispatcher.subscriptions); i++ {
		group := &dispatcher.subscriptions[i]
		if group.name == groupName {
			dispatcher.mx.Unlock()
			group.mutex.Lock()
			group.handlers = append(group.handlers, handler)
			group.topHandler++
			group.mutex.Unlock()
			return
		}
	}

	dispatcher.subscriptions = append(dispatcher.subscriptions, subscriptionGroups{name: groupName, handlers: []interface{}{handler}})
	dispatcher.mx.Unlock()
}

// Subscribe an event eventName, you can subscribe to multiple event groups by separating eventName names with |
// example "group1|group2" will subscribe group1 and group2, take per example app.run and app.run.tls
//
//	subscribing to app.run event groups with be as simples as app.Subscribe("app.run|app.run.tls",func(e *event.Event,a *app.App){
//		println("App is starting a server ", e.EventName())
//	})
func (dispatcher *Dispatcher) Subscribe(events string, handler interface{}) *Dispatcher {
	dispatcher.assert()
	eventNames := strings.Split(events, "|")
	for i := 0; i < len(eventNames); i++ {
		dispatcher.subscribe(eventNames[i], handler)
	}
	return dispatcher
}

func (dispatcher *Dispatcher) group(groupName string) (group *subscriptionGroups, ok bool) {
	dispatcher.mx.RLock()
	numOfSubscriptions := len(dispatcher.subscriptions)
	for i := 0; i < numOfSubscriptions; i++ {
		group = &dispatcher.subscriptions[i]
		if ok = group.name == groupName; ok {
			break
		}
	}
	dispatcher.mx.RUnlock()
	return
}

// assert valid emitter
func (dispatcher *Dispatcher) assert() {
	if dispatcher.parent == nil && dispatcher != sub {
		panic(errors.New("All emitters are required to inherit from root,\ncheck if you'are not using a zero value or\nif you are not a using a Dispatcher struct instead of a pointer."))
	}
}

func (dispatcher *Dispatcher) emit(eventName string, event Payload) (canceled bool, err error) {
	dispatcher.assert()

	c := reflect.ValueOf(event)
	var _type = c.Type()
	var _arg = []reflect.Value{c}
	var hasUnsubscribes = false

	if group, ok := dispatcher.group(eventName); ok {
		group.mutex.RLock()
		defer group.mutex.RUnlock()

		for i := group.topHandler; i >= 0; i-- {
			if group.handlers[i] != nil {
				v := reflect.ValueOf(group.handlers[i])
				if _type.AssignableTo(v.Type().In(0)) {
					v.Call(_arg)
					if event.unsubscribed() {
						hasUnsubscribes = true
						group.handlers[i] = nil
					}
					canceled, err = event.WasCanceled(), event.error()
					if canceled || err != nil {
						return
					}
				}
			}
		}

		if hasUnsubscribes {
			group.mutex.RUnlock()
			group.clearNilHandlers()
			group.mutex.RLock()
		}
	}

	if dispatcher.parent != nil {
		canceled, err = dispatcher.parent.emit(eventName, event)
	}
	return
}

// Dispatch emits an event in the given eventName with the specified key,
// calling *Event.Cancel() will stop the event propagation, calling *Event.CancelWithError(err) will flag an error
// and cancel the event propagation
func (dispatcher *Dispatcher) Dispatch(registry *container.Registry, eventName string, event Payload) (bool, error) {
	event.init(registry, eventName)
	return dispatcher.emit(eventName, event)
}
