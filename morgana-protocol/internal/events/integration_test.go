package events

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// MockTaskClient implements a simple task client for testing
type MockTaskClient struct {
	delay time.Duration
}

func (m *MockTaskClient) RunWithContext(ctx context.Context, taskType string, prompt string, options map[string]interface{}) (*TaskResult, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	return &TaskResult{
		Output: "Mock task output",
	}, nil
}

// MockPromptLoader implements a simple prompt loader for testing
type MockPromptLoader struct{}

func (m *MockPromptLoader) Load(agentType string) (string, error) {
	return "Mock system prompt for " + agentType, nil
}

// MockTracer implements a simple tracer for testing
type MockTracer struct{}

func (m *MockTracer) Start(ctx context.Context, name string, opts ...interface{}) (context.Context, MockSpan) {
	return ctx, MockSpan{}
}

type MockSpan struct{}

func (m MockSpan) End()                          {}
func (m MockSpan) RecordError(error)             {}
func (m MockSpan) SetStatus(interface{}, string) {}
func (m MockSpan) SetAttributes(...interface{})  {}

// Define TaskResult since it's not defined in the adapter package
type TaskResult struct {
	Output string
}

// Update the MockTaskClient to return the correct type
func (m *MockTaskClient) RunWithContext2(ctx context.Context, taskType string, prompt string, options map[string]interface{}) (TaskResult, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	return TaskResult{
		Output: "Mock task output",
	}, nil
}

func TestEventIntegration(t *testing.T) {
	// Create event bus
	config := DefaultBusConfig()
	eventBus := NewEventBus(config)
	defer eventBus.Close()

	// Track events
	var events []Event
	var mu sync.Mutex

	// Subscribe to all events
	eventBus.SubscribeAll(func(event Event) {
		mu.Lock()
		events = append(events, event)
		mu.Unlock()
	})

	// Create and publish test events
	ctx := context.Background()
	taskID := GenerateTaskID()

	startEvent := NewTaskStartedEvent(ctx, taskID, "code-implementer", "Test prompt", nil, 0, "", "", time.Minute)
	progressEvent := NewTaskProgressEvent(ctx, taskID, "code-implementer", "execution", "Running task", 0.5, time.Second)
	completeEvent := NewTaskCompletedEvent(ctx, taskID, "code-implementer", "Task completed", time.Second*2, "test-model")

	eventBus.Publish(startEvent)
	eventBus.Publish(progressEvent)
	eventBus.Publish(completeEvent)

	// Wait for events to be processed
	time.Sleep(50 * time.Millisecond)

	// Analyze events
	mu.Lock()
	defer mu.Unlock()

	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	expectedTypes := []EventType{EventTaskStarted, EventTaskProgress, EventTaskCompleted}
	for i, expectedType := range expectedTypes {
		if i >= len(events) || events[i].Type() != expectedType {
			t.Errorf("Event %d: expected %s, got %s", i, expectedType, events[i].Type())
		}
	}

	// Check that all events have the same task ID
	for _, event := range events {
		if event.TaskID() != taskID {
			t.Errorf("Event has wrong task ID: expected %s, got %s", taskID, event.TaskID())
		}
	}
}

func TestEventBusPerformanceImpact(t *testing.T) {
	// Test performance impact of event system
	config := DefaultBusConfig()
	eventBus := NewEventBus(config)
	defer eventBus.Close()

	// Add several subscribers to create realistic load
	var eventCount int64
	for i := 0; i < 5; i++ {
		eventBus.Subscribe(EventTaskStarted, func(event Event) {
			atomic.AddInt64(&eventCount, 1)
			// Simulate minimal processing
			_ = event.TaskID()
			_ = event.Timestamp()
		})
	}

	// Measure time to publish many events
	numEvents := 10000
	start := time.Now()

	for i := 0; i < numEvents; i++ {
		event := &TaskStartedEvent{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        GenerateTaskID(),
				Ctx:       context.Background(),
			},
			AgentType: "performance-test",
		}

		if !eventBus.PublishAsync(event) {
			t.Error("Failed to publish event")
		}
	}

	publishDuration := time.Since(start)

	// Wait for all events to be processed
	deadline := time.Now().Add(5 * time.Second)
	expectedEvents := int64(numEvents * 5) // 5 subscribers

	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&eventCount) >= expectedEvents {
			break
		}
		time.Sleep(time.Millisecond)
	}

	totalDuration := time.Since(start)

	// Performance assertions
	publishThroughput := float64(numEvents) / publishDuration.Seconds()
	if publishThroughput < 100000 { // 100k events/sec minimum
		t.Errorf("Publish throughput too low: %.0f events/sec", publishThroughput)
	}

	if atomic.LoadInt64(&eventCount) != expectedEvents {
		t.Errorf("Not all events processed: expected %d, got %d", expectedEvents, atomic.LoadInt64(&eventCount))
	}

	t.Logf("Published %d events in %v (%.0f events/sec)", numEvents, publishDuration, publishThroughput)
	t.Logf("Total processing time: %v", totalDuration)
	t.Logf("Events processed: %d", atomic.LoadInt64(&eventCount))
}

func TestEventBusMemoryUsage(t *testing.T) {
	// Test that event bus doesn't cause memory leaks
	config := DefaultBusConfig()
	config.BufferSize = 1000
	eventBus := NewEventBus(config)
	defer eventBus.Close()

	// Add subscribers
	eventBus.Subscribe(EventTaskStarted, func(event Event) {})
	eventBus.Subscribe(EventTaskCompleted, func(event Event) {})

	// Publish and process many events
	numRounds := 100
	eventsPerRound := 1000

	for round := 0; round < numRounds; round++ {
		for i := 0; i < eventsPerRound; i++ {
			event := &TaskStartedEvent{
				BaseEvent: BaseEvent{
					EventType: EventTaskStarted,
					Time:      time.Now(),
					ID:        GenerateTaskID(),
					Ctx:       context.Background(),
				},
			}
			eventBus.PublishAsync(event)
		}

		// Allow processing
		time.Sleep(10 * time.Millisecond)

		// Check stats
		stats := eventBus.Stats()
		if stats.QueueSize > config.BufferSize/2 {
			t.Errorf("Queue size growing too large: %d", stats.QueueSize)
		}
	}

	// Final check - queue should be mostly empty
	time.Sleep(100 * time.Millisecond)
	stats := eventBus.Stats()
	if stats.QueueSize > 100 {
		t.Errorf("Final queue size too large: %d", stats.QueueSize)
	}

	totalPublished := stats.TotalPublished
	expectedTotal := int64(numRounds * eventsPerRound)
	if totalPublished < expectedTotal {
		t.Errorf("Not all events published: expected %d, got %d", expectedTotal, totalPublished)
	}

	t.Logf("Final stats: Published=%d, Dropped=%d, QueueSize=%d",
		stats.TotalPublished, stats.TotalDropped, stats.QueueSize)
}

func TestEventBusErrorHandling(t *testing.T) {
	config := DefaultBusConfig()
	config.RecoverPanics = true
	eventBus := NewEventBus(config)
	defer eventBus.Close()

	var normalEventCount int64

	// Add a subscriber that panics
	eventBus.Subscribe(EventTaskStarted, func(event Event) {
		if event.TaskID() == "panic-task" {
			panic("Test panic")
		}
		atomic.AddInt64(&normalEventCount, 1)
	})

	// Add another subscriber to ensure others still work
	eventBus.Subscribe(EventTaskStarted, func(event Event) {
		atomic.AddInt64(&normalEventCount, 1)
	})

	// Publish normal event
	normalEvent := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "normal-task",
			Ctx:       context.Background(),
		},
	}
	eventBus.Publish(normalEvent)

	// Publish panic-inducing event
	panicEvent := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "panic-task",
			Ctx:       context.Background(),
		},
	}
	eventBus.Publish(panicEvent)

	// Publish another normal event
	eventBus.Publish(normalEvent)

	time.Sleep(50 * time.Millisecond)

	// Should have processed normal events despite panic
	if atomic.LoadInt64(&normalEventCount) < 3 { // 2 normal events * 2 subscribers + 1 from non-panicking subscriber
		t.Errorf("Expected at least 3 normal events processed, got %d", atomic.LoadInt64(&normalEventCount))
	}
}

func TestEventBusSubscriptionManagement(t *testing.T) {
	config := DefaultBusConfig()
	eventBus := NewEventBus(config)
	defer eventBus.Close()

	// Test dynamic subscription management
	var events []Event
	var mu sync.Mutex

	collector := func(event Event) {
		mu.Lock()
		events = append(events, event)
		mu.Unlock()
	}

	// Add multiple subscriptions
	sub1 := eventBus.Subscribe(EventTaskStarted, collector)
	sub2 := eventBus.Subscribe(EventTaskCompleted, collector)
	sub3 := eventBus.SubscribeAll(collector)

	// Publish test events
	startEvent := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "test-task",
			Ctx:       context.Background(),
		},
	}

	completeEvent := &TaskCompletedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskCompleted,
			Time:      time.Now(),
			ID:        "test-task",
			Ctx:       context.Background(),
		},
	}

	eventBus.Publish(startEvent)
	eventBus.Publish(completeEvent)

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	initialCount := len(events)
	events = nil // Reset
	mu.Unlock()

	// Should have received events from multiple subscriptions
	if initialCount < 4 { // TaskStarted: specific + all, TaskCompleted: specific + all
		t.Errorf("Expected at least 4 events, got %d", initialCount)
	}

	// Remove specific subscriptions
	if !eventBus.Unsubscribe(sub1) {
		t.Error("Failed to unsubscribe sub1")
	}
	if !eventBus.Unsubscribe(sub2) {
		t.Error("Failed to unsubscribe sub2")
	}

	// Publish again
	eventBus.Publish(startEvent)
	eventBus.Publish(completeEvent)

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	finalCount := len(events)
	mu.Unlock()

	// Should only receive events from SubscribeAll
	if finalCount != 2 {
		t.Errorf("Expected 2 events after unsubscribe, got %d", finalCount)
	}

	// Remove remaining subscription
	if !eventBus.Unsubscribe(sub3) {
		t.Error("Failed to unsubscribe sub3")
	}

	// Verify stats
	stats := eventBus.Stats()
	if stats.ActiveSubscribers != 0 {
		t.Errorf("Expected 0 active subscribers, got %d", stats.ActiveSubscribers)
	}
}
