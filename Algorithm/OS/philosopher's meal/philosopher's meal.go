package main

import (
	"fmt"
	"sync"
	"time"
)

// 哲学家进餐问题

func main()  {
	run := func(philoCnt int) { // 开始执行i个哲学家进餐
		var (
			thinking = 0
			hungry = 1
			eating = 2
		)
		mutex := sync.Mutex{}
		philoMutex := make([]sync.Mutex, philoCnt)
		state := make([]int, philoCnt)

		for i := range philoMutex {
			philoMutex[i].Lock() // 先全部上锁
		}

		leftPos := func(i int) int {
			return (philoCnt + i - 1) % philoCnt
		}
		rightPos := func(i int) int {
			return (philoCnt + i + 1) % philoCnt
		}

		test := func(i int) {
			if state[i] == hungry && state[leftPos(i)] != eating && state[rightPos(i)] != eating {
				state[i] = eating
				philoMutex[i].Unlock()
			}
		}

		getForks := func(i int) {
			mutex.Lock()
			state[i] = hungry
			test(i)
			mutex.Unlock()
			philoMutex[i].Lock()
		}

		putForks := func(i int) {
			mutex.Lock()
			state[i] = thinking
			test(leftPos(i))
			test(rightPos(i))
			mutex.Unlock()
		}

		philosopher := func(i int) {
			for {
				time.Sleep(time.Second) // think
				getForks(i)
				time.Sleep(time.Second) // eat
				putForks(i)
			}
		}

		// 部署i个哲学家
		for i := 0; i < philoCnt; i++ {
			go philosopher(i)
		}

		// 监视进餐,用来显示就餐情况
		func() {
			for {
				mutex.Lock()
				fmt.Println(state)
				mutex.Unlock()
				time.Sleep(100 * time.Millisecond)
			}
		} ()
	}

	run(5)
}

