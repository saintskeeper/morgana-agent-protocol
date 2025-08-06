package events

import (
	"sync/atomic"
	"unsafe"
)

// CircularBuffer implements a lock-free circular buffer for high-performance event storage
// Uses atomic operations for thread-safe access without locks
type CircularBuffer struct {
	buffer []unsafe.Pointer // Store Event interface{} as unsafe.Pointer
	size   int64            // Buffer size (power of 2 for efficient masking)
	mask   int64            // Bitmask for wrap-around (size - 1)
	head   int64            // Write position
	tail   int64            // Read position
}

// NewCircularBuffer creates a new circular buffer with the specified capacity
// Capacity will be rounded up to the next power of 2 for efficient indexing
func NewCircularBuffer(capacity int) *CircularBuffer {
	// Round up to next power of 2
	size := nextPowerOf2(capacity)

	return &CircularBuffer{
		buffer: make([]unsafe.Pointer, size),
		size:   int64(size),
		mask:   int64(size - 1),
		head:   0,
		tail:   0,
	}
}

// Push adds an event to the buffer. Returns false if buffer is full.
// This is a lock-free operation using atomic compare-and-swap
func (cb *CircularBuffer) Push(event Event) bool {
	for {
		head := atomic.LoadInt64(&cb.head)
		tail := atomic.LoadInt64(&cb.tail)

		// Check if buffer is full
		if head-tail >= cb.size {
			return false
		}

		// Try to reserve the slot
		if atomic.CompareAndSwapInt64(&cb.head, head, head+1) {
			// Successfully reserved slot, store the event
			index := head & cb.mask
			atomic.StorePointer(&cb.buffer[index], unsafe.Pointer(&event))
			return true
		}
		// CAS failed, retry
	}
}

// Pop removes and returns an event from the buffer. Returns nil if buffer is empty.
// This is a lock-free operation using atomic compare-and-swap
func (cb *CircularBuffer) Pop() Event {
	for {
		tail := atomic.LoadInt64(&cb.tail)
		head := atomic.LoadInt64(&cb.head)

		// Check if buffer is empty
		if tail >= head {
			return nil
		}

		// Try to acquire the slot
		if atomic.CompareAndSwapInt64(&cb.tail, tail, tail+1) {
			// Successfully acquired slot, load the event
			index := tail & cb.mask
			eventPtr := atomic.LoadPointer(&cb.buffer[index])
			if eventPtr == nil {
				continue // Slot not yet written, retry
			}

			// Clear the slot for reuse
			atomic.StorePointer(&cb.buffer[index], nil)

			// Convert back to Event interface
			event := *(*Event)(eventPtr)
			return event
		}
		// CAS failed, retry
	}
}

// PopBatch removes up to maxEvents from the buffer and returns them as a slice
// This is more efficient than calling Pop() multiple times
func (cb *CircularBuffer) PopBatch(maxEvents int, events []Event) int {
	count := 0

	for count < maxEvents && count < len(events) {
		event := cb.Pop()
		if event == nil {
			break
		}
		events[count] = event
		count++
	}

	return count
}

// Size returns the current number of events in the buffer
func (cb *CircularBuffer) Size() int {
	head := atomic.LoadInt64(&cb.head)
	tail := atomic.LoadInt64(&cb.tail)
	size := head - tail
	if size < 0 {
		return 0
	}
	return int(size)
}

// Capacity returns the maximum capacity of the buffer
func (cb *CircularBuffer) Capacity() int {
	return int(cb.size)
}

// IsFull returns true if the buffer is full
func (cb *CircularBuffer) IsFull() bool {
	head := atomic.LoadInt64(&cb.head)
	tail := atomic.LoadInt64(&cb.tail)
	return head-tail >= cb.size
}

// IsEmpty returns true if the buffer is empty
func (cb *CircularBuffer) IsEmpty() bool {
	head := atomic.LoadInt64(&cb.head)
	tail := atomic.LoadInt64(&cb.tail)
	return head == tail
}

// nextPowerOf2 returns the next power of 2 greater than or equal to n
func nextPowerOf2(n int) int {
	if n <= 1 {
		return 1
	}

	// Handle specific case where n is already a power of 2
	if n&(n-1) == 0 {
		return n
	}

	// Find the next power of 2
	power := 1
	for power < n {
		power <<= 1
	}
	return power
}
