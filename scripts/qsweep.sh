#!/bin/bash

# qsweep - Quick sweep function for organizing AI documentation
# Usage: ./qsweep.sh [--id TICKET-ID] [--type TYPE] [--dry-run]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AI_DOCS_DIR="docs"
COMPLETED_DIR="$AI_DOCS_DIR"
ACTIVE_DIR="$AI_DOCS_DIR/active"

# Search paths for documentation
DOC_PATHS=(
    "back-end-go/docs"
    "front-end-next/docs"
    "operator/docs"
    "manifests/helm-charts/spellcarver-app"
)

# Valid commit types
VALID_TYPES=("feat" "fix" "docs" "chore" "refactor" "test" "build" "ci" "perf" "style")

# Parse command line arguments
DRY_RUN=false
TICKET_ID=""
DOC_TYPE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --id)
            TICKET_ID="$2"
            shift 2
            ;;
        --type)
            DOC_TYPE="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help)
            echo "Usage: $0 [--id TICKET-ID] [--type TYPE] [--dry-run]"
            echo ""
            echo "Options:"
            echo "  --id TICKET-ID    Sweep specific ticket (e.g., SPE-249)"
            echo "  --type TYPE       Sweep by type (feat, fix, docs, chore, refactor, test)"
            echo "  --dry-run         Show what would be moved without actually moving"
            echo ""
            echo "Valid types: ${VALID_TYPES[*]}"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Validate type if provided
if [[ -n "$DOC_TYPE" ]] && [[ ! " ${VALID_TYPES[@]} " =~ " ${DOC_TYPE} " ]]; then
    echo -e "${RED}Invalid type: $DOC_TYPE${NC}"
    echo "Valid types: ${VALID_TYPES[*]}"
    exit 1
fi

# Function to extract ticket ID from filename
extract_ticket_id() {
    local filename="$1"
    if [[ $filename =~ (SPE-[0-9]+|spe-[0-9]+) ]]; then
        echo "${BASH_REMATCH[1],,}" # Convert to lowercase
    fi
}

# Function to determine document type from content
determine_doc_type() {
    local file="$1"
    local content=$(head -n 50 "$file" 2>/dev/null || echo "")

    # Check for common patterns
    if [[ $content =~ "new feature"|"feature implementation"|"Feature:" ]]; then
        echo "feat"
    elif [[ $content =~ "bug fix"|"Bug:"|"Issue:"|"Fix:" ]]; then
        echo "fix"
    elif [[ $content =~ "documentation"|"README"|"guide"|"tutorial" ]]; then
        echo "docs"
    elif [[ $content =~ "refactor"|"code improvement"|"optimization" ]]; then
        echo "refactor"
    elif [[ $content =~ "test"|"testing"|"unit test"|"integration test" ]]; then
        echo "test"
    else
        # Default to feat for feature-related docs
        echo "feat"
    fi
}

# Function to move a document
move_document() {
    local source="$1"
    local dest_dir="$2"
    local filename=$(basename "$source")

    if [[ $DRY_RUN == true ]]; then
        echo -e "${YELLOW}[DRY RUN]${NC} Would move: $source -> $dest_dir/$filename"
    else
        mkdir -p "$dest_dir"
        mv "$source" "$dest_dir/"
        echo -e "${GREEN}Moved:${NC} $source -> $dest_dir/$filename"
    fi
}

# Main sweep logic
echo -e "${BLUE}Starting qsweep...${NC}"

# Find all markdown files in doc paths
found_files=()
for path in "${DOC_PATHS[@]}"; do
    if [[ -d "$path" ]]; then
        while IFS= read -r -d '' file; do
            found_files+=("$file")
        done < <(find "$path" -type f -name "*.md" -print0)
    fi
done

if [[ ${#found_files[@]} -eq 0 ]]; then
    echo -e "${YELLOW}No documentation files found to sweep.${NC}"
    exit 0
fi

echo -e "${BLUE}Found ${#found_files[@]} documentation files${NC}"

# Process each file
moved_count=0
for file in "${found_files[@]}"; do
    filename=$(basename "$file")
    ticket_id=$(extract_ticket_id "$filename")

    # Skip if filtering by ticket ID and doesn't match
    if [[ -n "$TICKET_ID" ]] && [[ "$ticket_id" != "${TICKET_ID,,}" ]]; then
        continue
    fi

    # Determine document type
    if [[ -n "$DOC_TYPE" ]]; then
        doc_type="$DOC_TYPE"
    else
        doc_type=$(determine_doc_type "$file")
    fi

    # Skip if filtering by type and doesn't match
    if [[ -n "$DOC_TYPE" ]] && [[ "$doc_type" != "$DOC_TYPE" ]]; then
        continue
    fi

    # Create destination directory
    if [[ -n "$ticket_id" ]]; then
        dest_dir="$COMPLETED_DIR/$doc_type/$ticket_id"
    else
        # For files without ticket IDs, use a generic name
        dest_dir="$COMPLETED_DIR/$doc_type/misc"
    fi

    # Move the document
    move_document "$file" "$dest_dir"
    ((moved_count++))
done

# Summary
if [[ $DRY_RUN == true ]]; then
    echo -e "\n${YELLOW}[DRY RUN] Would move $moved_count files${NC}"
else
    echo -e "\n${GREEN}Successfully moved $moved_count files${NC}"
fi

# Show current structure
echo -e "\n${BLUE}Current docs structure:${NC}"
tree "$AI_DOCS_DIR" 2>/dev/null || find "$AI_DOCS_DIR" -type d | sort
