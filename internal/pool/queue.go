package pool

import (
	"fmt"
	"sync"
)

const (
	// Describes how much the slice size should replicate.
	sliceReplicator = 4

	// Describes the size of the head slice.
	maxFirstSlice = 64

	// Describes the slice size of each internal used node.
	sliceSize = 256
)

type node[T any] struct {
	values []T
	next   *node[T]
}

type Queue[T any] struct {
	head             *node[T]
	tail             *node[T]
	headPointer      uint
	headSlicePointer uint
	tailPointer      uint
	len              uint
}

func (q *Queue[T]) push(data T) {
	q.len++

	if q.head == nil {
		// Initialize head node when no nodes are in the queue.
		head := &node[T]{values: make([]T, sliceReplicator)}
		head.next = head
		q.head = head
		q.tail = head
		q.tail.values[0] = data
		q.headSlicePointer = sliceReplicator - 1
		q.tailPointer = 1
		return
	}

	if q.tailPointer < uint(len(q.tail.values)) {
		// Append data to tail slice when it is not full.
		q.tail.values[q.tailPointer] = data
		q.tailPointer++
		return
	}

	if q.tailPointer < maxFirstSlice {
		// Still wants to add data to the first slice, but the tail is not large enough yet.
		newValues := make([]T, len(q.tail.values)*sliceReplicator)
		copy(newValues, q.tail.values)
		q.tail.values = newValues
		q.tail.values[q.tailPointer] = data
		q.tailPointer++
		q.headSlicePointer = uint(len(newValues) - 1)
		return
	}

	if q.tail.next != q.head {
		// The `next` node in the tail is not the head and there is still room within the tail.
		next := q.tail.next
		q.tail = next
		q.tail.values[0] = data
		q.tailPointer = 1
		return
	}

	// No available node is present in the queue.
	node := &node[T]{values: make([]T, sliceSize)}
	node.next = q.head
	q.tail.next = node
	q.tail = node
	q.tail.values[0] = data
	q.tailPointer = 1
}

func (q *Queue[T]) pop() (T, error) {
	if q.len == 0 {
		return *new(T), fmt.Errorf("queue is empty")
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

func (q *Queue[T]) empty() bool {
	return q.len <= 0
}

type ConcurrentQueue[T any] struct {
	lock     *sync.Mutex
	notEmpty *sync.Cond
	backend  *Queue[T]
}

func NewConcurrentQueue[T any](maxSize uint32) *ConcurrentQueue[T] {
	queue := ConcurrentQueue[T]{
		lock: &sync.Mutex{},
	}
	queue.notEmpty = sync.NewCond(queue.lock)

	queue.backend = &Queue[T]{}
	return &queue
}

func (c *ConcurrentQueue[T]) enqueue(data T) {
	c.lock.Lock()

	c.backend.push(data)
	c.notEmpty.Signal()
	c.lock.Unlock()
}

func (c *ConcurrentQueue[T]) dequeue() (interface{}, error) {
	c.lock.Lock()

	for c.backend.empty() {
		c.notEmpty.Wait()
	}

	data, err := c.backend.pop()
	c.lock.Unlock()
	return data, err
}
