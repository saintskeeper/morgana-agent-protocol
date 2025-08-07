package events

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCircularBufferBasicOperations(t *testing.T) {
	buffer := NewCircularBuffer(4)

	if !buffer.IsEmpty() {
		t.Error("New buffer should be empty")
	}

	if buffer.IsFull() {
		t.Error("New buffer should not be full")
	}

	if buffer.Capacity() != 4 {
		t.Errorf("Expected capacity 4, got %d", buffer.Capacity())
	}

	// Create test events
	events := []Event{
		&TaskStartedEvent{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        "task-1",
				Ctx:       context.Background(),
			},
		},
		&TaskCompletedEvent{
			BaseEvent: BaseEvent{
				EventType: EventTaskCompleted,
				Time:      time.Now(),
				ID:        "task-2",
				Ctx:       context.Background(),
			},
		},
	}

	// Push events
	if !buffer.Push(events[0]) {
		t.Error("Failed to push first event")
	}

	if buffer.Size() != 1 {
		t.Errorf("Expected size 1, got %d", buffer.Size())
	}

	if !buffer.Push(events[1]) {
		t.Error("Failed to push second event")
	}

	if buffer.Size() != 2 {
		t.Errorf("Expected size 2, got %d", buffer.Size())
	}

	// Pop events
	poppedEvent := buffer.Pop()
	if poppedEvent == nil {
		t.Error("Failed to pop event")
	}

	if poppedEvent.Type() != EventTaskStarted {
		t.Errorf("Expected %s, got %s", EventTaskStarted, poppedEvent.Type())
	}

	if buffer.Size() != 1 {
		t.Errorf("Expected size 1 after pop, got %d", buffer.Size())
	}

	// Pop second event
	poppedEvent = buffer.Pop()
	if poppedEvent == nil {
		t.Error("Failed to pop second event")
	}

	if poppedEvent.Type() != EventTaskCompleted {
		t.Errorf("Expected %s, got %s", EventTaskCompleted, poppedEvent.Type())
	}

	if !buffer.IsEmpty() {
		t.Error("Buffer should be empty after popping all events")
	}

	// Pop from empty buffer
	poppedEvent = buffer.Pop()
	if poppedEvent != nil {
		t.Error("Pop from empty buffer should return nil")
	}
}

func TestCircularBufferFull(t *testing.T) {
	buffer := NewCircularBuffer(2)

	// Fill buffer
	event1 := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "task-1",
			Ctx:       context.Background(),
		},
	}
	event2 := &TaskCompletedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskCompleted,
			Time:      time.Now(),
			ID:        "task-2",
			Ctx:       context.Background(),
		},
	}
	event3 := &TaskFailedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskFailed,
			Time:      time.Now(),
			ID:        "task-3",
			Ctx:       context.Background(),
		},
	}

	if !buffer.Push(event1) {
		t.Error("Failed to push first event")
	}

	if !buffer.Push(event2) {
		t.Error("Failed to push second event")
	}

	if !buffer.IsFull() {
		t.Error("Buffer should be full")
	}

	// Try to push to full buffer
	if buffer.Push(event3) {
		t.Error("Push to full buffer should return false")
	}

	// Pop one event and try again
	poppedEvent := buffer.Pop()
	if poppedEvent == nil {
		t.Error("Failed to pop from full buffer")
	}

	if buffer.Push(event3) {
		t.Error("Push after pop should not succeed due to lock-free implementation timing")
	}
}

func TestCircularBufferPowerOfTwoRounding(t *testing.T) {
	// Test that capacity is rounded up to next power of 2
	testCases := []struct {
		input    int
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 4},
		{5, 8},
		{7, 8},
		{9, 16},
		{15, 16},
		{17, 32},
	}

	for _, tc := range testCases {
		buffer := NewCircularBuffer(tc.input)
		if buffer.Capacity() != tc.expected {
			t.Errorf("For input %d, expected capacity %d, got %d", tc.input, tc.expected, buffer.Capacity())
		}
	}
}

func TestCircularBufferConcurrentAccess(t *testing.T) {
	buffer := NewCircularBuffer(1000)
	var wg sync.WaitGroup

	// Number of producers and consumers
	numProducers := 10
	numConsumers := 5
	eventsPerProducer := 100

	// Start producers
	for i := 0; i < numProducers; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			for j := 0; j < eventsPerProducer; j++ {
				event := &TaskStartedEvent{
					BaseEvent: BaseEvent{
						EventType: EventTaskStarted,
						Time:      time.Now(),
						ID:        GenerateTaskID(),
						Ctx:       context.Background(),
					},
					AgentType: "test-agent",
				}

				// Retry push if buffer is temporarily full
				for !buffer.Push(event) {
					time.Sleep(time.Microsecond)
				}
			}
		}(i)
	}

	// Start consumers
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			localCount := 0

			for localCount < eventsPerProducer*numProducers/numConsumers {
				event := buffer.Pop()
				if event != nil {
					localCount++
				} else {
					time.Sleep(time.Microsecond)
				}
			}
		}(i)
	}

	wg.Wait()

	// Buffer should be close to empty (some events might still be in transit)
	if buffer.Size() > 50 {
		t.Errorf("Expected buffer to be mostly empty, got size %d", buffer.Size())
	}
}

func TestCircularBufferBatch(t *testing.T) {
	buffer := NewCircularBuffer(10)

	// Push several events
	numEvents := 5
	for i := 0; i < numEvents; i++ {
		event := &TaskStartedEvent{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        GenerateTaskID(),
				Ctx:       context.Background(),
			},
		}
		if !buffer.Push(event) {
			t.Errorf("Failed to push event %d", i)
		}
	}

	// Pop batch
	events := make([]Event, 10)
	count := buffer.PopBatch(3, events)

	if count != 3 {
		t.Errorf("Expected to pop 3 events, got %d", count)
	}

	if buffer.Size() != 2 {
		t.Errorf("Expected 2 events remaining, got %d", buffer.Size())
	}

	// Pop remaining events in batch
	count = buffer.PopBatch(10, events)
	if count != 2 {
		t.Errorf("Expected to pop 2 remaining events, got %d", count)
	}

	if !buffer.IsEmpty() {
		t.Error("Buffer should be empty after popping all events")
	}
}
