# Makefile for Morgana Protocol
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.Version=$(VERSION)
BINARY := morgana
MONITOR_BINARY := morgana-monitor

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOGET := $(GOCMD) get

# Directories
CMD_DIR := ./cmd/morgana
MONITOR_CMD_DIR := ./cmd/morgana-monitor
DIST_DIR := ./dist

.PHONY: all build test clean deps install dev run help

# Default target
all: test build

# Help target
help:
	@echo "Morgana Protocol - Agent Orchestration System"
	@echo ""
	@echo "Build targets:"
	@echo "  make build       - Build for all platforms"
	@echo "  make dev         - Build for current platform only"
	@echo "  make clean       - Clean build artifacts"
	@echo ""
	@echo "Install targets:"
	@echo "  make install     - Install to ~/.claude/bin"
	@echo "  sudo make install-system - Install to /usr/local/bin"
	@echo ""
	@echo "Test targets:"
	@echo "  make test        - Run unit tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-integration-coverage - Run integration tests with coverage"
	@echo "  make test-all    - Run all tests"
	@echo ""
	@echo "Integration test suites:"
	@echo "  make test-events - Event stream tests"
	@echo "  make test-tui    - TUI integration tests"  
	@echo "  make test-monitor - Monitor system tests"
	@echo "  make test-commands - Command polling tests"
	@echo "  make test-resources - Resource leak tests"
	@echo "  make test-suite  - Complete integration suite"
	@echo "  make test-performance - Performance validation tests"
	@echo ""
	@echo "Monitor targets:"
	@echo "  make up          - Start monitor in headless mode"
	@echo "  make down        - Stop monitor daemon"
	@echo "  make attach      - Attach TUI to running monitor"
	@echo "  make status      - Check monitor status"
	@echo "  make monitor     - Build and start monitor"
	@echo ""
	@echo "Development targets:"
	@echo "  make deps        - Download dependencies"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Lint code"
	@echo "  make mocks       - Generate mocks for testing"
	@echo "  make run         - Run with example"
	@echo ""

# Development build (current platform only)
dev:
	@echo "Building for development..."
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY) $(CMD_DIR)
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(MONITOR_BINARY) $(MONITOR_CMD_DIR)

# Build for all platforms
build: build-darwin build-linux build-windows

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-darwin-arm64 $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(MONITOR_BINARY)-darwin-amd64 $(MONITOR_CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(MONITOR_BINARY)-darwin-arm64 $(MONITOR_CMD_DIR)

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-linux-arm64 $(CMD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(MONITOR_BINARY)-linux-amd64 $(MONITOR_CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(MONITOR_BINARY)-linux-arm64 $(MONITOR_CMD_DIR)

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-windows-amd64.exe $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(MONITOR_BINARY)-windows-amd64.exe $(MONITOR_CMD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -cover -short ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out -short ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration -timeout=30m ./tests/integration

# Run integration tests with coverage
test-integration-coverage:
	@echo "Running integration tests with coverage..."
	$(GOTEST) -v -tags=integration -timeout=30m -coverprofile=integration-coverage.out ./tests/integration
	$(GOCMD) tool cover -html=integration-coverage.out -o integration-coverage.html
	@echo "Integration coverage report generated: integration-coverage.html"

# Run specific integration test suites
test-events:
	@echo "Running event stream integration tests..."
	$(GOTEST) -v -tags=integration -run TestEventStream ./tests/integration

test-tui:
	@echo "Running TUI integration tests..."
	$(GOTEST) -v -tags=integration -run TestTUI ./tests/integration

test-monitor:
	@echo "Running monitor integration tests..."
	$(GOTEST) -v -tags=integration -run TestMonitor ./tests/integration

test-commands:
	@echo "Running command polling tests..."
	$(GOTEST) -v -tags=integration -run TestCommand ./tests/integration

test-resources:
	@echo "Running resource leak tests..."
	$(GOTEST) -v -tags=integration -run TestResource ./tests/integration

test-suite:
	@echo "Running comprehensive integration suite..."
	$(GOTEST) -v -tags=integration -run TestMorganaIntegrationSuite ./tests/integration

# Run performance validation tests only
test-performance:
	@echo "Running performance validation tests..."
	$(GOTEST) -v -tags=integration -run Performance ./tests/integration

# Run all tests including integration
test-all: test test-integration

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install locally
install: dev
	@echo "Installing to ~/.claude/bin..."
	@mkdir -p ~/.claude/bin
	@cp $(DIST_DIR)/$(BINARY) ~/.claude/bin/
	@cp $(DIST_DIR)/$(MONITOR_BINARY) ~/.claude/bin/
	@chmod +x ~/.claude/bin/$(BINARY)
	@chmod +x ~/.claude/bin/$(MONITOR_BINARY)
	@echo "Installed to ~/.claude/bin/$(BINARY) and ~/.claude/bin/$(MONITOR_BINARY)"

# Install to /usr/local/bin (requires sudo)
install-system: dev
	@echo "Installing to /usr/local/bin (requires sudo)..."
	@mkdir -p /usr/local/bin
	@cp $(DIST_DIR)/$(BINARY) /usr/local/bin/
	@cp $(DIST_DIR)/$(MONITOR_BINARY) /usr/local/bin/
	@chmod +x /usr/local/bin/$(BINARY)
	@chmod +x /usr/local/bin/$(MONITOR_BINARY)
	@echo "Installed to /usr/local/bin/$(BINARY) and /usr/local/bin/$(MONITOR_BINARY)"

# Run example
run: dev
	@echo "Running example..."
	$(DIST_DIR)/$(BINARY) --agent code-implementer --prompt "Hello, Morgana!"

# Initialize project (first time setup)
init:
	@echo "Initializing Morgana Protocol..."
	$(GOMOD) init github.com/saintskeeper/claude-code-configs/morgana-protocol
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Generate mocks for testing
mocks:
	@echo "Generating mocks..."
	mockgen -source=internal/adapter/adapter.go -destination=internal/adapter/mocks/adapter_mock.go
	mockgen -source=pkg/task/client.go -destination=pkg/task/mocks/client_mock.go

# Monitor control targets
up:
	@echo "🚀 Starting Morgana Monitor in headless mode..."
	@$(DIST_DIR)/$(MONITOR_BINARY) --headless > /tmp/morgana-monitor.log 2>&1 &
	@echo "✅ Monitor started (PID: $$!)"
	@echo "📝 Logs: /tmp/morgana-monitor.log"

down:
	@echo "🛑 Stopping Morgana Monitor..."
	@pkill -f $(MONITOR_BINARY) 2>/dev/null || echo "Monitor not running"
	@rm -f /tmp/morgana.sock
	@echo "✅ Monitor stopped"

attach:
	@echo "📺 Attaching to Morgana Monitor TUI..."
	@echo "   Press Q, Esc, or Ctrl+C to disconnect gracefully"
	@$(DIST_DIR)/$(MONITOR_BINARY) --client || echo "❌ Failed to attach. Is monitor running? Try 'make up' first"
	@stty sane 2>/dev/null || true

status:
	@if pgrep -f $(MONITOR_BINARY) > /dev/null; then \
		echo "✅ Morgana Monitor is running"; \
		echo "Socket: /tmp/morgana.sock"; \
		ls -la /tmp/morgana.sock 2>/dev/null || echo "⚠️  Socket not found"; \
	else \
		echo "❌ Morgana Monitor is not running"; \
	fi

monitor: dev up
	@echo "Monitor started in headless mode. Use 'make attach' to view TUI"