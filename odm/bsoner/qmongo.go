package bsoner

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Doc primitive.D
}

func (doc *Document) MarshalBSON() ([]byte, error) {
	if doc == nil {
		return bson.Marshal(nil)
	}
	return bson.Marshal(doc.Doc)
}

func NewDocumentBuilder() *Document {
	return &Document{}
}

func NewDocumentSet(k string, v interface{}) *Document {
	return NewDocumentBuilder().Set(k, v)
}

// Set sets bson field k with the value v row[k]=v
func (doc *Document) Set(k string, v interface{}) *Document {
	numDocs := len(doc.Doc)
	for i := 0; i < numDocs; i++ {
		d := &doc.Doc[i]
		if d.Key == k {
			d.Value = v
			return doc
		}
	}
	doc.Doc = append(doc.Doc, primitive.E{Key: k, Value: v})
	return doc
}

func (doc *Document) DocSet(k string, v interface{}) *Document {
	numDocs := len(doc.Doc)
	for i := 0; i < numDocs; i++ {
		d := &doc.Doc[i]
		if d.Key == "$set" {
			d.Value.(*Document).Set(k, v)
			return doc
		}
	}

	doc.Doc = append(doc.Doc, primitive.E{Key: "$set", Value: NewDocumentBuilder().Set(k, v)})
	return doc
}

func DocSet(k string, v interface{}) *Document {
	return NewDocumentBuilder().DocSet(k, v)
}

func (doc *Document) SetOnInsert(k string, v interface{}) *Document {
	numDocs := len(doc.Doc)
	for i := 0; i < numDocs; i++ {
		d := &doc.Doc[i]
		if d.Key == "$setOnInsert" {
			d.Value.(*Document).Set(k, v)
			return doc
		}
	}

	doc.Doc = append(doc.Doc, primitive.E{Key: "$setOnInsert", Value: NewDocumentBuilder().Set(k, v)})
	return doc
}

// SetOnInsert is equivalent to !expr0 || !expr1 || !expr...
func SetOnInsert(k string, v interface{}) *Document {
	return NewDocumentBuilder().SetOnInsert(k, v)
}
