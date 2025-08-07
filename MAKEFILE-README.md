# Morgana System Makefile

Convenient commands for managing the Morgana Protocol system and TUI monitoring.

## Quick Start

```bash
# Install binaries and start monitoring
make install
make up

# Check system status
make status
make check

# View live TUI
make attach
```

## Available Commands

### Core Operations

| Command        | Description                                   |
| -------------- | --------------------------------------------- |
| `make up`      | Start morgana-monitor daemon (TUI monitoring) |
| `make down`    | Stop morgana-monitor daemon                   |
| `make status`  | Show morgana-monitor status                   |
| `make restart` | Restart the monitor (down + up)               |
| `make attach`  | Attach to running TUI session (screen/tmux)   |

### Installation

| Command               | Description                                           |
| --------------------- | ----------------------------------------------------- |
| `make install`        | Build and install to `/usr/local/bin` (requires sudo) |
| `make install-user`   | Install to `~/.claude/bin` (no sudo required)         |
| `make uninstall`      | Remove from `/usr/local/bin` (requires sudo)          |
| `make uninstall-user` | Remove from `~/.claude/bin`                           |

### Development

| Command         | Description                               |
| --------------- | ----------------------------------------- |
| `make build`    | Build morgana binaries                    |
| `make clean`    | Clean build artifacts and stop services   |
| `make dev-up`   | Build, install locally, and start monitor |
| `make dev-down` | Stop monitor and clean                    |

### Testing & Debugging

| Command              | Description                       |
| -------------------- | --------------------------------- |
| `make test`          | Run a test agent to verify system |
| `make test-parallel` | Run parallel test agents          |
| `make check`         | Full system health check          |
| `make logs`          | Show morgana-monitor logs         |
| `make tail`          | Tail monitor log in real-time     |
| `make ps`            | Show running morgana processes    |

## Installation Methods

### System-wide Installation (Recommended)

```bash
# Build and install to /usr/local/bin
make install

# This installs:
# - /usr/local/bin/morgana
# - /usr/local/bin/morgana-monitor
```

### User Installation (No sudo)

```bash
# Install to ~/.claude/bin
make install-user

# Add to PATH if needed:
export PATH=$HOME/.claude/bin:$PATH
```

## Usage Examples

### Basic Workflow

```bash
# 1. Install the system
make install

# 2. Start monitoring
make up

# 3. Check status
make status

# 4. View TUI
make attach
# Press Ctrl+A then D to detach from screen
# Or Ctrl+B then D to detach from tmux

# 5. Run agents (monitor auto-connects)
morgana -- --agent code-implementer --prompt "Build feature"

# 6. Stop when done
make down
```

### Development Workflow

```bash
# Quick development setup
make dev-up

# Make changes to code...

# Rebuild and restart
make restart

# Run tests
make test
make test-parallel

# Check everything
make check

# Clean up
make dev-down
```

### Troubleshooting

```bash
# Check system health
make check

# View processes
make ps

# Check logs
make logs
make tail

# Force clean restart
make clean
make up
```

## System Architecture

The Makefile manages two main components:

1. **morgana** - The main agent orchestration CLI
2. **morgana-monitor** - Persistent daemon with TUI for monitoring

The monitor runs as a background daemon and receives events from all morgana
executions via IPC (Unix domain socket at `/tmp/morgana.sock`).

## File Locations

- **Binaries**: `/usr/local/bin/` or `~/.claude/bin/`
- **IPC Socket**: `/tmp/morgana.sock`
- **PID File**: `/tmp/morgana-monitor.pid`
- **Logs**: `/tmp/morgana-monitor.log`
- **Source**: `./morgana-protocol/`

## Requirements

- Go 1.21+ (for building)
- screen or tmux (for TUI sessions)
- sudo (for system-wide installation)

## Tips

- Use `make attach` to view the live TUI monitoring
- The monitor auto-starts when agents run via AgentAdapter
- Use `make check` to verify system health
- Run `make test` to verify everything works

## Uninstalling

```bash
# Remove from /usr/local/bin
make uninstall

# Remove from ~/.claude/bin
make uninstall-user

# Clean all artifacts
make clean
```
