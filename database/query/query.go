package query

type Result struct {
	count int
}

func (r *Result) NumOfRecords() int {
	return r.count
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

type Query struct {
}
