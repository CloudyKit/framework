package events

import (
	"github.com/CloudyKit/framework/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	InsertOneKey         = "odm.Manager.InsertOne"
	InsertManyKey        = "odm.Manager.InsertMany"
	DeleteOneKey         = "odm.Manager.DeleteOne"
	DeleteManyKey        = "odm.Manager.DeleteMany"
	UpdateOneKey         = "odm.Manager.UpdateOne"
	UpdateManyKey        = "odm.Manager.UpdateMany"
	ReplaceOneKey        = "odm.Manager.ReplaceOne"
	CountDocumentsKey    = "odm.Manager.CountDocuments"
	DistinctKey          = "odm.Manager.Distinct"
	FindOneKey           = "odm.Manager.FindOne"
	FindManyKey          = "odm.Manager.FindMany"
	FindOneAndDeleteKey  = "odm.Manager.FindOneAndDelete"
	FindOneAndReplaceKey = "odm.Manager.FindOneAndReplace"
	FindOneAndUpdateKey  = "odm.Manager.FindOneAndUpdate"
	FilterKey            = "odm.Manager.Filter"

	UpdateKeys = UpdateOneKey + "|" + UpdateManyKey
	DeleteKeys = DeleteOneKey + "|" + DeleteManyKey
)

type InsertOneEvent struct {
	event.Event
	CollectionName string
	Options        []*options.InsertOneOptions
	Document       interface{}
}

type InsertManyEvent struct {
	event.Event
	CollectionName string
	Documents      []interface{}
	Options        []*options.InsertManyOptions
}

type DeleteEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.DeleteOptions
}

type UpdateEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Document       interface{}
	Options        []*options.UpdateOptions
}

type ReplaceEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Document       interface{}
	Options        []*options.ReplaceOptions
}

type CountDocumentsEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.CountOptions
}

type DistinctEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	FieldName      string
	Options        []*options.DistinctOptions
}

type FindManyEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.FindOptions
}

type FindOneEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.FindOneOptions
}

type FindOneAndDeleteEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.FindOneAndDeleteOptions
}

type FindOneAndReplaceEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.FindOneAndReplaceOptions
	Document       interface{}
}

type FindOneAndUpdateEvent struct {
	event.Event
	CollectionName string
	Filter         interface{}
	Options        []*options.FindOneAndUpdateOptions
	Document       interface{}
}

type FilterEvent struct {
	event.Event
	Filter         interface{}
	CollectionName string
}
