#!/bin/bash
# Setup script to download and install WC3 sounds for Claude Code

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SOUNDS_DIR="$HOME/Sounds/game_samples"
TEMP_DIR="/tmp/wc3-sounds-$$"
SOUNDS_URL="https://www.hiveworkshop.com/attachments/sound-sets-zip.328586/"

echo -e "${BLUE}🎵 WC3 Sounds Setup for Claude Code${NC}"
echo -e "${BLUE}====================================${NC}\n"

# Check if sounds directory already exists
if [[ -d "$SOUNDS_DIR" ]] && [[ -n "$(ls -A "$SOUNDS_DIR" 2>/dev/null)" ]]; then
    echo -e "${YELLOW}⚠️  Sounds directory already exists and contains files.${NC}"
    read -p "Do you want to overwrite existing sounds? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}Setup cancelled.${NC}"
        exit 1
    fi
fi

# Create temporary directory
echo -e "${BLUE}Creating temporary directory...${NC}"
mkdir -p "$TEMP_DIR"
cd "$TEMP_DIR"

# Download the sounds archive
echo -e "${BLUE}Downloading WC3 sounds from HiveWorkshop...${NC}"
echo -e "${YELLOW}Note: This may take a moment depending on your connection speed.${NC}"

# Try to download using curl or wget
if command -v curl &> /dev/null; then
    curl -L -o "sound-sets.zip" "$SOUNDS_URL" || {
        echo -e "${RED}❌ Failed to download sounds archive.${NC}"
        echo -e "${YELLOW}Please download manually from:${NC}"
        echo -e "${BLUE}https://www.hiveworkshop.com/threads/sound-sets-from-campaign-dialog.328584/${NC}"
        rm -rf "$TEMP_DIR"
        exit 1
    }
elif command -v wget &> /dev/null; then
    wget -O "sound-sets.zip" "$SOUNDS_URL" || {
        echo -e "${RED}❌ Failed to download sounds archive.${NC}"
        echo -e "${YELLOW}Please download manually from:${NC}"
        echo -e "${BLUE}https://www.hiveworkshop.com/threads/sound-sets-from-campaign-dialog.328584/${NC}"
        rm -rf "$TEMP_DIR"
        exit 1
    }
else
    echo -e "${RED}❌ Neither curl nor wget found. Please install one of them.${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Check if download was successful
if [[ ! -f "sound-sets.zip" ]] || [[ ! -s "sound-sets.zip" ]]; then
    echo -e "${RED}❌ Download failed or file is empty.${NC}"
    echo -e "${YELLOW}Please download manually from:${NC}"
    echo -e "${BLUE}https://www.hiveworkshop.com/threads/sound-sets-from-campaign-dialog.328584/${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo -e "${GREEN}✅ Download complete!${NC}"

# Extract the archive
echo -e "${BLUE}Extracting sounds...${NC}"
unzip -q "sound-sets.zip" || {
    echo -e "${RED}❌ Failed to extract archive.${NC}"
    rm -rf "$TEMP_DIR"
    exit 1
}

# Create sounds directory
echo -e "${BLUE}Creating sounds directory...${NC}"
mkdir -p "$SOUNDS_DIR"

# Find and move sound files
echo -e "${BLUE}Installing sound files...${NC}"

# The archive might have different structures, so let's find all mp3 files
find . -name "*.mp3" -o -name "*.wav" | while read -r file; do
    # Get the relative path from the current directory
    rel_path="${file#./}"
    # Create the directory structure in the target
    target_dir="$SOUNDS_DIR/$(dirname "$rel_path")"
    mkdir -p "$target_dir"
    # Copy the file
    cp "$file" "$target_dir/"
done

# Count installed files
SOUND_COUNT=$(find "$SOUNDS_DIR" -type f \( -name "*.mp3" -o -name "*.wav" \) | wc -l)

# Clean up
echo -e "${BLUE}Cleaning up temporary files...${NC}"
cd /
rm -rf "$TEMP_DIR"

# Verify installation
if [[ $SOUND_COUNT -gt 0 ]]; then
    echo -e "${GREEN}✅ Successfully installed $SOUND_COUNT sound files!${NC}"
    echo -e "${GREEN}📁 Sounds installed to: $SOUNDS_DIR${NC}"
    
    # Test a sound if afplay is available
    if command -v afplay &> /dev/null && [[ -f "$SOUNDS_DIR/jobs-done_1.mp3" ]]; then
        echo -e "\n${BLUE}Testing sound playback...${NC}"
        afplay "$SOUNDS_DIR/jobs-done_1.mp3"
        echo -e "${GREEN}✅ Sound test complete!${NC}"
    fi
    
    echo -e "\n${GREEN}🎉 WC3 sounds are now ready for Claude Code!${NC}"
else
    echo -e "${RED}❌ No sound files were installed.${NC}"
    echo -e "${YELLOW}The archive structure might be different than expected.${NC}"
    echo -e "${YELLOW}Please extract manually to: $SOUNDS_DIR${NC}"
    exit 1
fi

# Show some example sounds that were installed
echo -e "\n${BLUE}Sample sounds installed:${NC}"
find "$SOUNDS_DIR" -type f \( -name "*.mp3" -o -name "*.wav" \) | head -10 | while read -r file; do
    echo -e "  ${GREEN}✓${NC} ${file#$SOUNDS_DIR/}"
done

echo -e "\n${YELLOW}Note: Claude Code settings.json is already configured to use these sounds.${NC}"