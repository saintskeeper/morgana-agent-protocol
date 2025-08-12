#!/usr/bin/env python3
"""
Demonstration of claude_agent_executor.py usage

This script shows how to use the Claude Code Agent Executor both with
and without the Task function available.
"""

import sys
import json
from pathlib import Path

# Add the scripts directory to the path
sys.path.append(str(Path(__file__).parent))

from claude_agent_executor import (
    ClaudeAgentExecutor,
    execute_code_implementer,
    execute_test_specialist,
    execute_validation_expert,
    execute_sprint_planner
)


def demo_single_execution():
    """Demonstrate single agent execution."""
    print("=== Single Agent Execution Demo ===")
    
    # Create executor
    executor = ClaudeAgentExecutor()
    
    # Execute code implementer
    result = executor.execute_agent(
        agent_type="code-implementer",
        prompt="Create a simple user authentication function",
        task_id="demo-auth",
        model="claude-sonnet-4"
    )
    
    print(f"Task ID: {result['task_id']}")
    print(f"Success: {result['success']}")
    print(f"Duration: {result['duration_ms']}ms")
    print(f"Execution Mode: {result['execution_mode']}")
    print()


def demo_convenience_functions():
    """Demonstrate convenience functions."""
    print("=== Convenience Functions Demo ===")
    
    # Test each agent type
    agents = [
        ("code-implementer", execute_code_implementer),
        ("test-specialist", execute_test_specialist),
        ("validation-expert", execute_validation_expert),
        ("sprint-planner", execute_sprint_planner)
    ]
    
    for agent_name, agent_func in agents:
        result = agent_func(f"Demo task for {agent_name}")
        print(f"{agent_name}: {result['success']} ({result['execution_mode']})")
    
    print()


def demo_parallel_execution():
    """Demonstrate parallel agent execution."""
    print("=== Parallel Execution Demo ===")
    
    executor = ClaudeAgentExecutor()
    
    tasks = [
        {
            "agent_type": "code-implementer",
            "prompt": "Implement user registration",
            "task_id": "demo-register"
        },
        {
            "agent_type": "test-specialist", 
            "prompt": "Create tests for user registration",
            "task_id": "demo-register-tests"
        },
        {
            "agent_type": "validation-expert",
            "prompt": "Review user registration code",
            "task_id": "demo-register-review"
        }
    ]
    
    results = executor.execute_parallel(tasks)
    
    for i, result in enumerate(results, 1):
        print(f"Task {i}: {result['agent_type']} - {result['success']} ({result['execution_mode']})")
    
    print()


def demo_task_function_simulation():
    """Simulate what happens when Task function is available."""
    print("=== Task Function Simulation ===")
    print("This shows what would happen if the Task function were available:")
    
    # Mock Task function for demonstration
    def mock_task(prompt, **kwargs):
        return {
            "status": "completed",
            "result": f"Task executed with prompt: {prompt[:50]}...",
            "model": kwargs.get("model", "default"),
            "parameters": kwargs
        }
    
    # Temporarily add to globals to simulate Claude Code environment
    import builtins
    original_task = getattr(builtins, 'Task', None)
    setattr(builtins, 'Task', mock_task)
    
    try:
        result = execute_code_implementer(
            "Create authentication with Task function available",
            model="claude-sonnet-4"
        )
        
        print(f"Success: {result['success']}")
        print(f"Execution Mode: {result['execution_mode']}")
        print(f"Result: {result['result']['status']}")
        
    finally:
        # Clean up
        if original_task is None:
            delattr(builtins, 'Task')
        else:
            setattr(builtins, 'Task', original_task)
    
    print()


def demo_error_handling():
    """Demonstrate error handling."""
    print("=== Error Handling Demo ===")
    
    executor = ClaudeAgentExecutor()
    
    # Test invalid agent type
    try:
        result = executor.execute_agent(
            agent_type="invalid-agent",
            prompt="This should fail"
        )
    except ValueError as e:
        print(f"Caught expected error: {e}")
    
    # Test empty parallel tasks
    empty_results = executor.execute_parallel([])
    print(f"Empty parallel execution results: {len(empty_results)} tasks")
    
    # Test malformed parallel task
    malformed_results = executor.execute_parallel([
        {"agent_type": "code-implementer"},  # Missing prompt
        {"prompt": "Missing agent type"}      # Missing agent_type
    ])
    
    for i, result in enumerate(malformed_results, 1):
        print(f"Malformed task {i}: {result['success']} - {result.get('error', 'No error')}")
    
    print()


def main():
    """Run all demonstrations."""
    print("Claude Code Agent Executor - Demonstration\n")
    
    demo_single_execution()
    demo_convenience_functions()
    demo_parallel_execution()
    demo_task_function_simulation()
    demo_error_handling()
    
    print("=== Event Logging ===")
    print("All task executions are logged to /tmp/morgana/events.jsonl")
    print("Use morgana-monitor-ctl.sh to view real-time events")


if __name__ == "__main__":
    main()