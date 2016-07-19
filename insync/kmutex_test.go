package insync_test

import (
	"github.com/CloudyKit/framework/insync"
	"sync"
	"testing"
	"time"
)

func TestSynchronizer_WorkON(t *testing.T) {
	var kMX = insync.KMutex{}

	var mx sync.RWMutex
	var counters = map[string]map[string]int{}
	var keys = []string{
		"key124123412312431",
		"key12341234123412412",
		"key12341234123421",
		"key123423143421",
		"key12341212",
		"key121431",
		"key123",
		"key1",
	}

	for _, key := range keys {
		counters[key] = make(map[string]int)
	}

	doWork := func(key string) {
		defer kMX.Lock(key).Unlock()

		mx.RLock()
		add := counters[key][key] + 1
		mx.RUnlock()

		mx.Lock()
		counters[key][key] = add
		mx.Unlock()
	}

	awaiter := &sync.WaitGroup{}
	const numOfIterations = 10000
	awaiter.Add(numOfIterations * len(keys))

	for i := 0; i < numOfIterations; i++ {
		for _, key := range keys {
			go func(key string) {
				doWork(key)
				awaiter.Done()
			}(key)
		}
	}
	awaiter.Wait()
	for _, key := range keys {
		counter := counters[key][key]
		if counter != numOfIterations {
			t.Errorf("Unexpected value %d on key %s", counter, key)
		}
	}
}

var BenchKeys = []string{
	"key124123412312431",
	"key12341234123412412",
	"key12341234123421",
	"key123423143421",
	"key12341212",
	"key121431",
	"key123",
	"key1",
}

var BenchKMutex = insync.KMutex{}
var BenchMutex = sync.Mutex{}

const durationOf = time.Microsecond

//go:noinline
func benchmarkKMutex(key string) {
	defer BenchKMutex.Lock(key).Unlock()
	time.Sleep(durationOf)
}

//go:noinline
func benchmarkMutex(key string) {
	BenchMutex.Lock()
	defer BenchMutex.Unlock()
	time.Sleep(durationOf)
}

func BenchmarkKeyMutex(b *testing.B) {
	awaiter := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1; i++ {
			for _, key := range BenchKeys {
				awaiter.Add(1)
				go func(key string) {
					benchmarkKMutex(key)
					awaiter.Done()
				}(key)
			}
		}
	}
	awaiter.Wait()
}

func BenchmarkMutex(b *testing.B) {
	awaiter := &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		for i := 0; i < 1; i++ {
			for _, key := range BenchKeys {
				awaiter.Add(1)
				go func(key string) {
					benchmarkMutex(key)
					awaiter.Done()
				}(key)
			}
		}
	}

	awaiter.Wait()
}
