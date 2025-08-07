# Sprint Plan: Fix Parallel Task Logging Visibility

## Sprint Overview

- **Sprint ID**: sprint-2025-01-07-logging-visibility
- **Duration**: 8-12 hours (1.5 days)
- **Goal**: Enable visible logging output when running morgana parallel tasks
  through make attach
- **Success Criteria**:
  - Logs from parallel morgana tasks appear in attached TUI session
  - Both headless and TUI modes show proper logging
  - Event forwarding works correctly between morgana CLI and monitor daemon

## Problem Statement

When using `make up` to start morgana-monitor in headless mode and then
`make attach` to view the TUI, no logs appear when morgana is running parallel
tasks. The monitor starts in headless mode but doesn't properly display
forwarded events from parallel task execution.

## Task Definitions

### Task 1: LOGGING_ANALYSIS

- **Priority**: P0 (Critical)
- **Type**: Investigation
- **Duration**: 2 hours
- **Dependencies**: None
- **Description**: Deep dive into current logging pipeline
- **Exit Criteria**:
  - Document complete event flow from morgana CLI → monitor daemon
  - Identify where logs are being lost in headless → TUI transition
  - Map all logging configurations in morgana-monitor headless vs TUI modes
  - Verify IPC socket communication between components

### Task 2: HEADLESS_FIX

- **Priority**: P0 (Critical)
- **Type**: Implementation
- **Duration**: 3 hours
- **Dependencies**: LOGGING_ANALYSIS
- **Description**: Fix headless mode event forwarding and logging
- **Exit Criteria**:
  - Headless mode properly receives and buffers events from IPC socket
  - Event buffer persists when TUI client attaches
  - Logs written to /tmp/morgana-monitor.log in headless mode
  - No events lost during mode transitions

### Task 3: TUI_CLIENT_IMPL

- **Priority**: P0 (Critical)
- **Type**: Implementation
- **Duration**: 2 hours
- **Dependencies**: HEADLESS_FIX
- **Description**: Implement proper TUI client mode for attach command
- **Exit Criteria**:
  - `morgana-monitor --client` properly connects to running daemon
  - Client receives buffered events from headless mode
  - Real-time event streaming works after attachment
  - Clean disconnection without affecting daemon

### Task 4: EVENT_BUFFER_IMPL

- **Priority**: P1 (High)
- **Type**: Implementation
- **Duration**: 2 hours
- **Dependencies**: HEADLESS_FIX
- **Description**: Implement event buffering for late-joining TUI clients
- **Exit Criteria**:
  - Circular buffer stores last 1000 events in daemon
  - New TUI clients receive buffered event history on connect
  - Buffer size configurable via environment variable
  - Memory-efficient implementation with proper cleanup

### Task 5: PARALLEL_TEST

- **Priority**: P1 (High)
- **Type**: Testing
- **Duration**: 1 hour
- **Dependencies**: TUI_CLIENT_IMPL, EVENT_BUFFER_IMPL
- **Description**: Create comprehensive test suite for parallel logging
- **Exit Criteria**:
  - Test script that runs parallel morgana tasks
  - Validates logs appear in attached TUI
  - Tests mode transitions (headless → TUI → headless)
  - Stress test with 10+ parallel tasks

### Task 6: CTL_SCRIPT_UPDATE

- **Priority**: P2 (Medium)
- **Type**: Implementation
- **Duration**: 1 hour
- **Dependencies**: TUI_CLIENT_IMPL
- **Description**: Update morgana-monitor-ctl.sh for better attach support
- **Exit Criteria**:
  - Script detects if daemon is in headless mode
  - Provides clear instructions for attaching
  - Handles screen/tmux sessions properly
  - Better error messages for connection failures

### Task 7: MAKEFILE_ENHANCE

- **Priority**: P2 (Medium)
- **Type**: Implementation
- **Duration**: 0.5 hours
- **Dependencies**: CTL_SCRIPT_UPDATE
- **Description**: Enhance Makefile attach target
- **Exit Criteria**:
  - `make attach` automatically uses correct client mode
  - Shows buffered events immediately on attach
  - Provides fallback to log tailing if TUI unavailable
  - Clear status messages during attachment

### Task 8: DOCS_UPDATE

- **Priority**: P3 (Low)
- **Type**: Documentation
- **Duration**: 0.5 hours
- **Dependencies**: All implementation tasks
- **Description**: Update documentation for new logging behavior
- **Exit Criteria**:
  - README updated with headless/TUI mode explanation
  - Troubleshooting section for logging issues
  - Examples of viewing parallel task logs
  - Architecture diagram updated if needed

## Technical Approach

### Root Cause Analysis

The issue appears to be that morgana-monitor starts in headless mode (line 113
in morgana-monitor-ctl.sh) but doesn't properly forward events to the TUI when
attached later. The headless mode likely isn't:

1. Listening on the IPC socket for events
2. Buffering events for later TUI attachment
3. Providing a proper client mode for attachment

### Solution Architecture

1. **Dual-mode daemon**: Support both headless (event collection) and TUI
   (display) modes
2. **Event buffering**: Store recent events for late-joining clients
3. **Client mode**: Separate mode for attaching to running daemon
4. **IPC enhancement**: Bidirectional communication for event replay

## Parallel Execution Opportunities

The following tasks can run in parallel after dependencies are met:

- EVENT_BUFFER_IMPL and PARALLEL_TEST preparation
- CTL_SCRIPT_UPDATE and MAKEFILE_ENHANCE
- All documentation can proceed alongside testing

## Risk Mitigation

**Risk 1**: Breaking existing TUI functionality

- **Mitigation**: Maintain backward compatibility, test both modes extensively

**Risk 2**: Memory leak from event buffering

- **Mitigation**: Implement circular buffer with size limits, proper cleanup

**Risk 3**: Race conditions in IPC communication

- **Mitigation**: Use proper synchronization, mutex locks for shared resources

## Definition of Done

- [ ] Parallel task logs visible in attached TUI session
- [ ] No events lost during mode transitions
- [ ] All tests passing with parallel execution
- [ ] Documentation updated
- [ ] Code reviewed and merged
- [ ] Verified on macOS with screen/tmux

## Validation Commands

```bash
# Start monitor in headless mode
make up

# Run parallel tasks
echo '[{"agent_type":"code-implementer","prompt":"Test 1"},{"agent_type":"test-specialist","prompt":"Test 2"}]' | morgana --parallel

# Attach and verify logs appear
make attach

# Check log file for completeness
tail -f /tmp/morgana-monitor.log
```

## Team Allocation Suggestion

- **Backend Developer**: Tasks 1-4 (core logging fixes)
- **DevOps Engineer**: Tasks 5-7 (testing and scripts)
- **Technical Writer**: Task 8 (documentation)

## Notes

This sprint focuses on fixing a critical observability issue that impacts the
developer experience when using Morgana Protocol's parallel execution features.
The fix will significantly improve debugging and monitoring capabilities.
