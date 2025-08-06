# Makefile for Morgana Protocol
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.Version=$(VERSION)
BINARY := morgana

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOGET := $(GOCMD) get

# Directories
CMD_DIR := ./cmd/morgana
DIST_DIR := ./dist

.PHONY: all build test clean deps install dev run help

# Default target
all: test build

# Help target
help:
	@echo "Morgana Protocol - Agent Orchestration System"
	@echo ""
	@echo "Available targets:"
	@echo "  make build       - Build for all platforms"
	@echo "  make dev         - Build for current platform only"
	@echo "  make test        - Run all tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deps        - Download dependencies"
	@echo "  make install     - Install locally to ~/.claude/bin"
	@echo "  make run         - Run with example"
	@echo ""

# Development build (current platform only)
dev:
	@echo "Building for development..."
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY) $(CMD_DIR)

# Build for all platforms
build: build-darwin build-linux build-windows

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-darwin-arm64 $(CMD_DIR)

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-linux-arm64 $(CMD_DIR)

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-windows-amd64.exe $(CMD_DIR)

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
	$(GOTEST) -v -tags=integration ./...

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
	@chmod +x ~/.claude/bin/$(BINARY)
	@echo "Installed to ~/.claude/bin/$(BINARY)"

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