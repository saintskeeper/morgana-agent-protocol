package tui

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// Model represents the main TUI model for bubbletea
type Model struct {
	// Configuration and context
	config    TUIConfig
	ctx       context.Context
	eventBus  events.EventBus
	bridge    *EventBridge
	processor *EventProcessor

	// State management
	state *ModelState

	// UI Components
	dashboard  *StatusDashboard
	logViewer  *LogViewer
	statsPanel *StatisticsPanel
	statistics *ExecutionStatistics

	// Layout and interaction
	focused  ComponentType
	layout   LayoutType
	keybinds *KeyBindings

	// Performance tracking
	startTime     time.Time
	lastRender    time.Time
	renderCounter int64

	// Shutdown management
	shutdownRequested bool
}

// LayoutType represents different layout configurations
type LayoutType int

const (
	LayoutSplit      LayoutType = iota // Dashboard top, logs bottom
	LayoutDashboard                    // Dashboard only
	LayoutLogs                         // Logs only
	LayoutStatistics                   // Statistics only
	LayoutHelp                         // Help screen
)

// KeyBindings defines keyboard shortcuts
type KeyBindings struct {
	Quit         []string
	SwitchFocus  []string
	ScrollUp     []string
	ScrollDown   []string
	ToggleLayout []string
	ShowHelp     []string
	ShowStats    []string
	Filter       []string
	Search       []string
	ClearFilter  []string
}

// DefaultKeyBindings returns default key bindings
func DefaultKeyBindings() *KeyBindings {
	return &KeyBindings{
		Quit:         []string{"q", "ctrl+c", "esc"},
		SwitchFocus:  []string{"tab", "shift+tab"},
		ScrollUp:     []string{"up", "k"},
		ScrollDown:   []string{"down", "j"},
		ToggleLayout: []string{"l", "space"},
		ShowHelp:     []string{"h", "?"},
		ShowStats:    []string{"s"},
		Filter:       []string{"f"},
		Search:       []string{"/"},
		ClearFilter:  []string{"c"},
	}
}

// NewModel creates a new TUI model
func NewModel(ctx context.Context, eventBus events.EventBus, config TUIConfig) *Model {
	// Initialize statistics tracking
	statistics := NewExecutionStatistics(100) // Keep 100 history snapshots

	// Initialize components
	dashboard := NewStatusDashboard(config.Theme)
	logViewer := NewLogViewer(config.Theme, config.ShowTimeStamps)
	statsPanel := NewStatisticsPanel(statistics, config.Theme)

	// Initialize state
	state := &ModelState{
		TaskStates: make(map[string]*TaskState),
		LogEntries: make([]*LogEntry, 0),
		StatusInfo: &StatusInfo{
			SystemLoad: SystemLoad{},
		},
		Focused: ComponentDashboard,
	}

	// Set initial focus
	dashboard.Focus()

	model := &Model{
		config:     config,
		ctx:        ctx,
		eventBus:   eventBus,
		processor:  NewEventProcessor(config),
		state:      state,
		dashboard:  dashboard,
		logViewer:  logViewer,
		statsPanel: statsPanel,
		statistics: statistics,
		focused:    ComponentDashboard,
		layout:     LayoutSplit,
		keybinds:   DefaultKeyBindings(),
		startTime:  time.Now(),
	}

	// Create and start event bridge
	model.bridge = NewEventBridge(eventBus, config)

	return model
}

// Init implements the tea.Model interface
func (m *Model) Init() tea.Cmd {
	// Start the event bridge
	return func() tea.Msg {
		// This will be called when the program starts
		return EventMessage{
			Event:     nil, // Special message to indicate initialization
			Timestamp: time.Now(),
		}
	}
}

// Update implements the tea.Model interface
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.state.Width = msg.Width
		m.state.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle keyboard input
		return m.handleKeyPress(msg)

	case EventMessage:
		// Handle events from the bridge
		if msg.Event != nil {
			m.processEvent(msg)
		} else {
			// Initialization message - start the bridge
			return m, m.startBridge()
		}

	case TickMessage:
		// Handle periodic updates for FPS control
		m.updatePerformanceMetrics()

		// Update components that need periodic refreshing
		if dashboard, updated := m.dashboard.Update(msg); updated {
			m.dashboard = dashboard.(*StatusDashboard)
		}
		if logViewer, updated := m.logViewer.Update(msg); updated {
			m.logViewer = logViewer.(*LogViewer)
		}
		if statsPanel, updated := m.statsPanel.Update(msg); updated {
			m.statsPanel = statsPanel.(*StatisticsPanel)
		}

		// Take periodic statistics snapshots (every 10 seconds)
		elapsed := time.Since(m.startTime)
		if int(elapsed.Seconds())%10 == 0 && elapsed.Milliseconds()%1000 < 100 {
			m.statistics.TakeSnapshot()
		}

	case tea.QuitMsg:
		// Handle quit message
		return m, m.shutdown()
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View implements the tea.Model interface
func (m *Model) View() string {
	if m.state.Width == 0 || m.state.Height == 0 {
		return "Initializing TUI..."
	}

	m.lastRender = time.Now()
	m.renderCounter++

	switch m.layout {
	case LayoutDashboard:
		return m.renderDashboardOnly()
	case LayoutLogs:
		return m.renderLogsOnly()
	case LayoutStatistics:
		return m.renderStatisticsOnly()
	case LayoutHelp:
		return m.renderHelp()
	default: // LayoutSplit
		return m.renderSplitLayout()
	}
}

// startBridge starts the event bridge
func (m *Model) startBridge() tea.Cmd {
	return func() tea.Msg {
		// This is a bit of a hack - we need the program instance to start the bridge
		// In real usage, this would be handled differently
		return nil
	}
}

// SetProgram sets the tea program for the bridge
func (m *Model) SetProgram(program *tea.Program) error {
	return m.bridge.Start(program)
}

// handleKeyPress processes keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Check for quit keys
	for _, quitKey := range m.keybinds.Quit {
		if key == quitKey {
			m.shutdownRequested = true
			return m, tea.Quit
		}
	}

	// Check for help keys
	for _, helpKey := range m.keybinds.ShowHelp {
		if key == helpKey {
			if m.layout == LayoutHelp {
				m.layout = LayoutSplit
			} else {
				m.layout = LayoutHelp
			}
			return m, nil
		}
	}

	// Check for statistics keys
	for _, statsKey := range m.keybinds.ShowStats {
		if key == statsKey {
			if m.layout == LayoutStatistics {
				m.layout = LayoutSplit
			} else {
				m.layout = LayoutStatistics
				m.focused = ComponentStatistics
				m.statsPanel.Focus()
				m.dashboard.Blur()
				m.logViewer.Blur()
			}
			return m, nil
		}
	}

	// Check for layout toggle
	for _, layoutKey := range m.keybinds.ToggleLayout {
		if key == layoutKey {
			m.toggleLayout()
			return m, nil
		}
	}

	// If in help mode, only handle quit and help keys
	if m.layout == LayoutHelp {
		return m, nil
	}

	// Check for focus switching
	for _, focusKey := range m.keybinds.SwitchFocus {
		if key == focusKey {
			m.switchFocus()
			return m, nil
		}
	}

	// Handle component-specific keys based on focus
	switch m.focused {
	case ComponentLogs:
		return m.handleLogViewerKeys(key)
	case ComponentDashboard:
		return m.handleDashboardKeys(key)
	case ComponentStatistics:
		return m.handleStatisticsKeys(key)
	}

	return m, nil
}

// handleLogViewerKeys handles keys specific to the log viewer
func (m *Model) handleLogViewerKeys(key string) (tea.Model, tea.Cmd) {
	// Scroll up
	for _, scrollKey := range m.keybinds.ScrollUp {
		if key == scrollKey {
			m.logViewer.ScrollUp(1)
			return m, nil
		}
	}

	// Scroll down
	for _, scrollKey := range m.keybinds.ScrollDown {
		if key == scrollKey {
			m.logViewer.ScrollDown(1)
			return m, nil
		}
	}

	// Page up/down
	if key == "pgup" {
		m.logViewer.ScrollUp(10)
		return m, nil
	}
	if key == "pgdown" {
		m.logViewer.ScrollDown(10)
		return m, nil
	}

	// Home/End
	if key == "home" {
		m.logViewer.scrollOffset = 0
		return m, nil
	}
	if key == "end" {
		// Scroll to bottom
		filtered := m.logViewer.getFilteredLogs()
		m.logViewer.scrollOffset = len(filtered)
		return m, nil
	}

	// Filter keys
	for _, filterKey := range m.keybinds.Filter {
		if key == filterKey {
			// In a real implementation, this would open a filter dialog
			// For now, cycle through common filters
			m.cycleFilter()
			return m, nil
		}
	}

	// Clear filter
	for _, clearKey := range m.keybinds.ClearFilter {
		if key == clearKey {
			m.logViewer.SetFilter("")
			m.state.CurrentFilter = ""
			return m, nil
		}
	}

	return m, nil
}

// handleDashboardKeys handles keys specific to the dashboard
func (m *Model) handleDashboardKeys(key string) (tea.Model, tea.Cmd) {
	// Dashboard might have specific key handlers in the future
	return m, nil
}

// handleStatisticsKeys handles keys specific to the statistics panel
func (m *Model) handleStatisticsKeys(key string) (tea.Model, tea.Cmd) {
	// Scroll up
	for _, scrollKey := range m.keybinds.ScrollUp {
		if key == scrollKey {
			m.statsPanel.ScrollUp(1)
			return m, nil
		}
	}

	// Scroll down
	for _, scrollKey := range m.keybinds.ScrollDown {
		if key == scrollKey {
			m.statsPanel.ScrollDown(1)
			return m, nil
		}
	}

	// Page up/down
	if key == "pgup" {
		m.statsPanel.ScrollUp(10)
		return m, nil
	}
	if key == "pgdown" {
		m.statsPanel.ScrollDown(10)
		return m, nil
	}

	// Pass key to statistics panel for internal navigation
	if updated, changed := m.statsPanel.Update(key); changed {
		m.statsPanel = updated.(*StatisticsPanel)
		return m, nil
	}

	return m, nil
}

// processEvent processes events from the event bridge
func (m *Model) processEvent(eventMsg EventMessage) {
	// Update state using the event processor
	m.processor.ProcessEvent(eventMsg.Event, m.state)

	// Update statistics tracking
	if eventMsg.Event != nil {
		m.statistics.ProcessEvent(eventMsg.Event)
		// Update statistics panel with new data
		stats, agentStats := m.statistics.GetStatistics()
		m.statsPanel.SetData(stats, agentStats)
	}

	// Update component data
	m.dashboard.SetData(m.state.TaskStates, m.state.StatusInfo)
	m.logViewer.SetLogs(m.state.LogEntries)

	// Update system load metrics
	m.updateSystemMetrics()
}

// updatePerformanceMetrics updates performance tracking
func (m *Model) updatePerformanceMetrics() {
	if m.state.StatusInfo == nil {
		return
	}

	now := time.Now()

	// Calculate FPS
	if !m.lastRender.IsZero() {
		frameDuration := now.Sub(m.lastRender)
		if frameDuration > 0 {
			fps := 1.0 / frameDuration.Seconds()
			m.state.StatusInfo.SystemLoad.RenderFPS = fps
		}
	}

	// Update uptime
	m.state.StatusInfo.Uptime = now.Sub(m.startTime)

	// Get memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.state.StatusInfo.SystemLoad.MemoryMB = float64(memStats.Alloc) / 1024 / 1024

	// Update event counts
	m.state.StatusInfo.TotalEvents = m.state.EventCount

	// Get bridge stats for event queue length
	if bridgeStats := m.bridge.GetStats(); bridgeStats.IsRunning {
		m.state.StatusInfo.SystemLoad.EventsPerSec = bridgeStats.FPS
	}
}

// updateSystemMetrics updates system performance metrics
func (m *Model) updateSystemMetrics() {
	// In a full implementation, this would collect real system metrics
	// For now, we'll use placeholder calculations

	if m.state.StatusInfo == nil {
		return
	}

	// Calculate CPU usage (placeholder)
	m.state.StatusInfo.SystemLoad.CPUPercent = 0.5 + (float64(m.renderCounter%100) / 200.0)

	// Event processing rate
	eventStats := m.eventBus.Stats()
	if eventStats.QueueSize > 0 {
		m.state.StatusInfo.SystemLoad.EventQueueLen = eventStats.QueueSize
	}
}

// switchFocus switches focus between components
func (m *Model) switchFocus() {
	// Blur current component
	switch m.focused {
	case ComponentDashboard:
		m.dashboard.Blur()
	case ComponentLogs:
		m.logViewer.Blur()
	case ComponentStatistics:
		m.statsPanel.Blur()
	}

	// Switch to next component
	switch m.layout {
	case LayoutSplit:
		if m.focused == ComponentDashboard {
			m.focused = ComponentLogs
			m.logViewer.Focus()
		} else {
			m.focused = ComponentDashboard
			m.dashboard.Focus()
		}
	case LayoutDashboard:
		// Only dashboard, no switching
	case LayoutLogs:
		// Only logs, no switching
	case LayoutStatistics:
		// Only statistics, no switching
	}
}

// toggleLayout cycles through different layout modes
func (m *Model) toggleLayout() {
	switch m.layout {
	case LayoutSplit:
		m.layout = LayoutDashboard
		m.focused = ComponentDashboard
		m.dashboard.Focus()
		m.logViewer.Blur()
		m.statsPanel.Blur()
	case LayoutDashboard:
		m.layout = LayoutLogs
		m.focused = ComponentLogs
		m.logViewer.Focus()
		m.dashboard.Blur()
		m.statsPanel.Blur()
	case LayoutLogs:
		m.layout = LayoutStatistics
		m.focused = ComponentStatistics
		m.statsPanel.Focus()
		m.dashboard.Blur()
		m.logViewer.Blur()
	case LayoutStatistics:
		m.layout = LayoutSplit
		m.focused = ComponentDashboard
		m.dashboard.Focus()
		m.logViewer.Blur()
		m.statsPanel.Blur()
	}
}

// cycleFilter cycles through common log filters
func (m *Model) cycleFilter() {
	filters := []string{"", "error", "warning", "info", "code-implementer", "sprint-planner", "test-specialist", "validation-expert"}

	currentIndex := 0
	for i, filter := range filters {
		if filter == m.state.CurrentFilter {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(filters)
	nextFilter := filters[nextIndex]

	m.state.CurrentFilter = nextFilter
	m.logViewer.SetFilter(nextFilter)
}

// renderSplitLayout renders the split dashboard and logs layout
func (m *Model) renderSplitLayout() string {
	width := m.state.Width
	height := m.state.Height

	// Calculate component heights
	headerHeight := 1
	statusBarHeight := 1
	availableHeight := height - headerHeight - statusBarHeight

	dashboardHeight := availableHeight / 2
	logsHeight := availableHeight - dashboardHeight

	// Render header
	header := m.renderHeader(width)

	// Render components
	dashboardView := m.dashboard.Render(width, dashboardHeight, m.config.Theme)
	logsView := m.logViewer.Render(width, logsHeight, m.config.Theme)

	// Render status bar
	statusBar := m.renderStatusBar(width)

	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		dashboardView,
		logsView,
		statusBar,
	)
}

// renderDashboardOnly renders only the dashboard
func (m *Model) renderDashboardOnly() string {
	width := m.state.Width
	height := m.state.Height

	headerHeight := 1
	statusBarHeight := 1
	availableHeight := height - headerHeight - statusBarHeight

	header := m.renderHeader(width)
	dashboardView := m.dashboard.Render(width, availableHeight, m.config.Theme)
	statusBar := m.renderStatusBar(width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		dashboardView,
		statusBar,
	)
}

// renderLogsOnly renders only the logs
func (m *Model) renderLogsOnly() string {
	width := m.state.Width
	height := m.state.Height

	headerHeight := 1
	statusBarHeight := 1
	availableHeight := height - headerHeight - statusBarHeight

	header := m.renderHeader(width)
	logsView := m.logViewer.Render(width, availableHeight, m.config.Theme)
	statusBar := m.renderStatusBar(width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		logsView,
		statusBar,
	)
}

// renderStatisticsOnly renders only the statistics
func (m *Model) renderStatisticsOnly() string {
	width := m.state.Width
	height := m.state.Height

	headerHeight := 1
	statusBarHeight := 1
	availableHeight := height - headerHeight - statusBarHeight

	header := m.renderHeader(width)
	statsView := m.statsPanel.Render(width, availableHeight, m.config.Theme)
	statusBar := m.renderStatusBar(width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		statsView,
		statusBar,
	)
}

// renderHelp renders the help screen
func (m *Model) renderHelp() string {
	width := m.state.Width
	height := m.state.Height

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.config.Theme.Foreground)).
		Padding(2)

	var help strings.Builder
	help.WriteString("🧙‍♂️ Morgana Protocol TUI - Keyboard Shortcuts\n\n")

	help.WriteString("Navigation:\n")
	help.WriteString("  Tab/Shift+Tab  - Switch between components\n")
	help.WriteString("  Space/L        - Toggle layout mode\n")
	help.WriteString("  S              - Show/hide statistics view\n")
	help.WriteString("  H/?            - Show/hide this help\n")
	help.WriteString("  Q/Ctrl+C/Esc   - Quit application\n\n")

	help.WriteString("Log Viewer (when focused):\n")
	help.WriteString("  ↑/K, ↓/J       - Scroll up/down one line\n")
	help.WriteString("  PgUp/PgDn      - Scroll up/down one page\n")
	help.WriteString("  Home/End       - Go to top/bottom\n")
	help.WriteString("  F              - Cycle through filters\n")
	help.WriteString("  C              - Clear current filter\n\n")

	help.WriteString("Statistics (when focused):\n")
	help.WriteString("  ↑/K, ↓/J       - Scroll up/down one line\n")
	help.WriteString("  PgUp/PgDn      - Scroll up/down one page\n")
	help.WriteString("  S              - Toggle detailed view\n")
	help.WriteString("  B              - Back to overview (in detail)\n\n")

	help.WriteString("Layout Modes:\n")
	help.WriteString("  Split          - Dashboard + Logs (default)\n")
	help.WriteString("  Dashboard Only - Agent status cards only\n")
	help.WriteString("  Logs Only      - Log viewer only\n")
	help.WriteString("  Statistics     - Execution statistics view\n\n")

	help.WriteString("Features:\n")
	help.WriteString("  • Real-time agent execution tracking\n")
	help.WriteString("  • Color-coded status indicators\n")
	help.WriteString("  • Animated progress bars\n")
	help.WriteString("  • Log filtering and search\n")
	help.WriteString("  • Performance monitoring\n")
	help.WriteString("  • Execution statistics and metrics\n")
	help.WriteString("  • Success/failure rate tracking\n")
	help.WriteString("  • Performance trend analysis\n")
	help.WriteString("  • Responsive terminal resizing\n\n")

	help.WriteString("Press H or ? to return to the main view")

	content := helpStyle.Render(help.String())

	// Center the help content
	if width > 0 && height > 0 {
		style := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center)
		content = style.Render(content)
	}

	return content
}

// renderHeader renders the application header with proper width
func (m *Model) renderHeader(width int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Foreground(lipgloss.Color(m.config.Theme.Primary)).
		Background(lipgloss.Color(m.config.Theme.Background)).
		Bold(true).
		Padding(0, 1)

	title := "🧙‍♂️ Morgana Protocol TUI"
	layoutInfo := fmt.Sprintf("[%s]", m.getLayoutString())
	focusInfo := fmt.Sprintf("Focus: %s", m.getFocusString())

	// Performance info if debug enabled
	var perfInfo string
	if m.config.ShowDebugInfo && m.state.StatusInfo != nil {
		fps := m.state.StatusInfo.SystemLoad.RenderFPS
		perfInfo = fmt.Sprintf(" | %.1f FPS", fps)
	}

	header := fmt.Sprintf("%s %s | %s%s", title, layoutInfo, focusInfo, perfInfo)
	return style.Render(header)
}

// renderStatusBar renders the bottom status bar
func (m *Model) renderStatusBar(width int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Foreground(lipgloss.Color(m.config.Theme.Muted)).
		Background(lipgloss.Color(m.config.Theme.Background)).
		Padding(0, 1)

	var parts []string

	// Runtime info
	uptime := ""
	if m.state.StatusInfo != nil {
		uptime = m.state.StatusInfo.Uptime.Round(time.Second).String()
	}
	parts = append(parts, fmt.Sprintf("Uptime: %s", uptime))

	// Event count
	events := int64(0)
	if m.state.StatusInfo != nil {
		events = m.state.StatusInfo.TotalEvents
	}
	parts = append(parts, fmt.Sprintf("Events: %d", events))

	// Help hint
	parts = append(parts, "Press H for help, Q to quit")

	statusText := strings.Join(parts, " | ")
	return style.Render(statusText)
}

// getLayoutString returns a string representation of the current layout
func (m *Model) getLayoutString() string {
	switch m.layout {
	case LayoutSplit:
		return "Split"
	case LayoutDashboard:
		return "Dashboard"
	case LayoutLogs:
		return "Logs"
	case LayoutStatistics:
		return "Statistics"
	case LayoutHelp:
		return "Help"
	default:
		return "Unknown"
	}
}

// getFocusString returns a string representation of the current focus
func (m *Model) getFocusString() string {
	switch m.focused {
	case ComponentDashboard:
		return "Dashboard"
	case ComponentLogs:
		return "Logs"
	case ComponentStatus:
		return "Status"
	case ComponentStatistics:
		return "Statistics"
	case ComponentHelp:
		return "Help"
	default:
		return "Unknown"
	}
}

// shutdown performs cleanup when shutting down
func (m *Model) shutdown() tea.Cmd {
	if m.bridge != nil {
		m.bridge.Stop()
	}
	return tea.Quit
}
