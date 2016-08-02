package database

import (
	"github.com/CloudyKit/framework/database/change"
	"github.com/CloudyKit/framework/database/query"
	"github.com/CloudyKit/framework/database/scheme"
)

type IDB interface {
	MigrateScheme(s *scheme.Scheme) error

	Begin() error
	Commit() error
	RowBack() error

	Search(s *scheme.Scheme, q *query.Query) query.Result

	Insert(s *scheme.Scheme, doc interface{}) (primaryKey string, err error)
	Update(s *scheme.Scheme, q *query.Query, doc interface{}) (numofmodified int, err error)

	Modify(s *scheme.Scheme, q *query.Query, operations ...change.Operation) (numofmodified int, err error)
	Delete(s *scheme.Scheme, q *query.Query) (numofmodified int, err error)
}
