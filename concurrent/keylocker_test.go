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

package concurrent_test

import (
	"github.com/CloudyKit/framework/concurrent"
	"sync"
	"testing"
	"time"
)

func TestKeyLocker(t *testing.T) {
	var kMX = concurrent.NewKeyLocker()

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

var BenchKMutex = concurrent.NewKeyLocker()
var BenchMutex = sync.Mutex{}

const durationOf = time.Microsecond

//go:noinline
func benchmarkKMutex(key string) {
	lock := BenchKMutex.Lock(key)
	time.Sleep(durationOf)
	lock.Unlock()
}

//go:noinline
func benchmarkMutex(key string) {
	BenchMutex.Lock()
	time.Sleep(durationOf)
	BenchMutex.Unlock()
}

func BenchmarkKeyMutex(b *testing.B) {
	wgroup := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1; i++ {
			for _, key := range BenchKeys {
				wgroup.Add(1)
				go func(key string) {
					benchmarkKMutex(key)
					wgroup.Done()
				}(key)
			}
		}
	}
	wgroup.Wait()
}

func BenchmarkMutex(b *testing.B) {
	wgroup := &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		for i := 0; i < 1; i++ {
			for _, key := range BenchKeys {
				wgroup.Add(1)
				go func(key string) {
					benchmarkMutex(key)
					wgroup.Done()
				}(key)
			}
		}
	}

	wgroup.Wait()
}
