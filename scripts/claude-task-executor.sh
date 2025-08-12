#!/bin/bash
#
# Claude Task Executor - Bridge between Morgana and Claude Code's Task tool
# This script captures morgana output and executes Task tool when needed
#

# Function to execute Task tool from within Claude Code
execute_claude_task() {
    local agent_type="$1"
    local prompt="$2"
    
    # This function will be called from within Claude Code
    # where the Task tool is available
    echo "CLAUDE_TASK_REQUEST:${agent_type}:${prompt}"
}

# Function to run morgana and capture mock output
run_morgana_with_task() {
    local agent_type="$1"
    local prompt="$2"
    
    # Run morgana in mock mode and capture output
    local morgana_output=$(morgana -- --agent "$agent_type" --prompt "$prompt" 2>&1)
    
    # Check if we got mock output
    if echo "$morgana_output" | grep -q "\[MOCK\]"; then
        # Signal Claude to execute real Task
        echo "===CLAUDE_TASK_REQUIRED==="
        echo "Agent: $agent_type"
        echo "Prompt: $prompt"
        echo "===END_CLAUDE_TASK==="
        
        # The actual Task execution will happen in Claude Code
        # This is just a placeholder for the request
        return 0
    else
        # Real output, just return it
        echo "$morgana_output"
    fi
}

# Export for use in Claude Code
export -f execute_claude_task
export -f run_morgana_with_task

# If called directly, show usage
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    echo "Claude Task Executor - Bridge for Morgana Protocol"
    echo ""
    echo "This script should be sourced in Claude Code:"
    echo "  source $0"
    echo ""
    echo "Then use the Task tool when you see CLAUDE_TASK_REQUEST markers"
fi