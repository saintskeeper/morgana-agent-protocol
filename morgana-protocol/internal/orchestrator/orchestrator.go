package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/adapter"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Orchestrator manages the execution of multiple agent tasks
type Orchestrator struct {
	adapter        *adapter.Adapter
	maxConcurrency int
	tracer         trace.Tracer
	eventBus       events.EventBus
}

// New creates a new Orchestrator
func New(adapter *adapter.Adapter, maxConcurrency int, tracer trace.Tracer) *Orchestrator {
	if maxConcurrency <= 0 {
		maxConcurrency = 5 // Default
	}
	return &Orchestrator{
		adapter:        adapter,
		maxConcurrency: maxConcurrency,
		tracer:         tracer,
		eventBus:       nil, // Will be set via SetEventBus
	}
}

// SetEventBus sets the event bus for publishing orchestrator events
func (o *Orchestrator) SetEventBus(eventBus events.EventBus) {
	o.eventBus = eventBus
}

// publishEvent publishes an event to the event bus if available
func (o *Orchestrator) publishEvent(event events.Event) {
	if o.eventBus != nil {
		o.eventBus.PublishAsync(event)
	}
}

// RunSequential executes tasks one after another
func (o *Orchestrator) RunSequential(ctx context.Context, tasks []adapter.Task) []adapter.Result {
	// Generate orchestrator task ID
	orchID := events.GenerateTaskID()
	ctx = events.SetTaskIDInContext(ctx, orchID)

	ctx, span := o.tracer.Start(ctx, "orchestrator.sequential",
		trace.WithAttributes(
			attribute.Int("task.count", len(tasks)),
		),
	)
	defer span.End()

	startTime := time.Now()

	// Publish orchestrator started event
	o.publishEvent(&events.OrchestratorStartedEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventOrchestratorStarted,
			Time:      time.Now(),
			ID:        orchID,
			Ctx:       ctx,
		},
		Mode:      "sequential",
		TaskCount: len(tasks),
	})

	results := make([]adapter.Result, len(tasks))
	successCount := 0
	failureCount := 0

	for i, task := range tasks {
		taskCtx, taskSpan := o.tracer.Start(ctx, fmt.Sprintf("task.%d", i),
			trace.WithAttributes(
				attribute.String("agent.type", task.AgentType),
				attribute.Int("task.index", i),
			),
		)
		results[i] = o.adapter.Execute(taskCtx, task)
		if results[i].Error == "" {
			successCount++
		} else {
			failureCount++
		}
		taskSpan.End()
	}

	// Publish orchestrator completed event
	duration := time.Since(startTime)
	o.publishEvent(&events.OrchestratorCompletedEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventOrchestratorCompleted,
			Time:      time.Now(),
			ID:        orchID,
			Ctx:       ctx,
		},
		Mode:         "sequential",
		TaskCount:    len(tasks),
		SuccessCount: successCount,
		FailureCount: failureCount,
		Duration:     duration,
	})

	return results
}

// RunParallel executes tasks concurrently with controlled concurrency
func (o *Orchestrator) RunParallel(ctx context.Context, tasks []adapter.Task) []adapter.Result {
	// Generate orchestrator task ID
	orchID := events.GenerateTaskID()
	ctx = events.SetTaskIDInContext(ctx, orchID)

	ctx, span := o.tracer.Start(ctx, "orchestrator.parallel",
		trace.WithAttributes(
			attribute.Int("task.count", len(tasks)),
			attribute.Int("concurrency.limit", o.maxConcurrency),
		),
	)
	defer span.End()

	startTime := time.Now()

	// Publish orchestrator started event
	o.publishEvent(&events.OrchestratorStartedEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventOrchestratorStarted,
			Time:      time.Now(),
			ID:        orchID,
			Ctx:       ctx,
		},
		Mode:           "parallel",
		TaskCount:      len(tasks),
		MaxConcurrency: o.maxConcurrency,
	})

	var wg sync.WaitGroup
	results := make([]adapter.Result, len(tasks))
	semaphore := make(chan struct{}, o.maxConcurrency)

	// Fill semaphore
	for i := 0; i < o.maxConcurrency; i++ {
		semaphore <- struct{}{}
	}

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t adapter.Task) {
			defer wg.Done()

			// Create span for this goroutine
			taskCtx, taskSpan := o.tracer.Start(ctx, fmt.Sprintf("parallel.task.%d", idx),
				trace.WithAttributes(
					attribute.String("agent.type", t.AgentType),
					attribute.Int("task.index", idx),
				),
			)
			defer taskSpan.End()

			// Wait for semaphore
			selectCtx, waitSpan := o.tracer.Start(taskCtx, "semaphore.wait")
			<-semaphore // Acquire
			waitSpan.End()

			// Execute task
			results[idx] = o.adapter.Execute(selectCtx, t)

			// Release semaphore
			semaphore <- struct{}{} // Release
		}(i, task)
	}

	wg.Wait()

	// Count successes and failures
	successCount := 0
	failureCount := 0
	for _, result := range results {
		if result.Error == "" {
			successCount++
		} else {
			failureCount++
		}
	}

	// Record parallel execution metrics
	execTime := time.Since(startTime)
	span.SetAttributes(
		attribute.Int64("execution.duration_ms", execTime.Milliseconds()),
		attribute.Float64("execution.avg_time_per_task_ms", float64(execTime.Milliseconds())/float64(len(tasks))),
	)

	// Publish orchestrator completed event
	o.publishEvent(&events.OrchestratorCompletedEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventOrchestratorCompleted,
			Time:      time.Now(),
			ID:        orchID,
			Ctx:       ctx,
		},
		Mode:         "parallel",
		TaskCount:    len(tasks),
		SuccessCount: successCount,
		FailureCount: failureCount,
		Duration:     execTime,
	})

	return results
}
