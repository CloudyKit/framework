// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package container

import (
	"testing"
)

type testHolder struct {
	*testing.T
}

func TestDiProvideAndInject(t *testing.T) {

	// creates a new DI context
	var newContext = New()
	defer newContext.Dispose()

	var tt testHolder

	// Provide t *testing.T into newContext
	newContext.WithValues(t)
	// Injects t *testing.T into testHolder struct{ *testing.T }
	newContext.Inject(&tt)
	// check value was injected successfully
	if tt.T != t {
		t.Fatal("Fail to inject value from the current context")
	}

	// creates a child context
	newChildContext := newContext.Fork()
	defer newChildContext.Dispose()

	// reset tt value
	tt = testHolder{}

	// Injects t *testing.T into testHolder struct{ *testing.T } again now using child context
	newChildContext.Inject(&tt)
	// check value was injected successfully
	if tt.T != t {
		t.Fatal("Fail to inject value from the child context")
	}

	var empty = New()
	defer empty.Dispose()
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
	var childContext = context.Fork()
	if context.references != 1 {
		t.Fatal("Inválid reference counting ", context.references)
	}
	if childContext.references != 0 {
		t.Fatal("Inválid reference counting ", childContext.references)
	}

	childContext.Dispose()
	if childContext.parent != nil {
		t.Fatal("Inválid reference counting ", childContext.references)
	}
	if context.references != 0 {
		t.Fatal("Inválid reference counting ", context.references)
	}

	context.Dispose()
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
		context.WithValues(b)
		context.Inject(&tt)
	}
}

func BenchmarkInjectChild(b *testing.B) {
	var tt struct {
		*testing.B
	}
	for i := 0; i < b.N; i++ {
		context := context.Fork()
		context.WithValues(b)
		context.Inject(&tt)
		context.Dispose()
	}
}
