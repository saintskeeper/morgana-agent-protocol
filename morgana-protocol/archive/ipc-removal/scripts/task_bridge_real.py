#!/usr/bin/env python3
"""
Real Task Bridge for Morgana Protocol
This version creates proper Task tool invocations for Claude Code
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
        
        # Check if we're in Claude Code environment
        # We can detect this by checking for specific environment variables
        # or by trying to import Claude-specific modules
        in_claude = os.getenv("CLAUDE_CODE") == "true" or os.getenv("ANTHROPIC_CLAUDE") == "true"
        
        if not in_claude:
            # Outside Claude Code - return a structured response that indicates
            # a Task tool invocation is needed
            
            # Read agent configuration to get the right prompt
            agent_dir = os.path.expanduser("~/.claude/agents")
            agent_file = os.path.join(agent_dir, f"{data['agent_type']}.md")
            
            agent_config = {}
            if os.path.exists(agent_file):
                with open(agent_file, 'r') as f:
                    content = f.read()
                    # Parse YAML frontmatter if present
                    if content.startswith('---'):
                        yaml_end = content.find('---', 3)
                        if yaml_end > 0:
                            # Extract the agent description after frontmatter
                            agent_config['description'] = content[yaml_end+3:].strip()
            
            # Create a response that signals Claude to use Task tool
            output = {
                "success": True,
                "output": f"[TASK_REQUIRED] Agent: {data['agent_type']} | Prompt: {data['prompt']}",
                "task_invocation": {
                    "tool": "Task",
                    "subagent_type": "general-purpose",
                    "description": f"Execute {data['agent_type']} agent",
                    "prompt": f"""You are acting as the {data['agent_type']} agent.

{agent_config.get('description', '')}

Task to complete:
{data['prompt']}

Please complete this task following the agent's specialized role and best practices."""
                },
                "error": None
            }
        else:
            # Inside Claude Code - we still can't directly call Task from Python
            # but we can signal that it should be called
            output = {
                "success": True,
                "output": f"[CLAUDE_ENV] Ready for Task execution: {data['agent_type']}",
                "requires_task": True,
                "agent_type": data['agent_type'],
                "prompt": data['prompt'],
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