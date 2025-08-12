#!/bin/bash
# migrate-config.sh - Clean up old Morgana configuration for event stream migration

set -e

CONFIG_FILE="$HOME/.claude/morgana-protocol/morgana.yaml"
BACKUP_FILE="${CONFIG_FILE}.backup-$(date +%Y%m%d-%H%M%S)"

echo "🔄 Migrating Morgana Protocol configuration..."
echo "📍 Config file: $CONFIG_FILE"

# Check if config file exists
if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "❌ Configuration file not found: $CONFIG_FILE"
    echo "   Run 'make dev' first to create default configuration"
    exit 1
fi

# Backup original configuration
cp "$CONFIG_FILE" "$BACKUP_FILE"
echo "📋 Backed up original config to: $BACKUP_FILE"

# Check if task_client section exists
if grep -q "^task_client:" "$CONFIG_FILE"; then
    echo "🧹 Removing deprecated task_client configuration..."
    
    # Create temporary file without task_client section
    # Remove from task_client: to the next top-level section
    awk '
        /^task_client:/ { in_task_client = 1; next }
        /^[a-zA-Z_][a-zA-Z0-9_]*:/ && in_task_client { in_task_client = 0 }
        !in_task_client { print }
    ' "$CONFIG_FILE" > "${CONFIG_FILE}.tmp"
    
    mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
    echo "✅ Removed deprecated task_client configuration"
else
    echo "✅ No deprecated task_client configuration found"
fi

# Verify TUI configuration exists
if grep -q "^tui:" "$CONFIG_FILE"; then
    echo "✅ TUI configuration already present"
else
    echo "➕ Adding default TUI configuration..."
    cat >> "$CONFIG_FILE" << 'EOF'

# TUI (Terminal User Interface) configuration
tui:
  enabled: true
  performance:
    refresh_rate: 50ms
    max_log_lines: 1000
    target_fps: 20
  visual:
    theme:
      name: dark
      primary: "#7C3AED"
  events:
    buffer_size: 1000
    enable_batching: true
    batch_size: 50
    batch_timeout: 100ms
EOF
    echo "✅ Added default TUI configuration"
fi

echo ""
echo "🎉 Migration complete!"
echo ""
echo "📊 Changes made:"
echo "   - Removed deprecated task_client configuration"
echo "   - Ensured TUI configuration is present"
echo "   - Preserved all other settings"
echo ""
echo "📝 To review changes:"
echo "   diff $BACKUP_FILE $CONFIG_FILE"
echo ""
echo "🚀 Next steps:"
echo "   1. Run 'make build && make install' to update binaries"
echo "   2. Run 'make up' to start the new monitor"
echo "   3. Run 'make attach' to view the TUI"
echo ""