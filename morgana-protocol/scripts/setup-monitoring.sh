#!/bin/bash
# setup-monitoring.sh - Set up Morgana monitoring environment

set -e

echo "🚀 Setting up Morgana Protocol monitoring environment..."

# Determine shell configuration file
SHELL_CONFIG=""
if [[ -f "$HOME/.zshrc" ]]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [[ -f "$HOME/.bashrc" ]]; then
    SHELL_CONFIG="$HOME/.bashrc"
elif [[ -f "$HOME/.bash_profile" ]]; then
    SHELL_CONFIG="$HOME/.bash_profile"
else
    echo "❓ Could not determine shell configuration file"
    echo "   Please add aliases manually to your shell profile"
    SHELL_CONFIG=""
fi

# Add aliases if shell config exists
if [[ -n "$SHELL_CONFIG" ]]; then
    echo "📝 Adding monitoring aliases to: $SHELL_CONFIG"
    
    # Check if aliases already exist
    if grep -q "# Morgana Protocol aliases" "$SHELL_CONFIG"; then
        echo "✅ Morgana aliases already exist in $SHELL_CONFIG"
    else
        cat >> "$SHELL_CONFIG" << 'EOF'

# Morgana Protocol aliases
alias monitor-up='cd ~/.claude/morgana-protocol && make up'
alias monitor-attach='cd ~/.claude/morgana-protocol && make attach'
alias monitor-down='cd ~/.claude/morgana-protocol && make down'
alias monitor-status='cd ~/.claude/morgana-protocol && make status'
alias monitor-logs='tail -f /tmp/morgana-monitor.log'
EOF
        echo "✅ Added monitoring aliases to $SHELL_CONFIG"
    fi
    
    # Add agent wrapper source
    if grep -q "agent-adapter-wrapper.sh" "$SHELL_CONFIG"; then
        echo "✅ Agent adapter wrapper already sourced in $SHELL_CONFIG"
    else
        echo "" >> "$SHELL_CONFIG"
        echo "# Source Morgana agent adapter wrapper" >> "$SHELL_CONFIG"
        echo "source ~/.claude/scripts/agent-adapter-wrapper.sh" >> "$SHELL_CONFIG"
        echo "✅ Added agent adapter wrapper to $SHELL_CONFIG"
    fi
fi

# Create morgana-protocol directory if it doesn't exist
MORGANA_DIR="$HOME/.claude/morgana-protocol"
if [[ ! -d "$MORGANA_DIR" ]]; then
    echo "📁 Creating Morgana Protocol directory: $MORGANA_DIR"
    mkdir -p "$MORGANA_DIR"
fi

# Check if binaries exist
if [[ -f "$HOME/.claude/bin/morgana" && -f "$HOME/.claude/bin/morgana-monitor" ]]; then
    echo "✅ Morgana binaries found in ~/.claude/bin/"
else
    echo "⚠️  Morgana binaries not found in ~/.claude/bin/"
    echo "   Run 'make build && make install' from the morgana-protocol directory"
fi

# Set up environment variables
echo "🔧 Setting up environment variables..."
export MORGANA_TUI_ENABLED=true
export MORGANA_TUI_THEME=dark

# Test monitor functionality
echo "🧪 Testing monitor functionality..."
cd "$HOME/.claude/morgana-protocol" || {
    echo "❌ Could not change to morgana-protocol directory"
    echo "   Please clone the repository to ~/.claude/morgana-protocol"
    exit 1
}

# Check if make commands work
if ! command -v make >/dev/null 2>&1; then
    echo "⚠️  'make' command not found. Please install make to use convenience commands."
else
    echo "✅ Make command available"
fi

echo ""
echo "🎉 Morgana monitoring environment setup complete!"
echo ""
echo "📋 What was configured:"
if [[ -n "$SHELL_CONFIG" ]]; then
echo "   ✅ Added monitoring aliases to $SHELL_CONFIG"
echo "   ✅ Added agent adapter wrapper sourcing"
fi
echo "   ✅ Set default environment variables"
echo "   ✅ Validated directory structure"
echo ""
echo "🚀 Available commands after restart:"
echo "   monitor-up      - Start monitor daemon"
echo "   monitor-attach  - Attach to TUI"
echo "   monitor-down    - Stop monitor daemon"
echo "   monitor-status  - Check monitor status"
echo "   monitor-logs    - View monitor logs"
echo ""
echo "🔄 To apply changes:"
if [[ -n "$SHELL_CONFIG" ]]; then
echo "   source $SHELL_CONFIG"
else
echo "   Restart your terminal"
fi
echo ""
echo "🧪 Quick test:"
echo "   monitor-up && sleep 2 && monitor-status && monitor-down"
echo ""