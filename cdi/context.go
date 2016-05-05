package cdi

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	DI struct {
		parent     *DI
		references int64
		values     map[reflect.Type]interface{}
	}

	finalizer interface {
		Finalize()
	}

	provider interface {
		Provide(c *DI) interface{}
	}

	providerSetter interface {
		Provide(c *DI, field reflect.Value)
	}
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			cc := new(DI)
			cc.values = make(map[reflect.Type]interface{})
			return cc
		},
	}
	walkableFields = map[reflect.Type]struct{}{}
)

// Walkable marks the type as field to be walked while searching for fields to inject
func Walkable(v ...interface{}) uint {
	for i := 0; i < len(v); i++ {
		typ := reflect.TypeOf(v[i])
	TRY:
		switch typ.Kind() {
		case reflect.Struct:
		case reflect.Ptr:
			typ = typ.Elem()
			goto TRY
		default:
			continue
		}

		walkableFields[typ] = struct{}{}
	}
	return 0
}

// New creates a new instance of context object
func New() (cc *DI) {
	cc = pool.Get().(*DI)
	return
}

// Context holds the dependency injection data

// Inject walks the target looking the for exported fields that types match injectable types in Global
func (c *DI) Inject(target interface{}) {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	c.InjectInStructValue(value)
}

// InjectInStructValue walks the struct value looking to injectable fields
func (c *DI) InjectInStructValue(value reflect.Value) {
	if value.Kind() != reflect.Struct {
		panic("Invalid value passed to inject, required kind is struct get " + value.Kind().String())
	}
	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldTyp := field.Type()

		if provided_value, wasSetted := c.val4TypeField(fieldTyp, field); provided_value != nil || wasSetted {
			if !wasSetted {
				field.Set(reflect.ValueOf(provided_value))
			}
			continue
		}

		if _, ok := walkableFields[fieldTyp]; ok {
			c.InjectInStructValue(field)
		}
	}
}

// Set sets a provider for the type of typ with value of val
func (c *DI) MapType(typ, val interface{}) {
	typOf, ok := typ.(reflect.Type)

	if !ok {
		typOf = reflect.TypeOf(typ)
	}

	if typOf.Kind() == reflect.Ptr && typOf.Elem().Kind() == reflect.Interface {
		typOf = typOf.Elem()
	}
	c.values[typOf] = val
}

// Put puts the list of values into the current context
func (c *DI) Map(value ...interface{}) {
	for i := 0; i < len(value); i++ {
		vof := value[i]
		v := reflect.ValueOf(vof)
		c.values[v.Type()] = vof
	}
}

// From put in to the context all child values from the provided value
// example all exported fields from a Struct or itens from a Slice
func (c *DI) From(st interface{}) {
	v := reflect.ValueOf(st)
RESTART:
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		goto RESTART
	case reflect.Slice:
		length := v.Len()
		for i := 0; i < length; i++ {
			field := v.Index(i)
			if field.CanInterface() {
				c.values[field.Type()] = field.Interface()
			}
		}
	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			field := v.Field(i)
			if field.CanInterface() {
				c.values[field.Type()] = field.Interface()
			}
		}
	default:
		panic("Invalid kind")
	}
}

// Child creates a new context using current values repository as provider for the new context
func (c *DI) Child() (child *DI) {
	c.references = atomic.AddInt64(&c.references, 1)
	child = New()
	child.parent = c
	return
}

// val4type walkings the context from the current to the top parent looking for the value with type typ
func (_context *DI) val4type(typ reflect.Type) (val interface{}) {
	for {
		val = _context.values[typ]
		if val != nil || _context.parent == nil {
			return
		}
		_context = _context.parent
	}
	return
}

// val4TypeField returns a value for the specified type typ
func (c *DI) val4TypeField(typ reflect.Type, valOf reflect.Value) (val interface{}, ok bool) {
	val = c.val4type(typ)
	switch provider := val.(type) {
	case func(*DI) interface{}:
		val = provider(c)
	case func(*DI, reflect.Value):
		provider(c, valOf)
		ok = true
	case provider:
		val = provider.Provide(c)
	case providerSetter:
		provider.Provide(c, valOf)
		ok = true
	}
	return
}

// Val4Type returns a value for the specified type typ
func (c *DI) Val4Type(typ reflect.Type) (val interface{}) {
	val = c.val4type(typ)
	if valFn, ok := val.(func(*DI) interface{}); ok {
		val = valFn(c)
	} else if _provider, isProvider := val.(provider); isProvider {
		val = _provider.Provide(c)
	}
	return
}

// Get returns a value for the type of typ
func (c *DI) Get(typ interface{}) interface{} {
	return c.Val4Type(reflect.TypeOf(typ))
}

// Done should be called when the context is not being used anymore
func (c *DI) Done() {
	// check if this is the last active reference
	c.references = atomic.AddInt64(&c.references, -1)

	if c.references == -1 {
		c.finalize()
	} else if c.references < -1 {
		panic(fmt.Errorf("Inválid reference counting expected value is -1 got %v", c.references))
	}
}

// finalize walks all values in the current context and invokes finalizers
// decrease reference counter into the parent
// and recycle the private data
func (c *DI) finalize() {
	// invokes parent Done method
	if c.parent != nil {
		defer c.parent.Done()
	}

	//runs recycle here
	for _typ, _val := range c.values {
		// not delete the keys
		delete(c.values, _typ)
		if _finalizer, isFinalizer := _val.(finalizer); isFinalizer {
			_finalizer.Finalize()
		}
	}
	c.references = 0
	c.parent = nil
	pool.Put(c)
}
