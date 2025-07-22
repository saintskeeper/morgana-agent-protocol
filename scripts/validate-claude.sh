#!/bin/bash
# Validate CLAUDE.md integrity

# Read JSON input from stdin
JSON_INPUT=$(cat)

# Extract file path from JSON using grep and sed (portable approach)
EDITED_FILE=$(echo "$JSON_INPUT" | grep -o '"file_path":[[:space:]]*"[^"]*"' | sed 's/.*"\([^"]*\)"$/\1/')

# Only run validation if CLAUDE.md was edited
if [[ ! "$EDITED_FILE" =~ CLAUDE\.md$ ]]; then
    # Not editing CLAUDE.md, exit silently with success
    exit 0
fi

CLAUDE_FILE="$EDITED_FILE"
MIN_LINES=200  # Minimum expected lines
REQUIRED_SECTIONS=(
    "Linear Integration"
    "Architecture Overview"
    "Development Commands"
    "Claude Code Guidelines"
    "important-instruction-reminders"
)

echo "🔍 Validating CLAUDE.md..."

# Check if file exists
if [[ ! -f "$CLAUDE_FILE" ]]; then
    echo "❌ CLAUDE.md not found!"
    exit 1
fi

# Check line count
LINE_COUNT=$(wc -l < "$CLAUDE_FILE")
if [[ $LINE_COUNT -lt $MIN_LINES ]]; then
    echo "⚠️  Warning: CLAUDE.md has only $LINE_COUNT lines (expected >$MIN_LINES)"
    echo "   File may have been truncated!"
fi

# Check required sections
MISSING_SECTIONS=()
for section in "${REQUIRED_SECTIONS[@]}"; do
    if ! grep -q "$section" "$CLAUDE_FILE"; then
        MISSING_SECTIONS+=("$section")
    fi
done

if [[ ${#MISSING_SECTIONS[@]} -gt 0 ]]; then
    echo "❌ Missing sections:"
    for section in "${MISSING_SECTIONS[@]}"; do
        echo "   - $section"
    done
    exit 1
fi

echo "✅ CLAUDE.md validation passed!"
echo "   Lines: $LINE_COUNT"
echo "   All required sections present"