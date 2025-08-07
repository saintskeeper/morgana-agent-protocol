//go:build integration
// +build integration

package monitor_test

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
)

func TestIPCServerClientIntegration(t *testing.T) {
	// Create temporary socket path
	socketPath := "/tmp/morgana_test.sock"
	defer os.Remove(socketPath)

	// Create event buses for server and client
	serverEventBus := events.NewEventBus(events.DefaultBusConfig())
	defer serverEventBus.Close()

	clientEventBus := events.NewEventBus(events.DefaultBusConfig())
	defer clientEventBus.Close()

	// Create and start IPC server
	server := monitor.NewIPCServer(socketPath, serverEventBus)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in goroutine
	go func() {
		if err := server.Start(ctx); err != nil {
			t.Errorf("Server start error: %v", err)
		}
	}()

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Subscribe to server events to verify forwarding
	var receivedEvent events.Event
	serverEventBus.SubscribeAll(func(event events.Event) {
		receivedEvent = event
	})

	// Create IPC client and connect
	client := monitor.NewIPCClient(socketPath, clientEventBus)
	if err := client.Connect(); err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer client.Close()

	// Publish an event on the client bus
	testEvent := events.NewTaskStartedEvent(
		context.Background(),
		"test-task-123",
		"test-agent",
		"test prompt",
		map[string]interface{}{"test": "value"},
		0,
		"test-model",
		"simple",
		time.Minute,
	)

	clientEventBus.Publish(testEvent)

	// Wait for event to be forwarded
	time.Sleep(100 * time.Millisecond)

	// Verify event was received by server
	if receivedEvent == nil {
		t.Fatal("No event was received by server")
	}

	if receivedEvent.Type() != events.EventTaskStarted {
		t.Errorf("Expected event type %s, got %s", events.EventTaskStarted, receivedEvent.Type())
	}

	if receivedEvent.TaskID() != "test-task-123" {
		t.Errorf("Expected task ID 'test-task-123', got '%s'", receivedEvent.TaskID())
	}

	// Stop server
	if err := server.Stop(); err != nil {
		t.Errorf("Server stop error: %v", err)
	}
}

func TestIPCMessageSerialization(t *testing.T) {
	// Test event serialization/deserialization
	testEvent := events.NewTaskProgressEvent(
		context.Background(),
		"task-456",
		"test-agent",
		"validation",
		"Testing progress",
		0.75,
		time.Second*30,
	)

	// Create IPC message
	msg := monitor.IPCMessage{
		Type:      string(testEvent.Type()),
		TaskID:    testEvent.TaskID(),
		Timestamp: testEvent.Timestamp(),
		Data:      testEvent,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Deserialize from JSON
	var deserializedMsg monitor.IPCMessage
	if err := json.Unmarshal(jsonData, &deserializedMsg); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify message fields
	if deserializedMsg.Type != string(events.EventTaskProgress) {
		t.Errorf("Expected type %s, got %s", events.EventTaskProgress, deserializedMsg.Type)
	}

	if deserializedMsg.TaskID != "task-456" {
		t.Errorf("Expected task ID 'task-456', got '%s'", deserializedMsg.TaskID)
	}

	// Verify data can be accessed
	if deserializedMsg.Data == nil {
		t.Error("Message data is nil")
	}
}

func TestIPCClientReconnection(t *testing.T) {
	socketPath := "/tmp/morgana_reconnect_test.sock"
	defer os.Remove(socketPath)

	clientEventBus := events.NewEventBus(events.DefaultBusConfig())
	defer clientEventBus.Close()

	client := monitor.NewIPCClient(socketPath, clientEventBus)

	// Test connection when server is not running
	if client.TryConnect() {
		t.Error("Client should not be able to connect when server is not running")
	}

	if client.IsConnected() {
		t.Error("Client should report as not connected")
	}

	// Start server
	serverEventBus := events.NewEventBus(events.DefaultBusConfig())
	defer serverEventBus.Close()

	server := monitor.NewIPCServer(socketPath, serverEventBus)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Start(ctx)
	time.Sleep(100 * time.Millisecond) // Wait for server to be ready

	// Now connection should succeed
	if !client.TryConnect() {
		t.Error("Client should be able to connect when server is running")
	}

	if !client.IsConnected() {
		t.Error("Client should report as connected")
	}

	// Clean up
	client.Close()
	server.Stop()
}

func TestMultipleClients(t *testing.T) {
	socketPath := "/tmp/morgana_multi_client_test.sock"
	defer os.Remove(socketPath)

	// Start server
	serverEventBus := events.NewEventBus(events.DefaultBusConfig())
	defer serverEventBus.Close()

	server := monitor.NewIPCServer(socketPath, serverEventBus)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Track received events
	var receivedEvents []events.Event
	serverEventBus.SubscribeAll(func(event events.Event) {
		receivedEvents = append(receivedEvents, event)
	})

	// Create multiple clients
	numClients := 3
	clients := make([]*monitor.IPCClient, numClients)

	for i := 0; i < numClients; i++ {
		clientEventBus := events.NewEventBus(events.DefaultBusConfig())
		defer clientEventBus.Close()

		client := monitor.NewIPCClient(socketPath, clientEventBus)
		clients[i] = client

		if err := client.Connect(); err != nil {
			t.Fatalf("Client %d failed to connect: %v", i, err)
		}

		// Publish test event from each client
		testEvent := events.NewTaskCompletedEvent(
			context.Background(),
			"task-"+string(rune('A'+i)),
			"test-agent",
			"Test output",
			time.Second,
			"test-model",
		)

		clientEventBus.Publish(testEvent)
	}

	// Wait for events to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify all events were received
	if len(receivedEvents) != numClients {
		t.Errorf("Expected %d events, got %d", numClients, len(receivedEvents))
	}

	// Check server client count
	if server.GetClientCount() != numClients {
		t.Errorf("Expected %d connected clients, got %d", numClients, server.GetClientCount())
	}

	// Clean up clients
	for i, client := range clients {
		if err := client.Close(); err != nil {
			t.Errorf("Client %d close error: %v", i, err)
		}
	}

	// Wait for disconnections to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify client count decreased
	if server.GetClientCount() != 0 {
		t.Errorf("Expected 0 connected clients after cleanup, got %d", server.GetClientCount())
	}

	server.Stop()
}

// Helper function to check if socket exists
func socketExists(path string) bool {
	_, err := net.Dial("unix", path)
	return err == nil
}
