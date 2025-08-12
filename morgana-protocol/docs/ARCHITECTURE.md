# Morgana Protocol Architecture

## Overview

Morgana Protocol has been transformed from a complex IPC orchestrator into a
simplified event monitoring system. The new architecture prioritizes real-time
event streaming, high-performance data processing, and beautiful visualization
through a Terminal User Interface (TUI).

## Design Philosophy

### Simplified Architecture Benefits

1. **Zero Configuration**: Works out of the box with sensible defaults
2. **High Performance**: 5M+ events/sec with minimal overhead
3. **Real-time Monitoring**: Live event streaming via Unix sockets
4. **Event Persistence**: Circular buffering prevents data loss
5. **Beautiful Visualization**: Rich TUI with customizable themes

### Key Design Decisions

- **Event-driven**: Everything flows through the high-performance event bus
- **Headless First**: Monitor daemon runs independently of visualization
- **Buffer-centric**: Circular buffer ensures no event loss
- **Unix Socket IPC**: Efficient local inter-process communication
- **Lock-free Design**: High concurrency with zero contention

## System Architecture

### High-Level Overview

```
┌─────────────────┐    Event Stream    ┌─────────────────┐
│  Agent Tasks    │ ──────────────────> │  Event Monitor  │
│  (morgana)      │   (High-perf bus)   │  (daemon)       │
└─────────────────┘                     └─────────────────┘
         │                                       │
         ▼                                       ▼
┌─────────────────┐                     ┌─────────────────┐
│   Event Bus     │                     │ Circular Buffer │
│ (5M+ evt/sec)   │                     │ (1000 events)   │
└─────────────────┘                     └─────────────────┘
                                                 │
                                         Unix Socket
                                                 │
                                        ┌─────────────────┐
                                        │   TUI Client    │
                                        │  (real-time)    │
                                        └─────────────────┘
```

### Component Details

#### 1. Event Bus (internal/events)

**Purpose**: High-performance event distribution system

**Key Features**:

- Lock-free circular buffer design
- 5M+ events/sec throughput (async mode)
- Thread-safe concurrent access
- Zero allocations after initialization
- Configurable worker pool

**Architecture**:

```go
type EventBus struct {
    buffer    *CircularBuffer
    workers   []*Worker
    config    BusConfig
    metrics   *Metrics
}
```

**Performance Characteristics**:

- Async Publishing: 5M+ events/sec (203ns/op)
- Sync Publishing: 17M+ events/sec (59ns/op)
- Memory: Zero allocations after warmup
- Overhead: <5% compared to direct calls

#### 2. Monitor Daemon (cmd/morgana-monitor)

**Purpose**: Headless event collection and distribution

**Key Features**:

- Unix socket server for TUI clients
- Circular event buffer (default: 1000 events)
- Automatic history replay for new clients
- Multi-client support
- Graceful shutdown handling

**Event Flow**:

1. Receives events from agent processes via event bus
2. Stores events in circular buffer
3. Forwards events to connected TUI clients
4. Provides history replay for late-joining clients

#### 3. Terminal User Interface (internal/tui)

**Purpose**: Real-time event visualization

**Key Features**:

- Beautiful, themeable interface
- Real-time event streaming
- Agent-specific progress tracking
- Event filtering and search
- Export capabilities

**UI Components**:

- Task list with progress bars
- Live event stream
- System statistics panel
- Agent-specific color coding

## Event System Design

### Event Types Hierarchy

```
Event (interface)
├── TaskEvents
│   ├── TaskStartedEvent
│   ├── TaskProgressEvent
│   ├── TaskCompletedEvent
│   └── TaskFailedEvent
├── OrchestratorEvents
│   ├── OrchestratorStartedEvent
│   ├── OrchestratorCompletedEvent
│   └── OrchestratorFailedEvent
└── AdapterEvents
    ├── AdapterValidationEvent
    ├── AdapterPromptLoadEvent
    └── AdapterExecutionEvent
```

### Event Flow Architecture

```
[Agent Task] → [Event Bus] → [Monitor Daemon] → [TUI Client]
     │              │              │                │
     ▼              ▼              ▼                ▼
[TaskStarted]  [Publish]     [Buffer]         [Display]
[TaskProgress] [Route]       [Forward]        [Progress]
[TaskComplete] [Batch]       [Replay]         [Update]
```

## Inter-Process Communication

### Unix Socket Protocol

Morgana uses a lightweight protocol over Unix sockets:

#### Message Types

- `event`: Regular event message from agent to monitor
- `request`: Client request (e.g., history replay)
- `replay`: Server response with buffered events

#### Protocol Flow

```
TUI Client                    Monitor Daemon
    │                              │
    ├─── Connect Unix Socket ────> │
    │                              │
    ├─── Request History ────────> │
    │                              │
    │ <──── Replay Events ──────── │
    │                              │
    │ <──── Live Events ────────── │
    │          (streaming)         │
```

### Event Buffering Strategy

#### Circular Buffer Design

- Fixed size buffer (default: 1000 events)
- Lock-free operations using atomic pointers
- Zero-copy event replay for new clients
- Automatic overflow handling (oldest events dropped)

#### Benefits

- **No data loss**: Events buffered even when no TUI attached
- **Late joining**: New clients receive full history
- **Multiple clients**: Each client gets same event stream
- **High performance**: Zero allocations during normal operation

## Configuration Architecture

### YAML-based Configuration

```yaml
tui:
  enabled: true
  performance:
    refresh_rate: 50ms
    max_log_lines: 1000
    target_fps: 20
  visual:
    theme:
      name: dark
      primary: "#7C3AED"
  events:
    buffer_size: 1000
    enable_batching: true
```

### Environment Overrides

All YAML settings can be overridden via environment variables:

- `MORGANA_TUI_ENABLED=true`
- `MORGANA_TUI_REFRESH_RATE=16ms`
- `MORGANA_TUI_THEME=dark`
- `MORGANA_EVENT_BUFFER_SIZE=2000`

## Performance Characteristics

### Benchmarks (Apple M1 Max)

| Operation           | Throughput      | Latency  | Memory   |
| ------------------- | --------------- | -------- | -------- |
| Async Event Publish | 5M+ events/sec  | 203ns/op | 0 allocs |
| Sync Event Publish  | 17M+ events/sec | 59ns/op  | 0 allocs |
| Event Buffer Read   | 50M+ reads/sec  | 20ns/op  | 0 allocs |
| Unix Socket IPC     | 1M+ msg/sec     | 1µs/op   | Minimal  |

### Resource Usage

- **Memory**: <50MB for daemon + TUI combined
- **CPU**: <5% overhead during normal operation
- **Network**: Local Unix sockets only
- **Storage**: No persistent storage required

## Scalability & Limits

### Current Limits

- Event buffer: 1000 events (configurable)
- Concurrent clients: Unlimited (practical limit ~100)
- Event throughput: 5M+ events/sec
- Event types: Unlimited (extensible design)

### Scalability Characteristics

- **Horizontal**: Multiple monitor daemons possible
- **Vertical**: Scales with available CPU cores
- **Memory**: O(1) memory usage (circular buffer)
- **Network**: Local-only (no network overhead)

## Security Model

### Local-Only Design

- Unix socket communication (no network exposure)
- File system permissions for socket access
- Process isolation between components
- No external dependencies or services

### Access Control

- Socket permissions control TUI access
- No authentication required (local system trust)
- Read-only event access (no command injection)
- Isolated process model

## Future Architecture Considerations

### Potential Enhancements

1. **Remote Monitoring**: Optional HTTP/WebSocket server
2. **Event Persistence**: Optional database backend
3. **Distributed Deployment**: Multi-host event aggregation
4. **Plugin System**: Custom event processors
5. **Metrics Export**: Prometheus/OpenTelemetry integration

### Backwards Compatibility

The current architecture is designed to maintain backwards compatibility:

- Existing agent interfaces unchanged
- Configuration migration supported
- Graceful degradation when components unavailable
- Optional features can be disabled

## Implementation Details

### Key Files

- `internal/events/`: Event system implementation
- `internal/monitor/`: Monitor daemon and IPC protocol
- `internal/tui/`: Terminal user interface
- `cmd/morgana-monitor/`: Monitor daemon entry point
- `cmd/morgana/`: Main agent runner

### Dependencies

- **Minimal**: Only standard Go libraries + minimal external deps
- **No Runtime**: No Python/Node.js dependencies for monitoring
- **Cross-platform**: Works on Linux, macOS, Windows

The simplified architecture achieves the goal of real-time agent monitoring
while maintaining high performance and minimal complexity. The event-driven
design provides excellent observability without the overhead of complex
orchestration systems.
