package database_test

import (
	"github.com/CloudyKit/framework/database/dbtest"
	"github.com/CloudyKit/framework/database/model"
	"github.com/CloudyKit/framework/database/scheme"
	"github.com/CloudyKit/framework/validation"
	"testing"
)

var (
	CompanyScheme  = scheme.New("companies", "ID")
	EmployeeScheme = scheme.New("employess", "ID")

	_ = scheme.Init(CompanyScheme, func(def *scheme.Def) {

		def.Field("Name", scheme.String{}, validation.MinLength("Please enter a válid name", 5))
		def.RefChildren("Employees", EmployeeScheme, "CompanyID")

	})

	_ = scheme.Init(EmployeeScheme, func(def *scheme.Def) {

		def.Field("Name", scheme.String{}, validation.MinLength("Please enter a válid name", 5))
		def.RefParent("Company", CompanyScheme, "CompanyID")

	})
)

type (
	Employee struct {
		model.Model

		Company *Company
		Name    string
	}

	Company struct {
		model.Model
		Name string

		Employees []*Employee
	}
)

func TestDB_SaveSimple(t *testing.T) {

	company := &Company{
		Name: "Cloudy Kit",
		Employees: []*Employee{
			{Name: "Henrique"},
			{Name: "Henrique"},
			{Name: "Henrique"},
			{Name: "Henrique"},
		},
	}

	fakeDB := dbtest.NewFakeDB()
	fakeDB.FakeDriver.PanicUpdate = true

	result, err := fakeDB.Save(CompanyScheme, company)
	if err != nil {
		t.Fatalf("Something bad happend can't insert the records to the database, err: %s", err)
	}

	if result.Bad() {
		t.Fatalf("Something bad happend can't insert the records to the database, validation: %v", result)
	}

	got := fakeDB.FakeDriver.OPLog().String()
	want := "INSERT: table(companies) key(1) set(Name)=\"Cloudy Kit\"\nINSERT: table(employess) key(2) set(Name)=\"Henrique\" set(CompanyID)=\"1\"\nINSERT: table(employess) key(3) set(Name)=\"Henrique\" set(CompanyID)=\"1\"\nINSERT: table(employess) key(4) set(Name)=\"Henrique\" set(CompanyID)=\"1\"\nINSERT: table(employess) key(5) set(Name)=\"Henrique\" set(CompanyID)=\"1\"\n"
	if got != want {
		t.Fatalf("OPLog mismatch want:\n%s\ngot:\n%s", want, got)
	}

	fakeDB.FakeDriver.PanicUpdate = false
	fakeDB.ResetOPLog()
	for _, e := range company.Employees {
		if e.Company != company {
			t.Error("Employee Company should be pointing to the parent Company")
		}
		e.Company = nil
		fakeDB.Save(EmployeeScheme, e)
	}

	want, got = fakeDB.OPLogExpect("")
	if want != got {
		t.Fatalf("OPLog mismatch want:\n%s\ngot:\n%s", want, got)
	}
}
