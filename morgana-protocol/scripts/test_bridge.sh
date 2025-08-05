#!/bin/bash
# Test script for Python bridge

echo "Testing Python bridge..."

# Test successful execution
echo '{"agent_type": "test", "prompt": "Hello"}' | python3 scripts/task_bridge.py
echo ""

# Test with debug mode
echo "Testing with debug mode..."
MORGANA_DEBUG=true echo '{"agent_type": "test", "prompt": "Debug test"}' | python3 scripts/task_bridge.py
echo ""

# Test error handling
echo "Testing error handling..."
echo 'invalid json' | python3 scripts/task_bridge.py
echo ""

echo "Bridge tests complete"