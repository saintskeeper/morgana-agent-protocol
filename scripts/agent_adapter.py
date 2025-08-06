#!/usr/bin/env python3
import subprocess
import json
import sys
import os

def AgentAdapter(agent_type, prompt, **kwargs):
    """
    Adapter using Morgana Protocol - bridges specialized agents with Claude Code's Task tool
    
    Args:
        agent_type: Type of specialized agent (code-implementer, sprint-planner, etc.)
        prompt: The task prompt for the agent
        **kwargs: Additional options passed to the agent
    
    Returns:
        str: Agent execution result
    """
    # Find morgana binary
    morgana_bin = os.path.expanduser("~/.claude/morgana-protocol/dist/morgana")
    if not os.path.exists(morgana_bin):
        # Try PATH
        morgana_bin = "morgana"
    
    # Build command
    cmd = [morgana_bin, "--"]
    cmd.extend(["--agent", agent_type])
    cmd.extend(["--prompt", prompt])
    
    # Add any additional options as JSON in environment
    env = os.environ.copy()
    if kwargs:
        env["MORGANA_OPTIONS"] = json.dumps(kwargs)
    
    # Execute
    try:
        proc = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            check=True,
            env=env
        )
        
        # Parse output
        result = json.loads(proc.stdout)
        if result.get("success"):
            return result["results"][0]["output"]
        else:
            raise Exception(f"Morgana execution failed: {result}")
            
    except subprocess.CalledProcessError as e:
        raise Exception(f"Morgana failed with code {e.returncode}: {e.stderr}")
    except json.JSONDecodeError as e:
        raise Exception(f"Failed to parse Morgana output: {e}")

# CLI usage
if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: agent_adapter.py <agent-type> <prompt>")
        sys.exit(1)
    
    try:
        result = AgentAdapter(sys.argv[1], sys.argv[2])
        print(result)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
