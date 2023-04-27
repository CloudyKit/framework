package odm

import (
	"context"
	"fmt"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

var ClientType = reflect.TypeOf((*mongo.Client)(nil))
var DatabaseType = reflect.TypeOf((*mongo.Database)(nil))
var CollectionType = reflect.TypeOf((*mongo.Collection)(nil))

func WithCollectionManager(registry *container.Registry, manager interface{}, collectionName string) error {

	typeOf, ok := manager.(reflect.Type)

	if !ok {
		typeOf = reflect.TypeOf(manager)
	}

	if typeOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("type %q is not of kind Struct", typeOf))
	}

	registry.MapInitializerFunc(typeOf, func(registry *container.Registry, value reflect.Value) {
		db := CurrentDatabase(registry)
		collection := registry.LoadType(CollectionType)                          // backups previous collection
		registry.WithTypeAndValue(CollectionType, db.Collection(collectionName)) // map new collection
		registry.InjectValue(value)                                              // injects new collection into the value
		registry.WithTypeAndValue(CollectionType, collection)                    // restore backed collection
	})

	return nil
}

func CurrentClient(registry *container.Registry) *mongo.Client {
	return registry.LoadType(ClientType).(*mongo.Client)
}

func CurrentDatabase(registry *container.Registry) *mongo.Database {
	return registry.LoadType(DatabaseType).(*mongo.Database)
}

func CurrentCollection(cdi *container.Registry) *mongo.Collection {
	return cdi.LoadType(CollectionType).(*mongo.Collection)
}

type Component struct {
	Database      string
	ClientOptions []*options.ClientOptions
	Client        *mongo.Client
}

func (component *Component) client() (*mongo.Client, error) {
	return mongo.Connect(context.TODO(), component.ClientOptions...)
}

func (component *Component) Bootstrap(app *app.Kernel) {
	app.Registry.WithTypeAndProviderFunc(ClientType, func(registry *container.Registry) interface{} {
		if component.Client == nil {
			var err error
			component.Client, err = component.client()
			if err != nil {
				panic(err)
			}
		}
		return component.Client
	})
	app.Registry.WithTypeAndProviderFunc(DatabaseType, func(registry *container.Registry) interface{} {
		return CurrentClient(registry).Database(component.Database)
	})
}

func NewComponent(databaseName string, options ...*options.ClientOptions) *Component {
	return &Component{Database: databaseName, ClientOptions: options}
}
