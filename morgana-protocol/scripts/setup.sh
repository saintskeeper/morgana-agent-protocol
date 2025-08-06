#!/bin/bash
set -e

# Morgana Protocol Setup Script
# This script builds and installs the Morgana Protocol CLI with TUI support

BINARY_NAME="morgana"
INSTALL_DIR="/usr/local/bin"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS")
            echo -e "${GREEN}✅ SUCCESS${NC}: $message"
            ;;
        "ERROR")
            echo -e "${RED}❌ ERROR${NC}: $message"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠️  WARNING${NC}: $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  INFO${NC}: $message"
            ;;
    esac
}

# Check if running as root for /usr/local/bin
check_permissions() {
    if [ "$EUID" -ne 0 ] && [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        print_status "ERROR" "Installation to $INSTALL_DIR requires root privileges"
        echo -e "${YELLOW}Please run with sudo:${NC} sudo $0"
        echo -e "${YELLOW}Or choose a different install location by setting INSTALL_DIR:${NC}"
        echo "  INSTALL_DIR=\$HOME/.local/bin $0"
        exit 1
    fi
}

echo -e "${BLUE}🧙‍♂️ Morgana Protocol Setup${NC}"
echo "=================================="
print_status "INFO" "Building and installing Morgana Protocol CLI with TUI support"
print_status "INFO" "Repository: $REPO_ROOT"
print_status "INFO" "Target: $INSTALL_DIR/$BINARY_NAME"

# Navigate to repository root
cd "$REPO_ROOT"

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/morgana" ]; then
    print_status "ERROR" "Not in a Morgana Protocol repository root"
    exit 1
fi

# Check permissions
check_permissions

# Verify Go is installed
if ! command -v go &> /dev/null; then
    print_status "ERROR" "Go is not installed or not in PATH"
    print_status "INFO" "Install Go from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_status "INFO" "Using Go version: $GO_VERSION"

# Check Go version (require 1.21+)
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
    print_status "WARNING" "Go 1.21+ recommended, you have $GO_VERSION"
fi

# Clean previous builds
print_status "INFO" "Cleaning previous builds..."
rm -f morgana morgana-*

# Update dependencies
print_status "INFO" "Updating dependencies..."
if ! go mod tidy; then
    print_status "ERROR" "Failed to update Go modules"
    exit 1
fi

# Detect OS and architecture for build tags
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    i386|i686) ARCH="386" ;;
    *) print_status "WARNING" "Unsupported architecture: $ARCH, using amd64" 
       ARCH="amd64" ;;
esac

print_status "INFO" "Building for $OS/$ARCH..."

# Build the binary with optimizations and version info
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION="v1.0.0-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")"

LDFLAGS="-w -s -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

print_status "INFO" "Version: $VERSION"
print_status "INFO" "Build flags: $LDFLAGS"

# Build with CGO disabled for maximum compatibility
if ! CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build \
    -ldflags "$LDFLAGS" \
    -trimpath \
    -tags "netgo osusergo static_build" \
    -o "$BINARY_NAME" \
    ./cmd/morgana; then
    print_status "ERROR" "Failed to build Morgana Protocol CLI"
    exit 1
fi

print_status "SUCCESS" "Build completed successfully"

# Verify the binary works
print_status "INFO" "Testing binary..."
if ! ./"$BINARY_NAME" --version; then
    print_status "ERROR" "Binary test failed"
    exit 1
fi

print_status "SUCCESS" "Binary test passed"

# Create install directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    print_status "INFO" "Creating install directory: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
fi

# Backup existing binary if it exists
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    print_status "INFO" "Backing up existing binary..."
    cp "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME.backup.$(date +%s)"
fi

# Install the binary
print_status "INFO" "Installing to $INSTALL_DIR..."
cp "$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

print_status "SUCCESS" "Morgana Protocol CLI installed to $INSTALL_DIR/$BINARY_NAME"

# Verify installation
if command -v morgana &> /dev/null; then
    INSTALLED_VERSION=$(morgana --version 2>/dev/null | head -n1 || echo "unknown")
    print_status "SUCCESS" "Installation verified: $INSTALLED_VERSION"
else
    print_status "WARNING" "$INSTALL_DIR is not in PATH"
    echo -e "${YELLOW}Add the following to your shell profile:${NC}"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
fi

# Clean up build artifacts
rm -f "$BINARY_NAME"

echo ""
echo -e "${GREEN}🎉 Installation Complete!${NC}"
echo "=================================="
echo ""
echo -e "${BLUE}Quick Start Commands:${NC}"
echo "  morgana --version                    # Check installation"
echo "  morgana --help                       # Show all options"
echo "  morgana --tui                        # Start with TUI interface"
echo "  morgana --config morgana.yaml --tui  # Use your config with TUI"
echo ""
echo -e "${BLUE}TUI Modes:${NC}"
echo "  morgana --tui --tui-mode dev                    # Development mode"
echo "  morgana --tui --tui-mode optimized             # Balanced (default)"
echo "  morgana --tui --tui-mode high-performance      # Production mode"
echo ""
echo -e "${BLUE}Example Usage:${NC}"
echo "  morgana --tui --agent code-implementer --prompt \"Build REST API\""
echo "  echo '[{\"agent_type\":\"test-specialist\",\"prompt\":\"Generate tests\"}]' | morgana --tui"
echo ""
echo -e "${BLUE}Configuration:${NC}"
echo "  Edit morgana.yaml to customize TUI settings, agents, and telemetry"
echo "  Set MORGANA_DEBUG=true for additional performance monitoring"
echo ""
echo -e "${BLUE}Documentation:${NC}"
echo "  📖 TUI User Guide: TUI_USER_GUIDE.md"
echo "  📋 Integration Guide: MONITORING_INTEGRATION.md" 
echo "  🔧 Repository: https://github.com/saintskeeper/claude-code-configs"
echo ""

# Optional: Install shell completions if the binary supports it
if morgana completion --help &>/dev/null; then
    print_status "INFO" "Shell completions available - run 'morgana completion --help'"
fi

# Optional: Check for updates
if command -v git &> /dev/null && [ -d ".git" ]; then
    if [ "$(git status --porcelain 2>/dev/null)" ]; then
        print_status "WARNING" "Repository has uncommitted changes"
    fi
    
    # Check if we're behind remote
    LOCAL=$(git rev-parse @ 2>/dev/null || echo "")
    REMOTE=$(git rev-parse @{u} 2>/dev/null || echo "")
    if [ "$LOCAL" != "$REMOTE" ] && [ -n "$REMOTE" ]; then
        print_status "INFO" "Updates available - run 'git pull' and reinstall"
    fi
fi

print_status "SUCCESS" "Morgana Protocol is ready to use! 🧙‍♂️"