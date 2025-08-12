# Morgana Command Polling System

A lightweight file-based command system for pause/resume control of Morgana
Protocol agents.

## Overview

The command polling system provides optional flow control without complex IPC,
using a simple file-based command system for basic pause/resume/stop
functionality.

## Components

### Core Module: `morgana_command_poll.py`

The main command polling module that provides:

- `MorganaCommandPoller`: Core polling class
- Command processing for pause, resume, stop, status
- Integration with existing morgana_events logging
- Rate limiting to prevent excessive file reads
- Thread-safe execution state management

### Command Utility: `morgana-command.py`

Command-line utility for sending commands:

```bash
# Pause execution
python morgana-command.py pause

# Resume execution
python morgana-command.py resume

# Stop execution
python morgana-command.py stop

# Request status
python morgana-command.py status

# Check current state
python morgana-command.py check
```

### Integration with Claude Agent Executor

The `claude_agent_executor.py` has been enhanced with optional command polling:

- Automatic integration via `integration_check_point()`
- Command checking before and during task execution
- Graceful handling of pause/resume/stop commands
- Logging of command-related state changes

## Usage

### Basic Command Polling

```python
from morgana_command_poll import MorganaCommandPoller

# Create poller
poller = MorganaCommandPoller()

# Poll for commands
while poller.poll_once():
    # Do work
    if poller.should_pause():
        # Handle pause state
        if not poller.wait_while_paused():
            break  # Stopped

    # Continue work...
```

### Integration Checkpoint (Recommended)

For long-running tasks, use the integration checkpoint:

```python
from morgana_command_poll import integration_check_point

for i in range(1000):
    # Check for pause/resume/stop commands
    if not integration_check_point():
        print("Task stopped by command")
        break

    # Do work...
    time.sleep(1)
```

### Claude Agent Executor Integration

Command polling is automatically enabled in `ClaudeAgentExecutor`:

```python
from claude_agent_executor import execute_code_implementer

# Command polling is enabled by default
result = execute_code_implementer("Implement new feature")

# Disable command polling if needed
from claude_agent_executor import ClaudeAgentExecutor
executor = ClaudeAgentExecutor(enable_command_polling=False)
```

## Commands File Format

Commands are written to `/tmp/morgana/commands.txt` as JSON:

```json
{
  "command": "pause",
  "timestamp": "2025-01-01T12:00:00.000Z"
}
```

Supported commands:

- `pause`: Pause execution
- `resume`: Resume execution
- `stop`: Stop execution completely
- `status`: Request current status (logged to events)

## Rate Limiting

The poller includes rate limiting:

- Default: Maximum 10 polls per second
- Configurable via `max_poll_rate` parameter
- Prevents excessive file system access

## Event Logging

All command-related activities are logged via `morgana_events.py`:

- Command processing
- State changes
- Integration checkpoint calls
- Error conditions

## Testing

Run the test script to verify functionality:

```bash
# Basic polling test
python test_command_polling.py basic

# Integration checkpoint test
python test_command_polling.py checkpoint

# Simulated work with command polling
python test_command_polling.py work
```

## Design Principles

1. **Optional**: System continues to work if command polling is disabled
2. **Lightweight**: Simple file-based IPC without complex dependencies
3. **Rate Limited**: Prevents excessive file system access
4. **Thread Safe**: Safe for concurrent usage
5. **Integrated**: Works seamlessly with existing logging and execution systems

## File Locations

- Commands file: `/tmp/morgana/commands.txt`
- Events log: `/tmp/morgana/events.jsonl`
- Module: `/Users/walterday/.claude/scripts/morgana_command_poll.py`
- Utility: `/Users/walterday/.claude/scripts/morgana-command.py`
- Test: `/Users/walterday/.claude/scripts/test_command_polling.py`

## Error Handling

- Commands file not found: Continues normally (optional feature)
- Invalid command format: Logged and ignored
- JSON parsing errors: Logged and ignored
- File access errors: Logged and handled gracefully

The system is designed to be robust and never interrupt execution due to command
system issues.
