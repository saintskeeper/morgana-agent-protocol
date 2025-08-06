//go:build integration
// +build integration

package task

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestClientIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get test bridge path
	testBridge := filepath.Join("testdata", "test_bridge.py")
	if _, err := os.Stat(testBridge); err != nil {
		t.Fatalf("Test bridge not found: %v", err)
	}

	tests := []struct {
		name          string
		testMode      string
		agentType     string
		prompt        string
		timeout       time.Duration
		expectError   bool
		errorContains string
	}{
		{
			name:        "successful execution",
			testMode:    "success",
			agentType:   "test-agent",
			prompt:      "test prompt",
			timeout:     2 * time.Second,
			expectError: false,
		},
		{
			name:          "error from bridge",
			testMode:      "error",
			agentType:     "test-agent",
			prompt:        "test prompt",
			timeout:       2 * time.Second,
			expectError:   true,
			errorContains: "Test error scenario",
		},
		{
			name:          "timeout handling",
			testMode:      "timeout",
			agentType:     "test-agent",
			prompt:        "test prompt",
			timeout:       1 * time.Second,
			expectError:   true,
			errorContains: "task execution timed out",
		},
		{
			name:          "invalid json response",
			testMode:      "invalid_json",
			agentType:     "test-agent",
			prompt:        "test prompt",
			timeout:       2 * time.Second,
			expectError:   true,
			errorContains: "parsing response",
		},
		{
			name:          "bridge crash",
			testMode:      "crash",
			agentType:     "test-agent",
			prompt:        "test prompt",
			timeout:       2 * time.Second,
			expectError:   true,
			errorContains: "task execution failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test mode
			os.Setenv("TEST_MODE", tt.testMode)
			defer os.Unsetenv("TEST_MODE")

			// Create client with test bridge
			client := NewClientWithConfig(testBridge, "python3", false)

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			// Execute
			result, err := client.RunWithContext(ctx, tt.agentType, tt.prompt, nil)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected result but got nil")
				} else if result.Output == "" {
					t.Error("Expected output but got empty string")
				}
			}
		})
	}
}

func TestClientMockMode(t *testing.T) {
	client := NewClientWithConfig("", "python3", true)

	ctx := context.Background()
	result, err := client.RunWithContext(ctx, "test-agent", "test prompt", nil)

	if err != nil {
		t.Fatalf("Unexpected error in mock mode: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result in mock mode")
	}

	if !contains(result.Output, "[MOCK]") {
		t.Errorf("Expected mock output, got: %s", result.Output)
	}
}

func TestClientBridgeDiscovery(t *testing.T) {
	// Test findBridgeScript method
	client := NewClient()

	// Save current directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Create temporary directory structure
	tmpDir := t.TempDir()
	scriptsDir := filepath.Join(tmpDir, "scripts")
	os.MkdirAll(scriptsDir, 0755)

	// Create test bridge script
	bridgePath := filepath.Join(scriptsDir, "task_bridge.py")
	os.WriteFile(bridgePath, []byte("#!/usr/bin/env python3\n"), 0755)

	// Change to tmp directory
	os.Chdir(tmpDir)

	// Test discovery
	found := client.findBridgeScript()
	if found == "" {
		t.Error("Failed to find bridge script in ./scripts/")
	}

	// Test with environment variable
	customPath := filepath.Join(tmpDir, "custom_bridge.py")
	os.WriteFile(customPath, []byte("#!/usr/bin/env python3\n"), 0755)
	os.Setenv("MORGANA_BRIDGE_PATH", customPath)
	defer os.Unsetenv("MORGANA_BRIDGE_PATH")

	// Create new client to test env var
	client2 := NewClient()
	if client2.bridgePath != customPath {
		t.Errorf("Expected client to use custom bridge from env var %s, got %s", customPath, client2.bridgePath)
	}
}

func TestClientContextCancellation(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testBridge := filepath.Join("testdata", "test_bridge.py")
	client := NewClientWithConfig(testBridge, "python3", false)

	// Set timeout mode
	os.Setenv("TEST_MODE", "timeout")
	defer os.Unsetenv("TEST_MODE")

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	// Execute
	_, err := client.RunWithContext(ctx, "test-agent", "test prompt", nil)

	if err == nil {
		t.Error("Expected error from cancelled context")
	}

	if !contains(err.Error(), "cancelled") {
		t.Errorf("Expected cancellation error, got: %v", err)
	}
}

func TestClientOptions(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testBridge := filepath.Join("testdata", "test_bridge.py")
	client := NewClientWithConfig(testBridge, "python3", false)

	os.Setenv("TEST_MODE", "success")
	defer os.Unsetenv("TEST_MODE")

	// Test with options
	options := map[string]interface{}{
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	ctx := context.Background()
	result, err := client.RunWithContext(ctx, "test-agent", "test prompt", options)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil || result.Output == "" {
		t.Error("Expected result with output")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}
