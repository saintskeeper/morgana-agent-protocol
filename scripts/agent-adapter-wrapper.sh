#!/bin/bash
# Direct wrapper for Morgana Protocol - no Python needed!

MORGANA_BIN="/Users/walterday/.claude/morgana-protocol/dist/morgana"

# Function that mimics AgentAdapter for qdirector-enhanced.md
AgentAdapter() {
    local agent_type="$1"
    local prompt="$2"
    
    # Create JSON task
    local task_json="[{\"agent_type\":\"$agent_type\",\"prompt\":\"$prompt\"}]"
    
    # Call Morgana and extract output
    echo "$task_json" | "$MORGANA_BIN" | jq -r '.results[0].output'
}

# Export the function so it can be used by other scripts
export -f AgentAdapter

# If called directly, execute the agent
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    if [ $# -lt 2 ]; then
        echo "Usage: $0 <agent-type> <prompt>"
        exit 1
    fi
    
    AgentAdapter "$1" "$2"
fi