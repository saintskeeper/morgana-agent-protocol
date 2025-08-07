package main

import (
	"context"
	"log"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
)

func main() {
	// Create context
	ctx := context.Background()

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Create TUI config
	config := tui.CreateDevelopmentConfig()

	// Start TUI asynchronously
	tuiInstance, err := tui.RunAsync(ctx, eventBus, config)
	if err != nil {
		log.Fatalf("Failed to start TUI: %v", err)
	}

	// Simulate some tasks for demonstration
	go simulateTasks(ctx, eventBus)

	// Wait for TUI to finish (Ctrl+C to quit)
	log.Println("TUI started. Press 'q' or Ctrl+C to quit.")

	// In a real application, this would wait for a signal or context cancellation
	select {
	case <-ctx.Done():
		break
	case <-time.After(5 * time.Minute):
		log.Println("Demo timeout reached")
		break
	}

	// Stop TUI
	if err := tuiInstance.Stop(); err != nil {
		log.Printf("Error stopping TUI: %v", err)
	}

	log.Println("TUI demo completed")
}

// simulateTasks creates fake task events for demonstration
func simulateTasks(ctx context.Context, eventBus events.EventBus) {
	taskTypes := []string{
		"code-implementer",
		"sprint-planner",
		"test-specialist",
		"validation-expert",
	}

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(i*500) * time.Millisecond):
			// Create a task
			agentType := taskTypes[i%len(taskTypes)]
			taskID := events.GenerateTaskID()

			// Task started
			startEvent := events.NewTaskStartedEvent(
				ctx, taskID, agentType,
				"Demo task implementation",
				nil, 0, "", "medium",
				2*time.Minute,
			)
			eventBus.PublishAsync(startEvent)

			// Simulate progress over time
			go simulateTaskProgress(ctx, eventBus, taskID, agentType)
		}
	}
}

// simulateTaskProgress simulates a task execution with progress updates
func simulateTaskProgress(ctx context.Context, eventBus events.EventBus, taskID, agentType string) {
	stages := []struct {
		stage    string
		message  string
		progress float64
		duration time.Duration
	}{
		{"validation", "Validating agent configuration", 0.1, 200 * time.Millisecond},
		{"prompt_load", "Loading prompt template", 0.2, 300 * time.Millisecond},
		{"model_selection", "Selecting optimal model", 0.3, 400 * time.Millisecond},
		{"execution", "Executing task logic", 0.6, 1 * time.Second},
		{"post_processing", "Processing results", 0.9, 500 * time.Millisecond},
	}

	startTime := time.Now()

	for _, stage := range stages {
		select {
		case <-ctx.Done():
			return
		case <-time.After(stage.duration):
			// Progress event
			progressEvent := events.NewTaskProgressEvent(
				ctx, taskID, agentType,
				stage.stage, stage.message,
				stage.progress,
				time.Since(startTime),
			)
			eventBus.PublishAsync(progressEvent)
		}
	}

	// Task completion (80% success rate)
	select {
	case <-ctx.Done():
		return
	case <-time.After(500 * time.Millisecond):
		totalDuration := time.Since(startTime)

		if time.Now().UnixNano()%5 == 0 {
			// Task failed (20% chance)
			failEvent := events.NewTaskFailedEvent(
				ctx, taskID, agentType,
				"Simulated task failure",
				totalDuration,
				"execution", 0,
			)
			eventBus.PublishAsync(failEvent)
		} else {
			// Task completed successfully
			completeEvent := events.NewTaskCompletedEvent(
				ctx, taskID, agentType,
				"Task completed successfully with demo output",
				totalDuration,
				"claude-3-sonnet",
			)
			eventBus.PublishAsync(completeEvent)
		}
	}
}
