package pool

import (
	"fmt"
	"sync"
)

const (
	// sliceReplicator is how much the slice size should replicate.
	sliceReplicator = 4
	// maxFirstSlice is the size of the head slice.
	maxFirstSlice = 64
	// sliceSize is the slice size of each internal used node.
	sliceSize = 256
)

// node represents an internally used queue node.
type node[T any] struct {
	// values represents the stored values inside the node.
	values []T
	// next points to the next node in the queue.
	next *node[T]
}

// Queue is a modified queue implementation duplicate of the queue design principles
// described in https://github.com/ef-ds/queue.
type Queue[T any] struct {
	// head points to the node that is the first node in the queue.
	head *node[T]
	// tail points to the node that is the last node in the queue.
	tail *node[T]
	// headPointer is the index to the node that is the first node in the queue.
	headPointer      uint
	headSlicePointer uint
	// tailPointer is the index to the node that is the last node in the queue.
	tailPointer uint
	// len is the length of the queue.
	len uint
}

// push adds any data to the end of the queue.
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

// pop removes and returns the element at the front of the queue.
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

// empty returns `true` if the queue is empty.
func (q *Queue[T]) empty() bool {
	return q.len <= 0
}

// ConcurrentQueue is a concurrent implementation of the designed `Queue`.
type ConcurrentQueue[T any] struct {
	// lock is a mutex which locks inserting or removing nodes from the queue.
	lock     *sync.Mutex
	notEmpty *sync.Cond
	// backend is the used queue.
	backend *Queue[T]
}

// NewConcurrentQueue creates a new queue with any value type and areturns it.
func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	queue := ConcurrentQueue[T]{
		lock: &sync.Mutex{},
	}
	queue.notEmpty = sync.NewCond(queue.lock)

	queue.backend = &Queue[T]{}
	return &queue
}

// enqueue pushes data to the queue and locks the mutex while pushing to the queue.
func (c *ConcurrentQueue[T]) enqueue(data T) {
	c.lock.Lock()

	c.backend.push(data)
	// Notify waiting goroutines that the queue is not empty.
	c.notEmpty.Signal()
	c.lock.Unlock()
}

// dequeue pops data from the queue and locks the mutex while popping from the queue.
func (c *ConcurrentQueue[T]) dequeue() (interface{}, error) {
	c.lock.Lock()

	// Wait for the queue to become non-empty.
	for c.backend.empty() {
		c.notEmpty.Wait()
	}

	data, err := c.backend.pop()
	c.lock.Unlock()
	return data, err
}
