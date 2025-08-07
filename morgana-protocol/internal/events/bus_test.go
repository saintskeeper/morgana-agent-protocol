package events

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestEventBusBasicFunctionality(t *testing.T) {
	config := DefaultBusConfig()
	config.BufferSize = 100
	config.Workers = 2
	bus := NewEventBus(config)
	defer bus.Close()

	// Test subscription
	var receivedEvents []Event
	var mu sync.Mutex

	subID := bus.Subscribe(EventTaskStarted, func(event Event) {
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
	})

	// Test publishing
	testEvent := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "test-task-1",
			Ctx:       context.Background(),
		},
		AgentType: "test-agent",
		Timeout:   time.Minute,
	}

	// Publish synchronously
	bus.Publish(testEvent)

	// Wait a bit for async processing
	time.Sleep(10 * time.Millisecond)

	// Check received events
	mu.Lock()
	if len(receivedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(receivedEvents))
	}

	if len(receivedEvents) > 0 && receivedEvents[0].Type() != EventTaskStarted {
		t.Errorf("Expected event type %s, got %s", EventTaskStarted, receivedEvents[0].Type())
	}
	mu.Unlock()

	// Test unsubscribe
	if !bus.Unsubscribe(subID) {
		t.Error("Failed to unsubscribe")
	}

	// Verify unsubscription worked
	mu.Lock()
	initialCount := len(receivedEvents)
	mu.Unlock()

	bus.Publish(testEvent)
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if len(receivedEvents) != initialCount {
		t.Errorf("Event received after unsubscription")
	}
	mu.Unlock()
}

func TestEventBusAsync(t *testing.T) {
	config := DefaultBusConfig()
	config.BufferSize = 1000
	bus := NewEventBus(config)
	defer bus.Close()

	var eventCount int64
	bus.Subscribe(EventTaskStarted, func(event Event) {
		atomic.AddInt64(&eventCount, 1)
	})

	// Publish many events asynchronously
	numEvents := 100
	for i := 0; i < numEvents; i++ {
		event := &TaskStartedEvent{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        GenerateTaskID(),
				Ctx:       context.Background(),
			},
			AgentType: "test-agent",
		}
		if !bus.PublishAsync(event) {
			t.Error("Failed to publish async event")
		}
	}

	// Wait for all events to be processed
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&eventCount) == int64(numEvents) {
			break
		}
		time.Sleep(time.Millisecond)
	}

	if atomic.LoadInt64(&eventCount) != int64(numEvents) {
		t.Errorf("Expected %d events, got %d", numEvents, atomic.LoadInt64(&eventCount))
	}
}

func TestEventBusFilter(t *testing.T) {
	config := DefaultBusConfig()
	bus := NewEventBus(config)
	defer bus.Close()

	var receivedTaskIDs []string
	var mu sync.Mutex

	// Subscribe with filter - only accept tasks with ID starting with "important"
	bus.SubscribeWithFilter(EventTaskStarted, SubscriberWithFilter{
		Handler: func(event Event) {
			mu.Lock()
			receivedTaskIDs = append(receivedTaskIDs, event.TaskID())
			mu.Unlock()
		},
		Filter: func(event Event) bool {
			return event.TaskID()[:9] == "important"
		},
	})

	// Publish events with different IDs
	events := []*TaskStartedEvent{
		{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        "important-task-1",
				Ctx:       context.Background(),
			},
		},
		{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        "regular-task-1",
				Ctx:       context.Background(),
			},
		},
		{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        "important-task-2",
				Ctx:       context.Background(),
			},
		},
	}

	for _, event := range events {
		bus.Publish(event)
	}

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(receivedTaskIDs) != 2 {
		t.Errorf("Expected 2 filtered events, got %d", len(receivedTaskIDs))
	}

	expectedIDs := []string{"important-task-1", "important-task-2"}
	for i, expectedID := range expectedIDs {
		if i >= len(receivedTaskIDs) || receivedTaskIDs[i] != expectedID {
			t.Errorf("Expected task ID %s, got %v", expectedID, receivedTaskIDs)
		}
	}
}

func TestEventBusSubscribeAll(t *testing.T) {
	config := DefaultBusConfig()
	bus := NewEventBus(config)
	defer bus.Close()

	var receivedEvents []Event
	var mu sync.Mutex

	// Subscribe to all events
	bus.SubscribeAll(func(event Event) {
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
	})

	// Publish different types of events
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
				ID:        "task-1",
				Ctx:       context.Background(),
			},
		},
	}

	for _, event := range events {
		bus.Publish(event)
	}

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(receivedEvents) != 2 {
		t.Errorf("Expected 2 events, got %d", len(receivedEvents))
	}

	if receivedEvents[0].Type() != EventTaskStarted || receivedEvents[1].Type() != EventTaskCompleted {
		t.Error("Events not received in correct order or type")
	}
}

func TestEventBusStats(t *testing.T) {
	config := DefaultBusConfig()
	config.BufferSize = 100
	bus := NewEventBus(config)
	defer bus.Close()

	// Subscribe to some events
	bus.Subscribe(EventTaskStarted, func(event Event) {})
	bus.Subscribe(EventTaskCompleted, func(event Event) {})
	bus.SubscribeAll(func(event Event) {})

	stats := bus.Stats()

	if stats.ActiveSubscribers != 3 {
		t.Errorf("Expected 3 active subscribers, got %d", stats.ActiveSubscribers)
	}

	if stats.QueueCapacity < 100 {
		t.Errorf("Expected queue capacity >= 100, got %d", stats.QueueCapacity)
	}

	if stats.SubscribersByType[EventTaskStarted] != 1 {
		t.Errorf("Expected 1 TaskStarted subscriber, got %d", stats.SubscribersByType[EventTaskStarted])
	}
}

func TestEventBusConcurrentAccess(t *testing.T) {
	config := DefaultBusConfig()
	bus := NewEventBus(config)
	defer bus.Close()

	var eventCount int64

	// Multiple subscribers
	for i := 0; i < 5; i++ {
		bus.Subscribe(EventTaskStarted, func(event Event) {
			atomic.AddInt64(&eventCount, 1)
		})
	}

	// Concurrent publishers
	var wg sync.WaitGroup
	numPublishers := 10
	eventsPerPublisher := 20

	for i := 0; i < numPublishers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < eventsPerPublisher; j++ {
				event := &TaskStartedEvent{
					BaseEvent: BaseEvent{
						EventType: EventTaskStarted,
						Time:      time.Now(),
						ID:        GenerateTaskID(),
						Ctx:       context.Background(),
					},
				}
				bus.PublishAsync(event)
			}
		}()
	}

	wg.Wait()

	// Wait for all events to be processed
	deadline := time.Now().Add(2 * time.Second)
	expectedEvents := int64(numPublishers * eventsPerPublisher * 5) // 5 subscribers per event

	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&eventCount) == expectedEvents {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if atomic.LoadInt64(&eventCount) != expectedEvents {
		t.Errorf("Expected %d events, got %d", expectedEvents, atomic.LoadInt64(&eventCount))
	}
}
