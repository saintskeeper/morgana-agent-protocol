#!/bin/bash
# Claude hook: Runs after creating a new branch
# This hook automatically sweeps documentation to ai-docs when starting new work

echo "🧹 Running qsweep to organize documentation..."

# Get the project root (assuming hook is in .claude/hooks/)
PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

# Run qsweep from project root
cd "$PROJECT_ROOT"

# Check if qsweep.sh exists
if [[ -f "$PROJECT_ROOT/.claude/scripts/qsweep.sh" ]]; then
    # Run qsweep to move completed docs
    "$PROJECT_ROOT/.claude/scripts/qsweep.sh"

    # Show summary
    echo "✅ Documentation sweep complete!"
    echo "📁 Check ai-docs/completed/ for organized documentation"
else
    echo "⚠️  qsweep.sh not found in project root"
fi
