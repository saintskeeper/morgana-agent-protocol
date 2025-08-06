package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete application configuration
type Config struct {
	// Agent configuration
	Agents AgentsConfig `yaml:"agents"`

	// Execution configuration
	Execution ExecutionConfig `yaml:"execution"`

	// Telemetry configuration
	Telemetry TelemetryConfig `yaml:"telemetry"`

	// Task client configuration
	TaskClient TaskClientConfig `yaml:"task_client"`
}

// AgentsConfig holds agent-related configuration
type AgentsConfig struct {
	// Directory containing agent prompts
	PromptDir string `yaml:"prompt_dir"`

	// Default timeout for agent execution
	DefaultTimeout time.Duration `yaml:"default_timeout"`

	// Agent-specific timeouts
	Timeouts map[string]time.Duration `yaml:"timeouts"`

	// Retry configuration
	Retry RetryConfig `yaml:"retry"`
}

// ExecutionConfig holds execution-related configuration
type ExecutionConfig struct {
	// Maximum concurrent executions
	MaxConcurrency int `yaml:"max_concurrency"`

	// Default execution mode (sequential/parallel)
	DefaultMode string `yaml:"default_mode"`

	// Queue size for parallel execution
	QueueSize int `yaml:"queue_size"`
}

// TelemetryConfig holds telemetry configuration
type TelemetryConfig struct {
	// Enable/disable telemetry
	Enabled bool `yaml:"enabled"`

	// Exporter type (stdout, otlp, none)
	Exporter string `yaml:"exporter"`

	// OTLP endpoint
	OTLPEndpoint string `yaml:"otlp_endpoint"`

	// Service name
	ServiceName string `yaml:"service_name"`

	// Environment (production, staging, development)
	Environment string `yaml:"environment"`

	// Sampling rate (0.0 to 1.0)
	SamplingRate float64 `yaml:"sampling_rate"`
}

// TaskClientConfig holds task client configuration
type TaskClientConfig struct {
	// Python bridge script path
	BridgePath string `yaml:"bridge_path"`

	// Python executable path
	PythonPath string `yaml:"python_path"`

	// Enable mock mode
	MockMode bool `yaml:"mock_mode"`

	// Command timeout
	Timeout time.Duration `yaml:"timeout"`
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	// Enable retries
	Enabled bool `yaml:"enabled"`

	// Maximum retry attempts
	MaxAttempts int `yaml:"max_attempts"`

	// Initial backoff duration
	InitialBackoff time.Duration `yaml:"initial_backoff"`

	// Maximum backoff duration
	MaxBackoff time.Duration `yaml:"max_backoff"`

	// Backoff multiplier
	Multiplier float64 `yaml:"multiplier"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		Agents: AgentsConfig{
			PromptDir:      filepath.Join(homeDir, ".claude", "agents"),
			DefaultTimeout: 2 * time.Minute,
			Timeouts:       make(map[string]time.Duration),
			Retry: RetryConfig{
				Enabled:        true,
				MaxAttempts:    3,
				InitialBackoff: 1 * time.Second,
				MaxBackoff:     30 * time.Second,
				Multiplier:     2.0,
			},
		},
		Execution: ExecutionConfig{
			MaxConcurrency: 5,
			DefaultMode:    "sequential",
			QueueSize:      100,
		},
		Telemetry: TelemetryConfig{
			Enabled:      true,
			Exporter:     "stdout",
			OTLPEndpoint: "localhost:4317",
			ServiceName:  "morgana-protocol",
			Environment:  "production",
			SamplingRate: 0.1,
		},
		TaskClient: TaskClientConfig{
			BridgePath: "", // Auto-discover
			PythonPath: "python3",
			MockMode:   false,
			Timeout:    5 * time.Minute,
		},
	}
}

// LoadFile loads configuration from a YAML file
func LoadFile(path string) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Apply environment overrides
	cfg.applyEnvOverrides()

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}

// applyEnvOverrides applies environment variable overrides
func (c *Config) applyEnvOverrides() {
	// Agent overrides
	if dir := os.Getenv("MORGANA_AGENT_DIR"); dir != "" {
		c.Agents.PromptDir = dir
	}

	// Execution overrides
	if conc := os.Getenv("MORGANA_MAX_CONCURRENCY"); conc != "" {
		if val, err := parseInt(conc); err == nil {
			c.Execution.MaxConcurrency = val
		}
	}

	// Telemetry overrides
	if exporter := os.Getenv("MORGANA_OTEL_EXPORTER"); exporter != "" {
		c.Telemetry.Exporter = exporter
	}
	if endpoint := os.Getenv("MORGANA_OTEL_ENDPOINT"); endpoint != "" {
		c.Telemetry.OTLPEndpoint = endpoint
	}

	// Task client overrides
	if bridge := os.Getenv("MORGANA_BRIDGE_PATH"); bridge != "" {
		c.TaskClient.BridgePath = bridge
	}
	if mock := os.Getenv("MORGANA_MOCK_MODE"); mock == "true" {
		c.TaskClient.MockMode = true
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate agents config
	if c.Agents.PromptDir == "" {
		return fmt.Errorf("agents.prompt_dir is required")
	}
	if c.Agents.DefaultTimeout <= 0 {
		return fmt.Errorf("agents.default_timeout must be positive")
	}

	// Validate execution config
	if c.Execution.MaxConcurrency <= 0 {
		return fmt.Errorf("execution.max_concurrency must be positive")
	}
	if c.Execution.DefaultMode != "sequential" && c.Execution.DefaultMode != "parallel" {
		return fmt.Errorf("execution.default_mode must be 'sequential' or 'parallel'")
	}

	// Validate telemetry config
	if c.Telemetry.Enabled {
		validExporters := []string{"stdout", "otlp", "none"}
		if !contains(validExporters, c.Telemetry.Exporter) {
			return fmt.Errorf("telemetry.exporter must be one of: %v", validExporters)
		}
		if c.Telemetry.SamplingRate < 0 || c.Telemetry.SamplingRate > 1 {
			return fmt.Errorf("telemetry.sampling_rate must be between 0.0 and 1.0")
		}
	}

	// Validate retry config
	if c.Agents.Retry.Enabled {
		if c.Agents.Retry.MaxAttempts <= 0 {
			return fmt.Errorf("agents.retry.max_attempts must be positive")
		}
		if c.Agents.Retry.Multiplier <= 1 {
			return fmt.Errorf("agents.retry.multiplier must be greater than 1")
		}
	}

	return nil
}

// GetAgentTimeout returns the timeout for a specific agent type
func (c *Config) GetAgentTimeout(agentType string) time.Duration {
	if timeout, ok := c.Agents.Timeouts[agentType]; ok {
		return timeout
	}
	return c.Agents.DefaultTimeout
}

// Helper functions

func parseInt(s string) (int, error) {
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	return val, err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
