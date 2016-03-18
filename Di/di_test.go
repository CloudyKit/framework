package Di_test

import (
	"github.com/CloudyKit/framework/Di"
	"testing"
)

var context = Di.New()

func TestDi(t *testing.T) {

	newContext := context.Child()
	defer newContext.Done()

	var tt struct {
		*testing.T
	}

	newContext.Put(t)
	newContext.Inject(&tt)

	if tt.T == nil {
		t.Fail()
	} else {
		tt.Log("Injector is working")
	}
}

func TestDiFromParent(t *testing.T) {
	context.Put(t)
	newContext := context.Child()
	defer newContext.Done()
	var tt struct {
		*testing.T
	}
	newContext.Inject(&tt)
	if tt.T == nil {
		newContext.Put(t)
		newContext.Inject(&tt)
		if tt.T != nil {
			t.Log("Found in current scope")
		}
		t.Fail()
	} else {
		tt.Log("Injector is working")
	}
}

func BenchmarkInject(b *testing.B) {
	var tt struct {
		*testing.B
	}
	for i := 0; i < b.N; i++ {
		context.Put(b)
		context.Inject(&tt)
	}
}

func BenchmarkInjectChild(b *testing.B) {
	var tt struct {
		*testing.B
	}
	for i := 0; i < b.N; i++ {
		context := context.Child()
		context.Put(b)
		context.Inject(&tt)
		context.Done()
	}
}
