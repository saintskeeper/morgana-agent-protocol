# Morgana Monitor Auto-Start Integration

This document describes the auto-start functionality implemented for the
morgana-monitor daemon, providing zero-configuration monitoring for Morgana
Protocol agents.

## Overview

The shell integration scripts have been updated to automatically start and
manage the morgana-monitor daemon when agents are executed. This ensures
monitoring capabilities are always available without manual intervention.

## Implementation

### Updated Files

1. **`~/.claude/scripts/morgana-adapter.sh`**

   - Added `ensure_morgana_monitor()` function for auto-start capability
   - Modified `AgentAdapter()` and `AgentAdapterParallel()` to call monitor
     check before execution
   - Uses `/tmp/morgana.sock` presence to determine if monitor is running
   - Supports multiple startup methods: screen, tmux, script, or direct
     execution

2. **`~/.claude/scripts/morgana-monitor-ctl.sh`** (NEW)
   - Control script for managing the morgana-monitor daemon
   - Commands: start, stop, restart, status, attach, logs
   - PID file management at `/tmp/morgana-monitor.pid`
   - Session management for screen/tmux integration

### Auto-Start Logic

The `ensure_morgana_monitor()` function implements the following logic:

1. **Detection**: Check if `/tmp/morgana.sock` exists and has an active process
2. **Binary Location**: Search for morgana-monitor binary in standard locations
3. **Startup Methods** (in priority order):

   - **screen**: Creates detached session `morgana-monitor`
   - **tmux**: Creates detached session `morgana-monitor`
   - **script**: Uses script command for pseudo-terminal
   - **fallback**: Direct execution with TERM environment variable

4. **Verification**: Wait up to 5 seconds for socket creation before declaring
   success

## Usage

### Automatic Usage

The monitor auto-starts automatically when using agent functions:

```bash
# Source the adapter (monitor starts if needed)
source ~/.claude/scripts/agent-adapter-wrapper.sh

# Execute agent (monitor auto-starts if not running)
AgentAdapter "code-implementer" "implement user service"

# Parallel execution (monitor auto-starts if not running)
AgentAdapterParallel
```

### Manual Control

Use the control script for manual management:

```bash
# Start monitor daemon
~/.claude/scripts/morgana-monitor-ctl.sh start

# Check status
~/.claude/scripts/morgana-monitor-ctl.sh status

# Attach to TUI (if screen/tmux available)
~/.claude/scripts/morgana-monitor-ctl.sh attach

# View logs
~/.claude/scripts/morgana-monitor-ctl.sh logs

# Stop daemon
~/.claude/scripts/morgana-monitor-ctl.sh stop
```

## Session Management

When screen or tmux are available, the monitor runs in detached sessions:

### Screen Session

```bash
# View TUI
screen -r morgana-monitor

# Detach with: Ctrl+A, D
```

### Tmux Session

```bash
# View TUI
tmux attach -t morgana-monitor

# Detach with: Ctrl+B, D
```

## Files and Paths

| File                       | Purpose                    |
| -------------------------- | -------------------------- |
| `/tmp/morgana.sock`        | Unix domain socket for IPC |
| `/tmp/morgana-monitor.pid` | Process ID file            |
| `/tmp/morgana-monitor.log` | Log output (fallback mode) |

## Zero Configuration

The system is designed for zero configuration:

1. **No manual daemon management** - Monitor starts automatically when needed
2. **No configuration files** - Uses optimal defaults
3. **Session persistence** - Monitor runs in background sessions when possible
4. **Graceful degradation** - Falls back to simpler methods if screen/tmux
   unavailable

## Troubleshooting

### Monitor Won't Start

- Ensure morgana-monitor binary is built and available
- Check permissions on `/tmp/` directory
- Verify TERM environment variable is set

### Socket Issues

- Remove stale socket: `rm /tmp/morgana.sock`
- Check for conflicting processes: `lsof /tmp/morgana.sock`

### Session Issues

- List screen sessions: `screen -list`
- Clean up sessions: `screen -wipe`
- For tmux: `tmux list-sessions` and `tmux kill-session -t morgana-monitor`

## Benefits

1. **Seamless Integration**: No manual daemon management required
2. **Visual Monitoring**: TUI available via screen/tmux sessions
3. **Persistent Sessions**: Monitor continues running between agent executions
4. **Graceful Handling**: Robust cleanup and session management
5. **Multiple Fallbacks**: Works across different system configurations

## Future Enhancements

- Configuration file support for custom socket paths
- Systemd integration for system-level daemon management
- Monitoring multiple projects simultaneously
- Log rotation and retention policies
