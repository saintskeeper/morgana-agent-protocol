# EventConsumer Integration Guide

The `EventConsumer` provides a simple, one-way event consumption mechanism that
reads events from `/tmp/morgana/events.jsonl` and publishes them to the existing
EventBus system.

## Features

- **File-based monitoring**: Monitors `/tmp/morgana/events.jsonl` using fsnotify
  or polling fallback
- **File rotation handling**: Gracefully handles log rotation scenarios
- **Flexible event format**: Supports both IPCMessage format and raw event JSON
- **High performance**: Uses buffered reading and asynchronous event publishing
- **Integration-friendly**: Drops into existing EventBus architecture seamlessly

## Usage

### Basic Integration

```go
import (
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
)

// Create event bus
eventBus := events.NewEventBus(events.DefaultBusConfig())
defer eventBus.Close()

// Create consumer with default config
config := monitor.DefaultConsumerConfig()
consumer, err := monitor.NewEventConsumer(eventBus, config)
if err != nil {
    log.Fatal(err)
}

// Start consuming events
if err := consumer.Start(); err != nil {
    log.Fatal(err)
}
defer consumer.Stop()
```

### Configuration Options

```go
config := monitor.ConsumerConfig{
    EventFile:    "/tmp/morgana/events.jsonl", // Path to events file
    PollInterval: 100 * time.Millisecond,     // Polling interval (fallback)
    BufferSize:   64 * 1024,                  // File read buffer size
}
```

### Adding to Existing Monitor

To integrate the consumer with the existing morgana-monitor, you can add a
command-line flag:

```go
// Add flag
consumerMode := flag.Bool("consumer", false, "Use file-based event consumer instead of IPC")

// In your main function
if *consumerMode {
    // Use EventConsumer
    config := monitor.DefaultConsumerConfig()
    consumer, err := monitor.NewEventConsumer(eventBus, config)
    if err != nil {
        log.Fatalf("Failed to create consumer: %v", err)
    }
    if err := consumer.Start(); err != nil {
        log.Fatalf("Failed to start consumer: %v", err)
    }
    defer consumer.Stop()
} else {
    // Use existing IPC server
    server := monitor.NewIPCServer(socketPath, eventBus)
    // ... existing server code
}
```

## Event Format

The consumer supports multiple event formats:

### Morgana Protocol Format

```json
{
  "event_type": "task.started",
  "task_id": "abc123",
  "timestamp": "2025-08-11T14:53:03-04:00",
  "agent_type": "code-implementer",
  "prompt": "Create a simple function"
}
```

### Legacy Format

```json
{
  "event_type": "task_started",
  "task_id": "abc123",
  "agent_type": "code-implementer",
  "timestamp": "2025-08-11T18:32:48.827862+00:00"
}
```

### IPCMessage Format

```json
{
  "type": "task.started",
  "task_id": "abc123",
  "timestamp": "2025-08-11T14:53:03-04:00",
  "data": { "agent_type": "code-implementer" }
}
```

## Monitoring and Debugging

The consumer provides several methods for monitoring:

```go
// Check if running
if consumer.IsRunning() {
    log.Printf("Consumer is active")
}

// Get current file position
offset := consumer.GetCurrentOffset()
log.Printf("Current read offset: %d", offset)

// Get monitored file path
file := consumer.GetEventFile()
log.Printf("Monitoring: %s", file)
```

## Performance Characteristics

- **Memory usage**: ~64KB buffer + event queue (configurable)
- **CPU usage**: Minimal when using fsnotify, low when polling
- **Latency**: Sub-100ms event processing (typically <10ms)
- **Throughput**: Handles thousands of events/second
- **File handling**: Graceful rotation, automatic recovery

## Comparison with IPC Server

| Feature          | EventConsumer       | IPC Server        |
| ---------------- | ------------------- | ----------------- |
| Communication    | One-way (read-only) | Bidirectional     |
| Setup complexity | Simple              | Moderate          |
| File handling    | Native support      | Via client        |
| Resource usage   | Lower               | Higher            |
| Scalability      | High                | Moderate          |
| History replay   | File-based          | Memory buffer     |
| Reliability      | High (file-based)   | Network dependent |

## Migration from IPC

If you're currently using the IPC server and want to migrate:

1. Replace `monitor.NewIPCServer()` with `monitor.NewEventConsumer()`
2. Remove socket-related cleanup code
3. Ensure events are written to `/tmp/morgana/events.jsonl`
4. Remove client-server connection handling
5. Update any code that relies on bidirectional communication

The EventConsumer is designed to be a drop-in replacement for the receiving side
of the IPC architecture while providing better performance and simpler
deployment.
