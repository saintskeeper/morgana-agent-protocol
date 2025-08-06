#!/bin/bash
set -e

# Morgana Protocol Local Development Setup Script
# This script sets up everything needed for local Morgana Protocol development

BINARY_NAME="morgana"
INSTALL_DIR="$HOME/.local/bin"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_DIR="$HOME/.claude"
AGENTS_DIR="$CONFIG_DIR/agents"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS")
            echo -e "${GREEN}✅ SUCCESS${NC}: $message"
            ;;
        "ERROR")
            echo -e "${RED}❌ ERROR${NC}: $message"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠️  WARNING${NC}: $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  INFO${NC}: $message"
            ;;
    esac
}

echo -e "${BLUE}🧙‍♂️ Morgana Protocol Local Setup${NC}"
echo "===================================="
print_status "INFO" "Setting up complete Morgana Protocol local development environment"
print_status "INFO" "Repository: $REPO_ROOT"

# Navigate to repository root
cd "$REPO_ROOT"

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/morgana" ]; then
    print_status "ERROR" "Not in a Morgana Protocol repository root"
    exit 1
fi

# 1. DEPENDENCY CHECKS
echo ""
echo -e "${BLUE}📋 Checking Dependencies${NC}"
echo "=========================="

# Check Go installation
if ! command -v go &> /dev/null; then
    print_status "ERROR" "Go is not installed or not in PATH"
    print_status "INFO" "Install Go from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_status "SUCCESS" "Go version: $GO_VERSION"

# Check Go version (require 1.21+)
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
    print_status "WARNING" "Go 1.21+ recommended, you have $GO_VERSION"
fi

# Check Git
if ! command -v git &> /dev/null; then
    print_status "WARNING" "Git not found - version info will be limited"
else
    print_status "SUCCESS" "Git available"
fi

# Check terminal support for TUI
if [[ -t 1 ]] && [[ -n "$TERM" ]] && [[ "$TERM" != "dumb" ]]; then
    COLORS=$(tput colors 2>/dev/null || echo "8")
    print_status "SUCCESS" "Terminal supports TUI (TERM: $TERM, Colors: $COLORS)"
else
    print_status "WARNING" "Terminal may not fully support TUI features"
fi

# 2. DIRECTORY SETUP
echo ""
echo -e "${BLUE}📁 Setting Up Directories${NC}"
echo "============================"

# Create necessary directories
directories=(
    "$CONFIG_DIR"
    "$AGENTS_DIR"
    "$INSTALL_DIR"
    "$CONFIG_DIR/logs"
    "$CONFIG_DIR/cache"
)

for dir in "${directories[@]}"; do
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        print_status "SUCCESS" "Created directory: $dir"
    else
        print_status "INFO" "Directory exists: $dir"
    fi
done

# 3. AGENT SETUP
echo ""
echo -e "${BLUE}🤖 Setting Up Agents${NC}"
echo "====================="

# Create sample agents if they don't exist
agents=(
    "code-implementer"
    "sprint-planner" 
    "test-specialist"
    "validation-expert"
)

for agent in "${agents[@]}"; do
    agent_file="$AGENTS_DIR/${agent}.md"
    if [ ! -f "$agent_file" ]; then
        case $agent in
            "code-implementer")
                cat > "$agent_file" << EOF
# Code Implementer Agent

You are a senior software developer specializing in implementing robust, maintainable code solutions.

## Your Role
- Implement features based on requirements
- Write clean, well-structured code
- Follow best practices and coding standards
- Consider performance and scalability
- Provide clear code documentation

## Guidelines
- Use appropriate design patterns
- Implement proper error handling
- Write testable code
- Consider security implications
- Optimize for readability and maintainability

## Output Format
Provide complete, working code implementations with:
- Clear file structure
- Proper imports/dependencies
- Implementation details
- Basic usage examples
EOF
                ;;
            "sprint-planner")
                cat > "$agent_file" << EOF
# Sprint Planner Agent

You are a technical project manager specialized in agile sprint planning and task decomposition.

## Your Role
- Break down large features into manageable tasks
- Estimate task complexity and effort
- Identify dependencies and risks
- Plan sprint iterations
- Prioritize work based on business value

## Guidelines
- Create clear, actionable tasks
- Provide realistic time estimates
- Identify technical dependencies
- Consider team capacity and skills
- Balance feature work with technical debt

## Output Format
Provide structured sprint plans with:
- Task breakdowns
- Priority assignments
- Time estimates
- Dependency mapping
- Risk assessments
EOF
                ;;
            "test-specialist")
                cat > "$agent_file" << EOF
# Test Specialist Agent

You are a quality assurance engineer specializing in comprehensive software testing.

## Your Role
- Design comprehensive test strategies
- Write unit, integration, and end-to-end tests
- Identify edge cases and failure scenarios
- Ensure proper test coverage
- Automate testing processes

## Guidelines
- Create thorough test coverage
- Test both happy paths and error cases
- Use appropriate testing frameworks
- Write maintainable test code
- Document test scenarios and rationale

## Output Format
Provide complete test implementations with:
- Test framework setup
- Comprehensive test cases
- Mock and fixture data
- Test documentation
- Coverage reports
EOF
                ;;
            "validation-expert")
                cat > "$agent_file" << EOF
# Validation Expert Agent

You are a senior technical reviewer specializing in code quality, security, and architecture validation.

## Your Role
- Review code for quality and standards compliance
- Identify security vulnerabilities
- Validate architectural decisions
- Ensure performance considerations
- Verify best practices implementation

## Guidelines
- Conduct thorough code reviews
- Check for security vulnerabilities
- Validate design patterns usage
- Assess performance implications
- Ensure maintainability and scalability

## Output Format
Provide detailed validation reports with:
- Code quality assessment
- Security analysis
- Performance considerations
- Architectural feedback
- Improvement recommendations
EOF
                ;;
        esac
        print_status "SUCCESS" "Created agent: $agent_file"
    else
        print_status "INFO" "Agent exists: $agent"
    fi
done

# 4. CONFIGURATION SETUP
echo ""
echo -e "${BLUE}⚙️  Setting Up Configuration${NC}"
echo "==============================="

# Copy default config if it doesn't exist in ~/.claude/
target_config="$CONFIG_DIR/morgana.yaml"
source_config="$REPO_ROOT/morgana.yaml"

if [ ! -f "$target_config" ]; then
    if [ -f "$source_config" ]; then
        cp "$source_config" "$target_config"
        print_status "SUCCESS" "Copied configuration to $target_config"
    else
        # Create a basic config
        cat > "$target_config" << EOF
# Morgana Protocol Configuration - Local Development Setup

# Agent configuration
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m
  timeouts:
    code-implementer: 5m
    test-specialist: 3m
    validation-expert: 2m
    sprint-planner: 1m

# Execution configuration
execution:
  max_concurrency: 3
  default_mode: sequential

# Telemetry configuration  
telemetry:
  enabled: true
  exporter: stdout
  service_name: morgana-local
  environment: development

# Task client configuration
task_client:
  mock_mode: false
  timeout: 5m

# TUI configuration
tui:
  enabled: true
  performance:
    refresh_rate: 16ms
    max_log_lines: 5000
  visual:
    show_debug_info: true
    show_timestamps: true
  features:
    enable_filtering: true
    enable_search: true
    enable_export: false
  events:
    buffer_size: 1000
EOF
        print_status "SUCCESS" "Created default configuration: $target_config"
    fi
else
    print_status "INFO" "Configuration exists: $target_config"
fi

# 5. BUILD AND INSTALL
echo ""
echo -e "${BLUE}🔨 Building and Installing${NC}"
echo "============================="

# Clean previous builds
print_status "INFO" "Cleaning previous builds..."
rm -f morgana morgana-*

# Update dependencies
print_status "INFO" "Updating Go dependencies..."
if ! go mod tidy; then
    print_status "ERROR" "Failed to update Go modules"
    exit 1
fi

# Build version info
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION="v1.0.0-local-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")"

LDFLAGS="-w -s -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

print_status "INFO" "Building Morgana CLI..."
print_status "INFO" "Version: $VERSION"

# Build with optimizations
if ! CGO_ENABLED=0 go build \
    -ldflags "$LDFLAGS" \
    -trimpath \
    -tags "netgo osusergo static_build" \
    -o "$BINARY_NAME" \
    ./cmd/morgana; then
    print_status "ERROR" "Failed to build Morgana Protocol CLI"
    exit 1
fi

print_status "SUCCESS" "Build completed successfully"

# Test the binary
print_status "INFO" "Testing binary..."
if ! ./"$BINARY_NAME" --version; then
    print_status "ERROR" "Binary test failed"
    exit 1
fi

print_status "SUCCESS" "Binary test passed"

# Backup existing binary if it exists
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    print_status "INFO" "Backing up existing binary..."
    cp "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME.backup.$(date +%s)"
fi

# Install the binary
print_status "INFO" "Installing to $INSTALL_DIR..."
cp "$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

print_status "SUCCESS" "Morgana Protocol CLI installed to $INSTALL_DIR/$BINARY_NAME"

# 6. PATH SETUP
echo ""
echo -e "${BLUE}🛤️  Setting Up PATH${NC}"
echo "==================="

# Check if PATH includes install directory
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_status "WARNING" "$INSTALL_DIR is not in your PATH"
    
    # Try to add to shell profile automatically
    SHELL_PROFILE=""
    SHELL_NAME=$(basename "$SHELL")
    
    case "$SHELL_NAME" in
        "zsh")
            if [ -f "$HOME/.zshrc" ]; then
                SHELL_PROFILE="$HOME/.zshrc"
            fi
            ;;
        "bash")
            if [ -f "$HOME/.bashrc" ]; then
                SHELL_PROFILE="$HOME/.bashrc"
            elif [ -f "$HOME/.bash_profile" ]; then
                SHELL_PROFILE="$HOME/.bash_profile"
            fi
            ;;
    esac
    
    if [ -n "$SHELL_PROFILE" ]; then
        echo ""
        echo -e "${YELLOW}Add $INSTALL_DIR to your PATH automatically? [y/N]${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            # Check if already in profile
            if ! grep -q "$INSTALL_DIR" "$SHELL_PROFILE" 2>/dev/null; then
                echo "" >> "$SHELL_PROFILE"
                echo "# Morgana Protocol CLI" >> "$SHELL_PROFILE"
                echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_PROFILE"
                print_status "SUCCESS" "Added $INSTALL_DIR to $SHELL_PROFILE"
                print_status "INFO" "Restart your shell or run: source $SHELL_PROFILE"
            else
                print_status "INFO" "PATH already configured in $SHELL_PROFILE"
            fi
        fi
    fi
else
    print_status "SUCCESS" "$INSTALL_DIR is already in PATH"
fi

# 7. VERIFICATION
echo ""
echo -e "${BLUE}✅ Verification${NC}"
echo "================="

# Test installation
if command -v morgana &> /dev/null; then
    INSTALLED_VERSION=$(morgana --version 2>/dev/null | head -n1 || echo "unknown")
    print_status "SUCCESS" "Installation verified: $INSTALLED_VERSION"
else
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        INSTALLED_VERSION=$("$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null | head -n1 || echo "unknown")
        print_status "SUCCESS" "Binary verified: $INSTALLED_VERSION"
        print_status "WARNING" "Binary not in PATH - restart shell or source profile"
    else
        print_status "ERROR" "Installation verification failed"
    fi
fi

# Test configuration
if morgana --config "$target_config" --version &>/dev/null || "$INSTALL_DIR/morgana" --config "$target_config" --version &>/dev/null; then
    print_status "SUCCESS" "Configuration file validated"
else
    print_status "WARNING" "Configuration file may have issues"
fi

# Test TUI support
if [[ -t 1 ]] && [[ -n "$TERM" ]] && [[ "$TERM" != "dumb" ]]; then
    print_status "SUCCESS" "TUI support verified"
else
    print_status "WARNING" "TUI may not work in this terminal"
fi

# Clean up build artifacts
rm -f "$BINARY_NAME"

# 8. COMPLETION
echo ""
echo -e "${GREEN}🎉 Local Setup Complete!${NC}"
echo "=========================="
echo ""
print_status "SUCCESS" "Morgana Protocol local development environment is ready!"
echo ""

echo -e "${BLUE}📁 Setup Summary:${NC}"
echo "  Binary:        $INSTALL_DIR/$BINARY_NAME"
echo "  Configuration: $target_config"
echo "  Agents:        $AGENTS_DIR/"
echo "  Logs:          $CONFIG_DIR/logs/"
echo ""

echo -e "${BLUE}🚀 Quick Start Commands:${NC}"
echo "  morgana --version                                    # Check installation"
echo "  morgana --config ~/.claude/morgana.yaml --tui       # Start with TUI"
echo "  morgana --tui --tui-mode dev                        # Development mode"
echo "  morgana --tui --agent code-implementer --prompt \"Hello\" # Test agent"
echo ""

echo -e "${BLUE}🧪 Test Your Setup:${NC}"
echo "  # Run demo tasks with TUI"
echo "  morgana --config ~/.claude/morgana.yaml --tui"
echo ""
echo "  # Test specific agent"
echo "  morgana --tui --agent code-implementer --prompt \"Create a hello world function in Go\""
echo ""

echo -e "${BLUE}📚 Next Steps:${NC}"
echo "  1. Restart your shell or run: source ~/.$(basename $SHELL)rc"  
echo "  2. Read TUI_USER_GUIDE.md for keyboard shortcuts and features"
echo "  3. Customize agents in ~/.claude/agents/"
echo "  4. Edit ~/.claude/morgana.yaml for your preferences"
echo ""

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠️  Don't forget to restart your shell or run:${NC}"
    echo "  source ~/.$(basename $SHELL)rc"
    echo ""
fi

print_status "SUCCESS" "Happy orchestrating with Morgana Protocol! 🧙‍♂️"