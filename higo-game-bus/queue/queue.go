package queue

import "sync"

type Queue struct {
	lock sync.Mutex
	cond *sync.Cond
	data []interface{}
}

func NewQueue(size int) *Queue {
	q := &Queue{
		data: make([]interface{}, 0, size),
	}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.data = append(q.data, item)
	q.cond.Signal()
}

func (q *Queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.data) == 0 {
		q.cond.Wait()
	}

	item := q.data[0]
	q.data = q.data[1:]

	return item
}
