# Morgana Event Stream Logger

A simple event logging system for Morgana Protocol that writes JSON lines to
`/tmp/morgana/events.jsonl` for real-time monitoring of agent activities.

## Core Components

### 1. `morgana_events.py` - Core Event Logger

The main event logging library with the `MorganaEventLogger` class.

**Key Features:**

- Writes JSON lines to `/tmp/morgana/events.jsonl`
- Includes timestamp, session_id, and event metadata
- Immediate flush for real-time monitoring
- Support for task lifecycle events

**Basic Usage:**

```python
from morgana_events import get_logger

logger = get_logger()
task_id = "my-task-001"

# Log task lifecycle
logger.task_started(task_id, "code-implementer", "Implement user auth")
logger.task_progress(task_id, "setup", "Setting up dependencies", 0.2)
logger.task_completed(task_id, "Auth system implemented", 1500, "claude-sonnet-4")
```

**Decorator Usage:**

```python
from morgana_events import log_claude_task

@log_claude_task(agent_type="code-implementer")
def my_coding_task():
    # Task automatically logged
    return "Task completed"
```

### 2. `morgana_events_viewer.py` - Event Viewer

Real-time event log viewer with filtering and session management.

**Usage:**

```bash
# View all events
python3 morgana_events_viewer.py

# Follow new events (like tail -f)
python3 morgana_events_viewer.py --follow

# Filter by session
python3 morgana_events_viewer.py --session abc12345

# List all sessions
python3 morgana_events_viewer.py --list-sessions
```

### 3. `morgana_events_integration.sh` - Bash Integration

Shell functions for easy event logging in bash scripts.

**Usage:**

```bash
# Source the integration functions
source /Users/walterday/.claude/scripts/morgana_events_integration.sh

# Generate task ID and log manually
task_id=$(morgana_generate_task_id)
morgana_log_task_start "$task_id" "test-specialist" "Running tests"
morgana_log_task_progress "$task_id" "unit-tests" "Running unit tests" 0.5
morgana_log_task_complete "$task_id" "All tests passed" 1200

# Or use the timing wrapper
morgana_time_task "$task_id" "build-agent" "Building project" make clean build
```

**Available Functions:**

- `morgana_generate_task_id` - Generate unique task ID
- `morgana_log_task_start <task_id> <agent_type> <prompt> [options]`
- `morgana_log_task_progress <task_id> <stage> <message> [progress]`
- `morgana_log_task_complete <task_id> <output> <duration_ms> [model]`
- `morgana_log_task_error <task_id> <error> <duration_ms> [stage]`
- `morgana_time_task <task_id> <agent_type> <prompt> <command...>` - Execute and
  time command
- `morgana_view_events [options]` - View events
- `morgana_follow_events [options]` - Follow events in real-time
- `morgana_list_sessions [options]` - List active sessions

## Event Types

### task_started

Logged when a task begins execution.

```json
{
  "event_type": "task_started",
  "task_id": "abc123",
  "agent_type": "code-implementer",
  "prompt": "Implement user authentication",
  "options": {},
  "timestamp": "2025-08-11T18:27:36.536008+00:00",
  "session_id": "343cadf6",
  "pid": 37731
}
```

### task_progress

Logged during task execution to show progress.

```json
{
  "event_type": "task_progress",
  "task_id": "abc123",
  "stage": "implementation",
  "message": "Writing authentication module",
  "progress": 0.6,
  "timestamp": "2025-08-11T18:27:37.123456+00:00",
  "session_id": "343cadf6",
  "pid": 37731
}
```

### task_completed

Logged when a task completes successfully.

```json
{
  "event_type": "task_completed",
  "task_id": "abc123",
  "output": "Successfully implemented authentication",
  "duration_ms": 1500,
  "model": "claude-sonnet-4",
  "timestamp": "2025-08-11T18:27:38.036008+00:00",
  "session_id": "343cadf6",
  "pid": 37731
}
```

### task_failed

Logged when a task fails.

```json
{
  "event_type": "task_failed",
  "task_id": "abc123",
  "error": "Authentication service unavailable",
  "duration_ms": 800,
  "stage": "testing",
  "timestamp": "2025-08-11T18:27:38.836008+00:00",
  "session_id": "343cadf6",
  "pid": 37731
}
```

## Integration with Claude Code Tasks

For Claude Code Task() integration, use the decorator:

```python
from morgana_events import log_claude_task

@log_claude_task(task_id="custom-id", agent_type="code-implementer")
def implement_feature():
    # Your implementation
    return result
```

Or use manual logging for more control:

```python
from morgana_events import get_logger
import time

logger = get_logger()
task_id = "feature-001"
start_time = time.time()

try:
    logger.task_started(task_id, "code-implementer", "Implement new feature")

    # Your implementation work
    result = implement_feature()

    duration_ms = int((time.time() - start_time) * 1000)
    logger.task_completed(task_id, str(result), duration_ms)

except Exception as e:
    duration_ms = int((time.time() - start_time) * 1000)
    logger.task_failed(task_id, str(e), duration_ms)
    raise
```

## Real-time Monitoring

To monitor events in real-time:

```bash
# Terminal 1: Follow events
python3 /Users/walterday/.claude/scripts/morgana_events_viewer.py --follow

# Terminal 2: Run your agent tasks
# Events will appear in Terminal 1 as they happen
```

## File Locations

- **Event Log:** `/tmp/morgana/events.jsonl`
- **Core Logger:** `/Users/walterday/.claude/scripts/morgana_events.py`
- **Event Viewer:** `/Users/walterday/.claude/scripts/morgana_events_viewer.py`
- **Bash Integration:**
  `/Users/walterday/.claude/scripts/morgana_events_integration.sh`
- **Examples:** `/Users/walterday/.claude/scripts/morgana_events_example.py`

## Architecture

The event logger is designed to be:

- **Simple:** No complex IPC or polling - just write JSON lines to a file
- **Real-time:** Immediate flush ensures events are written instantly
- **Session-aware:** Each execution gets a unique session ID
- **Language-agnostic:** Python core with bash integration, easily extensible
- **Non-intrusive:** Minimal overhead, continues working even if logging fails

This replaces the complex IPC mechanisms with a straightforward file-based
approach that's easy to monitor and debug.
