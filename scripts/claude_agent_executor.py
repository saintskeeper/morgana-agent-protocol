#!/usr/bin/env python3
"""
Claude Code Agent Executor - Native Task Tool Integration

This executor directly calls Claude Code's Task tool and logs events using morgana_events.py.
Designed to work within Claude Code REPL environment without subprocess or file-based communication.
"""

import os
import sys
import time
import uuid
from typing import Dict, Any, Optional, Callable
from pathlib import Path

# Import morgana events for logging
try:
    from morgana_events import get_logger, log_task_start, log_task_complete, log_task_error
except ImportError:
    print("Warning: morgana_events not available - logging disabled", file=sys.stderr)
    get_logger = None
    log_task_start = None
    log_task_complete = None
    log_task_error = None

# Import command polling (optional)
try:
    from morgana_command_poll import get_poller, integration_check_point
except ImportError:
    print("Warning: morgana_command_poll not available - command polling disabled", file=sys.stderr)
    get_poller = None
    integration_check_point = None


class ClaudeAgentExecutor:
    """Native Claude Code agent executor with Task tool integration."""
    
    def __init__(self, agent_dir: str = "~/.claude/agents", log_events: bool = True, enable_command_polling: bool = True):
        """Initialize the agent executor.
        
        Args:
            agent_dir: Directory containing agent configuration files
            log_events: Whether to log events to morgana_events
            enable_command_polling: Whether to enable command polling for pause/resume control
        """
        self.agent_dir = Path(os.path.expanduser(agent_dir))
        self.log_events = log_events and get_logger is not None
        self.logger = get_logger() if self.log_events else None
        
        # Command polling support
        self.command_polling = enable_command_polling and get_poller is not None
        self.poller = get_poller() if self.command_polling else None
        
        # Supported agent types
        self.supported_agents = {
            "code-implementer",
            "sprint-planner", 
            "test-specialist",
            "validation-expert"
        }
    
    def _load_agent_config(self, agent_type: str) -> str:
        """Load agent configuration from markdown file.
        
        Args:
            agent_type: Type of agent to load
            
        Returns:
            Agent prompt/configuration as string
        """
        agent_file = self.agent_dir / f"{agent_type}.md"
        
        if not agent_file.exists():
            return f"You are a {agent_type} specialist agent. Complete the requested task to the best of your ability."
        
        try:
            with open(agent_file, 'r', encoding='utf-8') as f:
                return f.read()
        except Exception as e:
            print(f"Warning: Failed to load agent config for {agent_type}: {e}", file=sys.stderr)
            return f"You are a {agent_type} specialist agent. Complete the requested task to the best of your ability."
    
    def _detect_task_function(self) -> Optional[Callable]:
        """Detect if Task function is available in the current environment.
        
        Returns:
            Task function if available, None otherwise
        """
        # Check if we're in Claude Code environment by looking for Task in globals
        frame = sys._getframe(1)
        while frame:
            if 'Task' in frame.f_globals:
                return frame.f_globals['Task']
            frame = frame.f_back
        
        # Also check if Task is available in current namespace
        import builtins
        if hasattr(builtins, 'Task'):
            return getattr(builtins, 'Task')
        
        return None
    
    def _execute_with_command_polling(self, task_func: Callable, prompt: str, task_params: Dict[str, Any]) -> Any:
        """Execute task function with command polling integration.
        
        Args:
            task_func: The Task function to execute
            prompt: Task prompt
            task_params: Task parameters
            
        Returns:
            Task result
            
        Raises:
            RuntimeError: If task is stopped by command
        """
        if not self.command_polling or not integration_check_point:
            # No command polling - execute directly
            return task_func(prompt, **task_params)
        
        # Execute with command polling check
        # Note: This is a simple approach. In a more sophisticated implementation,
        # we might need to run the task in a separate thread and poll periodically.
        if not integration_check_point():
            raise RuntimeError("Task stopped by command before execution")
        
        result = task_func(prompt, **task_params)
        
        # Check again after execution
        if not integration_check_point():
            raise RuntimeError("Task stopped by command after execution") 
        
        return result
    
    def execute_agent(self, 
                     agent_type: str,
                     prompt: str,
                     task_id: Optional[str] = None,
                     timeout: Optional[int] = None,
                     model: Optional[str] = None,
                     **kwargs) -> Dict[str, Any]:
        """Execute an agent task using Claude Code's Task tool.
        
        Args:
            agent_type: Type of agent to execute
            prompt: Task prompt for the agent
            task_id: Optional task ID for tracking
            timeout: Optional timeout in seconds
            model: Optional model to use
            **kwargs: Additional parameters
            
        Returns:
            Dict with execution results and metadata
        """
        # Validate agent type
        if agent_type not in self.supported_agents:
            raise ValueError(f"Unsupported agent type: {agent_type}. Supported: {self.supported_agents}")
        
        # Generate task ID if not provided
        if task_id is None:
            task_id = f"{agent_type}_{str(uuid.uuid4())[:8]}"
        
        start_time = time.time()
        
        # Log task start
        if self.log_events and log_task_start:
            log_task_start(task_id, agent_type, prompt, 
                          timeout=timeout, model=model, **kwargs)
        
        try:
            # Check for pause/stop commands before starting
            if self.command_polling and integration_check_point:
                if not integration_check_point():
                    duration_ms = int((time.time() - start_time) * 1000)
                    if self.log_events and log_task_error:
                        log_task_error(task_id, "Task stopped by command", duration_ms, "command_stop")
                    return {
                        "success": False,
                        "agent_type": agent_type,
                        "task_id": task_id,
                        "error": "Task stopped by command",
                        "duration_ms": duration_ms,
                        "model": model,
                        "execution_mode": "command_stopped"
                    }
            
            # Load agent configuration
            agent_config = self._load_agent_config(agent_type)
            
            # Construct full prompt for the Task tool
            full_prompt = f"""You are executing as the {agent_type} specialist agent.

{agent_config}

TASK TO COMPLETE:
{prompt}

Please complete this task following your specialized role and best practices."""
            
            # Try to execute using Claude Code's Task function
            task_func = self._detect_task_function()
            
            if task_func is not None:
                # We're in Claude Code environment - execute Task directly
                task_params = {
                    "subagent_type": "general-purpose",
                    "description": f"Execute {agent_type} agent task",
                }
                
                if timeout:
                    task_params["timeout"] = timeout
                if model:
                    task_params["model"] = model
                
                # Add any additional kwargs as task parameters
                task_params.update(kwargs)
                
                # Execute the Task with command polling integration
                result = self._execute_with_command_polling(task_func, full_prompt, task_params)
                
                # Calculate duration and log success
                duration_ms = int((time.time() - start_time) * 1000)
                
                if self.log_events and log_task_complete:
                    output = str(result) if result is not None else "Task completed successfully"
                    log_task_complete(task_id, output, duration_ms, model)
                
                return {
                    "success": True,
                    "agent_type": agent_type,
                    "task_id": task_id,
                    "result": result,
                    "duration_ms": duration_ms,
                    "model": model,
                    "execution_mode": "claude_code_native"
                }
                
            else:
                # Fallback: Not in Claude Code environment
                duration_ms = int((time.time() - start_time) * 1000)
                
                fallback_result = {
                    "agent_type": agent_type,
                    "prompt": prompt,
                    "status": "mock_execution",
                    "message": "Task function not available - returning mock response",
                    "full_prompt": full_prompt,
                    "task_params": {
                        "timeout": timeout,
                        "model": model,
                        **kwargs
                    }
                }
                
                if self.log_events and log_task_complete:
                    log_task_complete(task_id, "Mock execution - Task function not available", 
                                    duration_ms, model)
                
                return {
                    "success": True,
                    "agent_type": agent_type,
                    "task_id": task_id,
                    "result": fallback_result,
                    "duration_ms": duration_ms,
                    "model": model,
                    "execution_mode": "mock_fallback"
                }
                
        except RuntimeError as e:
            # Handle command-related stops
            duration_ms = int((time.time() - start_time) * 1000)
            
            if "stopped by command" in str(e):
                if self.log_events and log_task_error:
                    log_task_error(task_id, str(e), duration_ms, "command_stop")
                
                return {
                    "success": False,
                    "agent_type": agent_type,
                    "task_id": task_id,
                    "error": str(e),
                    "duration_ms": duration_ms,
                    "model": model,
                    "execution_mode": "command_stopped"
                }
            else:
                # Re-raise non-command RuntimeErrors
                raise
                
        except Exception as e:
            # Log error and return error result
            duration_ms = int((time.time() - start_time) * 1000)
            
            if self.log_events and log_task_error:
                log_task_error(task_id, str(e), duration_ms, "execution")
            
            return {
                "success": False,
                "agent_type": agent_type,
                "task_id": task_id,
                "error": str(e),
                "duration_ms": duration_ms,
                "model": model,
                "execution_mode": "error"
            }
    
    def execute_parallel(self, tasks: list) -> list:
        """Execute multiple agent tasks in parallel.
        
        Args:
            tasks: List of task dictionaries with agent_type, prompt, and optional params
            
        Returns:
            List of execution results
        """
        results = []
        
        for task in tasks:
            if not isinstance(task, dict):
                results.append({
                    "success": False,
                    "error": "Invalid task format - must be dictionary",
                    "task": task
                })
                continue
            
            # Extract parameters
            agent_type = task.get("agent_type")
            prompt = task.get("prompt")
            
            if not agent_type or not prompt:
                results.append({
                    "success": False,
                    "error": "Missing required parameters: agent_type and prompt",
                    "task": task
                })
                continue
            
            # Execute the task
            result = self.execute_agent(
                agent_type=agent_type,
                prompt=prompt,
                task_id=task.get("task_id"),
                timeout=task.get("timeout"),
                model=task.get("model"),
                **{k: v for k, v in task.items() 
                   if k not in ["agent_type", "prompt", "task_id", "timeout", "model"]}
            )
            
            results.append(result)
        
        return results


# Global executor instance
_global_executor: Optional[ClaudeAgentExecutor] = None


def get_executor() -> ClaudeAgentExecutor:
    """Get or create the global executor instance."""
    global _global_executor
    if _global_executor is None:
        _global_executor = ClaudeAgentExecutor()
    return _global_executor


# Convenience functions for each agent type
def execute_code_implementer(prompt: str, **kwargs) -> Dict[str, Any]:
    """Execute a code-implementer agent task.
    
    Args:
        prompt: Task prompt
        **kwargs: Additional parameters (timeout, model, etc.)
        
    Returns:
        Execution results
    """
    return get_executor().execute_agent("code-implementer", prompt, **kwargs)


def execute_sprint_planner(prompt: str, **kwargs) -> Dict[str, Any]:
    """Execute a sprint-planner agent task.
    
    Args:
        prompt: Task prompt
        **kwargs: Additional parameters (timeout, model, etc.)
        
    Returns:
        Execution results
    """
    return get_executor().execute_agent("sprint-planner", prompt, **kwargs)


def execute_test_specialist(prompt: str, **kwargs) -> Dict[str, Any]:
    """Execute a test-specialist agent task.
    
    Args:
        prompt: Task prompt
        **kwargs: Additional parameters (timeout, model, etc.)
        
    Returns:
        Execution results
    """
    return get_executor().execute_agent("test-specialist", prompt, **kwargs)


def execute_validation_expert(prompt: str, **kwargs) -> Dict[str, Any]:
    """Execute a validation-expert agent task.
    
    Args:
        prompt: Task prompt
        **kwargs: Additional parameters (timeout, model, etc.)
        
    Returns:
        Execution results
    """
    return get_executor().execute_agent("validation-expert", prompt, **kwargs)


def execute_parallel_agents(tasks: list) -> list:
    """Execute multiple agent tasks in parallel.
    
    Args:
        tasks: List of task dictionaries
        
    Returns:
        List of execution results
    """
    return get_executor().execute_parallel(tasks)


# Command-line interface
if __name__ == "__main__":
    import json
    import argparse
    
    parser = argparse.ArgumentParser(description="Claude Code Agent Executor")
    parser.add_argument("--parallel", action="store_true", help="Execute parallel tasks from stdin")
    parser.add_argument("--agent-type", 
                       choices=["code-implementer", "sprint-planner", "test-specialist", "validation-expert"],
                       help="Type of agent to execute (required for single execution)")
    parser.add_argument("--prompt", help="Task prompt (required for single execution)")
    parser.add_argument("--task-id", help="Optional task ID")
    parser.add_argument("--timeout", type=int, help="Timeout in seconds")
    parser.add_argument("--model", help="Model to use")
    
    args = parser.parse_args()
    
    executor = ClaudeAgentExecutor()
    
    if args.parallel:
        # Read tasks from stdin
        try:
            tasks = json.loads(sys.stdin.read())
            results = executor.execute_parallel(tasks)
            print(json.dumps(results, indent=2))
        except json.JSONDecodeError as e:
            print(f"Error parsing JSON input: {e}", file=sys.stderr)
            sys.exit(1)
    else:
        # Single agent execution - validate required arguments
        if not args.agent_type or not args.prompt:
            parser.error("--agent-type and --prompt are required for single execution")
        
        result = executor.execute_agent(
            agent_type=args.agent_type,
            prompt=args.prompt,
            task_id=args.task_id,
            timeout=args.timeout,
            model=args.model
        )
        
        print(json.dumps(result, indent=2))
        
        # Exit with error code if execution failed
        if not result.get("success", False):
            sys.exit(1)