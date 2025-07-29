#!/bin/bash

# Pre-API Call Hook for Token-Efficient Mode
# This hook checks if token-efficient mode is enabled and sets appropriate environment variables

SETTINGS_FILE="$HOME/.claude/settings.json"

# Check if jq is available
if command -v jq &> /dev/null; then
    # Check if token-efficient mode is enabled
    enabled=$(jq -r '.experimental.tokenEfficientTools.enabled // false' "$SETTINGS_FILE" 2>/dev/null)
    
    if [ "$enabled" = "true" ]; then
        # Get beta header value
        beta_header=$(jq -r '.experimental.tokenEfficientTools.betaHeader // "token-efficient-tools-2025-02-19"' "$SETTINGS_FILE" 2>/dev/null)
        
        # Export environment variables that Claude Code might use
        export CLAUDE_BETA_HEADER="$beta_header"
        export CLAUDE_TOKEN_EFFICIENT_MODE="true"
        
        # Log if verbose mode is set
        if [ -n "$CLAUDE_VERBOSE" ]; then
            echo "[Token-Efficient Mode] Beta header set: $beta_header" >&2
        fi
    fi
fi

# Continue with normal execution
exit 0