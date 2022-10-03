package pool

import (
	"fmt"
	"sync"
)

type node struct {
	values []interface{}
	next   *node
}

type Queue struct {
	head             *node
	tail             *node
	headPointer      uint
	headSlicePointer uint
	tailPointer      uint
	len              uint
}

func (q *Queue) push(data interface{}) {
	if q.head == nil {
		// Create head node when no nodes are in the queue.
		head := &node{values: make([]interface{}, 4)}
		head.next = head
		q.head = head
		q.tail = head
		q.tail.values[0] = data
		q.headSlicePointer = 3
		q.tailPointer = 1
	} else if q.tailPointer < uint(len(q.tail.values)) {
		q.tail.values[q.tailPointer] = data
		q.tailPointer++
	} else if q.tailPointer < 64 {
		newValues := make([]interface{}, len(q.tail.values)*4)
		copy(newValues, q.tail.values)
		q.tail.values = newValues
		q.tail.values[q.tailPointer] = data
		q.tailPointer++
		q.headSlicePointer = uint(len(newValues) - 1)
	} else if q.tail.next != q.head {
		next := q.tail.next
		q.tail = next
		q.tail.values[0] = data
		q.tailPointer = 1
	} else {
		// No available node is present in the queue.
		node := &node{values: make([]interface{}, 256)}
		node.next = q.head
		q.tail.next = node
		q.tail = node
		q.tail.values[0] = data
		q.tailPointer = 1
	}

	q.len++
}

func (q *Queue) pop() (interface{}, error) {
	if q.len == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	value := q.head.values[q.headPointer]

	if q.headPointer < q.headSlicePointer {
		q.headPointer++
	} else if q.head == q.tail {
		q.tailPointer = q.headPointer
	} else {
		q.headPointer = 0
		q.head = q.head.next
		q.headSlicePointer = uint(len(q.head.values) - 1)
	}

	q.len--
	return value, nil
}

func (q *Queue) empty() bool {
	return q.len <= 0
}

type ConcurrentQueue struct {
	lock     *sync.Mutex
	notEmpty *sync.Cond
	backend  *Queue
}

func NewConcurrentQueue(maxSize uint32) *ConcurrentQueue {
	queue := ConcurrentQueue{
		lock: &sync.Mutex{},
	}
	queue.notEmpty = sync.NewCond(queue.lock)

	queue.backend = &Queue{}
	return &queue
}

func (c *ConcurrentQueue) enqueue(data interface{}) {
	c.lock.Lock()

	c.backend.push(data)
	c.notEmpty.Signal()
	c.lock.Unlock()
}

func (c *ConcurrentQueue) dequeue() (interface{}, error) {
	c.lock.Lock()

	for c.backend.empty() {
		c.notEmpty.Wait()
	}

	data, err := c.backend.pop()
	c.lock.Unlock()
	return data, err
}
