# Morgana Protocol Integration Tests

This directory contains comprehensive integration tests for the Morgana
Protocol's simplified event stream architecture.

## Overview

The integration tests validate the complete event-driven system:

- **Event Stream Processing**: Event bus, publishers, subscribers
- **TUI Integration**: Terminal User Interface event display
- **Monitor System**: Event forwarding and monitoring
- **Command Polling**: Pause/resume/stop functionality
- **Performance Requirements**: <10ms latency, high throughput
- **Resource Management**: Memory leaks, file rotation
- **Error Handling**: Recovery, backpressure, reliability

## Test Structure

### Core Test Files

- **`test_utils.go`** - Test utilities, event generators, collectors, setup
  functions
- **`event_stream_test.go`** - Event stream reading/writing, performance, error
  handling
- **`tui_integration_test.go`** - TUI display, statistics tracking,
  configuration
- **`monitor_integration_test.go`** - Monitor client/server, event forwarding
- **`command_polling_test.go`** - Command file polling, pause/resume/stop
- **`resource_leak_test.go`** - Memory leaks, goroutine leaks, file rotation
- **`integration_suite_test.go`** - Complete system integration scenarios

### Test Categories

#### 1. Event Stream Tests

- **Basic Flow**: Event writing → processing → reading
- **Event Types**: TaskStarted, TaskProgress, TaskCompleted, TaskFailed
- **High Volume**: Throughput and latency validation
- **Error Recovery**: Subscriber panics, backpressure handling

#### 2. TUI Integration Tests

- **Display Validation**: Events appear in TUI correctly
- **Performance**: FPS, memory usage under load
- **Configuration**: Different TUI configs (optimized, high-performance)
- **Multi-Instance**: TUI manager with multiple instances

#### 3. Monitor Integration Tests

- **Client/Server**: IPC communication over Unix sockets
- **Multiple Clients**: Event forwarding from many sources
- **Reconnection**: Client reconnection after server restart
- **Event Filtering**: Different event types handled correctly

#### 4. Command Polling Tests

- **Basic Commands**: pause, resume, stop via file polling
- **Performance Impact**: Polling doesn't degrade event processing
- **File Cleanup**: Command files removed after processing
- **Error Handling**: Invalid commands handled gracefully

#### 5. Resource Leak Tests

- **Memory Monitoring**: Memory usage over time
- **Goroutine Tracking**: Goroutine leak detection
- **File Descriptors**: FD leak detection (Unix systems)
- **Long-Running**: Extended operation stability

#### 6. File Rotation Tests

- **Log Rotation**: Simulated log file rotation scenarios
- **Multiple Rotations**: Rapid successive rotations
- **Error Scenarios**: Failed rotation handling

## Performance Requirements Tested

1. **Latency**: Event processing <10ms average
2. **Throughput**: >1000 events/sec minimum
3. **Memory**: <200MB reasonable usage limit
4. **Event Loss**: <1% dropped events under normal load
5. **Resource Leaks**: No memory or goroutine leaks over time

## Running the Tests

### Prerequisites

```bash
# Ensure Go modules are up to date
go mod tidy

# Install dependencies if needed
go mod download
```

### Individual Test Suites

```bash
# Event stream tests
go test -v -tags=integration ./tests/integration -run TestEventStream

# TUI integration tests
go test -v -tags=integration ./tests/integration -run TestTUI

# Monitor integration tests
go test -v -tags=integration ./tests/integration -run TestMonitor

# Command polling tests
go test -v -tags=integration ./tests/integration -run TestCommand

# Resource leak tests
go test -v -tags=integration ./tests/integration -run TestResource

# File rotation tests
go test -v -tags=integration ./tests/integration -run TestFileRotation

# Complete integration suite
go test -v -tags=integration ./tests/integration -run TestMorganaIntegrationSuite
```

### All Integration Tests

```bash
# Run all integration tests
go test -v -tags=integration ./tests/integration

# Run with verbose output and no timeout
go test -v -tags=integration -timeout=30m ./tests/integration

# Run specific performance tests only
go test -v -tags=integration ./tests/integration -run Performance
```

### Makefile Targets

If using the project's Makefile:

```bash
# Run integration tests
make test-integration

# Run all tests (unit + integration)
make test-all

# Run with coverage
make test-integration-coverage
```

## Test Environment

### Terminal Support

TUI tests automatically skip if running in a non-terminal environment:

- CI/CD systems without TTY support
- Headless servers
- Screen/tmux sessions (depending on configuration)

### Platform Differences

Some tests adapt based on the operating system:

- **File Descriptor tests**: Unix/Linux only (skipped on Windows)
- **Signal handling**: Platform-specific behavior
- **Path separators**: Cross-platform file path handling

### Resource Requirements

Integration tests require:

- **Memory**: At least 512MB available
- **File Descriptors**: Several hundred for socket tests
- **Disk Space**: ~100MB for temporary files during rotation tests
- **Time**: Up to 30 minutes for complete suite

## Test Configuration

### Environment Variables

```bash
# Enable debug output in event system
export MORGANA_DEBUG=true

# Increase test timeout for slow systems
export TEST_TIMEOUT=30m

# Run tests with race detector
export GORACE="halt_on_error=1"
```

### Test Flags

```bash
# Verbose output
go test -v ...

# Enable race detector
go test -race ...

# Set timeout
go test -timeout=30m ...

# Run specific test pattern
go test -run="TestEventStream.*Performance" ...
```

## Interpreting Results

### Performance Metrics

Tests log detailed performance metrics:

```
Event latency - Max: 2.3ms, Avg: 1.1ms
Event throughput: 15420 events/sec
TUI Stats: Events=5000, Renders=120, FPS=24.3
Resource analysis: Mem=12.3MB, GR=0, Events=10000, GC=8
```

### Common Issues

1. **Timeout Failures**: Increase timeout for slow systems
2. **TUI Skipped**: Expected on headless systems
3. **FD Limit**: Increase ulimit if FD tests fail
4. **Memory Variance**: GC timing affects memory measurements

### Success Criteria

✅ All tests pass without errors ✅ Performance requirements met ✅ No resource
leaks detected ✅ System remains stable under load ✅ Error conditions handled
gracefully

## Extending the Tests

### Adding New Test Cases

1. Create test function in appropriate `*_test.go` file
2. Use `SetupIntegrationTest(t)` for common setup
3. Follow naming convention: `TestFeatureName`
4. Include performance and error scenarios
5. Use `t.Logf()` for informative output

### Test Utilities

The `test_utils.go` file provides:

- **`TestEventGenerator`**: Generate various event patterns
- **`EventCollector`**: Collect and analyze events
- **`ResourceMonitor`**: Track memory/goroutine usage
- **`PerformanceAssertions`**: Validate performance requirements
- **Helper functions**: Wait for events, setup, cleanup

### Example Test

```go
func TestNewFeature(t *testing.T) {
    setup := SetupIntegrationTest(t)
    defer setup.Cleanup()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    t.Run("BasicOperation", func(t *testing.T) {
        // Generate test events
        taskID := setup.Generator.GenerateTaskLifecycle(
            ctx, "test-agent", 10*time.Millisecond)

        // Wait for processing
        if !WaitForEvents(setup.Collector, 5, time.Second) {
            t.Fatal("Events not processed")
        }

        // Validate results
        events := setup.Collector.GetEventsForTask(taskID)
        if len(events) != 5 {
            t.Errorf("Expected 5 events, got %d", len(events))
        }
    })
}
```

## Continuous Integration

### GitHub Actions

```yaml
- name: Run Integration Tests
  run: |
    go test -v -tags=integration -timeout=20m ./tests/integration
```

### Test Reports

Consider generating test reports:

```bash
# JSON output for CI analysis
go test -tags=integration -json ./tests/integration > test-results.json

# Coverage report
go test -tags=integration -coverprofile=integration.out ./tests/integration
go tool cover -html=integration.out -o integration-coverage.html
```

## Support

For issues with integration tests:

1. Check test logs for specific failure details
2. Verify system requirements (memory, FDs, etc.)
3. Run individual test suites to isolate problems
4. Enable debug output with `MORGANA_DEBUG=true`
5. Increase timeouts for slower systems

The integration tests validate the complete Morgana Protocol architecture and
ensure the simplified event stream design meets all performance and reliability
requirements.
