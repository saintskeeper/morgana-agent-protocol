#!/bin/bash
# Claude hook: Runs after editing Go files
# Automatically formats Go code with gofmt

# Get the edited file path from Claude
EDITED_FILE="$1"

# Check if it's a Go file
if [[ "$EDITED_FILE" =~ \.go$ ]]; then
    echo "🔧 Formatting Go file: $EDITED_FILE"

    # Run gofmt
    gofmt -w "$EDITED_FILE"

    # Also run goimports if available (organizes imports)
    if command -v goimports &> /dev/null; then
        goimports -w "$EDITED_FILE"
        echo "✅ Go formatting complete (gofmt + goimports)"
    else
        echo "✅ Go formatting complete (gofmt)"
        echo "💡 Tip: Install goimports for automatic import organization"
    fi
fi
