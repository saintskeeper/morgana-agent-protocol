package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

func main() {
	fmt.Println("🚀 Morgana Protocol Event System Demo")
	fmt.Println("=====================================")

	// Create event bus with default configuration
	config := events.DefaultBusConfig()
	config.Debug = true
	eventBus := events.NewEventBus(config)
	defer eventBus.Close()

	// Set up event subscribers
	setupSubscribers(eventBus)

	// Simulate task execution with events
	ctx := context.Background()
	taskID := events.GenerateTaskID()
	ctx = events.SetTaskIDInContext(ctx, taskID)

	fmt.Printf("🎯 Starting task: %s\n", taskID)

	// Simulate a task lifecycle
	simulateTaskExecution(ctx, eventBus, taskID)

	// Wait for events to be processed
	time.Sleep(100 * time.Millisecond)

	// Show statistics
	stats := eventBus.Stats()
	fmt.Printf("\n📊 Event Bus Statistics:\n")
	fmt.Printf("  - Total Published: %d\n", stats.TotalPublished)
	fmt.Printf("  - Total Dropped: %d\n", stats.TotalDropped)
	fmt.Printf("  - Active Subscribers: %d\n", stats.ActiveSubscribers)
	fmt.Printf("  - Current Queue Size: %d\n", stats.QueueSize)
	fmt.Printf("  - Queue Capacity: %d\n", stats.QueueCapacity)

	for eventType, count := range stats.SubscribersByType {
		fmt.Printf("  - %s subscribers: %d\n", eventType, count)
	}

	fmt.Println("\n✅ Demo completed!")
}

func setupSubscribers(eventBus events.EventBus) {
	// Subscribe to task started events
	eventBus.Subscribe(events.EventTaskStarted, func(event events.Event) {
		if startEvent, ok := event.(*events.TaskStartedEvent); ok {
			fmt.Printf("🟢 Task Started: %s (Agent: %s, Timeout: %s)\n",
				startEvent.TaskID(), startEvent.AgentType, startEvent.Timeout)
		}
	})

	// Subscribe to progress events
	eventBus.Subscribe(events.EventTaskProgress, func(event events.Event) {
		if progressEvent, ok := event.(*events.TaskProgressEvent); ok {
			fmt.Printf("🔄 Task Progress: %s - %s (%.1f%%, Duration: %s)\n",
				progressEvent.TaskID(), progressEvent.Stage, progressEvent.Progress*100, progressEvent.Duration)
		}
	})

	// Subscribe to completion events
	eventBus.Subscribe(events.EventTaskCompleted, func(event events.Event) {
		if completeEvent, ok := event.(*events.TaskCompletedEvent); ok {
			fmt.Printf("✅ Task Completed: %s (Duration: %s, Output Length: %d bytes)\n",
				completeEvent.TaskID(), completeEvent.Duration, completeEvent.OutputLength)
		}
	})

	// Subscribe to failure events
	eventBus.Subscribe(events.EventTaskFailed, func(event events.Event) {
		if failEvent, ok := event.(*events.TaskFailedEvent); ok {
			fmt.Printf("❌ Task Failed: %s - %s (Stage: %s, Duration: %s)\n",
				failEvent.TaskID(), failEvent.Error, failEvent.Stage, failEvent.Duration)
		}
	})

	// Subscribe to all events for logging
	eventBus.SubscribeAll(func(event events.Event) {
		log.Printf("[ALL_EVENTS] %s: %s", event.Type(), event.TaskID())
	})
}

func simulateTaskExecution(ctx context.Context, eventBus events.EventBus, taskID string) {
	startTime := time.Now()

	// Task started
	startEvent := events.NewTaskStartedEvent(
		ctx, taskID, "code-implementer", "Implement event system demo",
		map[string]interface{}{"complexity": "medium"}, 0, "", "medium", 2*time.Minute)
	eventBus.PublishAsync(startEvent)

	// Simulate some delays and progress updates
	time.Sleep(50 * time.Millisecond)

	// Validation progress
	progressEvent := events.NewTaskProgressEvent(
		ctx, taskID, "code-implementer", "validation", "Validating requirements", 0.2, time.Since(startTime))
	eventBus.PublishAsync(progressEvent)

	time.Sleep(100 * time.Millisecond)

	// Execution progress
	progressEvent = events.NewTaskProgressEvent(
		ctx, taskID, "code-implementer", "execution", "Implementing solution", 0.7, time.Since(startTime))
	eventBus.PublishAsync(progressEvent)

	time.Sleep(80 * time.Millisecond)

	// Task completion
	completeEvent := events.NewTaskCompletedEvent(
		ctx, taskID, "code-implementer", "Implementation completed successfully with full test coverage",
		time.Since(startTime), "claude-sonnet-4")
	eventBus.PublishAsync(completeEvent)
}
