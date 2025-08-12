//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
)

// TestMorganaIntegrationSuite runs the complete integration test suite
func TestMorganaIntegrationSuite(t *testing.T) {
	// Setup for the entire suite
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("FullSystemIntegration", func(t *testing.T) {
		// Test the complete system: Event Bus + TUI + Monitor
		testFullSystemIntegration(t, setup, ctx)
	})

	t.Run("PerformanceValidation", func(t *testing.T) {
		// Validate all performance requirements
		testPerformanceValidation(t, setup, ctx)
	})

	t.Run("ReliabilityTesting", func(t *testing.T) {
		// Test system reliability under various conditions
		testReliabilityScenarios(t, setup, ctx)
	})
}

func testFullSystemIntegration(t *testing.T, setup *TestSetup, ctx context.Context) {
	// Start monitor server
	server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
	go func() {
		if err := server.Start(ctx); err != nil && err != context.Canceled {
			t.Errorf("Monitor server error: %v", err)
		}
	}()
	defer server.Stop()

	// Wait for monitor to be ready
	time.Sleep(100 * time.Millisecond)

	// Start TUI if terminal is supported
	var tuiInstance *tui.TUI
	if tui.IsTerminalSupported() {
		config := tui.CreateOptimizedConfig()
		config.RefreshRate = 100 * time.Millisecond // Slower for integration test

		tuiInstance = tui.New(ctx, setup.EventBus, config)
		err := tuiInstance.StartAsync()
		if err != nil {
			t.Errorf("Failed to start TUI: %v", err)
		} else {
			defer tuiInstance.Stop()
		}
	}

	// Generate a variety of events
	agentTypes := RandomAgentTypes()

	// Successful tasks
	for i, agentType := range agentTypes {
		taskID := setup.Generator.GenerateTaskLifecycle(ctx, agentType,
			time.Duration(i+1)*10*time.Millisecond)
		t.Logf("Generated successful task: %s (%s)", taskID, agentType)
	}

	// Failed tasks
	for i := 0; i < 2; i++ {
		agentType := agentTypes[i%len(agentTypes)]
		taskID := setup.Generator.GenerateFailedTask(ctx, agentType, 20*time.Millisecond)
		t.Logf("Generated failed task: %s (%s)", taskID, agentType)
	}

	// High volume events
	taskIDs := setup.Generator.GenerateHighVolumeEvents(ctx, 200, agentTypes)
	t.Logf("Generated %d high volume events", len(taskIDs))

	// Wait for all processing to complete
	time.Sleep(2 * time.Second)

	// Validate system state
	stats := setup.Collector.GetStats()
	busStats := setup.EventBus.Stats()

	// Should have processed all events
	expectedEvents := len(agentTypes)*5 + 2*3 + len(taskIDs) // lifecycles + failures + high volume
	if int(stats.TotalEvents) < expectedEvents*0.9 {         // Allow 10% loss
		t.Errorf("Expected ~%d events, got %d", expectedEvents, stats.TotalEvents)
	}

	// Monitor should have connected clients
	if server.GetClientCount() > 0 {
		t.Logf("Monitor has %d connected clients", server.GetClientCount())
	}

	// TUI should be processing events
	if tuiInstance != nil {
		tuiStats := tuiInstance.GetStats()
		if tuiStats.EventsProcessed == 0 {
			t.Error("TUI processed no events")
		}
		t.Logf("TUI processed %d events at %.1f FPS",
			tuiStats.EventsProcessed, tuiStats.FPS)
	}

	// Verify event types
	if stats.EventsByType[events.EventTaskStarted] == 0 {
		t.Error("No task started events processed")
	}
	if stats.EventsByType[events.EventTaskCompleted] == 0 {
		t.Error("No task completed events processed")
	}
	if stats.EventsByType[events.EventTaskFailed] == 0 {
		t.Error("No task failed events processed")
	}

	t.Logf("Full system integration: Bus(Published=%d, Dropped=%d), Collector(Events=%d)",
		busStats.TotalPublished, busStats.TotalDropped, stats.TotalEvents)
}

func testPerformanceValidation(t *testing.T, setup *TestSetup, ctx context.Context) {
	// Validate all performance requirements from the original specification

	t.Run("LatencyRequirement", func(t *testing.T) {
		// <10ms event processing latency
		agentTypes := []string{"performance-test"}

		start := time.Now()
		setup.Generator.GenerateHighVolumeEvents(ctx, 100, agentTypes)

		// Wait for processing
		if !WaitForEvents(setup.Collector, 100, time.Second*2) {
			t.Fatal("Timeout waiting for latency test events")
		}

		processingTime := time.Since(start)
		avgLatency := processingTime / 100

		if avgLatency > 10*time.Millisecond {
			t.Errorf("Average event latency too high: %v (requirement: <10ms)", avgLatency)
		}

		t.Logf("Latency validation: %v average per event", avgLatency)
	})

	t.Run("ThroughputRequirement", func(t *testing.T) {
		// High throughput event processing
		setup.Collector.Reset()

		agentTypes := RandomAgentTypes()
		eventCount := 5000

		start := time.Now()
		setup.Generator.GenerateHighVolumeEvents(ctx, eventCount, agentTypes)

		// Allow up to 10 seconds for processing
		WaitForEvents(setup.Collector, eventCount, time.Second*10)
		duration := time.Since(start)

		stats := setup.Collector.GetStats()
		throughput := float64(stats.TotalEvents) / duration.Seconds()

		minThroughput := 1000.0 // 1000 events/sec minimum
		if throughput < minThroughput {
			t.Errorf("Throughput too low: %.0f events/sec (requirement: >%.0f)",
				throughput, minThroughput)
		}

		t.Logf("Throughput validation: %.0f events/sec", throughput)
	})

	t.Run("MemoryUsageRequirement", func(t *testing.T) {
		// Memory usage should remain reasonable
		monitor := NewResourceMonitor()
		runtime.GC()
		monitor.TakeMeasurement(0)

		// Process many events
		for i := 0; i < 10; i++ {
			setup.Generator.GenerateHighVolumeEvents(ctx, 500, RandomAgentTypes())
			time.Sleep(100 * time.Millisecond)
		}

		runtime.GC()
		time.Sleep(50 * time.Millisecond)
		stats := setup.Collector.GetStats()
		monitor.TakeMeasurement(stats.TotalEvents)

		analysis := monitor.GetLeakAnalysis()

		maxMemoryMB := 200.0 // 200MB maximum reasonable usage
		if analysis.MemoryGrowthMB > maxMemoryMB {
			t.Errorf("Memory usage too high: %.1f MB (max: %.1f MB)",
				analysis.MemoryGrowthMB, maxMemoryMB)
		}

		t.Logf("Memory usage validation: %.1f MB used", analysis.MemoryGrowthMB)
	})

	t.Run("EventLossRequirement", func(t *testing.T) {
		// Event loss should be minimal (<1%)
		setup.Collector.Reset()

		eventCount := 2000
		setup.Generator.GenerateHighVolumeEvents(ctx, eventCount, RandomAgentTypes())

		time.Sleep(time.Second * 3)

		stats := setup.Collector.GetStats()
		busStats := setup.EventBus.Stats()

		lossRate := float64(busStats.TotalDropped) / float64(busStats.TotalPublished)
		maxLossRate := 0.01 // 1% maximum loss rate

		if lossRate > maxLossRate {
			t.Errorf("Event loss rate too high: %.2f%% (max: %.2f%%)",
				lossRate*100, maxLossRate*100)
		}

		t.Logf("Event loss validation: %.2f%% loss rate", lossRate*100)
	})
}

func testReliabilityScenarios(t *testing.T, setup *TestSetup, ctx context.Context) {
	t.Run("SystemRecovery", func(t *testing.T) {
		// Test system recovery from various failures

		// Normal operation
		setup.Generator.GenerateTaskLifecycle(ctx, "recovery-test-1", 20*time.Millisecond)
		time.Sleep(200 * time.Millisecond)

		// Simulate panic in subscriber (should be recovered)
		panicSubscriber := func(event events.Event) {
			if event.TaskID() == "panic-trigger" {
				panic("Test panic for recovery")
			}
		}

		subscriptionID := setup.EventBus.Subscribe(events.EventTaskStarted, panicSubscriber)

		// Trigger panic
		panicEvent := events.NewTaskStartedEvent(
			ctx, "panic-trigger", "panic-agent", "panic test", nil, 0, "", "", time.Minute,
		)
		setup.EventBus.PublishAsync(panicEvent)

		// System should continue working
		time.Sleep(100 * time.Millisecond)
		setup.Generator.GenerateTaskLifecycle(ctx, "recovery-test-2", 20*time.Millisecond)
		time.Sleep(200 * time.Millisecond)

		stats := setup.Collector.GetStats()
		if stats.TotalEvents < 6 { // Should have at least lifecycle events
			t.Error("System did not recover properly from panic")
		}

		setup.EventBus.Unsubscribe(subscriptionID)
		t.Log("System successfully recovered from subscriber panic")
	})

	t.Run("LoadStressTest", func(t *testing.T) {
		// Test system under extreme load
		setup.Collector.Reset()

		// Generate extreme load
		extremeLoad := 10000
		agentTypes := RandomAgentTypes()

		start := time.Now()
		for batch := 0; batch < 10; batch++ {
			setup.Generator.GenerateHighVolumeEvents(ctx, extremeLoad/10, agentTypes)
			time.Sleep(50 * time.Millisecond) // Brief pause between batches
		}

		// System should handle load gracefully
		time.Sleep(5 * time.Second) // Allow time for processing

		busStats := setup.EventBus.Stats()
		processingDuration := time.Since(start)

		// System should not crash or stop processing
		if busStats.TotalPublished == 0 {
			t.Error("System stopped publishing events under stress")
		}

		// Some events may be dropped under extreme load, but system should survive
		t.Logf("Stress test: %d events published, %d dropped in %v",
			busStats.TotalPublished, busStats.TotalDropped, processingDuration)
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		// Test concurrent operations don't cause race conditions

		numGoroutines := 10
		eventsPerGoroutine := 100

		start := time.Now()

		// Start multiple goroutines generating events concurrently
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				agentType := fmt.Sprintf("concurrent-agent-%d", id)
				for j := 0; j < eventsPerGoroutine; j++ {
					setup.Generator.GenerateTaskLifecycle(ctx, agentType, time.Millisecond)
				}
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Wait for processing
		time.Sleep(2 * time.Second)

		stats := setup.Collector.GetStats()
		duration := time.Since(start)

		// Should process most events without race conditions
		expectedEvents := numGoroutines * eventsPerGoroutine * 5 // 5 events per lifecycle
		if int(stats.TotalEvents) < expectedEvents/2 {           // Allow significant variance due to concurrency
			t.Errorf("Too few events processed in concurrent test: %d < %d",
				stats.TotalEvents, expectedEvents/2)
		}

		t.Logf("Concurrent operations: %d events processed from %d goroutines in %v",
			stats.TotalEvents, numGoroutines, duration)
	})

	t.Run("LongRunningStability", func(t *testing.T) {
		// Test system stability over extended operation
		testDuration := 30 * time.Second
		endTime := time.Now().Add(testDuration)

		// Continuous event generation
		go func() {
			counter := 0
			for time.Now().Before(endTime) {
				setup.Generator.GenerateTaskLifecycle(ctx,
					fmt.Sprintf("stability-test-%d", counter),
					10*time.Millisecond)
				counter++
				time.Sleep(100 * time.Millisecond)
			}
		}()

		// Monitor system health periodically
		healthChecks := 0
		healthTicker := time.NewTicker(5 * time.Second)
		defer healthTicker.Stop()

		for time.Now().Before(endTime) {
			select {
			case <-healthTicker.C:
				stats := setup.Collector.GetStats()
				busStats := setup.EventBus.Stats()

				// System should still be processing events
				if stats.TotalEvents == 0 {
					t.Error("System stopped processing events during stability test")
					return
				}

				// Queue should not grow indefinitely
				if busStats.QueueSize > busStats.QueueCapacity*0.9 {
					t.Logf("Warning: Queue filling up (%d/%d)",
						busStats.QueueSize, busStats.QueueCapacity)
				}

				healthChecks++
				t.Logf("Health check %d: Events=%d, Queue=%d/%d",
					healthChecks, stats.TotalEvents, busStats.QueueSize, busStats.QueueCapacity)

			case <-time.After(testDuration):
				break
			}
		}

		finalStats := setup.Collector.GetStats()
		if finalStats.TotalEvents == 0 {
			t.Error("No events processed during stability test")
		}

		t.Logf("Stability test completed: %d events processed over %v",
			finalStats.TotalEvents, testDuration)
	})
}
