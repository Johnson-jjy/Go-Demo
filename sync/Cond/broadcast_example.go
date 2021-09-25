// 三个协程调用 Wait() 等待, 另一个协程调用 Broadcast() 唤醒所有等待的协程
package main

import (
	"log"
	"sync"
	"time"
)

// done 即互斥锁需要保护的条件变量
var done = false

// read() 调用 Wait() 等待通知，直到 done 为 true
func read(name string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		// 当前协程被唤醒时, 条件不一定符合要求, 需要再次 Wait 等待下次被唤醒.
		// 为了保险起见, 使用 for 能够确保条件符合
		c.Wait()
	}
	log.Println(name, "starts reading")
	c.L.Unlock()
}

// write() 接收数据，接收完成后，将 done 置为 true，调用 Broadcast() 通知所有等待的协程
func write(name string, c *sync.Cond) {
	log.Println(name, "starts writing")
	time.Sleep(time.Second) // 确保前面的 3 个 read 协程都执行到 Wait(),处于等待状态
	c.L.Lock()
	done = true
	c.L.Unlock()
	log.Println(name, "wakes all")
	c.Broadcast()
}

func main()  {
	cond := sync.NewCond(&sync.Mutex{})

	go read("reader1", cond)
	go read("reader2", cond)
	go read("reader3", cond)

	write("writer", cond)

	time.Sleep(time.Second * 3)
}
