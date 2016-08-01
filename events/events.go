package events

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var eventTYPE = reflect.TypeOf((*Event)(nil))

type Event struct {
	i        int64
	group    string
	key      string
	err      error
	canceled bool
	*Emitter
}

func (e *Event) Group() string {
	return e.group
}

func (e *Event) Key() string {
	return e.key
}

func (e *Event) UnSubscribe() {
	if sub, found := e.Emitter.group(e.group); found {
		go func(i int64) {
			sub.rwmutex.Lock()
			sub.handlers = append(sub.handlers[0:i], sub.handlers[i+1:]...)
			sub.rwmutex.Unlock()
		}(e.i)
	}
}

func (e *Event) Cancel() {
	e.canceled = true
}

func (e *Event) Error(err error) {
	e.err = err
}

func (e *Event) Errorf(format string, v ...interface{}) {
	e.err = fmt.Errorf(format, v...)
}

type subscriptionGroups struct {
	rwmutex      sync.RWMutex
	name         string
	handlers     []interface{}
	topHandler   int64
	runningEmits int64
}

type Emitter struct {
	parent        *Emitter
	mx            sync.RWMutex
	subscriptions []subscriptionGroups
}

func (manager *Emitter) Inherit() *Emitter {
	return &Emitter{parent: manager}
}

var err = errors.New("unexpected handler signature: func(*Event,ContextType) is expected")

func validateHandler(h interface{}) error {
	t := reflect.TypeOf(h)

	if t.Kind() != reflect.Func || t.NumOut() != 0 {
		return err
	}

	numIn := t.NumIn()
	if numIn > 2 || eventTYPE.AssignableTo(t.In(0)) == false || t.IsVariadic() {
		return err
	}

	return nil
}

func (manager *Emitter) Reset(groupName string) bool {
	manager.assert()

	subsgroup, ok := manager.group(groupName)
	if ok {
		subsgroup.rwmutex.Lock()
		subsgroup.topHandler = -1
		subsgroup.handlers = nil
		subsgroup.rwmutex.Unlock()
	}
	return ok
}

func (manager *Emitter) subscribe(groupName string, handler interface{}) {
	err := validateHandler(handler)
	if err != nil {
		panic(err)
	}

	manager.mx.Lock()
	for i := 0; i < len(manager.subscriptions); i++ {
		group := &manager.subscriptions[i]
		if group.name == groupName {
			manager.mx.Unlock()
			group.rwmutex.Lock()
			group.handlers = append(group.handlers, handler)
			group.topHandler++
			group.rwmutex.Unlock()
			return
		}
	}

	manager.subscriptions = append(manager.subscriptions, subscriptionGroups{name: groupName, handlers: []interface{}{handler}})
	manager.mx.Unlock()
}

// Subscribe an event group, you can subscribe to multiple events groups by separating group names with |
// example "group1|group2" will subscribe group1 and group2, take per example app.run and app.run.tls
// subscribing to app.run events groups with be as simples as app.Subscribe("app.run|app.run.tls",func(e *events.Event,a *app.App){
// 	println("App is starting a server ", e.Group())
// })
func (manager *Emitter) Subscribe(groups string, handler interface{}) *Emitter {
	manager.assert()
	groupnames := strings.Split(groups, "|")
	for i := 0; i < len(groupnames); i++ {
		manager.subscribe(groupnames[i], handler)
	}
	return manager
}

func (manager *Emitter) group(groupName string) (group *subscriptionGroups, ok bool) {
	manager.mx.RLock()
	numOfSubscriptions := len(manager.subscriptions)
	for i := 0; i < numOfSubscriptions; i++ {
		group = &manager.subscriptions[i]
		if ok = group.name == groupName; ok {
			break
		}
	}
	manager.mx.RUnlock()
	return
}

// assert valid emitter
func (manager *Emitter) assert() {
	if manager.parent == nil && manager != sub {
		panic(errors.New("All emitters are required to inherit from root,\ncheck if you'are not using a zero value or\nif you are not a using a Emitter struct instead of a pointer."))
	}
}

func (manager *Emitter) emit(event *Event, c, e reflect.Value) (canceled bool, err error) {
	manager.assert()

	var _type = c.Type()
	var _arg [2]reflect.Value
	var _argslice = _arg[:]
	_argslice[0] = e
	_argslice[1] = c

	if group, ok := manager.group(event.group); ok {

		group.rwmutex.RLock()
		defer group.rwmutex.RUnlock()

		for event.i = group.topHandler; event.i >= 0; event.i-- {
			v := reflect.ValueOf(group.handlers[event.i])

			if _type.AssignableTo(v.Type().In(1)) {
				v.Call(_argslice)
				canceled, err = event.canceled, event.err
				if canceled || err != nil {
					return
				}
			}
		}
	}

	if manager.parent != nil {
		canceled, err = manager.parent.emit(event, c, e)
	}
	return
}

// Emit emits an event in the given group with the specified key,
// calling *Event.Cancel() will stop the event propagation, calling *Event.Error(err) will flag an error
// and cancel the event propagation
func (manager *Emitter) Emit(groupName, key string, context interface{}) (bool, error) {
	var event = &Event{group: groupName, key: key}
	return manager.emit(event, reflect.ValueOf(context), reflect.ValueOf(event))
}
