package events

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkEventBusPublishSync measures synchronous publish performance
func BenchmarkEventBusPublishSync(b *testing.B) {
	config := DefaultBusConfig()
	bus := NewEventBus(config)
	defer bus.Close()

	// Add a simple subscriber
	bus.Subscribe(EventTaskStarted, func(event Event) {
		// Minimal processing
	})

	event := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "benchmark-task",
			Ctx:       context.Background(),
		},
		AgentType: "benchmark-agent",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(event)
	}
}

// BenchmarkEventBusPublishAsync measures asynchronous publish performance
func BenchmarkEventBusPublishAsync(b *testing.B) {
	config := DefaultBusConfig()
	config.BufferSize = 100000 // Large buffer to avoid drops
	bus := NewEventBus(config)
	defer bus.Close()

	// Add a simple subscriber
	bus.Subscribe(EventTaskStarted, func(event Event) {
		// Minimal processing
	})

	event := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "benchmark-task",
			Ctx:       context.Background(),
		},
		AgentType: "benchmark-agent",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !bus.PublishAsync(event) {
			b.Error("Failed to publish async event")
		}
	}
}

// BenchmarkCircularBufferPush measures circular buffer push performance
func BenchmarkCircularBufferPush(b *testing.B) {
	buffer := NewCircularBuffer(100000)

	event := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "benchmark-task",
			Ctx:       context.Background(),
		},
	}

	// Start consumer to prevent buffer from filling up
	go func() {
		for {
			if buffer.Pop() == nil {
				runtime.Gosched()
			}
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for !buffer.Push(event) {
			runtime.Gosched()
		}
	}
}

// BenchmarkCircularBufferPop measures circular buffer pop performance
func BenchmarkCircularBufferPop(b *testing.B) {
	buffer := NewCircularBuffer(100000)

	event := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "benchmark-task",
			Ctx:       context.Background(),
		},
	}

	// Pre-fill buffer
	for i := 0; i < 50000; i++ {
		buffer.Push(event)
	}

	// Start producer to keep buffer filled
	go func() {
		for {
			if !buffer.Push(event) {
				runtime.Gosched()
			}
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for buffer.Pop() == nil {
			runtime.Gosched()
		}
	}
}

// BenchmarkEventBusConcurrentPublish measures concurrent publish performance
func BenchmarkEventBusConcurrentPublish(b *testing.B) {
	config := DefaultBusConfig()
	config.BufferSize = 100000
	config.Workers = runtime.NumCPU()
	bus := NewEventBus(config)
	defer bus.Close()

	// Add multiple subscribers to create realistic load
	for i := 0; i < 5; i++ {
		bus.Subscribe(EventTaskStarted, func(event Event) {
			// Simulate minimal processing
			_ = event.TaskID()
		})
	}

	numGoroutines := runtime.NumCPU()
	var wg sync.WaitGroup

	b.ResetTimer()
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			event := &TaskStartedEvent{
				BaseEvent: BaseEvent{
					EventType: EventTaskStarted,
					Time:      time.Now(),
					ID:        GenerateTaskID(),
					Ctx:       context.Background(),
				},
			}

			for i := 0; i < b.N/numGoroutines; i++ {
				bus.PublishAsync(event)
			}
		}()
	}
	wg.Wait()
}

// BenchmarkTaskIDGeneration measures task ID generation performance
func BenchmarkTaskIDGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateTaskID()
	}
}

// BenchmarkEventCreation measures event creation performance
func BenchmarkEventCreation(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &TaskStartedEvent{
			BaseEvent: BaseEvent{
				EventType: EventTaskStarted,
				Time:      time.Now(),
				ID:        "benchmark-task",
				Ctx:       ctx,
			},
			AgentType: "benchmark-agent",
		}
	}
}

// BenchmarkEventBusOverhead measures overhead compared to direct function calls
func BenchmarkEventBusOverhead(b *testing.B) {
	// First, measure direct function call
	var counter int64
	directFunc := func() {
		atomic.AddInt64(&counter, 1)
	}

	b.Run("DirectCall", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			directFunc()
		}
	})

	// Then measure via event bus
	config := DefaultBusConfig()
	bus := NewEventBus(config)
	defer bus.Close()

	bus.Subscribe(EventTaskStarted, func(event Event) {
		atomic.AddInt64(&counter, 1)
	})

	event := &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        "benchmark-task",
			Ctx:       context.Background(),
		},
	}

	b.Run("EventBus", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bus.Publish(event)
		}
	})
}

// BenchmarkHighThroughput simulates high-throughput scenario
func BenchmarkHighThroughput(b *testing.B) {
	config := DefaultBusConfig()
	config.BufferSize = 1000000
	config.Workers = runtime.NumCPU() * 2
	bus := NewEventBus(config)
	defer bus.Close()

	// Add realistic subscribers that do some work
	var processedEvents int64
	for i := 0; i < 10; i++ {
		bus.Subscribe(EventTaskStarted, func(event Event) {
			// Simulate some processing time
			_ = event.TaskID()
			_ = event.Timestamp()
			atomic.AddInt64(&processedEvents, 1)
		})
	}

	numPublishers := runtime.NumCPU()
	eventsPerPublisher := b.N / numPublishers

	var wg sync.WaitGroup

	b.ResetTimer()
	start := time.Now()

	for p := 0; p < numPublishers; p++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < eventsPerPublisher; i++ {
				event := &TaskStartedEvent{
					BaseEvent: BaseEvent{
						EventType: EventTaskStarted,
						Time:      time.Now(),
						ID:        GenerateTaskID(),
						Ctx:       context.Background(),
					},
					AgentType: "high-throughput-agent",
				}

				// Retry until published
				for !bus.PublishAsync(event) {
					runtime.Gosched()
				}
			}
		}()
	}

	wg.Wait()

	// Wait for all events to be processed
	deadline := time.Now().Add(5 * time.Second)
	expectedEvents := int64(b.N * 10) // 10 subscribers per event

	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&processedEvents) >= expectedEvents {
			break
		}
		time.Sleep(time.Millisecond)
	}

	duration := time.Since(start)
	throughput := float64(b.N) / duration.Seconds()

	b.ReportMetric(throughput, "events/sec")
	b.ReportMetric(float64(atomic.LoadInt64(&processedEvents)), "total_processed")
}
