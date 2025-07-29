#!/bin/bash

# Token-Efficient Tool Use Configuration Script
# This script manages the token-efficient tools beta feature for Claude Code

SETTINGS_FILE="$HOME/.claude/settings.json"

# Function to check if jq is installed
check_jq() {
    if ! command -v jq &> /dev/null; then
        echo "Error: jq is required but not installed. Install it with: brew install jq"
        exit 1
    fi
}

# Function to enable token-efficient mode
enable_token_efficient() {
    check_jq
    # Ensure env object exists
    if ! jq -e '.env' "$SETTINGS_FILE" > /dev/null 2>&1; then
        jq '. + {env: {}}' "$SETTINGS_FILE" > "$SETTINGS_FILE.tmp" && mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"
    fi
    # Set the environment variables
    jq '.env.CLAUDE_TOKEN_EFFICIENT_MODE = "true" | .env.CLAUDE_BETA_HEADER = "token-efficient-tools-2025-02-19"' "$SETTINGS_FILE" > "$SETTINGS_FILE.tmp" && mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"
    echo "✅ Token-efficient mode enabled"
}

# Function to disable token-efficient mode
disable_token_efficient() {
    check_jq
    # Ensure env object exists
    if ! jq -e '.env' "$SETTINGS_FILE" > /dev/null 2>&1; then
        jq '. + {env: {}}' "$SETTINGS_FILE" > "$SETTINGS_FILE.tmp" && mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"
    fi
    jq '.env.CLAUDE_TOKEN_EFFICIENT_MODE = "false"' "$SETTINGS_FILE" > "$SETTINGS_FILE.tmp" && mv "$SETTINGS_FILE.tmp" "$SETTINGS_FILE"
    echo "❌ Token-efficient mode disabled"
}

# Function to check status
check_status() {
    check_jq
    enabled=$(jq -r '.env.CLAUDE_TOKEN_EFFICIENT_MODE // "false"' "$SETTINGS_FILE")
    if [ "$enabled" = "true" ]; then
        echo "Token-efficient mode is: ENABLED ✅"
        echo "Beta header: $(jq -r '.env.CLAUDE_BETA_HEADER // "token-efficient-tools-2025-02-19"' "$SETTINGS_FILE")"
        echo ""
        echo "Note: Token-efficient mode works with Claude 3.7 Sonnet models"
        echo "      when used through the Anthropic API"
    else
        echo "Token-efficient mode is: DISABLED ❌"
    fi
}

# Main script logic
case "$1" in
    "enable")
        enable_token_efficient
        ;;
    "disable")
        disable_token_efficient
        ;;
    "status")
        check_status
        ;;
    *)
        echo "Usage: $0 {enable|disable|status}"
        echo ""
        echo "This script manages the token-efficient tools beta feature for Claude Code."
        echo "When enabled, it reduces output tokens by an average of 14% (up to 70%)."
        echo ""
        echo "Note: This feature is only compatible with Claude Sonnet 3.7 models."
        exit 1
        ;;
esac