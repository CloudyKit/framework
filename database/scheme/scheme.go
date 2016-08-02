package scheme

import (
	"fmt"
	"github.com/CloudyKit/framework/validation"
	"reflect"
	"regexp"
	"sync"
	"sync/atomic"
)

type RefKind int

const (
	RefTo       RefKind = 1 << iota // references
	RefFrom                         // references a parent, source scheme has field point to target scheme
	RefChild                        // references a child, target scheme has a field point to source scheme
	RefChildren                     // references children, target scheme elements has a field point to source scheme
	RefParent                       // references a parent, source scheme has field point to target scheme
	//RefParents                      // references multiple parents, source scheme has field point to target scheme
)

type Type interface {
	Value(v reflect.Value) (new reflect.Value, err error)
}

type Field struct {
	Name string
	Kind int
	Type Type

	RefKind   RefKind
	RefField  string
	RefScheme *Scheme
	RefTable  string
	testers   []validation.Tester
}

type Index struct {
	Fields []string
	Type   string
}

type Scheme struct {
	primaryKey string
	name       string

	done uint32
	mx   sync.Mutex

	indexes []Index

	fields map[string]*Field
}

func (scheme *Scheme) check() {
	var foundPrimary bool
	for name, def := range scheme.fields {

		if !foundPrimary && name == scheme.primaryKey {
			foundPrimary = true
		}

		switch def.RefKind {
		case RefTo:
			if tdef, ok := def.RefScheme.Fields()[def.RefField]; ok {
				if tdef.Name != def.RefScheme.primaryKey {
					panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to a non primary key field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
				}
			} else {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefFrom:
			if tdef, ok := def.RefScheme.Fields()[def.RefField]; ok {
				if tdef.Name != def.RefScheme.primaryKey {
					panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to a non primary key field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
				}
				def.RefTable = tdef.RefTable
			} else {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefChildren:
			if _, ok := def.RefScheme.Fields()[def.RefField]; !ok {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefChild:
			if _, ok := def.RefScheme.Fields()[def.RefField]; !ok {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, def.RefScheme.name))
			}
		case RefParent:
			if _, ok := scheme.fields[def.RefField]; !ok {
				panic(fmt.Errorf("Scheme Ref: field %s from scheme %s is referecing to an inexistent field %s in scheme %s", def.Name, scheme.name, def.RefField, scheme.name))
			}
		case 0: // no ref
		default:
			panic(fmt.Errorf("unexpected ref kind %d", def.RefKind))
		}
	}
}

func (scheme *Scheme) PrimaryKey() string {
	return scheme.primaryKey
}

func (scheme *Scheme) Entity() string {
	return scheme.name
}

func (scheme *Scheme) Fields() map[string]*Field {
	if atomic.LoadUint32(&scheme.done) == 1 {
		return scheme.fields
	}

	scheme.mx.Lock()
	defer scheme.mx.Unlock()
	if scheme.done == 0 {
		scheme.check()
		atomic.StoreUint32(&scheme.done, 1)
	}

	return scheme.fields
}

type Def Scheme

func (def *Def) assertField(fieldName string) {
	def.assertDone()
	//todo: add way to not disallow updates in schema after first run
	if _, found := def.fields[fieldName]; found {
		panic(fmt.Errorf("Scheme.Field: field %s on scheme %s was already mapped.", fieldName, def.name))
	}
}

func (def *Def) Refs(fieldName string, refScheme *Scheme, refTable, refField string, testers ...validation.Tester) *Def {
	def.assertField(fieldName)
	if refTable == "" {
		refTable = fmt.Sprintf("%s_to_%s", def.name, refScheme.name)
	}
	def.fields[fieldName] = &Field{Name: fieldName, RefKind: RefTo, RefScheme: refScheme, RefField: refField, RefTable: refTable, testers: testers}
	return def
}

func (def *Def) RefsFrom(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) *Def {
	def.assertField(fieldName)
	def.fields[fieldName] = &Field{Name: fieldName, RefKind: RefFrom, RefField: refField, RefScheme: refScheme, testers: testers}
	return def
}

func (def *Def) RefsChild(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) *Def {
	def.assertField(fieldName)
	def.fields[fieldName] = &Field{Name: fieldName, RefKind: RefChild, RefField: refField, RefScheme: refScheme, testers: testers}
	return def
}

func (def *Def) RefsChildren(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) *Def {
	def.assertField(fieldName)
	def.fields[fieldName] = &Field{Name: fieldName, RefKind: RefChildren, RefField: refField, RefScheme: refScheme, testers: testers}
	return def
}

func (def *Def) RefsParent(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) *Def {
	def.assertField(fieldName)
	def.fields[fieldName] = &Field{Name: fieldName, RefKind: RefParent, RefField: refField, RefScheme: refScheme, testers: testers}
	return def
}

//func (def *Def) RefsParents(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) *Def {
//	def.assertField(fieldName)
//	def.fields[fieldName] = &Field{Name: fieldName, RefKind: RefParents, RefField: refField, RefScheme: refScheme, testers: testers}
//	return def
//}

func (def *Def) Field(fieldName string, fieldType Type, t ...validation.Tester) *Def {
	def.assertField(fieldName)
	def.fields[fieldName] = &Field{Name: fieldName, Type: fieldType, testers: t}
	return def
}

var validIndex = regexp.MustCompile("^[+-]{0,1}[a-zA-Z_][a-zA-Z_0-9]*$")

func (def *Def) assertDone() {
	if atomic.LoadUint32(&def.done) == 1 {
		panic(fmt.Errorf("Scheme Def: scheme %s can't be modified after first use", def.name))
	}
}

func (def *Def) Index(typ string, fields ...string) *Def {
	def.assertDone()
	if len(fields) == 0 {
		panic(fmt.Errorf("Scheme Def: defining index type %s on scheme %s without fields", typ, def.name))
	}

	for _, field := range fields {
		if !validIndex.MatchString(field) {
			panic(fmt.Errorf("Scheme Def: defining index type %s on scheme %s with an invalid name %s", typ, def.name, field))
		}
	}

	def.indexes = append(def.indexes, Index{Type: typ, Fields: fields})
	return def
}
