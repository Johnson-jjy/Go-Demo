package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func xadd(val *int32, delta int32) (new int32) {
	for {
		v := *val;
		if atomic.CompareAndSwapInt32(val, v, v+delta) {
			return v+delta;
		}
	}
	panic("unreached")
}

func main()  {
	var test int32
	test = 0
	count := 0
	go func() {
		for  {
			fmt.Println("A")
			if xadd(&test, 1) == 1 {
				fmt.Println(count)
				count++
				return
			}
		}
	}()

	go func() {
		for  {
			fmt.Println("B")
			if xadd(&test, 1) == 1 {
				fmt.Println(count)
				count++
				return
			}
		}
	}()

	go func() {
		for  {
			fmt.Println("C")
			if xadd(&test, -1) == 0 {
				fmt.Println(count)
				count--
				return
			}
		}
	}()

	time.Sleep(3 * time.Second)
}