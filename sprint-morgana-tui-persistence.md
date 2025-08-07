# Sprint Plan: Morgana TUI Persistent Monitoring

## Sprint Goal

Enable persistent, global TUI monitoring for all Morgana Protocol agent
executions, ensuring real-time visibility whenever the director model is active.

## Problem Statement

The TUI monitoring system shows no live data because:

1. Each morgana execution creates an isolated EventBus and ephemeral TUI
2. AgentAdapter doesn't pass --tui flag to morgana
3. No persistent daemon exists for continuous monitoring
4. No IPC mechanism for cross-process event sharing

## Success Criteria

- [ ] TUI monitor runs as persistent daemon showing all agent activity
- [ ] Works globally whenever director model executes agents
- [ ] Zero configuration required - auto-starts when needed
- [ ] Aggregates events from all morgana processes
- [ ] Maintains history across multiple executions
- [ ] Clean shutdown and resource management

## Technical Architecture

### Component Overview

```
┌─────────────────┐     IPC Socket      ┌──────────────────┐
│ morgana client  │────────────────────▶│ morgana-monitor  │
│   (agent exec)  │  /tmp/morgana.sock  │    (daemon)      │
│                 │                      │                  │
│ - EventBus      │     JSON Events     │ - EventBus       │
│ - IPCForwarder  │────────────────────▶│ - IPCReceiver    │
│ - Orchestrator  │                      │ - TUI Display    │
└─────────────────┘                      └──────────────────┘
     ▲                                            │
     │                                            ▼
┌─────────────────┐                      ┌──────────────────┐
│ AgentAdapter.sh │                      │   Terminal UI    │
│                 │                      │                  │
│ - Auto-start    │                      │ - Live Updates   │
│ - Health check  │                      │ - Event History  │
└─────────────────┘                      └──────────────────┘
```

## Implementation Tasks

### Phase 1: Core Infrastructure (Day 1)

#### Task 1.1: Create morgana-monitor daemon binary

**Priority**: P0-Critical **Estimate**: 3 hours **Files to create**:

- `morgana-protocol/cmd/morgana-monitor/main.go`
- `morgana-protocol/internal/monitor/server.go`

**Implementation**:

```go
// cmd/morgana-monitor/main.go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // Create persistent event bus
    eventBus := events.NewEventBus(events.DefaultBusConfig())
    defer eventBus.Close()

    // Start TUI with optimized config
    tuiConfig := tui.CreateOptimizedConfig()
    tuiInstance, _ := tui.RunAsync(ctx, eventBus, tuiConfig)
    defer tuiInstance.Stop()

    // Start IPC server
    server := monitor.NewIPCServer("/tmp/morgana.sock", eventBus)
    go server.Start(ctx)

    // Wait for shutdown
    <-sigChan
}
```

#### Task 1.2: Implement IPC server

**Priority**: P0-Critical **Estimate**: 2 hours **Files to create**:

- `morgana-protocol/internal/monitor/ipc_server.go`
- `morgana-protocol/internal/monitor/protocol.go`

**Implementation**:

```go
// internal/monitor/ipc_server.go
package monitor

import (
    "context"
    "encoding/json"
    "net"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

type IPCServer struct {
    socketPath string
    eventBus   events.EventBus
    listener   net.Listener
}

func NewIPCServer(socketPath string, eventBus events.EventBus) *IPCServer {
    return &IPCServer{
        socketPath: socketPath,
        eventBus:   eventBus,
    }
}

func (s *IPCServer) Start(ctx context.Context) error {
    // Remove existing socket
    os.Remove(s.socketPath)

    // Create Unix domain socket
    listener, err := net.Listen("unix", s.socketPath)
    if err != nil {
        return err
    }
    s.listener = listener

    // Accept connections
    go s.acceptConnections(ctx)

    return nil
}

func (s *IPCServer) acceptConnections(ctx context.Context) {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            continue
        }
        go s.handleConnection(ctx, conn)
    }
}

func (s *IPCServer) handleConnection(ctx context.Context, conn net.Conn) {
    defer conn.Close()

    decoder := json.NewDecoder(conn)
    for {
        var msg IPCMessage
        if err := decoder.Decode(&msg); err != nil {
            return
        }

        // Reconstruct and publish event
        event := s.reconstructEvent(msg)
        s.eventBus.Publish(event)
    }
}
```

### Phase 2: Client Integration (Day 1-2)

#### Task 2.1: Create IPC client forwarder

**Priority**: P0-Critical **Estimate**: 2 hours **Files to create**:

- `morgana-protocol/internal/monitor/ipc_client.go`

**Implementation**:

```go
// internal/monitor/ipc_client.go
package monitor

import (
    "encoding/json"
    "net"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

type IPCClient struct {
    socketPath string
    eventBus   events.EventBus
    conn       net.Conn
}

func NewIPCClient(socketPath string, eventBus events.EventBus) *IPCClient {
    return &IPCClient{
        socketPath: socketPath,
        eventBus:   eventBus,
    }
}

func (c *IPCClient) Connect() error {
    conn, err := net.Dial("unix", c.socketPath)
    if err != nil {
        return err // Monitor not running
    }
    c.conn = conn

    // Subscribe to all events and forward
    c.eventBus.SubscribeAll(c.forwardEvent)

    return nil
}

func (c *IPCClient) forwardEvent(event events.Event) {
    if c.conn == nil {
        return
    }

    msg := IPCMessage{
        Type:      string(event.Type()),
        TaskID:    event.TaskID(),
        Timestamp: event.Timestamp(),
        Data:      event,
    }

    encoder := json.NewEncoder(c.conn)
    encoder.Encode(msg)
}
```

#### Task 2.2: Integrate IPC client in main morgana

**Priority**: P0-Critical **Estimate**: 1 hour **Files to modify**:

- `morgana-protocol/cmd/morgana/main.go`

**Changes**:

```go
// Add after EventBus creation (line 146)
// Try to connect to monitor daemon
if !*enableTUI { // Only forward if not running own TUI
    ipcClient := monitor.NewIPCClient("/tmp/morgana.sock", eventBus)
    if err := ipcClient.Connect(); err == nil {
        log.Println("📡 Connected to Morgana Monitor for live monitoring")
        defer ipcClient.Close()
    }
}
```

### Phase 3: Shell Integration (Day 2)

#### Task 3.1: Update AgentAdapter for auto-start

**Priority**: P1-High **Estimate**: 1 hour **Files to modify**:

- `scripts/morgana-adapter.sh`

**Changes**:

```bash
# Add monitor management functions
function ensure_morgana_monitor() {
    local socket_path="/tmp/morgana.sock"

    # Check if monitor is running
    if [ ! -S "$socket_path" ]; then
        echo "🚀 Starting Morgana Monitor for live TUI..." >&2

        # Start monitor daemon
        nohup morgana-monitor > /tmp/morgana-monitor.log 2>&1 &
        local monitor_pid=$!

        # Wait for socket to be created (max 2 seconds)
        local counter=0
        while [ ! -S "$socket_path" ] && [ $counter -lt 20 ]; do
            sleep 0.1
            counter=$((counter + 1))
        done

        if [ -S "$socket_path" ]; then
            echo "✅ Morgana Monitor started (PID: $monitor_pid)" >&2
            echo "📺 View live TUI: tail -f /tmp/morgana-monitor.log" >&2
        else
            echo "⚠️ Monitor failed to start, continuing without TUI" >&2
        fi
    fi
}

# Update AgentAdapter function
function AgentAdapter() {
    local agent_type="$1"
    local prompt="$2"
    shift 2

    # Ensure monitor is running
    ensure_morgana_monitor

    # Execute with IPC enabled
    "$morgana_cmd" -- --agent "$agent_type" --prompt "$prompt" $@
}
```

#### Task 3.2: Create monitor control script

**Priority**: P1-High **Estimate**: 30 minutes **Files to create**:

- `scripts/morgana-monitor-ctl.sh`

**Implementation**:

```bash
#!/bin/bash
# Morgana Monitor Control Script

case "$1" in
    start)
        if pgrep -f morgana-monitor > /dev/null; then
            echo "Monitor already running"
        else
            morgana-monitor &
            echo "Monitor started"
        fi
        ;;
    stop)
        pkill -f morgana-monitor
        rm -f /tmp/morgana.sock
        echo "Monitor stopped"
        ;;
    status)
        if [ -S /tmp/morgana.sock ]; then
            echo "Monitor is running"
        else
            echo "Monitor is not running"
        fi
        ;;
    attach)
        # Attach to running monitor's terminal
        screen -r morgana-monitor || tmux attach -t morgana-monitor
        ;;
    *)
        echo "Usage: $0 {start|stop|status|attach}"
        exit 1
        ;;
esac
```

### Phase 4: Testing & Polish (Day 2-3)

#### Task 4.1: Integration tests

**Priority**: P1-High **Estimate**: 2 hours **Files to create**:

- `morgana-protocol/internal/monitor/ipc_test.go`
- `morgana-protocol/cmd/morgana-monitor/main_test.go`

#### Task 4.2: Update documentation

**Priority**: P2-Medium **Estimate**: 1 hour **Files to modify**:

- `README.md`
- `docs/monitoring.md`

#### Task 4.3: Add graceful shutdown

**Priority**: P2-Medium **Estimate**: 1 hour

- Implement PID file management
- Clean socket removal on shutdown
- Save TUI state for recovery

## Risk Mitigation

### Risk 1: Socket permission issues

**Mitigation**: Use user-specific socket path: `$HOME/.morgana/morgana.sock`

### Risk 2: Monitor crashes

**Mitigation**: Add systemd service file for auto-restart:

```ini
[Unit]
Description=Morgana Protocol Monitor
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/morgana-monitor
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Risk 3: Performance impact

**Mitigation**:

- Batch event forwarding (send every 100ms)
- Use circular buffer for event history
- Implement event filtering by severity

## Validation Checklist

### Functional Requirements

- [ ] Monitor starts automatically when first agent runs
- [ ] TUI displays events from multiple concurrent agents
- [ ] History persists across agent executions
- [ ] Clean shutdown without data loss
- [ ] Works with parallel agent execution

### Performance Requirements

- [ ] < 100ms latency for event display
- [ ] < 5% CPU overhead for monitoring
- [ ] Handles 1000+ events/second
- [ ] Memory usage < 50MB for monitor daemon

### User Experience

- [ ] Zero configuration required
- [ ] Clear status indicators
- [ ] Intuitive keyboard shortcuts
- [ ] Helpful error messages
- [ ] Documentation and examples

## Timeline

**Day 1**:

- Morning: Create morgana-monitor daemon (Tasks 1.1, 1.2)
- Afternoon: Implement IPC client/server (Tasks 2.1, 2.2)

**Day 2**:

- Morning: Update shell integration (Tasks 3.1, 3.2)
- Afternoon: Initial testing and debugging

**Day 3**:

- Morning: Integration tests (Task 4.1)
- Afternoon: Documentation and polish (Tasks 4.2, 4.3)

## Definition of Done

- All unit tests pass
- Integration tests verify end-to-end flow
- Documentation updated with examples
- Code reviewed and approved
- Performance benchmarks meet requirements
- Successfully deployed and tested in production environment
