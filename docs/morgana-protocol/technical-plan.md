# Morgana Protocol: Technical Implementation Plan

## Enhanced Sprint Plan for Go-Based Agent Adapter

### Sprint Overview

**Goal**: Implement AgentAdapter in Go for cross-platform deployment
**Repository**: https://github.com/saintskeeper/claude-code-configs
**Priority**: P0-Critical **Estimated Effort**: 6-10 hours

### 1. Project Structure

```
claude-code-configs/
└── morgana-protocol/
    ├── cmd/
    │   └── morgana/
    │       └── main.go              # Entry point
    ├── internal/
    │   ├── adapter/
    │   │   ├── adapter.go           # Core adapter logic
    │   │   └── adapter_test.go
    │   ├── orchestrator/
    │   │   ├── parallel.go          # Goroutine orchestration
    │   │   └── parallel_test.go
    │   └── prompt/
    │       ├── loader.go            # Agent prompt loader
    │       └── loader_test.go
    ├── pkg/
    │   └── task/
    │       └── client.go            # Task tool interface
    ├── scripts/
    │   ├── build.sh                 # Multi-platform build
    │   └── install.sh               # Installation script
    ├── go.mod
    ├── go.sum
    ├── Makefile                     # Build automation
    └── README.md
```

### 2. Task Breakdown

#### TASK_001: Core Adapter Implementation

**Priority**: P0-Critical **Dependencies**: None **Complexity**: Medium

**Implementation** (`internal/adapter/adapter.go`):

```go
package adapter

import (
    "fmt"
    "os/exec"
    "encoding/json"
)

type Task struct {
    AgentType string                 `json:"agent_type"`
    Prompt    string                 `json:"prompt"`
    Options   map[string]interface{} `json:"options,omitempty"`
}

type Result struct {
    Output string `json:"output"`
    Error  string `json:"error,omitempty"`
}

type Adapter struct {
    promptLoader *PromptLoader
    taskClient   *TaskClient
    logger       *Logger
}

func (a *Adapter) Execute(task Task) (Result, error) {
    // Load agent prompt
    agentPrompt, err := a.promptLoader.Load(task.AgentType)
    if err != nil {
        return Result{}, fmt.Errorf("loading agent prompt: %w", err)
    }

    // Combine prompts
    fullPrompt := fmt.Sprintf("%s\n\nTask: %s", agentPrompt, task.Prompt)

    // Execute via Task tool
    return a.taskClient.Run("general-purpose", fullPrompt, task.Options)
}
```

#### TASK_002: Parallel Orchestration

**Priority**: P0-Critical **Dependencies**: TASK_001 **Complexity**: Medium

**Implementation** (`internal/orchestrator/parallel.go`):

```go
package orchestrator

import (
    "sync"
    "github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/adapter"
)

type Orchestrator struct {
    adapter        *adapter.Adapter
    maxConcurrency int
}

func (o *Orchestrator) RunParallel(tasks []adapter.Task) []adapter.Result {
    var wg sync.WaitGroup
    results := make([]adapter.Result, len(tasks))
    semaphore := make(chan struct{}, o.maxConcurrency)

    for i, task := range tasks {
        wg.Add(1)
        go func(idx int, t adapter.Task) {
            defer wg.Done()
            semaphore <- struct{}{}        // Acquire
            defer func() { <-semaphore }() // Release

            results[idx] = o.adapter.Execute(t)
        }(i, task)
    }

    wg.Wait()
    return results
}
```

#### TASK_003: Prompt Loader

**Priority**: P0-Critical **Dependencies**: None **Complexity**: Simple

**Implementation** (`internal/prompt/loader.go`):

```go
package prompt

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"
    "sync"
)

type Loader struct {
    agentDir string
    cache    map[string]string
    mu       sync.RWMutex
}

func (l *Loader) Load(agentType string) (string, error) {
    l.mu.RLock()
    if prompt, ok := l.cache[agentType]; ok {
        l.mu.RUnlock()
        return prompt, nil
    }
    l.mu.RUnlock()

    // Load from file
    path := filepath.Join(l.agentDir, agentType+".md")
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("reading agent file: %w", err)
    }

    // Parse YAML frontmatter
    prompt := l.extractPrompt(string(content))

    // Cache it
    l.mu.Lock()
    l.cache[agentType] = prompt
    l.mu.Unlock()

    return prompt, nil
}
```

#### TASK_004: Task Client

**Priority**: P1-High **Dependencies**: None **Complexity**: Medium

**Implementation** (`pkg/task/client.go`):

```go
package task

import (
    "bytes"
    "encoding/json"
    "os/exec"
)

type Client struct {
    claudePath string
}

func (c *Client) Run(agentType, prompt string, options map[string]interface{}) (Result, error) {
    // Prepare command
    cmd := exec.Command(c.claudePath, "task",
        "--subagent-type", agentType,
        "--prompt", prompt)

    // Set options as environment variables
    for k, v := range options {
        cmd.Env = append(cmd.Env, fmt.Sprintf("TASK_%s=%v",
            strings.ToUpper(k), v))
    }

    // Execute
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return Result{Error: stderr.String()}, err
    }

    return Result{Output: stdout.String()}, nil
}
```

### 3. Build Configuration

**Makefile** (root level):

```makefile
# Makefile for Morgana Protocol
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X main.Version=$(VERSION)
BINARY := morgana

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOGET := $(GOCMD) get

# Directories
CMD_DIR := ./cmd/morgana
DIST_DIR := ./dist

.PHONY: all build test clean

all: test build

build: build-darwin build-linux build-windows

build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-darwin-arm64 $(CMD_DIR)

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-linux-arm64 $(CMD_DIR)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY)-windows-amd64.exe $(CMD_DIR)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

clean:
	@echo "Cleaning..."
	rm -rf $(DIST_DIR)

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

install: build
	@echo "Installing to ~/.claude/bin..."
	./scripts/install.sh
```

### 4. Integration Scripts

**Installation Script** (`scripts/install.sh`):

```bash
#!/bin/bash
set -e

INSTALL_DIR="$HOME/.claude/bin"
BINARY_NAME="morgana"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Source binary
SOURCE="dist/${BINARY_NAME}-${OS}-${ARCH}"
if [[ "$OS" == "windows" ]]; then
    SOURCE="${SOURCE}.exe"
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

# Copy binary
cp "$SOURCE" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "Morgana Protocol installed to $INSTALL_DIR/$BINARY_NAME"
echo "Add $INSTALL_DIR to your PATH if not already done"
```

### 5. Testing Strategy

**Unit Tests** (`internal/adapter/adapter_test.go`):

```go
func TestParallelExecution(t *testing.T) {
    tasks := []Task{
        {AgentType: "code-implementer", Prompt: "test1"},
        {AgentType: "test-specialist", Prompt: "test2"},
    }

    orch := NewOrchestrator(2)
    results := orch.RunParallel(tasks)

    assert.Len(t, results, 2)
    for _, r := range results {
        assert.NotEmpty(t, r.Output)
        assert.Empty(t, r.Error)
    }
}
```

### 6. Usage Examples

**Shell Integration**:

```bash
# Single task
morgana --agent code-implementer --prompt "implement auth service"

# Parallel tasks via JSON
echo '[
  {"agent_type": "code-implementer", "prompt": "implement auth"},
  {"agent_type": "test-specialist", "prompt": "create tests"}
]' | morgana --parallel

# From markdown command
!morgana --parallel \
  --agent code-implementer --prompt "implement feature" \
  --agent test-specialist --prompt "create tests"
```

### 7. Deployment Plan

1. **Phase 1**: Local testing with mock Task client
2. **Phase 2**: Integration with real Claude Code Task tool
3. **Phase 3**: GitHub Actions for automated builds
4. **Phase 4**: Distribution via GitHub releases

### 8. Success Metrics

- [ ] Compiles for macOS (Intel/ARM), Linux, Windows
- [ ] Executes Task tool with proper agent type translation
- [ ] Handles 5+ parallel tasks without race conditions
- [ ] Sub-100ms overhead per task
- [ ] Binary size <10MB compressed
- [ ] 80%+ test coverage
- [ ] Zero dependencies at runtime

### 9. Risk Mitigation

1. **Task Tool API Changes**: Abstract interface allows easy updates
2. **Performance Issues**: Configurable concurrency limits
3. **Cross-Platform Bugs**: Comprehensive CI/CD testing
4. **Agent Prompt Changes**: File-based system allows hot updates

### 10. Next Steps

1. Create go.mod and initialize project
2. Implement core adapter logic
3. Add parallel orchestration
4. Create build pipeline
5. Integration testing
6. Documentation and examples
