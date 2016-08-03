package driver

import (
	"errors"
	"github.com/CloudyKit/framework/database/change"
	"github.com/CloudyKit/framework/database/query"
	"github.com/CloudyKit/framework/database/scheme"
)

type NoTransactions struct{}

func (NoTransactions) Begin() error {
	return errors.New("Transactions are not supported by this driver")
}
func (NoTransactions) Commit() error {
	return errors.New("Transactions are not supported by this driver")
}
func (NoTransactions) RowBack() error {
	return errors.New("Transactions are not supported by this driver")
}

type Result interface {
	NumOfRecords() int
	ScanRow(v ...interface{}) error
}

type Driver interface {
	UseScheme(s *scheme.Scheme) error

	Begin() error
	Commit() error
	RowBack() error

	Search(s *scheme.Scheme, q *query.Query) Result

	New(s *scheme.Scheme, primaryKeyField string, operations ...change.Set) (primaryKey string, err error)
	Update(s *scheme.Scheme, primaryKeyField, primaryKey string, operations ...change.Set) (numofmodified int, err error)

	Modify(s *scheme.Scheme, q *query.Query, operations ...change.Operation) (numofmodified int, err error)
	Remove(s *scheme.Scheme, q *query.Query) (numofmodified int, err error)
}
