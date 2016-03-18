package Di

import (
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

type context struct {
	parent     *context
	references int64
	values     map[reflect.Type]interface{}
}

var __context_pool = sync.Pool{
	New: func() interface{} {
		cc := new(context)
		cc.values = make(map[reflect.Type]interface{})
		return cc
	},
}

var __walkable = map[reflect.Type]struct{}{}

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

		__walkable[typ] = struct{}{}
	}
	return 0
}

// New creates a new instance of context object
func New() (cc Context) {
	cc.context = __context_pool.Get().(*context)
	return
}

type Context struct {
	*context
}

func (c Context) Inject(target interface{}) {
	value := reflect.ValueOf(target)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		panic("Invalid value passed to inject, required kind is struct get " + value.Kind().String())
	}

	c.injectFields(value)
}

func (c Context) injectFields(value reflect.Value) {
	numFields := value.NumField()

	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldTyp := field.Type()
		if provided_value := c.Val4Typ(fieldTyp); provided_value != nil {
			field.Set(reflect.ValueOf(provided_value))
		} else if _, ok := __walkable[fieldTyp]; ok {
			c.injectFields(field)
		}
	}
}

func (c Context) Set(typ, val interface{}) {
	c.values[reflect.TypeOf(typ)] = val
}

func (c Context) Put(value ...interface{}) {
	for i := 0; i < len(value); i++ {
		vof := value[i]
		v := reflect.ValueOf(vof)
		c.context.values[v.Type()] = vof
	}
}

//From put in to the context all child values from the provided value
//example all exported fields from a Struct or itens from a Slice
func (c Context) From(st interface{}) {
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
				c.context.values[field.Type()] = field.Interface()
			}
		}
	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			field := v.Field(i)
			if field.CanInterface() {
				c.context.values[field.Type()] = field.Interface()
			}
		}
	default:
		panic("Invalid kind")
	}
}

//Child creates a new context using current values repository as provider for the new context
func (c Context) Child() (child Context) {
	atomic.AddInt64(&c.references, 1)
	child.context = __context_pool.Get().(*context)
	child.parent = c.context
	return
}

func (_context *context) val4typ(typ reflect.Type) (val interface{}) {
	for {
		val = _context.values[typ]
		if val != nil || _context.parent == nil {
			return
		}
		_context = _context.parent
	}
	return
}

//Val4Typ returns an value for the specified type
func (c Context) Val4Typ(typ reflect.Type) (val interface{}) {
	val = c.val4typ(typ)
	if valFn, ok := val.(func(Context) interface{}); ok {
		val = valFn(c)
	}
	return
}

func (c Context) Get(typ interface{}) interface{} {
	return c.Val4Typ(reflect.TypeOf(typ))
}

//Done should be called when the context is not being used anymore
func (c *context) Done() {
	// check if this is the last active reference
	referenceCounting := atomic.AddInt64(&c.references, -1)
	if referenceCounting == -1 {
		c.done()
	}
}

func (c *context) done() {
	type (
		closer interface {
			Close()
		}

		done interface {
			Done()
		}
	)

	//runs recycle here
	for _typ, _val := range c.values {
		switch _val := _val.(type) {
		case io.Closer:
			_val.Close()
		case closer:
			_val.Close()
		case done:
			_val.Done()
		}
		// not delete keys for caching
		c.values[_typ] = nil
	}
	c.parent = nil
	c.references = 0
	__context_pool.Put(c)

	// invokes parent Done method
	if c.parent != nil {
		c.parent.Done()
	}
}
