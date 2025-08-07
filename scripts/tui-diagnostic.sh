#!/bin/bash
#
# TUI Diagnostic Script
# Helps diagnose TUI performance and display issues
#

echo "🔍 Morgana TUI Diagnostic Report"
echo "================================"
echo ""

# Check monitor status
echo "📊 Monitor Status:"
if pgrep -f morgana-monitor > /dev/null; then
    PID=$(pgrep -f morgana-monitor | head -1)
    echo "✅ Running (PID: $PID)"
    
    # Get CPU and memory usage
    ps -p $PID -o pid,ppid,%cpu,%mem,etime,command | tail -1
else
    echo "❌ Not running"
fi
echo ""

# Check socket
echo "🔌 IPC Socket:"
if [ -S /tmp/morgana.sock ]; then
    echo "✅ Active"
    ls -la /tmp/morgana.sock
else
    echo "❌ Not found"
fi
echo ""

# Check event bus stats (if available)
echo "📈 Event Statistics:"
if [ -f /tmp/morgana-monitor.log ]; then
    echo "Recent events (last 10):"
    tail -10 /tmp/morgana-monitor.log | grep -E "Event|Task|Agent"
else
    echo "No log file found"
fi
echo ""

# Test event sending
echo "🧪 Testing Event Flow:"
echo "Sending test event..."

# Create a simple test
TEST_OUTPUT=$(morgana -- --agent sprint-planner --prompt "diagnostic test" 2>&1)

if echo "$TEST_OUTPUT" | grep -q "Connected to Morgana Monitor"; then
    echo "✅ IPC connection successful"
else
    echo "❌ IPC connection failed"
fi

if echo "$TEST_OUTPUT" | grep -q "Disconnected from Morgana Monitor"; then
    echo "✅ Clean disconnect"
else
    echo "⚠️ Disconnect not confirmed"
fi
echo ""

# Performance recommendations
echo "💡 Performance Tuning Options:"
echo "1. Reduce refresh rate: Change from 16ms to 33ms in morgana.yaml"
echo "2. Disable debug mode: unset MORGANA_DEBUG"
echo "3. Reduce max log lines: Set to 1000 instead of 10000"
echo "4. Use performance mode: morgana-monitor --performance"
echo ""

# Display current config
echo "⚙️ Current Configuration:"
if [ -f ~/.claude/morgana-protocol/morgana.yaml ]; then
    echo "TUI Settings from morgana.yaml:"
    grep -A 5 "tui:" ~/.claude/morgana-protocol/morgana.yaml | grep -E "enabled|refresh_rate|max_log_lines"
else
    echo "Config file not found"
fi
echo ""

echo "📋 Diagnostic complete!"
echo ""
echo "To share with support:"
echo "1. Run: $0 > /tmp/tui-diagnostic.txt"
echo "2. Take a screenshot of the TUI issue"
echo "3. Describe: What you expected vs what you see"