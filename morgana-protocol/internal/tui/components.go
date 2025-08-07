package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Component interface for all UI components
type Component interface {
	Render(width, height int, theme Theme) string
	Update(msg interface{}) (Component, bool)
	Focus() Component
	Blur() Component
	IsFocused() bool
}

// ProgressBar component for displaying task progress
type ProgressBar struct {
	progress float64
	width    int
	label    string
	theme    Theme
	focused  bool
	animated bool
	animStep int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(width int, theme Theme) *ProgressBar {
	return &ProgressBar{
		width:    width,
		theme:    theme,
		animated: true,
	}
}

// SetProgress sets the progress value (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64, label string) {
	p.progress = progress
	p.label = label
}

// Render renders the progress bar
func (p *ProgressBar) Render(width, height int, theme Theme) string {
	if width != p.width {
		p.width = width
	}
	p.theme = theme

	// Calculate available width for progress bar
	labelWidth := len(p.label)
	percentWidth := 6                                       // " 100% "
	availableWidth := width - labelWidth - percentWidth - 2 // -2 for brackets

	if availableWidth < 10 {
		availableWidth = 10
	}

	// Calculate filled and empty portions
	filled := int(p.progress * float64(availableWidth))
	empty := availableWidth - filled

	// Create progress bar style
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ProgressBar.Complete))

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.ProgressBar.Incomplete))

	// Build progress bar
	var filledChar, emptyChar string
	if p.animated && p.progress > 0 && p.progress < 1 {
		// Animation characters for active progress
		animChars := []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
		filledChar = "█"
		emptyChar = animChars[p.animStep%len(animChars)]
	} else {
		filledChar = "█"
		emptyChar = "░"
	}

	progressBar := "[" +
		progressStyle.Render(strings.Repeat(filledChar, filled)) +
		emptyStyle.Render(strings.Repeat(emptyChar, empty)) +
		"]"

	// Format percentage
	percent := fmt.Sprintf("%3.0f%%", p.progress*100)

	// Combine label, progress bar, and percentage
	result := fmt.Sprintf("%s %s %s", p.label, progressBar, percent)

	// Apply focus styling if needed
	if p.focused {
		focusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Highlight)).
			Bold(true)
		result = focusStyle.Render(result)
	}

	return result
}

// Update handles updates to the progress bar
func (p *ProgressBar) Update(msg interface{}) (Component, bool) {
	updated := false

	switch m := msg.(type) {
	case TickMessage:
		if p.animated && p.progress > 0 && p.progress < 1 {
			p.animStep++
			updated = true
		}
	case float64:
		if m != p.progress {
			p.progress = m
			updated = true
		}
	}

	return p, updated
}

// Focus sets the component as focused
func (p *ProgressBar) Focus() Component {
	p.focused = true
	return p
}

// Blur removes focus from the component
func (p *ProgressBar) Blur() Component {
	p.focused = false
	return p
}

// IsFocused returns whether the component is focused
func (p *ProgressBar) IsFocused() bool {
	return p.focused
}

// StatusDashboard component for displaying multi-agent status
type StatusDashboard struct {
	taskStates map[string]*TaskState
	statusInfo *StatusInfo
	theme      Theme
	focused    bool
	width      int
	height     int
}

// NewStatusDashboard creates a new status dashboard
func NewStatusDashboard(theme Theme) *StatusDashboard {
	return &StatusDashboard{
		taskStates: make(map[string]*TaskState),
		theme:      theme,
	}
}

// SetData updates the dashboard data
func (s *StatusDashboard) SetData(taskStates map[string]*TaskState, statusInfo *StatusInfo) {
	s.taskStates = taskStates
	s.statusInfo = statusInfo
}

// Render renders the status dashboard
func (s *StatusDashboard) Render(width, height int, theme Theme) string {
	s.width = width
	s.height = height
	s.theme = theme

	var output strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 1)

	output.WriteString(headerStyle.Render("📊 Agent Status Dashboard"))
	output.WriteString("\n\n")

	// System status summary
	if s.statusInfo != nil {
		summaryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground)).
			Padding(0, 1)

		summary := fmt.Sprintf("Active: %d | Completed: %d | Failed: %d | Events: %d",
			s.statusInfo.ActiveTasks,
			s.statusInfo.CompletedTasks,
			s.statusInfo.FailedTasks,
			s.statusInfo.TotalEvents,
		)
		output.WriteString(summaryStyle.Render(summary))
		output.WriteString("\n\n")
	}

	// Task status cards
	if len(s.taskStates) == 0 {
		noTasksStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true).
			Padding(0, 1)
		output.WriteString(noTasksStyle.Render("No active tasks"))
	} else {
		// Render task cards in a grid
		cardWidth := (width - 4) / 2 // 2 columns with padding
		if cardWidth < 30 {
			cardWidth = width - 2 // Single column if too narrow
		}

		cardCount := 0
		for _, task := range s.taskStates {
			if cardCount > 0 && cardCount%2 == 0 {
				output.WriteString("\n")
			}

			card := s.renderTaskCard(task, cardWidth, theme)
			output.WriteString(card)

			// Add spacing between columns
			if cardCount%2 == 0 && len(s.taskStates) > 1 {
				output.WriteString("  ")
			}

			cardCount++
		}
	}

	// Apply focus styling
	result := output.String()
	if s.focused {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Primary)).
			Padding(0, 1)
		result = borderStyle.Render(result)
	}

	return result
}

// renderTaskCard renders an individual task status card
func (s *StatusDashboard) renderTaskCard(task *TaskState, width int, theme Theme) string {
	var cardStyle lipgloss.Style
	var statusColor string

	// Choose color based on status
	switch task.Status {
	case StatusRunning:
		statusColor = theme.Info
	case StatusCompleted:
		statusColor = theme.Success
	case StatusFailed:
		statusColor = theme.Error
	case StatusPending:
		statusColor = theme.Warning
	default:
		statusColor = theme.Muted
	}

	cardStyle = lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(statusColor)).
		Padding(1)

	var content strings.Builder

	// Agent type header
	agentColor, exists := theme.AgentColors[task.AgentType]
	if !exists {
		agentColor = theme.AgentColors["default"]
	}

	agentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(agentColor)).
		Bold(true)
	content.WriteString(agentStyle.Render(task.AgentType))
	content.WriteString("\n")

	// Task ID
	idStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted)).
		Italic(true)
	content.WriteString(idStyle.Render(fmt.Sprintf("ID: %s", truncateString(task.ID, 12))))
	content.WriteString("\n")

	// Status and stage
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(statusColor)).
		Bold(true)
	content.WriteString(statusStyle.Render(fmt.Sprintf("Status: %s", s.getStatusString(task.Status))))
	content.WriteString("\n")

	if task.Stage != "" {
		stageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground))
		content.WriteString(stageStyle.Render(fmt.Sprintf("Stage: %s", task.Stage)))
		content.WriteString("\n")
	}

	// Progress bar for running tasks
	if task.Status == StatusRunning && task.Progress > 0 {
		progressBar := NewProgressBar(width-4, theme)
		progressBar.SetProgress(task.Progress, "")
		content.WriteString(progressBar.Render(width-4, 1, theme))
		content.WriteString("\n")
	}

	// Duration
	if task.Duration > 0 {
		durationStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted))
		content.WriteString(durationStyle.Render(fmt.Sprintf("Duration: %v", task.Duration.Round(time.Millisecond))))
		content.WriteString("\n")
	}

	// Message or error
	if task.Error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Error)).
			Italic(true)
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", truncateString(task.Error, 40))))
	} else if task.Message != "" {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground))
		content.WriteString(messageStyle.Render(truncateString(task.Message, 40)))
	}

	return cardStyle.Render(content.String())
}

// getStatusString returns a human-readable status string
func (s *StatusDashboard) getStatusString(status TaskStatus) string {
	switch status {
	case StatusPending:
		return "Pending"
	case StatusRunning:
		return "Running"
	case StatusCompleted:
		return "Completed"
	case StatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// Update handles updates to the dashboard
func (s *StatusDashboard) Update(msg interface{}) (Component, bool) {
	// Dashboard doesn't handle direct updates, data is set via SetData
	return s, false
}

// Focus sets the component as focused
func (s *StatusDashboard) Focus() Component {
	s.focused = true
	return s
}

// Blur removes focus from the component
func (s *StatusDashboard) Blur() Component {
	s.focused = false
	return s
}

// IsFocused returns whether the component is focused
func (s *StatusDashboard) IsFocused() bool {
	return s.focused
}

// LogViewer component for displaying scrollable logs
type LogViewer struct {
	logEntries    []*LogEntry
	theme         Theme
	focused       bool
	width         int
	height        int
	scrollOffset  int
	filter        string
	searchQuery   string
	searchIndex   int
	showTimestamp bool
}

// NewLogViewer creates a new log viewer
func NewLogViewer(theme Theme, showTimestamp bool) *LogViewer {
	return &LogViewer{
		theme:         theme,
		showTimestamp: showTimestamp,
		logEntries:    make([]*LogEntry, 0),
	}
}

// SetLogs updates the log entries
func (l *LogViewer) SetLogs(entries []*LogEntry) {
	l.logEntries = entries
}

// SetFilter sets the current filter
func (l *LogViewer) SetFilter(filter string) {
	l.filter = filter
	l.scrollOffset = 0 // Reset scroll when filter changes
}

// SetSearch sets the current search query
func (l *LogViewer) SetSearch(query string, index int) {
	l.searchQuery = query
	l.searchIndex = index
}

// Render renders the log viewer
func (l *LogViewer) Render(width, height int, theme Theme) string {
	l.width = width
	l.height = height
	l.theme = theme

	var output strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 1)

	header := "📝 Task Logs"
	if l.filter != "" {
		header += fmt.Sprintf(" (filtered: %s)", l.filter)
	}
	if l.searchQuery != "" {
		header += fmt.Sprintf(" (search: %s)", l.searchQuery)
	}

	output.WriteString(headerStyle.Render(header))
	output.WriteString("\n")

	// Filter and display log entries
	filteredLogs := l.getFilteredLogs()
	visibleHeight := height - 2 // Account for header and border

	if len(filteredLogs) == 0 {
		noLogsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true).
			Padding(1)
		output.WriteString(noLogsStyle.Render("No logs available"))
	} else {
		// Calculate visible range
		start := l.scrollOffset
		end := start + visibleHeight
		if end > len(filteredLogs) {
			end = len(filteredLogs)
		}
		if start >= len(filteredLogs) {
			start = len(filteredLogs) - visibleHeight
			if start < 0 {
				start = 0
			}
		}

		// Render visible log entries
		for i := start; i < end; i++ {
			entry := filteredLogs[i]
			logLine := l.renderLogEntry(entry, width-2, theme)
			output.WriteString(logLine)
			if i < end-1 {
				output.WriteString("\n")
			}
		}

		// Add scroll indicators
		if len(filteredLogs) > visibleHeight {
			scrollInfo := fmt.Sprintf(" (%d-%d of %d)", start+1, end, len(filteredLogs))
			scrollStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Muted)).
				Italic(true)
			output.WriteString("\n")
			output.WriteString(scrollStyle.Render(scrollInfo))
		}
	}

	// Apply focus styling
	result := output.String()
	if l.focused {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Primary)).
			Padding(0, 1)
		result = borderStyle.Render(result)
	}

	return result
}

// renderLogEntry renders a single log entry
func (l *LogViewer) renderLogEntry(entry *LogEntry, width int, theme Theme) string {
	var parts []string

	// Timestamp
	if l.showTimestamp {
		timestampStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted))
		timestamp := entry.Timestamp.Format("15:04:05.000")
		parts = append(parts, timestampStyle.Render(timestamp))
	}

	// Log level
	levelColor := l.getLogLevelColor(entry.Level, theme)
	levelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(levelColor)).
		Bold(true)
	levelStr := l.getLogLevelString(entry.Level)
	parts = append(parts, levelStyle.Render(fmt.Sprintf("[%s]", levelStr)))

	// Agent type (if available)
	if entry.AgentType != "" {
		agentColor, exists := theme.AgentColors[entry.AgentType]
		if !exists {
			agentColor = theme.AgentColors["default"]
		}
		agentStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(agentColor))
		parts = append(parts, agentStyle.Render(fmt.Sprintf("[%s]", entry.AgentType)))
	}

	// Stage (if available)
	if entry.Stage != "" && entry.Stage != entry.AgentType {
		stageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Secondary))
		parts = append(parts, stageStyle.Render(fmt.Sprintf("[%s]", entry.Stage)))
	}

	// Message
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground))
	message := entry.Message

	// Highlight search terms
	if l.searchQuery != "" && strings.Contains(strings.ToLower(message), strings.ToLower(l.searchQuery)) {
		highlightStyle := lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Warning)).
			Foreground(lipgloss.Color(theme.Background))
		// Simple highlight - in production might use more sophisticated highlighting
		message = strings.ReplaceAll(message, l.searchQuery, highlightStyle.Render(l.searchQuery))
	}

	parts = append(parts, messageStyle.Render(message))

	// Join parts with spaces and truncate if necessary
	line := strings.Join(parts, " ")
	if len(line) > width {
		line = line[:width-3] + "..."
	}

	return line
}

// getFilteredLogs returns logs filtered by current filter and search
func (l *LogViewer) getFilteredLogs() []*LogEntry {
	if l.filter == "" && l.searchQuery == "" {
		return l.logEntries
	}

	var filtered []*LogEntry
	for _, entry := range l.logEntries {
		// Apply filter
		if l.filter != "" {
			switch l.filter {
			case "error":
				if entry.Level != LogLevelError {
					continue
				}
			case "warning":
				if entry.Level != LogLevelWarn {
					continue
				}
			case "info":
				if entry.Level != LogLevelInfo {
					continue
				}
			case "debug":
				if entry.Level != LogLevelDebug {
					continue
				}
			default:
				// Filter by agent type or stage
				if entry.AgentType != l.filter && entry.Stage != l.filter {
					continue
				}
			}
		}

		// Apply search
		if l.searchQuery != "" {
			queryLower := strings.ToLower(l.searchQuery)
			if !strings.Contains(strings.ToLower(entry.Message), queryLower) &&
				!strings.Contains(strings.ToLower(entry.AgentType), queryLower) &&
				!strings.Contains(strings.ToLower(entry.Stage), queryLower) {
				continue
			}
		}

		filtered = append(filtered, entry)
	}

	return filtered
}

// getLogLevelColor returns the color for a log level
func (l *LogViewer) getLogLevelColor(level LogLevel, theme Theme) string {
	switch level {
	case LogLevelError:
		return theme.Error
	case LogLevelWarn:
		return theme.Warning
	case LogLevelInfo:
		return theme.Info
	case LogLevelDebug:
		return theme.Muted
	default:
		return theme.Foreground
	}
}

// getLogLevelString returns the string representation of a log level
func (l *LogViewer) getLogLevelString(level LogLevel) string {
	switch level {
	case LogLevelError:
		return "ERR"
	case LogLevelWarn:
		return "WRN"
	case LogLevelInfo:
		return "INF"
	case LogLevelDebug:
		return "DBG"
	default:
		return "UNK"
	}
}

// ScrollUp scrolls the log viewer up
func (l *LogViewer) ScrollUp(lines int) {
	l.scrollOffset -= lines
	if l.scrollOffset < 0 {
		l.scrollOffset = 0
	}
}

// ScrollDown scrolls the log viewer down
func (l *LogViewer) ScrollDown(lines int) {
	filtered := l.getFilteredLogs()
	maxOffset := len(filtered) - (l.height - 2) // Account for header
	if maxOffset < 0 {
		maxOffset = 0
	}

	l.scrollOffset += lines
	if l.scrollOffset > maxOffset {
		l.scrollOffset = maxOffset
	}
}

// Update handles updates to the log viewer
func (l *LogViewer) Update(msg interface{}) (Component, bool) {
	return l, false
}

// Focus sets the component as focused
func (l *LogViewer) Focus() Component {
	l.focused = true
	return l
}

// Blur removes focus from the component
func (l *LogViewer) Blur() Component {
	l.focused = false
	return l
}

// IsFocused returns whether the component is focused
func (l *LogViewer) IsFocused() bool {
	return l.focused
}

// Utility functions

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
}
