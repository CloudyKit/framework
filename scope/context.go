package scope

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	Variables struct {
		parent     *Variables
		references int64
		values     map[reflect.Type]interface{}
	}

	finalizer interface {
		Finalize()
	}

	provider interface {
		Provide(c *Variables) interface{}
	}

	setter interface {
		Provide(c *Variables, field reflect.Value)
	}
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			cc := new(Variables)
			cc.values = make(map[reflect.Type]interface{})
			return cc
		},
	}
	walkableFields = map[reflect.Type]struct{}{}
)

func TypeOfElem(i interface{}) reflect.Type {
	return reflect.TypeOf(i).Elem()
}

func TypeOf(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}

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
func New() (cc *Variables) {
	cc = pool.Get().(*Variables)
	cc.references = 0
	return
}

//Inject walks the target looking the for exported fields that types match injectable types in Global
func (c *Variables) Inject(target interface{}) {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	c.InjectInStructValue(value)
}

var cdiType = reflect.TypeOf((*Variables)(nil))

//InjectInStructValue walks the struct value looking to injectable fields
func (c *Variables) InjectInStructValue(value reflect.Value) {
	if value.Kind() != reflect.Struct {
		panic("Invalid value passed to inject, required kind is struct get " + value.Kind().String())
	}
	numFields := value.NumField()
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldTyp := field.Type()

		if provided_value, wasSetted := c.getByTypeField(fieldTyp, field); provided_value != nil || wasSetted {
			if !wasSetted {
				field.Set(reflect.ValueOf(provided_value))
			}
		} else if cdiType == fieldTyp {
			field.Set(reflect.ValueOf(c))
		} else if _, ok := walkableFields[fieldTyp]; ok {
			c.InjectInStructValue(field)
		}
	}
}

//MapType sets a provider for the type of typ with value of val
func (c *Variables) MapType(typOf reflect.Type, val interface{}) {
	c.values[typOf] = val
}

//Map puts the list of values into the current context
func (c *Variables) Map(value ...interface{}) {
	for i := 0; i < len(value); i++ {
		vof := value[i]
		v := reflect.ValueOf(vof)
		c.values[v.Type()] = vof
	}
}

//Inherit creates a new context using current values repository as provider for the new context
func (c *Variables) Inherit() (child *Variables) {
	if atomic.LoadInt64(&c.references) > -1 {
		c.references = atomic.AddInt64(&c.references, 1)
		child = New()
		child.parent = c
	} else {
		panic(errors.New("Invoking Child in a context already recycled"))
	}
	return
}

//getByType searchs for value of type typ, walking the context tree from the current to the top parent looking for the value with type typ
func (_context *Variables) getByType(typ reflect.Type) (val interface{}) {
	for {
		val = _context.values[typ]
		if val != nil || _context.parent == nil {
			return
		}
		_context = _context.parent
	}
	return
}

//getByTypeField returns a value for the specified type typ
func (c *Variables) getByTypeField(typ reflect.Type, valOf reflect.Value) (val interface{}, ok bool) {
	val = c.getByType(typ)
	switch provider := val.(type) {
	case func(*Variables) interface{}:
		val = provider(c)
	case func(*Variables, reflect.Value):
		provider(c, valOf)
		ok = true
	case provider:
		val = provider.Provide(c)
	case setter:
		provider.Provide(c, valOf)
		ok = true
	}
	return
}

//GetByType returns a value for the specified type typ
func (c *Variables) GetByType(typ reflect.Type) (val interface{}) {
	val = c.getByType(typ)
	switch provider := val.(type) {
	case func(*Variables) interface{}:
		val = provider(c)
	case func(*Variables, reflect.Value):
		valOf := reflect.New(typ).Elem()
		provider(c, valOf)
		val = valOf.Interface()
	case provider:
		val = provider.Provide(c)
	case setter:
		valOf := reflect.New(typ).Elem()
		provider.Provide(c, valOf)
		val = valOf.Interface()
	}
	return
}

//GetByPtr searches a value which has same typ that typ points to and returns the value from the context tree
func (c *Variables) GetByPtr(typ interface{}) interface{} {
	return c.GetByType(reflect.TypeOf(typ))
}

// End call end when the request is not need any more, this will cause all finalizers to run,
// this will my cause parent scopes to end also, this will happen case the parent scopes already ended but
// is waiting all children to end
func (c *Variables) End() int64 {
	// check if this is the last active reference
	c.references = atomic.AddInt64(&c.references, -1)

	if c.references == -1 {
		c.finalize()
	} else if c.references < -1 {
		panic(fmt.Errorf("Inválid reference counting expected value is -1 got %v", c.references))
	}
	return c.references
}

var err = errors.New("scope.Variables.EndForce: requested that at this point all references to this context are previous cleared")

// EndForce works same as End, but require that all children to be terminated when called
func (c *Variables) EndForce() {
	if c.End() > -1 {
		panic(err)
	}
}

func (c *Variables) recycle() {
	if c.parent != nil {
		c.parent.End()
		c.parent = nil
	}
	pool.Put(c)
}

// finalize walks all values in the current context and invokes finalizers
// decrease reference counter into the parent
// and recycle the private data
func (c *Variables) finalize() {
	// invokes parent Done method
	defer c.recycle()
	//runs recycle here
	for _typ, _val := range c.values {
		// not delete the keys
		delete(c.values, _typ)
		if _finalizer, isFinalizer := _val.(finalizer); isFinalizer {
			_finalizer.Finalize()
		}
	}
}

func Checkpoint(c **Variables) func() {
	*c = (*c).Inherit()
	return func() {
		cc := (*c)
		*c = cc.parent
		cc.End()
	}
}
