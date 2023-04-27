package odm

import (
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/event"
	"github.com/CloudyKit/framework/odm/events"
)

type Trait interface {
	OnFilter(filter interface{}) interface{}
	OnInsert(doc interface{}) interface{}
	OnUpdateOrReplace(filter, doc interface{}) interface{}
	OnDelete(filter interface{}) interface{}
}

func WithTraits(registry *container.Registry, traits ...Trait) error {

	dispatcher := event.GetDispatcher(registry)

	dispatcher.Subscribe(events.FilterKey, func(event *events.FilterEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Filter = trait.OnFilter(event.Filter)
		}
	})

	dispatcher.Subscribe(events.InsertOneKey, func(event *events.InsertOneEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Document = trait.OnInsert(event.Document)
		}
	})

	dispatcher.Subscribe(events.InsertManyKey, func(event *events.InsertManyEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			for i := 0; i < len(event.Documents); i++ {
				event.Documents[i] = trait.OnInsert(event.Documents[i])
			}
		}
	})

	dispatcher.Subscribe(events.ReplaceOneKey, func(event *events.ReplaceEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Document = trait.OnUpdateOrReplace(event.Filter, event.Document)
		}
	})

	dispatcher.Subscribe(events.FindOneAndReplaceKey, func(event *events.FindOneAndReplaceEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Document = trait.OnUpdateOrReplace(event.Filter, event.Document)
		}
	})

	dispatcher.Subscribe(events.FindOneAndUpdateKey, func(event *events.FindOneAndUpdateEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Document = trait.OnUpdateOrReplace(event.Filter, event.Document)
		}
	})

	dispatcher.Subscribe(events.UpdateKeys, func(event *events.UpdateEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Document = trait.OnUpdateOrReplace(event.Filter, event.Document)
		}
	})

	dispatcher.Subscribe(events.DeleteKeys, func(event *events.DeleteEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Filter = trait.OnDelete(event.Filter)
		}
	})

	dispatcher.Subscribe(events.FindOneAndDeleteKey, func(event *events.FindOneAndDeleteEvent) {
		for i := 0; i < len(traits); i++ {
			trait := traits[i]
			event.Filter = trait.OnDelete(event.Filter)
		}
	})

	return nil
}
