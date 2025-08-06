# Morgana Protocol TUI - User Guide

This guide covers the Terminal User Interface (TUI) for Morgana Protocol,
providing real-time monitoring and interactive control of agent orchestration.

## Quick Start

### Basic Usage

```bash
# Start with TUI enabled in optimized mode (default)
morgana --tui

# Enable TUI with development mode (more features, debug info)
morgana --tui --tui-mode dev

# High-performance mode (lower CPU usage, fewer features)
morgana --tui --tui-mode high-performance

# Execute specific tasks with TUI monitoring
morgana --tui --agent code-implementer --prompt "Build REST API"

# Parallel execution with TUI
echo '[{"agent_type":"code-implementer","prompt":"Task 1"}, {"agent_type":"test-specialist","prompt":"Task 2"}]' | morgana --tui --parallel
```

### TUI Modes

The TUI supports three operational modes optimized for different use cases:

#### Development Mode (`--tui-mode dev`)

- **Refresh Rate**: 30 FPS
- **Features**: All features enabled (filtering, search, export, debug info)
- **Buffer Size**: 500 events
- **Best For**: Interactive development, debugging, feature exploration

#### Optimized Mode (`--tui-mode optimized`) - Default

- **Refresh Rate**: 60 FPS
- **Features**: Essential features enabled (filtering, search)
- **Buffer Size**: 1000 events
- **Best For**: Regular monitoring, balanced performance/features

#### High-Performance Mode (`--tui-mode high-performance`)

- **Refresh Rate**: 30 FPS (reduced CPU usage)
- **Features**: Minimal features, optimized for throughput
- **Buffer Size**: 2000 events
- **Best For**: Production monitoring, high-load scenarios

## Interface Overview

The TUI provides a real-time dashboard with three main layout modes:

### Split Layout (Default)

```
┌─────────────────────────────────────────────────┐
│ 🧙‍♂️ Morgana Protocol TUI [Split] | Focus: Dashboard │
├─────────────────────────────────────────────────┤
│                                                 │
│           AGENT STATUS DASHBOARD                │
│                                                 │
├─────────────────────────────────────────────────┤
│                                                 │
│              REAL-TIME LOGS                     │
│                                                 │
└─────────────────────────────────────────────────┘
│ Uptime: 2m15s | Events: 1,247 | Press H for help, Q to quit │
```

### Dashboard-Only Layout

```
┌─────────────────────────────────────────────────┐
│ 🧙‍♂️ Morgana Protocol TUI [Dashboard] | Focus: Dashboard │
├─────────────────────────────────────────────────┤
│                                                 │
│                                                 │
│           AGENT STATUS DASHBOARD                │
│                                                 │
│                                                 │
└─────────────────────────────────────────────────┘
│ Uptime: 2m15s | Events: 1,247 | Press H for help, Q to quit │
```

### Logs-Only Layout

```
┌─────────────────────────────────────────────────┐
│ 🧙‍♂️ Morgana Protocol TUI [Logs] | Focus: Logs      │
├─────────────────────────────────────────────────┤
│                                                 │
│                                                 │
│              REAL-TIME LOGS                     │
│                                                 │
│                                                 │
└─────────────────────────────────────────────────┘
│ Uptime: 2m15s | Events: 1,247 | Press H for help, Q to quit │
```

## Components

### Agent Status Dashboard

The dashboard displays real-time status cards for each active agent:

```
┌─ code-implementer ─────┐ ┌─ test-specialist ─────────┐
│ Status: Running        │ │ Status: Completed         │
│ Progress: [████████░] │ │ Progress: [██████████]   │
│ Stage: execution       │ │ Stage: completed          │
│ Duration: 2m 15s       │ │ Duration: 1m 32s          │
│ Model: o3              │ │ Model: gpt-4              │
└────────────────────────┘ └───────────────────────────┘
```

**Status Indicators:**

- 🟢 **Running** - Agent is actively executing
- ✅ **Completed** - Agent finished successfully
- ❌ **Failed** - Agent encountered an error
- ⏸️ **Pending** - Agent waiting to start

**Progress Bars:**

- Animated progress indication
- Color-coded by status (green=success, red=error, blue=running)
- Real-time updates during execution

### Real-Time Log Viewer

The log viewer provides streaming updates of all agent activities:

```
2024-08-06 14:32:15 INFO  code-implementer  starting     Task started: code-implementer
2024-08-06 14:32:16 DEBUG code-implementer  validation   Agent validation successful
2024-08-06 14:32:16 INFO  code-implementer  prompt_load  Prompt loaded successfully
2024-08-06 14:32:18 DEBUG code-implementer  execution    Execution phase: planning
2024-08-06 14:32:45 DEBUG code-implementer  execution    Execution phase: implementation
2024-08-06 14:33:12 INFO  code-implementer  completed    Task completed successfully
```

**Log Levels:**

- `DEBUG` - Detailed execution steps
- `INFO` - Important milestones
- `WARN` - Non-critical issues
- `ERROR` - Failures and errors

**Columns:**

- **Timestamp** - Precise event timing
- **Level** - Log severity level
- **Agent** - Agent type performing the action
- **Stage** - Current execution phase
- **Message** - Human-readable description

## Keyboard Shortcuts

### Navigation

- `Tab` / `Shift+Tab` - Switch focus between dashboard and logs
- `Space` / `L` - Toggle between layout modes (Split → Dashboard → Logs → Split)
- `H` / `?` - Show/hide help screen
- `Q` / `Ctrl+C` / `Esc` - Quit application

### Log Viewer (when focused)

- `↑` / `K` - Scroll up one line
- `↓` / `J` - Scroll down one line
- `Page Up` - Scroll up one page (10 lines)
- `Page Down` - Scroll down one page (10 lines)
- `Home` - Jump to top of logs
- `End` - Jump to bottom of logs

### Filtering and Search

- `F` - Cycle through log filters
- `C` - Clear current filter
- `/` - Open search (development mode only)

**Available Filters:**

- `All` - Show all log entries
- `error` - Show only errors
- `warning` - Show warnings and errors
- `info` - Show info level and above
- `code-implementer` - Show only code-implementer events
- `sprint-planner` - Show only sprint-planner events
- `test-specialist` - Show only test-specialist events
- `validation-expert` - Show only validation-expert events

## Performance Features

### Real-Time Metrics

The TUI continuously monitors and displays system performance:

**Header Information** (debug mode):

```
🧙‍♂️ Morgana Protocol TUI [Split] | Focus: Dashboard | 59.2 FPS
```

**Status Bar Information:**

```
Uptime: 5m 32s | Events: 3,421 | Press H for help, Q to quit
```

**System Performance** (development mode):

- **FPS** - Rendering frame rate
- **CPU Usage** - System CPU utilization
- **Memory Usage** - Current memory consumption
- **Event Throughput** - Events processed per second
- **Queue Length** - Pending events in buffer

### Event System Stats

When `MORGANA_DEBUG=true`, additional performance information is logged:

```bash
MORGANA_DEBUG=true morgana --tui --tui-mode dev
```

Debug output includes:

```
📊 Event throughput: 1543 events/sec
🚀 Task started: task_abc123 [code-implementer]
✅ Task completed: task_abc123 [code-implementer] in 2.3s
```

## Configuration

### Environment Variables

```bash
# Enable debug logging and performance monitoring
export MORGANA_DEBUG=true

# OpenTelemetry configuration for monitoring integration
export OTEL_EXPORTER_OTLP_ENDPOINT="http://collector:4317"
export OTEL_SERVICE_NAME="morgana-tui"
export OTEL_RESOURCE_ATTRIBUTES="environment=production"
```

### Command-Line Options

```bash
# TUI-specific options
--tui                    # Enable TUI interface
--tui-mode MODE          # Set TUI mode: dev, optimized, high-performance

# Agent execution options
--parallel               # Execute tasks in parallel
--max-concurrency N      # Maximum concurrent tasks
--agent-dir PATH         # Custom agent directory
--config FILE            # Configuration file

# Telemetry options
--otel-exporter TYPE     # OpenTelemetry exporter: stdout, otlp, none
--otel-endpoint URL      # OTLP collector endpoint
```

### Configuration Files

Create a YAML configuration file for persistent settings:

```yaml
# morgana-config.yaml
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 30s

execution:
  max_concurrency: 4
  default_mode: sequential

telemetry:
  enabled: true
  exporter: otlp
  otlp_endpoint: localhost:4317
  service_name: morgana-production
```

Use with:

```bash
morgana --tui --config morgana-config.yaml
```

## Integration Examples

### Basic Task Monitoring

Monitor a single agent task:

```bash
morgana --tui --agent code-implementer --prompt "Implement user authentication system with JWT tokens and password hashing"
```

### Parallel Task Execution

Monitor multiple agents working in parallel:

```bash
morgana --tui --parallel \
  --agent code-implementer --prompt "Build REST API" \
  --agent test-specialist --prompt "Create API tests" \
  --agent validation-expert --prompt "Security review"
```

### JSON Input Processing

Process complex task definitions:

```bash
cat tasks.json | morgana --tui
```

Where `tasks.json` contains:

```json
[
  {
    "agent_type": "code-implementer",
    "prompt": "Build user authentication system",
    "options": { "complexity": "high", "language": "go" },
    "model_hint": "o3"
  },
  {
    "agent_type": "test-specialist",
    "prompt": "Generate comprehensive tests",
    "options": { "coverage": "high", "framework": "testify" }
  }
]
```

### Production Monitoring

Optimized for production environments:

```bash
# High-performance monitoring with OTLP export
MORGANA_DEBUG=false morgana --tui --tui-mode high-performance \
  --otel-exporter otlp \
  --otel-endpoint monitoring.company.com:4317 \
  --max-concurrency 8
```

## Troubleshooting

### Terminal Compatibility

**Issue**: TUI doesn't display properly

**Solutions**:

- Ensure terminal supports TUI mode: `echo $TERM`
- Use a modern terminal (iTerm2, Terminal.app, gnome-terminal)
- Check terminal size: `tput cols && tput lines`
- Verify TTY support: `test -t 1 && echo "TTY supported"`

**Unsupported terminals**:

- `TERM=dumb` (basic terminals)
- `TERM=unknown` (unrecognized terminals)
- Non-interactive shells (scripts, CI/CD)

### Performance Issues

**Issue**: High CPU usage or lag

**Solutions**:

- Switch to high-performance mode: `--tui-mode high-performance`
- Reduce event volume: `--max-concurrency 2`
- Disable debug features: `MORGANA_DEBUG=false`
- Check system resources: `top` or `htop`

**Issue**: Memory usage growth

**Solutions**:

- Restart for long-running sessions
- Reduce log buffer: Use high-performance mode
- Monitor memory: Watch for memory leaks

### Display Issues

**Issue**: Colors not showing

**Solutions**:

- Check color support: `tput colors`
- Use terminal with 256-color support
- Set TERM properly: `export TERM=xterm-256color`

**Issue**: Layout breaks on resize

**Solutions**:

- Terminal will automatically adjust
- Try toggling layout mode: Press `L` or `Space`
- Restart if layout remains broken

### Event System Issues

**Issue**: Missing events or updates

**Solutions**:

- Check event bus connection
- Increase buffer size (development mode)
- Enable debug mode: `MORGANA_DEBUG=true`
- Verify no competing processes

## Advanced Usage

### Custom Agent Integration

The TUI automatically detects and displays any Morgana Protocol agent:

```bash
# Custom agent example
morgana --tui --agent my-custom-agent --prompt "Custom task"
```

The TUI will:

- Display the agent in the dashboard
- Show real-time progress and status
- Log all execution phases
- Provide performance metrics

### Integration with Monitoring Systems

For production monitoring, integrate with your existing observability stack:

```bash
# With Jaeger tracing
export OTEL_EXPORTER_JAEGER_ENDPOINT="http://jaeger:14268/api/traces"
morgana --tui --otel-exporter jaeger

# With Prometheus metrics
export OTEL_EXPORTER_PROMETHEUS_PORT=9090
morgana --tui --otel-exporter prometheus

# With custom OTLP collector
export OTEL_EXPORTER_OTLP_ENDPOINT="http://otel-collector:4317"
morgana --tui --otel-exporter otlp
```

### Session Management

For long-running monitoring sessions:

```bash
# Start in background with logging
nohup morgana --tui --tui-mode high-performance > morgana.log 2>&1 &

# Monitor with screen/tmux
screen -S morgana-tui
morgana --tui --tui-mode optimized
# Ctrl+A, D to detach

# Reconnect later
screen -r morgana-tui
```

## API Integration

While the TUI is primarily a terminal interface, it integrates with the Morgana
Protocol event system, allowing programmatic access to all displayed information
through the event bus.

## Support

For issues, questions, or feature requests:

1. Check this user guide first
2. Review the main [MONITORING_INTEGRATION.md](MONITORING_INTEGRATION.md) for
   technical details
3. Enable debug mode and check logs: `MORGANA_DEBUG=true`
4. Report issues with:
   - Terminal type and version
   - Operating system
   - Morgana Protocol version
   - Steps to reproduce

## Changelog

### v1.0.0 - Initial TUI Release

- Real-time agent monitoring
- Multiple layout modes
- Event-driven architecture
- Performance optimization modes
- OpenTelemetry integration
- Comprehensive keyboard controls
- Terminal compatibility detection

---

_This guide covers the Morgana Protocol TUI interface. For complete system
documentation, see the main repository README and integration guides._
