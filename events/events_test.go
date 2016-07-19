package events

import "testing"

type TestContext struct {
	Counter int
}

func TestManager_Emit(t *testing.T) {

	events := new(Emitter)
	testcontext := new(TestContext)

	events.Subscribe("testRunning", func(e *Event, tc *TestContext) {
		tc.Counter++
	})

	events.Subscribe("testNotRunning", func(e *Event, tc *TestContext) {
		tc.Counter++
	})

	events.Emit("testRunning", "", testcontext)
	t.Logf("Counter %d", testcontext.Counter)

	if testcontext.Counter != 1 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}

	events.Subscribe("testRunning", func(e *Event, tc *TestContext) {
		tc.Counter++
	})

	events.Emit("testRunning", "", testcontext)
	t.Logf("Counter %d", testcontext.Counter)
	if testcontext.Counter != 3 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}
}

func TestManager_EmitParent(t *testing.T) {

	testcontext := new(TestContext)

	events := new(Emitter)
	events.Subscribe("testRunning", func(e *Event, tc *TestContext) {
		tc.Counter++
	})

	// in
	events = events.Inherit()

	events.Subscribe("testNotRunning", func(e *Event, tc *TestContext) {
		tc.Counter++
	})

	events.Emit("testRunning", "", testcontext)
	t.Logf("Counter %d", testcontext.Counter)

	if testcontext.Counter != 1 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}

	events.Subscribe("testRunning", func(e *Event, tc *TestContext) {
		tc.Counter++
	})

	events.Emit("testRunning", "", testcontext)
	t.Logf("Counter %d", testcontext.Counter)
	if testcontext.Counter != 3 {
		t.Fatalf("Subscribe func for testRunning was not called %#v", events)
	}
}

func TestEmitOrderANDCancellation(t *testing.T) {
	events := new(Emitter)

	testcontext := new(TestContext)

	events.Subscribe("cancelation", func(v *Event, c *TestContext) {
		c.Counter = 1
		t.Fail()
	})

	events.Subscribe("cancelation", func(v *Event, c *TestContext) {
		c.Counter = 2
		v.Cancel()
	})
	events.Emit("cancelation", "", testcontext)
	if testcontext.Counter != 2 {
		t.Fail()
	}
}

var bench_events = new(Emitter)
var bench_context = new(TestContext)

//go:noinline
func bench_EventHandler(v *Event, c *TestContext) {
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
			bench_events.Emit("benchmark", "Emit", bench_context)
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
			bench_events.Subscribe("benchmark", func(e *Event, k *TestContext) {
				e.Cancel()
			})
			bench_events.Subscribe("benchmark", bench_EventHandler)
			bench_events.Subscribe("benchmark", bench_EventHandler)
			bench_events.Subscribe("benchmark", bench_EventHandler)
			bench_events.Emit("benchmark", "Emit", bench_context)
		}
	})
	b.Log("Number of handlers", len(bench_events.subscriptions[0].handlers))
}
