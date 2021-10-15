package cqueue

import "sync"

type ConcurrentQueue interface {
	Dequeue() interface{}
	Enqueue(v interface{})
}

type SliceQueue struct {
	data []interface{}
	mu sync.Mutex
}

// 暂未考虑err等, 未来补充; 同理,方法中的很多设计都还很不严谨
func NewConcurrentQueue(n int) ConcurrentQueue {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

// Enqueue 把值放在队尾
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue 移去队头并返回
func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}

