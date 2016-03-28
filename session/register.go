package session

import (
	"github.com/CloudyKit/framework/di"
	"encoding/gob"
	"reflect"
	"sync"
	"fmt"
)

var (
	sessionsTypes = map[reflect.Type]string{}
	rwMx = sync.Mutex{}
)

func persistPtr(typOf reflect.Type, diContext *di.Context, mapto string, i interface{}) {
	structTyp := typOf.Elem()
	if structTyp.Kind() != reflect.Struct {
		panic(fmt.Errorf("Type %q is not a pointer to struct", typOf))
	}

	sessionsTypes[typOf] = mapto
	diContext.Set(i, func(c *di.Context) (ret interface{}) {
		sess := c.Get((*Session)(nil)).(*Session)
		ret = sess.Get(mapto)
		if ret == nil {
			ret = reflect.New(structTyp).Interface()
			sess.Set(mapto, ret)
		}
		c.Map(ret)
		return
	})
}

func persistStruct(typOf reflect.Type, diContext *di.Context, mapto string, i interface{}) {
	diContext.Set(i, func(c *di.Context, t reflect.Value) {
		sess := c.Get((*Session)(nil)).(*Session)
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

func Persist(diContext *di.Context, i interface{}) error {
	return PersistKey(diContext, "", i)
}

func PersistKey(diContext *di.Context, key string, i interface{}) error {
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
			persistPtr(typOf, diContext, key, i)
		case reflect.Struct:
			persistStruct(typOf, diContext, key, i)
		default:
			panic(fmt.Errorf("Type %q is not a válid typ", typOf))
		}

		return nil
	}
	return fmt.Errorf("Type %q is already persistent", typOf)
}
