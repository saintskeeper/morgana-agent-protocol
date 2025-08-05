#!/bin/bash
set -e

INSTALL_DIR="$HOME/.claude/bin"
BINARY_NAME="morgana"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Installing Morgana Protocol..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Source binary
SOURCE="$REPO_ROOT/dist/${BINARY_NAME}-${OS}-${ARCH}"
if [[ "$OS" == "windows" ]]; then
    SOURCE="${SOURCE}.exe"
fi

# Check if binary exists
if [ ! -f "$SOURCE" ]; then
    echo "Binary not found: $SOURCE"
    echo "Please run 'make build' first"
    exit 1
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

# Copy binary
cp "$SOURCE" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "✓ Morgana Protocol installed to $INSTALL_DIR/$BINARY_NAME"

# Check if PATH includes install directory
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "⚠️  Add $INSTALL_DIR to your PATH:"
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
    echo "   Add this to your ~/.bashrc or ~/.zshrc to make it permanent"
fi

echo ""
echo "Test installation with: morgana --version"