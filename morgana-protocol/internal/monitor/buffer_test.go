package monitor

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEventBuffer(t *testing.T) {
	t.Run("AddAndRetrieveEvents", func(t *testing.T) {
		server := &IPCServer{
			eventBuffer:   make([]IPCMessage, 1000),
			bufferMutex:   sync.RWMutex{},
			maxBufferSize: 1000,
			bufferIndex:   0,
			bufferCount:   0,
		}

		// Add some events
		for i := 0; i < 10; i++ {
			msg := IPCMessage{
				Type:      "test_event",
				TaskID:    fmt.Sprintf("task_%d", i),
				Timestamp: time.Now(),
				Data:      map[string]interface{}{"index": i},
			}
			server.addToBuffer(msg)
		}

		// Retrieve buffered events
		events := server.getBufferedEvents()
		if len(events) != 10 {
			t.Errorf("Expected 10 events, got %d", len(events))
		}

		// Verify order (oldest to newest)
		for i, event := range events {
			if data, ok := event.Data.(map[string]interface{}); ok {
				if idx, ok := data["index"].(int); ok {
					if idx != i {
						t.Errorf("Event %d has wrong index: %d", i, idx)
					}
				}
			}
		}
	})

	t.Run("CircularBufferOverwrite", func(t *testing.T) {
		server := &IPCServer{
			eventBuffer:   make([]IPCMessage, 10), // Small buffer for testing
			bufferMutex:   sync.RWMutex{},
			maxBufferSize: 10,
			bufferIndex:   0,
			bufferCount:   0,
		}

		// Add more events than buffer size
		for i := 0; i < 15; i++ {
			msg := IPCMessage{
				Type:      "test_event",
				TaskID:    fmt.Sprintf("task_%d", i),
				Timestamp: time.Now(),
				Data:      map[string]interface{}{"index": i},
			}
			server.addToBuffer(msg)
		}

		// Should only have the last 10 events (5-14)
		events := server.getBufferedEvents()
		if len(events) != 10 {
			t.Errorf("Expected 10 events, got %d", len(events))
		}

		// First event should be index 5 (oldest that wasn't overwritten)
		if data, ok := events[0].Data.(map[string]interface{}); ok {
			if idx, ok := data["index"].(int); ok {
				if idx != 5 {
					t.Errorf("First event should have index 5, got %d", idx)
				}
			}
		}

		// Last event should be index 14
		if data, ok := events[9].Data.(map[string]interface{}); ok {
			if idx, ok := data["index"].(int); ok {
				if idx != 14 {
					t.Errorf("Last event should have index 14, got %d", idx)
				}
			}
		}
	})

	t.Run("ThreadSafety", func(t *testing.T) {
		server := &IPCServer{
			eventBuffer:   make([]IPCMessage, 1000),
			bufferMutex:   sync.RWMutex{},
			maxBufferSize: 1000,
			bufferIndex:   0,
			bufferCount:   0,
		}

		var wg sync.WaitGroup
		numGoroutines := 10
		eventsPerGoroutine := 100

		// Concurrent writes
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for i := 0; i < eventsPerGoroutine; i++ {
					msg := IPCMessage{
						Type:      "concurrent_event",
						TaskID:    fmt.Sprintf("task_%d_%d", goroutineID, i),
						Timestamp: time.Now(),
						Data:      map[string]interface{}{"goroutine": goroutineID, "index": i},
					}
					server.addToBuffer(msg)
				}
			}(g)
		}

		// Concurrent reads
		for r := 0; r < 5; r++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 50; i++ {
					_ = server.getBufferedEvents()
					time.Sleep(time.Microsecond) // Small delay to increase contention
				}
			}()
		}

		wg.Wait()

		// Verify we have the expected number of events
		events := server.getBufferedEvents()
		expectedCount := numGoroutines * eventsPerGoroutine
		if expectedCount > 1000 {
			expectedCount = 1000 // Buffer size limit
		}
		if len(events) != expectedCount {
			t.Errorf("Expected %d events, got %d", expectedCount, len(events))
		}
	})

	t.Run("EmptyBuffer", func(t *testing.T) {
		server := &IPCServer{
			eventBuffer:   make([]IPCMessage, 1000),
			bufferMutex:   sync.RWMutex{},
			maxBufferSize: 1000,
			bufferIndex:   0,
			bufferCount:   0,
		}

		events := server.getBufferedEvents()
		if len(events) != 0 {
			t.Errorf("Expected empty buffer, got %d events", len(events))
		}
	})
}

func BenchmarkEventBuffer(b *testing.B) {
	server := &IPCServer{
		eventBuffer:   make([]IPCMessage, 1000),
		bufferMutex:   sync.RWMutex{},
		maxBufferSize: 1000,
		bufferIndex:   0,
		bufferCount:   0,
	}

	msg := IPCMessage{
		Type:      "benchmark_event",
		TaskID:    "bench_task",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"test": "data"},
	}

	b.Run("AddToBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			server.addToBuffer(msg)
		}
	})

	// Pre-fill buffer for read benchmark
	for i := 0; i < 1000; i++ {
		server.addToBuffer(msg)
	}

	b.Run("GetBufferedEvents", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = server.getBufferedEvents()
		}
	})
}
