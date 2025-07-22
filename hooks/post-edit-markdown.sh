#!/bin/bash
# Claude hook: Runs after editing Markdown files
# Fixes end-of-file newlines and formatting issues

# Get the edited file path from Claude
EDITED_FILE="$1"

# Check if it's a Markdown file
if [[ "$EDITED_FILE" =~ \.(md|markdown)$ ]]; then
    echo "🔧 Formatting Markdown file: $EDITED_FILE"

    # Fix end-of-file newline
    if [ -s "$EDITED_FILE" ] && [ -z "$(tail -c 1 "$EDITED_FILE")" ]; then
        echo "✓ File already has proper EOF newline"
    else
        echo "" >> "$EDITED_FILE"
        echo "✓ Added EOF newline"
    fi

    # Fix trailing whitespace
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS version
        sed -i '' 's/[[:space:]]*$//' "$EDITED_FILE"
    else
        # Linux version
        sed -i 's/[[:space:]]*$//' "$EDITED_FILE"
    fi
    echo "✓ Removed trailing whitespace"

    # Optional: Run prettier if available
    if command -v prettier &> /dev/null; then
        prettier --write "$EDITED_FILE" --prose-wrap always
        echo "✓ Prettier formatting applied"
    fi

    echo "✅ Markdown formatting complete"
fi
