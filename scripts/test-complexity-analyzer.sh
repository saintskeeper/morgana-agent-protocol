#!/bin/bash

# Test script for code complexity analyzer

echo "=== Code Complexity Analyzer Test Suite ==="
echo ""

# Test simple tasks
echo "1. Testing SIMPLE tasks:"
echo "   Task: 'Create a utility function to format dates'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Create a utility function to format dates")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Create a utility function to format dates" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

echo "   Task: 'Add a button component'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Add a button component")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Add a button component" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

# Test moderate tasks
echo "2. Testing MODERATE tasks:"
echo "   Task: 'Implement REST API with authentication'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Implement REST API with authentication")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Implement REST API with authentication" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

echo "   Task: 'Create a caching service with Redis'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Create a caching service with Redis")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Create a caching service with Redis" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

# Test complex tasks
echo "3. Testing COMPLEX tasks:"
echo "   Task: 'Design distributed caching system with failover'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Design distributed caching system with failover")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Design distributed caching system with failover" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

echo "   Task: 'Refactor entire payment processing architecture for microservices'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Refactor entire payment processing architecture for microservices")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Refactor entire payment processing architecture for microservices" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

# Test with token-efficient mode disabled
echo "4. Testing with token-efficient mode DISABLED:"
echo "   Task: 'Create a simple validation function' (token-efficient=false)"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Create a simple validation function")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Create a simple validation function" false)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

# Test edge cases
echo "5. Testing edge cases:"
echo "   Task: 'Implement concurrent distributed blockchain consensus algorithm'"
complexity=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh analyze "Implement concurrent distributed blockchain consensus algorithm")
model=$(/Users/walterday/.claude/scripts/code-complexity-analyzer.sh recommend "Implement concurrent distributed blockchain consensus algorithm" true)
echo "   Complexity: $complexity"
echo "   Recommended model: $model"
echo ""

echo "=== Test Complete ==="