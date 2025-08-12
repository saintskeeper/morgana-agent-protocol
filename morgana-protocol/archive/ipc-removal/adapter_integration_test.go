//go:build integration
// +build integration

package adapter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/prompt"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/pkg/task"
	"go.opentelemetry.io/otel/trace/noop"
)

// Define interface that our mock will implement
type TaskClient interface {
	RunWithContext(ctx context.Context, agentType, prompt string, options map[string]interface{}) (*task.Result, error)
	Run(agentType, prompt string, options map[string]interface{}) (*task.Result, error)
}

// MockTaskClient for testing
type MockTaskClient struct {
	responses map[string]*task.Result
	errors    map[string]error
	delays    map[string]time.Duration
}

func NewMockTaskClient() *MockTaskClient {
	return &MockTaskClient{
		responses: make(map[string]*task.Result),
		errors:    make(map[string]error),
		delays:    make(map[string]time.Duration),
	}
}

func (m *MockTaskClient) RunWithContext(ctx context.Context, agentType, prompt string, options map[string]interface{}) (*task.Result, error) {
	// Check for configured delay
	if delay, ok := m.delays[agentType]; ok {
		select {
		case <-time.After(delay):
			// Continue after delay
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Check for configured error
	if err, ok := m.errors[agentType]; ok {
		return nil, err
	}

	// Return configured response
	if resp, ok := m.responses[agentType]; ok {
		return resp, nil
	}

	// Default response
	return &task.Result{
		Output: "Mock response for " + agentType,
	}, nil
}

func (m *MockTaskClient) Run(agentType, prompt string, options map[string]interface{}) (*task.Result, error) {
	return m.RunWithContext(context.Background(), agentType, prompt, options)
}

func TestAdapterIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test prompt loader with testdata directory
	testDataDir := filepath.Join("testdata")
	promptLoader := prompt.NewPromptLoader(testDataDir)

	// Create task client in mock mode
	taskClient := task.NewClientWithConfig("", "python3", true)

	// Create adapter with noop tracer
	tracer := noop.NewTracerProvider().Tracer("test")
	adapter := New(promptLoader, taskClient, tracer)

	// Set default timeout
	adapter.SetTimeouts(5*time.Second, map[string]time.Duration{
		"code-implementer": 2 * time.Second,
	})

	tests := []struct {
		name          string
		task          Task
		expectError   bool
		errorContains string
	}{
		{
			name: "successful execution",
			task: Task{
				AgentType: "code-implementer",
				Prompt:    "Test prompt",
			},
			expectError: false,
		},
		{
			name: "invalid agent type",
			task: Task{
				AgentType: "invalid-agent",
				Prompt:    "Test prompt",
			},
			expectError:   true,
			errorContains: "unknown agent type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result := adapter.Execute(ctx, tt.task)

			if tt.expectError {
				if result.Error == "" {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !containsString(result.Error, tt.errorContains) {
					t.Errorf("Expected error containing '%s', got: %s", tt.errorContains, result.Error)
				}
			} else {
				if result.Error != "" {
					t.Errorf("Unexpected error: %s", result.Error)
				}
				if result.Output == "" {
					t.Error("Expected output but got empty string")
				}
			}
		})
	}
}

func TestAdapterTimeout(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	testDataDir := filepath.Join("testdata")
	promptLoader := prompt.NewPromptLoader(testDataDir)

	// For timeout testing, we'll use the real Python bridge with a slow script
	testBridge := filepath.Join("..", "..", "pkg", "task", "testdata", "test_bridge.py")
	taskClient := task.NewClientWithConfig(testBridge, "python3", false)

	tracer := noop.NewTracerProvider().Tracer("test")
	adapter := New(promptLoader, taskClient, tracer)

	// Set very short timeout
	adapter.SetTimeouts(1*time.Second, map[string]time.Duration{
		"code-implementer": 500 * time.Millisecond,
	})

	// Set TEST_MODE to timeout
	os.Setenv("TEST_MODE", "timeout")
	defer os.Unsetenv("TEST_MODE")

	// Execute task that will timeout
	ctx := context.Background()
	result := adapter.Execute(ctx, Task{
		AgentType: "code-implementer",
		Prompt:    "Test prompt",
	})

	// Verify timeout error - the error message confirms it's working
	if result.Error == "" {
		t.Error("Expected timeout error")
		return
	}

	// The actual error message is "executing task: task execution timed out"
	// which confirms the timeout is working correctly
	if !containsString(result.Error, "timed out") && !containsString(result.Error, "timeout") {
		t.Errorf("Expected timeout-related error, got: %s", result.Error)
	}
}

func TestAdapterConcurrentExecution(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	testDataDir := filepath.Join("testdata")
	promptLoader := prompt.NewPromptLoader(testDataDir)

	// Use mock mode for concurrent testing
	taskClient := task.NewClientWithConfig("", "python3", true)

	tracer := noop.NewTracerProvider().Tracer("test")
	adapter := New(promptLoader, taskClient, tracer)
	adapter.SetTimeouts(5*time.Second, nil)

	// Execute multiple tasks concurrently
	numTasks := 10
	results := make(chan Result, numTasks)

	ctx := context.Background()
	for i := 0; i < numTasks; i++ {
		go func(id int) {
			result := adapter.Execute(ctx, Task{
				AgentType: "code-implementer",
				Prompt:    fmt.Sprintf("Test prompt %d", id),
			})
			results <- result
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < numTasks; i++ {
		result := <-results
		if result.Error == "" {
			successCount++
		}
	}

	if successCount != numTasks {
		t.Errorf("Expected %d successful executions, got %d", numTasks, successCount)
	}
}

func TestAdapterAgentSpecificTimeouts(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	testDataDir := filepath.Join("testdata")
	promptLoader := prompt.NewPromptLoader(testDataDir)

	// Use test bridge for timeout testing
	testBridge := filepath.Join("..", "..", "pkg", "task", "testdata", "test_bridge.py")
	taskClient := task.NewClientWithConfig(testBridge, "python3", false)

	tracer := noop.NewTracerProvider().Tracer("test")
	adapter := New(promptLoader, taskClient, tracer)

	// Configure different timeouts
	adapter.SetTimeouts(3*time.Second, map[string]time.Duration{
		"code-implementer": 1 * time.Second, // Will timeout
		"sprint-planner":   5 * time.Second, // Won't timeout
	})

	tests := []struct {
		name          string
		agentType     string
		expectTimeout bool
	}{
		{
			name:          "agent with short timeout",
			agentType:     "code-implementer",
			expectTimeout: true,
		},
		{
			name:          "agent with long timeout",
			agentType:     "sprint-planner",
			expectTimeout: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set TEST_MODE to timeout for testing
			if tt.expectTimeout {
				os.Setenv("TEST_MODE", "timeout")
				defer os.Unsetenv("TEST_MODE")
			} else {
				os.Setenv("TEST_MODE", "success")
				defer os.Unsetenv("TEST_MODE")
			}

			ctx := context.Background()
			result := adapter.Execute(ctx, Task{
				AgentType: tt.agentType,
				Prompt:    "Test prompt",
			})

			if tt.expectTimeout {
				if result.Error == "" {
					t.Error("Expected timeout error")
				}
			} else {
				if result.Error != "" {
					t.Errorf("Unexpected error: %s", result.Error)
				}
			}
		})
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && containsString(s[1:], substr)
}
