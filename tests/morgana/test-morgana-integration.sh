#!/bin/bash
# Test script for Morgana Protocol Integration Testing Framework

echo "🧪 Morgana Protocol Integration Test Suite"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in the right directory
if [ ! -d "../../morgana-protocol" ]; then
    echo -e "${RED}❌ Error: morgana-protocol directory not found${NC}"
    echo "Please run this script from the ~/.claude/tests/morgana directory"
    exit 1
fi

cd ../../morgana-protocol

echo -e "${BLUE}📦 Installing dependencies...${NC}"
make deps

echo -e "\n${BLUE}🔨 Building Morgana...${NC}"
make build
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo -e "\n${BLUE}🧪 Running unit tests...${NC}"
make test
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Unit tests passed${NC}"
else
    echo -e "${RED}❌ Unit tests failed${NC}"
fi

echo -e "\n${BLUE}🧪 Running integration tests...${NC}"
make test-integration
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Integration tests passed${NC}"
else
    echo -e "${RED}❌ Integration tests failed${NC}"
fi

echo -e "\n${BLUE}📊 Running test coverage...${NC}"
go test -v -tags=integration -cover ./... | grep -E "(coverage:|PASS|FAIL)"

echo -e "\n${BLUE}🐍 Testing Python bridge...${NC}"
if [ -f "scripts/task_bridge.py" ]; then
    echo '{"agent_type":"test","prompt":"Hello from test"}' | python3 scripts/task_bridge.py
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Python bridge working${NC}"
    else
        echo -e "${RED}❌ Python bridge failed${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Python bridge script not found${NC}"
fi

echo -e "\n${BLUE}🚀 Testing Morgana CLI...${NC}"
if [ -f "dist/morgana" ]; then
    ./dist/morgana --version
    echo -e "${GREEN}✅ Morgana CLI accessible${NC}"
    
    # Test mock mode
    echo -e "\n${BLUE}Testing mock mode...${NC}"
    echo '[{"agent_type":"code-implementer","prompt":"test"}]' | ./dist/morgana --mock
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Mock mode working${NC}"
    else
        echo -e "${RED}❌ Mock mode failed${NC}"
    fi
else
    echo -e "${RED}❌ Morgana binary not found${NC}"
fi

echo -e "\n${BLUE}🔍 Test Summary${NC}"
echo "================"
echo "Integration tests cover:"
echo "  • Task client (Python bridge execution)"
echo "  • Adapter (agent orchestration)"
echo "  • Timeout handling"
echo "  • Error propagation"
echo "  • Concurrent execution"
echo "  • Mock mode for testing"

echo -e "\n${GREEN}✨ To run specific tests:${NC}"
echo "  go test -v -tags=integration ./pkg/task/..."
echo "  go test -v -tags=integration ./internal/adapter/..."
echo "  go test -v -run TestClientTimeout ./pkg/task/..."

echo -e "\n${GREEN}✨ To debug tests:${NC}"
echo "  export TEST_MODE=timeout"
echo "  go test -v -tags=integration -run TestAdapterTimeout ./internal/adapter"

cd - >/dev/null