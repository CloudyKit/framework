package database

import (
	"github.com/CloudyKit/framework/database/driver"
	"github.com/CloudyKit/framework/database/query"
	"github.com/CloudyKit/framework/database/scheme"
)

type Result struct {
	q *query.Query
	s *scheme.Scheme
	r driver.Result
}

func (r *Result) NumOfRecords() int {
	return r.r.NumOfRecords()
}

func (r *Result) Fetch(target interface{}) error {
	return nil
}

func (r *Result) FetchNext(target interface{}) bool {
	return true
}

func (r *Result) FetchAll(target interface{}) error {
	return nil
}
