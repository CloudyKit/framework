package scope

import (
	"testing"
)

type testHolder struct {
	*testing.T
}

func TestDiProvideAndInject(t *testing.T) {

	// creates a new DI context
	var newContext = New()
	defer newContext.End()

	var tt testHolder

	// Provide t *testing.T into newContext
	newContext.Map(t)
	// Injects t *testing.T into testHolder struct{ *testing.T }
	newContext.Inject(&tt)
	// check value was injected successfully
	if tt.T != t {
		t.Fatal("Fail to inject value from the current context")
	}

	// creates a child context
	newChildContext := newContext.Inherit()
	defer newChildContext.End()

	// reset tt value
	tt = testHolder{}

	// Injects t *testing.T into testHolder struct{ *testing.T } again now using child context
	newChildContext.Inject(&tt)
	// check value was injected successfully
	if tt.T != t {
		t.Fatal("Fail to inject value from the child context")
	}

	var empty = New()
	defer empty.End()
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
	var childContext = context.Inherit()
	if context.references != 1 {
		t.Fatal("Inválid reference counting ", context.references)
	}
	if childContext.references != 0 {
		t.Fatal("Inválid reference counting ", childContext.references)
	}

	childContext.End()
	if childContext.parent != nil {
		t.Fatal("Inválid reference counting ", childContext.references)
	}
	if context.references != 0 {
		t.Fatal("Inválid reference counting ", context.references)
	}

	context.End()
	if context.parent != nil {
		t.Fatal("Inválid reference counting ", context.references)
	}

}

func TestGlobal_Checkpoint(t *testing.T) {
	var context = New()
	var contextcopy = context
	(func() {
		defer Checkpoint(&context)()
		if context == contextcopy {
			t.Fail()
		}
	})()

	if context != contextcopy {
		t.Fail()
	}
}

var context = New()

func BenchmarkInject(b *testing.B) {
	var tt struct {
		*testing.B
	}
	for i := 0; i < b.N; i++ {
		context.Map(b)
		context.Inject(&tt)
	}
}

func BenchmarkInjectChild(b *testing.B) {
	var tt struct {
		*testing.B
	}
	for i := 0; i < b.N; i++ {
		context := context.Inherit()
		context.Map(b)
		context.Inject(&tt)
		context.End()
	}
}
