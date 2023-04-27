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

package event

import "testing"

type TestContext struct {
	Event
	Counter int
}

func TestManager_Emit(t *testing.T) {

	events := NewDispatcher()
	testcontext := new(TestContext)

	events.Subscribe("testRunning", func(tc *TestContext) {
		tc.Counter++
	})

	events.Subscribe("testNotRunning", func(tc *TestContext) {
		tc.Counter++
	})

	events.Dispatch(nil, "testRunning", testcontext)
	t.Logf("Counter %d", testcontext.Counter)

	if testcontext.Counter != 1 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}

	events.Subscribe("testRunning", func(tc *TestContext) {
		tc.Counter++
	})

	events.Dispatch(nil, "testRunning", testcontext)
	t.Logf("Counter %d", testcontext.Counter)
	if testcontext.Counter != 3 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}
}

func TestManager_EmitParent(t *testing.T) {

	testcontext := new(TestContext)

	events := NewDispatcher()
	events.Subscribe("testRunning", func(tc *TestContext) {
		tc.Counter++
	})

	// in
	events = events.Inherit()

	events.Subscribe("testNotRunning", func(tc *TestContext) {
		tc.Counter++
	})

	events.Dispatch(nil, "testRunning", testcontext)
	t.Logf("Counter %d", testcontext.Counter)

	if testcontext.Counter != 1 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}

	events.Subscribe("testRunning", func(tc *TestContext) {
		tc.Counter++
	})

	events.Dispatch(nil, "testRunning", testcontext)
	t.Logf("Counter %d", testcontext.Counter)
	if testcontext.Counter != 3 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}
}

func TestEmitOrderANDCancellation(t *testing.T) {
	events := NewDispatcher()

	testcontext := new(TestContext)

	events.Subscribe("cancelation", func(c *TestContext) {
		c.Counter = 1
		t.Fail()
	})

	events.Subscribe("cancelation", func(c *TestContext) {
		c.Counter = 2
		c.Cancel()
	})
	events.Dispatch(nil, "cancelation", testcontext)
	if testcontext.Counter != 2 {
		t.Fail()
	}
}

var bench_events = NewDispatcher()
var bench_context = new(TestContext)

//go:noinline
func bench_EventHandler(c *TestContext) {
	if c.Counter == -1 {
		c.Counter++
	}
}

var (
	_ = bench_events.Subscribe("benchmark", bench_EventHandler)
	_ = bench_events.Subscribe("benchmark", bench_EventHandler)
	_ = bench_events.Subscribe("benchmark", bench_EventHandler)
	_ = bench_events.Subscribe("benchmark", bench_EventHandler)
	_ = bench_events.Subscribe("benchmark", bench_EventHandler)
	_ = bench_events.Subscribe("benchmark", bench_EventHandler)
)

func BenchmarkManager_Emit(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bench_events.Dispatch(nil, "benchmark", bench_context)
		}
	})
}

func BenchmarkManager_Subscribe(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bench_events.Subscribe("benchmark", bench_EventHandler)
		}
	})
}

func BenchmarkManager_SubscribeEmit(b *testing.B) {
	b.SetParallelism(10000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bench_events.Reset("benchmark")
			bench_events.Subscribe("benchmark", func(k *TestContext) {
				k.Cancel()
			})
			bench_events.Subscribe("benchmark", bench_EventHandler)
			bench_events.Subscribe("benchmark", bench_EventHandler)
			bench_events.Subscribe("benchmark", bench_EventHandler)
			bench_events.Dispatch(nil, "benchmark", bench_context)
		}
	})
	b.Log("NumericId of handlers", len(bench_events.subscriptions[0].handlers))
}
