# Morgana System Makefile
# Convenient commands for managing Morgana Protocol and monitoring

.PHONY: help up down status install build clean logs attach

# Default target - show help
help:
	@echo "Morgana System Management"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  up        Start morgana-monitor daemon (TUI monitoring)"
	@echo "  down      Stop morgana-monitor daemon"
	@echo "  status    Show morgana-monitor status"
	@echo "  attach    Attach to running TUI session (screen/tmux)"
	@echo "  logs      Show morgana-monitor logs"
	@echo "  install   Build and install morgana binaries to /usr/local/bin"
	@echo "  build     Build morgana binaries (development)"
	@echo "  clean     Clean build artifacts and stop all services"
	@echo "  kill      Force kill all morgana processes"
	@echo "  test      Run a test agent to verify system"
	@echo ""
	@echo "Examples:"
	@echo "  make up       # Start monitoring"
	@echo "  make install  # Install binaries system-wide"
	@echo "  make status   # Check if monitor is running"

# Start morgana-monitor daemon
up:
	@echo "🚀 Starting Morgana Monitor..."
	@~/.claude/scripts/morgana-monitor-ctl.sh start

# Stop morgana-monitor daemon
down:
	@echo "🛑 Stopping Morgana Monitor..."
	@~/.claude/scripts/morgana-monitor-ctl.sh stop

# Check morgana-monitor status
status:
	@~/.claude/scripts/morgana-monitor-ctl.sh status

# Attach to running TUI session
attach:
	@echo "📺 Attaching to Morgana Monitor TUI..."
	@~/.claude/scripts/morgana-monitor-ctl.sh attach

# Show morgana-monitor logs
logs:
	@echo "📝 Morgana Monitor Logs:"
	@~/.claude/scripts/morgana-monitor-ctl.sh logs

# Build morgana binaries (both morgana and morgana-monitor)
build:
	@echo "🔨 Building Morgana binaries..."
	@cd morgana-protocol && make dev
	@echo "✅ Build complete:"
	@echo "  - morgana-protocol/dist/morgana"
	@echo "  - morgana-protocol/dist/morgana-monitor"

# Install morgana binaries to /usr/local/bin
install: build
	@echo "📦 Installing Morgana binaries to /usr/local/bin..."
	@echo "⚠️  This may require sudo privileges"
	@# Try to copy without sudo first (in case user has write permissions)
	@cp morgana-protocol/dist/morgana /usr/local/bin/morgana 2>/dev/null || \
		(echo "📝 Need sudo to install morgana..." && sudo cp morgana-protocol/dist/morgana /usr/local/bin/morgana)
	@chmod +x /usr/local/bin/morgana 2>/dev/null || sudo chmod +x /usr/local/bin/morgana
	@cp morgana-protocol/dist/morgana-monitor /usr/local/bin/morgana-monitor 2>/dev/null || \
		(echo "📝 Need sudo to install morgana-monitor..." && sudo cp morgana-protocol/dist/morgana-monitor /usr/local/bin/morgana-monitor)
	@chmod +x /usr/local/bin/morgana-monitor 2>/dev/null || sudo chmod +x /usr/local/bin/morgana-monitor
	@echo "✅ Installed to /usr/local/bin:"
	@echo "  - /usr/local/bin/morgana"
	@echo "  - /usr/local/bin/morgana-monitor"
	@echo ""
	@echo "🎯 Verifying installation..."
	@which morgana && morgana --version || echo "⚠️  morgana not found in PATH"
	@which morgana-monitor && echo "✅ morgana-monitor installed" || echo "⚠️  morgana-monitor not found in PATH"

# Install to user directory (no sudo required)
install-user: build
	@echo "📦 Installing Morgana binaries to ~/.claude/bin..."
	@mkdir -p ~/.claude/bin
	@cp morgana-protocol/dist/morgana ~/.claude/bin/morgana
	@chmod +x ~/.claude/bin/morgana
	@cp morgana-protocol/dist/morgana-monitor ~/.claude/bin/morgana-monitor
	@chmod +x ~/.claude/bin/morgana-monitor
	@echo "✅ Installed to ~/.claude/bin:"
	@echo "  - ~/.claude/bin/morgana"
	@echo "  - ~/.claude/bin/morgana-monitor"
	@echo ""
	@echo "💡 Add to PATH if needed: export PATH=$$HOME/.claude/bin:$$PATH"

# Clean build artifacts and stop services
clean: down
	@echo "🧹 Cleaning up..."
	@# Kill all morgana-monitor processes
	@pkill -f morgana-monitor 2>/dev/null || true
	@# Clean up screen sessions
	@screen -S morgana-monitor -X quit 2>/dev/null || true
	@# Clean up tmux sessions
	@tmux kill-session -t morgana-monitor 2>/dev/null || true
	@# Remove artifacts
	@cd morgana-protocol && make clean
	@rm -f /tmp/morgana.sock /tmp/morgana-monitor.pid /tmp/morgana-monitor.log
	@echo "✅ Cleanup complete - all morgana-monitor instances killed"

# Run a test agent to verify system
test:
	@echo "🧪 Running test agent..."
	@source ~/.claude/scripts/agent-adapter-wrapper.sh && \
		AgentAdapter "sprint-planner" "Create a test plan for verifying the monitoring system"

# Quick test with parallel agents
test-parallel:
	@echo "🧪 Running parallel test agents..."
	@echo '[{"agent_type":"code-implementer","prompt":"Test implementation"},{"agent_type":"test-specialist","prompt":"Test validation"}]' | morgana --parallel

# Full system check
check: status
	@echo ""
	@echo "🔍 System Check:"
	@echo -n "  morgana binary: "
	@which morgana > /dev/null 2>&1 && echo "✅ Found at $$(which morgana)" || echo "❌ Not found"
	@echo -n "  morgana-monitor binary: "
	@which morgana-monitor > /dev/null 2>&1 && echo "✅ Found at $$(which morgana-monitor)" || echo "❌ Not found"
	@echo -n "  IPC socket: "
	@test -S /tmp/morgana.sock && echo "✅ Active at /tmp/morgana.sock" || echo "❌ Not found"
	@echo -n "  Agent adapter: "
	@test -f ~/.claude/scripts/morgana-adapter.sh && echo "✅ Installed" || echo "❌ Not found"
	@echo ""

# Development shortcuts
dev-up: build install-user up
	@echo "✅ Development environment ready"

dev-down: down clean
	@echo "✅ Development environment stopped"

# Restart the monitor
restart: down up
	@echo "✅ Monitor restarted"

# Force kill all morgana processes
kill:
	@echo "☠️  Force killing all morgana processes..."
	@pkill -9 -f morgana-monitor 2>/dev/null || true
	@pkill -9 -f morgana 2>/dev/null || true
	@screen -S morgana-monitor -X quit 2>/dev/null || true
	@tmux kill-session -t morgana-monitor 2>/dev/null || true
	@rm -f /tmp/morgana.sock /tmp/morgana-monitor.pid
	@echo "✅ All morgana processes terminated"

# Show running morgana processes
ps:
	@echo "📊 Morgana Processes:"
	@ps aux | grep -E "morgana|morgana-monitor" | grep -v grep || echo "No morgana processes running"

# Tail monitor log in real-time
tail:
	@echo "📜 Tailing morgana-monitor log (Ctrl+C to stop)..."
	@tail -f /tmp/morgana-monitor.log

# Uninstall from /usr/local/bin
uninstall:
	@echo "🗑️  Uninstalling Morgana from /usr/local/bin..."
	@echo "⚠️  This requires sudo privileges"
	@sudo rm -f /usr/local/bin/morgana /usr/local/bin/morgana-monitor
	@echo "✅ Uninstalled from /usr/local/bin"

# Uninstall from user directory
uninstall-user:
	@echo "🗑️  Uninstalling Morgana from ~/.claude/bin..."
	@rm -f ~/.claude/bin/morgana ~/.claude/bin/morgana-monitor
	@echo "✅ Uninstalled from ~/.claude/bin"