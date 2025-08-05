#!/bin/bash
#
# Morgana Adapter - Wrapper script for Claude Code integration
# This script bridges the AgentAdapter concept with the Morgana Protocol
#

MORGANA_BIN="$HOME/.claude/morgana-protocol/dist/morgana"

# Check if Morgana is built
if [ ! -f "$MORGANA_BIN" ]; then
    echo "Error: Morgana Protocol not found at $MORGANA_BIN"
    echo "Please build it first: cd ~/.claude/morgana-protocol && make dev"
    exit 1
fi

# Function to execute parallel tasks
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
                    exit 1
                fi
                ;;
            *)
                echo "Unknown argument: $1"
                exit 1
                ;;
        esac
    done
    
    tasks_json+="]"
    
    # Execute with Morgana
    echo "$tasks_json" | "$MORGANA_BIN" --parallel
}

# Function to execute single task
run_single_agent() {
    local agent_type="$1"
    local prompt="$2"
    
    "$MORGANA_BIN" --agent "$agent_type" --prompt "$prompt"
}

# Main execution
if [ $# -eq 0 ]; then
    echo "Morgana Adapter - Agent orchestration for Claude Code"
    echo ""
    echo "Usage:"
    echo "  Single agent:    $0 <agent-type> <prompt>"
    echo "  Parallel agents: $0 --parallel --agent <type> --prompt <prompt> [--agent <type> --prompt <prompt>...]"
    echo ""
    echo "Available agents:"
    echo "  - code-implementer"
    echo "  - sprint-planner"
    echo "  - test-specialist"
    echo "  - validation-expert"
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