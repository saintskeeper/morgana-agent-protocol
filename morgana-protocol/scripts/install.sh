#!/bin/bash
set -e

# Morgana Protocol Install Script (for pre-built binaries)
# This script installs pre-built binaries from dist/ directory

INSTALL_DIR="$HOME/.claude/bin"
BINARY_NAME="morgana"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS") echo -e "${GREEN}✅ SUCCESS${NC}: $message" ;;
        "ERROR") echo -e "${RED}❌ ERROR${NC}: $message" ;;
        "WARNING") echo -e "${YELLOW}⚠️  WARNING${NC}: $message" ;;
        "INFO") echo -e "${BLUE}ℹ️  INFO${NC}: $message" ;;
    esac
}

echo -e "${BLUE}🧙‍♂️ Morgana Protocol Install${NC}"
echo "================================="
print_status "INFO" "Installing pre-built Morgana Protocol CLI with TUI support"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    i386|i686) ARCH="386" ;;
    *) print_status "ERROR" "Unsupported architecture: $ARCH"; exit 1 ;;
esac

print_status "INFO" "Detected platform: $OS/$ARCH"

# Source binary
SOURCE="$REPO_ROOT/dist/${BINARY_NAME}-${OS}-${ARCH}"
if [[ "$OS" == "windows" ]]; then
    SOURCE="${SOURCE}.exe"
fi

# Check if binary exists
if [ ! -f "$SOURCE" ]; then
    print_status "ERROR" "Pre-built binary not found: $SOURCE"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  1. Run 'make build' to build binaries first"
    echo "  2. Use './scripts/setup.sh' to build and install from source"
    echo "  3. Use './scripts/setup-user.sh' for user installation"
    exit 1
fi

print_status "SUCCESS" "Found pre-built binary: $SOURCE"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Backup existing binary if it exists
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    print_status "INFO" "Backing up existing binary..."
    cp "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME.backup.$(date +%s)"
fi

# Copy binary
cp "$SOURCE" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

print_status "SUCCESS" "Morgana Protocol installed to $INSTALL_DIR/$BINARY_NAME"

# Test the installation
if "$INSTALL_DIR/$BINARY_NAME" --version &>/dev/null; then
    VERSION=$("$INSTALL_DIR/$BINARY_NAME" --version | head -n1)
    print_status "SUCCESS" "Installation verified: $VERSION"
else
    print_status "WARNING" "Binary installed but version check failed"
fi

# Check if PATH includes install directory
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_status "WARNING" "$INSTALL_DIR is not in your PATH"
    echo ""
    echo -e "${YELLOW}Add $INSTALL_DIR to your PATH:${NC}"
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
    echo -e "${YELLOW}Add this to your ~/.bashrc or ~/.zshrc to make it permanent${NC}"
    echo ""
fi

echo ""
echo -e "${GREEN}🎉 Installation Complete!${NC}"
echo "=========================="
echo ""
echo -e "${BLUE}Quick Start:${NC}"
echo "  morgana --version                    # Check installation"
echo "  morgana --tui                        # Start with TUI interface"
echo "  morgana --config morgana.yaml --tui  # Use config with TUI"
echo ""
echo -e "${BLUE}TUI Modes:${NC}"
echo "  morgana --tui --tui-mode dev                    # Development mode"
echo "  morgana --tui --tui-mode optimized             # Balanced (default)"
echo "  morgana --tui --tui-mode high-performance      # Production mode"
echo ""
echo -e "${BLUE}Documentation:${NC}"
echo "  📖 TUI User Guide: TUI_USER_GUIDE.md"
echo "  📋 Integration Guide: MONITORING_INTEGRATION.md"
echo ""

print_status "SUCCESS" "Morgana Protocol is ready to use! 🧙‍♂️"