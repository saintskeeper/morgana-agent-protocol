package config

import (
	"os"
	"testing"
	"time"
)

func TestTUIConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      TUIConfig
		expectError bool
	}{
		{
			name: "valid default config",
			config: TUIConfig{
				Enabled: true,
				Performance: TUIPerformanceConfig{
					RefreshRate:        16 * time.Millisecond,
					MaxLogLines:        10000,
					OptimizedRendering: true,
					TargetFPS:          60,
				},
				Visual: TUIVisualConfig{
					Theme:          createDarkTheme(),
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
					SubscriptionFilters: []string{},
					EnableBatching:      true,
					BatchSize:           50,
					BatchTimeout:        100 * time.Millisecond,
					ProcessingTimeout:   5 * time.Second,
				},
			},
			expectError: false,
		},
		{
			name: "invalid refresh rate - too fast",
			config: TUIConfig{
				Enabled: true,
				Performance: TUIPerformanceConfig{
					RefreshRate:        500 * time.Microsecond, // Too fast
					MaxLogLines:        10000,
					OptimizedRendering: true,
					TargetFPS:          60,
				},
				Visual: TUIVisualConfig{
					MinWidth:  80,
					MinHeight: 24,
				},
				Events: TUIEventConfig{
					BufferSize:        1000,
					BatchSize:         50,
					BatchTimeout:      100 * time.Millisecond,
					ProcessingTimeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
		{
			name: "invalid terminal size",
			config: TUIConfig{
				Enabled: true,
				Performance: TUIPerformanceConfig{
					RefreshRate:        16 * time.Millisecond,
					MaxLogLines:        10000,
					OptimizedRendering: true,
					TargetFPS:          60,
				},
				Visual: TUIVisualConfig{
					MinWidth:  30, // Too small
					MinHeight: 5,  // Too small
				},
				Events: TUIEventConfig{
					BufferSize:        1000,
					BatchSize:         50,
					BatchTimeout:      100 * time.Millisecond,
					ProcessingTimeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
		{
			name: "invalid color format",
			config: TUIConfig{
				Enabled: true,
				Performance: TUIPerformanceConfig{
					RefreshRate:        16 * time.Millisecond,
					MaxLogLines:        10000,
					OptimizedRendering: true,
					TargetFPS:          60,
				},
				Visual: TUIVisualConfig{
					Theme: TUIThemeConfig{
						Name:    "custom",
						Primary: "invalid-color", // Invalid hex color
					},
					MinWidth:  80,
					MinHeight: 24,
				},
				Events: TUIEventConfig{
					BufferSize:        1000,
					BatchSize:         50,
					BatchTimeout:      100 * time.Millisecond,
					ProcessingTimeout: 5 * time.Second,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{TUI: tt.config}
			err := cfg.validateTUIConfig()

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	// Set up environment variables
	envVars := map[string]string{
		"MORGANA_TUI_ENABLED":       "false",
		"MORGANA_TUI_REFRESH_RATE":  "33ms",
		"MORGANA_TUI_MAX_LOG_LINES": "5000",
		"MORGANA_TUI_BUFFER_SIZE":   "2000",
		"MORGANA_TUI_THEME":         "light",
		"MORGANA_TUI_COMPACT_MODE":  "true",
		"MORGANA_TUI_SHOW_DEBUG":    "true",
	}

	// Set environment variables
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	cfg := DefaultConfig()
	cfg.applyEnvOverrides()

	// Verify environment overrides were applied
	if cfg.TUI.Enabled {
		t.Error("expected TUI to be disabled via environment variable")
	}
	if cfg.TUI.Performance.RefreshRate != 33*time.Millisecond {
		t.Errorf("expected refresh rate 33ms, got %v", cfg.TUI.Performance.RefreshRate)
	}
	if cfg.TUI.Performance.MaxLogLines != 5000 {
		t.Errorf("expected max log lines 5000, got %d", cfg.TUI.Performance.MaxLogLines)
	}
	if cfg.TUI.Events.BufferSize != 2000 {
		t.Errorf("expected buffer size 2000, got %d", cfg.TUI.Events.BufferSize)
	}
	if cfg.TUI.Visual.Theme.Name != "light" {
		t.Errorf("expected theme 'light', got %s", cfg.TUI.Visual.Theme.Name)
	}
	if !cfg.TUI.Visual.CompactMode {
		t.Error("expected compact mode to be enabled")
	}
	if !cfg.TUI.Visual.ShowDebugInfo {
		t.Error("expected debug info to be enabled")
	}
}

func TestTUIConfigPresets(t *testing.T) {
	presets := CreateTUIConfigPresets()

	expectedPresets := []string{"performance", "development", "minimal"}
	for _, preset := range expectedPresets {
		if _, exists := presets[preset]; !exists {
			t.Errorf("expected preset %s to exist", preset)
		}
	}

	// Test performance preset characteristics
	perfConfig := presets["performance"]
	if !perfConfig.Enabled {
		t.Error("performance preset should be enabled")
	}
	if perfConfig.Performance.RefreshRate != 33*time.Millisecond {
		t.Error("performance preset should have 33ms refresh rate for lower CPU usage")
	}
	if perfConfig.Features.EnableFiltering {
		t.Error("performance preset should disable filtering for better performance")
	}

	// Test development preset characteristics
	devConfig := presets["development"]
	if !devConfig.Visual.ShowDebugInfo {
		t.Error("development preset should show debug info")
	}
	if !devConfig.Features.EnableExport {
		t.Error("development preset should enable export functionality")
	}

	// Test minimal preset characteristics
	minConfig := presets["minimal"]
	if minConfig.Performance.MaxLogLines != 1000 {
		t.Error("minimal preset should have low memory usage")
	}
	if minConfig.Features.EnableSearch {
		t.Error("minimal preset should disable search to keep things simple")
	}
}

func TestHexColorValidation(t *testing.T) {
	tests := []struct {
		color    string
		expected bool
	}{
		{"#FF0000", true},   // Valid hex
		{"#ff0000", true},   // Valid hex lowercase
		{"#123ABC", true},   // Valid hex mixed case
		{"", true},          // Empty is valid (uses default)
		{"FF0000", false},   // Missing #
		{"#FF00", false},    // Too short
		{"#FF00000", false}, // Too long
		{"#GGGGGG", false},  // Invalid hex characters
		{"red", false},      // Named color not supported
	}

	for _, tt := range tests {
		t.Run(tt.color, func(t *testing.T) {
			result := isValidHexColor(tt.color)
			if result != tt.expected {
				t.Errorf("isValidHexColor(%s) = %v, expected %v", tt.color, result, tt.expected)
			}
		})
	}
}

func TestTUIConfigConversion(t *testing.T) {
	cfg := DefaultConfig()

	// Test conversion to TUI config
	tuiConfig := cfg.ToTUIConfig()

	// Verify basic fields are converted correctly
	if tuiConfig.RefreshRate != cfg.TUI.Performance.RefreshRate {
		t.Errorf("refresh rate not converted correctly")
	}
	if tuiConfig.EventBufferSize != cfg.TUI.Events.BufferSize {
		t.Errorf("event buffer size not converted correctly")
	}
	if tuiConfig.MaxLogLines != cfg.TUI.Performance.MaxLogLines {
		t.Errorf("max log lines not converted correctly")
	}

	// Verify theme conversion
	if tuiConfig.Theme.Primary != cfg.TUI.Visual.Theme.Primary {
		t.Errorf("theme primary color not converted correctly")
	}

	// Verify agent colors are copied
	for agent, expectedColor := range cfg.TUI.Visual.Theme.AgentColors {
		if actualColor := tuiConfig.Theme.AgentColors[agent]; actualColor != expectedColor {
			t.Errorf("agent color for %s not converted correctly: expected %s, got %s", agent, expectedColor, actualColor)
		}
	}
}

func TestGetTUIAgentColor(t *testing.T) {
	cfg := DefaultConfig()

	// Test existing agent color
	color := cfg.GetTUIAgentColor("code-implementer")
	expectedColor := cfg.TUI.Visual.Theme.AgentColors["code-implementer"]
	if color != expectedColor {
		t.Errorf("expected color %s for code-implementer, got %s", expectedColor, color)
	}

	// Test non-existing agent color (should return default)
	unknownColor := cfg.GetTUIAgentColor("unknown-agent")
	expectedDefault := cfg.TUI.Visual.Theme.AgentColors["default"]
	if unknownColor != expectedDefault {
		t.Errorf("expected default color %s for unknown agent, got %s", expectedDefault, unknownColor)
	}
}
