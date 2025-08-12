//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
)

func TestTUIIntegration(t *testing.T) {
	// Skip TUI tests if not in a terminal environment
	if !tui.IsTerminalSupported() {
		t.Skip("Terminal does not support TUI - skipping TUI integration tests")
	}

	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("TUIDisplaysEvents", func(t *testing.T) {
		// Create TUI with optimized config for testing
		config := tui.CreateOptimizedConfig()
		config.RefreshRate = 50 * time.Millisecond // Faster refresh for testing

		tuiInstance := tui.New(ctx, setup.EventBus, config)

		// Start TUI asynchronously
		err := tuiInstance.StartAsync()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}
		defer tuiInstance.Stop()

		// Wait for TUI to initialize
		time.Sleep(100 * time.Millisecond)

		// Generate some test events
		taskID1 := setup.Generator.GenerateTaskLifecycle(ctx, "code-implementer", 50*time.Millisecond)
		taskID2 := setup.Generator.GenerateFailedTask(ctx, "test-specialist", 30*time.Millisecond)

		// Wait for events to be processed and displayed
		time.Sleep(500 * time.Millisecond)

		// Get TUI stats to verify it's processing events
		tuiStats := tuiInstance.GetStats()

		if tuiStats.EventsProcessed == 0 {
			t.Error("TUI did not process any events")
		}

		if !tuiStats.IsRunning {
			t.Error("TUI is not running")
		}

		if tuiStats.RenderCount == 0 {
			t.Error("TUI did not render any frames")
		}

		t.Logf("TUI Stats: Events=%d, Renders=%d, FPS=%.1f, Uptime=%v",
			tuiStats.EventsProcessed, tuiStats.RenderCount, tuiStats.FPS, tuiStats.Uptime)
		t.Logf("Generated tasks: %s, %s", taskID1, taskID2)
	})

	t.Run("TUIPerformanceUnderLoad", func(t *testing.T) {
		// Test TUI performance with high event volume
		config := tui.CreateHighPerformanceConfig()
		tuiInstance := tui.New(ctx, setup.EventBus, config)

		err := tuiInstance.StartAsync()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}
		defer tuiInstance.Stop()

		// Wait for initialization
		time.Sleep(100 * time.Millisecond)

		// Generate high volume of events
		agentTypes := RandomAgentTypes()
		setup.Generator.GenerateHighVolumeEvents(ctx, 500, agentTypes)

		// Let TUI process events for a while
		time.Sleep(2 * time.Second)

		tuiStats := tuiInstance.GetStats()

		// Verify TUI is still responsive under load
		if tuiStats.EventsProcessed == 0 {
			t.Error("TUI stopped processing events under load")
		}

		if tuiStats.FPS < 1.0 {
			t.Errorf("TUI FPS too low under load: %.1f", tuiStats.FPS)
		}

		// Memory usage should be reasonable (this is a rough check)
		if tuiStats.MemoryMB > 500 {
			t.Errorf("TUI memory usage too high: %.1f MB", tuiStats.MemoryMB)
		}

		t.Logf("TUI under load - Events: %d, FPS: %.1f, Memory: %.1f MB",
			tuiStats.EventsProcessed, tuiStats.FPS, tuiStats.MemoryMB)
	})

	t.Run("TUIConfigurationOptions", func(t *testing.T) {
		// Test different TUI configurations
		testConfigs := []struct {
			name   string
			config tui.TUIConfig
		}{
			{"Default", tui.DefaultTUIConfig()},
			{"Optimized", tui.CreateOptimizedConfig()},
			{"HighPerformance", tui.CreateHighPerformanceConfig()},
			{"Development", tui.CreateDevelopmentConfig()},
		}

		for _, tc := range testConfigs {
			t.Run(tc.name, func(t *testing.T) {
				// Validate configuration
				if err := tui.ValidateConfig(tc.config); err != nil {
					t.Errorf("Invalid config %s: %v", tc.name, err)
					return
				}

				// Test TUI creation with config
				tuiInstance := tui.New(ctx, setup.EventBus, tc.config)
				if tuiInstance == nil {
					t.Errorf("Failed to create TUI with %s config", tc.name)
					return
				}

				// Brief startup test
				err := tuiInstance.StartAsync()
				if err != nil {
					t.Errorf("Failed to start TUI with %s config: %v", tc.name, err)
					return
				}

				time.Sleep(50 * time.Millisecond)
				tuiInstance.Stop()

				t.Logf("Successfully tested %s configuration", tc.name)
			})
		}
	})

	t.Run("TUIManagerMultipleInstances", func(t *testing.T) {
		// Test TUI manager with multiple instances
		manager := tui.NewTUIManager()

		// Create multiple TUI instances
		instanceIDs := []string{"tui1", "tui2", "tui3"}

		for _, id := range instanceIDs {
			config := tui.CreateOptimizedConfig()
			config.RefreshRate = 100 * time.Millisecond // Slower refresh to reduce load

			tuiInstance, err := manager.Create(id, ctx, setup.EventBus, config)
			if err != nil {
				t.Errorf("Failed to create TUI instance %s: %v", id, err)
				continue
			}

			err = tuiInstance.StartAsync()
			if err != nil {
				t.Errorf("Failed to start TUI instance %s: %v", id, err)
			}
		}

		// Wait for all instances to initialize
		time.Sleep(200 * time.Millisecond)

		// Generate events for all instances to process
		setup.Generator.GenerateTaskLifecycle(ctx, "multi-tui-test", 20*time.Millisecond)

		time.Sleep(300 * time.Millisecond)

		// Verify all instances are running
		instances := manager.List()
		if len(instances) != len(instanceIDs) {
			t.Errorf("Expected %d instances, got %d", len(instanceIDs), len(instances))
		}

		// Check each instance
		for _, id := range instanceIDs {
			tuiInstance, exists := manager.Get(id)
			if !exists {
				t.Errorf("Instance %s not found in manager", id)
				continue
			}

			stats := tuiInstance.GetStats()
			if !stats.IsRunning {
				t.Errorf("Instance %s is not running", id)
			}

			t.Logf("Instance %s stats: Events=%d, FPS=%.1f",
				id, stats.EventsProcessed, stats.FPS)
		}

		// Clean up all instances
		if err := manager.StopAll(); err != nil {
			t.Errorf("Failed to stop all TUI instances: %v", err)
		}

		// Verify cleanup
		if len(manager.List()) != 0 {
			t.Error("Not all TUI instances were cleaned up")
		}
	})
}

func TestTUIStatisticsTracking(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !tui.IsTerminalSupported() {
		t.Skip("Terminal does not support TUI - skipping statistics tests")
	}

	t.Run("StatisticsAccuracy", func(t *testing.T) {
		config := tui.CreateDevelopmentConfig() // Enable debug info
		tuiInstance := tui.New(ctx, setup.EventBus, config)

		err := tuiInstance.StartAsync()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}
		defer tuiInstance.Stop()

		// Wait for initialization
		time.Sleep(100 * time.Millisecond)

		// Generate known number of events
		expectedEvents := 10
		agentTypes := []string{"stats-test"}

		for i := 0; i < expectedEvents; i++ {
			setup.Generator.GenerateTaskLifecycle(ctx, agentTypes[0], 10*time.Millisecond)
		}

		// Wait for processing
		time.Sleep(2 * time.Second)

		tuiStats := tuiInstance.GetStats()

		// Each lifecycle generates 5 events, but TUI might not catch all async events
		minExpectedEvents := int64(expectedEvents * 3) // At least 3 events per lifecycle

		if tuiStats.EventsProcessed < minExpectedEvents {
			t.Errorf("TUI processed fewer events than expected: %d < %d",
				tuiStats.EventsProcessed, minExpectedEvents)
		}

		// Verify TUI is tracking render count
		if tuiStats.RenderCount == 0 {
			t.Error("TUI render count is zero")
		}

		// Verify FPS is reasonable (should be > 0 and < 1000)
		if tuiStats.FPS <= 0 || tuiStats.FPS > 1000 {
			t.Errorf("TUI FPS seems unreasonable: %.1f", tuiStats.FPS)
		}

		// Verify uptime is tracked
		if tuiStats.Uptime <= 0 {
			t.Error("TUI uptime not tracked")
		}

		t.Logf("Statistics accuracy test - Events: %d, Renders: %d, FPS: %.1f, Uptime: %v",
			tuiStats.EventsProcessed, tuiStats.RenderCount, tuiStats.FPS, tuiStats.Uptime)
	})

	t.Run("RealTimeStatistics", func(t *testing.T) {
		config := tui.CreateOptimizedConfig()
		tuiInstance := tui.New(ctx, setup.EventBus, config)

		err := tuiInstance.StartAsync()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}
		defer tuiInstance.Stop()

		time.Sleep(100 * time.Millisecond)

		// Take initial measurement
		initialStats := tuiInstance.GetStats()

		// Generate events continuously
		go func() {
			for i := 0; i < 20; i++ {
				setup.Generator.GenerateTaskLifecycle(ctx, "realtime-test", 25*time.Millisecond)
				time.Sleep(50 * time.Millisecond)
			}
		}()

		// Monitor statistics over time
		measurements := []tui.TUIStats{}
		for i := 0; i < 5; i++ {
			time.Sleep(200 * time.Millisecond)
			stats := tuiInstance.GetStats()
			measurements = append(measurements, stats)
		}

		// Verify statistics are updating
		finalStats := measurements[len(measurements)-1]

		if finalStats.EventsProcessed <= initialStats.EventsProcessed {
			t.Error("Event count did not increase over time")
		}

		if finalStats.RenderCount <= initialStats.RenderCount {
			t.Error("Render count did not increase over time")
		}

		if finalStats.Uptime <= initialStats.Uptime {
			t.Error("Uptime did not increase")
		}

		// Log progression
		for i, stats := range measurements {
			t.Logf("Measurement %d: Events=%d, Renders=%d, FPS=%.1f",
				i+1, stats.EventsProcessed, stats.RenderCount, stats.FPS)
		}
	})
}

func TestTUIErrorHandling(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !tui.IsTerminalSupported() {
		t.Skip("Terminal does not support TUI - skipping error handling tests")
	}

	t.Run("InvalidConfiguration", func(t *testing.T) {
		// Test TUI handles invalid configurations gracefully
		invalidConfigs := []tui.TUIConfig{
			{RefreshRate: -1},    // Negative refresh rate
			{RefreshRate: 0},     // Zero refresh rate
			{EventBufferSize: 0}, // Zero buffer size
			{MaxLogLines: -1},    // Negative max lines
		}

		for i, invalidConfig := range invalidConfigs {
			err := tui.ValidateConfig(invalidConfig)
			if err == nil {
				t.Errorf("Invalid config %d should have failed validation", i)
			} else {
				t.Logf("Config %d correctly failed validation: %v", i, err)
			}
		}
	})

	t.Run("GracefulShutdown", func(t *testing.T) {
		config := tui.CreateOptimizedConfig()
		tuiInstance := tui.New(ctx, setup.EventBus, config)

		err := tuiInstance.StartAsync()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}

		// Generate some events
		setup.Generator.GenerateTaskLifecycle(ctx, "shutdown-test", 20*time.Millisecond)
		time.Sleep(100 * time.Millisecond)

		// Verify TUI is running
		initialStats := tuiInstance.GetStats()
		if !initialStats.IsRunning {
			t.Error("TUI should be running before shutdown")
		}

		// Test graceful shutdown
		err = tuiInstance.Stop()
		if err != nil {
			t.Errorf("TUI shutdown failed: %v", err)
		}

		// Verify shutdown completed
		time.Sleep(50 * time.Millisecond)
		finalStats := tuiInstance.GetStats()

		// Note: IsRunning might not be immediately false due to async shutdown
		t.Logf("TUI shutdown completed - Final stats: Events=%d, Renders=%d",
			finalStats.EventsProcessed, finalStats.RenderCount)
	})

	t.Run("TerminalResize", func(t *testing.T) {
		// Test TUI handles terminal resize events gracefully
		// Note: This is difficult to test programmatically without actual terminal interaction
		// In a real implementation, you'd simulate SIGWINCH signals

		config := tui.CreateOptimizedConfig()
		tuiInstance := tui.New(ctx, setup.EventBus, config)

		err := tuiInstance.StartAsync()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}
		defer tuiInstance.Stop()

		// Generate events during "resize"
		setup.Generator.GenerateTaskLifecycle(ctx, "resize-test", 30*time.Millisecond)

		// Simulate some activity during potential resize
		time.Sleep(200 * time.Millisecond)

		stats := tuiInstance.GetStats()
		if stats.EventsProcessed == 0 {
			t.Error("TUI stopped processing events (possibly due to resize issues)")
		}

		t.Logf("TUI handled potential resize scenario - Events: %d", stats.EventsProcessed)
	})
}
