package multitenant

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/event"
	"github.com/CloudyKit/framework/odm/bsoner"
	"github.com/CloudyKit/framework/odm/events"
	"github.com/CloudyKit/framework/registry"
)

var defaultManager = &Manager{}

var _ = registry.Provider(app.Default.Registry, func(r *container.Registry) *Manager {
	defaultManager.SubscribeEvents(r)
	r.WithValues(defaultManager)
	return defaultManager
})

type Handler struct {
	FieldName   string
	Collections map[string]bool
	Resolver    Resolver
}

type Manager struct {
	eventsRegistered bool
	Handlers         []*Handler
}

type Resolver interface {
	GetCurrentTenantKey(registry *container.Registry) any
	SetCurrentTenantKey(registry *container.Registry, document any)
}

func (m *Manager) WithTenantKey(tenantResolver Resolver, fieldName string, collectionNames ...string) *Manager {
	handler := &Handler{
		FieldName: fieldName,
		Resolver:  tenantResolver,
	}
	if len(collectionNames) > 0 {
		handler.Collections = map[string]bool{}
		if len(collectionNames) != 0 {
			for i := 0; i < len(collectionNames); i++ {
				handler.Collections[collectionNames[i]] = true
			}
		}
	}

	m.Handlers = append(m.Handlers, handler)
	return m
}

func (m *Manager) SubscribeEvents(r *container.Registry) {
	if !m.eventsRegistered {
		m.eventsRegistered = true
		dispatcher := registry.Get[*event.Dispatcher](r)

		dispatcher.Subscribe(events.FilterKey, func(filterEvent *events.FilterEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				if len(handler.Collections) == 0 || handler.Collections[filterEvent.CollectionName] {
					companyFilter := bsoner.Eq(handler.FieldName, handler.Resolver.GetCurrentTenantKey(filterEvent.Registry()))
					if filterEvent.Filter == nil {
						filterEvent.Filter = companyFilter
					} else {
						filterEvent.Filter = bsoner.And(
							filterEvent.Filter,
							companyFilter,
						)
					}
				}
			}
		})

		dispatcher.Subscribe(events.InsertOneKey, func(insertOneEvent *events.InsertOneEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				handler.Resolver.SetCurrentTenantKey(insertOneEvent.Registry(), insertOneEvent.Document)
			}
		})

		dispatcher.Subscribe(events.InsertManyKey, func(insertManyEvent *events.InsertManyEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				for _, document := range insertManyEvent.Documents {
					handler.Resolver.SetCurrentTenantKey(insertManyEvent.Registry(), document)
				}
			}
		})

		dispatcher.Subscribe(events.ReplaceOneKey, func(replaceEvent *events.ReplaceEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				handler.Resolver.SetCurrentTenantKey(replaceEvent.Registry(), replaceEvent.Document)
			}
		})

		dispatcher.Subscribe(events.UpdateOneKey, func(replaceEvent *events.UpdateEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				handler.Resolver.SetCurrentTenantKey(replaceEvent.Registry(), replaceEvent.Document)
			}
		})

		dispatcher.Subscribe(events.UpdateManyKey, func(replaceEvent *events.UpdateEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				handler.Resolver.SetCurrentTenantKey(replaceEvent.Registry(), replaceEvent.Document)
			}
		})

		dispatcher.Subscribe(events.FindOneAndReplaceKey, func(replaceEvent *events.FindOneAndReplaceEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				handler.Resolver.SetCurrentTenantKey(replaceEvent.Registry(), replaceEvent.Document)
			}
		})

		dispatcher.Subscribe(events.FindOneAndUpdateKey, func(replaceEvent *events.FindOneAndUpdateEvent) {
			for i := 0; i < len(m.Handlers); i++ {
				handler := m.Handlers[i]
				handler.Resolver.SetCurrentTenantKey(replaceEvent.Registry(), replaceEvent.Document)
			}
		})

	}
}

func WithTenantKey(r *container.Registry, tenantResolver Resolver, fieldName string, collectionNames ...string) *Manager {
	manager := registry.Get[*Manager](r)
	return manager.WithTenantKey(tenantResolver, fieldName, collectionNames...)
}
