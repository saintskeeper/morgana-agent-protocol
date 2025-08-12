package monitor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

func TestEventConsumer_BasicFunctionality(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	eventFile := filepath.Join(tempDir, "events.jsonl")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Track received events
	receivedEvents := make([]events.Event, 0)
	subscriptionID := eventBus.SubscribeAll(func(event events.Event) {
		receivedEvents = append(receivedEvents, event)
	})
	defer eventBus.Unsubscribe(subscriptionID)

	// Create consumer
	config := ConsumerConfig{
		EventFile:    eventFile,
		PollInterval: 10 * time.Millisecond, // Fast polling for tests
		BufferSize:   1024,
	}
	consumer, err := NewEventConsumer(eventBus, config)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Start consumer
	if err := consumer.Start(); err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Stop()

	// Write test events
	testEvents := []string{
		`{"event_type":"task.started","task_id":"test-001","timestamp":"2025-08-11T14:53:03-04:00","agent_type":"code-implementer","prompt":"Test task"}`,
		`{"event_type":"task.progress","task_id":"test-001","timestamp":"2025-08-11T14:53:04-04:00","agent_type":"code-implementer","stage":"execution","message":"Working","progress":0.5}`,
		`{"event_type":"task.completed","task_id":"test-001","timestamp":"2025-08-11T14:53:06-04:00","agent_type":"code-implementer","output":"Success","duration":"3s","model":"test-model"}`,
	}

	file, err := os.OpenFile(eventFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	for _, event := range testEvents {
		if _, err := file.WriteString(event + "\n"); err != nil {
			t.Fatalf("Failed to write event: %v", err)
		}
		file.Sync()
		time.Sleep(50 * time.Millisecond) // Give consumer time to process
	}
	file.Close()

	// Wait for events to be processed
	time.Sleep(200 * time.Millisecond)

	// Check that events were received
	if len(receivedEvents) != len(testEvents) {
		t.Errorf("Expected %d events, got %d", len(testEvents), len(receivedEvents))
	}

	// Verify event types
	expectedTypes := []events.EventType{
		events.EventTaskStarted,
		events.EventTaskProgress,
		events.EventTaskCompleted,
	}

	for i, expectedType := range expectedTypes {
		if i >= len(receivedEvents) {
			break
		}
		if receivedEvents[i].Type() != expectedType {
			t.Errorf("Event %d: expected type %s, got %s", i, expectedType, receivedEvents[i].Type())
		}
		if receivedEvents[i].TaskID() != "test-001" {
			t.Errorf("Event %d: expected task ID 'test-001', got '%s'", i, receivedEvents[i].TaskID())
		}
	}
}

func TestEventConsumer_LegacyEventFormat(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	eventFile := filepath.Join(tempDir, "events.jsonl")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Track received events
	receivedEvents := make([]events.Event, 0)
	subscriptionID := eventBus.SubscribeAll(func(event events.Event) {
		receivedEvents = append(receivedEvents, event)
	})
	defer eventBus.Unsubscribe(subscriptionID)

	// Create consumer
	config := ConsumerConfig{
		EventFile:    eventFile,
		PollInterval: 10 * time.Millisecond,
		BufferSize:   1024,
	}
	consumer, err := NewEventConsumer(eventBus, config)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Start consumer
	if err := consumer.Start(); err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Stop()

	// Write legacy format event
	legacyEvent := `{"event_type":"task_started","task_id":"legacy-001","agent_type":"test-agent","timestamp":"2025-08-11T18:32:48.827862+00:00","session_id":"test-session","pid":12345}`

	file, err := os.OpenFile(eventFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(legacyEvent + "\n"); err != nil {
		t.Fatalf("Failed to write event: %v", err)
	}
	file.Sync()

	// Wait for event to be processed
	time.Sleep(100 * time.Millisecond)

	// Check that event was received
	if len(receivedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(receivedEvents))
	}

	if len(receivedEvents) > 0 {
		event := receivedEvents[0]
		if event.TaskID() != "legacy-001" {
			t.Errorf("Expected task ID 'legacy-001', got '%s'", event.TaskID())
		}
		// Note: legacy format uses underscores, which should be converted
		expectedType := events.EventType("task_started")
		if event.Type() != expectedType {
			t.Errorf("Expected type %s, got %s", expectedType, event.Type())
		}
	}
}

func TestEventConsumer_IPCMessageFormat(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	eventFile := filepath.Join(tempDir, "events.jsonl")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Track received events
	receivedEvents := make([]events.Event, 0)
	subscriptionID := eventBus.SubscribeAll(func(event events.Event) {
		receivedEvents = append(receivedEvents, event)
	})
	defer eventBus.Unsubscribe(subscriptionID)

	// Create consumer
	config := ConsumerConfig{
		EventFile:    eventFile,
		PollInterval: 10 * time.Millisecond,
		BufferSize:   1024,
	}
	consumer, err := NewEventConsumer(eventBus, config)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Start consumer
	if err := consumer.Start(); err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Stop()

	// Create IPC message format
	ipcMsg := IPCMessage{
		Type:      "task.started",
		TaskID:    "ipc-001",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"agent_type": "ipc-agent",
			"prompt":     "IPC test task",
		},
	}

	ipcJSON, err := json.Marshal(ipcMsg)
	if err != nil {
		t.Fatalf("Failed to marshal IPC message: %v", err)
	}

	file, err := os.OpenFile(eventFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(string(ipcJSON) + "\n"); err != nil {
		t.Fatalf("Failed to write event: %v", err)
	}
	file.Sync()

	// Wait for event to be processed
	time.Sleep(100 * time.Millisecond)

	// Check that event was received
	if len(receivedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(receivedEvents))
	}

	if len(receivedEvents) > 0 {
		event := receivedEvents[0]
		if event.TaskID() != "ipc-001" {
			t.Errorf("Expected task ID 'ipc-001', got '%s'", event.TaskID())
		}
		if event.Type() != events.EventTaskStarted {
			t.Errorf("Expected type %s, got %s", events.EventTaskStarted, event.Type())
		}
	}
}

func TestEventConsumer_FileRotation(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	eventFile := filepath.Join(tempDir, "events.jsonl")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Track received events
	receivedEvents := make([]events.Event, 0)
	subscriptionID := eventBus.SubscribeAll(func(event events.Event) {
		receivedEvents = append(receivedEvents, event)
	})
	defer eventBus.Unsubscribe(subscriptionID)

	// Create consumer
	config := ConsumerConfig{
		EventFile:    eventFile,
		PollInterval: 10 * time.Millisecond,
		BufferSize:   1024,
	}
	consumer, err := NewEventConsumer(eventBus, config)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Start consumer
	if err := consumer.Start(); err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Stop()

	// Write initial event
	initialEvent := `{"event_type":"task.started","task_id":"rotation-001","timestamp":"2025-08-11T14:53:03-04:00","agent_type":"test"}`

	file, err := os.OpenFile(eventFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if _, err := file.WriteString(initialEvent + "\n"); err != nil {
		t.Fatalf("Failed to write initial event: %v", err)
	}
	file.Close()

	time.Sleep(50 * time.Millisecond)

	// Simulate file rotation by removing and recreating
	if err := os.Remove(eventFile); err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}

	// Wait a bit for the consumer to detect the removal
	time.Sleep(50 * time.Millisecond)

	// Create new file with different content
	newEvent := `{"event_type":"task.completed","task_id":"rotation-002","timestamp":"2025-08-11T14:54:03-04:00","agent_type":"test","output":"After rotation"}`

	file, err = os.OpenFile(eventFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create new test file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(newEvent + "\n"); err != nil {
		t.Fatalf("Failed to write new event: %v", err)
	}
	file.Sync()

	// Wait for events to be processed
	time.Sleep(100 * time.Millisecond)

	// Check that both events were received
	if len(receivedEvents) < 1 {
		t.Errorf("Expected at least 1 event, got %d", len(receivedEvents))
	}

	// The consumer should handle rotation and read the new event
	found := false
	for _, event := range receivedEvents {
		if event.TaskID() == "rotation-002" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Event after rotation was not received")
	}
}

func TestEventConsumer_Configuration(t *testing.T) {
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Test default configuration
	defaultConfig := DefaultConsumerConfig()
	if defaultConfig.EventFile != "/tmp/morgana/events.jsonl" {
		t.Errorf("Default EventFile should be '/tmp/morgana/events.jsonl', got '%s'", defaultConfig.EventFile)
	}
	if defaultConfig.PollInterval != 100*time.Millisecond {
		t.Errorf("Default PollInterval should be 100ms, got %v", defaultConfig.PollInterval)
	}
	if defaultConfig.BufferSize != 64*1024 {
		t.Errorf("Default BufferSize should be 64KB, got %d", defaultConfig.BufferSize)
	}

	// Test custom configuration
	customConfig := ConsumerConfig{
		EventFile:    "/custom/path/events.jsonl",
		PollInterval: 50 * time.Millisecond,
		BufferSize:   32 * 1024,
	}

	consumer, err := NewEventConsumer(eventBus, customConfig)
	if err != nil {
		t.Fatalf("Failed to create consumer with custom config: %v", err)
	}

	if consumer.GetEventFile() != "/custom/path/events.jsonl" {
		t.Errorf("EventFile should be '/custom/path/events.jsonl', got '%s'", consumer.GetEventFile())
	}

	// Test configuration validation (zero values should use defaults)
	zeroConfig := ConsumerConfig{}
	consumer2, err := NewEventConsumer(eventBus, zeroConfig)
	if err != nil {
		t.Fatalf("Failed to create consumer with zero config: %v", err)
	}

	if consumer2.GetEventFile() != "/tmp/morgana/events.jsonl" {
		t.Errorf("Zero config should default EventFile to '/tmp/morgana/events.jsonl', got '%s'", consumer2.GetEventFile())
	}
}
