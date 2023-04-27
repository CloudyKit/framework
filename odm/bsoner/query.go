package bsoner

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter struct {
	Doc primitive.D
}

func (filter *Filter) MarshalBSON() ([]byte, error) {
	return bson.Marshal(filter.Doc)
}

func NewFilterBuilder() *Filter {
	return &Filter{}
}

func FilterBy(k string, v interface{}) *Filter {
	return NewFilterBuilder().Set(k, v)
}

// Set sets bson field k with the value v row[k]=v
func (filter *Filter) Set(k string, v interface{}) *Filter {
	numDocs := len(filter.Doc)
	for i := 0; i < numDocs; i++ {
		d := &filter.Doc[i]
		if d.Key == k {
			d.Value = v
			return filter
		}
	}
	filter.Doc = append(filter.Doc, primitive.E{Key: k, Value: v})
	return filter
}

// ByID query by object _id
func ByID(id interface{}) bson.D {
	return primitive.D{{"_id", id}}
}

// lt less
type lt struct {
	Lt interface{} `bson:"$lt"`
}

// Lt sets bson field k with row[k][$lt]=v which can be express like row[k] == v
func (filter *Filter) Lt(k string, v interface{}) *Filter {
	filter.Set(k, lt{v})
	return filter
}

// Lt creates a new query object and sets bson field k with row[k][$eq]=v which can be express like row[k] == v
func Lt(k string, v interface{}) *Filter {
	return NewFilterBuilder().Lt(k, v)
}

// -- lte --
type lte struct {
	Lte interface{} `bson:"$lte"`
}

// Lte sets bson field k with row[k][$lte]=v which can be express like row[k] == v
func (filter *Filter) Lte(k string, v interface{}) *Filter {
	filter.Set(k, lte{v})
	return filter
}

// Lte creates a new query object and sets bson field k with row[k][$eq]=v which can be express like row[k] == v
func Lte(k string, v interface{}) *Filter {
	return NewFilterBuilder().Gt(k, v)
}

// -- gt --
type gt struct {
	Gt interface{} `bson:"$gt"`
}

// Gt sets bson field k with row[k][$gt]=v which can be express like row[k] == v
func (filter *Filter) Gt(k string, v interface{}) *Filter {
	filter.Set(k, gt{v})
	return filter
}

// Gt creates a new query object and sets bson field k with row[k][$eq]=v which can be express like row[k] == v
func Gt(k string, v interface{}) *Filter {
	return NewFilterBuilder().Gt(k, v)
}

// -- gte --
type gte struct {
	Gte interface{} `bson:"$gte"`
}

// Gte sets bson field k with row[k][$gte]=v which can be express like row[k] == v
func (filter *Filter) Gte(k string, v interface{}) *Filter {
	filter.Set(k, gte{v})
	return filter
}

// Gte creates a new query object and sets bson field k with row[k][$eq]=v which can be express like row[k] == v
func Gte(k string, v interface{}) *Filter {
	return NewFilterBuilder().Gt(k, v)
}

// -- eq --
type eq struct {
	Eq interface{} `bson:"$eq"`
}

// Eq sets bson field k with row[k][$eq]=v which can be express like row[k] == v
func (filter *Filter) Eq(k string, v interface{}) *Filter {
	filter.Set(k, eq{Eq: v})
	return filter
}

// Eq creates a new query object and sets bson field k with row[k][$eq]=v which can be express like row[k] == v
func Eq(k string, v interface{}) *Filter {
	return NewFilterBuilder().Eq(k, v)
}

// -- neq --
type neq struct {
	Neq interface{} `bson:"$ne"`
}

// Neq sets bson field k with row[k][$ne]=v which can be express like row[k] != v
func (filter *Filter) Neq(k string, v interface{}) *Filter {
	filter.Set(k, neq{v})
	return filter
}

// Neq creates a new query object and sets bson field k with row[k][$ne]=v which can be express like row[k] != v
func Neq(k string, v interface{}) *Filter {
	return NewFilterBuilder().Neq(k, v)
}

// -- in --
type in struct {
	In interface{} `bson:"$in"`
}

// In sets bson field k with row[k][$ne]=v which can be express like row[k] != v
func (filter *Filter) In(k string, v ...interface{}) *Filter {
	filter.Set(k, in{v})
	return filter
}

// In creates a new query object and sets bson field k with row[k][$ne]=v which can be express like row[k] != v
func In(k string, v ...interface{}) *Filter {
	return NewFilterBuilder().In(k, v)
}

// -- nin --
type nin struct {
	Nin interface{} `bson:"$nin"`
}

// Nin sets bson field k with row[k][$ne]=v which can be express like row[k] != v
func (filter *Filter) Nin(k string, v ...interface{}) *Filter {
	filter.Set(k, nin{v})
	return filter
}

// Nin creates a new query object and sets bson field k with row[k][$ne]=v which can be express like row[k] != v
func Nin(k string, v ...interface{}) *Filter {
	return NewFilterBuilder().Nin(k, v)
}

// RegEx sets bson field k with a regex pattern and options
func (filter *Filter) RegEx(k, pattern, options string) *Filter {
	filter.Set(k, primitive.Regex{Pattern: pattern, Options: options})
	return filter
}

// RegEx creates a new query object and sets bson field k with regex pattern and options
func RegEx(k, pattern, options string) *Filter {
	return NewFilterBuilder().RegEx(k, pattern, options)
}

// And is equivalent to expr0 && expr1 && expr...
func (filter *Filter) And(docs ...interface{}) *Filter {
	numDocs := len(filter.Doc)
	for i := 0; i < numDocs; i++ {
		d := &filter.Doc[i]
		if d.Key == "$and" {
			d.Value = append(d.Value.([]interface{}), docs...)
			return filter
		}
	}
	filter.Doc = append(filter.Doc, primitive.E{Key: "$and", Value: docs})
	return filter
}

// And is equivalent to expr0 && expr1 && expr...
func And(docs ...interface{}) *Filter {
	return NewFilterBuilder().And(docs...)
}

// Or is equivalent to expr0 || expr1 || expr...
func (filter *Filter) Or(docs ...interface{}) *Filter {
	numDocs := len(filter.Doc)
	for i := 0; i < numDocs; i++ {
		d := &filter.Doc[i]
		if d.Key == "$or" {
			d.Value = append(d.Value.([]interface{}), docs...)
			return filter
		}
	}
	filter.Doc = append(filter.Doc, primitive.E{Key: "$or", Value: docs})
	return filter
}

// Or is equivalent to expr0 || expr1 || expr...
func Or(docs ...interface{}) *Filter {
	return NewFilterBuilder().Or(docs...)
}

// NotOr is equivalent to !expr0 || !expr1 || !expr...
func (filter *Filter) NotOr(docs ...interface{}) *Filter {
	numDocs := len(filter.Doc)
	for i := 0; i < numDocs; i++ {
		d := &filter.Doc[i]
		if d.Key == "$nor" {
			d.Value = append(d.Value.([]interface{}), docs...)
			return filter
		}
	}
	filter.Doc = append(filter.Doc, primitive.E{Key: "$nor", Value: docs})
	return filter
}

// NotOr is equivalent to !expr0 || !expr1 || !expr...
func NotOr(docs ...interface{}) *Filter {
	return NewFilterBuilder().NotOr(docs...)
}

// -- exists --
type exists struct {
	Exists interface{} `bson:"$exists"`
}

// Exists mongo $exists query
func (filter *Filter) Exists(k string, v interface{}) *Filter {
	filter.Set(k, exists{v})
	return filter
}

// Exists is equivalent to !expr0 || !expr1 || !expr...
func Exists(k string, v interface{}) *Filter {
	return NewFilterBuilder().Exists(k, v)
}

func NewList(list ...interface{}) []interface{} {
	return list
}
