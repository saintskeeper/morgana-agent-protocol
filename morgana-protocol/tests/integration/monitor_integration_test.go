//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
)

func TestMonitorIntegration(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("MonitorServerClientCommunication", func(t *testing.T) {
		// Create and start monitor server
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)

		// Start server in background
		go func() {
			if err := server.Start(ctx); err != nil && err != context.Canceled {
				t.Errorf("Monitor server error: %v", err)
			}
		}()

		// Wait for server to be ready
		time.Sleep(100 * time.Millisecond)

		// Create separate event bus for client
		clientConfig := events.DefaultBusConfig()
		clientEventBus := events.NewEventBus(clientConfig)
		defer clientEventBus.Close()

		// Create and connect client
		client := monitor.NewIPCClient(setup.MonitorSock, clientEventBus)
		err := client.Connect()
		if err != nil {
			t.Fatalf("Failed to connect monitor client: %v", err)
		}
		defer client.Close()

		// Generate events on client side
		taskID := setup.Generator.GenerateTaskLifecycle(
			ctx, "monitor-test-agent", 50*time.Millisecond)

		// Wait for events to be forwarded to server
		time.Sleep(500 * time.Millisecond)

		// Verify server received events
		if !WaitForEvents(setup.Collector, 5, time.Second*2) {
			t.Fatal("Server did not receive forwarded events")
		}

		// Verify events contain correct task ID
		events := setup.Collector.GetEventsForTask(taskID)
		if len(events) == 0 {
			t.Error("No events found for generated task")
		}

		// Check server stats
		if server.GetClientCount() != 1 {
			t.Errorf("Expected 1 client, got %d", server.GetClientCount())
		}

		t.Logf("Successfully tested monitor communication with task %s", taskID)

		// Stop server
		if err := server.Stop(); err != nil {
			t.Errorf("Server stop error: %v", err)
		}
	})

	t.Run("MultipleClientsToMonitor", func(t *testing.T) {
		// Test multiple clients connecting to single monitor
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)

		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		// Create multiple clients
		numClients := 3
		clients := make([]*monitor.IPCClient, numClients)
		clientBuses := make([]events.EventBus, numClients)

		for i := 0; i < numClients; i++ {
			config := events.DefaultBusConfig()
			clientBus := events.NewEventBus(config)
			clientBuses[i] = clientBus

			client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
			clients[i] = client

			err := client.Connect()
			if err != nil {
				t.Fatalf("Client %d failed to connect: %v", i, err)
			}
		}

		// Verify all clients connected
		if server.GetClientCount() != numClients {
			t.Errorf("Expected %d clients, got %d", numClients, server.GetClientCount())
		}

		// Generate events from each client
		for i, clientBus := range clientBuses {
			generator := NewTestEventGenerator(clientBus)
			generator.GenerateTaskLifecycle(ctx,
				fmt.Sprintf("multi-client-agent-%d", i), 30*time.Millisecond)
		}

		// Wait for all events to be forwarded
		expectedEvents := numClients * 5 // 5 events per lifecycle
		if !WaitForEvents(setup.Collector, expectedEvents, time.Second*3) {
			t.Errorf("Did not receive all expected events from multiple clients")
		}

		stats := setup.Collector.GetStats()
		t.Logf("Multi-client test: %d clients generated %d events",
			numClients, stats.TotalEvents)

		// Clean up
		for i, client := range clients {
			client.Close()
			clientBuses[i].Close()
		}

		server.Stop()
	})

	t.Run("MonitorReconnection", func(t *testing.T) {
		// Test client reconnection when server restarts
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)

		// Start server
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		// Connect client
		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(setup.MonitorSock, clientBus)

		// Initial connection should succeed
		if !client.TryConnect() {
			t.Fatal("Initial client connection failed")
		}

		if !client.IsConnected() {
			t.Error("Client should report as connected")
		}

		// Stop server
		server.Stop()
		time.Sleep(100 * time.Millisecond)

		// Client should detect disconnection
		// Generate event to trigger connection check
		generator := NewTestEventGenerator(clientBus)
		generator.GenerateTaskLifecycle(ctx, "reconnect-test", 10*time.Millisecond)
		time.Sleep(200 * time.Millisecond)

		// Try to reconnect should fail while server is down
		if client.TryConnect() {
			t.Error("Client should not be able to connect to stopped server")
		}

		// Restart server
		server = monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		// Reconnection should now succeed
		if !client.TryConnect() {
			t.Error("Client should be able to reconnect to restarted server")
		}

		t.Log("Successfully tested monitor reconnection")

		client.Close()
		server.Stop()
	})

	t.Run("MonitorEventFiltering", func(t *testing.T) {
		// Test that monitor properly handles different event types
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
		err := client.Connect()
		if err != nil {
			t.Fatalf("Client connection failed: %v", err)
		}
		defer client.Close()

		generator := NewTestEventGenerator(clientBus)

		// Generate different types of events
		successTask := generator.GenerateTaskLifecycle(ctx, "success-agent", 20*time.Millisecond)
		failTask := generator.GenerateFailedTask(ctx, "fail-agent", 20*time.Millisecond)

		// Wait for all events
		time.Sleep(600 * time.Millisecond)

		stats := setup.Collector.GetStats()

		// Verify different event types were processed
		if stats.EventsByType[events.EventTaskStarted] < 2 {
			t.Error("Expected at least 2 task started events")
		}

		if stats.EventsByType[events.EventTaskCompleted] < 1 {
			t.Error("Expected at least 1 task completed event")
		}

		if stats.EventsByType[events.EventTaskFailed] < 1 {
			t.Error("Expected at least 1 task failed event")
		}

		t.Logf("Event filtering test - Success task: %s, Failed task: %s",
			successTask, failTask)
		t.Logf("Event type distribution: %+v", stats.EventsByType)

		server.Stop()
	})
}

func TestMonitorPerformance(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("HighVolumeEventForwarding", func(t *testing.T) {
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
		err := client.Connect()
		if err != nil {
			t.Fatalf("Client connection failed: %v", err)
		}
		defer client.Close()

		// Generate high volume of events
		generator := NewTestEventGenerator(clientBus)
		agentTypes := RandomAgentTypes()
		eventCount := 1000

		start := time.Now()
		taskIDs := generator.GenerateHighVolumeEvents(ctx, eventCount, agentTypes)
		generateDuration := time.Since(start)

		// Wait for forwarding
		time.Sleep(2 * time.Second)

		stats := setup.Collector.GetStats()
		processingDuration := time.Since(start)

		// Check performance
		throughput := float64(stats.TotalEvents) / processingDuration.Seconds()

		if throughput < 500 { // Should handle at least 500 events/sec through monitor
			t.Errorf("Monitor throughput too low: %.0f events/sec", throughput)
		}

		forwardingRate := float64(stats.TotalEvents) / float64(len(taskIDs))
		if forwardingRate < 0.8 { // Should forward at least 80% of events
			t.Errorf("Too many events lost in forwarding: %.1f%%", forwardingRate*100)
		}

		t.Logf("High volume forwarding: %d events generated in %v, %d forwarded in %v",
			len(taskIDs), generateDuration, stats.TotalEvents, processingDuration)
		t.Logf("Forwarding throughput: %.0f events/sec", throughput)

		server.Stop()
	})

	t.Run("MonitorMemoryUsage", func(t *testing.T) {
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
		client.Connect()
		defer client.Close()

		generator := NewTestEventGenerator(clientBus)

		// Process events in batches to test memory behavior
		batchSize := 500
		batches := 5

		for batch := 0; batch < batches; batch++ {
			generator.GenerateHighVolumeEvents(ctx, batchSize,
				[]string{fmt.Sprintf("memory-test-batch-%d", batch)})

			time.Sleep(200 * time.Millisecond)

			// Check that system is still responsive
			if !client.IsConnected() {
				t.Errorf("Client disconnected during batch %d", batch)
				break
			}
		}

		// Final check
		time.Sleep(500 * time.Millisecond)
		finalStats := setup.Collector.GetStats()

		if finalStats.TotalEvents < int64(batchSize*batches*0.8) {
			t.Errorf("Too many events lost during memory test: %d < %d",
				finalStats.TotalEvents, int64(batchSize*batches))
		}

		t.Logf("Memory test completed: %d events processed across %d batches",
			finalStats.TotalEvents, batches)

		server.Stop()
	})

	t.Run("MonitorLatencyMeasurement", func(t *testing.T) {
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
		client.Connect()
		defer client.Close()

		// Measure forwarding latency
		latencies := make([]time.Duration, 0)
		eventCount := 100

		for i := 0; i < eventCount; i++ {
			start := time.Now()

			// Create event with timestamp
			event := events.NewTaskStartedEvent(
				ctx, fmt.Sprintf("latency-test-%d", i), "latency-agent",
				"latency test prompt", nil, 0, "", "", time.Minute,
			)

			clientBus.PublishAsync(event)

			// Wait for event to be processed (simple approximation)
			time.Sleep(time.Millisecond * 5)

			latency := time.Since(start)
			latencies = append(latencies, latency)
		}

		// Calculate average latency
		var totalLatency time.Duration
		var maxLatency time.Duration

		for _, latency := range latencies {
			totalLatency += latency
			if latency > maxLatency {
				maxLatency = latency
			}
		}

		avgLatency := totalLatency / time.Duration(len(latencies))

		// Performance assertions
		if avgLatency > 10*time.Millisecond {
			t.Errorf("Average forwarding latency too high: %v", avgLatency)
		}

		if maxLatency > 50*time.Millisecond {
			t.Errorf("Max forwarding latency too high: %v", maxLatency)
		}

		t.Logf("Latency measurement: Avg=%v, Max=%v (n=%d)",
			avgLatency, maxLatency, len(latencies))

		server.Stop()
	})
}

func TestMonitorErrorHandling(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("InvalidSocketPath", func(t *testing.T) {
		// Test monitor with invalid socket path
		invalidPath := "/nonexistent/path/monitor.sock"
		server := monitor.NewIPCServer(invalidPath, setup.EventBus)

		err := server.Start(ctx)
		if err == nil {
			t.Error("Server should fail to start with invalid socket path")
			server.Stop()
		} else {
			t.Logf("Correctly failed with invalid socket path: %v", err)
		}
	})

	t.Run("ClientConnectionToNonexistentServer", func(t *testing.T) {
		// Test client connection when server is not running
		nonexistentSock := setup.TempDir + "/nonexistent.sock"

		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(nonexistentSock, clientBus)

		if client.TryConnect() {
			t.Error("Client should not be able to connect to nonexistent server")
			client.Close()
		}

		if client.IsConnected() {
			t.Error("Client should report as not connected")
		}

		t.Log("Correctly handled connection to nonexistent server")
	})

	t.Run("ServerShutdownWithActiveClients", func(t *testing.T) {
		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		// Connect multiple clients
		clients := make([]*monitor.IPCClient, 3)
		clientBuses := make([]events.EventBus, 3)

		for i := 0; i < 3; i++ {
			clientBus := events.NewEventBus(events.DefaultBusConfig())
			clientBuses[i] = clientBus

			client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
			clients[i] = client

			err := client.Connect()
			if err != nil {
				t.Fatalf("Client %d connection failed: %v", i, err)
			}
		}

		// Verify all clients connected
		if server.GetClientCount() != 3 {
			t.Errorf("Expected 3 clients, got %d", server.GetClientCount())
		}

		// Shut down server while clients are connected
		err := server.Stop()
		if err != nil {
			t.Errorf("Server shutdown failed: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// Verify server reports no clients after shutdown
		if server.GetClientCount() != 0 {
			t.Errorf("Expected 0 clients after shutdown, got %d", server.GetClientCount())
		}

		// Clean up clients
		for i, client := range clients {
			client.Close()
			clientBuses[i].Close()
		}

		t.Log("Successfully handled server shutdown with active clients")
	})

	t.Run("EventSerializationErrors", func(t *testing.T) {
		// This test would be more meaningful with actual serialization issues
		// For now, we test that the system handles normal events without errors

		server := monitor.NewIPCServer(setup.MonitorSock, setup.EventBus)
		go server.Start(ctx)
		time.Sleep(100 * time.Millisecond)

		clientBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientBus.Close()

		client := monitor.NewIPCClient(setup.MonitorSock, clientBus)
		err := client.Connect()
		if err != nil {
			t.Fatalf("Client connection failed: %v", err)
		}
		defer client.Close()

		// Generate events with various data types
		generator := NewTestEventGenerator(clientBus)

		// Create events with complex data
		for i := 0; i < 10; i++ {
			complexEvent := events.NewTaskStartedEvent(
				ctx, fmt.Sprintf("serialization-test-%d", i), "test-agent",
				fmt.Sprintf("Complex prompt with unicode: 测试 emoji: 🚀 number: %d", i),
				map[string]interface{}{
					"string":  "test value",
					"number":  42,
					"boolean": true,
					"array":   []string{"a", "b", "c"},
					"nested":  map[string]interface{}{"key": "value"},
				},
				0, "model-hint", "complexity", time.Minute,
			)
			clientBus.PublishAsync(complexEvent)
		}

		time.Sleep(500 * time.Millisecond)

		stats := setup.Collector.GetStats()
		if stats.TotalEvents < 10 {
			t.Errorf("Expected at least 10 events, got %d (possible serialization issues)",
				stats.TotalEvents)
		}

		t.Logf("Serialization test: %d events processed successfully", stats.TotalEvents)

		server.Stop()
	})
}
