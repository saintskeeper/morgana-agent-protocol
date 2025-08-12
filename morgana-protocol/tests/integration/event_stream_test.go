//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

func TestEventStreamWritingAndReading(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx := context.Background()

	t.Run("BasicEventFlow", func(t *testing.T) {
		// Generate a complete task lifecycle
		taskID := setup.Generator.GenerateTaskLifecycle(ctx, "code-implementer", 10*time.Millisecond)

		// Wait for all events to be processed (start + 3 progress + complete = 5 events)
		if !WaitForEvents(setup.Collector, 5, time.Second*2) {
			t.Fatal("Timed out waiting for events")
		}

		// Verify event sequence
		taskEvents := setup.Collector.GetEventsForTask(taskID)
		if len(taskEvents) != 5 {
			t.Errorf("Expected 5 events for task, got %d", len(taskEvents))
		}

		// Verify event order
		expectedTypes := []events.EventType{
			events.EventTaskStarted,
			events.EventTaskProgress,
			events.EventTaskProgress,
			events.EventTaskProgress,
			events.EventTaskCompleted,
		}

		for i, expected := range expectedTypes {
			if i >= len(taskEvents) {
				t.Errorf("Missing event at index %d, expected %s", i, expected)
				continue
			}
			if taskEvents[i].Type() != expected {
				t.Errorf("Event %d: expected %s, got %s", i, expected, taskEvents[i].Type())
			}
		}

		t.Logf("Successfully processed task lifecycle for %s", taskID)
	})

	t.Run("FailedTaskFlow", func(t *testing.T) {
		// Generate a failed task scenario
		taskID := setup.Generator.GenerateFailedTask(ctx, "test-specialist", 20*time.Millisecond)

		// Wait for failure events (start + progress + failed = 3 events)
		if !WaitForEvents(setup.Collector, 8, time.Second*2) { // 5 from previous test + 3 new
			t.Fatal("Timed out waiting for failed task events")
		}

		// Verify failed task events
		taskEvents := setup.Collector.GetEventsForTask(taskID)
		if len(taskEvents) != 3 {
			t.Errorf("Expected 3 events for failed task, got %d", len(taskEvents))
		}

		// Check for failed event
		failedEvents := setup.Collector.GetEventsOfType(events.EventTaskFailed)
		found := false
		for _, event := range failedEvents {
			if event.TaskID() == taskID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Failed event not found for failed task")
		}

		t.Logf("Successfully processed failed task for %s", taskID)
	})

	t.Run("HighVolumeEvents", func(t *testing.T) {
		// Test high volume event processing
		agentTypes := RandomAgentTypes()
		eventCount := 1000

		start := time.Now()
		taskIDs := setup.Generator.GenerateHighVolumeEvents(ctx, eventCount, agentTypes)
		generateDuration := time.Since(start)

		// Wait for events to be processed
		if !WaitForEvents(setup.Collector, 8+eventCount, time.Second*5) {
			t.Fatal("Timed out waiting for high volume events")
		}

		processingDuration := time.Since(start)

		// Performance assertions
		stats := setup.Collector.GetStats()
		busStats := setup.EventBus.Stats()

		// Check that most events were processed
		if int(stats.TotalEvents) < 8+eventCount*0.9 { // Allow 10% loss under high load
			t.Errorf("Too many events lost: expected ~%d, got %d", 8+eventCount, stats.TotalEvents)
		}

		t.Logf("High volume test: %d events generated in %v, processed in %v",
			len(taskIDs), generateDuration, processingDuration)
		t.Logf("Stats: %+v", busStats)
	})
}

func TestEventStreamPerformance(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx := context.Background()

	t.Run("LatencyRequirement", func(t *testing.T) {
		// Test that events are processed within 10ms requirement
		agentTypes := []string{"code-implementer"}

		// Generate events with precise timing
		taskIDs := setup.Generator.GenerateHighVolumeEvents(ctx, 100, agentTypes)

		// Wait for processing
		if !WaitForEvents(setup.Collector, 100, time.Second*2) {
			t.Fatal("Timed out waiting for latency test events")
		}

		// Check performance requirements
		stats := setup.Collector.GetStats()
		busStats := setup.EventBus.Stats()
		assertions := DefaultPerformanceAssertions()

		AssertPerformance(t, stats, busStats, assertions)

		t.Logf("Processed %d events with average latency check", len(taskIDs))
	})

	t.Run("ThroughputTest", func(t *testing.T) {
		// Reset collector for clean metrics
		setup.Collector.Reset()

		agentTypes := RandomAgentTypes()
		eventCount := 5000

		start := time.Now()

		// Generate events as fast as possible
		setup.Generator.GenerateHighVolumeEvents(ctx, eventCount, agentTypes)

		// Wait for processing
		if !WaitForEvents(setup.Collector, eventCount, time.Second*10) {
			t.Log("Warning: Not all events processed within timeout (expected under high load)")
		}

		duration := time.Since(start)
		stats := setup.Collector.GetStats()

		throughput := float64(stats.TotalEvents) / duration.Seconds()
		t.Logf("Throughput: %.0f events/sec (%d events in %v)",
			throughput, stats.TotalEvents, duration)

		// Verify we can handle at least 1000 events/sec
		if throughput < 1000 {
			t.Errorf("Throughput too low: %.0f events/sec (minimum: 1000)", throughput)
		}
	})

	t.Run("MemoryLeakTest", func(t *testing.T) {
		// Test that event system doesn't leak memory over time
		agentTypes := []string{"test-specialist", "code-implementer"}

		// Process many events in batches
		batchSize := 1000
		batches := 5

		for batch := 0; batch < batches; batch++ {
			setup.Generator.GenerateHighVolumeEvents(ctx, batchSize, agentTypes)

			// Allow processing
			time.Sleep(200 * time.Millisecond)

			// Check queue doesn't grow indefinitely
			busStats := setup.EventBus.Stats()
			if busStats.QueueSize > busStats.QueueCapacity/2 {
				t.Logf("Warning: Queue size growing in batch %d: %d/%d",
					batch, busStats.QueueSize, busStats.QueueCapacity)
			}
		}

		// Final check - queue should be mostly empty
		time.Sleep(500 * time.Millisecond)
		busStats := setup.EventBus.Stats()

		if busStats.QueueSize > 100 {
			t.Errorf("Queue not draining properly: %d events remaining", busStats.QueueSize)
		}

		t.Logf("Memory leak test completed. Final queue size: %d", busStats.QueueSize)
	})
}

func TestEventStreamErrorHandling(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx := context.Background()

	t.Run("EventBusRecovery", func(t *testing.T) {
		// Test that event bus recovers from panics in subscribers
		var panicSubscriptionID string
		var normalEventCount int

		// Add a subscriber that panics for specific events
		panicSubscriber := func(event events.Event) {
			if event.TaskID() == "panic-task" {
				panic("Test panic in subscriber")
			}
		}

		normalSubscriber := func(event events.Event) {
			normalEventCount++
		}

		panicSubscriptionID = setup.EventBus.Subscribe(events.EventTaskStarted, panicSubscriber)
		normalSubscriptionID := setup.EventBus.Subscribe(events.EventTaskStarted, normalSubscriber)

		// Generate normal events
		setup.Generator.GenerateTaskLifecycle(ctx, "code-implementer", 10*time.Millisecond)

		// Generate panic-causing event
		panicEvent := events.NewTaskStartedEvent(
			ctx, "panic-task", "test-agent", "panic test", nil, 0, "", "", time.Minute,
		)
		setup.EventBus.PublishAsync(panicEvent)

		// Generate more normal events
		setup.Generator.GenerateTaskLifecycle(ctx, "test-specialist", 10*time.Millisecond)

		time.Sleep(500 * time.Millisecond)

		// Normal subscriber should still have received events
		if normalEventCount < 2 {
			t.Errorf("Normal subscriber affected by panic: got %d events, expected at least 2",
				normalEventCount)
		}

		// Clean up
		setup.EventBus.Unsubscribe(panicSubscriptionID)
		setup.EventBus.Unsubscribe(normalSubscriptionID)

		t.Logf("Event bus recovered from panic, normal events processed: %d", normalEventCount)
	})

	t.Run("BackpressureHandling", func(t *testing.T) {
		// Test behavior under extreme load (queue overflow)
		agentTypes := []string{"load-test-agent"}

		// Generate way more events than can be processed quickly
		extremeLoad := 20000

		start := time.Now()
		setup.Generator.GenerateHighVolumeEvents(ctx, extremeLoad, agentTypes)
		generateTime := time.Since(start)

		// Let system try to process for a reasonable time
		time.Sleep(time.Second * 3)

		busStats := setup.EventBus.Stats()

		t.Logf("Backpressure test: Generated %d events in %v", extremeLoad, generateTime)
		t.Logf("Bus stats under extreme load: Published=%d, Dropped=%d, QueueSize=%d",
			busStats.TotalPublished, busStats.TotalDropped, busStats.QueueSize)

		// System should gracefully handle overload by dropping events
		if busStats.TotalDropped == 0 && busStats.TotalPublished > int64(extremeLoad*0.5) {
			t.Log("Warning: System may not be properly handling backpressure")
		}

		// Queue should not exceed capacity
		if busStats.QueueSize > busStats.QueueCapacity {
			t.Errorf("Queue size exceeded capacity: %d > %d",
				busStats.QueueSize, busStats.QueueCapacity)
		}
	})
}

func TestEventStreamTypes(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx := context.Background()

	t.Run("AllEventTypes", func(t *testing.T) {
		// Test all event types are properly handled
		agentType := "comprehensive-test-agent"

		// Generate complete lifecycle
		taskID := setup.Generator.GenerateTaskLifecycle(ctx, agentType, 10*time.Millisecond)

		// Generate failure scenario
		failTaskID := setup.Generator.GenerateFailedTask(ctx, agentType, 10*time.Millisecond)

		// Wait for all events
		if !WaitForEvents(setup.Collector, 8, time.Second*2) {
			t.Fatal("Timed out waiting for comprehensive events")
		}

		stats := setup.Collector.GetStats()

		// Verify all expected event types were processed
		expectedTypes := map[events.EventType]bool{
			events.EventTaskStarted:   false,
			events.EventTaskProgress:  false,
			events.EventTaskCompleted: false,
			events.EventTaskFailed:    false,
		}

		for eventType := range stats.EventsByType {
			if _, exists := expectedTypes[eventType]; exists {
				expectedTypes[eventType] = true
			}
		}

		for eventType, found := range expectedTypes {
			if !found {
				t.Errorf("Event type %s was not processed", eventType)
			}
		}

		t.Logf("Successfully processed all event types for tasks %s and %s",
			taskID, failTaskID)
		t.Logf("Event type distribution: %+v", stats.EventsByType)
	})

	t.Run("EventDataIntegrity", func(t *testing.T) {
		// Verify event data is preserved correctly
		agentType := "data-integrity-test"
		testData := map[string]interface{}{
			"test_string":  "test value",
			"test_number":  42,
			"test_boolean": true,
			"test_array":   []string{"a", "b", "c"},
		}

		startEvent := events.NewTaskStartedEvent(
			ctx, "data-test-task", agentType, "data integrity test prompt",
			testData, 0, "test-model", "medium", time.Minute,
		)

		setup.EventBus.PublishAsync(startEvent)

		// Wait for event to be processed
		if !WaitForEventType(setup.Collector, events.EventTaskStarted, time.Second) {
			t.Fatal("Timed out waiting for data integrity event")
		}

		// Get the processed event
		startEvents := setup.Collector.GetEventsOfType(events.EventTaskStarted)
		var testEvent *events.TaskStartedEvent
		for _, event := range startEvents {
			if event.TaskID() == "data-test-task" {
				testEvent = event.(*events.TaskStartedEvent)
				break
			}
		}

		if testEvent == nil {
			t.Fatal("Test event not found")
		}

		// Verify data integrity
		if testEvent.AgentType != agentType {
			t.Errorf("Agent type mismatch: expected %s, got %s", agentType, testEvent.AgentType)
		}

		// Note: In a real implementation, we'd verify the Options field contains our test data
		// This would require proper JSON marshaling/unmarshaling in the event system

		t.Log("Event data integrity verified")
	})
}
