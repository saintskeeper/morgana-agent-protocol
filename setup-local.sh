#!/bin/bash
# Script to set up qsweep as a git hook locally

echo "🔧 Setting up qsweep as a git post-checkout hook..."

# Create git hooks directory if it doesn't exist
mkdir -p .git/hooks

# Create post-checkout hook
cat > .git/hooks/post-checkout << 'EOF'
#!/bin/bash
# Git hook: Runs after checkout (including branch creation)
# $1: previous HEAD
# $2: new HEAD
# $3: 1 if checking out a branch, 0 if checking out files

# Only run on branch checkout, not file checkout
if [ "$3" = "1" ]; then
    # Check if this is a new branch (previous and current HEAD are different)
    if [ "$1" != "$2" ]; then
        echo "🧹 New branch detected - running qsweep..."

        # Run qsweep if it exists
        if [ -f "./qsweep.sh" ]; then
            ./qsweep.sh
            echo "✅ Documentation sweep complete!"
        fi
    fi
fi
EOF

# Make hook executable
chmod +x .git/hooks/post-checkout

echo "✅ Post-checkout hook installed!"

# Install pre-commit hook for auto-formatting
echo ""
echo "🔧 Setting up pre-commit auto-formatting hook..."

cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Git pre-commit hook: Auto-format modified files before commit

echo "🔧 Running pre-commit formatting..."

# Get list of modified files
FILES=$(git diff --cached --name-only --diff-filter=ACM)

# Format each file based on type
for file in $FILES; do
    if [ -f "$file" ]; then
        case "$file" in
            *.go)
                echo "  Formatting Go: $file"
                gofmt -w "$file"
                # Re-add file if it was modified
                git add "$file"
                ;;
            *.md|*.markdown)
                echo "  Formatting Markdown: $file"
                # Fix EOF newline
                if [ -s "$file" ] && [ -z "$(tail -c 1 "$file")" ]; then
                    :
                else
                    echo "" >> "$file"
                fi
                # Remove trailing whitespace
                sed -i '' 's/[[:space:]]*$//' "$file"
                git add "$file"
                ;;
        esac
    fi
done

# Continue with pre-commit checks
if [ -f ".pre-commit-config.yaml" ]; then
    pre-commit run --files $FILES
fi

echo "✅ Pre-commit formatting complete"
EOF

chmod +x .git/hooks/pre-commit

echo "✅ Pre-commit formatting hook installed!"
echo ""
echo "Setup complete! Your hooks are now active:"
echo "  • qsweep runs on branch creation/checkout"
echo "  • Go and Markdown files are auto-formatted before commit"
echo ""
echo "To test: git checkout -b test-branch"
