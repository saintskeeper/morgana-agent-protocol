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
    
    # Check if essential hooks are preserved (handle both old and new format)
    local has_validate=$(jq -r '
        if .hooks.PreToolUse | type == "array" then
            .hooks.PreToolUse[]?.hooks[]? | select(.command | contains("validate-claude.sh")) | .command
        else
            .hooks.PreToolUse | to_entries[] | select(.value | tostring | contains("validate-claude.sh")) | .value
        end' "$SETTINGS_FILE" 2>/dev/null | head -1)
    
    local has_post_edit=$(jq -r '
        if .hooks.PostToolUse | type == "array" then
            .hooks.PostToolUse[]?.hooks[]? | select(.command | contains("post-edit.sh")) | .command
        else
            .hooks.PostToolUse | to_entries[] | select(.value | if type == "array" then .[] else . end | tostring | contains("post-edit.sh")) | .value | if type == "array" then .[0] else . end
        end' "$SETTINGS_FILE" 2>/dev/null | head -1)
    
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
    
    # Check if backup already exists
    if [[ -f "$BACKUP_FILE" ]]; then
        echo -e "${YELLOW}⚠️  Existing backup found at: $BACKUP_FILE${NC}"
        
        # Count sounds in existing backup
        local backup_sounds=$(jq '
            [.. | 
                if type == "object" and has("command") then
                    .command | select(. != null and contains("afplay"))
                elif type == "string" then
                    . | select(contains("afplay"))
                elif type == "array" then
                    .[] | select(type == "string" and contains("afplay"))
                else
                    empty
                end
            ] | length' "$BACKUP_FILE" 2>/dev/null || echo "0")
        
        # Count sounds in current settings
        local current_sounds=$(jq '
            [.. | 
                if type == "object" and has("command") then
                    .command | select(. != null and contains("afplay"))
                elif type == "string" then
                    . | select(contains("afplay"))
                elif type == "array" then
                    .[] | select(type == "string" and contains("afplay"))
                else
                    empty
                end
            ] | length' "$SETTINGS_FILE" 2>/dev/null || echo "0")
        
        echo -e "${BLUE}Backup has $backup_sounds sounds, current settings have $current_sounds sounds${NC}"
        
        if [[ $current_sounds -gt 0 ]]; then
            # If current settings have sounds, merge them with backup
            echo -e "${GREEN}Merging current sounds with backup...${NC}"
            
            # Create a temporary merged backup
            local MERGE_FILE="$HOME/.claude/settings.json.sound-backup.merged"
            
            # Extract all sound hooks from both files and merge
            jq -s '
                # Function to extract sound commands with paths
                def extract_sounds:
                    [paths(scalars) as $p | 
                        getpath($p) as $val |
                        if ($val | type == "string" and contains("afplay")) then
                            {path: $p, value: $val}
                        else empty end
                    ] + 
                    [paths(objects) as $p | 
                        getpath($p) as $obj |
                        if ($obj | type == "object" and has("command") and (.command | type == "string" and contains("afplay"))) then
                            {path: $p, value: $obj}
                        else empty end
                    ];
                
                # Get sounds from both files
                (.[0] | extract_sounds) as $backup_sounds |
                (.[1] | extract_sounds) as $current_sounds |
                
                # Start with the backup structure
                .[0] as $base |
                
                # Add any new sounds from current that aren'\''t in backup
                reduce ($current_sounds[] | select(. as $curr | $backup_sounds | map(.value) | index($curr.value) | not)) as $sound
                    ($base; setpath($sound.path; $sound.value))
            ' "$BACKUP_FILE" "$SETTINGS_FILE" > "$MERGE_FILE"
            
            # Validate merged file
            if validate_json "$MERGE_FILE"; then
                mv "$MERGE_FILE" "$BACKUP_FILE"
                echo -e "${GREEN}✓ Backup updated with merged sounds${NC}"
            else
                echo -e "${RED}Warning: Failed to merge, keeping original backup${NC}"
                rm -f "$MERGE_FILE"
            fi
        else
            echo -e "${BLUE}Current settings have no sounds, keeping existing backup${NC}"
        fi
    else
        # No existing backup, create new one
        cp "$SETTINGS_FILE" "$BACKUP_FILE"
        echo -e "${BLUE}Backup created at: $BACKUP_FILE${NC}"
    fi
    
    # Use jq to remove only afplay commands while preserving all other functionality
    if command -v jq &> /dev/null; then
        # Handle both old array format and new object format
        jq '
        .hooks |= with_entries(
            .value |= 
                if type == "array" then
                    # Old format: array of matchers
                    map(
                        if .hooks then
                            .hooks |= map(
                                select(.command | type == "string" and (contains("afplay") | not))
                            ) |
                            if (.hooks | length) == 0 then empty else . end
                        else
                            .
                        end
                    )
                elif type == "object" then
                    # New format: object with tool matchers as keys
                    with_entries(
                        .value |= 
                            if type == "array" then
                                map(select(type == "string" and (contains("afplay") | not)))
                            elif type == "string" then
                                if contains("afplay") then empty else . end
                            else
                                .
                            end |
                        select(.value != [] and .value != null)
                    )
                else
                    # Simple string format for hooks like Stop, Notification
                    if type == "string" and contains("afplay") then
                        null
                    else
                        .
                    end
                end
        ) | .hooks |= with_entries(select(.value != null and .value != {} and .value != []))' "$SETTINGS_FILE" > "$TEMP_FILE"
        
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
    
    # Check if sounds are already enabled
    local current_sounds=$(jq '
        [.. | 
            if type == "object" and has("command") then
                .command | select(. != null and contains("afplay"))
            elif type == "string" then
                . | select(contains("afplay"))
            elif type == "array" then
                .[] | select(type == "string" and contains("afplay"))
            else
                empty
            end
        ] | length' "$SETTINGS_FILE" 2>/dev/null || echo "0")
    
    if [[ $current_sounds -gt 0 ]]; then
        echo -e "${YELLOW}⚠️  Sounds are already enabled (found $current_sounds sound effects)${NC}"
        echo -e "${BLUE}Current settings already contain sounds. Options:${NC}"
        echo "1. Keep current sounds as-is"
        echo "2. Replace with backup (may lose recently added sounds)"
        echo "3. Merge backup with current (combine all sounds)"
        
        read -p "Choose option [1-3]: " choice
        
        case $choice in
            1)
                echo -e "${GREEN}Keeping current sounds${NC}"
                test_functionality
                return 0
                ;;
            2)
                echo -e "${YELLOW}Replacing with backup...${NC}"
                ;;
            3)
                echo -e "${GREEN}Merging sounds...${NC}"
                # Create a temporary file for merge
                local MERGE_FILE="$HOME/.claude/settings.json.merged"
                
                # Merge sounds from backup into current settings
                jq -s '
                    # Function to extract non-sound hooks
                    def extract_non_sounds:
                        . as $root |
                        if .hooks then
                            .hooks |= with_entries(
                                .value |= 
                                    if type == "object" then
                                        with_entries(
                                            .value |= 
                                                if type == "array" then
                                                    map(select(type == "string" and (contains("afplay") | not)))
                                                elif type == "string" then
                                                    if contains("afplay") then empty else . end
                                                else .
                                                end |
                                            select(.value != [] and .value != null)
                                        )
                                    elif type == "array" then
                                        map(
                                            if .hooks then
                                                .hooks |= map(
                                                    select(.command | type == "string" and (contains("afplay") | not))
                                                )
                                            else . end
                                        )
                                    elif type == "string" then
                                        if contains("afplay") then null else . end
                                    else .
                                    end
                            ) | .hooks |= with_entries(select(.value != null and .value != {} and .value != []))
                        else . end;
                    
                    # Start with current settings without sounds
                    .[1] | extract_non_sounds as $base |
                    
                    # Get all hooks from backup
                    .[0] as $backup |
                    
                    # Merge backup hooks into base
                    $base * $backup
                ' "$BACKUP_FILE" "$SETTINGS_FILE" > "$MERGE_FILE"
                
                if validate_json "$MERGE_FILE"; then
                    mv "$MERGE_FILE" "$SETTINGS_FILE"
                    echo -e "${GREEN}✓ Successfully merged sounds${NC}"
                else
                    echo -e "${RED}Error: Failed to merge, aborting${NC}"
                    rm -f "$MERGE_FILE"
                    exit 1
                fi
                
                # Count merged sounds
                local merged_sounds=$(jq '
                    [.. | 
                        if type == "object" and has("command") then
                            .command | select(. != null and contains("afplay"))
                        elif type == "string" then
                            . | select(contains("afplay"))
                        elif type == "array" then
                            .[] | select(type == "string" and contains("afplay"))
                        else
                            empty
                        end
                    ] | length' "$SETTINGS_FILE" 2>/dev/null || echo "0")
                echo -e "${BLUE}Merged to $merged_sounds total sound effects${NC}"
                test_functionality
                return 0
                ;;
            *)
                echo -e "${RED}Invalid choice, aborting${NC}"
                exit 1
                ;;
        esac
    fi
    
    # Validate backup before restoring
    if validate_json "$BACKUP_FILE"; then
        cp "$BACKUP_FILE" "$SETTINGS_FILE"
        echo -e "${GREEN}✓ Sounds have been restored from backup${NC}"
        
        # Count restored sounds (handle both old and new format)
        local sound_count=$(jq '
            [.. | 
                if type == "object" and has("command") then
                    .command | select(. != null and contains("afplay"))
                elif type == "string" then
                    . | select(contains("afplay"))
                elif type == "array" then
                    .[] | select(type == "string" and contains("afplay"))
                else
                    empty
                end
            ] | length' "$SETTINGS_FILE" 2>/dev/null || echo "0")
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