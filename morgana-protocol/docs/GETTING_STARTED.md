# Getting Started with Morgana Protocol

Morgana Protocol is a simplified event monitoring system for specialized agents
in Claude Code. This guide will help you get up and running quickly.

## Prerequisites

- Go 1.21+ (for building from source)
- Python 3.8+ (for Claude Code integration)
- Terminal with Unicode support (for TUI)

## Quick Setup

### 1. Installation

```bash
# Clone the repository
git clone <repo-url>
cd morgana-protocol

# Build and install
make build && make install

# Verify installation
morgana --version
```

### 2. Start the Monitor

```bash
# Start the event monitor daemon
make up

# Check that it's running
make status
```

### 3. Run Your First Agent

```bash
# In one terminal, attach the TUI monitor
make attach

# In another terminal, run an agent
morgana -- --agent code-implementer --prompt "Write a hello world function in Python"
```

You should see real-time events flowing through the TUI interface!

## Basic Usage Patterns

### Single Agent Execution

```bash
# Simple agent task
morgana -- --agent code-implementer --prompt "Implement user authentication"

# With custom timeout
morgana --config custom.yaml -- --agent validation-expert --prompt "Review this code"
```

### Parallel Agent Execution

```bash
# JSON input for multiple agents
echo '[
  {"agent_type":"code-implementer","prompt":"Write the login function"},
  {"agent_type":"test-specialist","prompt":"Create unit tests"},
  {"agent_type":"validation-expert","prompt":"Review implementation"}
]' | morgana --parallel
```

### Live Monitoring

```bash
# Start monitor daemon
make up

# Attach TUI in background
make attach &

# Run multiple tasks
for i in {1..3}; do
  morgana -- --agent code-implementer --prompt "Task $i" &
done
wait

# Stop monitor when done
make down
```

## Configuration

### Basic Configuration (morgana.yaml)

```yaml
# Event system configuration
tui:
  enabled: true
  performance:
    refresh_rate: 50ms
    max_log_lines: 1000
    target_fps: 20

  # Visual theme
  visual:
    theme:
      name: dark
      primary: "#7C3AED"
      secondary: "#06B6D4"

# Agent configuration
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m
  timeouts:
    code-implementer: 5m
    test-specialist: 3m

# Event buffer size
tui:
  events:
    buffer_size: 1000
    enable_batching: true
```

### Environment Variables

```bash
# Enable TUI mode
export MORGANA_TUI_ENABLED=true

# Set refresh rate
export MORGANA_TUI_REFRESH_RATE=16ms

# Choose theme
export MORGANA_TUI_THEME=dark

# Set buffer size
export MORGANA_TUI_BUFFER_SIZE=2000
```

## Integration with Claude Code

### Shell Integration

Add to your shell profile (`.bashrc`, `.zshrc`):

```bash
# Source the agent wrapper
source ~/.claude/scripts/agent-adapter-wrapper.sh

# Convenience aliases
alias monitor-up='cd ~/.claude/morgana-protocol && make up'
alias monitor-attach='cd ~/.claude/morgana-protocol && make attach'
alias monitor-down='cd ~/.claude/morgana-protocol && make down'
```

### Python Integration (Claude Code)

```python
# The existing AgentAdapter interface works unchanged
from agent_adapter import AgentAdapter

result = AgentAdapter("code-implementer", "Write a function")
```

## TUI Interface Guide

### Main Interface

The TUI shows:

- **Task List**: Active and completed agent tasks
- **Event Stream**: Real-time event log
- **Progress Bars**: Per-agent execution progress
- **Status Panel**: System statistics

### Keyboard Shortcuts

- `q` or `Ctrl+C`: Quit TUI
- `↑/↓`: Navigate task list
- `Enter`: View task details
- `r`: Refresh display
- `f`: Toggle filtering
- `h`: Show help

### Event Types

- `task.started`: Agent task begins
- `task.progress`: Execution progress update
- `task.completed`: Successful completion
- `task.failed`: Task failure
- `orchestrator.*`: Batch operations

## Troubleshooting

### Common Issues

**Monitor won't start:**

```bash
# Check for existing processes
make status

# Force cleanup if needed
pkill -f morgana-monitor
make up
```

**TUI display issues:**

```bash
# Check terminal compatibility
echo $TERM

# Try different refresh rate
export MORGANA_TUI_REFRESH_RATE=100ms
make attach
```

**Events not showing:**

```bash
# Verify monitor is running
make status

# Check socket permissions
ls -la /tmp/morgana-monitor.sock

# Restart monitor
make down && make up
```

### Performance Tuning

**High event volume:**

```yaml
# morgana.yaml
tui:
  performance:
    refresh_rate: 100ms # Reduce for better performance
    max_log_lines: 500 # Lower memory usage
    target_fps: 10 # Less CPU usage

  events:
    buffer_size: 2000 # Larger buffer
    batch_size: 100 # Batch processing
```

**Low-resource systems:**

```bash
# Disable TUI features
export MORGANA_TUI_SHOW_DEBUG=false
export MORGANA_TUI_COMPACT_MODE=true
```

## Next Steps

1. **Explore Configuration**: Customize themes and performance settings
2. **Try Advanced Features**: Event filtering, export capabilities
3. **Integration**: Use with your existing Claude Code workflows
4. **Monitoring**: Set up the full observability stack

## Getting Help

- Check the [Architecture Guide](ARCHITECTURE.md) for system design
- Read [Integration Guide](INTEGRATION.md) for advanced usage
- View [Event System Documentation](EVENT_SYSTEM.md) for technical details

The event monitoring system is designed to be zero-configuration but highly
customizable when needed. Start simple and add complexity as your monitoring
needs grow.
