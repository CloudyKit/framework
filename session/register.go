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

package session

import (
	"encoding/gob"
	"fmt"
	"github.com/CloudyKit/framework/container"
	"reflect"
	"sync"
)

var (
	sessionsTypes = map[reflect.Type]string{}
	rwMx          = sync.Mutex{}
)

func persistPtr(typOf reflect.Type, c *container.IoC, mapto string) {
	structTyp := typOf.Elem()

	if structTyp.Kind() != reflect.Struct {
		panic(fmt.Errorf("Type %q is not a pointer to struct", typOf))
	}

	sessionsTypes[typOf] = mapto
	c.MapValue(typOf, func(c *container.IoC) (ret interface{}) {
		sess := LoadSession(c)
		ret = sess.Get(mapto)
		if ret == nil {
			ret = reflect.New(structTyp).Interface()
			sess.Set(mapto, ret)
		}
		c.MapValue(typOf, ret)
		return
	})
}

func persistStruct(typOf reflect.Type, c *container.IoC, mapto string) {
	c.MapValue(typOf, func(c *container.IoC, t reflect.Value) {
		sess := LoadSession(c)
		val := sess.Get(mapto)
		if val != nil {
			valueOf := reflect.ValueOf(val)
			if valueOf.Kind() == reflect.Ptr {
				valueOf = valueOf.Elem()
			}
			t.Set(valueOf)
		}
		sess.Set(mapto, t.Addr().Interface())
	})
}

func Persist(c *container.IoC, i interface{}) error {
	return PersistKey(c, "", i)
}

func PersistKey(c *container.IoC, key string, i interface{}) error {
	rwMx.Lock()
	defer rwMx.Unlock()
	typOf := reflect.TypeOf(i)

	if key == "" {
		key = typOf.String()
	}
	if _, exists := sessionsTypes[typOf]; !exists {

		//maps type to key
		sessionsTypes[typOf] = key

		//register gob type
		gob.Register(i)

		switch typOf.Kind() {
		case reflect.Ptr:
			persistPtr(typOf, c, key)
		case reflect.Struct:
			persistStruct(typOf, c, key)
		default:
			panic(fmt.Errorf("Type %q is not a válid typ", typOf))
		}

		return nil
	}
	return fmt.Errorf("Type %q is already persistent", typOf)
}
