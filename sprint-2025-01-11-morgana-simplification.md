# Sprint Plan: Morgana Protocol Simplification

## Sprint Overview

- **Sprint ID**: sprint-2025-01-11
- **Duration**: 5 days (40 hours)
- **Goal**: Transform Morgana from complex orchestrator to simple event
  monitoring layer
- **Success Criteria**:
  - Remove file-based IPC complexity
  - Implement simple event stream monitoring
  - Maintain TUI visualization
  - Enable native Claude Code execution

## Background

Based on consensus analysis from O3-mini, Gemini-2.5-pro, and GPT-4.1:

- Current file-based IPC is overly complex and creates technical debt
- Morgana should be a pure monitoring/visualization layer
- Execution should remain within Claude Code's native environment
- Simple event stream pattern aligns with industry best practices

## Task Definitions

### Phase 1: Core Refactoring (Day 1-2)

#### MORGANA_STREAM_IMPL

- **Priority**: P0-Critical
- **Type**: Implementation
- **Complexity**: Medium (8 hours)
- **Dependencies**: None
- **Description**: Create simple event stream writer for Claude Code
- **Exit Criteria**:
  - Event logger function writes JSON lines to `/tmp/morgana/events.jsonl`
  - Supports agent_start, agent_complete, agent_error event types
  - Includes timestamp, agent_type, and duration fields
  - No polling or IPC complexity
- **Files to Create**:
  - `/Users/walterday/.claude/scripts/morgana_events.py`

#### MORGANA_CONSUMER_IMPL

- **Priority**: P0-Critical
- **Type**: Implementation
- **Complexity**: Medium (8 hours)
- **Dependencies**: MORGANA_STREAM_IMPL
- **Description**: Refactor Go monitor to consume event stream
- **Exit Criteria**:
  - Monitor tails `/tmp/morgana/events.jsonl` file
  - Publishes events to existing EventBus
  - No bidirectional communication
  - Maintains existing TUI display functionality
- **Files to Modify**:
  - `/Users/walterday/.claude/morgana-protocol/internal/monitor/consumer.go`

### Phase 2: Remove Complexity (Day 2-3)

#### REMOVE_IPC_TASK

- **Priority**: P0-Critical
- **Type**: Refactoring
- **Complexity**: Medium (6 hours)
- **Dependencies**: MORGANA_CONSUMER_IMPL
- **Description**: Remove file-based IPC and task bridge code
- **Exit Criteria**:
  - Delete `task_bridge.py` and related files
  - Remove `pkg/task/client.go` subprocess execution
  - Clean up `repl_adapter.go` and polling logic
  - Update config to remove bridge_path settings
- **Files to Remove**:
  - `/Users/walterday/.claude/scripts/task_bridge*.py`
  - `/Users/walterday/.claude/morgana-protocol/pkg/task/repl_adapter.go`

#### CLAUDE_NATIVE_IMPL

- **Priority**: P1-High
- **Type**: Implementation
- **Complexity**: Simple (4 hours)
- **Dependencies**: MORGANA_STREAM_IMPL
- **Description**: Create native Claude Code agent executor
- **Exit Criteria**:
  - Simple function that executes Task tool directly
  - Logs events before and after execution
  - No subprocess or file-based communication
  - Works within Claude Code REPL
- **Files to Create**:
  - `/Users/walterday/.claude/scripts/claude_agent_executor.py`

### Phase 3: Optional Enhancements (Day 3-4)

#### COMMAND_POLL_IMPL

- **Priority**: P2-Medium
- **Type**: Implementation
- **Complexity**: Simple (4 hours)
- **Dependencies**: CLAUDE_NATIVE_IMPL
- **Description**: Add lightweight command polling for pause/resume
- **Exit Criteria**:
  - Claude agents check `/tmp/morgana/commands.txt` periodically
  - Supports pause, resume, and stop commands
  - Single file, no complex IPC
  - Optional feature, not required for core functionality

#### MONITOR_ENHANCE_IMPL

- **Priority**: P2-Medium
- **Type**: Enhancement
- **Complexity**: Medium (6 hours)
- **Dependencies**: MORGANA_CONSUMER_IMPL
- **Description**: Enhance TUI with execution statistics
- **Exit Criteria**:
  - Display agent execution count
  - Show average duration per agent type
  - Add simple performance metrics
  - No new communication channels

### Phase 4: Testing & Documentation (Day 4-5)

#### INTEGRATION_TEST

- **Priority**: P1-High
- **Type**: Testing
- **Complexity**: Medium (6 hours)
- **Dependencies**: All implementation tasks
- **Description**: Create integration tests for new architecture
- **Exit Criteria**:
  - Test event stream writing and reading
  - Verify TUI displays events correctly
  - Ensure no resource leaks or polling issues
  - Document test results

#### DOCS_UPDATE

- **Priority**: P1-High
- **Type**: Documentation
- **Complexity**: Simple (4 hours)
- **Dependencies**: All implementation tasks
- **Description**: Update documentation for simplified architecture
- **Exit Criteria**:
  - README reflects new event-stream architecture
  - Remove references to file-based IPC
  - Add simple usage examples
  - Update architecture diagrams

#### MIGRATION_GUIDE

- **Priority**: P2-Medium
- **Type**: Documentation
- **Complexity**: Simple (2 hours)
- **Dependencies**: DOCS_UPDATE
- **Description**: Create migration guide for existing users
- **Exit Criteria**:
  - Step-by-step migration instructions
  - Backward compatibility notes
  - Configuration changes documented

## Critical Path

```
MORGANA_STREAM_IMPL → MORGANA_CONSUMER_IMPL → REMOVE_IPC_TASK
                  ↓
           CLAUDE_NATIVE_IMPL → INTEGRATION_TEST → DOCS_UPDATE
```

## Risk Analysis

### Identified Risks

1. **Risk**: Existing users depend on current IPC mechanism

   - **Mitigation**: Keep old code in legacy branch, provide migration guide

2. **Risk**: Event stream file grows unbounded

   - **Mitigation**: Implement log rotation (separate task if needed)

3. **Risk**: Loss of bidirectional control features
   - **Mitigation**: Optional command polling provides basic control

## Resource Allocation

### Recommended Team Structure

- **1 Backend Developer**: Focus on Go refactoring (MORGANA_CONSUMER_IMPL,
  REMOVE_IPC_TASK)
- **1 Full-Stack Developer**: Handle Python/Claude integration
  (MORGANA_STREAM_IMPL, CLAUDE_NATIVE_IMPL)
- **1 QA/DevOps**: Testing and documentation (INTEGRATION_TEST, DOCS_UPDATE)

### Solo Developer Approach

Execute tasks in critical path order, prioritizing P0 tasks first.

## Success Metrics

- **Code Reduction**: Remove >500 lines of IPC complexity
- **Performance**: Zero polling overhead, <10ms event latency
- **Reliability**: No race conditions or synchronization issues
- **Simplicity**: New developers understand architecture in <30 minutes

## Post-Sprint Considerations

### Future Enhancements (Not in scope)

- WebSocket support for real-time monitoring
- Event persistence in SQLite
- Advanced filtering in TUI
- Multi-agent execution graphs

### Technical Debt Addressed

- Removes file-based IPC anti-pattern
- Eliminates polling inefficiency
- Resolves Go/Claude execution boundary issues
- Aligns with industry monitoring standards

## Implementation Notes

### Key Design Decisions

1. **Event format**: JSON lines (jsonl) for simplicity and debugging
2. **No message bus**: Direct file append is sufficient
3. **No orchestration**: Claude Code handles all execution
4. **Passive monitoring**: Morgana only observes, never controls

### Code Example - New Event Logger

```python
import json
import time
from pathlib import Path

def log_agent_event(event_type, agent_type, **kwargs):
    event = {
        "timestamp": time.time(),
        "type": event_type,
        "agent": agent_type,
        **kwargs
    }

    event_file = Path("/tmp/morgana/events.jsonl")
    event_file.parent.mkdir(exist_ok=True)

    with open(event_file, "a") as f:
        f.write(json.dumps(event) + "\n")
```

---

## Sprint Summary

**Total Tasks**: 9 (5 critical, 2 high priority, 2 medium priority) **Estimated
Duration**: 40 hours (5 days) **Critical Success Factor**: Simplify radically -
remove IPC complexity

This sprint transforms Morgana from an over-engineered orchestrator into a
clean, simple monitoring tool that respects platform constraints and delivers
immediate value.
