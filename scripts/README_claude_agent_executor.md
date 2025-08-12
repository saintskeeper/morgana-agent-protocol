# Claude Code Agent Executor

A native Claude Code agent executor that directly calls the Task tool and
integrates with the morgana events system. Designed to work within Claude Code
REPL environment without subprocess or file-based communication.

## Features

- **Native Task Tool Integration**: Directly calls Claude Code's Task function
  when available
- **Morgana Events Logging**: Automatically logs all task execution events
- **Fallback Mode**: Works outside Claude Code with mock responses
- **Multi-Agent Support**: Supports all four agent types
- **Parallel Execution**: Can execute multiple agents concurrently
- **Timing and Error Handling**: Comprehensive logging and error handling
- **Convenience Functions**: Easy-to-use wrapper functions for each agent

## Supported Agent Types

- `code-implementer`: Expert code implementation specialist
- `sprint-planner`: Sprint planning and task organization
- `test-specialist`: Comprehensive test generation and validation
- `validation-expert`: Code review and validation specialist

## Usage

### Command Line Interface

```bash
# Single agent execution
python3 claude_agent_executor.py --agent-type code-implementer --prompt "Implement user authentication"

# With optional parameters
python3 claude_agent_executor.py \
    --agent-type code-implementer \
    --prompt "Implement user authentication" \
    --task-id "auth-001" \
    --timeout 300 \
    --model "claude-sonnet-4"

# Parallel execution
echo '[
    {"agent_type": "code-implementer", "prompt": "Task 1"},
    {"agent_type": "test-specialist", "prompt": "Task 2"}
]' | python3 claude_agent_executor.py --parallel
```

### Python API

```python
from claude_agent_executor import ClaudeAgentExecutor

# Create executor instance
executor = ClaudeAgentExecutor()

# Execute single agent
result = executor.execute_agent(
    agent_type="code-implementer",
    prompt="Implement user authentication",
    task_id="auth-001",
    timeout=300,
    model="claude-sonnet-4"
)

# Execute multiple agents in parallel
tasks = [
    {"agent_type": "code-implementer", "prompt": "Task 1"},
    {"agent_type": "test-specialist", "prompt": "Task 2"}
]
results = executor.execute_parallel(tasks)
```

### Convenience Functions

```python
from claude_agent_executor import (
    execute_code_implementer,
    execute_sprint_planner,
    execute_test_specialist,
    execute_validation_expert
)

# Use convenience functions for each agent type
result = execute_code_implementer("Implement user authentication")
result = execute_test_specialist("Generate tests for authentication")
result = execute_validation_expert("Review authentication implementation")
result = execute_sprint_planner("Plan authentication feature sprint")
```

## Response Format

All execution functions return a dictionary with the following structure:

```python
{
    "success": True,                    # Boolean indicating success/failure
    "agent_type": "code-implementer",   # Type of agent executed
    "task_id": "auth-001",              # Unique task identifier
    "result": {...},                    # Agent execution result
    "duration_ms": 1500,               # Execution duration in milliseconds
    "model": "claude-sonnet-4",         # Model used (if specified)
    "execution_mode": "claude_code_native"  # Execution mode
}
```

## Execution Modes

- `claude_code_native`: Task function available, native Claude Code execution
- `mock_fallback`: Task function not available, returns mock response
- `error`: Execution failed with error

## Event Logging

All task executions are automatically logged to `/tmp/morgana/events.jsonl` with
the following event types:

- `task_started`: Task execution begins
- `task_progress`: Task progress updates (if supported)
- `task_completed`: Task completes successfully
- `task_failed`: Task fails with error

Each event includes:

- Timestamp (ISO 8601 format)
- Session ID
- Process ID
- Task ID
- Agent type
- Event-specific data

## Environment Detection

The executor automatically detects whether it's running in a Claude Code
environment by searching for the `Task` function in the call stack and global
namespace. When the Task function is available, it executes natively. Otherwise,
it falls back to mock mode for testing and development outside Claude Code.

## Error Handling

- Input validation for agent types and required parameters
- Comprehensive exception handling with logging
- Graceful fallback when Task function is unavailable
- Detailed error messages in response format

## Agent Configuration

Agent configurations are loaded from `~/.claude/agents/*.md` files. Each agent
has its own specialized prompt and capabilities defined in markdown format with
YAML frontmatter.

## Integration with Morgana Protocol

The executor integrates seamlessly with the Morgana Protocol monitoring system:

- Events are logged to the standard morgana events stream
- Compatible with existing morgana monitoring tools
- Supports session tracking and parallel execution monitoring

## Requirements

- Python 3.6+
- Access to `~/.claude/agents/` directory with agent configurations
- morgana_events.py for event logging (optional, graceful degradation if not
  available)
- Claude Code environment for native Task tool execution (falls back to mock
  mode if not available)
