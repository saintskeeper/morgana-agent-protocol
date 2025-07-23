#!/bin/bash

# qsound - Quick sound toggle for Claude Code
# Usage: ./qsound.sh [silence|enable]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SETTINGS_FILE="$HOME/.claude/settings.json"
BACKUP_FILE="$HOME/.claude/settings.json.sound-backup"
TEMP_FILE="$HOME/.claude/settings.json.tmp"

# Function to display usage
usage() {
    echo "Usage: $0 [silence|enable]"
    echo ""
    echo "Commands:"
    echo "  silence  - Disable all sounds (backs up current settings)"
    echo "  enable   - Re-enable sounds (restores from backup)"
    echo ""
    exit 1
}

# Check if settings file exists
if [[ ! -f "$SETTINGS_FILE" ]]; then
    echo -e "${RED}Error: settings.json not found at $SETTINGS_FILE${NC}"
    exit 1
fi

# Function to validate JSON
validate_json() {
    local file="$1"
    if command -v jq &> /dev/null; then
        if jq empty "$file" 2>/dev/null; then
            return 0
        else
            return 1
        fi
    else
        # Fallback to python if jq not available
        if command -v python3 &> /dev/null; then
            python3 -m json.tool "$file" > /dev/null 2>&1
            return $?
        fi
    fi
    return 1
}

# Function to test Claude functionality
test_functionality() {
    echo -e "${BLUE}Testing Claude Code functionality...${NC}"
    
    # Check if essential hooks are preserved
    local has_validate=$(jq -r '.hooks.PreToolUse[]?.hooks[]? | select(.command | contains("validate-claude.sh")) | .command' "$SETTINGS_FILE" 2>/dev/null)
    local has_post_edit=$(jq -r '.hooks.PostToolUse[]?.hooks[]? | select(.command | contains("post-edit.sh")) | .command' "$SETTINGS_FILE" 2>/dev/null)
    
    if [[ -n "$has_validate" ]]; then
        echo -e "${GREEN}✓ Validation hooks preserved${NC}"
    else
        echo -e "${YELLOW}⚠ Warning: Validation hooks might be missing${NC}"
    fi
    
    if [[ -n "$has_post_edit" ]]; then
        echo -e "${GREEN}✓ Post-edit hooks preserved${NC}"
    else
        echo -e "${YELLOW}⚠ Warning: Post-edit hooks might be missing${NC}"
    fi
}

# Function to remove all sound commands from settings.json
silence_sounds() {
    echo -e "${YELLOW}🔇 Silencing all Claude Code sounds...${NC}"
    
    # Create backup
    cp "$SETTINGS_FILE" "$BACKUP_FILE"
    echo -e "${BLUE}Backup created at: $BACKUP_FILE${NC}"
    
    # Use jq to remove only afplay commands while preserving all other functionality
    if command -v jq &> /dev/null; then
        # Create a more sophisticated filter that preserves hook structure
        jq '
        .hooks |= with_entries(
            .value |= map(
                if .hooks then
                    .hooks |= map(
                        select(.command | type == "string" and (contains("afplay") | not))
                    ) |
                    if (.hooks | length) == 0 then empty else . end
                else
                    .
                end
            )
        )' "$SETTINGS_FILE" > "$TEMP_FILE"
        
        # Validate the new JSON
        if validate_json "$TEMP_FILE"; then
            mv "$TEMP_FILE" "$SETTINGS_FILE"
            echo -e "${GREEN}✓ All sounds have been disabled${NC}"
            
            # Test functionality
            test_functionality
            
            echo -e "${YELLOW}Your settings are backed up. Use '$0 enable' to restore sounds${NC}"
        else
            echo -e "${RED}Error: Failed to create valid JSON${NC}"
            rm -f "$TEMP_FILE"
            exit 1
        fi
    else
        echo -e "${RED}Error: jq is required but not installed${NC}"
        echo "Install with: brew install jq"
        exit 1
    fi
}

# Function to restore sounds from backup
enable_sounds() {
    echo -e "${GREEN}🔊 Re-enabling Claude Code sounds...${NC}"
    
    # Check if backup exists
    if [[ ! -f "$BACKUP_FILE" ]]; then
        echo -e "${RED}Error: No backup found at $BACKUP_FILE${NC}"
        echo "Cannot restore sounds without a backup"
        exit 1
    fi
    
    # Validate backup before restoring
    if validate_json "$BACKUP_FILE"; then
        cp "$BACKUP_FILE" "$SETTINGS_FILE"
        echo -e "${GREEN}✓ Sounds have been restored from backup${NC}"
        
        # Count restored sounds
        local sound_count=$(jq '[.. | .command? | select(. != null and contains("afplay"))] | length' "$SETTINGS_FILE" 2>/dev/null || echo "0")
        echo -e "${BLUE}Restored $sound_count sound effects${NC}"
        
        # Test functionality
        test_functionality
    else
        echo -e "${RED}Error: Backup file is corrupted${NC}"
        exit 1
    fi
    
    # Keep backup for safety
    echo -e "${YELLOW}Backup preserved at: $BACKUP_FILE${NC}"
}

# Main logic
case "${1:-}" in
    silence)
        silence_sounds
        ;;
    enable)
        enable_sounds
        ;;
    *)
        usage
        ;;
esac