#!/usr/bin/env python3
"""
Morgana Command Polling Integration Example

Demonstrates how to use the command polling system with claude_agent_executor.
"""

import time
import sys
from pathlib import Path

# Add current directory to Python path
sys.path.insert(0, str(Path(__file__).parent))

from claude_agent_executor import execute_code_implementer, ClaudeAgentExecutor
from morgana_command_poll import integration_check_point


def example_with_integration_checkpoint():
    """Example using integration_check_point for manual control."""
    print("=== Example: Long-running task with command polling ===")
    print("Run in another terminal:")
    print("  ./morgana-cmd pause    # Pause execution") 
    print("  ./morgana-cmd resume   # Resume execution")
    print("  ./morgana-cmd stop     # Stop execution")
    print("  ./morgana-cmd check    # Check current state")
    print()
    
    for i in range(10):
        print(f"Processing item {i+1}/10...")
        
        # Integration checkpoint - handles pause/resume/stop
        if not integration_check_point():
            print("Task stopped by command!")
            return False
        
        # Simulate some work
        time.sleep(2)
    
    print("Task completed successfully!")
    return True


def example_with_agent_executor():
    """Example using ClaudeAgentExecutor with command polling."""
    print("\n=== Example: Agent execution with command polling ===")
    
    # Create executor with command polling enabled (default)
    executor = ClaudeAgentExecutor(enable_command_polling=True)
    
    print("Executing code-implementer task...")
    print("You can send commands during execution using:")
    print("  ./morgana-cmd pause")
    print("  ./morgana-cmd resume") 
    print("  ./morgana-cmd stop")
    print()
    
    # Execute a task - command polling is integrated automatically
    result = executor.execute_agent(
        agent_type="code-implementer",
        prompt="Create a simple hello world function with proper error handling",
        timeout=60
    )
    
    print(f"Task result: {result['success']}")
    if result.get('error'):
        print(f"Error: {result['error']}")
    else:
        print("Task completed successfully!")


def example_with_convenience_function():
    """Example using convenience functions."""
    print("\n=== Example: Using convenience functions ===")
    
    print("Executing test-specialist task with command polling...")
    
    # Use convenience function - command polling is automatic
    result = execute_code_implementer(
        "Write unit tests for a calculator function",
        timeout=30
    )
    
    print(f"Task completed: {result['success']}")
    if result.get('execution_mode') == 'command_stopped':
        print("Task was stopped by command!")
    elif result.get('error'):
        print(f"Task failed: {result['error']}")


def main():
    """Main function to run examples."""
    print("Morgana Command Polling Integration Examples")
    print("=" * 50)
    
    if len(sys.argv) > 1:
        example_type = sys.argv[1]
    else:
        example_type = "checkpoint"
    
    if example_type == "checkpoint":
        example_with_integration_checkpoint()
    elif example_type == "executor":
        example_with_agent_executor()
    elif example_type == "convenience":
        example_with_convenience_function()
    elif example_type == "all":
        example_with_integration_checkpoint()
        example_with_agent_executor()
        example_with_convenience_function()
    else:
        print(f"Unknown example type: {example_type}")
        print("Available types: checkpoint, executor, convenience, all")
        sys.exit(1)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nExample interrupted by user")
    except Exception as e:
        print(f"\nExample failed with error: {e}")
        raise