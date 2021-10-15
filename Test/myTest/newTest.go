package main

import (
	"fmt"
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

var wg sync.WaitGroup

func main()  {
	sp := NewSpinLock()
	cnt := 0
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				sp.Lock()
				cnt++
				sp.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(cnt)
}