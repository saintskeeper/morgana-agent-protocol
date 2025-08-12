#!/bin/bash
# validate-migration.sh - Validate Morgana Protocol migration to event stream architecture

set -e

MORGANA_DIR="$HOME/.claude/morgana-protocol"
ERRORS=0
WARNINGS=0

echo "🔍 Validating Morgana Protocol migration..."
echo "📍 Morgana directory: $MORGANA_DIR"
echo ""

# Color codes for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

error() {
    echo -e "${RED}❌ ERROR: $1${NC}"
    ERRORS=$((ERRORS + 1))
}

warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    WARNINGS=$((WARNINGS + 1))
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if we're in the right directory
if [[ ! -d "$MORGANA_DIR" ]]; then
    error "Morgana directory not found: $MORGANA_DIR"
    echo "   Please ensure Morgana Protocol is cloned to ~/.claude/morgana-protocol"
    exit 1
fi

cd "$MORGANA_DIR"

echo "📋 1. Checking binary installation..."
if [[ -f "$HOME/.claude/bin/morgana" ]]; then
    success "morgana binary found"
    MORGANA_VERSION=$("$HOME/.claude/bin/morgana" --version 2>/dev/null | head -1 || echo "unknown")
    info "Version: $MORGANA_VERSION"
else
    error "morgana binary not found in ~/.claude/bin/"
    echo "      Run: make build && make install"
fi

if [[ -f "$HOME/.claude/bin/morgana-monitor" ]]; then
    success "morgana-monitor binary found"
else
    error "morgana-monitor binary not found in ~/.claude/bin/"
    echo "      Run: make build && make install"
fi

echo ""
echo "📋 2. Checking configuration..."
CONFIG_FILE="$MORGANA_DIR/morgana.yaml"

if [[ -f "$CONFIG_FILE" ]]; then
    success "Configuration file exists: $CONFIG_FILE"
    
    # Check for deprecated sections
    if grep -q "^task_client:" "$CONFIG_FILE"; then
        warning "Deprecated task_client configuration found"
        echo "      Run: ./scripts/migrate-config.sh"
    else
        success "No deprecated task_client configuration"
    fi
    
    # Check for required TUI configuration
    if grep -q "^tui:" "$CONFIG_FILE"; then
        success "TUI configuration found"
        
        # Check specific TUI settings
        if grep -A 20 "^tui:" "$CONFIG_FILE" | grep -q "events:"; then
            success "TUI events configuration found"
        else
            warning "TUI events configuration missing"
            echo "      Add events section to TUI configuration"
        fi
    else
        warning "TUI configuration missing"
        echo "      Run: ./scripts/migrate-config.sh"
    fi
else
    error "Configuration file not found: $CONFIG_FILE"
fi

echo ""
echo "📋 3. Testing core functionality..."

# Test morgana command
if command -v morgana >/dev/null 2>&1; then
    success "morgana command available"
    
    # Quick dry run test
    if morgana --help >/dev/null 2>&1; then
        success "morgana command working"
    else
        error "morgana command fails basic test"
    fi
else
    error "morgana command not in PATH"
    echo "      Ensure ~/.claude/bin is in your PATH"
fi

# Test monitor binary
if command -v morgana-monitor >/dev/null 2>&1; then
    success "morgana-monitor command available"
else
    error "morgana-monitor command not in PATH"
    echo "      Ensure ~/.claude/bin is in your PATH"
fi

echo ""
echo "📋 4. Testing monitor system..."

# Check if monitor is already running
if pgrep -f morgana-monitor >/dev/null; then
    info "Monitor already running - stopping for test"
    pkill -f morgana-monitor 2>/dev/null || true
    sleep 2
fi

# Test monitor startup
if [[ -f "$HOME/.claude/bin/morgana-monitor" ]]; then
    info "Starting monitor for test..."
    "$HOME/.claude/bin/morgana-monitor" --headless --test >/dev/null 2>&1 &
    MONITOR_PID=$!
    sleep 3
    
    if kill -0 $MONITOR_PID 2>/dev/null; then
        success "Monitor starts successfully"
        
        # Test socket creation
        if [[ -e "/tmp/morgana-monitor.sock" ]]; then
            success "Monitor socket created"
        else
            warning "Monitor socket not found"
        fi
        
        # Clean up test monitor
        kill $MONITOR_PID 2>/dev/null || true
        rm -f /tmp/morgana-monitor.sock
    else
        error "Monitor fails to start"
    fi
else
    error "Cannot test monitor - binary not found"
fi

echo ""
echo "📋 5. Testing integration..."

# Test event generation
info "Testing basic agent execution..."
if timeout 10s morgana -- --agent validation-expert --prompt "Test migration validation" --mock >/dev/null 2>&1; then
    success "Basic agent execution works"
else
    warning "Agent execution test failed - may need Claude Code setup"
    echo "      This is normal if Claude Code is not configured"
fi

echo ""
echo "📋 6. Checking shell integration..."

# Check for shell aliases
SHELL_FILES=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile")
FOUND_ALIASES=false

for shell_file in "${SHELL_FILES[@]}"; do
    if [[ -f "$shell_file" ]] && grep -q "monitor-up" "$shell_file"; then
        success "Monitor aliases found in $(basename "$shell_file")"
        FOUND_ALIASES=true
        break
    fi
done

if ! $FOUND_ALIASES; then
    warning "Monitor aliases not found in shell configuration"
    echo "      Run: ./scripts/setup-monitoring.sh"
fi

# Check for agent adapter wrapper
FOUND_WRAPPER=false
for shell_file in "${SHELL_FILES[@]}"; do
    if [[ -f "$shell_file" ]] && grep -q "agent-adapter-wrapper.sh" "$shell_file"; then
        success "Agent adapter wrapper sourced in $(basename "$shell_file")"
        FOUND_WRAPPER=true
        break
    fi
done

if ! $FOUND_WRAPPER; then
    warning "Agent adapter wrapper not sourced"
    echo "      Run: ./scripts/setup-monitoring.sh"
fi

echo ""
echo "📋 7. Checking for old IPC artifacts..."

# Check for archived files
if [[ -d "$MORGANA_DIR/archive/ipc-removal" ]]; then
    success "IPC removal archive found (good for rollback if needed)"
else
    info "No IPC removal archive found"
fi

# Check for old Python bridge files (should be archived)
OLD_BRIDGE_FILES=(
    "scripts/task_bridge.py"
    "scripts/task_bridge_claude.py"
    "scripts/morgana_repl.py"
    "pkg/task/client.go"
)

OLD_FILES_EXIST=false
for file in "${OLD_BRIDGE_FILES[@]}"; do
    if [[ -f "$MORGANA_DIR/$file" ]]; then
        warning "Old IPC file still exists: $file"
        OLD_FILES_EXIST=true
    fi
done

if ! $OLD_FILES_EXIST; then
    success "No old IPC files found in main codebase"
fi

echo ""
echo "📊 Validation Summary"
echo "===================="

if [[ $ERRORS -eq 0 && $WARNINGS -eq 0 ]]; then
    echo -e "${GREEN}🎉 Migration validation PASSED!${NC}"
    echo ""
    echo "✅ Your Morgana Protocol is fully migrated to event stream architecture"
    echo "✅ All components are working correctly"
    echo "✅ Ready for production use"
    echo ""
    echo "🚀 Quick start:"
    echo "   monitor-up        # Start monitor daemon"
    echo "   monitor-attach    # Attach TUI interface"
    echo "   morgana -- --agent code-implementer --prompt 'Hello world'"
elif [[ $ERRORS -eq 0 ]]; then
    echo -e "${YELLOW}⚠️  Migration validation completed with $WARNINGS warning(s)${NC}"
    echo ""
    echo "✅ Migration is functional but has minor issues"
    echo "📝 Review warnings above and run suggested commands"
    echo "🚀 System should work correctly for basic usage"
elif [[ $ERRORS -gt 0 ]]; then
    echo -e "${RED}❌ Migration validation FAILED with $ERRORS error(s) and $WARNINGS warning(s)${NC}"
    echo ""
    echo "🔧 Please fix errors above before proceeding"
    echo "📖 See docs/MIGRATION_GUIDE.md for detailed instructions"
    exit 1
fi

echo ""
echo "📖 For more information:"
echo "   • Migration Guide: docs/MIGRATION_GUIDE.md"
echo "   • Getting Started: docs/GETTING_STARTED.md"
echo "   • Architecture: docs/ARCHITECTURE.md"
echo ""