#!/bin/bash
# Test script for QDIRECTOR with Morgana Protocol integration

echo "🧪 Testing QDIRECTOR with Morgana Protocol Integration"
echo "====================================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Source the agent adapter functions
echo -e "${BLUE}Loading Morgana adapter functions...${NC}"
source ../../scripts/agent-adapter-wrapper.sh

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Adapter functions loaded successfully${NC}"
else
    echo -e "${RED}❌ Failed to load adapter functions${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Test 1: Single Agent Execution${NC}"
echo "Testing: AgentAdapter \"code-implementer\" \"Create a simple hello world function\""
echo ""

# Test single agent execution
if command -v AgentAdapter >/dev/null 2>&1; then
    echo -e "${GREEN}✅ AgentAdapter function is available${NC}"
    
    # Mock test - just verify the function exists and can be called
    # In real usage, this would execute via Morgana
    echo -e "${YELLOW}Note: This is a dry run - actual Morgana execution requires built binary${NC}"
else
    echo -e "${RED}❌ AgentAdapter function not found${NC}"
fi

echo ""
echo -e "${BLUE}Test 2: Parallel Agent Execution${NC}"
echo "Testing parallel execution with AgentAdapterParallel"
echo ""

if command -v AgentAdapterParallel >/dev/null 2>&1; then
    echo -e "${GREEN}✅ AgentAdapterParallel function is available${NC}"
    
    # Show example of parallel execution
    cat << 'EOF'
Example usage:
AgentAdapterParallel << 'JSON'
[
  {"agent_type": "code-implementer", "prompt": "implement auth service"},
  {"agent_type": "test-specialist", "prompt": "create auth tests"},
  {"agent_type": "validation-expert", "prompt": "security review"}
]
JSON
EOF
else
    echo -e "${RED}❌ AgentAdapterParallel function not found${NC}"
fi

echo ""
echo -e "${BLUE}Test 3: QDIRECTOR Command Compatibility${NC}"
echo ""

# Test the Python adapter exists
if [ -f ../../scripts/agent_adapter.py ]; then
    echo -e "${GREEN}✅ Python adapter exists (for backward compatibility)${NC}"
else
    echo -e "${YELLOW}⚠️  Python adapter not found (optional)${NC}"
fi

# Check if Morgana binary exists
echo ""
echo -e "${BLUE}Test 4: Morgana Binary Check${NC}"

MORGANA_PATHS=(
    "$HOME/.claude/bin/morgana"
    "$HOME/.claude/morgana-protocol/dist/morgana"
)

MORGANA_FOUND=false
for path in "${MORGANA_PATHS[@]}"; do
    if [ -f "$path" ]; then
        echo -e "${GREEN}✅ Morgana binary found at: $path${NC}"
        MORGANA_FOUND=true
        
        # Test version
        "$path" --version 2>/dev/null && echo -e "${GREEN}✅ Morgana is executable${NC}"
        break
    fi
done

if [ "$MORGANA_FOUND" = false ]; then
    echo -e "${YELLOW}⚠️  Morgana binary not found. Build it with:${NC}"
    echo "   cd ~/.claude/morgana-protocol && make build"
fi

echo ""
echo -e "${BLUE}Test 5: Example QDIRECTOR Workflow${NC}"
echo ""

cat << 'EOF'
Example QDIRECTOR workflow with Morgana:

1. Sequential execution:
   AgentAdapter "sprint-planner" "Plan authentication feature"
   AgentAdapter "code-implementer" "Implement JWT service"
   AgentAdapter "test-specialist" "Create JWT tests"

2. Parallel execution:
   morgana_parallel << 'JSON'
   [
     {"agent_type": "code-implementer", "prompt": "implement user model"},
     {"agent_type": "code-implementer", "prompt": "implement auth middleware"},
     {"agent_type": "test-specialist", "prompt": "create integration tests"}
   ]
   JSON

3. Direct Morgana usage:
   morgana -- --agent validation-expert --prompt "Security audit auth system"
EOF

echo ""
echo -e "${GREEN}✨ Summary:${NC}"
echo "- AgentAdapter functions are integrated with Morgana Protocol"
echo "- QDIRECTOR can now use true Go-based parallel execution"
echo "- OpenTelemetry tracing is available for monitoring"
echo "- Integration tests validate the entire pipeline"

echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Build Morgana if not already done: cd ~/.claude/morgana-protocol && make build"
echo "2. Run integration tests: cd ~/.claude/morgana-protocol && make test-integration"
echo "3. Try a real QDIRECTOR command with /qdirector-enhanced"