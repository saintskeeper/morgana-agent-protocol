# Morgana Protocol - Monitoring Integration

This document describes the complete integration of the event-driven monitoring
system with Morgana orchestration, connecting the TUI renderer with the event
system for real-time monitoring.

## Overview

The monitoring integration provides real-time visibility into agent execution
through a sophisticated event-driven architecture that delivers <100ms latency
from agent action to TUI update while maintaining zero performance impact on
core agent execution.

## Architecture

### Core Components

1. **Event Bus** - High-performance pub/sub system

   - 5M+ events/sec capability
   - Configurable buffer size and worker pool
   - Async publishing with drop protection
   - Comprehensive statistics tracking

2. **TUI Renderer** - Real-time interface

   - Bubbletea-based responsive UI
   - Multiple layout modes (split, dashboard-only, logs-only)
   - Configurable refresh rates (dev, optimized, high-performance)
   - Real-time performance metrics display

3. **Event Bridge** - TUI-Event integration

   - Converts events to bubbletea messages
   - FPS control and performance monitoring
   - Event filtering and processing
   - Statistics collection

4. **Orchestrator Integration** - Task management
   - Publishes lifecycle events
   - Progress tracking
   - Error propagation
   - Performance metrics

## Integration Flow

```
Agent Execution → Adapter Events → Event Bus → Event Bridge → TUI Updates
                             ↓
                    Orchestrator Events → Event Bus → Performance Monitoring
                             ↓
                      OpenTelemetry → OTEL Collector → Monitoring Stack
```

## Event Types

### Task Lifecycle Events

- `task.started` - Task begins execution
- `task.progress` - Progress updates during execution
- `task.completed` - Task completes successfully
- `task.failed` - Task execution fails

### Orchestrator Events

- `orchestrator.started` - Orchestration begins
- `orchestrator.completed` - Orchestration completes
- `orchestrator.failed` - Orchestration fails

### Adapter Events

- `adapter.validation` - Agent validation
- `adapter.prompt_load` - Prompt loading
- `adapter.execution` - Execution phases

## Usage

### Command Line Interface

#### Basic TUI Mode

```bash
# Enable TUI with optimized configuration
morgana --tui --tui-mode optimized

# Run with demo tasks
morgana --tui --tui-mode dev

# High-performance mode for production
morgana --tui --tui-mode high-performance
```

#### Task Execution with Monitoring

```bash
# Sequential execution with TUI
morgana --tui --agent code-implementer --prompt "Implement authentication"

# Parallel execution with monitoring
morgana --tui --parallel --agent code-implementer --prompt "Task 1" --agent test-specialist --prompt "Task 2"

# JSON input with TUI
echo '[{"agent_type":"code-implementer","prompt":"Build API"}]' | morgana --tui
```

### Configuration

#### TUI Modes

**Development Mode** (`--tui-mode dev`)

- Show debug information
- Enable all features (filtering, search, export)
- Higher refresh rate for responsive development

**Optimized Mode** (`--tui-mode optimized`) - Default

- 60 FPS refresh rate
- Balanced performance and features
- Suitable for most use cases

**High-Performance Mode** (`--tui-mode high-performance`)

- 30 FPS for lower CPU usage
- Larger event buffers
- Disabled expensive features
- Optimized for production

#### Event Bus Configuration

```yaml
# Environment variables
MORGANA_DEBUG=true # Enable debug logging and performance monitoring
```

## Performance Characteristics

### Latency Targets

- **Event Publishing**: <1ms
- **Event Processing**: <10ms
- **TUI Update**: <100ms total latency
- **Agent Execution**: Zero performance impact

### Throughput Capabilities

- **Event Bus**: 5M+ events/sec theoretical
- **TUI Processing**: 1000+ events/sec sustained
- **Concurrent Tasks**: Up to configured max concurrency
- **Parallel Execution**: Linear scaling with concurrency

### Resource Usage

- **Memory**: ~10-50MB depending on log buffer size
- **CPU**: <5% overhead for monitoring
- **Network**: Zero impact (local event bus)
- **Disk**: Only for OTEL telemetry export

## TUI Features

### Real-time Dashboard

- **Agent Status Cards**: Live progress and status
- **Task Queue**: Current and completed tasks
- **Performance Metrics**: CPU, memory, event throughput
- **System Health**: Error rates, success rates

### Interactive Log Viewer

- **Real-time Streaming**: Live log updates
- **Filtering**: By agent type, log level, keywords
- **Search**: Interactive log search
- **Navigation**: Scroll, jump to top/bottom, page up/down

### Keyboard Shortcuts

- `Tab/Shift+Tab` - Switch between components
- `Space/L` - Toggle layout mode
- `H/?` - Show/hide help
- `Q/Ctrl+C/Esc` - Quit application
- `↑/↓` or `K/J` - Scroll logs
- `F` - Cycle through filters
- `C` - Clear current filter

## Integration Points

### Adapter Integration

```go
// Connect adapter to event bus
adapter.SetEventBus(eventBus)

// Events are automatically published during execution
result := adapter.Execute(ctx, task)
```

### Orchestrator Integration

```go
// Connect orchestrator to event bus
orchestrator.SetEventBus(eventBus)

// Events published during parallel/sequential execution
results := orchestrator.RunParallel(ctx, tasks)
```

### TUI Integration

```go
// Start TUI with event bus connection
tui, err := tui.RunAsync(ctx, eventBus, tuiConfig)
if err != nil {
    log.Fatal(err)
}
defer tui.Stop()
```

## Error Handling

### Event System Resilience

- **Async Publishing**: Non-blocking event publishing
- **Buffer Overflow**: Graceful event dropping with statistics
- **Subscriber Panics**: Automatic recovery and logging
- **Bus Shutdown**: Graceful cleanup of resources

### TUI Error Recovery

- **Terminal Resize**: Automatic layout adaptation
- **Event Processing Errors**: Graceful degradation
- **Component Failures**: Isolated failure handling
- **Rendering Errors**: Continue with reduced functionality

## Testing

### Integration Tests

```bash
# Run monitoring integration tests
go test -tags=integration ./cmd/morgana -run TestMonitoringIntegration

# Performance benchmarks
go test -tags=integration ./cmd/morgana -run TestEventPerformance

# Concurrent execution tests
go test -tags=integration ./cmd/morgana -run TestConcurrentExecution

# TUI integration tests
go test -tags=integration ./cmd/morgana -run TestTUIIntegration
```

### Performance Validation

- Event throughput >= 1000 events/sec
- Drop rate < 1% under normal load
- TUI responsiveness maintained under high event volume
- Memory usage stable over extended runs

## Monitoring and Observability

### Built-in Metrics

- **Event Bus Statistics**: Published, dropped, queue size
- **TUI Performance**: FPS, render count, memory usage
- **Task Execution**: Success rate, duration, throughput
- **System Resources**: CPU, memory usage

### OpenTelemetry Integration

- **Spans**: Task execution, orchestration phases
- **Metrics**: Performance counters, resource usage
- **Traces**: End-to-end request tracing
- **Export**: Jaeger, Prometheus, OTLP collectors

### Debug Information

```bash
# Enable debug logging
MORGANA_DEBUG=true morgana --tui --tui-mode dev

# Performance monitoring output:
# 📊 Event throughput: 1543 events/sec
# 🚀 Task started: task_abc123 [code-implementer]
# ✅ Task completed: task_abc123 [code-implementer] in 2.3s
```

## Production Deployment

### Recommended Configuration

```bash
# Production TUI usage
morgana --tui --tui-mode high-performance \
        --max-concurrency 8 \
        --otel-exporter otlp \
        --otel-endpoint collector:4317
```

### Environment Variables

```bash
export MORGANA_DEBUG=false           # Disable debug logging
export OTEL_EXPORTER_OTLP_ENDPOINT="http://collector:4317"
export OTEL_SERVICE_NAME="morgana-production"
export OTEL_RESOURCE_ATTRIBUTES="environment=production"
```

### Monitoring Stack Integration

- **Jaeger**: Distributed tracing
- **Prometheus**: Metrics collection
- **Grafana**: Visualization dashboards
- **Alertmanager**: Error and performance alerts

## Troubleshooting

### Common Issues

**High Event Drop Rate**

- Increase event bus buffer size
- Add more worker goroutines
- Check subscriber performance

**TUI Performance Issues**

- Switch to high-performance mode
- Reduce refresh rate
- Disable expensive features

**Memory Usage Growth**

- Reduce log buffer size
- Check for event processing bottlenecks
- Monitor for memory leaks in subscribers

**Missing Events**

- Verify event bus connections
- Check subscriber error logs
- Validate event publishing code

### Debug Commands

```bash
# Check TUI terminal support
morgana --tui --version

# Validate configuration
morgana --config config.yaml --version

# Test event system performance
go test -tags=integration -run TestEventPerformance -v
```

## Future Enhancements

### Planned Features

- **Web UI**: Browser-based monitoring interface
- **Remote Monitoring**: Multi-instance orchestrator tracking
- **Advanced Filtering**: Complex query language for logs
- **Export Functions**: Save logs and metrics to files
- **Plugin System**: Custom TUI components and event handlers

### Performance Improvements

- **Event Batching**: Reduce syscall overhead
- **Memory Pooling**: Reduce GC pressure
- **Compression**: Compress event payloads
- **Persistence**: Event replay and historical analysis

## Contributing

See the main repository README for contribution guidelines. The monitoring
integration follows the same patterns:

1. **Event-Driven**: All interactions through event bus
2. **Zero Impact**: No performance degradation to core functionality
3. **Graceful Degradation**: Continue working if monitoring fails
4. **Comprehensive Testing**: Integration and performance tests required
5. **Documentation**: Update this document for any changes

## References

- [Event System Documentation](docs/EVENT_SYSTEM.md)
- [TUI Implementation Guide](docs/TUI_IMPLEMENTATION.md)
- [Integration Examples](examples/)
- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/go/)
- [Bubbletea Framework](https://github.com/charmbracelet/bubbletea)
