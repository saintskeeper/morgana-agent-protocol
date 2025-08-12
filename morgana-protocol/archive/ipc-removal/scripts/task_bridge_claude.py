#!/usr/bin/env python3
"""
Task Bridge for Claude Code - Interfaces between Go orchestrator and Claude's Task tool
This version is designed to work within Claude Code's environment
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
        
        # Get the agent prompt file path
        agent_dir = os.path.expanduser("~/.claude/agents")
        agent_file = os.path.join(agent_dir, f"{data['agent_type']}.md")
        
        # Read the agent prompt if it exists
        agent_context = ""
        if os.path.exists(agent_file):
            with open(agent_file, 'r') as f:
                agent_context = f.read()
        
        # Prepare the full prompt with agent context
        full_prompt = f"""
You are executing as the {data['agent_type']} agent.

{agent_context}

Task: {data['prompt']}
"""
        
        # For Claude Code, we need to output a structured response that indicates
        # the Task tool should be invoked by Claude
        output = {
            "success": True,
            "output": {
                "action": "invoke_task",
                "agent_type": data['agent_type'],
                "prompt": full_prompt,
                "original_prompt": data['prompt']
            },
            "error": None
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