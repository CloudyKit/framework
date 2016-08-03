package scheme

import (
	"fmt"
	"github.com/CloudyKit/framework/validation"
	"reflect"
	"sync"
	"sync/atomic"
)

// RefKind represent the kind of referencing between entities and they fields
type RefKind int

const (
	RefRelatesTo   RefKind = 1 << iota // references
	RefRelatesBack                     // references a parent, source scheme has field point to target scheme
	RefChild                           // references a child, target scheme has a field point to source scheme
	RefChildren                        // references children, target scheme elements has a field point to source scheme
	RefParent                          // references a parent, source scheme has field point to target scheme
)

// Type any type implementing this interface is able to transform the types
// in the way that the fields is valid to insert into a database, ex: Int type will
// convert string type to int, before parsing to the database Driver
type Type interface {
	Value(v reflect.Value) (new reflect.Value, err error)
}

// Field holds metadata about the fields and the references
type Field struct {
	Name string
	Type Type

	Required bool

	RefKind   RefKind
	RefField  string
	RefScheme *Scheme
	RefTable  string
	Testers   []validation.Tester
}

// Index abstract index definition
type Index struct {
	Fields []string
	Type   string
}

// Scheme represents a database entity, fields and the references between different entities
type Scheme struct {
	primaryKey string
	name       string

	done uint32
	mx   sync.Mutex

	indexes []Index

	fields []*Field
}

func (scheme *Scheme) FieldByName(name string) (field *Field, found bool) {
	for _, f := range scheme.fields {
		if found = f.Name == name; found {
			field = f
			return
		}
	}
	return
}

func (scheme *Scheme) check() {

	for _, def := range scheme.fields {

		switch def.RefKind {
		case RefRelatesTo:
			if tdef, ok := def.RefScheme.FieldByName(def.RefField); ok {
				if tdef.Name != def.RefScheme.primaryKey {
					panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to a non primary key field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
				}
			} else {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefRelatesBack:
			if tdef, ok := def.RefScheme.FieldByName(def.RefField); ok {
				if tdef.Name != def.RefScheme.primaryKey {
					panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to a non primary key field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
				}
				def.RefTable = tdef.RefTable
			} else {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefChildren:
			// todo: fix this
			//if _, ok := def.RefScheme.FieldByName(def.RefField); !ok {
			//	panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			//}
		case RefChild:
			if _, ok := def.RefScheme.FieldByName(def.RefField); !ok {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefParent:
			if _, ok := scheme.FieldByName(def.RefField); ok {
				// todo: duplicate field
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, scheme.name))
			}
		case 0: // no ref
		default:
			panic(fmt.Errorf("unexpected ref kind %d", def.RefKind))
		}
	}
}

// KeyField returns the name of the key field
func (scheme *Scheme) KeyField() string {
	scheme.markUsed()
	return scheme.primaryKey
}

// Entity returns the name of the entity mapped in this scheme
func (scheme *Scheme) Entity() string {
	scheme.markUsed()
	return scheme.name
}

func (scheme *Scheme) markUsed() {
	if atomic.LoadUint32(&scheme.done) == 0 {
		scheme.mx.Lock()
		defer scheme.mx.Unlock()
		if scheme.done == 0 {
			scheme.check()
			atomic.StoreUint32(&scheme.done, 1)
		}
	}
}

// Indexes returns the list of indexes in the scheme
func (scheme *Scheme) Indexes() []Index {
	scheme.markUsed()
	return scheme.indexes
}

// Fields returns the list of fields in the scheme
func (scheme *Scheme) Fields() []*Field {
	scheme.markUsed()
	return scheme.fields
}

// Copy creates a copy of the scheme mapped to a new entity
func (scheme *Scheme) Copy(entityName string, def func(*Def)) *Scheme {
	return nil
}

// Extends creates a copy of the scheme with extra fields and references
func (scheme *Scheme) Extends(def func(*Def)) *Scheme {
	return nil
}
