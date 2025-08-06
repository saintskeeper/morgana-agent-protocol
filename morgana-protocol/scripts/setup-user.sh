#!/bin/bash
set -e

# Morgana Protocol User Setup Script
# This script builds and installs the Morgana Protocol CLI to ~/.local/bin (no sudo required)

BINARY_NAME="morgana"
INSTALL_DIR="$HOME/.local/bin"
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

echo -e "${BLUE}🧙‍♂️ Morgana Protocol User Setup${NC}"
echo "===================================="
print_status "INFO" "Building and installing Morgana Protocol CLI with TUI support"
print_status "INFO" "Repository: $REPO_ROOT"
print_status "INFO" "Target: $INSTALL_DIR/$BINARY_NAME (user installation)"

# Navigate to repository root
cd "$REPO_ROOT"

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/morgana" ]; then
    print_status "ERROR" "Not in a Morgana Protocol repository root"
    exit 1
fi

# Verify Go is installed
if ! command -v go &> /dev/null; then
    print_status "ERROR" "Go is not installed or not in PATH"
    print_status "INFO" "Install Go from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_status "INFO" "Using Go version: $GO_VERSION"

# Clean previous builds
print_status "INFO" "Cleaning previous builds..."
rm -f morgana morgana-*

# Update dependencies
print_status "INFO" "Updating dependencies..."
if ! go mod tidy; then
    print_status "ERROR" "Failed to update Go modules"
    exit 1
fi

# Detect OS and architecture
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

# Build version info
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION="v1.0.0-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")"

LDFLAGS="-w -s -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

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

# Test the binary
print_status "INFO" "Testing binary..."
if ! ./"$BINARY_NAME" --version; then
    print_status "ERROR" "Binary test failed"
    exit 1
fi

print_status "SUCCESS" "Binary test passed"

# Create install directory
mkdir -p "$INSTALL_DIR"

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

# Check PATH and update shell profile if needed
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_status "WARNING" "$INSTALL_DIR is not in your PATH"
    
    # Try to add to shell profile
    SHELL_PROFILE=""
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_PROFILE="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        if [ -f "$HOME/.bashrc" ]; then
            SHELL_PROFILE="$HOME/.bashrc"
        elif [ -f "$HOME/.bash_profile" ]; then
            SHELL_PROFILE="$HOME/.bash_profile"
        fi
    fi
    
    if [ -n "$SHELL_PROFILE" ]; then
        echo ""
        echo -e "${YELLOW}Would you like to add $INSTALL_DIR to your PATH automatically? [y/N]${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_PROFILE"
            print_status "SUCCESS" "Added $INSTALL_DIR to $SHELL_PROFILE"
            print_status "INFO" "Restart your shell or run: source $SHELL_PROFILE"
        else
            echo -e "${YELLOW}Add the following to your shell profile manually:${NC}"
            echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
        fi
    else
        echo -e "${YELLOW}Add the following to your shell profile:${NC}"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
else
    # Verify installation
    if command -v morgana &> /dev/null; then
        INSTALLED_VERSION=$(morgana --version 2>/dev/null | head -n1 || echo "unknown")
        print_status "SUCCESS" "Installation verified: $INSTALLED_VERSION"
    fi
fi

# Clean up build artifacts
rm -f "$BINARY_NAME"

echo ""
echo -e "${GREEN}🎉 User Installation Complete!${NC}"
echo "======================================"
echo ""
echo -e "${BLUE}Quick Start Commands:${NC}"
echo "  $INSTALL_DIR/morgana --version                    # Check installation"
echo "  $INSTALL_DIR/morgana --tui                        # Start with TUI interface"
echo "  $INSTALL_DIR/morgana --config morgana.yaml --tui  # Use your config with TUI"
echo ""
echo -e "${BLUE}If PATH is updated:${NC}"
echo "  morgana --tui --tui-mode dev                    # Development mode"
echo "  morgana --tui --tui-mode optimized             # Balanced (default)"
echo "  morgana --tui --tui-mode high-performance      # Production mode"
echo ""
echo -e "${BLUE}Documentation:${NC}"
echo "  📖 TUI User Guide: TUI_USER_GUIDE.md"
echo "  📋 Integration Guide: MONITORING_INTEGRATION.md"
echo ""

print_status "SUCCESS" "Morgana Protocol is ready to use! 🧙‍♂️"