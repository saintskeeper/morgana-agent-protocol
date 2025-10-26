#!/bin/bash
# Claude Code Hook: post-write
# Intelligent logging agent integration
# Automatically adds logging to code files using Ollama Cloud AI

# Read JSON input from stdin
JSON_INPUT=$(cat)

# Extract file path from JSON using grep and sed
EDITED_FILE=$(echo "$JSON_INPUT" | grep -o '"file_path":[[:space:]]*"[^"]*"' | sed 's/.*"\([^"]*\)"$/\1/')

# Exit silently if no file path found
if [ -z "$EDITED_FILE" ]; then
    exit 0
fi

# Only process code files that might benefit from logging
case "$EDITED_FILE" in
    *.py|*.go|*.js|*.jsx|*.ts|*.tsx)
        # Found a code file - process it
        ;;
    *)
        # Not a code file we handle - exit silently
        exit 0
        ;;
esac

# Project directory
PROJECT_DIR="/Users/walterday/Git/MerlinMines/Agent/Tools/tempest_search"

# Check if uv is available
if ! command -v uv &> /dev/null; then
    echo "⚠️  Warning: uv not found - logging agent skipped" >&2
    exit 0
fi

# Run the logging agent
# Environment variables (OLLAMA_URL, OLLAMA_API_KEY, DEFAULT_MODEL) are expected to be set in shell RC
# Change to project directory and run agent with --apply flag
cd "$PROJECT_DIR" && uv run agents/logging_debugging/logging_agent.py "$EDITED_FILE" --apply > /dev/null 2>&1 || true

# Always exit successfully so we don't block Claude Code
exit 0
