#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

DOCS_DIR="${1:-docs}"
URLS_FILE="${2:-urls.txt}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOG_FILE="logs/refresh_${TIMESTAMP}.log"
JSON_OUTPUT="logs/refresh_${TIMESTAMP}.json"

mkdir -p "$DOCS_DIR" logs

echo "📅 Daily Documentation Refresh - $(date)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create default urls.txt if it doesn't exist
if [ ! -f "$URLS_FILE" ]; then
    cat > "$URLS_FILE" << 'EOF'
# Documentation URLs to refresh daily
# Add one URL per line
https://docs.anthropic.com/en/api/openai-sdk
https://docs.anthropic.com/en/api/messages
https://docs.anthropic.com/en/api/client-sdks
https://docs.anthropic.com/en/api/streaming
EOF
    echo "✅ Created $URLS_FILE with default URLs"
fi

# Read URLs and create refresh prompt
URLS=$(grep -v '^#' "$URLS_FILE" | grep -v '^$')
URL_COUNT=$(echo "$URLS" | wc -l | tr -d ' ')

if [ -z "$URLS" ]; then
    echo "❌ No URLs found in $URLS_FILE"
    exit 1
fi

echo "📋 Found $URL_COUNT URLs to refresh"
echo ""

# Create the prompt for Claude
PROMPT="Please fetch and update the following documentation pages using WebFetch. For each URL:
1. Fetch the current content
2. Extract key information, API changes, and code examples
3. Save to $DOCS_DIR/ with filename based on the URL path
4. Include fetch timestamp and source URL in each file

URLs to refresh:
$URLS

Format each document as clean markdown with:
- Source URL at top
- Fetch timestamp
- Key sections preserved
- Code examples formatted properly"

# Execute with Claude CLI in headless mode
echo "🤖 Sending request to Claude (headless mode)..."
echo "$PROMPT" | tee "$LOG_FILE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >> "$LOG_FILE"

# Save prompt to temporary file to avoid stdin issues
TEMP_PROMPT_FILE="/tmp/claude_refresh_prompt_$$.txt"
echo "$PROMPT" > "$TEMP_PROMPT_FILE"

# Use Claude CLI with proper headless flags
# Using a file for input is more reliable than piping
if command -v claude >/dev/null 2>&1; then
    echo "🤖 Sending request to Claude..."
    claude --print "$(cat "$TEMP_PROMPT_FILE")" \
           --max-turns 10 \
           2>&1 | tee -a "$LOG_FILE"
    
    CLAUDE_EXIT_CODE=$?
    
    if [ $CLAUDE_EXIT_CODE -eq 0 ]; then
        echo "✅ Claude processed the request successfully"
    else
        echo "⚠️  Claude exited with code: $CLAUDE_EXIT_CODE"
    fi
else
    echo "❌ Claude CLI not found. Please install it first."
    echo "   Run: npm install -g @anthropic-ai/claude-cli"
    rm -f "$TEMP_PROMPT_FILE"
    exit 1
fi

# Clean up temp file
rm -f "$TEMP_PROMPT_FILE"

# Create a simple JSON summary of what was done
echo ""
echo "📊 Creating summary..."
cat > "$JSON_OUTPUT" << EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "urls_processed": $URL_COUNT,
  "docs_dir": "$DOCS_DIR",
  "log_file": "$LOG_FILE"
}
EOF

echo ""
echo "✅ Documents refreshed via Claude"
echo ""

# Show recent updates
if [ -d "$DOCS_DIR" ] && [ "$(ls -A $DOCS_DIR)" ]; then
    echo "📁 Recent documents in $DOCS_DIR:"
    ls -lt "$DOCS_DIR" | head -6
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Refresh complete"
echo "📝 Log saved to: $LOG_FILE"
echo "📊 JSON output: $JSON_OUTPUT"