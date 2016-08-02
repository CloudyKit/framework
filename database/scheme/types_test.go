package scheme

import (
	"reflect"
	"testing"
)

type typeTest struct {
	Type              Type
	Val, ValConverted interface{}
}

func TestValidType(t *testing.T) {
	var _typeTest = []typeTest{
		{Type: Int{}, Val: "12", ValConverted: 12},
		{Type: Int{}, Val: uint(12), ValConverted: 12},
	}

	for _, test := range _typeTest {
		val, err := test.Type.Value(reflect.ValueOf(test.Val))
		if err != nil {
			t.Errorf("scheme type %s failed to convert %s", reflect.TypeOf(test.Type), reflect.TypeOf(test.Val))
		} else if !reflect.DeepEqual(test.ValConverted, val.Interface()) {
			t.Errorf("unexpected value returned from conversion with %s to %s  want: %v got: %v",
				reflect.TypeOf(test.Type),
				reflect.TypeOf(test.Val),
				val,
				test.ValConverted,
			)
		}
	}
}
