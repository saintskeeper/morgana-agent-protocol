#!/bin/bash

# Claude Code Wrapper Script for Token-Efficient Mode
# This wrapper intercepts Claude Code calls and adds beta headers when needed

SETTINGS_FILE="$HOME/.claude/settings.json"
ORIGINAL_CLAUDE="/opt/homebrew/bin/claude"

# Function to check if token-efficient mode is enabled
is_token_efficient_enabled() {
    if command -v jq &> /dev/null; then
        enabled=$(jq -r '.experimental.tokenEfficientTools.enabled // false' "$SETTINGS_FILE" 2>/dev/null)
        [ "$enabled" = "true" ]
    else
        return 1
    fi
}

# Function to get beta header
get_beta_header() {
    if command -v jq &> /dev/null; then
        jq -r '.experimental.tokenEfficientTools.betaHeader // "token-efficient-tools-2025-02-19"' "$SETTINGS_FILE" 2>/dev/null
    else
        echo "token-efficient-tools-2025-02-19"
    fi
}

# Check if token-efficient mode should be activated
if is_token_efficient_enabled; then
    # Set environment variables for beta mode
    export ANTHROPIC_BETA_HEADER=$(get_beta_header)
    export CLAUDE_TOKEN_EFFICIENT_MODE="true"
    
    # Show notification if verbose
    if [ -n "$CLAUDE_VERBOSE" ] || [ -n "$DEBUG" ]; then
        echo "[Token-Efficient Mode] Activated with header: $ANTHROPIC_BETA_HEADER" >&2
    fi
fi

# Execute the original Claude command with all arguments
exec "$ORIGINAL_CLAUDE" "$@"