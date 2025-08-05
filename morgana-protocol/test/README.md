# Morgana Protocol Integration Tests

This directory contains integration tests for the Morgana Protocol, focusing on
the Task tool integration and multi-process architecture.

## Test Structure

### Task Client Tests (`pkg/task/client_integration_test.go`)

Tests the Go client that interfaces with the Python bridge:

- Successful Python bridge execution
- Error handling and propagation
- Timeout and context cancellation
- Bridge script discovery
- Mock mode functionality
- JSON marshaling/unmarshaling

### Adapter Tests (`internal/adapter/adapter_integration_test.go`)

Tests the adapter that orchestrates agent execution:

- Agent type validation
- Agent prompt loading
- Timeout configuration (default and per-agent)
- Concurrent execution
- Telemetry integration
- Error propagation from task client

### End-to-End Tests (`cmd/morgana/main_integration_test.go`)

Tests the complete system from CLI to output:

- Command-line argument parsing
- JSON input via stdin
- Parallel vs sequential execution
- Configuration file loading
- Timeout enforcement
- Version command

## Running Tests

### Unit Tests Only (Fast)

```bash
make test
```

### Integration Tests Only

```bash
make test-integration
```

### All Tests

```bash
make test-all
```

### Individual Test Suites

```bash
# Task client tests
go test -v -tags=integration ./pkg/task

# Adapter tests
go test -v -tags=integration ./internal/adapter

# E2E tests
go test -v -tags=integration ./cmd/morgana
```

## Test Fixtures

### Python Bridge Test Script

`pkg/task/testdata/test_bridge.py` - Simulates various bridge scenarios:

- `TEST_MODE=success` - Normal successful execution
- `TEST_MODE=error` - Returns error response
- `TEST_MODE=timeout` - Sleeps to trigger timeout
- `TEST_MODE=invalid_json` - Returns malformed JSON
- `TEST_MODE=crash` - Simulates Python crash

### Test Agent Prompts

`internal/adapter/testdata/test-agent.md` - Simple test agent prompt

## Testing the Python Bridge

A manual test script is provided:

```bash
./scripts/test_bridge.sh
```

This tests the Python bridge directly with various inputs.

## Integration Test Scenarios

### 1. Multi-Language Bridge

- Go → Python subprocess execution
- JSON communication between processes
- Error propagation across language boundary
- Signal handling (SIGTERM, SIGINT)

### 2. Context and Timeouts

- Context deadline propagation to subprocess
- Per-agent timeout configuration
- Graceful cancellation handling
- Timeout error reporting

### 3. Concurrent Execution

- Multiple agents running in parallel
- Resource limits (semaphore)
- Race condition prevention
- Result ordering

### 4. Configuration

- YAML file parsing
- Environment variable overrides
- CLI flag precedence
- Dynamic agent discovery

## Known Limitations

1. **Python Dependency**: Tests require Python 3 to be installed
2. **Platform Differences**: Signal handling may vary between OS
3. **Timing Sensitivity**: Timeout tests may be flaky on slow systems
4. **File System**: Tests create temporary files and directories

## Future Improvements

1. **Mock MCP Server**: Test with actual Task() function simulation
2. **Performance Benchmarks**: Add benchmarks for subprocess overhead
3. **Stress Testing**: Test with hundreds of concurrent agents
4. **Integration with CI/CD**: GitHub Actions integration test workflow
5. **Coverage Analysis**: Measure integration test coverage
