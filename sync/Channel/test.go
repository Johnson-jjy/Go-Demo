package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			fmt.Println("A:", v)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c) // 不关闭则造成死锁
	}()
	return c
}

func main()  {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	in := asChan(nums...)
	for v := range in {
		fmt.Println("B:", v)
	}
}
