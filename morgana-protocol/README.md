# Morgana Protocol

A Go-based agent orchestration system that enables parallel execution of
specialized agents in Claude Code.

## Quick Start

```bash
# Build
make dev

# Test
./dist/morgana --version

# Install
make install
```

## Usage

### Direct CLI Usage

```bash
# Single agent
echo '[{"agent_type":"code-implementer","prompt":"implement feature"}]' | morgana

# Parallel agents
echo '[
  {"agent_type":"code-implementer","prompt":"implement auth"},
  {"agent_type":"test-specialist","prompt":"create tests"}
]' | morgana --parallel
```

### Integration with Claude Code

1. Source the wrapper in your shell:

```bash
source ~/.claude/scripts/agent-adapter-wrapper.sh
```

2. Use in markdown commands:

```bash
# Single agent
AgentAdapter "code-implementer" "implement auth service"

# Parallel execution (in bash)
tasks='[
  {"agent_type":"code-implementer","prompt":"implement auth"},
  {"agent_type":"test-specialist","prompt":"create tests"}
]'
echo "$tasks" | morgana --parallel
```

## Architecture

- **Go-based**: Single binary, no runtime dependencies
- **True Parallelism**: Goroutines with concurrency control
- **Cross-platform**: Compiles for macOS, Linux, Windows
- **JSON API**: Easy integration with any tool

## Available Agents

- `code-implementer` - Code implementation specialist
- `sprint-planner` - Sprint planning and decomposition
- `test-specialist` - Test creation with edge cases
- `validation-expert` - Code review and validation

## Implementation Status

✅ Core adapter infrastructure ✅ Parallel orchestration with goroutines ✅
Agent prompt loading with caching ✅ JSON input/output API ✅ Shell integration
scripts

## Why Go?

- Single binary distribution
- Real concurrent execution
- Fast startup time
- Cross-platform support
- No dependency management issues
