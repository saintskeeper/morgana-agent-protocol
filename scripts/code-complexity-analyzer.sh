#!/bin/bash

# Code Complexity Analyzer for Model Selection
# Analyzes task complexity to determine optimal model for code generation

analyze_task_complexity() {
    local task_description="$1"
    local context_files="$2"
    
    # Initialize complexity score
    local complexity_score=0
    
    # Convert to lowercase for analysis
    local task_lower=$(echo "$task_description" | tr '[:upper:]' '[:lower:]')
    
    # High complexity indicators (+3 points each)
    local high_complexity_patterns=(
        "architect"
        "design system"
        "refactor.*entire"
        "migrate"
        "implement.*framework"
        "distributed"
        "concurrent"
        "parallel"
        "async.*complex"
        "security.*critical"
        "performance.*critical"
        "algorithm.*complex"
        "data structure.*custom"
        "state machine"
        "parser"
        "compiler"
        "interpreter"
        "real-time"
        "blockchain"
        "cryptograph"
        "machine learning"
        "neural network"
    )
    
    # Medium complexity indicators (+2 points each)
    local medium_complexity_patterns=(
        "api"
        "service"
        "integration"
        "authentication"
        "authorization"
        "database.*schema"
        "cache"
        "queue"
        "websocket"
        "graphql"
        "rest.*endpoint"
        "middleware"
        "validation.*complex"
        "business logic"
        "workflow"
        "state management"
        "error handling.*comprehensive"
    )
    
    # Low complexity indicators (+1 point each)
    local low_complexity_patterns=(
        "component"
        "function"
        "utility"
        "helper"
        "convert"
        "format"
        "validate.*simple"
        "crud"
        "form"
        "button"
        "modal"
        "tooltip"
        "config"
        "constant"
        "type.*definition"
        "interface"
        "model"
        "schema.*simple"
    )
    
    # Check high complexity patterns
    for pattern in "${high_complexity_patterns[@]}"; do
        if [[ "$task_lower" =~ $pattern ]]; then
            ((complexity_score+=3))
        fi
    done
    
    # Check medium complexity patterns
    for pattern in "${medium_complexity_patterns[@]}"; do
        if [[ "$task_lower" =~ $pattern ]]; then
            ((complexity_score+=2))
        fi
    done
    
    # Check low complexity patterns
    for pattern in "${low_complexity_patterns[@]}"; do
        if [[ "$task_lower" =~ $pattern ]]; then
            ((complexity_score+=1))
        fi
    done
    
    # Analyze context files if provided
    if [ -n "$context_files" ]; then
        local file_count=$(echo "$context_files" | wc -w)
        if [ "$file_count" -gt 10 ]; then
            ((complexity_score+=3))  # Many files = higher complexity
        elif [ "$file_count" -gt 5 ]; then
            ((complexity_score+=2))
        elif [ "$file_count" -gt 2 ]; then
            ((complexity_score+=1))
        fi
    fi
    
    # Check for specific keywords that indicate complexity
    if [[ "$task_lower" =~ (optimize|performance|scale|distributed|concurrent) ]]; then
        ((complexity_score+=2))
    fi
    
    # Check for security-related tasks
    if [[ "$task_lower" =~ (security|encrypt|auth|permission|access.*control) ]]; then
        ((complexity_score+=2))
    fi
    
    # Determine complexity level based on score
    if [ "$complexity_score" -ge 8 ]; then
        echo "complex"
    elif [ "$complexity_score" -ge 4 ]; then
        echo "moderate"
    else
        echo "simple"
    fi
}

# Function to recommend model based on complexity
recommend_model() {
    local complexity="$1"
    local token_efficient_enabled="$2"
    
    case "$complexity" in
        "complex")
            echo "claude-4-opus"  # Use Claude 4 for complex tasks
            [ -n "$VERBOSE" ] && echo "Reason: Complex architectural or algorithmic task requiring deep reasoning" >&2
            ;;
        "moderate")
            if [ "$token_efficient_enabled" = "true" ]; then
                echo "claude-4-sonnet"  # Use Claude 4 Sonnet for moderate complexity
                [ -n "$VERBOSE" ] && echo "Reason: Moderate complexity with good balance of capability and efficiency" >&2
            else
                echo "gpt-4.1"  # Fallback to GPT-4.1
                [ -n "$VERBOSE" ] && echo "Reason: Moderate complexity requiring solid reasoning" >&2
            fi
            ;;
        "simple")
            if [ "$token_efficient_enabled" = "true" ]; then
                echo "claude-3-7-sonnet-20250219"  # Use Claude 3.7 for simple tasks
                [ -n "$VERBOSE" ] && echo "Reason: Simple task suitable for token-efficient generation" >&2
            else
                echo "gemini-2.5-flash"  # Fast model for simple tasks
                [ -n "$VERBOSE" ] && echo "Reason: Simple task suitable for fast generation" >&2
            fi
            ;;
    esac
}

# Function to analyze code structure complexity
analyze_code_structure() {
    local file_path="$1"
    local complexity_indicators=0
    
    if [ -f "$file_path" ]; then
        # Count classes/interfaces
        local class_count=$(grep -E "^(export\s+)?(class|interface)\s+" "$file_path" 2>/dev/null | wc -l)
        if [ "$class_count" -gt 5 ]; then
            ((complexity_indicators+=2))
        elif [ "$class_count" -gt 2 ]; then
            ((complexity_indicators+=1))
        fi
        
        # Count functions
        local function_count=$(grep -E "^(export\s+)?(async\s+)?function\s+|^\s*(async\s+)?\w+\s*\(" "$file_path" 2>/dev/null | wc -l)
        if [ "$function_count" -gt 10 ]; then
            ((complexity_indicators+=2))
        elif [ "$function_count" -gt 5 ]; then
            ((complexity_indicators+=1))
        fi
        
        # Check for complex patterns
        if grep -qE "(Promise\.all|async.*await|Observable|Subject|Stream)" "$file_path" 2>/dev/null; then
            ((complexity_indicators+=1))
        fi
        
        # Check for design patterns
        if grep -qE "(Factory|Singleton|Observer|Strategy|Decorator)" "$file_path" 2>/dev/null; then
            ((complexity_indicators+=2))
        fi
    fi
    
    echo "$complexity_indicators"
}

# Main execution
main() {
    local action="$1"
    
    case "$action" in
        "analyze")
            local task_description="$2"
            local context_files="$3"
            analyze_task_complexity "$task_description" "$context_files"
            ;;
        "recommend")
            local task_description="$2"
            local token_efficient_enabled="${3:-false}"
            local complexity=$(analyze_task_complexity "$task_description" "")
            recommend_model "$complexity" "$token_efficient_enabled"
            ;;
        "analyze-file")
            local file_path="$2"
            local score=$(analyze_code_structure "$file_path")
            if [ "$score" -ge 5 ]; then
                echo "complex"
            elif [ "$score" -ge 2 ]; then
                echo "moderate"
            else
                echo "simple"
            fi
            ;;
        *)
            echo "Usage: $0 {analyze|recommend|analyze-file} <args>"
            echo ""
            echo "Actions:"
            echo "  analyze <task_description> [context_files]"
            echo "    Analyzes task complexity and returns: simple|moderate|complex"
            echo ""
            echo "  recommend <task_description> [token_efficient_enabled]"
            echo "    Recommends optimal model based on task complexity"
            echo ""
            echo "  analyze-file <file_path>"
            echo "    Analyzes existing code file complexity"
            echo ""
            echo "Examples:"
            echo "  $0 analyze \"implement user authentication\""
            echo "  $0 recommend \"refactor entire payment system\" true"
            echo "  $0 analyze-file src/services/auth.service.ts"
            exit 1
            ;;
    esac
}

# Execute main function with all arguments
main "$@"