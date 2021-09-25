package main

import (
	"fmt"
	"sync"
	"time"
)

var condition int

// 消费者
func consumer(name string, cond *sync.Cond) {
	for {
		// 消费者开始消费时, 锁住
		cond.L.Lock()
		// 如果没有可消费的值, 则等待
		for condition == 0 {
			cond.Wait()
		}
		// 消费
		condition--
		fmt.Printf("%s: %d\n", name, condition)

		// 唤醒一个消费者
		cond.Signal()
		// 解锁
		cond.L.Unlock()
	}
}


// 生产者
func producer(name string, cond *sync.Cond)  {
	for {
		// 生产者开始生产
		cond.L.Lock()

		// 当生产太多时, 等待消费者消费
		for condition == 100 {
			cond.Wait()
		}

		// 生产
		condition++
		fmt.Printf("%s: %d\n", name, condition)

		// 通知消费者可以开始消费了
		cond.Signal()
		// 解锁
		cond.L.Unlock()
	}

}
func main()  {
	cond := sync.NewCond(new(sync.Mutex))

	go consumer("consumer1", cond)
	go consumer("consumer2", cond)

	go producer("producer1", cond)
	go producer("producer2", cond)
	go producer("producer3", cond)

	time.Sleep(time.Millisecond)
}
