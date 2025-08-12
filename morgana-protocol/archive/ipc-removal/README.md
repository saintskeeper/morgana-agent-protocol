# IPC Bridge Removal Archive

This directory contains the archived files from the Morgana Protocol's
file-based IPC (Inter-Process Communication) bridge system that was removed to
simplify the codebase.

## Date Removed

August 11, 2025

## Rationale for Removal

The complex file-based IPC mechanism has been replaced with a simple event
stream logger for the following reasons:

1. **Simplification**: The IPC bridge added significant complexity (474+ lines
   of code) that was difficult to maintain and debug.

2. **Event Stream Integration**: Morgana now uses a simple event stream logger
   that provides better real-time monitoring without the overhead of file-based
   communication.

3. **Reduced Dependencies**: Removing the Python bridge dependencies simplifies
   the build and deployment process.

4. **Better Maintainability**: The new approach is easier to understand, test,
   and extend.

## Files Archived

### Python Scripts (`scripts/`)

- `task_bridge.py` - Main task bridge implementation
- `task_bridge_claude.py` - Claude-specific bridge adapter
- `task_bridge_real.py` - Real implementation bridge
- `morgana_repl.py` - REPL adapter for Morgana
- `morgana_subagent_bridge.py` - Sub-agent bridge functionality

### Go Packages (`pkg/task/`)

- `client.go` - Task client implementation
- `repl_adapter.go` - REPL adapter for Go
- `client_integration_test.go` - Integration tests
- `testdata/` - Test data directory

## Configuration Changes

The following configuration sections were removed from `morgana.yaml`:

```yaml
task_client:
  bridge_path: ~/.claude/scripts/task_bridge.py
  python_path: python3
  mock_mode: false
  timeout: 5m
```

And from `internal/config/config.go`, the `TaskClientConfig` struct and related
functionality.

## Code Changes

### Main Application (`cmd/morgana/main.go`)

- Removed task client import
- Removed task client initialization
- Updated adapter creation to pass `nil` instead of task client

### Adapter (`internal/adapter/adapter.go`)

- Removed task client dependency
- Replaced `taskClient.RunWithContext()` calls with stub responses
- Updated constructor to accept interface{} instead of concrete task client

## Migration Path

The task execution functionality now returns informational messages indicating
that the IPC bridge has been removed. This maintains compatibility with existing
orchestration code while simplifying the architecture.

## Future Improvements

With the IPC bridge removed, future enhancements can focus on:

1. Enhanced event stream logging
2. Direct integration with Claude Code APIs
3. Simplified agent execution models
4. Better real-time monitoring capabilities

## Restoration

If these files need to be restored for any reason, they can be found in this
archive directory and the git history contains the full implementation details.
