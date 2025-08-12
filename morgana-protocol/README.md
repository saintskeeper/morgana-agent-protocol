# Morgana Protocol

A simplified event monitoring system for specialized agents in Claude Code.
Morgana transforms from a complex IPC orchestrator to a streamlined event
monitoring platform with real-time visualization and high-performance event
streaming.

## Why the New Architecture?

The simplified Morgana Protocol offers significant benefits:

- **🚀 Zero Configuration**: Works immediately without complex setup
- **⚡ High Performance**: 5M+ events/sec with <5% overhead
- **📺 Beautiful Monitoring**: Rich TUI with real-time event visualization
- **🔄 Event Replay**: Never miss events with circular buffering
- **🎯 Focused Purpose**: Specialized for agent monitoring vs. general
  orchestration

## Features

- 📡 **Event Stream Architecture** - Real-time event monitoring via Unix sockets
- 🖥️ **Terminal User Interface** - Beautiful TUI for live agent monitoring
- 🔄 **Event Buffering & Replay** - Never miss events with circular buffering
- 🚀 **High-Performance Event Bus** - 5M+ events/sec with lock-free design
- ⚡ **Zero-Configuration** - Works out of the box with sensible defaults
- 🎯 **Specialized Agent Support** - Optimized for Claude Code agent workflows

## Quick Start

```bash
# Build and install
make build && make install

# Start the event monitor daemon
make up

# Run an agent and monitor in real-time
make attach &
morgana -- --agent code-implementer --prompt "Implement a hello world function"

# Check system status
make status
```

## Configuration

Morgana supports flexible configuration through YAML files and environment
variables:

```yaml
# morgana.yaml
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m
  timeouts:
    code-implementer: 5m
    test-specialist: 3m

execution:
  max_concurrency: 5
  default_mode: sequential

telemetry:
  enabled: true
  exporter: otlp # Options: stdout, otlp, none
  otlp_endpoint: localhost:4317
```

Environment variables override config file settings:

```bash
export MORGANA_DEBUG=true
export MORGANA_BRIDGE_PATH=/custom/path/to/task_bridge.py
```

## Usage

### Command Line Interface

```bash
# Start the monitor daemon
make up

# Single agent execution
morgana -- --agent code-implementer --prompt "implement auth service"

# Multiple agents via JSON input
echo '[
  {"agent_type":"code-implementer","prompt":"implement auth"},
  {"agent_type":"test-specialist","prompt":"create tests"}
]' | morgana --parallel

# Live monitoring with TUI
make attach

# Check system status
make status
```

### Integration with Claude Code

1. Start the monitor for live observation:

```bash
# In terminal 1: Start monitor and TUI
make up && make attach

# In terminal 2: Source the wrapper
source ~/.claude/scripts/agent-adapter-wrapper.sh
```

2. Use agents with live monitoring:

```bash
# Single agent execution (monitored)
AgentAdapter "code-implementer" "implement user authentication"

# Parallel execution (visible in TUI)
morgana_parallel << 'EOF'
[
  {"agent_type":"code-implementer","prompt":"implement feature"},
  {"agent_type":"test-specialist","prompt":"write comprehensive tests"}
]
EOF
```

## Architecture

Morgana uses a simplified event stream architecture:

```
┌─────────────────┐     Event Stream     ┌─────────────────┐
│   Agent Tasks   │ ─────────────────────▶│  Event Monitor  │
│   (morgana)     │      Unix Socket      │   (daemon)      │
└─────────────────┘                       └─────────────────┘
         │                                         │
         ▼                                         ▼
┌─────────────────┐                       ┌─────────────────┐
│  Event Bus      │                       │ Circular Buffer │
│  (5M+ evt/sec)  │                       │  (1000 events) │
└─────────────────┘                       └─────────────────┘
                                                   │
                                                   ▼
                                          ┌─────────────────┐
                                          │   TUI Client    │
                                          │   (real-time)   │
                                          └─────────────────┘
```

### Key Components

- **Event Stream**: Real-time event publishing via high-performance event bus
- **Monitor Daemon**: Headless event collection with circular buffering
- **TUI Interface**: Beautiful terminal interface for live monitoring
- **Unix Socket IPC**: Efficient inter-process communication
- **Zero Dependencies**: No complex orchestration or file-based IPC

## Available Agents

- `code-implementer` - Code implementation specialist
- `sprint-planner` - Sprint planning and decomposition
- `test-specialist` - Test creation with edge cases
- `validation-expert` - Code review and validation

## Testing

```bash
# Run unit tests only
make test

# Run integration tests
make test-integration

# Run all tests
make test-all

# Test Python bridge manually
./scripts/test_bridge.sh
```

## Monitoring

### Distributed Tracing

Start the monitoring stack to visualize agent execution:

```bash
cd monitoring
docker-compose up -d

# Access services:
# - Jaeger UI: http://localhost:16686
# - Grafana: http://localhost:3000
# - Prometheus: http://localhost:9090
```

### TUI Monitor

Morgana includes a built-in Terminal User Interface (TUI) for real-time
monitoring of agent execution with event buffering and replay capabilities.

#### Starting the Monitor

```bash
# Build the monitor
make dev

# Start monitor in headless mode (runs as daemon)
make up

# Check monitor status
make status

# Attach TUI client to view events
make attach

# Stop the monitor
make down
```

#### Event Buffering & Replay

The monitor features a circular event buffer that stores the last 1000 events in
memory. This enables:

- **Late-joining clients**: Connect to a running monitor and receive historical
  events
- **No data loss**: Events are buffered even when no TUI is attached
- **Headless operation**: Monitor runs as a daemon, collecting events for later
  viewing
- **Multiple clients**: Multiple TUI clients can connect and receive the same
  event stream

#### Event Flow Architecture

```
┌─────────────┐     Event Stream     ┌──────────────┐
│   Morgana   │ ──────────────────> │   Monitor    │
│   Process   │   (Real-time IPC)    │   Daemon     │
└─────────────┘                      │              │
                                     │ ┌──────────┐ │
                                     │ │ Events   │ │
                                     │ │ Buffer   │ │
                                     │ │ (1000)   │ │
                                     │ └──────────┘ │
                                     └──────────────┘
                                            │
                                     Unix Socket
                                            │
                                     ┌──────────────┐
                                     │  TUI Client  │
                                     │  (Live View) │
                                     └──────────────┘
```

#### Usage Examples

```bash
# Start monitor and run parallel tasks
make up

# In another terminal, attach TUI
make attach &

# Run parallel tasks (events visible in TUI)
echo '[
  {"agent_type":"code-implementer","prompt":"Implement user login"},
  {"agent_type":"test-specialist","prompt":"Write login tests"}
]' | morgana --parallel

# Later, new TUI clients see buffered history
make attach  # Will show buffered events from earlier execution
```

#### Event Stream Protocol

The monitor uses a simple event streaming protocol:

- **Event Publishing**: Morgana processes publish events to the event bus
- **Stream Forwarding**: Monitor daemon forwards events via Unix socket
- **History Replay**: New clients receive buffered events automatically
- **Live Updates**: Real-time event streaming for connected TUI clients
- **Zero Configuration**: No complex handshakes or protocol negotiation

## Development

```bash
# Install dependencies
make deps

# Format code
make fmt

# Build for all platforms
make build

# Clean build artifacts
make clean
```

## Migration Guide

If you're upgrading from the old IPC-based Morgana Protocol, see the
**[Migration Guide](docs/MIGRATION_GUIDE.md)** for step-by-step instructions.

**Quick migration:**

```bash
# 1. Clean up old configuration
./scripts/migrate-config.sh

# 2. Set up monitoring environment
./scripts/setup-monitoring.sh

# 3. Validate migration
./scripts/validate-migration.sh

# 4. Test new system
make up && make attach
```

## Troubleshooting

1. **Monitor won't start**: Run `make down && make up` to restart
2. **TUI display issues**: Try `export MORGANA_TUI_REFRESH_RATE=100ms`
3. **Events not showing**: Check `make status` and socket permissions
4. **Configuration errors**: Run `./scripts/validate-migration.sh`
