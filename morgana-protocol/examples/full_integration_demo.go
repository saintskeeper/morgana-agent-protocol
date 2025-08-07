package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/adapter"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/orchestrator"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/prompt"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/pkg/task"
	"go.opentelemetry.io/otel"
)

func main() {
	fmt.Println("🔗 Morgana Protocol Full Integration Demo")
	fmt.Println("========================================")

	// Create event bus
	config := events.DefaultBusConfig()
	config.Debug = true
	config.BufferSize = 1000
	eventBus := events.NewEventBus(config)
	defer eventBus.Close()

	// Set up comprehensive event monitoring
	setupEventMonitoring(eventBus)

	// Create components
	promptLoader := prompt.NewLoader("internal/adapter/testdata") // Mock path
	taskClient := &task.Client{}                                  // This would be configured properly in real usage
	tracer := otel.Tracer("morgana-demo")

	// Create adapter with event integration
	adapterInstance := adapter.New(promptLoader, taskClient, tracer)
	adapterInstance.SetEventBus(eventBus)

	// Create orchestrator with event integration
	orchestratorInstance := orchestrator.New(adapterInstance, 3, tracer)
	orchestratorInstance.SetEventBus(eventBus)

	// Demo: Execute tasks sequentially with event tracking
	fmt.Println("\n🔄 Running Sequential Task Demo...")
	runSequentialDemo(orchestratorInstance)

	// Demo: Execute tasks in parallel with event tracking
	fmt.Println("\n⚡ Running Parallel Task Demo...")
	runParallelDemo(orchestratorInstance)

	// Wait for all events to be processed
	time.Sleep(200 * time.Millisecond)

	// Show final statistics
	showFinalStats(eventBus)

	fmt.Println("\n✅ Full integration demo completed!")
}

func setupEventMonitoring(eventBus events.EventBus) {
	// Task lifecycle monitoring
	eventBus.Subscribe(events.EventTaskStarted, func(event events.Event) {
		if e, ok := event.(*events.TaskStartedEvent); ok {
			log.Printf("📋 TASK_STARTED: %s [%s] - %s", e.TaskID(), e.AgentType, truncateString(e.Prompt, 50))
		}
	})

	eventBus.Subscribe(events.EventTaskProgress, func(event events.Event) {
		if e, ok := event.(*events.TaskProgressEvent); ok {
			log.Printf("⏳ TASK_PROGRESS: %s [%s] %.0f%% - %s", e.TaskID(), e.Stage, e.Progress*100, e.Message)
		}
	})

	eventBus.Subscribe(events.EventTaskCompleted, func(event events.Event) {
		if e, ok := event.(*events.TaskCompletedEvent); ok {
			log.Printf("✅ TASK_COMPLETED: %s [%s] %s - %d bytes output",
				e.TaskID(), e.AgentType, e.Duration, e.OutputLength)
		}
	})

	eventBus.Subscribe(events.EventTaskFailed, func(event events.Event) {
		if e, ok := event.(*events.TaskFailedEvent); ok {
			log.Printf("❌ TASK_FAILED: %s [%s] %s - %s", e.TaskID(), e.AgentType, e.Duration, e.Error)
		}
	})

	// Orchestrator monitoring
	eventBus.Subscribe(events.EventOrchestratorStarted, func(event events.Event) {
		if e, ok := event.(*events.OrchestratorStartedEvent); ok {
			log.Printf("🚀 ORCHESTRATOR_STARTED: %s mode - %d tasks", e.Mode, e.TaskCount)
		}
	})

	eventBus.Subscribe(events.EventOrchestratorCompleted, func(event events.Event) {
		if e, ok := event.(*events.OrchestratorCompletedEvent); ok {
			log.Printf("🏁 ORCHESTRATOR_COMPLETED: %s - %d/%d succeeded in %s",
				e.Mode, e.SuccessCount, e.TaskCount, e.Duration)
		}
	})

	// Performance monitoring - track events per second
	var eventCount int64
	var lastCheck time.Time = time.Now()

	eventBus.SubscribeAll(func(event events.Event) {
		eventCount++
		if time.Since(lastCheck) >= time.Second {
			log.Printf("📊 Events/sec: %d", eventCount)
			eventCount = 0
			lastCheck = time.Now()
		}
	})
}

func runSequentialDemo(orch *orchestrator.Orchestrator) {
	tasks := []adapter.Task{
		{
			AgentType: "code-implementer",
			Prompt:    "Implement a simple calculator function",
			Options:   map[string]interface{}{"language": "go"},
		},
		{
			AgentType: "test-specialist",
			Prompt:    "Write unit tests for the calculator function",
			Options:   map[string]interface{}{"coverage": "90%"},
		},
		{
			AgentType: "validation-expert",
			Prompt:    "Review code quality and suggest improvements",
			Options:   map[string]interface{}{"standards": "go-best-practices"},
		},
	}

	ctx := context.Background()
	results := orch.RunSequential(ctx, tasks)

	fmt.Printf("Sequential execution completed: %d results\n", len(results))
	for i, result := range results {
		if result.Error != "" {
			fmt.Printf("  Task %d: ERROR - %s\n", i+1, result.Error)
		} else {
			fmt.Printf("  Task %d: SUCCESS - %d bytes output\n", i+1, len(result.Output))
		}
	}
}

func runParallelDemo(orch *orchestrator.Orchestrator) {
	tasks := []adapter.Task{
		{
			AgentType: "code-implementer",
			Prompt:    "Implement user authentication system",
			Options:   map[string]interface{}{"security": "jwt"},
		},
		{
			AgentType: "code-implementer",
			Prompt:    "Implement data validation middleware",
			Options:   map[string]interface{}{"framework": "gin"},
		},
		{
			AgentType: "sprint-planner",
			Prompt:    "Plan next sprint for authentication feature",
			Options:   map[string]interface{}{"team_size": 4},
		},
		{
			AgentType: "test-specialist",
			Prompt:    "Design integration test strategy",
			Options:   map[string]interface{}{"scope": "auth-system"},
		},
	}

	ctx := context.Background()
	results := orch.RunParallel(ctx, tasks)

	fmt.Printf("Parallel execution completed: %d results\n", len(results))
	successCount := 0
	for i, result := range results {
		if result.Error != "" {
			fmt.Printf("  Task %d: ERROR - %s\n", i+1, result.Error)
		} else {
			fmt.Printf("  Task %d: SUCCESS - %d bytes output\n", i+1, len(result.Output))
			successCount++
		}
	}
	fmt.Printf("Success rate: %d/%d (%.1f%%)\n", successCount, len(results),
		float64(successCount)/float64(len(results))*100)
}

func showFinalStats(eventBus events.EventBus) {
	stats := eventBus.Stats()
	fmt.Printf("\n📈 Final Event Bus Statistics:\n")
	fmt.Printf("  Total Events Published: %d\n", stats.TotalPublished)
	fmt.Printf("  Total Events Dropped: %d\n", stats.TotalDropped)
	fmt.Printf("  Drop Rate: %.2f%%\n", float64(stats.TotalDropped)/float64(stats.TotalPublished)*100)
	fmt.Printf("  Active Subscribers: %d\n", stats.ActiveSubscribers)
	fmt.Printf("  Queue Utilization: %d/%d (%.1f%%)\n",
		stats.QueueSize, stats.QueueCapacity,
		float64(stats.QueueSize)/float64(stats.QueueCapacity)*100)

	fmt.Printf("\n  Subscriber Distribution:\n")
	for eventType, count := range stats.SubscribersByType {
		fmt.Printf("    %s: %d subscribers\n", eventType, count)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func init() {
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}
