#!/usr/bin/env python3
"""
Morgana Sub-Agent Bridge - Manages sub-agent execution and communication
Runs inside Claude Code REPL with access to Task tool
"""
import json
import time
import asyncio
from pathlib import Path
from typing import Dict, Any, List, Optional
from dataclasses import dataclass, asdict
from enum import Enum

class AgentState(Enum):
    """Sub-agent execution states"""
    PENDING = "pending"
    RUNNING = "running"
    WAITING = "waiting"  # Waiting for user input
    PAUSED = "paused"
    COMPLETED = "completed"
    FAILED = "failed"

@dataclass
class SubAgentTask:
    """Represents a sub-agent task"""
    id: str
    parent_id: Optional[str]
    agent_type: str
    prompt: str
    state: AgentState
    result: Optional[Dict[str, Any]] = None
    children: List[str] = None
    metadata: Dict[str, Any] = None
    
    def __post_init__(self):
        if self.children is None:
            self.children = []
        if self.metadata is None:
            self.metadata = {}

class SubAgentOrchestrator:
    """Orchestrates sub-agent execution with message passing"""
    
    def __init__(self, message_dir: str = "/tmp/morgana/messages"):
        self.message_dir = Path(message_dir)
        self.inbox = self.message_dir / "inbox"
        self.outbox = self.message_dir / "outbox"
        self.state_dir = self.message_dir / "state"
        
        # Create directories
        for dir in [self.inbox, self.outbox, self.state_dir]:
            dir.mkdir(parents=True, exist_ok=True)
        
        # Task tracking
        self.tasks: Dict[str, SubAgentTask] = {}
        self.execution_stack: List[str] = []
        
    def create_sub_agent(self, 
                        agent_type: str, 
                        prompt: str,
                        parent_id: Optional[str] = None) -> SubAgentTask:
        """Create a new sub-agent task"""
        import uuid
        task_id = str(uuid.uuid4())[:8]
        
        task = SubAgentTask(
            id=task_id,
            parent_id=parent_id,
            agent_type=agent_type,
            prompt=prompt,
            state=AgentState.PENDING
        )
        
        self.tasks[task_id] = task
        
        # If this has a parent, add to parent's children
        if parent_id and parent_id in self.tasks:
            self.tasks[parent_id].children.append(task_id)
        
        # Publish creation event
        self.publish_event({
            "type": "agent_created",
            "task_id": task_id,
            "agent_type": agent_type,
            "parent_id": parent_id,
            "timestamp": time.time()
        })
        
        return task
    
    def execute_with_subagents(self, task: SubAgentTask) -> Dict[str, Any]:
        """Execute a task that can spawn sub-agents"""
        
        # Update state
        task.state = AgentState.RUNNING
        self.execution_stack.append(task.id)
        
        # Publish execution started event
        self.publish_event({
            "type": "agent_started",
            "task_id": task.id,
            "agent_type": task.agent_type,
            "timestamp": time.time()
        })
        
        try:
            # Special handling for director/orchestrator agents
            if task.agent_type in ["morgana-director", "orchestrator"]:
                result = self.execute_director(task)
            else:
                # Regular agent execution
                result = self.execute_agent(task)
            
            task.state = AgentState.COMPLETED
            task.result = result
            
            # Publish completion event
            self.publish_event({
                "type": "agent_completed",
                "task_id": task.id,
                "result": result.get("output", "")[:200],  # First 200 chars
                "timestamp": time.time()
            })
            
            return result
            
        except Exception as e:
            task.state = AgentState.FAILED
            task.result = {"error": str(e)}
            
            # Publish failure event
            self.publish_event({
                "type": "agent_failed",
                "task_id": task.id,
                "error": str(e),
                "timestamp": time.time()
            })
            
            raise
        
        finally:
            self.execution_stack.pop()
    
    def execute_director(self, task: SubAgentTask) -> Dict[str, Any]:
        """Execute a director agent that spawns sub-agents"""
        
        # Parse the prompt to identify sub-tasks
        # In real implementation, the Task tool would handle this
        
        # Example: Director spawns multiple sub-agents
        sub_tasks = []
        
        # Create sub-agent for planning
        planner = self.create_sub_agent(
            agent_type="sprint-planner",
            prompt=f"Plan the implementation for: {task.prompt}",
            parent_id=task.id
        )
        sub_tasks.append(planner)
        
        # Execute planner and get result
        planner_result = self.execute_agent(planner)
        
        # Based on plan, create implementation sub-agents
        implementer = self.create_sub_agent(
            agent_type="code-implementer",
            prompt=f"Implement based on plan: {planner_result.get('output', '')}",
            parent_id=task.id
        )
        sub_tasks.append(implementer)
        
        # Execute implementer
        impl_result = self.execute_agent(implementer)
        
        # Create test sub-agent
        tester = self.create_sub_agent(
            agent_type="test-specialist",
            prompt=f"Create tests for: {impl_result.get('output', '')}",
            parent_id=task.id
        )
        sub_tasks.append(tester)
        
        # Execute tester
        test_result = self.execute_agent(tester)
        
        # Aggregate results
        return {
            "success": True,
            "output": f"Completed orchestration with {len(sub_tasks)} sub-agents",
            "sub_agents": [t.id for t in sub_tasks],
            "results": {
                "planning": planner_result,
                "implementation": impl_result,
                "testing": test_result
            }
        }
    
    def execute_agent(self, task: SubAgentTask) -> Dict[str, Any]:
        """Execute a single agent using Task tool"""
        
        # Check for pause request
        if self.check_pause_request(task.id):
            task.state = AgentState.PAUSED
            self.publish_event({
                "type": "agent_paused",
                "task_id": task.id,
                "timestamp": time.time()
            })
            
            # Wait for resume
            self.wait_for_resume(task.id)
            task.state = AgentState.RUNNING
        
        # Check for user input request
        user_input = self.check_user_input(task.id)
        if user_input:
            task.prompt += f"\n\nUser input: {user_input}"
        
        # Load agent-specific prompt
        agent_prompt = self.load_agent_prompt(task.agent_type)
        
        # Execute using Task tool (in Claude Code)
        if 'Task' in globals():
            result = Task(
                subagent_type="general-purpose",
                description=f"Execute {task.agent_type} agent",
                prompt=f"{agent_prompt}\n\nTask: {task.prompt}"
            )
            return {"success": True, "output": result}
        else:
            # Mock for testing
            return {
                "success": True,
                "output": f"[Mock] Executed {task.agent_type}: {task.prompt[:100]}..."
            }
    
    def publish_event(self, event: Dict[str, Any]):
        """Publish event to outbox for monitoring"""
        event_file = self.outbox / f"{event['timestamp']:.6f}_{event['type']}.json"
        with open(event_file, 'w') as f:
            json.dump(event, f, indent=2)
    
    def check_pause_request(self, task_id: str) -> bool:
        """Check if there's a pause request for this task"""
        pause_file = self.inbox / f"pause_{task_id}.json"
        return pause_file.exists()
    
    def wait_for_resume(self, task_id: str):
        """Wait for resume signal"""
        resume_file = self.inbox / f"resume_{task_id}.json"
        pause_file = self.inbox / f"pause_{task_id}.json"
        
        while not resume_file.exists():
            time.sleep(0.5)
        
        # Clean up
        pause_file.unlink(missing_ok=True)
        resume_file.unlink(missing_ok=True)
    
    def check_user_input(self, task_id: str) -> Optional[str]:
        """Check for user input for this task"""
        input_file = self.inbox / f"input_{task_id}.json"
        
        if input_file.exists():
            with open(input_file, 'r') as f:
                data = json.load(f)
            input_file.unlink()
            return data.get("input", "")
        
        return None
    
    def load_agent_prompt(self, agent_type: str) -> str:
        """Load agent-specific prompt"""
        agent_file = Path.home() / ".claude" / "agents" / f"{agent_type}.md"
        
        if agent_file.exists():
            with open(agent_file, 'r') as f:
                return f.read()
        
        return f"You are the {agent_type} agent."
    
    def get_execution_tree(self) -> Dict[str, Any]:
        """Get the current execution tree"""
        def build_tree(task_id: str) -> Dict[str, Any]:
            task = self.tasks.get(task_id)
            if not task:
                return {}
            
            return {
                "id": task.id,
                "agent_type": task.agent_type,
                "state": task.state.value,
                "children": [build_tree(child_id) for child_id in task.children]
            }
        
        # Find root tasks (no parent)
        roots = [t for t in self.tasks.values() if t.parent_id is None]
        return {
            "execution_tree": [build_tree(r.id) for r in roots],
            "active_stack": self.execution_stack,
            "total_tasks": len(self.tasks)
        }

# Interactive Control Functions
def pause_agent(task_id: str):
    """Pause a specific agent"""
    inbox = Path("/tmp/morgana/messages/inbox")
    inbox.mkdir(parents=True, exist_ok=True)
    
    pause_file = inbox / f"pause_{task_id}.json"
    with open(pause_file, 'w') as f:
        json.dump({"action": "pause", "task_id": task_id}, f)
    
    print(f"⏸️  Paused agent {task_id}")

def resume_agent(task_id: str):
    """Resume a paused agent"""
    inbox = Path("/tmp/morgana/messages/inbox")
    inbox.mkdir(parents=True, exist_ok=True)
    
    resume_file = inbox / f"resume_{task_id}.json"
    with open(resume_file, 'w') as f:
        json.dump({"action": "resume", "task_id": task_id}, f)
    
    print(f"▶️  Resumed agent {task_id}")

def send_input(task_id: str, user_input: str):
    """Send input to an agent"""
    inbox = Path("/tmp/morgana/messages/inbox")
    inbox.mkdir(parents=True, exist_ok=True)
    
    input_file = inbox / f"input_{task_id}.json"
    with open(input_file, 'w') as f:
        json.dump({"input": user_input, "task_id": task_id}, f)
    
    print(f"📝 Sent input to agent {task_id}")

# Usage in Claude Code REPL
if __name__ != "__main__":
    print("🎯 Morgana Sub-Agent Bridge loaded")
    print("Functions available:")
    print("  - SubAgentOrchestrator() - Main orchestrator")
    print("  - pause_agent(task_id) - Pause an agent")
    print("  - resume_agent(task_id) - Resume an agent")
    print("  - send_input(task_id, text) - Send input to agent")