#!/usr/bin/env python3
"""
Morgana Event Stream Logger

A simple event logger that writes JSON lines to /tmp/morgana/events.jsonl
for real-time monitoring of Morgana Protocol agent activities.
"""

import json
import os
import time
import uuid
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, Optional


class MorganaEventLogger:
    """Simple event logger for Morgana Protocol activities."""
    
    def __init__(self, log_file: str = "/tmp/morgana/events.jsonl"):
        """Initialize the event logger.
        
        Args:
            log_file: Path to the JSON lines log file
        """
        self.log_file = Path(log_file)
        self.session_id = str(uuid.uuid4())[:8]
        
        # Ensure the directory exists
        self.log_file.parent.mkdir(parents=True, exist_ok=True)
        
        # Initialize log file if it doesn't exist
        if not self.log_file.exists():
            self.log_file.touch()
    
    def _write_event(self, event_data: Dict[str, Any]) -> None:
        """Write an event to the log file with immediate flush.
        
        Args:
            event_data: The event data to log
        """
        # Add standard metadata
        event_data.update({
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "session_id": self.session_id,
            "pid": os.getpid()
        })
        
        # Write as JSON line with immediate flush
        with open(self.log_file, "a", encoding="utf-8") as f:
            json.dump(event_data, f, separators=(',', ':'))
            f.write("\n")
            f.flush()
            os.fsync(f.fileno())  # Force write to disk
    
    def log_agent_event(self, event_type: str, agent_type: str, **kwargs) -> None:
        """Log a generic agent event.
        
        Args:
            event_type: Type of event (e.g., 'started', 'progress', 'completed', 'failed')
            agent_type: Type of agent (e.g., 'code-implementer', 'test-specialist')
            **kwargs: Additional event-specific data
        """
        event_data = {
            "event_type": event_type,
            "agent_type": agent_type,
            **kwargs
        }
        self._write_event(event_data)
    
    def task_started(self, task_id: str, agent_type: str, prompt: str, 
                    options: Optional[Dict[str, Any]] = None) -> None:
        """Log task started event.
        
        Args:
            task_id: Unique identifier for the task
            agent_type: Type of agent handling the task
            prompt: The task prompt/description
            options: Optional configuration options
        """
        event_data = {
            "event_type": "task_started",
            "task_id": task_id,
            "agent_type": agent_type,
            "prompt": prompt[:500] + "..." if len(prompt) > 500 else prompt,  # Truncate long prompts
            "options": options or {}
        }
        self._write_event(event_data)
    
    def task_progress(self, task_id: str, stage: str, message: str, 
                     progress: Optional[float] = None) -> None:
        """Log task progress event.
        
        Args:
            task_id: Unique identifier for the task
            stage: Current stage of the task
            message: Progress message
            progress: Optional progress percentage (0.0-1.0)
        """
        event_data = {
            "event_type": "task_progress",
            "task_id": task_id,
            "stage": stage,
            "message": message,
        }
        if progress is not None:
            event_data["progress"] = progress
        
        self._write_event(event_data)
    
    def task_completed(self, task_id: str, output: str, duration_ms: int, 
                      model: Optional[str] = None) -> None:
        """Log task completed event.
        
        Args:
            task_id: Unique identifier for the task
            output: Task output/result
            duration_ms: Task duration in milliseconds
            model: Optional model used for the task
        """
        event_data = {
            "event_type": "task_completed",
            "task_id": task_id,
            "output": output[:1000] + "..." if len(output) > 1000 else output,  # Truncate long output
            "duration_ms": duration_ms,
        }
        if model:
            event_data["model"] = model
        
        self._write_event(event_data)
    
    def task_failed(self, task_id: str, error: str, duration_ms: int, 
                   stage: Optional[str] = None) -> None:
        """Log task failed event.
        
        Args:
            task_id: Unique identifier for the task
            error: Error message
            duration_ms: Task duration in milliseconds before failure
            stage: Optional stage where the failure occurred
        """
        event_data = {
            "event_type": "task_failed",
            "task_id": task_id,
            "error": error,
            "duration_ms": duration_ms,
        }
        if stage:
            event_data["stage"] = stage
        
        self._write_event(event_data)


# Global logger instance
_global_logger: Optional[MorganaEventLogger] = None


def get_logger() -> MorganaEventLogger:
    """Get or create the global event logger instance."""
    global _global_logger
    if _global_logger is None:
        _global_logger = MorganaEventLogger()
    return _global_logger


def log_claude_task(task_id: Optional[str] = None, agent_type: str = "claude-code") -> callable:
    """Decorator for Claude Code Task() integration.
    
    Args:
        task_id: Optional task ID, will generate one if not provided
        agent_type: Type of agent executing the task
    
    Returns:
        Decorator function
    """
    def decorator(func):
        def wrapper(*args, **kwargs):
            logger = get_logger()
            actual_task_id = task_id or str(uuid.uuid4())[:8]
            start_time = time.time()
            
            try:
                # Log task started
                prompt = f"Executing {func.__name__}"
                if args:
                    prompt += f" with args: {args}"
                if kwargs:
                    prompt += f" with kwargs: {kwargs}"
                
                logger.task_started(actual_task_id, agent_type, prompt)
                
                # Execute the task
                result = func(*args, **kwargs)
                
                # Log completion
                duration_ms = int((time.time() - start_time) * 1000)
                output = str(result) if result is not None else "Task completed successfully"
                logger.task_completed(actual_task_id, output, duration_ms)
                
                return result
                
            except Exception as e:
                # Log failure
                duration_ms = int((time.time() - start_time) * 1000)
                logger.task_failed(actual_task_id, str(e), duration_ms)
                raise
        
        return wrapper
    return decorator


# Convenience functions for quick logging
def log_task_start(task_id: str, agent_type: str, prompt: str, **options):
    """Quick function to log task start."""
    get_logger().task_started(task_id, agent_type, prompt, options)


def log_task_progress(task_id: str, stage: str, message: str, progress: Optional[float] = None):
    """Quick function to log task progress."""
    get_logger().task_progress(task_id, stage, message, progress)


def log_task_complete(task_id: str, output: str, duration_ms: int, model: Optional[str] = None):
    """Quick function to log task completion."""
    get_logger().task_completed(task_id, output, duration_ms, model)


def log_task_error(task_id: str, error: str, duration_ms: int, stage: Optional[str] = None):
    """Quick function to log task failure."""
    get_logger().task_failed(task_id, error, duration_ms, stage)


if __name__ == "__main__":
    # Demo usage
    logger = MorganaEventLogger()
    
    # Test basic logging
    task_id = str(uuid.uuid4())[:8]
    logger.task_started(task_id, "code-implementer", "Create morgana event logger")
    logger.task_progress(task_id, "implementation", "Writing core logger class", 0.5)
    logger.task_progress(task_id, "testing", "Running basic tests", 0.8)
    logger.task_completed(task_id, "Successfully created morgana_events.py", 1500, "claude-sonnet-4")
    
    print(f"Demo events logged to {logger.log_file}")
    print(f"Session ID: {logger.session_id}")