#!/usr/bin/env python3
"""
Slow Task Bridge - Simulates slow execution for timeout testing
"""
import json
import sys
import os
import signal
import time
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
    
    # Debug print to stderr to confirm script is running
    print(f"SLOW BRIDGE: Started, PID={os.getpid()}", file=sys.stderr)
    
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
        
        # Debug print to confirm processing
        print(f"SLOW BRIDGE: Processing {data['agent_type']} task", file=sys.stderr)
        
        # Simulate slow execution - all tasks take 2 seconds
        print("SLOW BRIDGE: Sleeping for 2 seconds...", file=sys.stderr)
        time.sleep(2)
        print("SLOW BRIDGE: Woke up after sleep", file=sys.stderr)
        
        # Return mock response
        output = {
            "success": True,
            "output": f"[SLOW MOCK] Executed {data['agent_type']} agent after delay",
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