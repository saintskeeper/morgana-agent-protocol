#!/bin/bash
#
# Morgana Monitor Control Script
# Provides start/stop/status/attach commands for the morgana-monitor daemon
#

SCRIPT_NAME="$(basename "$0")"
SOCKET_PATH="/tmp/morgana.sock"
PID_FILE="/tmp/morgana-monitor.pid"
LOG_FILE="/tmp/morgana-monitor.log"

# Find morgana-monitor binary (prioritize /usr/local/bin)
find_monitor_binary() {
    local monitor_cmd=""
    
    if command -v morgana-monitor >/dev/null 2>&1; then
        monitor_cmd="morgana-monitor"  # Will use /usr/local/bin/morgana-monitor if in PATH
    elif [ -f "/usr/local/bin/morgana-monitor" ]; then
        monitor_cmd="/usr/local/bin/morgana-monitor"
    elif [ -f "$HOME/.claude/bin/morgana-monitor" ]; then
        monitor_cmd="$HOME/.claude/bin/morgana-monitor"
    elif [ -f "$HOME/.claude/morgana-protocol/morgana-monitor" ]; then
        monitor_cmd="$HOME/.claude/morgana-protocol/morgana-monitor"
    elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana-monitor" ]; then
        monitor_cmd="$HOME/.claude/morgana-protocol/dist/morgana-monitor"
    else
        echo "Error: morgana-monitor binary not found" >&2
        echo "Please run 'make install' from ~/.claude directory" >&2
        return 1
    fi
    
    echo "$monitor_cmd"
    return 0
}

# Check if monitor is running
is_monitor_running() {
    local pid=""
    
    # Check if PID file exists and process is running
    if [ -f "$PID_FILE" ]; then
        pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            # Process is running
            return 0
        fi
        # Clean up stale PID file
        rm -f "$PID_FILE"
    fi
    
    # Check for any morgana-monitor process
    pid=$(pgrep -f "morgana-monitor.*-headless" | head -1)
    if [ -n "$pid" ]; then
        # Found a running monitor, update PID file
        echo "$pid" > "$PID_FILE"
        return 0
    fi
    
    # Check for socket without PID file (manual start)
    if [ -S "$SOCKET_PATH" ]; then
        # Try to find process by socket
        local socket_pid=$(lsof -t "$SOCKET_PATH" 2>/dev/null | head -1)
        if [ -n "$socket_pid" ]; then
            echo "$socket_pid" > "$PID_FILE"
            return 0
        else
            # Socket exists but no process, clean it up
            rm -f "$SOCKET_PATH"
        fi
    fi
    
    return 1
}

# Start the monitor daemon
start_monitor() {
    if is_monitor_running; then
        echo "✅ Morgana monitor is already running"
        return 0
    fi
    
    local monitor_cmd
    if ! monitor_cmd=$(find_monitor_binary); then
        return 1
    fi
    
    echo "🚀 Starting morgana monitor daemon..."
    
    # Clean up any stale socket/pid files
    rm -f "$SOCKET_PATH" "$PID_FILE"
    
    # Remove old log if it's too large (>10MB)
    if [ -f "$LOG_FILE" ] && [ $(stat -f%z "$LOG_FILE" 2>/dev/null || echo 0) -gt 10485760 ]; then
        echo "📝 Rotating large log file"
        mv "$LOG_FILE" "${LOG_FILE}.old"
    fi
    
    # Start monitor directly with nohup (simpler and more reliable than screen/tmux)
    echo "🔧 Starting monitor in background (headless mode)..."
    nohup "$monitor_cmd" -headless >> "$LOG_FILE" 2>&1 &
    local monitor_pid=$!
    
    # Store PID immediately
    echo $monitor_pid > "$PID_FILE"
    
    # Check if process started successfully
    sleep 0.5
    if ! kill -0 "$monitor_pid" 2>/dev/null; then
        echo "❌ Failed to start monitor process"
        echo "📝 Check logs at: $LOG_FILE"
        tail -n 10 "$LOG_FILE" 2>/dev/null
        rm -f "$PID_FILE"
        return 1
    fi
    
    # Wait for socket to be created (with shorter timeout)
    local wait_count=0
    local max_wait=10  # 5 seconds total (10 * 0.5s)
    echo "⏳ Waiting for socket creation..."
    
    while [ $wait_count -lt $max_wait ]; do
        if [ -S "$SOCKET_PATH" ]; then
            echo "✅ Morgana monitor started successfully (PID: $monitor_pid)"
            echo "📋 Socket: $SOCKET_PATH"
            echo "📝 Logs: tail -f $LOG_FILE"
            return 0
        fi
        
        # Check if process is still running
        if ! kill -0 "$monitor_pid" 2>/dev/null; then
            echo "❌ Monitor process died unexpectedly"
            echo "📝 Last log entries:"
            tail -n 20 "$LOG_FILE" 2>/dev/null
            rm -f "$PID_FILE"
            return 1
        fi
        
        sleep 0.5
        wait_count=$((wait_count + 1))
    done
    
    # If socket still doesn't exist but process is running, it might be okay
    if kill -0 "$monitor_pid" 2>/dev/null; then
        echo "⚠️  Monitor is running (PID: $monitor_pid) but socket not yet created"
        echo "📝 The monitor may still be initializing. Check logs: tail -f $LOG_FILE"
        echo "📋 Expected socket path: $SOCKET_PATH"
        return 0
    else
        echo "❌ Failed to start morgana monitor"
        echo "📝 Last log entries:"
        tail -n 20 "$LOG_FILE" 2>/dev/null
        rm -f "$PID_FILE"
        return 1
    fi
}

# Stop the monitor daemon
stop_monitor() {
    if ! is_monitor_running; then
        echo "ℹ️  Morgana monitor is not running"
        # Clean up any stale files
        rm -f "$PID_FILE" "$SOCKET_PATH"
        return 0
    fi
    
    local pid=$(cat "$PID_FILE")
    echo "🛑 Stopping morgana monitor (PID: $pid)..."
    
    # Send TERM signal for graceful shutdown
    if kill -TERM "$pid" 2>/dev/null; then
        # Wait for graceful shutdown
        local wait_count=0
        while [ $wait_count -lt 10 ]; do
            if ! kill -0 "$pid" 2>/dev/null; then
                break
            fi
            sleep 0.5
            wait_count=$((wait_count + 1))
        done
        
        # Force kill if still running
        if kill -0 "$pid" 2>/dev/null; then
            echo "⚡ Force stopping monitor..."
            kill -KILL "$pid" 2>/dev/null
        fi
    fi
    
    # Clean up files
    rm -f "$PID_FILE" "$SOCKET_PATH"
    echo "✅ Morgana monitor stopped"
}

# Show monitor status
status_monitor() {
    if is_monitor_running; then
        local pid=$(cat "$PID_FILE")
        echo "✅ Morgana monitor is running (PID: $pid)"
        echo "📋 Socket: $SOCKET_PATH"
        echo "📝 Logs: $LOG_FILE"
        
        # Show basic process info
        if command -v ps >/dev/null 2>&1; then
            echo "📊 Process info:"
            ps -p "$pid" -o pid,ppid,etime,cmd 2>/dev/null || echo "   Process details unavailable"
        fi
        
        return 0
    else
        echo "❌ Morgana monitor is not running"
        return 1
    fi
}

# Attach to monitor logs
attach_monitor() {
    if ! is_monitor_running; then
        echo "❌ Morgana monitor is not running"
        return 1
    fi
    
    if [ ! -f "$LOG_FILE" ]; then
        echo "❌ Log file not found: $LOG_FILE"
        return 1
    fi
    
    echo "📝 Monitoring logs (Ctrl+C to exit):"
    echo "📋 Log file: $LOG_FILE"
    echo "$(printf '%.0s-' {1..50})"
    
    # Follow the log file
    tail -f "$LOG_FILE"
}

# Show usage
show_usage() {
    echo "Morgana Monitor Control - Daemon management for monitoring"
    echo ""
    echo "Usage: $SCRIPT_NAME <command>"
    echo ""
    echo "Commands:"
    echo "  start     Start the morgana-monitor daemon"
    echo "  stop      Stop the morgana-monitor daemon"
    echo "  restart   Restart the morgana-monitor daemon"
    echo "  status    Show daemon status"
    echo "  attach    Attach to monitor logs (live tail)"
    echo "  logs      Show recent log output"
    echo ""
    echo "Files:"
    echo "  PID file: $PID_FILE"
    echo "  Socket:   $SOCKET_PATH"
    echo "  Logs:     $LOG_FILE"
    echo ""
    echo "The monitor provides a TUI interface and IPC socket for agent monitoring."
}

# Show recent logs
show_logs() {
    if [ ! -f "$LOG_FILE" ]; then
        echo "❌ Log file not found: $LOG_FILE"
        return 1
    fi
    
    echo "📝 Recent morgana-monitor logs:"
    echo "📋 Log file: $LOG_FILE"
    echo "$(printf '%.0s-' {1..50})"
    
    # Show last 50 lines
    tail -n 50 "$LOG_FILE"
}

# Main execution
case "${1:-}" in
    start)
        start_monitor
        ;;
    stop)
        stop_monitor
        ;;
    restart)
        stop_monitor
        sleep 1
        start_monitor
        ;;
    status)
        status_monitor
        ;;
    attach)
        attach_monitor
        ;;
    logs)
        show_logs
        ;;
    -h|--help|help|"")
        show_usage
        ;;
    *)
        echo "❌ Unknown command: $1"
        echo ""
        show_usage
        exit 1
        ;;
esac