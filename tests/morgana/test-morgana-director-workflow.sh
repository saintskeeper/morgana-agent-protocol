#!/bin/bash
# Test script to verify /morgana-director workflow integration

echo "🧪 Testing Morgana-Director Workflow Integration"
echo "=============================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Test 1: Check if Morgana binary is available${NC}"
if command -v morgana >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Morgana binary found in PATH${NC}"
    morgana --version
elif [ -f "$HOME/.claude/bin/morgana" ]; then
    echo -e "${GREEN}✅ Morgana binary found in ~/.claude/bin${NC}"
    "$HOME/.claude/bin/morgana" --version
elif [ -f "$HOME/.claude/morgana-protocol/dist/morgana" ]; then
    echo -e "${GREEN}✅ Morgana binary found in morgana-protocol/dist${NC}"
    "$HOME/.claude/morgana-protocol/dist/morgana" --version
else
    echo -e "${RED}❌ Morgana binary not found${NC}"
    echo "Build it with: cd ~/.claude/morgana-protocol && make build"
fi

echo ""
echo -e "${BLUE}Test 2: Check AgentAdapter functions${NC}"
if source ../../scripts/morgana-adapter.sh 2>/dev/null; then
    echo -e "${GREEN}✅ Morgana adapter functions loaded${NC}"
    
    if command -v AgentAdapter >/dev/null 2>&1; then
        echo -e "${GREEN}✅ AgentAdapter function available${NC}"
    else
        echo -e "${RED}❌ AgentAdapter function not found${NC}"
    fi
    
    if command -v AgentAdapterParallel >/dev/null 2>&1; then
        echo -e "${GREEN}✅ AgentAdapterParallel function available${NC}"
    else
        echo -e "${RED}❌ AgentAdapterParallel function not found${NC}"
    fi
else
    echo -e "${RED}❌ Failed to load morgana-adapter.sh${NC}"
fi

echo ""
echo -e "${BLUE}Test 3: Check morgana-director command exists${NC}"
if [ -f ../../commands/morgana-director.md ]; then
    echo -e "${GREEN}✅ morgana-director.md exists${NC}"
    
    # Check if it contains the key integration points
    if grep -q "AgentAdapter" ../../commands/morgana-director.md; then
        echo -e "${GREEN}✅ Contains AgentAdapter integration${NC}"
    else
        echo -e "${YELLOW}⚠️  Missing AgentAdapter integration${NC}"
    fi
    
    if grep -q "Morgana Protocol" ../../commands/morgana-director.md; then
        echo -e "${GREEN}✅ Contains Morgana Protocol references${NC}"
    else
        echo -e "${YELLOW}⚠️  Missing Morgana Protocol references${NC}"
    fi
else
    echo -e "${RED}❌ morgana-director.md not found${NC}"
fi

echo ""
echo -e "${BLUE}Test 4: Verify workflow integration${NC}"

echo -e "${BLUE}Expected workflow:${NC}"
cat << 'EOF'
1. /morgana-plan -> Creates sprint plan
2. /morgana-validate -> Validates against codebase  
3. /morgana-director -> Orchestrates execution using:
   - AgentAdapter("sprint-planner", "task")
   - AgentAdapterParallel([multiple tasks])
   - Direct morgana binary calls
EOF

echo ""
echo -e "${BLUE}Test 5: Mock workflow test${NC}"
echo "Testing parallel agent call format..."

# Test the JSON format that would be sent to Morgana
test_json='[
  {"agent_type": "code-implementer", "prompt": "implement hello world function"},
  {"agent_type": "test-specialist", "prompt": "create tests for hello world"}
]'

echo "JSON format test:"
echo "$test_json" | jq . 2>/dev/null && echo -e "${GREEN}✅ Valid JSON format${NC}" || echo -e "${RED}❌ Invalid JSON format${NC}"

echo ""
echo -e "${BLUE}Test 6: Check integration test results${NC}"
if [ -f ../../morgana-protocol/Makefile ]; then
    cd ../../morgana-protocol
    echo "Running quick integration test check..."
    
    # Check if integration tests exist and can run
    if make test-integration >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Integration tests pass${NC}"
    else
        echo -e "${YELLOW}⚠️  Integration tests failed or not built${NC}"
        echo "Try: cd ~/.claude/morgana-protocol && make build && make test-integration"
    fi
    cd - >/dev/null
else
    echo -e "${YELLOW}⚠️  Morgana protocol source not found${NC}"
fi

echo ""
echo -e "${GREEN}✨ How to test the actual workflow:${NC}"
cat << 'EOF'
1. Try a simple morgana-director command:
   /morgana-director implement a simple calculator function

2. Check the output for:
   - [🤖 Executing agent: ...] messages
   - Parallel execution indicators  
   - Morgana Protocol usage
   - AgentAdapter function calls

3. Look for structured output with color markers:
   - [✓] Success indicators
   - [!] Warning messages
   - [→] Next actions

4. Verify multiple agents are spawned for complex tasks
EOF

echo ""
echo -e "${BLUE}Debug Commands:${NC}"
echo "If issues found:"
echo "- Check logs: tail -f ~/.claude/logs/*"  
echo "- Manual test: morgana -- --agent code-implementer --prompt 'test'"
echo "- Parallel test: echo '[{\"agent_type\":\"code-implementer\",\"prompt\":\"test\"}]' | morgana --parallel"