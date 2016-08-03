package dbtest

import (
	"github.com/CloudyKit/framework/database/change"
	"github.com/CloudyKit/framework/database/driver"
	"github.com/CloudyKit/framework/database/query"
	"github.com/CloudyKit/framework/database/scheme"

	"bytes"
	"fmt"
	"reflect"
)

var _ = driver.Driver(&FakeDriver{})

type FakeRecord map[string]reflect.Value
type fakeTable map[string]FakeRecord
type fakeDB map[string]fakeTable

func NewFakeDriver() *FakeDriver {
	return &FakeDriver{}
}

type FakeDriver struct {
	driver.NoTransactions
	ids   int
	db    fakeDB
	oplog bytes.Buffer

	PanicNew    bool
	PanicUpdate bool
	PanicDelete bool
	PanicModify bool
	PanicRemove bool
}

func (fk *FakeDriver) ResetOPLog() {
	fk.oplog.Reset()
}

func (fk *FakeDriver) OPLog() *bytes.Buffer {
	return &fk.oplog
}

func (d *FakeDriver) UseScheme(s *scheme.Scheme) error {
	return nil
}

func (d *FakeDriver) Search(s *scheme.Scheme, q *query.Query) driver.Result {
	return nil
}

func (d *FakeDriver) getTable(tableName string) fakeTable {
	if d.db == nil {
		d.db = make(fakeDB)
	}

	table, ok := d.db[tableName]
	if !ok {
		table = make(fakeTable)
		d.db[tableName] = table
	}
	return table
}

func (d *FakeDriver) printf(format string, v ...interface{}) {
	fmt.Fprintf(&d.oplog, format, v...)
}

func (d *FakeDriver) New(s *scheme.Scheme, keyField string, operations ...change.Set) (key string, err error) {
	if d.PanicNew {
		panic(fmt.Errorf("Panic on New enabled"))
	}

	d.ids++
	key = fmt.Sprint(d.ids)

	d.printf("INSERT: table(%s) key(%s)", s.Entity(), key)
	table := d.getTable(s.Entity())
	record := make(FakeRecord)

	for _, set := range operations {
		record[set.Field] = set.Value
		d.printf(" set(%s)=%q", set.Field, set.Value.Interface())
	}

	//record[keyField] = reflect.ValueOf(key)
	table[key] = record

	d.oplog.WriteString("\n")
	return
}

func (d *FakeDriver) getRecord(table, key string) (record FakeRecord, found bool) {
	record, found = d.getTable(table)[key]
	return
}

func (d *FakeDriver) Update(s *scheme.Scheme, keyField, key string, operations ...change.Set) (numofmodified int, err error) {
	if d.PanicUpdate {
		panic(fmt.Errorf("Panic on Update enabled"))
	}
	d.printf("UPDATE: table(%s) key(%s)", s.Entity(), key)

	record, found := d.getRecord(s.Entity(), key)
	if found {
		numofmodified++
		for _, set := range operations {
			record[set.Field] = set.Value
			d.printf(" set(%s)=%q", set.Field, set.Value.Interface())
		}
		//record[keyField] = reflect.ValueOf(key)
	} else {
		d.printf(" NOT FOUND")
	}
	d.oplog.WriteString("\n")
	return
}

func (d *FakeDriver) Modify(s *scheme.Scheme, q *query.Query, operations ...change.Operation) (numofmodified int, err error) {
	if d.PanicModify {
		panic(fmt.Errorf("Panic on Modify enabled"))
	}

	//record, found := d.getRecord(s.Entity(), primaryKey)
	//
	//if found {
	//	numofmodified++
	//	for _, op := range operations {
	//		switch op := op.(type) {
	//		case change.Set:
	//			record[op.Field] = op.Value
	//		default:
	//			err = errors.New("operation is not supported")
	//		}
	//	}
	//}

	return
}

func (d *FakeDriver) Delete(s *scheme.Scheme, keyField, key string) (numofmodified int, err error) {
	if d.PanicDelete {
		panic(fmt.Errorf("Panic on Delete enabled"))
	}

	_, found := d.getRecord(s.Entity(), key)
	if found {
		numofmodified++
		delete(d.getTable(s.Entity()), key)
	}
	return
}

func (d *FakeDriver) Remove(s *scheme.Scheme, q *query.Query) (numofmodified int, err error) {
	if d.PanicRemove {
		panic(fmt.Errorf("Panic on Remove enabled"))
	}

	return
}
