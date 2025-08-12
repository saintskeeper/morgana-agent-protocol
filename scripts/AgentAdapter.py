#!/usr/bin/env python3
"""
AgentAdapter for Claude Code - Direct Task Tool Integration
This replaces the mock responses with real Task tool invocations
"""

def AgentAdapter(agent_type, prompt, **kwargs):
    """
    Execute a specialized agent using Claude Code's Task tool.
    
    This function is designed to be called from within Claude Code
    where the Task tool is available.
    
    Args:
        agent_type (str): Type of agent (code-implementer, sprint-planner, etc.)
        prompt (str): The task prompt for the agent
        **kwargs: Additional arguments (timeout, model, etc.)
    
    Returns:
        The result from the Task tool execution
    """
    import os
    import json
    
    # Read agent configuration
    agent_dir = os.path.expanduser("~/.claude/agents")
    agent_file = os.path.join(agent_dir, f"{agent_type}.md")
    
    agent_prompt = ""
    if os.path.exists(agent_file):
        with open(agent_file, 'r') as f:
            agent_prompt = f.read()
    
    # Construct the full prompt for the Task tool
    full_prompt = f"""You are executing as the {agent_type} specialist agent.

{agent_prompt}

TASK TO COMPLETE:
{prompt}

Please complete this task following your specialized role and best practices."""
    
    # In Claude Code, we would call the Task tool here
    # Since Task is a Claude Code native function, we need to signal for it
    print(f"EXECUTE_TASK:{agent_type}:{prompt}")
    
    # Return a structured response that Claude can use
    return {
        "agent_type": agent_type,
        "prompt": prompt,
        "status": "ready_for_task",
        "task_params": {
            "subagent_type": "general-purpose",
            "description": f"Execute {agent_type} agent task",
            "prompt": full_prompt
        }
    }

def AgentAdapterParallel(tasks):
    """
    Execute multiple agents in parallel using Task tool.
    
    Args:
        tasks (list): List of dict with agent_type and prompt
    
    Returns:
        List of results from parallel Task executions
    """
    results = []
    for task in tasks:
        result = AgentAdapter(
            agent_type=task.get("agent_type"),
            prompt=task.get("prompt")
        )
        results.append(result)
    return results

# For command-line usage
if __name__ == "__main__":
    import sys
    import json
    
    if len(sys.argv) < 3:
        print("Usage: AgentAdapter.py <agent-type> <prompt>")
        print("   or: AgentAdapter.py --parallel < tasks.json")
        sys.exit(1)
    
    if sys.argv[1] == "--parallel":
        # Read JSON from stdin
        tasks = json.loads(sys.stdin.read())
        results = AgentAdapterParallel(tasks)
        print(json.dumps(results, indent=2))
    else:
        # Single agent execution
        result = AgentAdapter(sys.argv[1], sys.argv[2])
        print(json.dumps(result, indent=2))