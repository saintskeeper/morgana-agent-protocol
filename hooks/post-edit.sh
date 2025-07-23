#!/bin/bash
# Claude hook: Main post-edit dispatcher
# Routes to specific formatters based on file type

# Read JSON input from stdin
JSON_INPUT=$(cat)

# Extract file path from JSON using grep and sed
EDITED_FILE=$(echo "$JSON_INPUT" | grep -o '"file_path":[[:space:]]*"[^"]*"' | sed 's/.*"\([^"]*\)"$/\1/')

# Exit silently if no file path found
if [ -z "$EDITED_FILE" ]; then
    exit 0
fi

# Only process files we have formatters for
case "$EDITED_FILE" in
    *.go|*.md|*.markdown|*.ts|*.tsx|*.js|*.jsx|*.json|*.yaml|*.yml|*.rs)
        # Continue with formatting
        ;;
    *)
        # No formatter for this file type, exit silently
        exit 0
        ;;
esac

# Get the directory of this script
HOOKS_DIR="$(cd "$(dirname "$0")" && pwd)"

# Run appropriate formatter based on file extension
case "$EDITED_FILE" in
    *.go)
        "$HOOKS_DIR/post-edit-go.sh" "$EDITED_FILE"
        # Play Orc Grunt sound for Go files
        nohup afplay ~/Sounds/game_samples/Orc/Grunt/Yes1.mp3 > /dev/null 2>&1 &
        ;;
    *.md|*.markdown)
        "$HOOKS_DIR/post-edit-markdown.sh" "$EDITED_FILE"
        # Play Night Elf Tyrande sound for documentation
        nohup afplay ~/Sounds/game_samples/NightElf/Tyrande/Yes1.mp3 > /dev/null 2>&1 &
        ;;
    *.ts|*.tsx|*.js|*.jsx|*.json)
        # TypeScript/JavaScript formatting
        if command -v prettier &> /dev/null; then
            echo "🔧 Formatting $EDITED_FILE with Prettier..."
            prettier --write "$EDITED_FILE"
            echo "✅ Prettier formatting complete"
            # Play Human Knight sound for web files
            nohup afplay ~/Sounds/game_samples/Human/Knight/Yes1.mp3 > /dev/null 2>&1 &
        fi
        ;;
    *.yaml|*.yml)
        # YAML formatting
        if command -v prettier &> /dev/null; then
            echo "🔧 Formatting YAML file: $EDITED_FILE"
            prettier --write "$EDITED_FILE"
            echo "✅ YAML formatting complete"
            # Play Medivh sound for config files
            nohup afplay ~/Sounds/game_samples/Human/Medivh/Yes2.mp3 > /dev/null 2>&1 &
        fi
        ;;
    *.rs)
        # Rust formatting
        if command -v rustfmt &> /dev/null; then
            echo "🔧 Formatting Rust file: $EDITED_FILE"
            rustfmt "$EDITED_FILE"
            echo "✅ Rust formatting complete"
            # Play Rexxar sound for Rust (tough/strong)
            nohup afplay ~/Sounds/game_samples/Orc/Rexxar/Yes1.mp3 > /dev/null 2>&1 &
        fi
        ;;
    *)
        # No specific formatter for this file type
        ;;
esac
