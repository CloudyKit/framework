package session

import (
	"encoding/gob"
	"fmt"
	"github.com/CloudyKit/framework/scope"
	"reflect"
	"sync"
)

var (
	sessionsTypes = map[reflect.Type]string{}
	rwMx          = sync.Mutex{}
)

func persistPtr(typOf reflect.Type, c *scope.Variables, mapto string) {
	structTyp := typOf.Elem()

	if structTyp.Kind() != reflect.Struct {
		panic(fmt.Errorf("Type %q is not a pointer to struct", typOf))
	}

	sessionsTypes[typOf] = mapto
	c.MapType(typOf, func(c *scope.Variables) (ret interface{}) {
		sess := GetSession(c)
		ret = sess.Get(mapto)
		if ret == nil {
			ret = reflect.New(structTyp).Interface()
			sess.Set(mapto, ret)
		}
		c.MapType(typOf, ret)
		return
	})
}

func persistStruct(typOf reflect.Type, c *scope.Variables, mapto string) {
	c.MapType(typOf, func(c *scope.Variables, t reflect.Value) {
		sess := GetSession(c)
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

func Persist(c *scope.Variables, i interface{}) error {
	return PersistKey(c, "", i)
}

func PersistKey(c *scope.Variables, key string, i interface{}) error {
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
