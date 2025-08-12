#!/usr/bin/env python3
"""
Example usage of Morgana Event Logger

This script demonstrates how to integrate the Morgana event logger
with various agent tasks and monitoring scenarios.
"""

import sys
import time
from pathlib import Path

# Add the scripts directory to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from morgana_events import (
    get_logger,
    log_claude_task,
    log_task_start,
    log_task_progress,
    log_task_complete,
    log_task_error
)


@log_claude_task(agent_type="code-implementer")
def example_coding_task(feature_name: str, complexity: str = "medium"):
    """Example coding task with automatic event logging."""
    print(f"Implementing {feature_name} (complexity: {complexity})")
    
    # Simulate work
    time.sleep(0.5)
    
    if complexity == "high":
        # Simulate potential failure
        import random
        if random.random() < 0.3:
            raise Exception("Complex feature implementation failed")
    
    return f"Successfully implemented {feature_name}"


def manual_task_example():
    """Example of manual event logging for more control."""
    logger = get_logger()
    task_id = "manual-task-001"
    
    try:
        start_time = time.time()
        
        # Start logging
        logger.task_started(
            task_id=task_id,
            agent_type="test-specialist",
            prompt="Run comprehensive test suite",
            options={"coverage": True, "integration": True}
        )
        
        # Progress updates
        logger.task_progress(task_id, "setup", "Initializing test environment", 0.1)
        time.sleep(0.2)
        
        logger.task_progress(task_id, "unit-tests", "Running unit tests", 0.4)
        time.sleep(0.3)
        
        logger.task_progress(task_id, "integration-tests", "Running integration tests", 0.7)
        time.sleep(0.4)
        
        logger.task_progress(task_id, "coverage", "Generating coverage report", 0.9)
        time.sleep(0.1)
        
        # Complete
        duration_ms = int((time.time() - start_time) * 1000)
        logger.task_completed(
            task_id=task_id,
            output="All tests passed: 147 unit tests, 23 integration tests. Coverage: 94%",
            duration_ms=duration_ms,
            model="test-runner-v2"
        )
        
        print(f"Manual task completed in {duration_ms}ms")
        
    except Exception as e:
        duration_ms = int((time.time() - start_time) * 1000)
        logger.task_failed(task_id, str(e), duration_ms, "integration-tests")
        raise


def quick_logging_example():
    """Example using convenience functions."""
    task_id = "quick-task-002"
    
    log_task_start(task_id, "validation-expert", "Validate API responses", timeout=30)
    
    log_task_progress(task_id, "validation", "Checking response schemas", 0.3)
    time.sleep(0.2)
    
    log_task_progress(task_id, "validation", "Testing error handling", 0.7)
    time.sleep(0.2)
    
    log_task_complete(task_id, "All API validations passed", 400, "validator-gpt-4")
    print("Quick logging example completed")


if __name__ == "__main__":
    print("Running Morgana Event Logger examples...")
    print(f"Events will be logged to: /tmp/morgana/events.jsonl")
    
    # Example 1: Automatic logging with decorator
    print("\n1. Testing automatic logging with decorator:")
    try:
        result = example_coding_task("user authentication", "medium")
        print(f"Result: {result}")
    except Exception as e:
        print(f"Task failed: {e}")
    
    # Example 2: Manual logging for full control
    print("\n2. Testing manual event logging:")
    try:
        manual_task_example()
    except Exception as e:
        print(f"Manual task failed: {e}")
    
    # Example 3: Quick convenience functions
    print("\n3. Testing convenience functions:")
    quick_logging_example()
    
    print(f"\nAll examples completed. Check /tmp/morgana/events.jsonl for logged events.")
    print(f"Session ID: {get_logger().session_id}")