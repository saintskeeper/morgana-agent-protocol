#!/bin/bash
#
# Morgana Events Integration Helper
#
# Source this file in your shell scripts to easily log Morgana Protocol events
# Usage: source /Users/walterday/.claude/scripts/morgana_events_integration.sh
#

# Ensure the events directory exists
mkdir -p /tmp/morgana

# Event logging functions for bash scripts
morgana_log_task_start() {
    local task_id="$1"
    local agent_type="$2"
    local prompt="$3"
    local options="${4:-{}}"
    
    python3 -c "
import sys
sys.path.insert(0, '/Users/walterday/.claude/scripts')
from morgana_events import log_task_start
log_task_start('$task_id', '$agent_type', '$prompt', **eval('$options'))
"
}

morgana_log_task_progress() {
    local task_id="$1"
    local stage="$2" 
    local message="$3"
    local progress="${4:-}"
    
    local progress_arg=""
    if [[ -n "$progress" ]]; then
        progress_arg=", $progress"
    fi
    
    python3 -c "
import sys
sys.path.insert(0, '/Users/walterday/.claude/scripts')
from morgana_events import log_task_progress
log_task_progress('$task_id', '$stage', '$message'$progress_arg)
"
}

morgana_log_task_complete() {
    local task_id="$1"
    local output="$2"
    local duration_ms="$3"
    local model="${4:-}"
    
    local model_arg=""
    if [[ -n "$model" ]]; then
        model_arg=", '$model'"
    fi
    
    # Use environment variable to pass potentially multiline output
    export MORGANA_OUTPUT="$output"
    python3 -c "
import sys, os
sys.path.insert(0, '/Users/walterday/.claude/scripts')
from morgana_events import log_task_complete
output = os.environ.get('MORGANA_OUTPUT', '')
log_task_complete('$task_id', output, $duration_ms$model_arg)
"
    unset MORGANA_OUTPUT
}

morgana_log_task_error() {
    local task_id="$1"
    local error="$2" 
    local duration_ms="$3"
    local stage="${4:-}"
    
    local stage_arg=""
    if [[ -n "$stage" ]]; then
        stage_arg=", '$stage'"
    fi
    
    # Use environment variable to pass potentially multiline error
    export MORGANA_ERROR="$error"
    python3 -c "
import sys, os
sys.path.insert(0, '/Users/walterday/.claude/scripts')
from morgana_events import log_task_error
error = os.environ.get('MORGANA_ERROR', '')
log_task_error('$task_id', error, $duration_ms$stage_arg)
"
    unset MORGANA_ERROR
}

# Generate a random task ID
morgana_generate_task_id() {
    python3 -c "import uuid; print(str(uuid.uuid4())[:8])"
}

# View events 
morgana_view_events() {
    python3 /Users/walterday/.claude/scripts/morgana_events_viewer.py "$@"
}

# Follow events in real-time
morgana_follow_events() {
    python3 /Users/walterday/.claude/scripts/morgana_events_viewer.py --follow "$@"
}

# List sessions
morgana_list_sessions() {
    python3 /Users/walterday/.claude/scripts/morgana_events_viewer.py --list-sessions "$@"
}

# Utility function for timing bash tasks
morgana_time_task() {
    local task_id="$1"
    local agent_type="$2"
    local prompt="$3"
    shift 3
    
    local start_time=$(date +%s%N)
    morgana_log_task_start "$task_id" "$agent_type" "$prompt"
    
    # Execute the command
    local output
    local exit_code=0
    
    if output=$("$@" 2>&1); then
        local end_time=$(date +%s%N)
        local duration_ms=$(( (end_time - start_time) / 1000000 ))
        morgana_log_task_complete "$task_id" "$output" "$duration_ms"
        echo "$output"
    else
        exit_code=$?
        local end_time=$(date +%s%N)
        local duration_ms=$(( (end_time - start_time) / 1000000 ))
        morgana_log_task_error "$task_id" "Command failed with exit code $exit_code: $output" "$duration_ms"
        echo "$output" >&2
        return $exit_code
    fi
}

echo "🧙‍♂️ Morgana Event Logging functions loaded:"
echo "   - morgana_log_task_start <task_id> <agent_type> <prompt> [options]"
echo "   - morgana_log_task_progress <task_id> <stage> <message> [progress]"
echo "   - morgana_log_task_complete <task_id> <output> <duration_ms> [model]"
echo "   - morgana_log_task_error <task_id> <error> <duration_ms> [stage]"
echo "   - morgana_generate_task_id"
echo "   - morgana_view_events [options]"
echo "   - morgana_follow_events [options]"
echo "   - morgana_list_sessions [options]"
echo "   - morgana_time_task <task_id> <agent_type> <prompt> <command...>"