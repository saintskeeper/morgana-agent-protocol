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
    
    # Check screen session first
    if command -v screen >/dev/null 2>&1 && screen -list | grep -q morgana-monitor; then
        # Get screen session PID and update PID file
        pid=$(screen -list | grep morgana-monitor | awk '{print $1}' | cut -d. -f1)
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            echo "$pid" > "$PID_FILE"
            return 0
        fi
    fi
    
    # Check tmux session
    if command -v tmux >/dev/null 2>&1 && tmux has-session -t morgana-monitor 2>/dev/null; then
        # Get tmux session PID and update PID file
        pid=$(tmux list-sessions -F '#{session_name} #{pane_pid}' | grep morgana-monitor | awk '{print $2}')
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            echo "$pid" > "$PID_FILE"
            return 0
        fi
    fi
    
    # Check if PID file exists and process is running
    if [ -f "$PID_FILE" ]; then
        pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            # Process is running, check if socket is responsive
            if [ -S "$SOCKET_PATH" ]; then
                return 0
            fi
        fi
        # Clean up stale PID file
        rm -f "$PID_FILE"
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
    
    # Remove old log if it's too large (>10MB)
    if [ -f "$LOG_FILE" ] && [ $(stat -f%z "$LOG_FILE" 2>/dev/null || echo 0) -gt 10485760 ]; then
        echo "📝 Rotating large log file"
        mv "$LOG_FILE" "${LOG_FILE}.old"
    fi
    
    # Start monitor with proper terminal support
    local monitor_pid=""
    if command -v screen >/dev/null 2>&1; then
        # CRITICAL: Always force headless mode when using screen to prevent TUI attempts
        # The --headless flag MUST be present to avoid terminal detection issues
        screen -dmS morgana-monitor bash -c "exec $monitor_cmd --headless >> $LOG_FILE 2>&1"
        # Get screen session PID
        monitor_pid=$(screen -list | grep morgana-monitor | awk '{print $1}' | cut -d. -f1)
        echo "🗺 Started in screen session (headless mode)"
        echo "📺 To view TUI: make attach or screen -r morgana-monitor"
    elif command -v tmux >/dev/null 2>&1; then
        # Use tmux for proper terminal emulation with headless mode
        tmux new-session -d -s morgana-monitor "$monitor_cmd --headless"
        # Get tmux session PID
        monitor_pid=$(tmux list-sessions -F '#{session_name} #{pane_pid}' | grep morgana-monitor | awk '{print $2}')
        echo "🗺 Started in tmux session (headless mode)"
        echo "📺 To view TUI: make attach or tmux attach -t morgana-monitor"
    else
        # Fallback approaches
        if command -v script >/dev/null 2>&1; then
            echo "📝 Using script for pseudo-terminal"
            script -q "$LOG_FILE" "$monitor_cmd" &
            monitor_pid=$!
        else
            echo "⚠️  No screen/tmux available, using TERM workaround"
            TERM=xterm-256color nohup "$monitor_cmd" >> "$LOG_FILE" 2>&1 &
            monitor_pid=$!
        fi
    fi
    
    # Store PID
    if [ -n "$monitor_pid" ]; then
        echo $monitor_pid > "$PID_FILE"
    fi
    
    # Wait for socket to be created (increased timeout for TUI initialization)
    local wait_count=0
    local max_wait=20  # 10 seconds total (20 * 0.5s) - reduced from 30s
    echo "⏳ Waiting for socket creation..."
    
    while [ $wait_count -lt $max_wait ]; do
        if [ -S "$SOCKET_PATH" ]; then
            echo "✅ Morgana monitor started successfully (PID: $monitor_pid)"
            echo "📋 Socket: $SOCKET_PATH"
            if command -v screen >/dev/null 2>&1 && screen -list | grep -q morgana-monitor; then
                echo "🗺️ Session: screen -r morgana-monitor (to view TUI)"
            elif command -v tmux >/dev/null 2>&1 && tmux has-session -t morgana-monitor 2>/dev/null; then
                echo "🗺️ Session: tmux attach -t morgana-monitor (to view TUI)"
            else
                echo "📝 Logs: $LOG_FILE"
            fi
            return 0
        fi
        
        # Check if process is still running
        if [ -n "$monitor_pid" ] && ! kill -0 "$monitor_pid" 2>/dev/null; then
            echo "❌ Monitor process died unexpectedly"
            echo "📝 Check logs at: $LOG_FILE"
            tail -n 20 "$LOG_FILE" 2>/dev/null
            rm -f "$PID_FILE"
            return 1
        fi
        
        # Show progress indicator every 2 seconds
        if [ $((wait_count % 4)) -eq 0 ] && [ $wait_count -gt 0 ]; then
            echo "   ⏳ Still waiting... ($((wait_count / 2))s elapsed)"
            # Check if monitor is actually running but socket creation failed
            if pgrep -f "morgana-monitor.*--headless" >/dev/null 2>&1; then
                echo "   ℹ️  Monitor process is running, but socket not created yet"
            fi
        fi
        
        sleep 0.5
        wait_count=$((wait_count + 1))
    done
    
    echo "⚠️  Socket creation timeout - attempting direct start..."
    
    # Try starting directly without screen/tmux as fallback
    echo "🔄 Falling back to direct execution..."
    pkill -f "morgana-monitor.*--headless" 2>/dev/null
    sleep 1
    
    # Start directly in background
    nohup "$monitor_cmd" --headless >> "$LOG_FILE" 2>&1 &
    monitor_pid=$!
    echo $monitor_pid > "$PID_FILE"
    
    # Give it a moment to start
    sleep 2
    
    if [ -S "$SOCKET_PATH" ]; then
        echo "✅ Monitor started successfully with fallback method (PID: $monitor_pid)"
        echo "📋 Socket: $SOCKET_PATH"
        echo "📝 Logs: $LOG_FILE"
        return 0
    else
        echo "❌ Failed to start morgana monitor"
        echo "📝 Last log entries:"
        tail -n 20 "$LOG_FILE" 2>/dev/null
        kill "$monitor_pid" 2>/dev/null
        rm -f "$PID_FILE"
        return 1
    fi
}

# Stop the monitor daemon
stop_monitor() {
    if ! is_monitor_running; then
        echo "ℹ️  Morgana monitor is not running"
        # Clean up any leftover sessions
        if command -v screen >/dev/null 2>&1 && screen -list | grep -q morgana-monitor; then
            screen -S morgana-monitor -X quit 2>/dev/null
        fi
        if command -v tmux >/dev/null 2>&1 && tmux has-session -t morgana-monitor 2>/dev/null; then
            tmux kill-session -t morgana-monitor 2>/dev/null
        fi
        rm -f "$PID_FILE" "$SOCKET_PATH"
        return 0
    fi
    
    local pid=$(cat "$PID_FILE")
    echo "🛑 Stopping morgana monitor (PID: $pid)..."
    
    # Stop screen session if it exists
    if command -v screen >/dev/null 2>&1 && screen -list | grep -q morgana-monitor; then
        echo "🗺 Stopping screen session..."
        screen -S morgana-monitor -X quit
    # Stop tmux session if it exists
    elif command -v tmux >/dev/null 2>&1 && tmux has-session -t morgana-monitor 2>/dev/null; then
        echo "🗺 Stopping tmux session..."
        tmux kill-session -t morgana-monitor
    else
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

# Attach to monitor logs or TUI session
attach_monitor() {
    if ! is_monitor_running; then
        echo "❌ Morgana monitor is not running"
        return 1
    fi
    
    # Try to attach to screen session first
    if command -v screen >/dev/null 2>&1 && screen -list | grep -q morgana-monitor; then
        echo "🗺️ Attaching to screen session with TUI..."
        echo "ℹ️  Use Ctrl+A, D to detach from session"
        screen -r morgana-monitor
        return 0
    fi
    
    # Try to attach to tmux session
    if command -v tmux >/dev/null 2>&1 && tmux has-session -t morgana-monitor 2>/dev/null; then
        echo "🗺️ Attaching to tmux session with TUI..."
        echo "ℹ️  Use Ctrl+B, D to detach from session"
        tmux attach -t morgana-monitor
        return 0
    fi
    
    # Fallback to log tailing if no session available
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