#!/bin/bash
# Script to set up local development environment for Claude Code configs

echo "🔧 Setting up local development environment..."

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

# Set up Morgana Protocol
echo ""
echo "🔧 Setting up Morgana Protocol..."

# Create necessary directories
mkdir -p ~/.claude/scripts
mkdir -p ~/.claude/bin

# Check if we're in the right directory
if [ -d "morgana-protocol" ]; then
    echo "  Found morgana-protocol directory"
    
    # Build Morgana if not already built
    if [ ! -f "morgana-protocol/dist/morgana" ]; then
        echo "  Building Morgana..."
        cd morgana-protocol
        make build
        cd ..
    fi
    
    # Install Morgana binary
    if [ -f "morgana-protocol/dist/morgana" ]; then
        echo "  Installing Morgana binary..."
        cp morgana-protocol/dist/morgana ~/.claude/bin/
        chmod +x ~/.claude/bin/morgana
    fi
    
    # Copy scripts
    echo "  Installing helper scripts..."
    [ -f "scripts/agent-adapter-wrapper.sh" ] && cp scripts/agent-adapter-wrapper.sh ~/.claude/scripts/
    [ -f "scripts/agent_adapter.py" ] && cp scripts/agent_adapter.py ~/.claude/scripts/
    [ -f "morgana-protocol/scripts/task_bridge.py" ] && cp morgana-protocol/scripts/task_bridge.py ~/.claude/scripts/
    
    # Make scripts executable
    chmod +x ~/.claude/scripts/*.sh ~/.claude/scripts/*.py 2>/dev/null
    
    # Create default Morgana config if it doesn't exist
    if [ ! -f "~/.claude/morgana.yaml" ]; then
        echo "  Creating default Morgana configuration..."
        cat > ~/.claude/morgana.yaml << 'MORGANA_CONFIG'
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m
  timeouts:
    code-implementer: 5m
    test-specialist: 3m
    validation-expert: 2m
    sprint-planner: 2m

execution:
  max_concurrency: 5
  default_mode: sequential
  queue_size: 100

telemetry:
  enabled: false
  exporter: none
  service_name: morgana-protocol
  environment: local

task_client:
  bridge_path: ~/.claude/scripts/task_bridge.py
  python_path: python3
  mock_mode: false
  timeout: 30s
MORGANA_CONFIG
    fi
    
    echo "✅ Morgana Protocol setup complete!"
else
    echo "⚠️  morgana-protocol directory not found - skipping Morgana setup"
fi

# Add PATH update suggestion
echo ""
echo "📝 Add the following to your shell configuration (.bashrc/.zshrc):"
echo ""
echo "# Claude Code configurations"
echo "export PATH=\"\$HOME/.claude/bin:\$PATH\""
echo "source ~/.claude/scripts/agent-adapter-wrapper.sh"

echo ""
echo "Setup complete! Your environment is now configured with:"
echo "  • Git hooks for qsweep and auto-formatting"
echo "  • Go and Markdown files are auto-formatted before commit"
echo "  • Morgana Protocol for parallel agent execution"
echo ""
echo "To test Morgana: morgana --version"
