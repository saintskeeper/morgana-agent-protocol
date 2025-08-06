package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	// TUI configuration
	TUI TUIConfig `yaml:"tui"`
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

// TUIConfig holds TUI-related configuration
type TUIConfig struct {
	// Enable TUI mode
	Enabled bool `yaml:"enabled"`

	// Performance settings
	Performance TUIPerformanceConfig `yaml:"performance"`

	// Visual settings
	Visual TUIVisualConfig `yaml:"visual"`

	// Feature flags
	Features TUIFeatureConfig `yaml:"features"`

	// Event system configuration
	Events TUIEventConfig `yaml:"events"`
}

// TUIPerformanceConfig holds performance-related TUI settings
type TUIPerformanceConfig struct {
	// Target refresh rate (e.g., "16ms" for 60fps)
	RefreshRate time.Duration `yaml:"refresh_rate"`

	// Maximum log lines to keep in memory
	MaxLogLines int `yaml:"max_log_lines"`

	// Enable performance optimizations
	OptimizedRendering bool `yaml:"optimized_rendering"`

	// FPS target (calculated from refresh rate, but can be overridden)
	TargetFPS int `yaml:"target_fps"`
}

// TUIVisualConfig holds visual appearance settings
type TUIVisualConfig struct {
	// Color theme configuration
	Theme TUIThemeConfig `yaml:"theme"`

	// Show debug information panel
	ShowDebugInfo bool `yaml:"show_debug_info"`

	// Show timestamps in logs
	ShowTimeStamps bool `yaml:"show_timestamps"`

	// Use compact display mode
	CompactMode bool `yaml:"compact_mode"`

	// Terminal size requirements
	MinWidth  int `yaml:"min_width"`
	MinHeight int `yaml:"min_height"`
}

// TUIThemeConfig holds color theme configuration
type TUIThemeConfig struct {
	// Theme name (dark, light, or custom)
	Name string `yaml:"name"`

	// Base colors
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary"`
	Success   string `yaml:"success"`
	Warning   string `yaml:"warning"`
	Error     string `yaml:"error"`
	Info      string `yaml:"info"`

	// UI element colors
	Background string `yaml:"background"`
	Foreground string `yaml:"foreground"`
	Border     string `yaml:"border"`
	Highlight  string `yaml:"highlight"`
	Muted      string `yaml:"muted"`

	// Component-specific colors
	ProgressBar TUIProgressBarTheme `yaml:"progress_bar"`

	// Agent-specific colors
	AgentColors map[string]string `yaml:"agent_colors"`
}

// TUIProgressBarTheme holds progress bar color configuration
type TUIProgressBarTheme struct {
	Complete   string `yaml:"complete"`
	Incomplete string `yaml:"incomplete"`
	Background string `yaml:"background"`
}

// TUIFeatureConfig holds feature flags for TUI functionality
type TUIFeatureConfig struct {
	// Enable log filtering capabilities
	EnableFiltering bool `yaml:"enable_filtering"`

	// Enable log search functionality
	EnableSearch bool `yaml:"enable_search"`

	// Enable export functionality
	EnableExport bool `yaml:"enable_export"`

	// Enable keyboard shortcuts
	EnableKeyboardShortcuts bool `yaml:"enable_keyboard_shortcuts"`

	// Enable mouse interaction
	EnableMouse bool `yaml:"enable_mouse"`

	// Enable auto-scroll for logs
	EnableAutoScroll bool `yaml:"enable_auto_scroll"`
}

// TUIEventConfig holds event system configuration for TUI
type TUIEventConfig struct {
	// Buffer size for event processing
	BufferSize int `yaml:"buffer_size"`

	// Event subscription filters
	SubscriptionFilters []string `yaml:"subscription_filters"`

	// Enable event batching for performance
	EnableBatching bool `yaml:"enable_batching"`

	// Batch size for event processing
	BatchSize int `yaml:"batch_size"`

	// Batch timeout for event processing
	BatchTimeout time.Duration `yaml:"batch_timeout"`

	// Event processing timeout
	ProcessingTimeout time.Duration `yaml:"processing_timeout"`
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
		TUI: TUIConfig{
			Enabled: true,
			Performance: TUIPerformanceConfig{
				RefreshRate:        16 * time.Millisecond, // 60fps
				MaxLogLines:        10000,
				OptimizedRendering: true,
				TargetFPS:          60,
			},
			Visual: TUIVisualConfig{
				Theme: TUIThemeConfig{
					Name:       "dark",
					Primary:    "#7C3AED", // Violet
					Secondary:  "#06B6D4", // Cyan
					Success:    "#10B981", // Emerald
					Warning:    "#F59E0B", // Amber
					Error:      "#EF4444", // Red
					Info:       "#3B82F6", // Blue
					Background: "#0F172A", // Slate 900
					Foreground: "#F1F5F9", // Slate 100
					Border:     "#334155", // Slate 700
					Highlight:  "#1E293B", // Slate 800
					Muted:      "#64748B", // Slate 500
					ProgressBar: TUIProgressBarTheme{
						Complete:   "#10B981", // Success color
						Incomplete: "#64748B", // Muted color
						Background: "#334155", // Border color
					},
					AgentColors: map[string]string{
						"code-implementer":  "#10B981", // Green
						"sprint-planner":    "#3B82F6", // Blue
						"test-specialist":   "#F59E0B", // Amber
						"validation-expert": "#EF4444", // Red
						"default":           "#64748B", // Slate
					},
				},
				ShowDebugInfo:  false,
				ShowTimeStamps: true,
				CompactMode:    false,
				MinWidth:       80,
				MinHeight:      24,
			},
			Features: TUIFeatureConfig{
				EnableFiltering:         true,
				EnableSearch:            true,
				EnableExport:            false,
				EnableKeyboardShortcuts: true,
				EnableMouse:             true,
				EnableAutoScroll:        true,
			},
			Events: TUIEventConfig{
				BufferSize:          1000,
				SubscriptionFilters: []string{}, // Subscribe to all events by default
				EnableBatching:      true,
				BatchSize:           50,
				BatchTimeout:        100 * time.Millisecond,
				ProcessingTimeout:   5 * time.Second,
			},
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

	// TUI overrides
	if enabled := os.Getenv("MORGANA_TUI_ENABLED"); enabled != "" {
		c.TUI.Enabled = enabled == "true"
	}
	if refreshRate := os.Getenv("MORGANA_TUI_REFRESH_RATE"); refreshRate != "" {
		if val, err := time.ParseDuration(refreshRate); err == nil {
			c.TUI.Performance.RefreshRate = val
		}
	}
	if maxLogLines := os.Getenv("MORGANA_TUI_MAX_LOG_LINES"); maxLogLines != "" {
		if val, err := strconv.Atoi(maxLogLines); err == nil && val > 0 {
			c.TUI.Performance.MaxLogLines = val
		}
	}
	if bufferSize := os.Getenv("MORGANA_TUI_BUFFER_SIZE"); bufferSize != "" {
		if val, err := strconv.Atoi(bufferSize); err == nil && val > 0 {
			c.TUI.Events.BufferSize = val
		}
	}
	if theme := os.Getenv("MORGANA_TUI_THEME"); theme != "" {
		c.TUI.Visual.Theme.Name = theme
	}
	if compactMode := os.Getenv("MORGANA_TUI_COMPACT_MODE"); compactMode == "true" {
		c.TUI.Visual.CompactMode = true
	}
	if showDebug := os.Getenv("MORGANA_TUI_SHOW_DEBUG"); showDebug == "true" {
		c.TUI.Visual.ShowDebugInfo = true
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

	// Validate TUI config
	if c.TUI.Enabled {
		if err := c.validateTUIConfig(); err != nil {
			return fmt.Errorf("tui configuration invalid: %w", err)
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

// validateTUIConfig validates TUI-specific configuration
func (c *Config) validateTUIConfig() error {
	// Validate performance settings
	if c.TUI.Performance.RefreshRate <= 0 {
		return fmt.Errorf("performance.refresh_rate must be positive, got %v", c.TUI.Performance.RefreshRate)
	}
	if c.TUI.Performance.RefreshRate < time.Millisecond {
		return fmt.Errorf("performance.refresh_rate too high, minimum is 1ms, got %v", c.TUI.Performance.RefreshRate)
	}
	if c.TUI.Performance.MaxLogLines <= 0 {
		return fmt.Errorf("performance.max_log_lines must be positive, got %d", c.TUI.Performance.MaxLogLines)
	}
	if c.TUI.Performance.TargetFPS <= 0 || c.TUI.Performance.TargetFPS > 120 {
		return fmt.Errorf("performance.target_fps must be between 1 and 120, got %d", c.TUI.Performance.TargetFPS)
	}

	// Validate visual settings
	if c.TUI.Visual.MinWidth < 40 {
		return fmt.Errorf("visual.min_width must be at least 40, got %d", c.TUI.Visual.MinWidth)
	}
	if c.TUI.Visual.MinHeight < 12 {
		return fmt.Errorf("visual.min_height must be at least 12, got %d", c.TUI.Visual.MinHeight)
	}

	// Validate theme settings
	if err := c.validateTUITheme(); err != nil {
		return fmt.Errorf("theme validation failed: %w", err)
	}

	// Validate event settings
	if c.TUI.Events.BufferSize <= 0 {
		return fmt.Errorf("events.buffer_size must be positive, got %d", c.TUI.Events.BufferSize)
	}
	if c.TUI.Events.BatchSize <= 0 {
		return fmt.Errorf("events.batch_size must be positive, got %d", c.TUI.Events.BatchSize)
	}
	if c.TUI.Events.BatchTimeout <= 0 {
		return fmt.Errorf("events.batch_timeout must be positive, got %v", c.TUI.Events.BatchTimeout)
	}
	if c.TUI.Events.ProcessingTimeout <= 0 {
		return fmt.Errorf("events.processing_timeout must be positive, got %v", c.TUI.Events.ProcessingTimeout)
	}

	return nil
}

// validateTUITheme validates theme-specific configuration
func (c *Config) validateTUITheme() error {
	theme := &c.TUI.Visual.Theme

	// Validate theme name
	validThemes := []string{"dark", "light", "custom"}
	if !contains(validThemes, theme.Name) {
		return fmt.Errorf("theme.name must be one of: %v, got %s", validThemes, theme.Name)
	}

	// Validate hex colors (basic check)
	colors := map[string]string{
		"primary":    theme.Primary,
		"secondary":  theme.Secondary,
		"success":    theme.Success,
		"warning":    theme.Warning,
		"error":      theme.Error,
		"info":       theme.Info,
		"background": theme.Background,
		"foreground": theme.Foreground,
		"border":     theme.Border,
		"highlight":  theme.Highlight,
		"muted":      theme.Muted,
	}

	for name, color := range colors {
		if color != "" && !isValidHexColor(color) {
			return fmt.Errorf("theme.%s is not a valid hex color: %s", name, color)
		}
	}

	// Validate progress bar colors
	if theme.ProgressBar.Complete != "" && !isValidHexColor(theme.ProgressBar.Complete) {
		return fmt.Errorf("theme.progress_bar.complete is not a valid hex color: %s", theme.ProgressBar.Complete)
	}
	if theme.ProgressBar.Incomplete != "" && !isValidHexColor(theme.ProgressBar.Incomplete) {
		return fmt.Errorf("theme.progress_bar.incomplete is not a valid hex color: %s", theme.ProgressBar.Incomplete)
	}
	if theme.ProgressBar.Background != "" && !isValidHexColor(theme.ProgressBar.Background) {
		return fmt.Errorf("theme.progress_bar.background is not a valid hex color: %s", theme.ProgressBar.Background)
	}

	// Validate agent colors
	for agent, color := range theme.AgentColors {
		if color != "" && !isValidHexColor(color) {
			return fmt.Errorf("theme.agent_colors.%s is not a valid hex color: %s", agent, color)
		}
	}

	return nil
}

// GetTUIAgentColor returns the color for a specific agent type
func (c *Config) GetTUIAgentColor(agentType string) string {
	if color, ok := c.TUI.Visual.Theme.AgentColors[agentType]; ok && color != "" {
		return color
	}
	if defaultColor, ok := c.TUI.Visual.Theme.AgentColors["default"]; ok && defaultColor != "" {
		return defaultColor
	}
	return c.TUI.Visual.Theme.Muted // Fallback to muted color
}

// IsTUIEnabled returns whether TUI is enabled
func (c *Config) IsTUIEnabled() bool {
	return c.TUI.Enabled
}

// GetTUIRefreshRate returns the configured refresh rate
func (c *Config) GetTUIRefreshRate() time.Duration {
	return c.TUI.Performance.RefreshRate
}

// UpdateTUIConfig allows runtime updates to TUI configuration
func (c *Config) UpdateTUIConfig(updates TUIConfig) error {
	// Create a temporary config to validate changes
	tempConfig := *c
	tempConfig.TUI = updates

	// Validate the updated configuration
	if err := tempConfig.validateTUIConfig(); err != nil {
		return fmt.Errorf("invalid TUI config update: %w", err)
	}

	// Apply the updates
	c.TUI = updates
	return nil
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

func isValidHexColor(color string) bool {
	if color == "" {
		return true // Empty is valid (will use default)
	}

	if !strings.HasPrefix(color, "#") {
		return false
	}

	color = strings.TrimPrefix(color, "#")
	if len(color) != 6 {
		return false
	}

	for _, r := range color {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}

	return true
}
