#!/bin/bash

# Prompt Helper Functions for Token-Efficient Prompts
# Source this file in your shell: source ~/.claude/scripts/prompt-helpers.sh

# Function to analyze code
claude-analyze() {
    local file=$1
    local focus=${2:-"security, performance, quality"}
    
    claude "Analyze $file for:
- Security: injection risks, auth issues
- Performance: complexity, bottlenecks  
- Quality: code smells, maintainability

Report findings as:
CRITICAL: [issue] at line X
HIGH: [issue] at line Y
MEDIUM: [issue] at line Z"
}

# Function to implement features
claude-implement() {
    local feature=$1
    shift
    local requirements="$@"
    
    claude "Implement: $feature
Requirements:
$requirements

Constraints:
- Follow existing patterns
- No comments unless asked
- Match project style
Output: Production-ready code"
}

# Function to generate tests
claude-test() {
    local component=$1
    
    claude "Generate tests for: $component
Coverage: Happy path, edge cases, errors
Framework: [auto-detect from project]
Pattern: AAA (Arrange, Act, Assert)
Focus: Behavior not implementation"
}

# Function for quick tasks
claude-quick() {
    local task=$1
    local input=$2
    
    claude "Task: $task
Input: $input
Output: Direct result, no explanation"
}

# Aliases for common operations
alias ca='claude-analyze'
alias ci='claude-implement'
alias ct='claude-test'
alias cq='claude-quick'

echo "Prompt helpers loaded! Available commands:"
echo "  ca <file>           - Analyze code file"
echo "  ci <feature> <reqs> - Implement feature"
echo "  ct <component>      - Generate tests"
echo "  cq <task> <input>   - Quick task execution"