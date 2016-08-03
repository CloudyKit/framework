package database

import (
	"github.com/CloudyKit/framework/database/change"
	"github.com/CloudyKit/framework/database/driver"
	"github.com/CloudyKit/framework/database/model"
	"github.com/CloudyKit/framework/database/query"
	"github.com/CloudyKit/framework/database/scheme"
	"github.com/CloudyKit/framework/validation"

	"reflect"
)

type Model interface {
	model.IModel
}

func New(driver driver.Driver) *DB {
	db := new(DB)
	db.driver = driver
	db.dbID = reflect.ValueOf(db).Pointer()
	return db
}

type DB struct {
	dbID        uintptr
	operationID uintptr
	driver      driver.Driver
}

func (db *DB) Begin() error {
	return db.driver.Begin()
}

func (db *DB) Commit() error {
	return db.driver.Commit()
}

func (db *DB) RowBack() error {
	return db.driver.RowBack()
}

func (db *DB) Search(s *scheme.Scheme, q *query.Query) Result {
	return Result{s: s, q: q, r: db.driver.Search(s, q)}
}

func (db *DB) Load(s *scheme.Scheme, key string, m Model) error {
	return nil
}

func (db *DB) Save(s *scheme.Scheme, m Model) (validation.Result, error) {

	md := model.GetModelData(m)
	md.Scheme = s

	v, err := db.executeSave(m, nil, "", reflect.ValueOf(m), nil)
	return v.Done(), err
}

func (db *DB) Modify(s *scheme.Scheme, q *query.Query, changes ...change.Operation) (int, error) {
	return db.driver.Modify(s, q, changes...)
}

func (db *DB) Remove(s *scheme.Scheme, q *query.Query) (int, error) {
	return db.driver.Remove(s, q)
}

func (db *DB) Driver() driver.Driver {
	return db.driver
}
