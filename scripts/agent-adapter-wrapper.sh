#!/bin/bash
# Agent Adapter Wrapper - Sources Morgana Protocol adapter functions

# Source the Morgana adapter functions
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$SCRIPT_DIR/morgana-adapter.sh"

# The morgana-adapter.sh provides:
# - AgentAdapter() function for single agent execution
# - AgentAdapterParallel() function for parallel execution
# - morgana_parallel() helper function

# If called directly, show usage or execute
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    if [ $# -lt 2 ]; then
        echo "Agent Adapter Wrapper - Morgana Protocol Integration"
        echo ""
        echo "Usage:"
        echo "  Single agent:    $0 <agent-type> <prompt>"
        echo "  Source mode:     source $0  # To use AgentAdapter functions"
        echo ""
        echo "Available agents:"
        echo "  - code-implementer"
        echo "  - sprint-planner"
        echo "  - test-specialist"
        echo "  - validation-expert"
        echo ""
        echo "Examples:"
        echo "  $0 code-implementer \"implement user service\""
        echo "  source $0  # Then use: AgentAdapter \"test-specialist\" \"create tests\""
        exit 1
    fi
    
    # Execute single agent
    AgentAdapter "$1" "$2"
fi