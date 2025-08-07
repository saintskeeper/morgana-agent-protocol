#!/bin/bash
#
# Morgana Protocol Adapter Functions for QDIRECTOR
# Source this file to use AgentAdapter functions in shell/markdown
#

# Function to ensure morgana-monitor is running
function ensure_morgana_monitor() {
    local socket_path="/tmp/morgana.sock"
    local monitor_cmd=""
    
    # Check if monitor socket exists and is active
    if [ -S "$socket_path" ]; then
        # Check if there's a process using the socket
        if lsof "$socket_path" >/dev/null 2>&1; then
            echo "🔍 Morgana monitor already running" >&2
            return 0
        else
            # Socket exists but no process, clean it up
            echo "🧹 Cleaning up stale socket" >&2
            rm -f "$socket_path" 2>/dev/null
        fi
    fi
    
    # Find morgana-monitor binary
    if [ -f "$HOME/.claude/bin/morgana-monitor" ]; then
        monitor_cmd="$HOME/.claude/bin/morgana-monitor"
    elif [ -f "$HOME/.claude/morgana-protocol/morgana-monitor" ]; then
        monitor_cmd="$HOME/.claude/morgana-protocol/morgana-monitor"
    elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana-monitor" ]; then
        monitor_cmd="$HOME/.claude/morgana-protocol/dist/morgana-monitor"
    elif command -v morgana-monitor >/dev/null 2>&1; then
        monitor_cmd="morgana-monitor"
    else
        echo "⚠️  Morgana monitor not found, continuing without monitoring" >&2
        return 1
    fi
    
    # Start monitor in background with pseudo-terminal support
    echo "🚀 Starting morgana monitor daemon..." >&2
    
    # Try different approaches to start monitor in background
    local monitor_pid=""
    if command -v screen >/dev/null 2>&1; then
        # Use screen for proper terminal emulation
        screen -dmS morgana-monitor "$monitor_cmd"
        sleep 1  # Give screen time to start
        # Get screen session PID
        monitor_pid=$(screen -list | grep morgana-monitor | awk '{print $1}' | cut -d. -f1 2>/dev/null)
    elif command -v tmux >/dev/null 2>&1; then
        # Use tmux for proper terminal emulation
        tmux new-session -d -s morgana-monitor "$monitor_cmd"
        sleep 1  # Give tmux time to start
        # Get tmux session PID
        monitor_pid=$(tmux list-sessions -F '#{session_name} #{pane_pid}' 2>/dev/null | grep morgana-monitor | awk '{print $2}')
    else
        # Fallback: use script to create a pseudo-terminal
        if command -v script >/dev/null 2>&1; then
            # Use script with null output to create pty
            script -q /dev/null "$monitor_cmd" > /tmp/morgana-monitor.log 2>&1 &
            monitor_pid=$!
        else
            # Last resort: try with TERM set and background
            TERM=xterm-256color nohup "$monitor_cmd" > /tmp/morgana-monitor.log 2>&1 &
            monitor_pid=$!
        fi
    fi
    
    # Store PID for cleanup if we have one
    if [ -n "$monitor_pid" ]; then
        echo $monitor_pid > /tmp/morgana-monitor.pid
    fi
    
    # Wait briefly for monitor to initialize
    local wait_count=0
    while [ $wait_count -lt 10 ]; do
        if [ -S "$socket_path" ]; then
            echo "✅ Morgana monitor started (PID: ${monitor_pid:-unknown})" >&2
            if command -v screen >/dev/null 2>&1 && screen -list | grep -q morgana-monitor; then
                echo "🗺 Session: screen -r morgana-monitor" >&2
            elif command -v tmux >/dev/null 2>&1 && tmux has-session -t morgana-monitor 2>/dev/null; then
                echo "🗺 Session: tmux attach -t morgana-monitor" >&2
            fi
            return 0
        fi
        sleep 0.5
        wait_count=$((wait_count + 1))
    done
    
    echo "⚠️  Morgana monitor may not have started correctly" >&2
    return 1
}

# Single agent execution using Morgana Protocol
function AgentAdapter() {
    local agent_type="$1"
    local prompt="$2"
    shift 2
    local additional_args="$@"
    
    # Ensure monitor is running before executing agent
    ensure_morgana_monitor
    
    # Check if morgana is available (prioritize /usr/local/bin)
    local morgana_cmd=""
    if command -v morgana >/dev/null 2>&1; then
        morgana_cmd="morgana"  # Will use /usr/local/bin/morgana if in PATH
    elif [ -f "/usr/local/bin/morgana" ]; then
        morgana_cmd="/usr/local/bin/morgana"
    elif [ -f "$HOME/.claude/bin/morgana" ]; then
        morgana_cmd="$HOME/.claude/bin/morgana"
    elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana" ]; then
        morgana_cmd="$HOME/.claude/morgana-protocol/dist/morgana"
    else
        echo "Error: Morgana binary not found. Please run 'make install' from ~/.claude" >&2
        return 1
    fi
    
    # Execute agent via Morgana
    echo "🤖 Executing agent: $agent_type" >&2
    "$morgana_cmd" -- --agent "$agent_type" --prompt "$prompt" $additional_args
}

# Parallel agent execution using Morgana Protocol
function AgentAdapterParallel() {
    # Ensure monitor is running before executing agents
    ensure_morgana_monitor
    
    local morgana_cmd=""
    if command -v morgana >/dev/null 2>&1; then
        morgana_cmd="morgana"  # Will use /usr/local/bin/morgana if in PATH
    elif [ -f "/usr/local/bin/morgana" ]; then
        morgana_cmd="/usr/local/bin/morgana"
    elif [ -f "$HOME/.claude/bin/morgana" ]; then
        morgana_cmd="$HOME/.claude/bin/morgana"
    elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana" ]; then
        morgana_cmd="$HOME/.claude/morgana-protocol/dist/morgana"
    else
        echo "Error: Morgana binary not found. Please run 'make install' from ~/.claude" >&2
        return 1
    fi
    
    echo "🚀 Executing agents in parallel via Morgana Protocol" >&2
    "$morgana_cmd" --parallel
}

# Helper function to create parallel task JSON
function morgana_parallel() {
    # Usage: morgana_parallel << 'EOF'
    # [
    #   {"agent_type": "code-implementer", "prompt": "implement feature"},
    #   {"agent_type": "test-specialist", "prompt": "write tests"}
    # ]
    # EOF
    
    AgentAdapterParallel
}

# Legacy wrapper functions for backward compatibility
run_parallel_agents() {
    local tasks_json="["
    local first=true
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --agent)
                if [ "$first" = false ]; then
                    tasks_json+=","
                fi
                first=false
                
                agent_type="$2"
                shift 2
                
                if [ "$1" = "--prompt" ]; then
                    prompt="$2"
                    shift 2
                    tasks_json+="{\"agent_type\":\"$agent_type\",\"prompt\":\"$prompt\"}"
                else
                    echo "Error: --agent must be followed by --prompt"
                    return 1
                fi
                ;;
            *)
                echo "Unknown argument: $1"
                return 1
                ;;
        esac
    done
    
    tasks_json+="]"
    
    # Execute with Morgana
    echo "$tasks_json" | AgentAdapterParallel
}

run_single_agent() {
    AgentAdapter "$1" "$2"
}

# Export functions for use in subshells
export -f ensure_morgana_monitor
export -f AgentAdapter
export -f AgentAdapterParallel
export -f morgana_parallel

# Main execution when script is run directly (not sourced)
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    if [ $# -eq 0 ]; then
        echo "Morgana Protocol Adapter - Agent orchestration for Claude Code"
        echo ""
        echo "Usage:"
        echo "  Source mode:     source $0"
        echo "  Single agent:    $0 <agent-type> <prompt>"
        echo "  Parallel agents: $0 --parallel --agent <type> --prompt <prompt> [--agent <type> --prompt <prompt>...]"
        echo ""
        echo "Available agents:"
        echo "  - code-implementer"
        echo "  - sprint-planner"
        echo "  - test-specialist"
        echo "  - validation-expert"
        echo ""
        echo "When sourced, provides functions:"
        echo "  - AgentAdapter: Execute single agent"
        echo "  - AgentAdapterParallel: Execute multiple agents in parallel"
        echo "  - morgana_parallel: Helper for parallel execution"
        exit 0
    fi
    
    # Check for parallel execution
    if [ "$1" = "--parallel" ]; then
        shift
        run_parallel_agents "$@"
    else
        # Single agent execution
        if [ $# -lt 2 ]; then
            echo "Error: Single agent mode requires <agent-type> and <prompt>"
            exit 1
        fi
        run_single_agent "$1" "$2"
    fi
else
    # Script is being sourced
    echo "🧙‍♂️ Morgana Protocol adapter functions loaded" >&2
    echo "   - AgentAdapter: Execute single agent" >&2
    echo "   - AgentAdapterParallel: Execute multiple agents in parallel" >&2
    echo "   - morgana_parallel: Helper for parallel execution" >&2
fi