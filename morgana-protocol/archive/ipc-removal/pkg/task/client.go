package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Result represents the output from the Task tool
type Result struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// Client interfaces with the Claude Code Task tool
type Client struct {
	claudePath string
	debug      bool
	bridgePath string
	pythonPath string
	mockMode   bool
}

// NewClient creates a new Task client
func NewClient() *Client {
	// Try to find claude in PATH or common locations
	claudePath := "claude"
	if path, err := exec.LookPath("claude"); err == nil {
		claudePath = path
	}

	return &Client{
		claudePath: claudePath,
		debug:      os.Getenv("MORGANA_DEBUG") == "true",
		bridgePath: os.Getenv("MORGANA_BRIDGE_PATH"),
		pythonPath: "python3",
		mockMode:   false,
	}
}

// NewClientWithConfig creates a new Task client with configuration
func NewClientWithConfig(bridgePath, pythonPath string, mockMode bool) *Client {
	// Try to find claude in PATH or common locations
	claudePath := "claude"
	if path, err := exec.LookPath("claude"); err == nil {
		claudePath = path
	}

	// Override with environment variable if set
	if envBridge := os.Getenv("MORGANA_BRIDGE_PATH"); envBridge != "" {
		bridgePath = envBridge
	}

	return &Client{
		claudePath: claudePath,
		debug:      os.Getenv("MORGANA_DEBUG") == "true",
		bridgePath: bridgePath,
		pythonPath: pythonPath,
		mockMode:   mockMode,
	}
}

// RunWithContext executes a task using the Claude Code Task tool with context support
func (c *Client) RunWithContext(ctx context.Context, agentType, prompt string, options map[string]interface{}) (*Result, error) {
	if c.debug {
		fmt.Printf("DEBUG: Task request - Agent: %s, Prompt length: %d\n", agentType, len(prompt))
	}

	// If mock mode is enabled, return mock response
	if c.mockMode {
		return c.mockResponse(agentType, prompt), nil
	}

	// Prepare input for Python bridge
	input := map[string]interface{}{
		"agent_type": agentType,
		"prompt":     prompt,
	}

	// Add any additional options
	for k, v := range options {
		input[k] = v
	}

	// Marshal to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshaling input: %w", err)
	}

	// Find Python bridge script
	bridgePath := c.bridgePath
	if bridgePath == "" {
		bridgePath = c.findBridgeScript()
	}
	if bridgePath == "" {
		return nil, fmt.Errorf("task_bridge.py not found. Please ensure it exists in scripts/ directory")
	}

	// Execute Python bridge with context
	cmd := exec.CommandContext(ctx, c.pythonPath, bridgePath)
	cmd.Stdin = bytes.NewReader(inputJSON)

	// Set environment
	cmd.Env = os.Environ()
	if c.debug {
		cmd.Env = append(cmd.Env, "MORGANA_DEBUG=true")
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	if err := cmd.Run(); err != nil {
		// Check if context was cancelled
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("task execution timed out")
		}
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("task execution cancelled")
		}
		return nil, fmt.Errorf("executing bridge: %w (stderr: %s)", err, stderr.String())
	}

	// Parse response
	var response struct {
		Success bool   `json:"success"`
		Output  string `json:"output"`
		Error   string `json:"error"`
		Mock    bool   `json:"mock"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return nil, fmt.Errorf("parsing response: %w (output: %s)", err, stdout.String())
	}

	// Check for errors
	if !response.Success {
		return nil, fmt.Errorf("task execution failed: %s", response.Error)
	}

	if c.debug && response.Mock {
		fmt.Println("DEBUG: Using mock response (Task function not available)")
	}

	return &Result{
		Output: response.Output,
	}, nil
}

// buildCommand constructs the command to execute
func (c *Client) buildCommand(agentType, prompt string, options map[string]interface{}) *exec.Cmd {
	// Build command arguments
	args := []string{"task", "--subagent-type", agentType}

	// Add prompt via stdin to handle multiline prompts
	cmd := exec.Command(c.claudePath, args...)
	cmd.Stdin = strings.NewReader(prompt)

	// Set environment variables for options
	env := os.Environ()
	for k, v := range options {
		envVar := fmt.Sprintf("TASK_%s=%v", strings.ToUpper(k), v)
		env = append(env, envVar)
	}
	cmd.Env = env

	return cmd
}

// mockResponse returns a mock response for testing
func (c *Client) mockResponse(agentType, prompt string) *Result {
	// Add realistic delay to allow TUI to show progress
	// This simulates actual agent processing time
	// Vary delay based on prompt length for more realistic simulation
	baseDelay := 1500 * time.Millisecond
	extraDelay := time.Duration(len(prompt)*2) * time.Millisecond
	totalDelay := baseDelay + extraDelay

	time.Sleep(totalDelay)

	return &Result{
		Output: fmt.Sprintf("[MOCK] Executed %s agent with prompt length: %d (simulated %dms)", agentType, len(prompt), totalDelay.Milliseconds()),
	}
}

// findBridgeScript locates the Python bridge script
func (c *Client) findBridgeScript() string {
	// Check common locations
	locations := []string{
		// Relative to current directory
		"./scripts/task_bridge.py",
		"../scripts/task_bridge.py",
		// Relative to binary location
		"", // Will be filled dynamically
		// Standard morgana-protocol installation
		os.ExpandEnv("$HOME/.claude/morgana-protocol/scripts/task_bridge.py"),
		"/usr/local/share/morgana-protocol/scripts/task_bridge.py",
	}

	// Add location relative to binary
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		locations[2] = filepath.Join(dir, "..", "scripts", "task_bridge.py")
	}

	// Check each location
	for _, loc := range locations {
		if loc != "" {
			if _, err := os.Stat(loc); err == nil {
				return loc
			}
		}
	}

	// Check if MORGANA_BRIDGE_PATH is set
	if bridgePath := os.Getenv("MORGANA_BRIDGE_PATH"); bridgePath != "" {
		if _, err := os.Stat(bridgePath); err == nil {
			return bridgePath
		}
	}

	return ""
}

// executeCommand runs the command and captures output
func (c *Client) executeCommand(cmd *exec.Cmd) (*Result, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return &Result{
			Error: stderr.String(),
		}, fmt.Errorf("command failed: %w", err)
	}

	// Try to parse JSON response
	var result Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		// If not JSON, treat entire output as the result
		result.Output = stdout.String()
	}

	return &result, nil
}

// Run executes a task using the Claude Code Task tool (backward compatibility)
func (c *Client) Run(agentType, prompt string, options map[string]interface{}) (*Result, error) {
	return c.RunWithContext(context.Background(), agentType, prompt, options)
}
