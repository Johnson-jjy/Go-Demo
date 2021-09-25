package main

import (
	"fmt"
	"sync"
	"time"
)

// 只使用一个互斥锁

func main()  {
	var (
		thinking = 0
		hungry = 1
		eating = 2
	)

	run := func(philoCnt int) {
		mutex := sync.RWMutex{}
		state := make([]int, philoCnt)
		for i := range state {
			state[i] = thinking
		}

		leftPos := func(i int) int {
			if i == 0 {
				return philoCnt - 1
			}
			return i - 1
		}
		rightPos := func(i int) int {
			if i == philoCnt - 1 {
				return 0
			}
			return i + 1
		}

		protectTest := func(i int) bool {
			if state[leftPos(i)] == eating || state[rightPos(i)] == eating {
				return false
			}
			state[i] = eating
			return  true
		}

		test := func(i int) bool {
			mutex.Lock()
			defer mutex.Unlock()
			if state[leftPos(i)] == eating || state[rightPos(i)] == eating {
				return false
			}
			state[i] = eating
			return true
		}

		getForks := func(i int) {
			mutex.Lock()
			state[i] = hungry
			mutex.Unlock()

			for test(i) == false {

			}
		}

		putForks := func(i int) {
			mutex.Lock()
			defer mutex.Unlock()
			state[i] = thinking
			protectTest(leftPos(i))
			protectTest(rightPos(i))
		}

		philosopher := func(i int) {
			for {
				time.Sleep(time.Second)
				getForks(i)
				time.Sleep(time.Second)
				putForks(i)
			}
		}
		for i := 0; i < philoCnt; i++ {
			go philosopher(i)
		}

		func() {
			for {
				mutex.RLock()
				eatingCnt := 0
				for _, val := range state {
					if val == eating {
						eatingCnt++
					}
				}
				fmt.Println(eatingCnt, state)
				mutex.RUnlock()
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	run(5)
}
