#!/usr/bin/env python3
"""
Morgana Event Viewer

A simple utility to view and monitor Morgana Protocol event logs in real-time.
"""

import json
import sys
import time
from pathlib import Path
from datetime import datetime


def format_event(event_data):
    """Format an event for display."""
    timestamp = event_data.get("timestamp", "")
    if "+" in timestamp or timestamp.endswith("Z"):
        # Parse and format timestamp
        try:
            if timestamp.endswith("Z"):
                dt = datetime.fromisoformat(timestamp[:-1])
            else:
                dt = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            formatted_time = dt.strftime("%H:%M:%S")
        except:
            formatted_time = timestamp[:8] if len(timestamp) >= 8 else timestamp
    else:
        formatted_time = timestamp[:8] if len(timestamp) >= 8 else timestamp
    
    event_type = event_data.get("event_type", "unknown")
    agent_type = event_data.get("agent_type", "unknown")
    session_id = event_data.get("session_id", "")[:8]
    task_id = event_data.get("task_id", "")
    
    # Color coding for different event types
    colors = {
        "task_started": "\033[92m",  # Green
        "task_progress": "\033[94m", # Blue  
        "task_completed": "\033[93m", # Yellow
        "task_failed": "\033[91m",    # Red
    }
    reset = "\033[0m"
    color = colors.get(event_type, "")
    
    base_info = f"{color}[{formatted_time}] {event_type.upper()}{reset} | {agent_type} | {session_id}"
    
    if task_id:
        base_info += f" | {task_id}"
    
    # Add event-specific details
    details = []
    
    if event_type == "task_started":
        prompt = event_data.get("prompt", "")[:60]
        if len(event_data.get("prompt", "")) > 60:
            prompt += "..."
        details.append(f"Prompt: {prompt}")
    
    elif event_type == "task_progress":
        stage = event_data.get("stage", "")
        message = event_data.get("message", "")
        progress = event_data.get("progress")
        details.append(f"Stage: {stage}")
        details.append(f"Message: {message}")
        if progress is not None:
            details.append(f"Progress: {progress:.1%}")
    
    elif event_type == "task_completed":
        duration_ms = event_data.get("duration_ms", 0)
        model = event_data.get("model", "")
        output = event_data.get("output", "")[:60]
        if len(event_data.get("output", "")) > 60:
            output += "..."
        details.append(f"Duration: {duration_ms}ms")
        if model:
            details.append(f"Model: {model}")
        details.append(f"Output: {output}")
    
    elif event_type == "task_failed":
        duration_ms = event_data.get("duration_ms", 0)
        error = event_data.get("error", "")[:60]
        if len(event_data.get("error", "")) > 60:
            error += "..."
        stage = event_data.get("stage", "")
        details.append(f"Duration: {duration_ms}ms")
        if stage:
            details.append(f"Failed at: {stage}")
        details.append(f"Error: {error}")
    
    if details:
        detail_str = " | ".join(details)
        return f"{base_info}\n    {detail_str}"
    else:
        return base_info


def view_events(log_file="/tmp/morgana/events.jsonl", follow=False, filter_session=None):
    """View events from the log file."""
    log_path = Path(log_file)
    
    if not log_path.exists():
        print(f"Event log file not found: {log_file}")
        print("Run some Morgana agents to generate events first.")
        return
    
    print(f"Viewing events from: {log_file}")
    if filter_session:
        print(f"Filtering by session: {filter_session}")
    print("-" * 80)
    
    # Read existing events
    try:
        with open(log_path, "r") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                
                try:
                    event = json.loads(line)
                    if filter_session and event.get("session_id") != filter_session:
                        continue
                    print(format_event(event))
                except json.JSONDecodeError as e:
                    print(f"Error parsing event: {e}")
                    continue
    
    except Exception as e:
        print(f"Error reading log file: {e}")
        return
    
    # Follow mode (tail -f like behavior)
    if follow:
        print("\n" + "="*80)
        print("Following new events... (Press Ctrl+C to exit)")
        print("="*80)
        
        last_position = log_path.stat().st_size
        
        try:
            while True:
                current_size = log_path.stat().st_size
                
                if current_size > last_position:
                    with open(log_path, "r") as f:
                        f.seek(last_position)
                        for line in f:
                            line = line.strip()
                            if not line:
                                continue
                            
                            try:
                                event = json.loads(line)
                                if filter_session and event.get("session_id") != filter_session:
                                    continue
                                print(format_event(event))
                            except json.JSONDecodeError:
                                continue
                    
                    last_position = current_size
                
                time.sleep(0.5)
                
        except KeyboardInterrupt:
            print("\nStopped following events.")


def list_sessions(log_file="/tmp/morgana/events.jsonl"):
    """List all session IDs in the log file."""
    log_path = Path(log_file)
    
    if not log_path.exists():
        print(f"Event log file not found: {log_file}")
        return
    
    sessions = {}
    
    try:
        with open(log_path, "r") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                
                try:
                    event = json.loads(line)
                    session_id = event.get("session_id", "unknown")
                    timestamp = event.get("timestamp", "")
                    
                    if session_id not in sessions:
                        sessions[session_id] = {
                            "first_seen": timestamp,
                            "last_seen": timestamp,
                            "event_count": 0
                        }
                    else:
                        sessions[session_id]["last_seen"] = timestamp
                    
                    sessions[session_id]["event_count"] += 1
                    
                except json.JSONDecodeError:
                    continue
    
    except Exception as e:
        print(f"Error reading log file: {e}")
        return
    
    print("Active sessions:")
    print("-" * 80)
    for session_id, info in sessions.items():
        first_time = info["first_seen"][:19] if len(info["first_seen"]) >= 19 else info["first_seen"]
        last_time = info["last_seen"][:19] if len(info["last_seen"]) >= 19 else info["last_seen"]
        count = info["event_count"]
        print(f"Session: {session_id} | Events: {count:3d} | First: {first_time} | Last: {last_time}")


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="View Morgana Protocol event logs")
    parser.add_argument("--log-file", "-f", default="/tmp/morgana/events.jsonl",
                       help="Path to the event log file")
    parser.add_argument("--follow", "-F", action="store_true",
                       help="Follow new events (like tail -f)")
    parser.add_argument("--session", "-s", 
                       help="Filter events by session ID")
    parser.add_argument("--list-sessions", "-l", action="store_true",
                       help="List all session IDs")
    
    args = parser.parse_args()
    
    if args.list_sessions:
        list_sessions(args.log_file)
    else:
        view_events(args.log_file, args.follow, args.session)