package pool

import (
	"fmt"
	"sync"
)

type Node struct {
	value interface{}
}

type Queue struct {
	nodes   []*Node
	size    uint32
	maxSize uint32
}

func (q Queue) createNode(data interface{}) *Node {
	node := Node{}
	node.value = data
	return &node
}

func (q *Queue) put(data interface{}) error {
	if q.size >= q.maxSize {
		return fmt.Errorf("queue is full")
	}

	node := q.createNode(data)
	q.nodes = append(q.nodes[:q.size], node)
	q.size++

	return nil
}

func (q *Queue) pop() (interface{}, error) {
	if q.empty() {
		return nil, fmt.Errorf("queue is empty")
	}

	q.size--
	node := q.nodes[0]

	q.nodes = q.nodes[1:]

	return node.value, nil
}

func (q *Queue) empty() bool {
	return q.size <= 0
}

func (q *Queue) full() bool {
	return q.size >= q.maxSize
}

type ConcurrentQueue struct {
	lock     *sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	backend  *Queue
}

func NewConcurrentQueue(maxSize uint32) *ConcurrentQueue {
	queue := ConcurrentQueue{
		lock: &sync.Mutex{},
	}
	queue.notFull = sync.NewCond(queue.lock)
	queue.notEmpty = sync.NewCond(queue.lock)

	queue.backend = &Queue{
		size:    0,
		maxSize: maxSize,
	}
	return &queue
}

func (c *ConcurrentQueue) enqueue(data interface{}) error {
	c.lock.Lock()

	for c.backend.full() {
		c.notFull.Wait()
	}

	err := c.backend.put(data)
	c.notEmpty.Signal()
	c.lock.Unlock()
	return err
}

func (c *ConcurrentQueue) dequeue() (interface{}, error) {
	c.lock.Lock()

	for c.backend.empty() {
		c.notEmpty.Wait()
	}

	data, err := c.backend.pop()
	c.notFull.Signal()
	c.lock.Unlock()
	return data, err
}

func (c *ConcurrentQueue) getSize() uint32 {
	c.lock.Lock()
	size := c.backend.size
	c.lock.Unlock()
	return size
}

func (c *ConcurrentQueue) getMaxSize() uint32 {
	c.lock.Lock()
	maxSize := c.backend.maxSize
	c.lock.Unlock()
	return maxSize
}
