package sequence

import (
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/ensure"
	"github.com/CloudyKit/framework/odm"
	"github.com/CloudyKit/framework/odm/bsoner"
	"github.com/CloudyKit/framework/odm/multitenant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Number int64
type Manager struct {
	DB                 odm.Manager
	MultiTenantManager *multitenant.Manager
}

type Document struct {
	Name           string `bson:"name"`
	CollectionName string `bson:"collectionName"`
	Number         Number `bson:"number,omitempty" `
	TenantData     bson.M `bson:"tenantData"`
}

func (m *Manager) Increment(r *container.Registry, name, collection string, seq Number) (Number, error) {

	counter := m.sequenceDocument(r, name, collection)

	if seq > 0 {

		// todo implement better solution here
		result := m.DB.FindOneAndUpdate(
			counter,
			bsoner.NewDocumentBuilder().Set("$setOnInsert", bsoner.NewDocumentBuilder().Set("number", seq)),
			options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
		)

		ensure.NilErr(result.Decode(counter))
		if counter.Number < seq {
			_, err := m.DB.UpdateOne(bsoner.And(
				bsoner.NewFilterBuilder().Set("name", counter.Name).
					Set("collectionName", counter.CollectionName).
					Set("tenantData", counter.TenantData),
				bsoner.NewFilterBuilder().Lt("number", seq),
			), bsoner.NewDocumentBuilder().Set("$set", bsoner.NewDocumentBuilder().Set("number", seq)))
			if err != nil {
				return seq, err
			}
		}

		return seq, result.Err()
	}

	result := m.DB.FindOneAndUpdate(
		counter,
		bsoner.NewDocumentBuilder().Set("$inc", bsoner.NewDocumentBuilder().Set("number", 1)),
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)

	ensure.NilErr(result.Decode(counter))

	return counter.Number, result.Err()
}

func (m *Manager) sequenceDocument(r *container.Registry, name string, collection string) *Document {
	var sequenceDocument = &Document{
		Name:           name,
		CollectionName: collection,
	}
	for i := 0; i < len(m.MultiTenantManager.Handlers); i++ {
		handler := m.MultiTenantManager.Handlers[i]
		if handler.Collections[collection] {
			if sequenceDocument.TenantData == nil {
				sequenceDocument.TenantData = map[string]interface{}{}
			}
			sequenceDocument.TenantData[handler.FieldName] = handler.Resolver.GetCurrentTenantKey(r)
		}
	}
	return sequenceDocument
}
