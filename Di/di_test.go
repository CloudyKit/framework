package Di

import (
	"testing"
)

type testHolder struct {
	*testing.T
}

func TestDiProvideAndInject(t *testing.T) {

	// creates a new DI context
	var newContext = New()
	defer newContext.Done()

	var tt testHolder

	// Provide t *testing.T into newContext
	newContext.Put(t)
	// Injects t *testing.T into testHolder struct{ *testing.T }
	newContext.Inject(&tt)
	// check value was injected successfully
	if tt.T != t {
		t.Fatal("Fail to inject value from the current context")
	}

	// creates a child context
	newChildContext := newContext.Child()
	defer newChildContext.Done()

	// reset tt value
	tt = testHolder{}

	// Injects t *testing.T into testHolder struct{ *testing.T } again now using child context
	newChildContext.Inject(&tt)
	// check value was injected successfully
	if tt.T != t {
		t.Fatal("Fail to inject value from the child context")
	}

	var empty = New()
	defer empty.Done()
	tt = testHolder{}
	if tt.T != nil {
		t.Fail()
	}
}

func TestDiDone(t *testing.T) {
	var context = New()
	if context.references != 0 {
		t.Fatal("Inválid reference counting ", context.references)
	}
	var childContext = context.Child()
	if context.references != 1 {
		t.Fatal("Inválid reference counting ", context.references)
	}
	if childContext.references != 0 {
		t.Fatal("Inválid reference counting ", childContext.references)
	}

	childContext.Done()
	if childContext.parent != nil {
		t.Fatal("Inválid reference counting ", childContext.references)
	}
	if context.references != 0 {
		t.Fatal("Inválid reference counting ", context.references)
	}

	context.Done()
	if context.parent != nil {
		t.Fatal("Inválid reference counting ", context.references)
	}
}

var context = New()

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
