package Di

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	private_Context struct {
		parent     *private_Context
		references int64
		values     map[reflect.Type]interface{}
	}

	closer interface {
		Close()
	}

	done interface {
		Done()
	}
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			cc := new(private_Context)
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
func New() (cc Context) {
	cc.private_Context = pool.Get().(*private_Context)
	return
}

// Context holds the dependency injection data
type Context struct {
	*private_Context
}

// Inject walks the target looking the exported fields for injectable values
func (c Context) Inject(target interface{}) {
	value := reflect.ValueOf(target)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		panic("Invalid value passed to inject, required kind is struct get " + value.Kind().String())
	}

	c.injectfields(value)
}

// injectfields walks the struct value looking to injectable fields
func (c Context) injectfields(value reflect.Value) {
	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldTyp := field.Type()
		if provided_value := c.Val4Type(fieldTyp); provided_value != nil {
			field.Set(reflect.ValueOf(provided_value))
		} else {
			if _, ok := walkableFields[fieldTyp]; ok {
				c.injectfields(field)
			}
		}
	}
}

// Set sets a provider for the type of typ with value of val
func (c Context) Set(typ, val interface{}) {
	c.values[reflect.TypeOf(typ)] = val
}

// Put puts the list of values into the current context
func (c Context) Put(value ...interface{}) {
	for i := 0; i < len(value); i++ {
		vof := value[i]
		v := reflect.ValueOf(vof)
		c.private_Context.values[v.Type()] = vof
	}
}

// From put in to the context all child values from the provided value
// example all exported fields from a Struct or itens from a Slice
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
				c.private_Context.values[field.Type()] = field.Interface()
			}
		}
	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			field := v.Field(i)
			if field.CanInterface() {
				c.private_Context.values[field.Type()] = field.Interface()
			}
		}
	default:
		panic("Invalid kind")
	}
}

// Child creates a new context using current values repository as provider for the new context
func (c Context) Child() (child Context) {
	c.references = atomic.AddInt64(&c.references, 1)
	child.private_Context = pool.Get().(*private_Context)
	child.parent = c.private_Context
	return
}

// val4type walkings the context from the current to the top parent looking for the value with type typ
func (_context *private_Context) val4type(typ reflect.Type) (val interface{}) {
	for {
		val = _context.values[typ]
		if val != nil || _context.parent == nil {
			return
		}
		_context = _context.parent
	}
	return
}

// Val4Type returns a value for the specified type typ
func (c Context) Val4Type(typ reflect.Type) (val interface{}) {
	val = c.val4type(typ)
	if valFn, ok := val.(func(Context) interface{}); ok {
		val = valFn(c)
	}
	return
}

// Get returns a value for the type of typ
func (c Context) Get(typ interface{}) interface{} {
	return c.Val4Type(reflect.TypeOf(typ))
}

// Done should be called when the context is not being used anymore
func (c *private_Context) Done() {
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
func (c *private_Context) finalize() {
	// invokes parent Done method
	if c.parent != nil {
		defer c.parent.Done()
	}
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
	c.references = 0
	pool.Put(c)
	c.parent = nil
}
