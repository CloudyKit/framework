package scheme

import (
	"fmt"
	"github.com/CloudyKit/framework/validation"
	"regexp"
	"sync/atomic"
)

// Def is responsible the define the metadata for the scheme
type Def Scheme
type fieldDef struct {
	f *Field
}

func (f fieldDef) Required() fieldDef {
	f.f.Required = true
	return f
}

func (def *Def) assertField(fieldName string) {
	def.assertDone()
	for _, f := range def.fields {
		if f.Name == fieldName {
			panic(fmt.Errorf("Scheme.Field: field %s on scheme %s was already mapped.", fieldName, def.name))
		}
	}
}

// Refs creates reference between two schemes
func (def *Def) RefRelates(fieldName string, refScheme *Scheme, refTable, refField string, testers ...validation.Tester) fieldDef {
	def.assertField(fieldName)
	if refTable == "" {
		refTable = fmt.Sprintf("%s_to_%s", def.name, refScheme.name)
	}
	field := &Field{Name: fieldName, RefKind: RefRelatesTo, RefScheme: refScheme, RefField: refField, RefTable: refTable, Testers: testers}
	def.fields = append(def.fields, field)
	return fieldDef{field}
}

// Refs creates a back reference
func (def *Def) RefRelatesBack(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) fieldDef {
	def.assertField(fieldName)
	field := &Field{Name: fieldName, RefKind: RefRelatesBack, RefField: refField, RefScheme: refScheme, Testers: testers}
	def.fields = append(def.fields, field)
	return fieldDef{field}
}

// RefsChild creates a reference
func (def *Def) RefChild(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) fieldDef {
	def.assertField(fieldName)
	field := &Field{Name: fieldName, RefKind: RefChild, RefField: refField, RefScheme: refScheme, Testers: testers}
	def.fields = append(def.fields, field)
	return fieldDef{field}
}

func (def *Def) RefChildren(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) fieldDef {
	def.assertField(fieldName)
	field := &Field{Name: fieldName, RefKind: RefChildren, RefField: refField, RefScheme: refScheme, Testers: testers}
	def.fields = append(def.fields, field)
	return fieldDef{field}
}

func (def *Def) RefParent(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) fieldDef {
	def.assertField(fieldName)
	field := &Field{Name: fieldName, RefKind: RefParent, RefField: refField, RefScheme: refScheme, Testers: testers}
	def.fields = append(def.fields, field)
	return fieldDef{field}
}

//func (def *Def) RefsParents(fieldName string, refScheme *Scheme, refField string, testers ...validation.Tester) *Def {
//	def.assertField(fieldName)
//	field := &Field{Name: fieldName, RefKind: RefParents, RefField: refField, RefScheme: refScheme, testers: testers}
//def.fields[fieldName]=field
// return DefField{field}
//}

func (def *Def) Field(fieldName string, fieldType Type, t ...validation.Tester) fieldDef {
	def.assertField(fieldName)
	field := &Field{Name: fieldName, Type: fieldType, Testers: t}
	def.fields = append(def.fields, field)
	return fieldDef{field}
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
