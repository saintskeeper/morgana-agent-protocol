# Morgana Protocol Integration Guide

## Overview

The Morgana Protocol now includes a Python bridge that enables the Go
orchestrator to call Claude Code's Task function. This hybrid architecture keeps
all complex logic in Go while using Python only for the Task API boundary.

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Claude Code    │────▶│  agent_adapter.py│────▶│  morgana (Go)   │
│  Python Env     │     │                  │     │                 │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                                  │                        │
                                  ▼                        ▼
                        ┌──────────────────┐     ┌─────────────────┐
                        │ task_bridge.py   │◀────│  Task Client    │
                        │ (Python Bridge)  │     │  (Go)           │
                        └──────────────────┘     └─────────────────┘
                                  │
                                  ▼
                        ┌──────────────────┐
                        │  Task() Function │
                        │  (Claude Code)   │
                        └──────────────────┘
```

## Components

### 1. Go Orchestrator (`morgana`)

- Handles all orchestration logic
- Manages parallel execution with goroutines
- Routes specialized agents to general-purpose
- Cross-platform binary (darwin/linux/windows)

### 2. Python Bridge (`task_bridge.py`)

- Minimal Python script (< 100 lines)
- Interfaces with Claude Code's Task() function
- Communicates via JSON over stdin/stdout
- Returns mock responses when outside Claude Code

### 3. Agent Adapter (`agent_adapter.py`)

- Drop-in replacement for AgentAdapter function
- Calls the Go orchestrator
- Maintains backward compatibility

## Usage

### Direct CLI Usage

```bash
# Single agent
./dist/morgana -- --agent code-implementer --prompt "Write a function"

# Parallel agents
./dist/morgana --parallel -- \
  --agent code-implementer --prompt "Write code" \
  --agent test-specialist --prompt "Write tests"

# Debug mode
MORGANA_DEBUG=true ./dist/morgana -- --agent sprint-planner --prompt "Plan sprint"
```

### From Python (Claude Code)

```python
from agent_adapter import AgentAdapter

# Single agent
result = AgentAdapter("code-implementer", "Write a hello world function")

# With options
result = AgentAdapter(
    "validation-expert",
    "Validate the implementation",
    timeout=30,
    context="web-app"
)
```

### Parallel Execution (Python)

```python
import asyncio
from agent_adapter import AgentAdapter

async def parallel_agents():
    tasks = [
        asyncio.create_task(AgentAdapter("code-implementer", "Write function")),
        asyncio.create_task(AgentAdapter("test-specialist", "Write tests")),
        asyncio.create_task(AgentAdapter("validation-expert", "Review code"))
    ]
    results = await asyncio.gather(*tasks)
    return results
```

## Environment Variables

- `MORGANA_DEBUG`: Enable debug output (`true`/`false`)
- `MORGANA_BRIDGE_PATH`: Custom path to task_bridge.py
- `MORGANA_OPTIONS`: Additional options as JSON

## Installation

1. Build the Go binary:

   ```bash
   cd /Users/walterday/.claude/morgana-protocol
   make build
   ```

2. The Python bridge is automatically found in:
   - `./scripts/task_bridge.py` (relative to binary)
   - `$HOME/.claude/morgana-protocol/scripts/task_bridge.py`
   - Or set `MORGANA_BRIDGE_PATH`

## Testing

Outside Claude Code environment, the system returns mock responses:

```
[MOCK] Executed general-purpose agent with prompt length: 4299
```

Inside Claude Code, it will execute real Task() calls.

## Limitations

1. **Python Dependency**: Python 3 is required for the bridge
2. **Subprocess Overhead**: Each task incurs ~50-100ms overhead
3. **No Batching**: Tasks are executed one at a time (though in parallel via
   goroutines)

## Future Improvements

1. **gRPC/HTTP Server**: Replace subprocess with persistent server
2. **Connection Pooling**: Reuse Python processes
3. **Circuit Breaker**: Handle Task() failures gracefully
4. **TypeScript Alternative**: Support Node.js runtime
