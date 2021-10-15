package Test

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

func (sl *spinLock) Lock() {
	var iter uint32
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		if iter<5{iter++;continue}
		runtime.Gosched()
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

func NewSpinLock() sync.Locker {
	return new(spinLock)
}
