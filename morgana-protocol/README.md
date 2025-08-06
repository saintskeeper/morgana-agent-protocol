# Morgana Protocol

A Go-based agent orchestration system that enables parallel execution of
specialized agents in Claude Code with comprehensive observability and
integration testing.

## Features

- 🚀 **True Parallel Execution** - Goroutine-based concurrency with resource
  limits
- 🔍 **OpenTelemetry Tracing** - Full observability of agent execution lifecycle
- ⏱️ **Configurable Timeouts** - Per-agent timeout configuration
- 🧪 **Integration Testing** - Comprehensive test suite for multi-language
  bridge
- 📝 **YAML Configuration** - Flexible configuration with environment overrides
- 🐍 **Python Bridge** - Seamless integration with Claude Code's Task function

## Quick Start

```bash
# Build
make build

# Run tests (including integration tests)
make test-all

# Install locally
make install

# Run with configuration
morgana --config morgana.yaml -- --agent code-implementer --prompt "Hello"
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
# Single agent with command line arguments
morgana -- --agent code-implementer --prompt "implement auth service"

# Multiple agents via JSON input
echo '[
  {"agent_type":"code-implementer","prompt":"implement auth"},
  {"agent_type":"test-specialist","prompt":"create tests"}
]' | morgana --parallel

# With custom configuration
morgana --config custom.yaml -- --agent sprint-planner --prompt "plan feature"

# Mock mode for testing (no Python required)
morgana --config morgana.yaml  # with mock_mode: true in config
```

### Integration with Claude Code

1. Source the wrapper in your shell:

```bash
source ~/.claude/scripts/agent-adapter-wrapper.sh
```

2. Use in markdown commands:

```bash
# Single agent execution
AgentAdapter "code-implementer" "implement user authentication"

# Parallel execution with multiple agents
morgana_parallel << 'EOF'
[
  {"agent_type":"code-implementer","prompt":"implement feature"},
  {"agent_type":"test-specialist","prompt":"write comprehensive tests"}
]
EOF
```

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   CLI/Shell     │────▶│  Go Orchestrator │────▶│  Python Bridge  │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │                          │
                               ▼                          ▼
                        ┌─────────────────┐       ┌─────────────────┐
                        │   Goroutines    │       │   Task Tool     │
                        │  (Parallelism)  │       │  (Claude Code)  │
                        └─────────────────┘       └─────────────────┘
```

- **Go Orchestrator**: Handles concurrency, timeouts, and agent management
- **Python Bridge**: Interfaces with Claude Code's Task function
- **OpenTelemetry**: Provides distributed tracing and observability
- **YAML Config**: Flexible configuration with environment overrides

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

Start the monitoring stack to visualize agent execution:

```bash
cd monitoring
docker-compose up -d

# Access services:
# - Jaeger UI: http://localhost:16686
# - Grafana: http://localhost:3000
# - Prometheus: http://localhost:9090
```

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

## Troubleshooting

1. **Python not found**: Ensure Python 3 is installed and in PATH
2. **Bridge script not found**: Set `MORGANA_BRIDGE_PATH` environment variable
3. **Timeout errors**: Adjust timeouts in morgana.yaml
4. **Mock mode**: Use `mock_mode: true` in config for testing without Python
