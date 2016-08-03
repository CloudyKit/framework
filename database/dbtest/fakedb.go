package dbtest

import (
	"github.com/CloudyKit/framework/database"
	"github.com/CloudyKit/framework/database/scheme"
	"reflect"
)

type FakeDB struct {
	*database.DB
	FakeDriver *FakeDriver
}

func NewFakeDB() *FakeDB {
	driver := NewFakeDriver()
	return &FakeDB{FakeDriver: driver, DB: database.New(driver)}
}

func (db *FakeDB) Reset() *FakeDB {
	db.FakeDriver = NewFakeDriver()
	db.DB = database.New(db.FakeDriver)
	return db
}

func (db *FakeDB) ResetOPLog() *FakeDB {
	db.FakeDriver.ResetOPLog()
	return db
}

func (db *FakeDB) OPLogString() string {
	return db.FakeDriver.oplog.String()
}

func (db *FakeDB) OPLogExpect(_want string) (w string, g string) {
	w = _want
	g = db.OPLogString()
	return
}

func (db *FakeDB) Diff(s *scheme.Scheme, key string, expect FakeRecord) (diff, extraInExpect, extraInRecord FakeRecord) {
	diff, extraInExpect, extraInRecord = make(FakeRecord), make(FakeRecord), make(FakeRecord)
	record, found := db.FakeDriver.getRecord(s.Entity(), key)
	if !found {
		extraInExpect = expect
		return
	}

	for field, value := range expect {
		if v, found := record[field]; found {
			if !reflect.DeepEqual(v.Interface(), value.Interface()) {
				diff[field] = v
			}
		} else {
			extraInExpect[field] = value
		}
	}

	for field, value := range record {
		if _, found := expect[field]; !found {
			extraInRecord[field] = value
		}
	}
	return
}

func (db *FakeDB) Expect(s *scheme.Scheme, key string, expect FakeRecord) FakeRecord {
	diff, extraInExpect, extraInRecord := db.Diff(s, key, expect)

	nDiff := make(FakeRecord)

	for k, v := range diff {
		nDiff["!"+k] = v
	}

	for k, v := range extraInExpect {
		nDiff["-"+k] = v
	}

	for k, v := range extraInRecord {
		nDiff["+"+k] = v
	}

	return nDiff
}
