#!/usr/bin/env python3
"""
Morgana REPL Bridge - Executes agents using native REPL capabilities
This runs INSIDE the REPL environment (Claude Code, Jupyter, etc)
"""
import json
import time
import os
import asyncio
from pathlib import Path
from typing import Dict, Any, Optional, Callable

class MorganaREPL:
    """REPL-native execution bridge for Morgana Agent Protocol"""
    
    def __init__(self, task_dir: str = "/tmp/morgana/tasks"):
        self.task_dir = Path(task_dir)
        self.pending_dir = self.task_dir / "pending"
        self.processing_dir = self.task_dir / "processing"
        self.complete_dir = self.task_dir / "complete"
        self.failed_dir = self.task_dir / "failed"
        
        # Create directories
        for dir in [self.pending_dir, self.processing_dir, 
                    self.complete_dir, self.failed_dir]:
            dir.mkdir(parents=True, exist_ok=True)
        
        # Task executor (will be set based on REPL environment)
        self.executor: Optional[Callable] = None
        self.running = False
        
    def detect_environment(self):
        """Detect which REPL we're running in"""
        # Check for Claude Code
        if 'Task' in globals():
            return 'claude'
        
        # Check for Jupyter/IPython
        try:
            get_ipython()
            return 'jupyter'
        except NameError:
            pass
        
        # Check for Google Colab
        try:
            import google.colab
            return 'colab'
        except ImportError:
            pass
        
        return 'generic'
    
    def setup_executor(self):
        """Set up the appropriate executor for this REPL"""
        env = self.detect_environment()
        
        if env == 'claude':
            # Claude Code - use native Task tool
            def claude_executor(agent_type: str, prompt: str, **kwargs):
                """Execute using Claude's native Task tool"""
                # Task is available in global scope in Claude Code
                result = Task(
                    subagent_type="general-purpose",
                    description=f"Execute {agent_type} agent",
                    prompt=prompt
                )
                return {"success": True, "output": result}
            
            self.executor = claude_executor
            print("✅ Claude Code environment detected - using native Task tool")
            
        elif env == 'jupyter':
            # Jupyter - could use IPython magic or subprocess
            def jupyter_executor(agent_type: str, prompt: str, **kwargs):
                """Execute in Jupyter environment"""
                # Could integrate with Jupyter AI extensions
                return {
                    "success": True, 
                    "output": f"[Jupyter] Would execute {agent_type} with prompt: {prompt[:100]}..."
                }
            
            self.executor = jupyter_executor
            print("📓 Jupyter environment detected")
            
        else:
            # Generic fallback
            def generic_executor(agent_type: str, prompt: str, **kwargs):
                """Generic executor for testing"""
                return {
                    "success": True,
                    "output": f"[Generic] Executed {agent_type} agent"
                }
            
            self.executor = generic_executor
            print("🔧 Generic environment - using mock executor")
    
    def process_task(self, task_file: Path) -> Dict[str, Any]:
        """Process a single task file"""
        try:
            # Read task
            with open(task_file, 'r') as f:
                task = json.load(f)
            
            task_id = task.get('id', task_file.stem)
            agent_type = task.get('agent_type')
            prompt = task.get('prompt')
            
            print(f"📋 Processing task {task_id}: {agent_type}")
            
            # Move to processing
            processing_file = self.processing_dir / task_file.name
            task_file.rename(processing_file)
            
            # Execute using appropriate executor
            result = self.executor(agent_type, prompt, **task.get('options', {}))
            
            # Add metadata
            result['task_id'] = task_id
            result['agent_type'] = agent_type
            result['completed_at'] = time.time()
            
            # Write result
            result_file = self.complete_dir / f"{task_id}.result.json"
            with open(result_file, 'w') as f:
                json.dump(result, f, indent=2)
            
            # Clean up processing file
            processing_file.unlink()
            
            print(f"✅ Completed task {task_id}")
            return result
            
        except Exception as e:
            print(f"❌ Error processing task: {e}")
            # Move to failed directory
            if task_file.exists():
                failed_file = self.failed_dir / task_file.name
                task_file.rename(failed_file)
            
            return {
                "success": False,
                "error": str(e),
                "task_id": task_id if 'task_id' in locals() else 'unknown'
            }
    
    async def poll_loop(self, interval: float = 0.5):
        """Async polling loop for new tasks"""
        self.running = True
        print(f"🔄 Starting Morgana REPL polling (interval: {interval}s)")
        print(f"📁 Watching: {self.pending_dir}")
        
        while self.running:
            try:
                # Check for pending tasks
                task_files = list(self.pending_dir.glob("*.json"))
                
                if task_files:
                    # Process oldest first
                    task_file = min(task_files, key=lambda f: f.stat().st_mtime)
                    self.process_task(task_file)
                
                await asyncio.sleep(interval)
                
            except KeyboardInterrupt:
                print("\n⏹️  Stopping Morgana REPL bridge")
                self.running = False
                break
            except Exception as e:
                print(f"❌ Polling error: {e}")
                await asyncio.sleep(interval)
    
    def start(self, interval: float = 0.5):
        """Start the REPL bridge (blocking)"""
        self.setup_executor()
        
        # Run async event loop
        try:
            asyncio.run(self.poll_loop(interval))
        except KeyboardInterrupt:
            print("\n👋 Morgana REPL bridge stopped")
    
    def start_background(self, interval: float = 0.5):
        """Start the REPL bridge in background (non-blocking)"""
        self.setup_executor()
        
        # Create background task
        import threading
        thread = threading.Thread(
            target=lambda: asyncio.run(self.poll_loop(interval)),
            daemon=True
        )
        thread.start()
        print("🚀 Morgana REPL bridge running in background")
        return thread
    
    def stop(self):
        """Stop the polling loop"""
        self.running = False

# Convenience functions for REPL usage
_morgana_instance = None

def morgana_start(background=False):
    """Start Morgana REPL bridge"""
    global _morgana_instance
    _morgana_instance = MorganaREPL()
    
    if background:
        return _morgana_instance.start_background()
    else:
        _morgana_instance.start()

def morgana_stop():
    """Stop Morgana REPL bridge"""
    global _morgana_instance
    if _morgana_instance:
        _morgana_instance.stop()
        print("✅ Morgana REPL bridge stopped")

def morgana_status():
    """Check Morgana REPL bridge status"""
    task_dir = Path("/tmp/morgana/tasks")
    
    if not task_dir.exists():
        print("❌ Morgana task directory not found")
        return
    
    pending = len(list((task_dir / "pending").glob("*.json")))
    complete = len(list((task_dir / "complete").glob("*.json")))
    failed = len(list((task_dir / "failed").glob("*.json")))
    
    print(f"📊 Morgana Status:")
    print(f"  Pending:  {pending}")
    print(f"  Complete: {complete}")
    print(f"  Failed:   {failed}")

# Auto-start in background if in Claude Code
if __name__ != "__main__" and 'Task' in globals():
    print("🎯 Claude Code detected - auto-starting Morgana REPL bridge")
    morgana_start(background=True)