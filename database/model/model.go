package model

import (
	"github.com/CloudyKit/framework/database/scheme"
)

type mark int

const (
	WasCreated mark = 1 + iota
	WasRemoved
	WasUpdated
	WasLoaded
	WasSetted
)

type modelData struct {
	Key        string
	Scheme     *scheme.Scheme
	Mark       mark
	LoadedData map[string]interface{}
}

type Model struct {
	data modelData
}

type IModel interface {
	modelData() *modelData
	Scheme() *scheme.Scheme
	Key() string
	SetKey(string)
}

// GetModelData returns the model data this func is not intended to be used by the
// user, the model data holds
func GetModelData(m IModel) *modelData {
	if m == nil {
		return nil
	}
	return m.modelData()
}

func (m *Model) SetKey(key string) {
	m.data.Key = key
	m.data.Mark = WasSetted
}

func (m *Model) Key() string {
	return m.data.Key
}

func (m *Model) Scheme() *scheme.Scheme {
	return m.data.Scheme
}

func (model *Model) modelData() *modelData {
	if model == nil {
		return nil
	}
	return &model.data
}
