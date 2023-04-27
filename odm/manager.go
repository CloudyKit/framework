package odm

import (
	"context"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/event"
	"github.com/CloudyKit/framework/odm/bsoner"
	"github.com/CloudyKit/framework/odm/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ = container.Injectable(&Manager{})

type Manager struct {
	Context    context.Context
	Collection *mongo.Collection
	Registry   *container.Registry
}

func (m *Manager) BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return m.Collection.BulkWrite(m.Context, models, opts...)
}

func (m *Manager) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	payload := &events.InsertOneEvent{Document: document, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.InsertOneKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.InsertOne(m.Context, payload.Document, payload.Options...)
}

func (m *Manager) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	payload := &events.InsertManyEvent{Documents: documents, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.InsertManyKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.InsertMany(m.Context, payload.Documents, payload.Options...)
}

func (m *Manager) filterEvent(filter interface{}) interface{} {
	payload := &events.FilterEvent{Filter: filter, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.FilterKey, payload)
	if err != nil {
		return nil
	}
	return payload.Filter
}

func (m *Manager) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	payload := &events.DeleteEvent{Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.DeleteOneKey, payload)
	if err != nil {
		return nil, err
	}

	return m.Collection.DeleteOne(m.Context, payload.Filter, payload.Options...)
}

func (m *Manager) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	payload := &events.DeleteEvent{Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.DeleteManyKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.DeleteMany(m.Context, filter, opts...)
}

func (m *Manager) UpdateByID(id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.UpdateOne(primitive.D{{"_id", id}}, update, opts...)
}

func (m *Manager) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	payload := &events.UpdateEvent{Filter: m.filterEvent(filter), Document: update, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.UpdateOneKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.UpdateOne(m.Context, payload.Filter, payload.Document, payload.Options...)
}

func (m *Manager) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	payload := &events.UpdateEvent{Filter: m.filterEvent(filter), Document: update, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.UpdateManyKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.UpdateMany(m.Context, payload.Filter, payload.Document, payload.Options...)
}

func (m *Manager) ReplaceOne(filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	payload := &events.ReplaceEvent{Filter: m.filterEvent(filter), Document: replacement, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.ReplaceOneKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.ReplaceOne(m.Context, payload.Filter, payload.Document, payload.Options...)
}

func (m *Manager) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return m.Collection.Aggregate(m.Context, pipeline, opts...)
}

func (m *Manager) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	payload := &events.CountDocumentsEvent{Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.CountDocumentsKey, payload)
	if err != nil {
		return 0, err
	}
	return m.Collection.CountDocuments(m.Context, payload.Filter, payload.Options...)
}

func (m *Manager) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return m.Collection.EstimatedDocumentCount(m.Context, opts...)
}

func (m *Manager) Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	payload := &events.DistinctEvent{FieldName: fieldName, Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.DistinctKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.Distinct(m.Context, payload.FieldName, payload.Filter, payload.Options...)
}

func (m *Manager) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	payload := &events.FindManyEvent{Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.FindManyKey, payload)
	if err != nil {
		return nil, err
	}
	return m.Collection.Find(m.Context, payload.Filter, payload.Options...)
}

func (m *Manager) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	payload := &events.FindOneEvent{Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.FindOneKey, payload)
	if err != nil {
		return nil
	}
	return m.Collection.FindOne(m.Context, payload.Filter, opts...)
}

func (m *Manager) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	payload := &events.FindOneAndDeleteEvent{Filter: m.filterEvent(filter), Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.FindOneAndDeleteKey, payload)
	if err != nil {
		return nil
	}
	return m.Collection.FindOneAndDelete(m.Context, payload.Filter, payload.Options...)
}

func (m *Manager) FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	payload := &events.FindOneAndReplaceEvent{Filter: m.filterEvent(filter), Document: replacement, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.FindOneAndReplaceKey, payload)
	if err != nil {
		return nil
	}
	return m.Collection.FindOneAndReplace(m.Context, payload.Filter, payload.Document, payload.Options...)
}

func (m *Manager) FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	payload := &events.FindOneAndUpdateEvent{Filter: m.filterEvent(filter), Document: update, Options: opts, CollectionName: m.Collection.Name()}
	_, err := event.Dispatch(m.Registry, events.FindOneAndUpdateKey, payload)
	if err != nil {
		return nil
	}
	return m.Collection.FindOneAndUpdate(m.Context, payload.Filter, payload.Document, payload.Options...)
}

func (m *Manager) FindOneByID(id interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.FindOne(bsoner.ByID(id), opts...)
}

func (m *Manager) FindPaged(q, sort interface{}, page, size int) (*mongo.Cursor, error) {
	page--
	if page < 0 {
		page = 0
	}
	findOptions := options.Find().SetSkip(int64(page * size)).SetLimit(int64(size))
	if sort != nil {
		findOptions.SetSort(sort)
	}
	return m.Find(q, findOptions)
}

func (m *Manager) FindPagedInto(q, sort interface{}, page, size int, dest interface{}) error {
	cursor, err := m.FindPaged(q, sort, page, size)
	if err != nil {
		return err
	}
	return cursor.All(m.Context, dest)
}

func (m *Manager) EnsureIndex(indexName string, keys primitive.D, unique bool, opts ...*options.CreateIndexesOptions) error {

	type IndexModel struct {
		Keys primitive.D `bson:"keys"`
		Name string      `bson:"name"`
	}

	indexes := m.Collection.Indexes()
	var indexModels []IndexModel
	list, err := indexes.List(m.Context)
	if err != nil {
		return err
	}

	err = list.All(m.Context, &indexModels)
	if err != nil {
		return err
	}

	for _, index := range indexModels {
		if index.Name == indexName {
			if len(index.Keys) != len(keys) {
				_, err = indexes.DropOne(nil, indexName)
				if err != nil {
					return err
				}
				break
			}
			for i, key := range keys {
				if index.Keys[i].Key != key.Key || index.Keys[i].Value != key.Value {
					_, err = indexes.DropOne(nil, indexName)
					if err != nil {
						return err
					}
					break
				}
			}
			return nil
		}
	}

	_, err = indexes.CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options.MergeIndexOptions(options.Index().SetName(indexName).SetUnique(unique)),
	}, opts...)

	return err
}
