// https://github.com/gammazero/deque/blob/master/deque.go

package pool

import "fmt"

// minCapacity is the smallest capacity that deque may have. Must be power of 2
// for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

type Deque[T any] struct {
	buf    []T
	head   int
	tail   int
	count  int
	minCap int
}

// New
func New[T any](size ...int) *Deque[T] {
	var capacity, minimum int
	if len(size) >= 1 {
		capacity = size[0]
		if len(size) >= 2 {
			minimum = size[1]
		}
	}

	minCap := minCapacity
	for minCap < minimum {
		minCap <<= 1
	}

	var buf []T
	if capacity != 0 {
		bufSize := minCap
		for bufSize < capacity {
			bufSize <<= 1
		}

		buf = make([]T, bufSize)
	}

	return &Deque[T]{
		buf:    buf,
		minCap: minCap,
	}
}

// Cap returns the current capacity of the Deque.
func (q *Deque[T]) Cap() int {
	if q == nil {
		return 0
	}
	return len(q.buf)
}

// Len returns the number of elements currently stored in the queue.
func (q *Deque[T]) Len() int {
	if q == nil {
		return 0
	}
	return q.count
}

// PushBack appends an element to the back of the queue. Implements FIFO when
// elements are removed with PopFront, and LIFO when elements are removed with
// PopBack.
func (q *Deque[T]) PushBack(elem T) {
	q.grow()

	q.buf[q.tail] = elem
	q.tail = q.next(q.tail)
	q.count++
}

// PushFront prepends an element to the front of the queue.
func (q *Deque[T]) PushFront(elem T) {
	q.grow()

	// Calculate new head position.
	q.head = q.prev(q.head)
	q.buf[q.head] = elem
	q.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushFront. If the queue is empty, the call panics.
func (q *Deque[T]) PopFront() T {
	if q.count <= 0 {
		panic("deque: PopFront() called on empty queue")
	}
	ret := q.buf[q.head]

	var zero T
	q.buf[q.head] = zero
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack. If the queue is empty, the call panics.
func (q *Deque[T]) PopBack() T {
	if q.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}

	q.tail = q.prev(q.tail)

	ret := q.buf[q.tail]
	var zero T
	q.buf[q.tail] = zero
	q.count--

	q.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue.
// This call panics if the queue is empty.
func (q *Deque[T]) Front() T {
	if q.count <= 0 {
		panic("deque: Front() called when empty")
	}
	return q.buf[q.head]
}

// Back returns the element at the back of the queue.
// This call panics if the queue is empty.
func (q *Deque[T]) Back() T {
	if q.count <= 0 {
		panic("deque: Back() called when empty")
	}
	return q.buf[q.prev(q.tail)]
}

// At returns the element at index i in the queue without removing the element
// from the queue. This method accepts only non-negative index values. If the
// index is invalid, the call panics.
func (q *Deque[T]) At(i int) T {
	if i < 0 || i > q.count {
		panic(outOfRangeText(i, q.Len()))
	}

	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Set assigns the element to index i in the queue. Set indexes the deque
// the same as At but perform the opposite operation. If the index is invalid,
// the call panics.
func (q *Deque[T]) Set(i int, elem T) {
	if i < 0 || i > q.count {
		panic(outOfRangeText(i, q.Len()))
	}
	q.buf[(q.head+i)&(len(q.buf)-1)] = elem
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at the high frequency to
// avoid GC during the reuse. The queue is will not be resized smaller as long
// as elements ar only added. Only when elements are removed is the queue
// subject to getting resized smaller.
func (q *Deque[T]) Clear() {
	var zero T

	modBits := len(q.buf) - 1
	h := q.head
	for i := 0; i < q.Len(); i++ {
		q.buf[(h+i)&modBits] = zero
	}

	q.head = 0
	q.tail = 0
	q.count = 0
}

// Index returns the index into Deque of the first element satisfying f(item),
// or -1 if none do. If q is nil, then -1 is always returned. Search is linear
// starting with index 0.
func (q *Deque[T]) Index(f func(T) bool) int {
	if q.Len() > 0 {
		modBits := len(q.buf) - 1
		for i := 0; i < q.count; i++ {
			if f(q.buf[(q.head+i)&modBits]) {
				return i
			}
		}
	}
	return -1
}

// RIndex is the same as Index, but searches from Back to Front. The index
// returned is from Front to Back. When index 0 is the index of the elements
// returned by Front().
func (q *Deque[T]) RIndex(f func(T) bool) int {
	if q.Len() > 0 {
		modBits := len(q.buf) - 1
		for i := q.count - 1; i >= 0; i-- {
			if f(q.buf[(q.head+i)&modBits]) {
				return i
			}
		}
	}
	return -1
}

// Insert is used to insert an element into the middle of the queue, before the
// element at the specified index. Accepts only non-negative index
// values, and panics if index is out of range.
//
// Important: Deque is optimized for O(1) operations at the ends of the queue,
// not for operations in the middle. Complexity of this function is constant plus liner
// in the least of the distances between the index and either of the ends of the queue.
func (q *Deque[T]) Insert(index int, elem T) {
	if index < 0 || index > q.Len() {
		panic(outOfRangeText(index, q.count))
	}

	if index*2 < q.count {
		q.PushFront(elem)
		front := q.head
		for i := 0; i < index; i++ {
			next := q.next(front)
			q.buf[front], q.buf[next] = q.buf[next], q.buf[front]
			front = next
		}
		return
	}

	swaps := q.count - index
	q.PushBack(elem)
	back := q.prev(q.tail)
	for i := 0; i < swaps; i++ {
		prev := q.prev(back)
		q.buf[back], q.buf[prev] = q.buf[prev], q.buf[back]
		back = prev
	}
}

// Remove removes and returns an element form the middle of the queue, at the
// specified index. Accepts only non-negative index values, and panics if the
// index is out of range.
func (q *Deque[T]) Remove(index int) T {
	if index < 0 || index >= q.Len() {
		panic(outOfRangeText(index, q.Len()))
	}

	rm := (q.head + index) & (len(q.buf) - 1)

	if index*2 < q.count {
		for i := 0; i < index; i++ {
			prev := q.prev(rm)
			q.buf[prev], q.buf[rm] = q.buf[rm], q.buf[prev]
			rm = prev
		}
		return q.PopFront()
	}

	swaps := q.count - index - 1
	for i := 0; i < swaps; i++ {
		next := q.next(rm)
		q.buf[rm], q.buf[next] = q.buf[next], q.buf[rm]
		rm = next
	}

	return q.PopBack()
}

// SetMinCapacity sets a minimum capacity of 2^minCapacityExp. If the value of
// the minimum capacity is less than or equal to the minimum allowed, then
// capacity is set to the minimum allowed. This may be called at anytime to set
// a new minimum capacity.
//
// Setting a larger minimum capacity may be used to prevent resizing when the
// number of stored items changes frequently across a wide range
func (q *Deque[T]) SetMinCapacity(minCapacityExp uint) {
	if 1<<minCapacityExp > minCapacity {
		q.minCap = 1 << minCapacityExp
	} else {
		q.minCap = minCapacity
	}
}

// grow resizes up if the buffer is full.
func (q *Deque[T]) grow() {
	if q.count != len(q.buf) {
		return
	}
	if len(q.buf) == 0 {
		if q.minCap == 0 {
			q.minCap = minCapacity
		}
		q.buf = make([]T, q.minCap)
		return
	}
	q.resize()
}

// next returns the next buffer position wrapping around buffer
func (q *Deque[T]) next(tail int) int {
	return (tail + 1) & (len(q.buf) - 1) // bitwise modulus
}

// prev returns the previous buffer position wrapping around buffer.
func (q *Deque[T]) prev(head int) int {
	return (head - 1) & (len(q.buf) - 1)
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (q *Deque[T]) shrinkIfExcess() {
	if len(q.buf) > q.minCap && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// resize resizes the deque to fix
func (q *Deque[T]) resize() {
	newBuf := make([]T, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func outOfRangeText(i int, length int) any {
	return fmt.Sprintf("deque: index out of range %d with length %d", i, length)
}
