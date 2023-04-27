package sequence

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/dynamic"
	"github.com/CloudyKit/framework/ensure"
	"github.com/CloudyKit/framework/event"
	"github.com/CloudyKit/framework/odm"
	"github.com/CloudyKit/framework/odm/events"
	"reflect"
	"strings"
)

var _ = WithSequences(app.Default.Registry, "sequences")

type Value int64
type controller struct {
	SequenceManager Manager
}

func (c *controller) findSequenceName(fieldData reflect.StructField) string {
	sequenceName := fieldData.Tag.Get("sequenceName")
	if sequenceName == "" {
		if fieldData.Name != "" {
			sequenceName = strings.ToLower(fieldData.Name)
		} else {
			sequenceName = "default"
		}
	}
	return sequenceName
}

func (c *controller) ensureSequence(r *container.Registry, field reflect.StructField, collection string, seq Number) Number {
	increment, err := c.SequenceManager.Increment(r, c.findSequenceName(field), collection, seq)
	ensure.NilErr(err)
	return increment
}

func WithSequences(r *container.Registry, collectionName string) error {

	dispatcher := event.GetDispatcher(r)

	_ = odm.WithCollectionManager(r, Manager{}, collectionName)

	dispatcher.Subscribe(events.InsertOneKey, func(insertOneEvent *events.InsertOneEvent) {
		controller := &controller{}
		insertOneEvent.Registry().Inject(controller)

		dynamic.StructVisitor(insertOneEvent.Document, func(seq Number, field reflect.StructField) Number {
			return controller.ensureSequence(insertOneEvent.Registry(), field, insertOneEvent.CollectionName, seq)
		})
	})

	dispatcher.Subscribe(events.InsertManyKey, func(insertManyEvent *events.InsertManyEvent) {

		controller := &controller{}
		insertManyEvent.Registry().Inject(controller)

		for _, document := range insertManyEvent.Documents {
			dynamic.StructVisitor(document, func(seq Number, field reflect.StructField) Number {
				return controller.ensureSequence(insertManyEvent.Registry(), field, insertManyEvent.CollectionName, seq)
			})
		}
	})

	return nil
}
