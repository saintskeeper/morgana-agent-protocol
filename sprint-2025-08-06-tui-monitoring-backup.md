# Sprint Plan: TUI Monitoring System for Morgana Protocol

## Sprint Overview

- **Sprint ID**: sprint-2025-08-06-tui-monitoring
- **Duration**: 2 weeks
- **Goal**: Create configurable TUI interfaces that provide real-time monitoring
  of Morgana Director and subagent progress with intuitive progress
  visualization
- **Success Criteria**:
  - Users can monitor agent execution status in real-time via TUI
  - Progress tracking shows agent lifecycle, task queues, and completion metrics
  - Configuration system allows customization of TUI display preferences
  - Integration with existing telemetry system for accurate progress data

## Task Definitions

### Task: TUI_CONFIG_DESIGN

- **Title**: Design TUI configuration system architecture
- **Priority**: P0-Critical
- **Type**: architecture
- **Estimated Complexity**: medium
- **Dependencies**: none
- **Exit Criteria**:
  - [ ] TUI configuration schema defined in YAML format
  - [ ] Display modes specified (compact, detailed, dashboard)
  - [ ] Update frequency and refresh rate options documented
  - [ ] Color theme and styling configuration designed
  - [ ] Integration points with existing config system identified
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/internal/config/config.go
  - /Users/walterday/.claude/morgana.yaml
- **Notes**: Extend existing config.go structure to support TUI settings

### Task: TUI_COMPONENTS_DESIGN

- **Title**: Design TUI component architecture and layouts
- **Priority**: P0-Critical
- **Type**: architecture
- **Estimated Complexity**: complex
- **Dependencies**: [TUI_CONFIG_DESIGN]
- **Exit Criteria**:
  - [ ] Progress bar components designed for agent execution
  - [ ] Status dashboard layout defined with agent states
  - [ ] Real-time log streaming component specified
  - [ ] Task queue visualization component designed
  - [ ] Error/alert notification system defined
  - [ ] Keyboard navigation and interaction patterns documented
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/internal/adapter/adapter.go
  - /Users/walterday/.claude/morgana-protocol/internal/telemetry/telemetry.go
- **Notes**: Focus on readability and non-intrusive monitoring experience

### Task: PROGRESS_TRACKING_IMPL

- **Title**: Implement progress tracking data structures
- **Priority**: P0-Critical
- **Type**: implementation
- **Estimated Complexity**: medium
- **Dependencies**: [TUI_COMPONENTS_DESIGN]
- **Exit Criteria**:
  - [ ] Agent execution progress tracking structures created
  - [ ] Task queue state management implemented
  - [ ] Metrics collection for completion rates and timing
  - [ ] Thread-safe progress update mechanisms
  - [ ] Integration with existing telemetry spans completed
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/internal/adapter/adapter.go
  - /Users/walterday/.claude/morgana-protocol/internal/orchestrator/orchestrator.go
- **Notes**: Leverage existing telemetry attributes and spans for data source

### Task: TUI_RENDERER_IMPL

- **Title**: Implement TUI rendering engine
- **Priority**: P1-High
- **Type**: implementation
- **Estimated Complexity**: complex
- **Dependencies**: [PROGRESS_TRACKING_IMPL]
- **Exit Criteria**:
  - [ ] Terminal-based UI library integrated (bubbletea/tcell)
  - [ ] Progress bars render with accurate percentages
  - [ ] Status dashboard updates in real-time
  - [ ] Log streaming with scrollback buffer implemented
  - [ ] Responsive layout that adapts to terminal size
  - [ ] Color themes and styling system working
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/cmd/morgana/main.go
- **Notes**: Consider bubbletea for elegant terminal UI framework

### Task: TUI_CONFIG_IMPL

- **Title**: Implement TUI configuration loading and validation
- **Priority**: P1-High
- **Type**: implementation
- **Estimated Complexity**: medium
- **Dependencies**: [TUI_CONFIG_DESIGN, TUI_RENDERER_IMPL]
- **Exit Criteria**:
  - [ ] TUI config section added to main config structure
  - [ ] Configuration validation for TUI settings
  - [ ] Runtime configuration changes supported
  - [ ] Environment variable overrides for TUI settings
  - [ ] Default TUI configuration provides good UX out of box
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/internal/config/config.go
  - /Users/walterday/.claude/morgana.yaml
- **Notes**: Maintain backward compatibility with existing config system

### Task: MONITORING_INTEGRATION

- **Title**: Integrate monitoring with Morgana Director orchestration
- **Priority**: P1-High
- **Type**: implementation
- **Estimated Complexity**: complex
- **Dependencies**: [TUI_RENDERER_IMPL, PROGRESS_TRACKING_IMPL]
- **Exit Criteria**:
  - [ ] Real-time updates from orchestrator to TUI
  - [ ] Agent lifecycle events captured and displayed
  - [ ] Task queue status synchronized with TUI
  - [ ] Error propagation from agents to TUI alerts
  - [ ] Performance metrics displayed (execution time, throughput)
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/internal/orchestrator/orchestrator.go
  - /Users/walterday/.claude/morgana-protocol/internal/adapter/adapter.go
- **Notes**: Ensure minimal performance overhead on agent execution

### Task: TUI_DISPLAY_MODES

- **Title**: Implement multiple TUI display modes
- **Priority**: P2-Medium
- **Type**: implementation
- **Estimated Complexity**: medium
- **Dependencies**: [MONITORING_INTEGRATION]
- **Exit Criteria**:
  - [ ] Compact mode for minimal terminal space usage
  - [ ] Detailed mode with comprehensive agent information
  - [ ] Dashboard mode with multi-pane layout
  - [ ] Runtime switching between display modes
  - [ ] Mode-specific keyboard shortcuts implemented
- **Context Files**:
  - TUI renderer implementation files
- **Notes**: Allow users to toggle between modes based on needs and terminal
  size

### Task: TUI_TESTING

- **Title**: Create comprehensive test suite for TUI components
- **Priority**: P1-High
- **Type**: testing
- **Estimated Complexity**: medium
- **Dependencies**: [TUI_DISPLAY_MODES, TUI_CONFIG_IMPL]
- **Exit Criteria**:
  - [ ] Unit tests for progress tracking components ≥ 85%
  - [ ] Integration tests for TUI config loading
  - [ ] Mock agent execution scenarios for TUI testing
  - [ ] Terminal output validation in CI environment
  - [ ] Performance tests for real-time update handling
- **Context Files**:
  - Test files in morgana-protocol/internal/
- **Notes**: Focus on testable components, mock terminal interactions where
  needed

### Task: TUI_DOCUMENTATION

- **Title**: Create user documentation for TUI monitoring system
- **Priority**: P2-Medium
- **Type**: documentation
- **Estimated Complexity**: simple
- **Dependencies**: [TUI_TESTING]
- **Exit Criteria**:
  - [ ] Configuration reference documentation
  - [ ] User guide with screenshot examples
  - [ ] Keyboard shortcuts and navigation guide
  - [ ] Troubleshooting section for common issues
  - [ ] Integration examples with morgana.yaml
- **Context Files**:
  - /Users/walterday/.claude/morgana-protocol/README.md
- **Notes**: Include visual examples and configuration snippets

## Dependency Graph

```
TUI_CONFIG_DESIGN ──→ TUI_COMPONENTS_DESIGN ──→ PROGRESS_TRACKING_IMPL ──→ TUI_RENDERER_IMPL
                                                                              ↓
TUI_CONFIG_DESIGN ──→ TUI_CONFIG_IMPL ────────────────────────────────────→ MONITORING_INTEGRATION
                                                                              ↓
                                                                     TUI_DISPLAY_MODES
                                                                              ↓
                                                                         TUI_TESTING
                                                                              ↓
                                                                     TUI_DOCUMENTATION
```

## Validation Rules

- All P0 tasks must complete before P1 tasks begin
- TUI rendering must not impact agent execution performance
- Configuration changes require validation and backward compatibility
- All TUI components must be testable in headless environments

## Risk Mitigation

- **Risk**: Terminal compatibility issues across different environments

  - **Mitigation**: Test on multiple terminal emulators and CI environments
  - **Fallback**: Provide fallback text-only mode for incompatible terminals

- **Risk**: Performance overhead from real-time TUI updates

  - **Mitigation**: Implement configurable update intervals and efficient diff
    rendering
  - **Fallback**: Disable TUI monitoring for performance-critical deployments

- **Risk**: Complex UI state management affecting reliability
  - **Mitigation**: Keep UI state separate from core orchestrator logic
  - **Fallback**: TUI failures should not impact agent execution

## Technical Considerations

- **Library Choice**: Use bubbletea/charm for elegant terminal UI with good Go
  integration
- **Data Flow**: Leverage existing telemetry system to avoid duplicate
  monitoring code
- **Configuration**: Extend existing YAML config structure with new `tui`
  section
- **Performance**: Implement efficient diff-based rendering to minimize terminal
  updates
- **Accessibility**: Support standard terminal accessibility features and
  color-blind friendly themes
