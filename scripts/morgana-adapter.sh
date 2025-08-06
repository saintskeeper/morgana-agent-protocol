#!/bin/bash
#
# Morgana Protocol Adapter Functions for QDIRECTOR
# Source this file to use AgentAdapter functions in shell/markdown
#

# Single agent execution using Morgana Protocol
function AgentAdapter() {
    local agent_type="$1"
    local prompt="$2"
    shift 2
    local additional_args="$@"
    
    # Check if morgana is available
    local morgana_cmd=""
    if [ -f "$HOME/.claude/bin/morgana" ]; then
        morgana_cmd="$HOME/.claude/bin/morgana"
    elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana" ]; then
        morgana_cmd="$HOME/.claude/morgana-protocol/dist/morgana"
    elif command -v morgana >/dev/null 2>&1; then
        morgana_cmd="morgana"
    else
        echo "Error: Morgana binary not found. Please run setup-local.sh" >&2
        return 1
    fi
    
    # Execute agent via Morgana
    echo "🤖 Executing agent: $agent_type" >&2
    "$morgana_cmd" -- --agent "$agent_type" --prompt "$prompt" $additional_args
}

# Parallel agent execution using Morgana Protocol
function AgentAdapterParallel() {
    local morgana_cmd=""
    if [ -f "$HOME/.claude/bin/morgana" ]; then
        morgana_cmd="$HOME/.claude/bin/morgana"
    elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana" ]; then
        morgana_cmd="$HOME/.claude/morgana-protocol/dist/morgana"
    elif command -v morgana >/dev/null 2>&1; then
        morgana_cmd="morgana"
    else
        echo "Error: Morgana binary not found. Please run setup-local.sh" >&2
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