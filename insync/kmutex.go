package insync

import (
	"sync"
	"sync/atomic"
)

//KMutex is a synchronization mechanism that's allows the creation of lock by key,
type KMutex struct {
	locker sync.Mutex
	keys   map[string]*LockedMutex
}

//LockedMutex represents a locked mutex as the name says, call Unlock to release the
//lock on the key
type LockedMutex struct {
	key       string
	awaiting  int64
	keyLocker *KMutex

	locker sync.Mutex
}

//Lock locks a specific key, and returns an locked mutex
//you need to call Unlock to release the lock
func (kmutex *KMutex) Lock(key string) *LockedMutex {
	//lock the main lock
	kmutex.locker.Lock()
	//try to get a stored lock for the key
	mutex, ok := kmutex.keys[key]
	if ok {
		//if a stored locker is found increase the number of awaiting goroutines
		mutex.awaiting++
	} else {
		//if there is no locker with that key, we create and store a new one
		mutex = workPool.Get().(*LockedMutex)
		mutex.key = key
		mutex.keyLocker = kmutex

		if kmutex.keys == nil {
			kmutex.keys = make(map[string]*LockedMutex)
		}
		kmutex.keys[key] = mutex
	}

	//release the main locker
	kmutex.locker.Unlock()
	//lock the the key locker
	mutex.locker.Lock()
	return mutex
}

//Unlocks the Locker
func (lockedMutex *LockedMutex) Unlock() {
	if atomic.LoadInt64(&lockedMutex.awaiting) == 0 {
		lockedMutex.keyLocker.locker.Lock()
		if lockedMutex.awaiting == 0 {
			delete(lockedMutex.keyLocker.keys, lockedMutex.key)
			workPool.Put(lockedMutex)
		}
		lockedMutex.keyLocker.locker.Unlock()
	} else {
		atomic.AddInt64(&lockedMutex.awaiting, -1)
	}
	lockedMutex.locker.Unlock()
}

var workPool = sync.Pool{
	New: func() interface{} {
		return &LockedMutex{}
	},
}
