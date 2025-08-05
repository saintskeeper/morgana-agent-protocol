#!/usr/bin/env python3
"""
Task Bridge - Interfaces between Go orchestrator and Claude Code's Task function
"""
import json
import sys
import os
import signal
import traceback

def timeout_handler(signum, frame):
    """Handle timeout signal from parent process"""
    output = {
        "success": False,
        "output": None,
        "error": "Task execution timed out",
        "timeout": True
    }
    print(json.dumps(output))
    sys.exit(1)

def main():
    # Set up timeout handler
    signal.signal(signal.SIGTERM, timeout_handler)
    signal.signal(signal.SIGINT, timeout_handler)
    
    try:
        # Read JSON input from stdin
        input_data = sys.stdin.read()
        if not input_data.strip():
            raise ValueError("No input data received")
        
        data = json.loads(input_data)
        
        # Validate required fields
        if "agent_type" not in data:
            raise ValueError("Missing required field: agent_type")
        if "prompt" not in data:
            raise ValueError("Missing required field: prompt")
        
        # Import Task function - this will only work within Claude Code environment
        try:
            # In Claude Code, Task is available in the global scope
            result = Task(
                subagent_type="general-purpose",  # Always use general-purpose
                prompt=f"{data['agent_type']} agent prompt:\n{data['prompt']}",
                description=f"Task for {data['agent_type']} agent"
            )
            
            # Return successful result
            output = {
                "success": True,
                "output": result,
                "error": None
            }
            
        except NameError:
            # Task function not available - we're outside Claude Code
            # Return mock response for testing
            output = {
                "success": True,
                "output": f"[MOCK] Executed {data['agent_type']} agent with prompt length: {len(data['prompt'])}",
                "error": None,
                "mock": True
            }
    
    except Exception as e:
        # Return error response
        output = {
            "success": False,
            "output": None,
            "error": str(e),
            "traceback": traceback.format_exc() if os.getenv("MORGANA_DEBUG") == "true" else None
        }
    
    # Write JSON response to stdout
    print(json.dumps(output))

if __name__ == "__main__":
    main()