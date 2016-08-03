package scheme

import (
	"fmt"
	"reflect"
	"testing"
)

type typeTest struct {
	Type      Type
	Val, Want interface{}
	Err       bool
}

func RunTypeTest(t *testing.T, _typeTest []typeTest) {
	for _, test := range _typeTest {
		oldVal := reflect.ValueOf(test.Val)
		val, err := test.Type.Value(oldVal)
		valInterface := val.Interface()

		if test.Err {
			if err == nil {
				t.Errorf("\nHas: %s(%v) expect fail \nGot: %s(%v)", oldVal.Kind(), test.Val, val.Kind(), valInterface)
			}
		} else {
			typWant := reflect.TypeOf(test.Want)
			if err != nil {
				t.Errorf("scheme type %s failed to convert %s: %s", reflect.TypeOf(test.Type), reflect.TypeOf(test.Val), err)
			} else if val.Kind() != typWant.Kind() || fmt.Sprint(valInterface) != fmt.Sprint(test.Want) {
				t.Errorf("\nHas: %s(%v) \nGot: %s(%v) \nWant: %s(%v)", oldVal.Kind(), test.Val, val.Kind(), valInterface, typWant.Kind(), test.Want)
			}
		}
	}
}

func TestInt_Value(t *testing.T) {
	RunTypeTest(t, []typeTest{
		{Type: Int{}, Val: "12", Want: int64(12)},
		{Type: Int{}, Val: uint(12), Want: int64(12)},
		{Type: Int{}, Val: 12.3, Want: int64(12)},
		{Type: Int{}, Val: true, Want: int64(1)},
		{Type: Int{}, Val: false, Want: int64(0)},
	})
}

func TestUint_Value(t *testing.T) {
	RunTypeTest(t, []typeTest{
		{Type: Uint{}, Val: "12", Want: uint64(12)},
		{Type: Uint{}, Val: -12, Want: uint64(12), Err: true},
		{Type: Uint{}, Val: 12, Want: uint64(12)},
		{Type: Uint{}, Val: 12.3, Want: uint64(12)},
		{Type: Uint{}, Val: -12.3, Want: uint64(12), Err: true},
		{Type: Uint{}, Val: true, Want: uint64(1)},
		{Type: Uint{}, Val: false, Want: uint64(0)},
	})
}

func TestFloat_Value(t *testing.T) {
	RunTypeTest(t, []typeTest{
		{Type: Float{}, Val: "12", Want: float64(12)},
		{Type: Float{}, Val: "12.5", Want: float64(12.5)},
		{Type: Float{}, Val: ".5", Want: float64(.5)},
		{Type: Float{}, Val: "1.5", Want: float64(1.5)},
		{Type: Float{}, Val: -12, Want: float64(-12)},
		{Type: Float{}, Val: 12, Want: float64(12)},
		{Type: Float{}, Val: 12.3, Want: float64(12.3)},
		{Type: Float{}, Val: -12.3, Want: float64(-12.3)},
		{Type: Float{}, Val: true, Want: float64(1)},
		{Type: Float{}, Val: false, Want: float64(0)},
	})
}

type myString string

func TestString_Value(t *testing.T) {
	RunTypeTest(t, []typeTest{
		{Type: String{}, Val: "Galaxy S7", Want: "Galaxy S7"},
		{Type: String{}, Val: myString("My String Galaxy S7"), Want: "My String Galaxy S7"},
		{Type: String{}, Val: -12, Want: "-12"},
		{Type: String{}, Val: 12, Want: "12"},
		{Type: String{}, Val: 12.3, Want: "12.3"},
		{Type: String{}, Val: -12.3, Want: "-12.3"},
		{Type: String{}, Val: true, Want: "true"},
		{Type: String{}, Val: false, Want: "false"},
	})
}

func TestBool_Value(t *testing.T) {
	RunTypeTest(t, []typeTest{
		{Type: Bool{}, Val: "true", Want: true},
		{Type: Bool{}, Val: "false", Want: false},
		{Type: Bool{}, Val: myString("true"), Want: true},
		{Type: Bool{}, Val: "dtrue", Want: true, Err: true},
		{Type: Bool{}, Val: -12, Want: false},
		{Type: Bool{}, Val: 12, Want: true},
		{Type: Bool{}, Val: 12.3, Want: true},
		{Type: Bool{}, Val: -12.3, Want: false},
		{Type: Bool{}, Val: true, Want: true},
		{Type: Bool{}, Val: false, Want: false},
	})
}
