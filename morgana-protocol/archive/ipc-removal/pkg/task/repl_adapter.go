package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// REPLAdapter implements task execution via filesystem IPC for REPL environments
type REPLAdapter struct {
	TaskDir      string
	PollInterval time.Duration
	Timeout      time.Duration
}

// NewREPLAdapter creates a new REPL adapter
func NewREPLAdapter() *REPLAdapter {
	taskDir := os.Getenv("MORGANA_TASK_DIR")
	if taskDir == "" {
		taskDir = "/tmp/morgana/tasks"
	}

	return &REPLAdapter{
		TaskDir:      taskDir,
		PollInterval: 500 * time.Millisecond,
		Timeout:      5 * time.Minute,
	}
}

// TaskRequest represents a task to be executed by the REPL
type TaskRequest struct {
	ID        string                 `json:"id"`
	AgentType string                 `json:"agent_type"`
	Prompt    string                 `json:"prompt"`
	Options   map[string]interface{} `json:"options,omitempty"`
	CreatedAt int64                  `json:"created_at"`
}

// TaskResponse represents the result from REPL execution
type TaskResponse struct {
	Success     bool   `json:"success"`
	Output      string `json:"output"`
	Error       string `json:"error,omitempty"`
	TaskID      string `json:"task_id"`
	AgentType   string `json:"agent_type"`
	CompletedAt int64  `json:"completed_at"`
}

// ensureDirectories creates the necessary directory structure
func (r *REPLAdapter) ensureDirectories() error {
	dirs := []string{
		filepath.Join(r.TaskDir, "pending"),
		filepath.Join(r.TaskDir, "processing"),
		filepath.Join(r.TaskDir, "complete"),
		filepath.Join(r.TaskDir, "failed"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// Execute submits a task to the REPL and waits for completion
func (r *REPLAdapter) Execute(ctx context.Context, agentType, prompt string, options map[string]interface{}) (*Result, error) {
	// Ensure directories exist
	if err := r.ensureDirectories(); err != nil {
		return nil, err
	}

	// Generate task ID
	taskID := uuid.New().String()

	// Create task request
	request := TaskRequest{
		ID:        taskID,
		AgentType: agentType,
		Prompt:    prompt,
		Options:   options,
		CreatedAt: time.Now().Unix(),
	}

	// Write to pending directory
	pendingPath := filepath.Join(r.TaskDir, "pending", fmt.Sprintf("%s.json", taskID))
	requestData, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	if err := ioutil.WriteFile(pendingPath, requestData, 0644); err != nil {
		return nil, fmt.Errorf("writing task file: %w", err)
	}

	// Wait for result
	resultPath := filepath.Join(r.TaskDir, "complete", fmt.Sprintf("%s.result.json", taskID))

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	// Poll for result
	ticker := time.NewTicker(r.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			// Clean up pending file on timeout
			os.Remove(pendingPath)
			return nil, fmt.Errorf("task execution timed out after %v", r.Timeout)

		case <-ticker.C:
			// Check if result exists
			if _, err := os.Stat(resultPath); err == nil {
				// Read result
				resultData, err := ioutil.ReadFile(resultPath)
				if err != nil {
					return nil, fmt.Errorf("reading result: %w", err)
				}

				var response TaskResponse
				if err := json.Unmarshal(resultData, &response); err != nil {
					return nil, fmt.Errorf("parsing result: %w", err)
				}

				// Clean up result file
				os.Remove(resultPath)

				// Convert to Result
				result := &Result{
					Output: response.Output,
				}

				if !response.Success {
					result.Error = response.Error
				}

				return result, nil
			}

			// Check if task failed
			failedPath := filepath.Join(r.TaskDir, "failed", fmt.Sprintf("%s.json", taskID))
			if _, err := os.Stat(failedPath); err == nil {
				// Read failed task for error details
				failedData, _ := ioutil.ReadFile(failedPath)
				os.Remove(failedPath)

				return &Result{
					Error: fmt.Sprintf("Task failed: %s", string(failedData)),
				}, nil
			}
		}
	}
}

// ExecuteAsync submits a task without waiting for completion
func (r *REPLAdapter) ExecuteAsync(ctx context.Context, agentType, prompt string, options map[string]interface{}) (string, error) {
	// Ensure directories exist
	if err := r.ensureDirectories(); err != nil {
		return "", err
	}

	// Generate task ID
	taskID := uuid.New().String()

	// Create task request
	request := TaskRequest{
		ID:        taskID,
		AgentType: agentType,
		Prompt:    prompt,
		Options:   options,
		CreatedAt: time.Now().Unix(),
	}

	// Write to pending directory
	pendingPath := filepath.Join(r.TaskDir, "pending", fmt.Sprintf("%s.json", taskID))
	requestData, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	if err := ioutil.WriteFile(pendingPath, requestData, 0644); err != nil {
		return "", fmt.Errorf("writing task file: %w", err)
	}

	return taskID, nil
}

// GetResult retrieves the result of an async task
func (r *REPLAdapter) GetResult(taskID string) (*Result, error) {
	resultPath := filepath.Join(r.TaskDir, "complete", fmt.Sprintf("%s.result.json", taskID))

	// Check if result exists
	if _, err := os.Stat(resultPath); err != nil {
		// Check if failed
		failedPath := filepath.Join(r.TaskDir, "failed", fmt.Sprintf("%s.json", taskID))
		if _, err := os.Stat(failedPath); err == nil {
			return &Result{
				Error: "Task failed - check failed directory",
			}, nil
		}

		// Still pending or processing
		return nil, fmt.Errorf("result not ready")
	}

	// Read result
	resultData, err := ioutil.ReadFile(resultPath)
	if err != nil {
		return nil, fmt.Errorf("reading result: %w", err)
	}

	var response TaskResponse
	if err := json.Unmarshal(resultData, &response); err != nil {
		return nil, fmt.Errorf("parsing result: %w", err)
	}

	// Clean up result file
	os.Remove(resultPath)

	// Convert to Result
	result := &Result{
		Output: response.Output,
	}

	if !response.Success {
		result.Error = response.Error
	}

	return result, nil
}

// Cleanup removes old task files
func (r *REPLAdapter) Cleanup(olderThan time.Duration) error {
	dirs := []string{
		filepath.Join(r.TaskDir, "complete"),
		filepath.Join(r.TaskDir, "failed"),
	}

	cutoff := time.Now().Add(-olderThan)

	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.ModTime().Before(cutoff) {
				os.Remove(filepath.Join(dir, file.Name()))
			}
		}
	}

	return nil
}
