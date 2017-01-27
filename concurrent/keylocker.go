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

package concurrent

import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

// NewKeyLocker returns a new KeyLocker
func NewKeyLocker() *KeyLocker {
	return &KeyLocker{keys: make(map[string]*lockedKey)}
}

// KeyLocker is a synchronization mechanism that's allows the creation of lock by key,
type KeyLocker struct {
	locker sync.Mutex
	keys   map[string]*lockedKey
}

// lockedKey represents a locked mutex as the name says, call Unlock to release the
// lock on the key
type lockedKey struct {
	key       string
	awaiting  int64
	keyLocker *KeyLocker
	locker    sync.Mutex
}

// Unlock unlocks key locker
func (locker *lockedKey) Unlock() {
	if atomic.LoadInt64(&locker.awaiting) == 0 {
		locker.keyLocker.locker.Lock()
		if locker.awaiting == 0 {
			delete(locker.keyLocker.keys, locker.key)
		}
		locker.keyLocker.locker.Unlock()
	} else {
		atomic.AddInt64(&locker.awaiting, -1)
	}
	locker.locker.Unlock()
}

// Lock panic always, used to implement the sync.Locker
func (locker *lockedKey) Lock() {
	panic(errors.New("Lock is unavailable in key locker, to lock the key use the key locker value"))
}

// Lock locks a specific key, and returns an locked mutex
// you need to call Unlock to release the lock
func (kmutex *KeyLocker) Lock(key string) sync.Locker {
	//lock the main lock
	kmutex.locker.Lock()
	//try to get a stored lock for the key
	mutex, ok := kmutex.keys[key]
	if ok {
		//if a stored locker is found increase the number of awaiting goroutines
		mutex.awaiting++
	} else {
		//if there is no locker with that key, we create and store a new one
		mutex = new(lockedKey)
		mutex.key = key
		mutex.keyLocker = kmutex
		kmutex.keys[key] = mutex
	}
	//release the main locker
	kmutex.locker.Unlock()
	//lock the the key locker
	mutex.locker.Lock()
	return mutex
}
