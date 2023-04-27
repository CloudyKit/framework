package dynamic

import (
	"reflect"
	"testing"
)

type SpecialValue struct {
	Int int
}
type SpecialValue2 struct {
	Int2 int
}
type SpecialValue3 struct {
	Int3 int
}

type TestStruct struct {
	SpecialValue
	Special2 SpecialValue2
	special3 SpecialValue3
}

func TestStructWalk(t *testing.T) {

	myStruct := &TestStruct{}

	StructVisitor(myStruct, func(value SpecialValue, _ reflect.StructField) SpecialValue {
		value.Int++
		return value
	})

	StructVisitor(myStruct, func(value *SpecialValue2, _ reflect.StructField) *SpecialValue2 {
		value.Int2++
		return value
	})

	StructVisitor(myStruct, func(value SpecialValue3, _ reflect.StructField) SpecialValue3 {
		value.Int3++
		return value
	})

	if myStruct.SpecialValue.Int == 0 {
		t.Fail()
	}

	if myStruct.Special2.Int2 == 0 {
		t.Fail()
	}

}

func TestPropertyGet(t *testing.T) {
	myStruct := &TestStruct{
		SpecialValue: SpecialValue{
			Int: 10,
		},
		Special2: SpecialValue2{
			Int2: 20,
		},
		special3: SpecialValue3{
			Int3: 30,
		},
	}
	propertyGetDefault, found := PropertyGetDefault(myStruct, "Int", -1)
	if !found {
		t.Fail()
	}
	if propertyGetDefault != 10 {
		t.Fail()
	}

	propertyGetDefault, found = PropertyGetDefault(myStruct, "SpecialValue.Int", -1)
	if !found {
		t.Fail()
	}
	if propertyGetDefault != 10 {
		t.Fail()
	}

	propertyGetDefault, found = PropertyGetDefault(myStruct, "Special2.Int2", -1)
	if !found {
		t.Fail()
	}
	if propertyGetDefault != 20 {
		t.Fail()
	}
}
