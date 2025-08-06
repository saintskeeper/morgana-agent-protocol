package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/adapter"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Orchestrator manages the execution of multiple agent tasks
type Orchestrator struct {
	adapter        *adapter.Adapter
	maxConcurrency int
	tracer         trace.Tracer
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
	}
}

// RunSequential executes tasks one after another
func (o *Orchestrator) RunSequential(ctx context.Context, tasks []adapter.Task) []adapter.Result {
	ctx, span := o.tracer.Start(ctx, "orchestrator.sequential",
		trace.WithAttributes(
			attribute.Int("task.count", len(tasks)),
		),
	)
	defer span.End()

	results := make([]adapter.Result, len(tasks))
	for i, task := range tasks {
		taskCtx, taskSpan := o.tracer.Start(ctx, fmt.Sprintf("task.%d", i),
			trace.WithAttributes(
				attribute.String("agent.type", task.AgentType),
				attribute.Int("task.index", i),
			),
		)
		results[i] = o.adapter.Execute(taskCtx, task)
		taskSpan.End()
	}
	return results
}

// RunParallel executes tasks concurrently with controlled concurrency
func (o *Orchestrator) RunParallel(ctx context.Context, tasks []adapter.Task) []adapter.Result {
	ctx, span := o.tracer.Start(ctx, "orchestrator.parallel",
		trace.WithAttributes(
			attribute.Int("task.count", len(tasks)),
			attribute.Int("concurrency.limit", o.maxConcurrency),
		),
	)
	defer span.End()

	startTime := time.Now()
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

	// Record parallel execution metrics
	execTime := time.Since(startTime)
	span.SetAttributes(
		attribute.Int64("execution.duration_ms", execTime.Milliseconds()),
		attribute.Float64("execution.avg_time_per_task_ms", float64(execTime.Milliseconds())/float64(len(tasks))),
	)

	return results
}
