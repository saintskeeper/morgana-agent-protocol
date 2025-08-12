#!/usr/bin/env python3
"""
Morgana Command Polling System

Lightweight file-based command system for pause/resume control.
Provides optional flow control without complex IPC.
"""

import json
import os
import time
import threading
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict, Any, Optional, Callable, Set
from enum import Enum

try:
    from morgana_events import get_logger
except ImportError:
    get_logger = None


class CommandType(Enum):
    """Supported command types."""
    PAUSE = "pause"
    RESUME = "resume" 
    STOP = "stop"
    STATUS = "status"


class ExecutionState(Enum):
    """Agent execution states."""
    RUNNING = "running"
    PAUSED = "paused"
    STOPPED = "stopped"


class MorganaCommandPoller:
    """Lightweight command polling for pause/resume control."""
    
    def __init__(self, 
                 commands_file: str = "/tmp/morgana/commands.txt",
                 poll_interval: float = 0.5,
                 max_poll_rate: float = 10.0):
        """Initialize command poller.
        
        Args:
            commands_file: Path to command file
            poll_interval: Polling interval in seconds (default: 0.5s)
            max_poll_rate: Maximum polls per second for rate limiting
        """
        self.commands_file = Path(commands_file)
        self.poll_interval = poll_interval
        self.min_poll_interval = 1.0 / max_poll_rate
        self.last_poll_time = 0.0
        
        # Execution state
        self.state = ExecutionState.RUNNING
        self.state_lock = threading.Lock()
        
        # Event logging
        self.logger = get_logger() if get_logger else None
        
        # Callbacks for state changes
        self.callbacks: Dict[ExecutionState, Set[Callable]] = {
            ExecutionState.PAUSED: set(),
            ExecutionState.RUNNING: set(),
            ExecutionState.STOPPED: set()
        }
        
        # Ensure command directory exists
        self.commands_file.parent.mkdir(parents=True, exist_ok=True)
        
        # Log poller initialization
        if self.logger:
            self.logger.log_agent_event("poller_initialized", "command_poller",
                                      commands_file=str(self.commands_file),
                                      poll_interval=poll_interval)
    
    def _enforce_rate_limit(self) -> None:
        """Enforce rate limiting between polls."""
        current_time = time.time()
        elapsed = current_time - self.last_poll_time
        
        if elapsed < self.min_poll_interval:
            sleep_time = self.min_poll_interval - elapsed
            time.sleep(sleep_time)
        
        self.last_poll_time = time.time()
    
    def _read_commands(self) -> Optional[Dict[str, Any]]:
        """Read and parse commands from file.
        
        Returns:
            Command data dict or None if no commands/error
        """
        if not self.commands_file.exists():
            return None
        
        try:
            with open(self.commands_file, 'r', encoding='utf-8') as f:
                content = f.read().strip()
                if not content:
                    return None
                return json.loads(content)
        except (json.JSONDecodeError, IOError) as e:
            if self.logger:
                self.logger.log_agent_event("command_read_error", "command_poller",
                                          error=str(e), file=str(self.commands_file))
            return None
    
    def _clear_commands(self) -> None:
        """Clear the commands file after processing."""
        try:
            if self.commands_file.exists():
                self.commands_file.unlink()
        except OSError as e:
            if self.logger:
                self.logger.log_agent_event("command_clear_error", "command_poller",
                                          error=str(e), file=str(self.commands_file))
    
    def _process_command(self, command_data: Dict[str, Any]) -> None:
        """Process a single command.
        
        Args:
            command_data: Command data dictionary
        """
        command = command_data.get("command")
        
        if not command:
            return
        
        try:
            cmd_type = CommandType(command)
        except ValueError:
            if self.logger:
                self.logger.log_agent_event("invalid_command", "command_poller",
                                          command=command)
            return
        
        with self.state_lock:
            old_state = self.state
            
            if cmd_type == CommandType.PAUSE and self.state == ExecutionState.RUNNING:
                self.state = ExecutionState.PAUSED
            elif cmd_type == CommandType.RESUME and self.state == ExecutionState.PAUSED:
                self.state = ExecutionState.RUNNING
            elif cmd_type == CommandType.STOP:
                self.state = ExecutionState.STOPPED
            elif cmd_type == CommandType.STATUS:
                # Status command doesn't change state, just logs current state
                if self.logger:
                    self.logger.log_agent_event("status_request", "command_poller",
                                              current_state=self.state.value)
                return
            
            # Log state changes
            if old_state != self.state and self.logger:
                self.logger.log_agent_event("state_changed", "command_poller",
                                          old_state=old_state.value,
                                          new_state=self.state.value,
                                          command=command)
            
            # Execute callbacks for new state
            if old_state != self.state:
                for callback in self.callbacks.get(self.state, set()):
                    try:
                        callback(old_state, self.state)
                    except Exception as e:
                        if self.logger:
                            self.logger.log_agent_event("callback_error", "command_poller",
                                                      error=str(e), state=self.state.value)
    
    def poll_once(self) -> bool:
        """Poll for commands once.
        
        Returns:
            True if should continue polling, False to stop
        """
        self._enforce_rate_limit()
        
        # Read commands
        command_data = self._read_commands()
        
        if command_data:
            self._process_command(command_data)
            self._clear_commands()
        
        # Return False if stopped
        with self.state_lock:
            return self.state != ExecutionState.STOPPED
    
    def should_pause(self) -> bool:
        """Check if execution should be paused.
        
        Returns:
            True if execution should pause
        """
        with self.state_lock:
            return self.state == ExecutionState.PAUSED
    
    def should_stop(self) -> bool:
        """Check if execution should stop.
        
        Returns:
            True if execution should stop
        """
        with self.state_lock:
            return self.state == ExecutionState.STOPPED
    
    def is_running(self) -> bool:
        """Check if execution is running normally.
        
        Returns:
            True if execution is running
        """
        with self.state_lock:
            return self.state == ExecutionState.RUNNING
    
    def get_state(self) -> ExecutionState:
        """Get current execution state.
        
        Returns:
            Current execution state
        """
        with self.state_lock:
            return self.state
    
    def wait_while_paused(self, check_interval: float = 0.1) -> bool:
        """Block while execution is paused.
        
        Args:
            check_interval: How often to check state while waiting
            
        Returns:
            True if resumed, False if stopped
        """
        while self.should_pause():
            if self.should_stop():
                return False
            
            # Poll for commands while paused
            if not self.poll_once():
                return False
            
            time.sleep(check_interval)
        
        return not self.should_stop()
    
    def add_state_callback(self, state: ExecutionState, callback: Callable) -> None:
        """Add callback for state changes.
        
        Args:
            state: State to trigger callback on
            callback: Function to call when entering state
        """
        self.callbacks[state].add(callback)
    
    def remove_state_callback(self, state: ExecutionState, callback: Callable) -> None:
        """Remove state callback.
        
        Args:
            state: State to remove callback from
            callback: Function to remove
        """
        self.callbacks[state].discard(callback)


# Global poller instance
_global_poller: Optional[MorganaCommandPoller] = None


def get_poller() -> MorganaCommandPoller:
    """Get or create the global command poller instance."""
    global _global_poller
    if _global_poller is None:
        _global_poller = MorganaCommandPoller()
    return _global_poller


def check_should_pause() -> bool:
    """Quick check if execution should pause."""
    return get_poller().should_pause()


def check_should_stop() -> bool:
    """Quick check if execution should stop."""
    return get_poller().should_stop()


def poll_commands() -> bool:
    """Poll for commands once.
    
    Returns:
        True to continue, False to stop
    """
    return get_poller().poll_once()


def wait_if_paused() -> bool:
    """Wait while paused, checking for commands.
    
    Returns:
        True if resumed, False if stopped
    """
    return get_poller().wait_while_paused()


def integration_check_point() -> bool:
    """Integration checkpoint for use in long-running tasks.
    
    This function should be called periodically in long-running tasks.
    It handles pausing and checking for stop commands.
    
    Returns:
        True to continue execution, False to stop
    """
    poller = get_poller()
    
    # Poll for new commands
    if not poller.poll_once():
        return False
    
    # Handle pause state
    if poller.should_pause():
        if poller.logger:
            poller.logger.log_agent_event("execution_paused", "command_poller")
        
        result = poller.wait_while_paused()
        
        if result and poller.logger:
            poller.logger.log_agent_event("execution_resumed", "command_poller")
        
        return result
    
    return True


# Command writer utility functions
def write_command(command: str, **kwargs) -> bool:
    """Write a command to the command file.
    
    Args:
        command: Command to write (pause, resume, stop, status)
        **kwargs: Additional command parameters
        
    Returns:
        True if successful, False otherwise
    """
    commands_file = Path("/tmp/morgana/commands.txt")
    
    # Ensure directory exists
    commands_file.parent.mkdir(parents=True, exist_ok=True)
    
    command_data = {
        "command": command,
        "timestamp": datetime.now(timezone.utc).isoformat(),
        **kwargs
    }
    
    try:
        with open(commands_file, 'w', encoding='utf-8') as f:
            json.dump(command_data, f)
        return True
    except OSError:
        return False


def pause_execution() -> bool:
    """Send pause command."""
    return write_command("pause")


def resume_execution() -> bool:
    """Send resume command."""  
    return write_command("resume")


def stop_execution() -> bool:
    """Send stop command."""
    return write_command("stop")


def status_request() -> bool:
    """Send status request."""
    return write_command("status")


if __name__ == "__main__":
    import sys
    import argparse
    
    parser = argparse.ArgumentParser(description="Morgana Command Poller")
    subparsers = parser.add_subparsers(dest="action", help="Available actions")
    
    # Polling mode
    poll_parser = subparsers.add_parser("poll", help="Start command polling")
    poll_parser.add_argument("--interval", type=float, default=0.5, 
                            help="Polling interval in seconds")
    
    # Command sending
    cmd_parser = subparsers.add_parser("send", help="Send command")
    cmd_parser.add_argument("command", choices=["pause", "resume", "stop", "status"],
                           help="Command to send")
    
    # Demo mode
    demo_parser = subparsers.add_parser("demo", help="Run demo with simulated work")
    
    args = parser.parse_args()
    
    if not args.action:
        parser.print_help()
        sys.exit(1)
    
    if args.action == "poll":
        # Start polling loop
        poller = MorganaCommandPoller(poll_interval=args.interval)
        print(f"Starting command poller (interval: {args.interval}s)")
        print(f"Commands file: {poller.commands_file}")
        print("Press Ctrl+C to exit")
        
        try:
            while poller.poll_once():
                time.sleep(args.interval)
        except KeyboardInterrupt:
            print("\nPoller stopped by user")
        
    elif args.action == "send":
        success = write_command(args.command)
        if success:
            print(f"Command '{args.command}' sent successfully")
        else:
            print(f"Failed to send command '{args.command}'")
            sys.exit(1)
    
    elif args.action == "demo":
        # Demo with simulated work
        print("Running demo with simulated work...")
        print("Try sending commands: python morgana_command_poll.py send pause")
        
        poller = MorganaCommandPoller()
        
        for i in range(100):
            print(f"Work iteration {i+1}")
            
            # Check for pause/stop commands
            if not integration_check_point():
                print("Execution stopped by command")
                break
            
            # Simulate work
            time.sleep(1)
        
        print("Demo completed")