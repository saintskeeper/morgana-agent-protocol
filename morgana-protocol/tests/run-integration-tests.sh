#!/bin/bash
set -e

# Morgana Protocol Integration Test Runner
# This script runs the complete integration test suite with proper setup and reporting

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_DIR="$PROJECT_ROOT/tests/integration"

# Configuration
TIMEOUT=${TEST_TIMEOUT:-30m}
VERBOSE=${TEST_VERBOSE:-true}
COVERAGE=${TEST_COVERAGE:-false}
PARALLEL=${TEST_PARALLEL:-false}
RACE_DETECTOR=${TEST_RACE:-true}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check Go version
    if ! command -v go &> /dev/null; then
        error "Go is not installed or not in PATH"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    log "Go version: $GO_VERSION"
    
    # Check if we're in the project root
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        error "Not in project root or go.mod not found"
        exit 1
    fi
    
    # Check for required directories
    if [[ ! -d "$TEST_DIR" ]]; then
        error "Integration test directory not found: $TEST_DIR"
        exit 1
    fi
    
    success "Prerequisites checked"
}

# Set up test environment
setup_environment() {
    log "Setting up test environment..."
    
    # Ensure test dependencies are available
    cd "$PROJECT_ROOT"
    go mod download
    go mod tidy
    
    # Set environment variables for testing
    export MORGANA_DEBUG=${MORGANA_DEBUG:-false}
    export GORACE="halt_on_error=1"
    
    # Create temporary directory for test artifacts
    export TEST_ARTIFACTS_DIR=$(mktemp -d)
    log "Test artifacts directory: $TEST_ARTIFACTS_DIR"
    
    success "Environment set up"
}

# Run individual test suite
run_test_suite() {
    local suite_name=$1
    local test_pattern=$2
    local description=$3
    
    log "Running $description..."
    
    local test_cmd="go test"
    local test_args="-tags=integration -timeout=$TIMEOUT"
    
    if [[ "$VERBOSE" == "true" ]]; then
        test_args="$test_args -v"
    fi
    
    if [[ "$RACE_DETECTOR" == "true" ]]; then
        test_args="$test_args -race"
    fi
    
    if [[ "$PARALLEL" == "true" ]]; then
        test_args="$test_args -parallel=4"
    fi
    
    if [[ "$COVERAGE" == "true" ]]; then
        local coverage_file="$TEST_ARTIFACTS_DIR/coverage-$suite_name.out"
        test_args="$test_args -coverprofile=$coverage_file"
    fi
    
    if [[ -n "$test_pattern" ]]; then
        test_args="$test_args -run=$test_pattern"
    fi
    
    # Run the test
    cd "$TEST_DIR"
    if $test_cmd $test_args . 2>&1 | tee "$TEST_ARTIFACTS_DIR/test-$suite_name.log"; then
        success "$description completed"
        return 0
    else
        error "$description failed"
        return 1
    fi
}

# Generate coverage report
generate_coverage_report() {
    if [[ "$COVERAGE" != "true" ]]; then
        return 0
    fi
    
    log "Generating coverage report..."
    
    cd "$PROJECT_ROOT"
    
    # Combine coverage files
    echo "mode: atomic" > "$TEST_ARTIFACTS_DIR/integration-coverage.out"
    for coverage_file in "$TEST_ARTIFACTS_DIR"/coverage-*.out; do
        if [[ -f "$coverage_file" ]]; then
            tail -n +2 "$coverage_file" >> "$TEST_ARTIFACTS_DIR/integration-coverage.out"
        fi
    done
    
    # Generate HTML report
    go tool cover -html="$TEST_ARTIFACTS_DIR/integration-coverage.out" -o "$TEST_ARTIFACTS_DIR/integration-coverage.html"
    
    # Calculate coverage percentage
    local coverage_pct=$(go tool cover -func="$TEST_ARTIFACTS_DIR/integration-coverage.out" | grep total | awk '{print $3}')
    
    success "Coverage report generated: $coverage_pct"
    log "HTML report: $TEST_ARTIFACTS_DIR/integration-coverage.html"
}

# Clean up
cleanup() {
    log "Cleaning up..."
    
    # Kill any remaining background processes
    pkill -f "morgana-monitor" 2>/dev/null || true
    
    # Clean up temporary sockets
    rm -f /tmp/morgana*.sock 2>/dev/null || true
    
    # Copy important artifacts to project directory if they exist
    if [[ -f "$TEST_ARTIFACTS_DIR/integration-coverage.html" ]]; then
        cp "$TEST_ARTIFACTS_DIR/integration-coverage.html" "$PROJECT_ROOT/"
        log "Coverage report copied to project root"
    fi
    
    if [[ -f "$TEST_ARTIFACTS_DIR/integration-coverage.out" ]]; then
        cp "$TEST_ARTIFACTS_DIR/integration-coverage.out" "$PROJECT_ROOT/"
    fi
    
    # Clean up artifacts directory
    rm -rf "$TEST_ARTIFACTS_DIR" 2>/dev/null || true
}

# Signal handlers
trap cleanup EXIT
trap 'error "Test interrupted"; exit 130' INT TERM

# Main execution
main() {
    local suite_filter=""
    local run_all=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --suite)
                suite_filter="$2"
                run_all=false
                shift 2
                ;;
            --coverage)
                COVERAGE=true
                shift
                ;;
            --no-race)
                RACE_DETECTOR=false
                shift
                ;;
            --parallel)
                PARALLEL=true
                shift
                ;;
            --timeout)
                TIMEOUT="$2"
                shift 2
                ;;
            --help|-h)
                echo "Usage: $0 [options]"
                echo "Options:"
                echo "  --suite <name>    Run specific test suite (events|tui|monitor|commands|resources|suite)"
                echo "  --coverage        Generate coverage report"
                echo "  --no-race         Disable race detector"
                echo "  --parallel        Run tests in parallel"
                echo "  --timeout <dur>   Set test timeout (default: 30m)"
                echo "  --help            Show this help"
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    log "Starting Morgana Protocol Integration Tests"
    log "Configuration: Coverage=$COVERAGE, Race=$RACE_DETECTOR, Parallel=$PARALLEL, Timeout=$TIMEOUT"
    
    check_prerequisites
    setup_environment
    
    local failed_suites=()
    local total_suites=0
    
    # Define test suites
    if [[ "$run_all" == "true" || "$suite_filter" == "events" ]]; then
        ((total_suites++))
        if ! run_test_suite "events" "TestEventStream" "Event Stream Tests"; then
            failed_suites+=("events")
        fi
    fi
    
    if [[ "$run_all" == "true" || "$suite_filter" == "tui" ]]; then
        ((total_suites++))
        if ! run_test_suite "tui" "TestTUI" "TUI Integration Tests"; then
            failed_suites+=("tui")
        fi
    fi
    
    if [[ "$run_all" == "true" || "$suite_filter" == "monitor" ]]; then
        ((total_suites++))
        if ! run_test_suite "monitor" "TestMonitor" "Monitor Integration Tests"; then
            failed_suites+=("monitor")
        fi
    fi
    
    if [[ "$run_all" == "true" || "$suite_filter" == "commands" ]]; then
        ((total_suites++))
        if ! run_test_suite "commands" "TestCommand" "Command Polling Tests"; then
            failed_suites+=("commands")
        fi
    fi
    
    if [[ "$run_all" == "true" || "$suite_filter" == "resources" ]]; then
        ((total_suites++))
        if ! run_test_suite "resources" "TestResource" "Resource Leak Tests"; then
            failed_suites+=("resources")
        fi
    fi
    
    if [[ "$run_all" == "true" || "$suite_filter" == "suite" ]]; then
        ((total_suites++))
        if ! run_test_suite "suite" "TestMorganaIntegrationSuite" "Complete Integration Suite"; then
            failed_suites+=("suite")
        fi
    fi
    
    # Generate coverage report
    generate_coverage_report
    
    # Report results
    local passed_suites=$((total_suites - ${#failed_suites[@]}))
    
    log "Test Results Summary:"
    log "  Total Suites: $total_suites"
    log "  Passed: $passed_suites"
    log "  Failed: ${#failed_suites[@]}"
    
    if [[ ${#failed_suites[@]} -eq 0 ]]; then
        success "All integration tests passed!"
        exit 0
    else
        error "Failed test suites: ${failed_suites[*]}"
        error "Check individual test logs in: $TEST_ARTIFACTS_DIR"
        exit 1
    fi
}

# Run main function
main "$@"