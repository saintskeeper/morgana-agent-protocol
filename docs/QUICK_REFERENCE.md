# Morgana Quick Reference

## Essential Make Commands

### Setup & Installation

```bash
make help           # Show all available commands
make build          # Build binaries (development)
make install-user   # Install to ~/.claude/bin (no sudo)
make install        # Install to /usr/local/bin (requires sudo)
make check          # Full system health check
```

### Daily Usage

```bash
make up             # Start monitor daemon
make down           # Stop monitor daemon
make status         # Check if monitor is running
make attach         # View live TUI monitor
make logs           # Show monitor logs
make restart        # Restart monitor (down + up)
```

### Testing

```bash
make test           # Run single test agent
make test-parallel  # Test parallel agent execution
```

### Development

```bash
make dev-up         # Build, install-user, and start monitor
make dev-down       # Stop and clean development environment
make clean          # Clean all artifacts and stop services
make kill           # Force kill all morgana processes
make ps             # Show running morgana processes
make tail           # Tail monitor log in real-time
```

### Troubleshooting

```bash
make status         # Check monitor status
make logs           # View recent logs
make ps             # List all morgana processes
make kill           # Force kill stuck processes
make clean          # Full cleanup and reset
```

## Common Workflows

### First Time Setup

```bash
cd ~/.claude
make install-user   # Install without sudo
make up             # Start monitoring
make test           # Verify everything works
```

### Daily Development Flow

```bash
# Morning startup
make up             # Start monitor
make status         # Verify it's running

# During work
make attach         # Watch live execution
# Run your agents...

# End of day
make down           # Stop monitor
```

### After Pulling Updates

```bash
git pull
make clean          # Clean old artifacts
make build          # Rebuild binaries
make install-user   # Reinstall
make restart        # Restart monitor
```

### Debugging Issues

```bash
make check          # Full system diagnostic
make ps             # Check for stuck processes
make kill           # Force kill if needed
make clean          # Full reset
make dev-up         # Fresh start
```

## Agent Execution

### Single Agent

```bash
morgana -- --agent code-implementer --prompt "Create a REST API"
```

### Parallel Agents

```bash
echo '[
  {"agent_type":"code-implementer","prompt":"Build feature"},
  {"agent_type":"test-specialist","prompt":"Write tests"}
]' | morgana --parallel
```

### With Custom Timeout

```bash
morgana -- --agent code-implementer --prompt "Complex task" --timeout 10m
```

### With Configuration

```bash
morgana --config custom.yaml -- --agent sprint-planner --prompt "Plan sprint"
```

## Monitoring Views

### Live Monitoring

```bash
make up             # Start daemon (runs in background)
make attach         # Connect TUI client to view events
# Press 'q' to detach (daemon keeps running)
```

### Log Analysis

```bash
make logs           # View recent logs
make tail           # Follow log in real-time
```

### Status Checks

```bash
make status         # Is monitor running?
make check          # Full system health check
make ps             # Show all processes
```

## Quick Tips

1. **Always use `make install-user`** for personal use - no sudo needed
2. **Monitor runs as daemon** - use `make attach` to view, not `make up`
   repeatedly
3. **Use `make check`** when something seems wrong - comprehensive diagnostic
4. **`make clean`** is your friend - resets everything when stuck
5. **`make dev-up`** for one-command development setup

## Environment Variables

```bash
export MORGANA_DEBUG=true              # Enable debug logging
export MORGANA_BRIDGE_PATH=/path/to/bridge.py  # Custom bridge
export PATH=$HOME/.claude/bin:$PATH    # Add to PATH if using install-user
```

## File Locations

| Item              | Location                   |
| ----------------- | -------------------------- |
| Binaries (user)   | `~/.claude/bin/`           |
| Binaries (system) | `/usr/local/bin/`          |
| Configuration     | `~/.claude/morgana.yaml`   |
| Agent definitions | `~/.claude/agents/`        |
| Monitor socket    | `/tmp/morgana.sock`        |
| Monitor PID       | `/tmp/morgana-monitor.pid` |
| Monitor log       | `/tmp/morgana-monitor.log` |

## Common Issues

### "morgana: command not found"

```bash
# If installed with make install-user:
export PATH=$HOME/.claude/bin:$PATH

# Or reinstall system-wide:
make install
```

### "Monitor not running"

```bash
make up         # Start it
make status     # Verify
make attach     # Connect to view
```

### "Address already in use"

```bash
make kill       # Force kill all
make clean      # Clean up
make up         # Start fresh
```

### "Permission denied"

```bash
# Use install-user instead of install
make install-user   # No sudo needed
```

---

**Need help?** Run `make help` for all available commands
