package tui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// ExecutionStatistics tracks comprehensive performance metrics and statistics
type ExecutionStatistics struct {
	// Agent execution statistics
	AgentStats map[string]*AgentStatistics `json:"agent_stats"`

	// Session-wide metrics
	SessionStats *SessionStatistics `json:"session_stats"`

	// Performance metrics
	PerformanceMetrics *PerformanceMetrics `json:"performance_metrics"`

	// Historical tracking for trends
	HistoryWindow []*StatSnapshot `json:"history_window"`

	// Thread safety
	mutex sync.RWMutex

	// Configuration
	maxHistorySize int
	sessionStart   time.Time
}

// AgentStatistics tracks statistics for a specific agent type
type AgentStatistics struct {
	AgentType string `json:"agent_type"`

	// Execution counts
	TotalExecutions      int64 `json:"total_executions"`
	SuccessfulExecutions int64 `json:"successful_executions"`
	FailedExecutions     int64 `json:"failed_executions"`
	CurrentlyRunning     int64 `json:"currently_running"`

	// Timing statistics
	TotalDuration         time.Duration `json:"total_duration"`
	AverageDuration       time.Duration `json:"average_duration"`
	MinDuration           time.Duration `json:"min_duration"`
	MaxDuration           time.Duration `json:"max_duration"`
	LastExecutionDuration time.Duration `json:"last_execution_duration"`

	// Success rates
	SuccessRate       float64 `json:"success_rate"`        // Percentage
	RecentSuccessRate float64 `json:"recent_success_rate"` // Last 10 executions

	// Performance tracking
	ExecutionsPerMinute    float64     `json:"executions_per_minute"`
	RecentExecutionTimes   []time.Time `json:"recent_execution_times"`
	RecentExecutionResults []bool      `json:"recent_execution_results"` // true = success, false = failure

	// Error tracking
	CommonErrors  map[string]int64 `json:"common_errors"`
	LastError     string           `json:"last_error"`
	LastErrorTime time.Time        `json:"last_error_time"`

	// Model usage tracking
	ModelUsage     map[string]int64 `json:"model_usage"`
	PreferredModel string           `json:"preferred_model"`

	// Stage performance (validation, prompt_load, execution)
	StagePerformance map[string]*StageStats `json:"stage_performance"`
}

// StageStats tracks performance for specific execution stages
type StageStats struct {
	StageName       string        `json:"stage_name"`
	TotalExecutions int64         `json:"total_executions"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`
	FailureRate     float64       `json:"failure_rate"`
}

// SessionStatistics tracks session-wide metrics
type SessionStatistics struct {
	SessionStart        time.Time     `json:"session_start"`
	TotalEvents         int64         `json:"total_events"`
	EventsPerSecond     float64       `json:"events_per_second"`
	PeakEventsPerSecond float64       `json:"peak_events_per_second"`
	TotalTasks          int64         `json:"total_tasks"`
	ActiveTasks         int64         `json:"active_tasks"`
	CompletedTasks      int64         `json:"completed_tasks"`
	FailedTasks         int64         `json:"failed_tasks"`
	OverallSuccessRate  float64       `json:"overall_success_rate"`
	AverageTaskDuration time.Duration `json:"average_task_duration"`
	SessionUptime       time.Duration `json:"session_uptime"`
}

// PerformanceMetrics tracks system performance indicators
type PerformanceMetrics struct {
	// Throughput metrics
	CurrentThroughput float64 `json:"current_throughput"` // events/second
	AverageThroughput float64 `json:"average_throughput"` // events/second
	PeakThroughput    float64 `json:"peak_throughput"`    // events/second

	// Latency metrics
	AverageLatency time.Duration `json:"average_latency"`
	P95Latency     time.Duration `json:"p95_latency"`
	P99Latency     time.Duration `json:"p99_latency"`

	// Queue metrics
	EventQueueLength    int     `json:"event_queue_length"`
	MaxQueueLength      int     `json:"max_queue_length"`
	QueueProcessingRate float64 `json:"queue_processing_rate"` // events processed per second

	// Resource utilization
	MemoryUsageMB   float64 `json:"memory_usage_mb"`
	CPUUsagePercent float64 `json:"cpu_usage_percent"`
	GoroutineCount  int     `json:"goroutine_count"`

	// Timing windows for trend analysis
	RecentLatencies   []time.Duration `json:"recent_latencies"`
	RecentThroughputs []float64       `json:"recent_throughputs"`
	LastUpdateTime    time.Time       `json:"last_update_time"`
}

// StatSnapshot represents a point-in-time statistics snapshot for trend analysis
type StatSnapshot struct {
	Timestamp             time.Time        `json:"timestamp"`
	EventsPerSecond       float64          `json:"events_per_second"`
	ActiveTasks           int64            `json:"active_tasks"`
	CompletedTasks        int64            `json:"completed_tasks"`
	FailedTasks           int64            `json:"failed_tasks"`
	AverageLatency        time.Duration    `json:"average_latency"`
	MemoryUsageMB         float64          `json:"memory_usage_mb"`
	AgentActivitySnapshot map[string]int64 `json:"agent_activity_snapshot"`
}

// NewExecutionStatistics creates a new ExecutionStatistics instance
func NewExecutionStatistics(maxHistorySize int) *ExecutionStatistics {
	return &ExecutionStatistics{
		AgentStats: make(map[string]*AgentStatistics),
		SessionStats: &SessionStatistics{
			SessionStart: time.Now(),
		},
		PerformanceMetrics: &PerformanceMetrics{
			RecentLatencies:   make([]time.Duration, 0),
			RecentThroughputs: make([]float64, 0),
			LastUpdateTime:    time.Now(),
		},
		HistoryWindow:  make([]*StatSnapshot, 0),
		maxHistorySize: maxHistorySize,
		sessionStart:   time.Now(),
	}
}

// ProcessEvent processes an event and updates statistics
func (es *ExecutionStatistics) ProcessEvent(event events.Event) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	// Update session event count
	es.SessionStats.TotalEvents++
	es.updateEventThroughput()

	switch event.Type() {
	case events.EventTaskStarted:
		es.processTaskStarted(event.(*events.TaskStartedEvent))
	case events.EventTaskProgress:
		es.processTaskProgress(event.(*events.TaskProgressEvent))
	case events.EventTaskCompleted:
		es.processTaskCompleted(event.(*events.TaskCompletedEvent))
	case events.EventTaskFailed:
		es.processTaskFailed(event.(*events.TaskFailedEvent))
	}

	// Update session uptime
	es.SessionStats.SessionUptime = time.Since(es.sessionStart)
}

// processTaskStarted handles task started events
func (es *ExecutionStatistics) processTaskStarted(event *events.TaskStartedEvent) {
	agentStats := es.getOrCreateAgentStats(event.AgentType)

	// Update counters
	agentStats.TotalExecutions++
	agentStats.CurrentlyRunning++
	es.SessionStats.TotalTasks++
	es.SessionStats.ActiveTasks++

	// Track recent execution time for throughput calculation
	agentStats.RecentExecutionTimes = append(agentStats.RecentExecutionTimes, event.Timestamp())
	if len(agentStats.RecentExecutionTimes) > 10 {
		agentStats.RecentExecutionTimes = agentStats.RecentExecutionTimes[1:]
	}

	// Update model usage
	if event.ModelHint != "" {
		if agentStats.ModelUsage == nil {
			agentStats.ModelUsage = make(map[string]int64)
		}
		agentStats.ModelUsage[event.ModelHint]++

		// Update preferred model (most used)
		maxUsage := int64(0)
		for model, usage := range agentStats.ModelUsage {
			if usage > maxUsage {
				maxUsage = usage
				agentStats.PreferredModel = model
			}
		}
	}
}

// processTaskProgress handles task progress events
func (es *ExecutionStatistics) processTaskProgress(event *events.TaskProgressEvent) {
	agentStats := es.getOrCreateAgentStats(event.AgentType)

	// Update stage performance
	if agentStats.StagePerformance == nil {
		agentStats.StagePerformance = make(map[string]*StageStats)
	}

	if _, exists := agentStats.StagePerformance[event.Stage]; !exists {
		agentStats.StagePerformance[event.Stage] = &StageStats{
			StageName: event.Stage,
		}
	}

	stageStats := agentStats.StagePerformance[event.Stage]
	stageStats.TotalExecutions++
	stageStats.TotalDuration += event.Duration
	stageStats.AverageDuration = time.Duration(int64(stageStats.TotalDuration) / stageStats.TotalExecutions)
}

// processTaskCompleted handles task completed events
func (es *ExecutionStatistics) processTaskCompleted(event *events.TaskCompletedEvent) {
	agentStats := es.getOrCreateAgentStats(event.AgentType)

	// Update counters
	agentStats.SuccessfulExecutions++
	agentStats.CurrentlyRunning--
	es.SessionStats.CompletedTasks++
	es.SessionStats.ActiveTasks--

	// Update timing statistics
	agentStats.TotalDuration += event.Duration
	agentStats.LastExecutionDuration = event.Duration

	// Calculate average duration
	if agentStats.TotalExecutions > 0 {
		agentStats.AverageDuration = time.Duration(int64(agentStats.TotalDuration) / agentStats.TotalExecutions)
	}

	// Update min/max duration
	if agentStats.MinDuration == 0 || event.Duration < agentStats.MinDuration {
		agentStats.MinDuration = event.Duration
	}
	if event.Duration > agentStats.MaxDuration {
		agentStats.MaxDuration = event.Duration
	}

	// Update success rates
	agentStats.SuccessRate = float64(agentStats.SuccessfulExecutions) / float64(agentStats.TotalExecutions) * 100

	// Track recent success for recent success rate
	agentStats.RecentExecutionResults = append(agentStats.RecentExecutionResults, true)
	if len(agentStats.RecentExecutionResults) > 10 {
		agentStats.RecentExecutionResults = agentStats.RecentExecutionResults[1:]
	}
	es.updateRecentSuccessRate(agentStats)

	// Update model usage
	if event.Model != "" {
		if agentStats.ModelUsage == nil {
			agentStats.ModelUsage = make(map[string]int64)
		}
		agentStats.ModelUsage[event.Model]++
	}

	// Update session averages
	es.updateSessionAverages()

	// Add latency measurement
	es.PerformanceMetrics.RecentLatencies = append(es.PerformanceMetrics.RecentLatencies, event.Duration)
	if len(es.PerformanceMetrics.RecentLatencies) > 100 {
		es.PerformanceMetrics.RecentLatencies = es.PerformanceMetrics.RecentLatencies[1:]
	}
	es.updateLatencyMetrics()
}

// processTaskFailed handles task failed events
func (es *ExecutionStatistics) processTaskFailed(event *events.TaskFailedEvent) {
	agentStats := es.getOrCreateAgentStats(event.AgentType)

	// Update counters
	agentStats.FailedExecutions++
	agentStats.CurrentlyRunning--
	es.SessionStats.FailedTasks++
	es.SessionStats.ActiveTasks--

	// Update timing statistics with partial execution time
	agentStats.TotalDuration += event.Duration
	agentStats.LastExecutionDuration = event.Duration

	// Calculate average duration
	if agentStats.TotalExecutions > 0 {
		agentStats.AverageDuration = time.Duration(int64(agentStats.TotalDuration) / agentStats.TotalExecutions)
	}

	// Update success rates
	if agentStats.TotalExecutions > 0 {
		agentStats.SuccessRate = float64(agentStats.SuccessfulExecutions) / float64(agentStats.TotalExecutions) * 100
	}

	// Track recent failure for recent success rate
	agentStats.RecentExecutionResults = append(agentStats.RecentExecutionResults, false)
	if len(agentStats.RecentExecutionResults) > 10 {
		agentStats.RecentExecutionResults = agentStats.RecentExecutionResults[1:]
	}
	es.updateRecentSuccessRate(agentStats)

	// Track error information
	agentStats.LastError = event.Error
	agentStats.LastErrorTime = event.Timestamp()

	if agentStats.CommonErrors == nil {
		agentStats.CommonErrors = make(map[string]int64)
	}
	// Simplify error message for grouping (take first 50 chars)
	errorKey := event.Error
	if len(errorKey) > 50 {
		errorKey = errorKey[:50]
	}
	agentStats.CommonErrors[errorKey]++

	// Update stage performance if applicable
	if event.Stage != "" {
		if agentStats.StagePerformance == nil {
			agentStats.StagePerformance = make(map[string]*StageStats)
		}

		if stageStats, exists := agentStats.StagePerformance[event.Stage]; exists {
			// Calculate failure rate for this stage
			totalStageExecutions := stageStats.TotalExecutions
			if totalStageExecutions > 0 {
				// This is an approximation - we'd need more detailed tracking for exact stage failure rates
				stageStats.FailureRate = float64(agentStats.FailedExecutions) / float64(totalStageExecutions) * 100
			}
		}
	}

	// Update session averages
	es.updateSessionAverages()
}

// getOrCreateAgentStats gets or creates agent statistics for a given agent type
func (es *ExecutionStatistics) getOrCreateAgentStats(agentType string) *AgentStatistics {
	if stats, exists := es.AgentStats[agentType]; exists {
		return stats
	}

	stats := &AgentStatistics{
		AgentType:              agentType,
		CommonErrors:           make(map[string]int64),
		ModelUsage:             make(map[string]int64),
		StagePerformance:       make(map[string]*StageStats),
		RecentExecutionTimes:   make([]time.Time, 0),
		RecentExecutionResults: make([]bool, 0),
	}
	es.AgentStats[agentType] = stats
	return stats
}

// updateRecentSuccessRate updates the recent success rate based on last 10 executions
func (es *ExecutionStatistics) updateRecentSuccessRate(stats *AgentStatistics) {
	if len(stats.RecentExecutionResults) == 0 {
		stats.RecentSuccessRate = 0
		return
	}

	successes := 0
	for _, success := range stats.RecentExecutionResults {
		if success {
			successes++
		}
	}
	stats.RecentSuccessRate = float64(successes) / float64(len(stats.RecentExecutionResults)) * 100
}

// updateEventThroughput updates events per second metrics
func (es *ExecutionStatistics) updateEventThroughput() {
	now := time.Now()
	duration := now.Sub(es.SessionStats.SessionStart)

	if duration.Seconds() > 0 {
		es.SessionStats.EventsPerSecond = float64(es.SessionStats.TotalEvents) / duration.Seconds()
	}

	// Update current throughput (events in last 5 seconds)
	es.PerformanceMetrics.RecentThroughputs = append(es.PerformanceMetrics.RecentThroughputs, 1.0) // 1 event
	if len(es.PerformanceMetrics.RecentThroughputs) > 300 {                                        // Keep last 5 minutes worth (assuming 1 event per second)
		es.PerformanceMetrics.RecentThroughputs = es.PerformanceMetrics.RecentThroughputs[1:]
	}

	// Calculate current throughput (events in last 5 seconds)
	recentWindow := 5.0 // seconds
	recentCount := 0
	cutoff := now.Add(-time.Duration(recentWindow) * time.Second)

	for i := range es.PerformanceMetrics.RecentThroughputs {
		if cutoff.Before(now.Add(-time.Duration(len(es.PerformanceMetrics.RecentThroughputs)-i) * time.Second)) {
			recentCount++
		}
	}

	es.PerformanceMetrics.CurrentThroughput = float64(recentCount) / recentWindow
	es.PerformanceMetrics.AverageThroughput = float64(es.SessionStats.TotalEvents) / duration.Seconds()

	// Track peak throughput
	if es.PerformanceMetrics.CurrentThroughput > es.PerformanceMetrics.PeakThroughput {
		es.PerformanceMetrics.PeakThroughput = es.PerformanceMetrics.CurrentThroughput
	}
	if es.SessionStats.EventsPerSecond > es.SessionStats.PeakEventsPerSecond {
		es.SessionStats.PeakEventsPerSecond = es.SessionStats.EventsPerSecond
	}
}

// updateLatencyMetrics calculates latency percentiles
func (es *ExecutionStatistics) updateLatencyMetrics() {
	latencies := es.PerformanceMetrics.RecentLatencies
	if len(latencies) == 0 {
		return
	}

	// Sort latencies for percentile calculation
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)

	// Simple bubble sort for small arrays
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Calculate average
	var total time.Duration
	for _, latency := range sorted {
		total += latency
	}
	es.PerformanceMetrics.AverageLatency = time.Duration(int64(total) / int64(len(sorted)))

	// Calculate P95 and P99
	if len(sorted) >= 20 { // Need reasonable sample size for percentiles
		p95Index := int(float64(len(sorted)) * 0.95)
		p99Index := int(float64(len(sorted)) * 0.99)

		if p95Index < len(sorted) {
			es.PerformanceMetrics.P95Latency = sorted[p95Index]
		}
		if p99Index < len(sorted) {
			es.PerformanceMetrics.P99Latency = sorted[p99Index]
		}
	}
}

// updateSessionAverages updates session-wide average calculations
func (es *ExecutionStatistics) updateSessionAverages() {
	totalTasks := es.SessionStats.CompletedTasks + es.SessionStats.FailedTasks
	if totalTasks > 0 {
		es.SessionStats.OverallSuccessRate = float64(es.SessionStats.CompletedTasks) / float64(totalTasks) * 100
	}

	// Calculate average task duration from all agent stats
	var totalDuration time.Duration
	var totalExecutions int64

	for _, agentStats := range es.AgentStats {
		totalDuration += agentStats.TotalDuration
		totalExecutions += agentStats.TotalExecutions
	}

	if totalExecutions > 0 {
		es.SessionStats.AverageTaskDuration = time.Duration(int64(totalDuration) / totalExecutions)
	}
}

// TakeSnapshot creates a point-in-time statistics snapshot
func (es *ExecutionStatistics) TakeSnapshot() *StatSnapshot {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	// Create agent activity snapshot
	agentActivity := make(map[string]int64)
	for agentType, stats := range es.AgentStats {
		agentActivity[agentType] = stats.TotalExecutions
	}

	snapshot := &StatSnapshot{
		Timestamp:             time.Now(),
		EventsPerSecond:       es.PerformanceMetrics.CurrentThroughput,
		ActiveTasks:           es.SessionStats.ActiveTasks,
		CompletedTasks:        es.SessionStats.CompletedTasks,
		FailedTasks:           es.SessionStats.FailedTasks,
		AverageLatency:        es.PerformanceMetrics.AverageLatency,
		MemoryUsageMB:         es.PerformanceMetrics.MemoryUsageMB,
		AgentActivitySnapshot: agentActivity,
	}

	// Add to history window
	es.HistoryWindow = append(es.HistoryWindow, snapshot)
	if len(es.HistoryWindow) > es.maxHistorySize {
		es.HistoryWindow = es.HistoryWindow[1:]
	}

	return snapshot
}

// GetStatistics returns a read-only copy of current statistics
func (es *ExecutionStatistics) GetStatistics() (*ExecutionStatistics, map[string]*AgentStatistics) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	// Create a deep copy for safe reading
	statsCopy := &ExecutionStatistics{
		SessionStats:       es.SessionStats,
		PerformanceMetrics: es.PerformanceMetrics,
		HistoryWindow:      make([]*StatSnapshot, len(es.HistoryWindow)),
		maxHistorySize:     es.maxHistorySize,
		sessionStart:       es.sessionStart,
	}

	copy(statsCopy.HistoryWindow, es.HistoryWindow)

	// Copy agent stats
	agentStatsCopy := make(map[string]*AgentStatistics)
	for agentType, stats := range es.AgentStats {
		agentStatsCopy[agentType] = &AgentStatistics{
			AgentType:             stats.AgentType,
			TotalExecutions:       stats.TotalExecutions,
			SuccessfulExecutions:  stats.SuccessfulExecutions,
			FailedExecutions:      stats.FailedExecutions,
			CurrentlyRunning:      stats.CurrentlyRunning,
			TotalDuration:         stats.TotalDuration,
			AverageDuration:       stats.AverageDuration,
			MinDuration:           stats.MinDuration,
			MaxDuration:           stats.MaxDuration,
			LastExecutionDuration: stats.LastExecutionDuration,
			SuccessRate:           stats.SuccessRate,
			RecentSuccessRate:     stats.RecentSuccessRate,
			ExecutionsPerMinute:   stats.ExecutionsPerMinute,
			LastError:             stats.LastError,
			LastErrorTime:         stats.LastErrorTime,
			PreferredModel:        stats.PreferredModel,
		}
	}

	return statsCopy, agentStatsCopy
}

// StatisticsPanel component for displaying execution statistics in the TUI
type StatisticsPanel struct {
	statistics     *ExecutionStatistics
	agentStats     map[string]*AgentStatistics
	theme          Theme
	focused        bool
	width          int
	height         int
	scrollOffset   int
	showDetailView bool
	selectedAgent  string
}

// NewStatisticsPanel creates a new statistics panel component
func NewStatisticsPanel(statistics *ExecutionStatistics, theme Theme) *StatisticsPanel {
	return &StatisticsPanel{
		statistics: statistics,
		theme:      theme,
	}
}

// SetData updates the statistics panel data
func (sp *StatisticsPanel) SetData(statistics *ExecutionStatistics, agentStats map[string]*AgentStatistics) {
	sp.statistics = statistics
	sp.agentStats = agentStats
}

// Render renders the statistics panel
func (sp *StatisticsPanel) Render(width, height int, theme Theme) string {
	sp.width = width
	sp.height = height
	sp.theme = theme

	if sp.statistics == nil {
		noStatsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true).
			Padding(1)
		return noStatsStyle.Render("No statistics available")
	}

	var output strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 1)

	output.WriteString(headerStyle.Render("📈 Execution Statistics"))
	output.WriteString("\n\n")

	if sp.showDetailView && sp.selectedAgent != "" {
		// Show detailed view for selected agent
		output.WriteString(sp.renderAgentDetailView(sp.selectedAgent, width-2, theme))
	} else {
		// Show overview
		output.WriteString(sp.renderOverview(width-2, theme))
	}

	// Apply focus styling
	result := output.String()
	if sp.focused {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Primary)).
			Padding(0, 1)
		result = borderStyle.Render(result)
	}

	return result
}

// renderOverview renders the statistics overview
func (sp *StatisticsPanel) renderOverview(width int, theme Theme) string {
	var content strings.Builder

	// Session Statistics
	sessionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true)
	content.WriteString(sessionStyle.Render("Session Overview"))
	content.WriteString("\n")

	if sp.statistics.SessionStats != nil {
		session := sp.statistics.SessionStats
		content.WriteString(fmt.Sprintf("  Uptime: %v\n", session.SessionUptime.Round(time.Second)))
		content.WriteString(fmt.Sprintf("  Total Tasks: %d (Active: %d, Completed: %d, Failed: %d)\n",
			session.TotalTasks, session.ActiveTasks, session.CompletedTasks, session.FailedTasks))
		content.WriteString(fmt.Sprintf("  Success Rate: %.1f%%\n", session.OverallSuccessRate))
		content.WriteString(fmt.Sprintf("  Events/sec: %.2f (Peak: %.2f)\n", session.EventsPerSecond, session.PeakEventsPerSecond))
		if session.AverageTaskDuration > 0 {
			content.WriteString(fmt.Sprintf("  Avg Task Duration: %v\n", session.AverageTaskDuration.Round(time.Millisecond)))
		}
	}

	content.WriteString("\n")

	// Performance Metrics
	perfStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Info)).
		Bold(true)
	content.WriteString(perfStyle.Render("Performance Metrics"))
	content.WriteString("\n")

	if sp.statistics.PerformanceMetrics != nil {
		perf := sp.statistics.PerformanceMetrics
		content.WriteString(fmt.Sprintf("  Throughput: %.2f events/sec (Peak: %.2f)\n", perf.CurrentThroughput, perf.PeakThroughput))
		if perf.AverageLatency > 0 {
			content.WriteString(fmt.Sprintf("  Latency - Avg: %v", perf.AverageLatency.Round(time.Millisecond)))
			if perf.P95Latency > 0 {
				content.WriteString(fmt.Sprintf(", P95: %v", perf.P95Latency.Round(time.Millisecond)))
			}
			if perf.P99Latency > 0 {
				content.WriteString(fmt.Sprintf(", P99: %v", perf.P99Latency.Round(time.Millisecond)))
			}
			content.WriteString("\n")
		}
		if perf.EventQueueLength > 0 {
			content.WriteString(fmt.Sprintf("  Queue: %d events (Max: %d)\n", perf.EventQueueLength, perf.MaxQueueLength))
		}
		if perf.MemoryUsageMB > 0 {
			content.WriteString(fmt.Sprintf("  Memory: %.1f MB\n", perf.MemoryUsageMB))
		}
	}

	content.WriteString("\n")

	// Agent Statistics Summary
	agentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Success)).
		Bold(true)
	content.WriteString(agentStyle.Render("Agent Performance"))
	content.WriteString("\n")

	if len(sp.agentStats) == 0 {
		noAgentsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true)
		content.WriteString(noAgentsStyle.Render("  No agent activity recorded"))
		content.WriteString("\n")
	} else {
		// Table header
		headerFormat := "  %-18s %8s %8s %8s %7s %12s\n"
		content.WriteString(fmt.Sprintf(headerFormat, "Agent", "Total", "Success", "Failed", "Rate", "Avg Duration"))

		// Sort agents by total executions for consistent display
		sortedAgents := make([]string, 0, len(sp.agentStats))
		for agentType := range sp.agentStats {
			sortedAgents = append(sortedAgents, agentType)
		}

		// Simple sort by total executions
		for i := 0; i < len(sortedAgents); i++ {
			for j := i + 1; j < len(sortedAgents); j++ {
				if sp.agentStats[sortedAgents[i]].TotalExecutions < sp.agentStats[sortedAgents[j]].TotalExecutions {
					sortedAgents[i], sortedAgents[j] = sortedAgents[j], sortedAgents[i]
				}
			}
		}

		for _, agentType := range sortedAgents {
			stats := sp.agentStats[agentType]

			// Get agent color
			agentColor, exists := theme.AgentColors[agentType]
			if !exists {
				agentColor = theme.AgentColors["default"]
			}

			agentNameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(agentColor))
			successRateStyle := lipgloss.NewStyle()
			if stats.SuccessRate >= 90 {
				successRateStyle = successRateStyle.Foreground(lipgloss.Color(theme.Success))
			} else if stats.SuccessRate >= 75 {
				successRateStyle = successRateStyle.Foreground(lipgloss.Color(theme.Warning))
			} else {
				successRateStyle = successRateStyle.Foreground(lipgloss.Color(theme.Error))
			}

			// Truncate agent name if too long
			displayName := agentType
			if len(displayName) > 16 {
				displayName = displayName[:15] + "…"
			}

			avgDuration := "N/A"
			if stats.AverageDuration > 0 {
				if stats.AverageDuration < time.Second {
					avgDuration = fmt.Sprintf("%dms", stats.AverageDuration.Milliseconds())
				} else {
					avgDuration = stats.AverageDuration.Round(time.Millisecond).String()
				}
			}

			line := fmt.Sprintf("  %-18s %8d %8d %8d %7.1f%% %12s\n",
				displayName,
				stats.TotalExecutions,
				stats.SuccessfulExecutions,
				stats.FailedExecutions,
				stats.SuccessRate,
				avgDuration,
			)

			// Apply styling to parts of the line
			parts := strings.Fields(line[2:]) // Remove leading spaces
			if len(parts) >= 6 {
				styledLine := fmt.Sprintf("  %s %8s %8s %8s %s %12s\n",
					agentNameStyle.Render(fmt.Sprintf("%-16s", parts[0])),
					parts[1], parts[2], parts[3],
					successRateStyle.Render(parts[4]),
					parts[5],
				)
				content.WriteString(styledLine)
			} else {
				content.WriteString(line)
			}
		}
	}

	// Recent Activity (if we have history)
	if len(sp.statistics.HistoryWindow) > 1 {
		content.WriteString("\n")
		historyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)).
			Bold(true)
		content.WriteString(historyStyle.Render("Recent Trends"))
		content.WriteString("\n")

		recent := sp.statistics.HistoryWindow[len(sp.statistics.HistoryWindow)-1]
		previous := sp.statistics.HistoryWindow[len(sp.statistics.HistoryWindow)-2]

		// Calculate trends
		eventsTrend := recent.EventsPerSecond - previous.EventsPerSecond
		tasksTrend := recent.CompletedTasks - previous.CompletedTasks

		trendSymbol := func(value float64) string {
			if value > 0 {
				return "↑"
			} else if value < 0 {
				return "↓"
			}
			return "→"
		}

		content.WriteString(fmt.Sprintf("  Events/sec: %.2f %s (%.2f)\n",
			recent.EventsPerSecond, trendSymbol(float64(eventsTrend)), eventsTrend))
		content.WriteString(fmt.Sprintf("  Completed: %d %s (+%d)\n",
			recent.CompletedTasks, trendSymbol(float64(tasksTrend)), tasksTrend))

		if recent.AverageLatency > 0 && previous.AverageLatency > 0 {
			latencyTrend := recent.AverageLatency - previous.AverageLatency
			content.WriteString(fmt.Sprintf("  Avg Latency: %v %s (%v)\n",
				recent.AverageLatency.Round(time.Millisecond),
				trendSymbol(float64(latencyTrend)),
				latencyTrend.Round(time.Millisecond)))
		}
	}

	// Usage hints
	content.WriteString("\n")
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted)).
		Italic(true)
	content.WriteString(hintStyle.Render("Press 's' to view detailed statistics, 'r' to reset counters"))

	return content.String()
}

// renderAgentDetailView renders detailed view for a specific agent
func (sp *StatisticsPanel) renderAgentDetailView(agentType string, width int, theme Theme) string {
	stats, exists := sp.agentStats[agentType]
	if !exists {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Error)).
			Italic(true)
		return errorStyle.Render(fmt.Sprintf("No statistics available for agent: %s", agentType))
	}

	var content strings.Builder

	// Agent header
	agentColor, exists := theme.AgentColors[agentType]
	if !exists {
		agentColor = theme.AgentColors["default"]
	}

	agentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(agentColor)).
		Bold(true)
	content.WriteString(agentStyle.Render(fmt.Sprintf("Detailed Statistics: %s", agentType)))
	content.WriteString("\n\n")

	// Execution Summary
	content.WriteString("Execution Summary:\n")
	content.WriteString(fmt.Sprintf("  Total Executions: %d\n", stats.TotalExecutions))
	content.WriteString(fmt.Sprintf("  Successful: %d (%.1f%%)\n", stats.SuccessfulExecutions, stats.SuccessRate))
	content.WriteString(fmt.Sprintf("  Failed: %d (%.1f%%)\n", stats.FailedExecutions, 100.0-stats.SuccessRate))
	content.WriteString(fmt.Sprintf("  Currently Running: %d\n", stats.CurrentlyRunning))
	content.WriteString(fmt.Sprintf("  Recent Success Rate: %.1f%% (last 10)\n", stats.RecentSuccessRate))
	content.WriteString("\n")

	// Timing Statistics
	content.WriteString("Timing Performance:\n")
	if stats.AverageDuration > 0 {
		content.WriteString(fmt.Sprintf("  Average Duration: %v\n", stats.AverageDuration.Round(time.Millisecond)))
	}
	if stats.MinDuration > 0 {
		content.WriteString(fmt.Sprintf("  Fastest: %v\n", stats.MinDuration.Round(time.Millisecond)))
	}
	if stats.MaxDuration > 0 {
		content.WriteString(fmt.Sprintf("  Slowest: %v\n", stats.MaxDuration.Round(time.Millisecond)))
	}
	if stats.LastExecutionDuration > 0 {
		content.WriteString(fmt.Sprintf("  Last Execution: %v\n", stats.LastExecutionDuration.Round(time.Millisecond)))
	}
	content.WriteString("\n")

	// Model Usage
	if len(stats.ModelUsage) > 0 {
		content.WriteString("Model Usage:\n")
		for model, count := range stats.ModelUsage {
			percentage := float64(count) / float64(stats.TotalExecutions) * 100
			content.WriteString(fmt.Sprintf("  %s: %d (%.1f%%)\n", model, count, percentage))
		}
		if stats.PreferredModel != "" {
			content.WriteString(fmt.Sprintf("  Preferred: %s\n", stats.PreferredModel))
		}
		content.WriteString("\n")
	}

	// Error Information
	if stats.FailedExecutions > 0 {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Error)).Bold(true)
		content.WriteString(errorStyle.Render("Error Analysis:"))
		content.WriteString("\n")

		if stats.LastError != "" {
			content.WriteString(fmt.Sprintf("  Last Error: %s\n", truncateString(stats.LastError, 80)))
			if !stats.LastErrorTime.IsZero() {
				content.WriteString(fmt.Sprintf("  Last Error Time: %s\n", stats.LastErrorTime.Format("15:04:05")))
			}
		}

		if len(stats.CommonErrors) > 0 {
			content.WriteString("  Common Errors:\n")
			// Sort errors by frequency
			type errorCount struct {
				error string
				count int64
			}

			var sortedErrors []errorCount
			for err, count := range stats.CommonErrors {
				sortedErrors = append(sortedErrors, errorCount{err, count})
			}

			// Simple sort by count
			for i := 0; i < len(sortedErrors); i++ {
				for j := i + 1; j < len(sortedErrors); j++ {
					if sortedErrors[i].count < sortedErrors[j].count {
						sortedErrors[i], sortedErrors[j] = sortedErrors[j], sortedErrors[i]
					}
				}
			}

			// Show top 3 errors
			maxShow := len(sortedErrors)
			if maxShow > 3 {
				maxShow = 3
			}

			for i := 0; i < maxShow; i++ {
				err := sortedErrors[i]
				percentage := float64(err.count) / float64(stats.FailedExecutions) * 100
				content.WriteString(fmt.Sprintf("    %s: %d (%.1f%%)\n",
					truncateString(err.error, 60), err.count, percentage))
			}
		}
		content.WriteString("\n")
	}

	// Stage Performance
	if len(stats.StagePerformance) > 0 {
		content.WriteString("Stage Performance:\n")
		for stageName, stageStats := range stats.StagePerformance {
			content.WriteString(fmt.Sprintf("  %s: %d executions, avg %v",
				stageName, stageStats.TotalExecutions, stageStats.AverageDuration.Round(time.Millisecond)))
			if stageStats.FailureRate > 0 {
				content.WriteString(fmt.Sprintf(" (%.1f%% failure rate)", stageStats.FailureRate))
			}
			content.WriteString("\n")
		}
	}

	// Back navigation hint
	content.WriteString("\n")
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted)).
		Italic(true)
	content.WriteString(hintStyle.Render("Press 'b' to go back to overview"))

	return content.String()
}

// Update handles updates to the statistics panel
func (sp *StatisticsPanel) Update(msg interface{}) (Component, bool) {
	switch m := msg.(type) {
	case TickMessage:
		// Periodic updates might trigger re-calculation of derived stats
		return sp, true
	case string:
		// Handle key presses for navigation
		switch m {
		case "s":
			sp.showDetailView = !sp.showDetailView
			if sp.showDetailView && len(sp.agentStats) > 0 {
				// Select first agent by default
				for agentType := range sp.agentStats {
					sp.selectedAgent = agentType
					break
				}
			}
			return sp, true
		case "b":
			sp.showDetailView = false
			sp.selectedAgent = ""
			return sp, true
		}
	}
	return sp, false
}

// Focus sets the component as focused
func (sp *StatisticsPanel) Focus() Component {
	sp.focused = true
	return sp
}

// Blur removes focus from the component
func (sp *StatisticsPanel) Blur() Component {
	sp.focused = false
	return sp
}

// IsFocused returns whether the component is focused
func (sp *StatisticsPanel) IsFocused() bool {
	return sp.focused
}

// ScrollUp scrolls the statistics panel up
func (sp *StatisticsPanel) ScrollUp(lines int) {
	sp.scrollOffset -= lines
	if sp.scrollOffset < 0 {
		sp.scrollOffset = 0
	}
}

// ScrollDown scrolls the statistics panel down
func (sp *StatisticsPanel) ScrollDown(lines int) {
	// Implementation would depend on content height calculation
	sp.scrollOffset += lines
}
