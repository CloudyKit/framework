package Validator

import "testing"

type User struct {
	FirstName string
	LastName  string
	Email     string
}

func runValidation(u interface{}) Result {
	return New(u).
		At("FirstName",
			NoEmpty("First Name can't be empty")).
		At("LastName",
			NoEmpty("Last Name can not be empty")).
		At("Email",
			NoEmpty("Email can not be empty"),
			Email("Please enter a valid email address")).
		Done()
}

func TestSimpleValidation(t *testing.T) {

	var testData = []struct {
		d interface{}
		v func(interface{}) Result
		e bool
	}{
		{User{}, runValidation, true},
		{User{
			FirstName: "TestName",
			LastName:  "LastName",
			Email:     "email@gmailcom",
		}, runValidation, true},
		{User{
			FirstName: "TestName",
			LastName:  "LastName",
			Email:     "email+name@gmail.com",
		}, runValidation, false},
		{User{
			FirstName: "TestName",
			LastName:  "LastName",
			Email:     "emailname@gmail.com",
		}, runValidation, false},
	}
	for key, value := range testData {
		errs := value.v(value.d)
		if value.e == (len(errs) == 0) {
			t.Errorf("Test:%v ExpectErr:%v Struct:%#v Errors:%v", key, value.e, value.d, errs)
		}
	}
}
