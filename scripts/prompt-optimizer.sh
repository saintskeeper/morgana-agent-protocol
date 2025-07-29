#!/bin/bash

# Prompt Optimizer - Converts verbose prompts to structured format
# Usage: echo "your verbose prompt" | prompt-optimizer.sh

prompt_type=""
verbose_prompt=""

# Read input
if [ -t 0 ]; then
    # Interactive mode
    echo "Enter your prompt (Ctrl+D when done):"
    verbose_prompt=$(cat)
else
    # Pipe mode
    verbose_prompt=$(cat)
fi

# Detect prompt type based on keywords
if echo "$verbose_prompt" | grep -qiE "analyze|review|check|validate|audit"; then
    prompt_type="analyze"
elif echo "$verbose_prompt" | grep -qiE "implement|create|build|add|feature"; then
    prompt_type="implement"
elif echo "$verbose_prompt" | grep -qiE "test|tests|coverage|unit|integration"; then
    prompt_type="test"
elif echo "$verbose_prompt" | grep -qiE "plan|sprint|tasks|breakdown"; then
    prompt_type="plan"
else
    prompt_type="simple"
fi

# Convert based on type
case $prompt_type in
    "analyze")
        echo "Converted to Analysis Pattern:"
        echo "---"
        echo "Analyze [TARGET] for:"
        echo "- Security: [specific checks]"
        echo "- Performance: [metrics]"
        echo "- Quality: [standards]"
        echo ""
        echo "Report findings as:"
        echo "SEVERITY: issue at location"
        ;;
    
    "implement")
        echo "Converted to Implementation Pattern:"
        echo "---"
        echo "Implement: [feature]"
        echo "Requirements:"
        echo "- [requirement 1]"
        echo "- [requirement 2]"
        echo ""
        echo "Constraints:"
        echo "- Follow existing patterns"
        echo "- Match project style"
        ;;
    
    "test")
        echo "Converted to Test Pattern:"
        echo "---"
        echo "Generate tests for: [component]"
        echo "Coverage: Happy path, edge cases, errors"
        echo "Framework: [detected]"
        echo "Pattern: AAA"
        ;;
    
    "plan")
        echo "Converted to Planning Pattern:"
        echo "---"
        echo "Task: Plan sprint for [feature]"
        echo "Output: QDIRECTOR YAML format"
        echo "Focus: Dependencies and exit criteria"
        echo "Constraints: 2-4 hour tasks"
        ;;
    
    *)
        echo "Converted to Simple Pattern:"
        echo "---"
        echo "Task: [action]"
        echo "Input: [data]"
        echo "Output: [expected format]"
        ;;
esac

echo ""
echo "Original prompt length: $(echo "$verbose_prompt" | wc -c) chars"
echo "Structured format: ~50-100 chars (est. 60-80% reduction)"