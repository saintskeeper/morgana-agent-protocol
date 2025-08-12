# Morgana Protocol Migration Guide

This guide helps you migrate from the old IPC-based Morgana Protocol to the new
simplified event stream architecture.

## Overview

On August 11, 2025, Morgana Protocol underwent a major architectural
simplification:

- **Removed**: Complex file-based IPC bridge system (474+ lines of code)
- **Added**: Simple, high-performance event stream monitoring
- **Result**: Easier maintenance, better real-time monitoring,
  zero-configuration setup

## Migration Summary

| Component         | Old Architecture      | New Architecture               |
| ----------------- | --------------------- | ------------------------------ |
| **Communication** | File-based IPC bridge | Event stream + Unix sockets    |
| **Dependencies**  | Python bridge scripts | Pure Go implementation         |
| **Configuration** | Complex IPC setup     | Zero-configuration defaults    |
| **Performance**   | File I/O overhead     | 5M+ events/sec in-memory       |
| **Monitoring**    | Limited visibility    | Real-time TUI with buffering   |
| **Maintenance**   | Complex debugging     | Simplified event-driven design |

## Step-by-Step Migration

### 1. Update Your Installation

First, update to the latest version:

```bash
# Navigate to your Morgana Protocol directory
cd ~/.claude/morgana-protocol

# Pull latest changes
git pull origin main

# Rebuild and reinstall
make clean && make build && make install

# Verify new version
morgana --version
```

### 2. Configuration Migration

#### A. Remove Old IPC Configuration

**Old configuration (remove these sections):**

```yaml
# ❌ REMOVE: This entire section is no longer needed
task_client:
  bridge_path: ~/.claude/scripts/task_bridge.py
  python_path: python3
  mock_mode: false
  timeout: 5m
```

#### B. Update to New Event-Based Configuration

**New configuration (already in place):**

```yaml
# ✅ NEW: Event monitoring configuration
tui:
  enabled: true
  performance:
    refresh_rate: 50ms
    max_log_lines: 1000
    target_fps: 20

  events:
    buffer_size: 1000
    enable_batching: true
    batch_size: 50
    batch_timeout: 100ms
```

#### C. Clean Up Old Files

**Remove deprecated configuration:**

```bash
# Check if you have old configuration files
ls -la ~/.claude/morgana-protocol/morgana.yaml*

# If you have custom configurations, migrate manually:
# - Remove task_client sections
# - Add tui.events configuration if customizing
# - Keep agent and execution settings (unchanged)
```

### 3. Update Your Workflows

#### A. Basic Usage (Unchanged)

The core agent execution interface remains the same:

```bash
# ✅ Still works exactly the same
morgana -- --agent code-implementer --prompt "Write a hello function"

# ✅ Parallel execution unchanged
echo '[{"agent_type":"code-implementer","prompt":"Task 1"}]' | morgana --parallel
```

#### B. Enable Real-Time Monitoring (New Feature)

**Start the monitor daemon:**

```bash
# Start monitor in background
make up

# Check status
make status

# View real-time events
make attach
```

**In your shell profile (`.bashrc`, `.zshrc`):**

```bash
# Add convenience aliases
alias monitor-up='cd ~/.claude/morgana-protocol && make up'
alias monitor-attach='cd ~/.claude/morgana-protocol && make attach'
alias monitor-down='cd ~/.claude/morgana-protocol && make down'
```

### 4. Integration Updates

#### A. Claude Code Integration (No Changes)

Your existing Claude Code integration continues to work:

```python
# ✅ AgentAdapter interface unchanged
from agent_adapter import AgentAdapter

result = AgentAdapter("validation-expert", "Review this code")
```

#### B. Shell Integration (Enhanced)

**Update your agent wrapper:**

```bash
# Source the enhanced adapter wrapper
source ~/.claude/scripts/agent-adapter-wrapper.sh

# New functions available:
# - AgentAdapter: Execute single agent (same as before)
# - AgentAdapterParallel: Execute multiple agents in parallel
# - morgana_parallel: Helper for parallel execution
```

### 5. Verify Migration

#### A. Test Basic Functionality

```bash
# Test single agent execution
morgana -- --agent code-implementer --prompt "Write a simple test function"

# Test parallel execution
echo '[
  {"agent_type":"code-implementer","prompt":"Write a function"},
  {"agent_type":"test-specialist","prompt":"Create tests"}
]' | morgana --parallel
```

#### B. Test Monitor System

```bash
# Start monitor
make up

# In another terminal, attach TUI
make attach

# In a third terminal, run a test task
morgana -- --agent validation-expert --prompt "Test task for monitoring"

# You should see real-time events in the TUI
```

#### C. Performance Test

```bash
# Run integration tests to verify performance
make test-integration

# Run performance validation
make test-performance
```

## Feature Comparison

### Removed Features

| Feature                   | Status     | Migration Path                |
| ------------------------- | ---------- | ----------------------------- |
| **File-based IPC**        | ❌ Removed | Replaced with event streaming |
| **Python bridge scripts** | ❌ Removed | Pure Go implementation        |
| **Complex task client**   | ❌ Removed | Simplified adapter interface  |
| **IPC configuration**     | ❌ Removed | Zero-configuration defaults   |

### New Features

| Feature             | Description                          | Benefits                             |
| ------------------- | ------------------------------------ | ------------------------------------ |
| **Event Stream**    | High-performance in-memory event bus | 5M+ events/sec, real-time monitoring |
| **TUI Monitor**     | Beautiful terminal interface         | Live progress, event history, themes |
| **Circular Buffer** | Event persistence without storage    | No data loss, late-joining clients   |
| **Unix Sockets**    | Efficient local IPC                  | Low latency, multi-client support    |
| **Zero Config**     | Works out of the box                 | No setup required                    |

### Enhanced Features

| Feature         | Old                 | New                  | Improvement          |
| --------------- | ------------------- | -------------------- | -------------------- |
| **Performance** | File I/O bottleneck | In-memory processing | 100x+ faster         |
| **Monitoring**  | Limited logging     | Real-time TUI        | Full visibility      |
| **Debugging**   | Complex file traces | Event stream logs    | Easy troubleshooting |
| **Concurrency** | File locking issues | Lock-free design     | Better parallelism   |

## Configuration Reference

### Environment Variables

All configuration can be overridden with environment variables:

```bash
# TUI Configuration
export MORGANA_TUI_ENABLED=true
export MORGANA_TUI_REFRESH_RATE=16ms
export MORGANA_TUI_THEME=dark
export MORGANA_TUI_BUFFER_SIZE=2000

# Performance Tuning
export MORGANA_TUI_MAX_LOG_LINES=500
export MORGANA_TUI_TARGET_FPS=30
export MORGANA_TUI_COMPACT_MODE=true
```

### YAML Configuration

**Complete example configuration:**

```yaml
# Morgana Protocol Configuration

agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m
  timeouts:
    code-implementer: 5m
    test-specialist: 3m
    validation-expert: 2m
  retry:
    enabled: true
    max_attempts: 3

execution:
  max_concurrency: 5
  default_mode: sequential

telemetry:
  enabled: true
  exporter: stdout
  service_name: morgana-protocol

tui:
  enabled: true
  performance:
    refresh_rate: 50ms
    max_log_lines: 1000
    target_fps: 20
  visual:
    theme:
      name: dark
      primary: "#7C3AED"
  events:
    buffer_size: 1000
    enable_batching: true
```

## Troubleshooting Migration

### Common Migration Issues

#### 1. Old Configuration Conflicts

**Problem**: Error messages about unknown `task_client` configuration.

**Solution**:

```bash
# Edit your morgana.yaml
vim ~/.claude/morgana-protocol/morgana.yaml

# Remove these sections:
# task_client:
#   bridge_path: ...
#   python_path: ...
#   mock_mode: ...
#   timeout: ...
```

#### 2. Monitor Won't Start

**Problem**: `make up` fails or monitor doesn't start.

**Solution**:

```bash
# Clean up old processes
pkill -f morgana-monitor

# Remove old socket files
rm -f /tmp/morgana*.sock

# Try starting again
make up
make status
```

#### 3. TUI Display Issues

**Problem**: TUI shows garbled text or doesn't refresh properly.

**Solution**:

```bash
# Check terminal compatibility
echo $TERM

# Try slower refresh rate
export MORGANA_TUI_REFRESH_RATE=100ms

# Reset terminal if corrupted
make down
reset
make up && make attach
```

#### 4. Events Not Showing

**Problem**: No events appear in TUI monitor.

**Solution**:

```bash
# Verify monitor is running
make status

# Check socket permissions
ls -la /tmp/morgana*.sock

# Restart monitor system
make down
sleep 2
make up
make attach
```

### Performance Issues

#### High Resource Usage

**Problem**: Monitor uses too much CPU/memory.

**Solution**:

```bash
# Reduce refresh rate
export MORGANA_TUI_REFRESH_RATE=200ms
export MORGANA_TUI_TARGET_FPS=5

# Use compact mode
export MORGANA_TUI_COMPACT_MODE=true
export MORGANA_TUI_MAX_LOG_LINES=500
```

#### Event Loss

**Problem**: Missing events during high load.

**Solution**:

```yaml
# Increase buffer size in morgana.yaml
tui:
  events:
    buffer_size: 5000
    batch_size: 200
```

## Rollback Instructions

If you need to rollback to the old IPC-based system:

### 1. Access Archived Files

```bash
# Navigate to archive directory
cd ~/.claude/morgana-protocol/archive/ipc-removal

# Review archived implementation
ls -la
cat README.md
```

### 2. Restore from Git History

```bash
# Find the commit before IPC removal
git log --oneline | grep -i ipc

# Create a rollback branch
git checkout -b rollback-to-ipc <commit-hash>

# Build old version
make clean && make build && make install
```

### 3. Restore Old Configuration

```bash
# Restore old configuration sections
cat >> morgana.yaml << EOF
task_client:
  bridge_path: ~/.claude/scripts/task_bridge.py
  python_path: python3
  mock_mode: false
  timeout: 5m
EOF
```

**Note**: The old IPC system had known reliability and performance issues.
Consider reporting specific problems with the new system rather than rolling
back.

## Migration Scripts

### Automated Configuration Cleanup

```bash
#!/bin/bash
# migrate-config.sh - Clean up old configuration

CONFIG_FILE="$HOME/.claude/morgana-protocol/morgana.yaml"
BACKUP_FILE="${CONFIG_FILE}.backup-$(date +%Y%m%d)"

echo "🔄 Migrating Morgana configuration..."

# Backup original
cp "$CONFIG_FILE" "$BACKUP_FILE"
echo "📋 Backed up to: $BACKUP_FILE"

# Remove old task_client section
sed -i.tmp '/^task_client:/,/^[a-z]/{ /^[a-z]/!d; /^task_client:/d; }' "$CONFIG_FILE"
rm -f "${CONFIG_FILE}.tmp"

echo "✅ Migration complete!"
echo "📝 Review changes: diff $BACKUP_FILE $CONFIG_FILE"
```

### Environment Setup Script

```bash
#!/bin/bash
# setup-monitoring.sh - Set up monitoring environment

echo "🚀 Setting up Morgana monitoring..."

# Add shell aliases
cat >> ~/.bashrc << 'EOF'
# Morgana Protocol aliases
alias monitor-up='cd ~/.claude/morgana-protocol && make up'
alias monitor-attach='cd ~/.claude/morgana-protocol && make attach'
alias monitor-down='cd ~/.claude/morgana-protocol && make down'
alias monitor-status='cd ~/.claude/morgana-protocol && make status'
EOF

# Source agent wrapper
echo "source ~/.claude/scripts/agent-adapter-wrapper.sh" >> ~/.bashrc

echo "✅ Setup complete! Reload shell or run: source ~/.bashrc"
```

## Getting Help

### Documentation

- [Architecture Guide](ARCHITECTURE.md) - System design details
- [Event System](EVENT_SYSTEM.md) - Technical event system documentation
- [Getting Started](GETTING_STARTED.md) - Quick start guide
- [Integration Guide](INTEGRATION.md) - Advanced integration patterns

### Support Channels

- **GitHub Issues**: Report bugs or request features
- **Integration Tests**: Run `make test-integration` to validate setup
- **Performance Tests**: Run `make test-performance` for benchmarks

### Diagnostic Commands

```bash
# System diagnostic
make status                    # Check monitor status
make test-integration         # Validate system integration
make test-performance        # Performance validation

# Event system diagnostic
morgana -- --agent validation-expert --prompt "Test event generation"
```

## Summary

The migration to event stream architecture provides:

✅ **Simplified Setup**: Zero-configuration defaults ✅ **Better Performance**:
5M+ events/sec vs file I/O bottlenecks ✅ **Real-time Monitoring**: Beautiful
TUI with live event streaming ✅ **Improved Reliability**: Lock-free design,
circular buffering ✅ **Easier Maintenance**: Pure Go, no Python dependencies ✅
**Enhanced Debugging**: Event stream logs, comprehensive monitoring

The core agent execution interface remains unchanged, ensuring your existing
workflows continue to work while gaining powerful new monitoring capabilities.

**Next Steps**:

1. Complete the migration steps above
2. Explore the new TUI monitoring features
3. Customize themes and performance settings as needed
4. Set up shell aliases for convenient monitor control

Welcome to the new, simplified Morgana Protocol! 🚀
