#!/bin/bash

# Pre-API Call Hook for Token-Efficient Mode
# This hook checks if token-efficient mode is enabled and sets appropriate environment variables
# NOTE: Token-efficient mode only works with Claude 3.7 Sonnet, not with Claude 4 models

SETTINGS_FILE="$HOME/.claude/settings.json"

# Compatible models for token-efficient mode
COMPATIBLE_MODELS=("claude-3-7-sonnet" "claude-3.7-sonnet")

# Check if jq is available
if command -v jq &> /dev/null; then
    # Check if token-efficient mode is enabled
    enabled=$(jq -r '.env.CLAUDE_TOKEN_EFFICIENT_MODE // "false"' "$SETTINGS_FILE" 2>/dev/null)
    
    if [ "$enabled" = "true" ]; then
        # Get beta header value
        beta_header=$(jq -r '.env.CLAUDE_BETA_HEADER // "token-efficient-tools-2025-02-19"' "$SETTINGS_FILE" 2>/dev/null)
        
        # Check if current model is compatible (if model info is available)
        model_compatible=true
        if [ -n "$CLAUDE_MODEL" ]; then
            model_compatible=false
            for compatible_model in "${COMPATIBLE_MODELS[@]}"; do
                if [[ "$CLAUDE_MODEL" == *"$compatible_model"* ]]; then
                    model_compatible=true
                    break
                fi
            done
            
            # Check for Claude 4 models (Opus/Sonnet) which are incompatible
            if [[ "$CLAUDE_MODEL" == *"claude-4"* ]] || [[ "$CLAUDE_MODEL" == *"opus"* ]]; then
                model_compatible=false
            fi
        fi
        
        # Check for disable_parallel_tool_use incompatibility
        if [ "$CLAUDE_DISABLE_PARALLEL_TOOL_USE" = "true" ]; then
            model_compatible=false
            if [ -n "$CLAUDE_VERBOSE" ]; then
                echo "[Token-Efficient Mode] WARNING: Disabled due to disable_parallel_tool_use" >&2
            fi
        fi
        
        if [ "$model_compatible" = "true" ]; then
            # Export environment variables that Claude Code might use
            export CLAUDE_BETA_HEADER="$beta_header"
            export CLAUDE_TOKEN_EFFICIENT_MODE="true"
            
            # Log if verbose mode is set
            if [ -n "$CLAUDE_VERBOSE" ]; then
                echo "[Token-Efficient Mode] Beta header set: $beta_header" >&2
                [ -n "$CLAUDE_MODEL" ] && echo "[Token-Efficient Mode] Model: $CLAUDE_MODEL" >&2
            fi
        elif [ -n "$CLAUDE_VERBOSE" ] && [ -n "$CLAUDE_MODEL" ]; then
            echo "[Token-Efficient Mode] Not applied - incompatible model: $CLAUDE_MODEL" >&2
        fi
    fi
fi

# Continue with normal execution
exit 0