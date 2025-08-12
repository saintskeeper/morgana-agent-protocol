# Morgana Protocol Integration Guide

## Overview

The Morgana Protocol is a simplified event monitoring system designed for
real-time observation of specialized agents in Claude Code. It provides a
light-weight, high-performance event streaming architecture with beautiful
Terminal User Interface (TUI) for live monitoring.

## Architecture

Morgana uses a simplified event streaming architecture:

```
┌─────────────────┐      Event Stream       ┌─────────────────┐
│  Agent Tasks   │ ──────────────────────> │  Event Monitor  │
│  (morgana)     │     (Unix Socket)      │  (daemon)      │
└─────────────────┘                        └─────────────────┘
         │                                      │
         ▼                                      ▼
┌─────────────────┐                        ┌─────────────────┐
│  Event Bus      │                        │  Event Buffer   │
│  (5M+ evt/sec)  │                        │  (Circular)     │
└─────────────────┘                        └─────────────────┘
                                                      │
                                                      ▼
                                             ┌─────────────────┐
                                             │   TUI Client    │
                                             │  (Real-time)   │
                                             └─────────────────┘
```

## Components

### 1. Event Bus (`internal/events`)

- High-performance event publishing (5M+ events/sec)
- Lock-free circular buffer design
- Thread-safe concurrent access
- Zero allocations after initialization

### 2. Monitor Daemon (`morgana-monitor`)

- Headless event collection daemon
- Unix socket server for TUI clients
- Circular event buffer (1000 events)
- Automatic event replay for new clients

### 3. Terminal Interface (TUI)

- Beautiful real-time event visualization
- Agent-specific progress tracking
- Customizable themes and layouts
- Export and filtering capabilities

## Usage

### Basic Usage

```bash
# Start the monitor daemon
make up

# Run a single agent task
morgana -- --agent code-implementer --prompt "Write a hello world function"

# Run parallel agents
echo '[
  {"agent_type":"code-implementer","prompt":"Write function"},
  {"agent_type":"test-specialist","prompt":"Write tests"}
]' | morgana --parallel

# View live monitoring
make attach
```

### Event Monitoring

```bash
# Start monitor daemon in background
make up

# Check daemon status
make status

# Attach TUI to view live events
make attach

# Run tasks while monitoring
morgana -- --agent code-implementer --prompt "Implement auth service"

# Stop the monitor daemon
make down
```

### Event Stream Integration

```bash
# Custom event processing
echo '{"agent_type":"validation-expert","prompt":"Review code"}' | \
  morgana --tui --parallel

# High-performance batch processing
for task in task1 task2 task3; do
  morgana -- --agent code-implementer --prompt "$task" &
done
wait
```

## Configuration

### Environment Variables

- `MORGANA_TUI_ENABLED`: Enable/disable TUI mode (`true`/`false`)
- `MORGANA_TUI_REFRESH_RATE`: TUI refresh rate (e.g., `16ms`)
- `MORGANA_TUI_THEME`: Color theme (`dark`, `light`, `custom`)
- `MORGANA_EVENT_BUFFER_SIZE`: Event buffer size (default: 1000)

### YAML Configuration

```yaml
# morgana.yaml
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

## Installation

```bash
# Clone and build
git clone <repo-url>
cd morgana-protocol
make build && make install

# Start monitoring
make up

# Test with sample task
morgana -- --agent code-implementer --prompt "Hello world"
```

## Performance

- **Event Throughput**: 5M+ events/sec (async mode)
- **Memory Footprint**: <50MB for daemon + TUI
- **Latency**: <100ns per event publish
- **Buffer Capacity**: 1000 events with zero-copy replay

## Advanced Features

1. **Event Filtering**: Subscribe to specific event types
2. **Historical Replay**: Access buffered events from earlier runs
3. **Multi-client Support**: Multiple TUI instances can connect
4. **Export Capabilities**: Save event logs for analysis
5. **Custom Themes**: Personalize TUI appearance
