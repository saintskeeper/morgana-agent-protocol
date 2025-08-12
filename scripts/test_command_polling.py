#!/usr/bin/env python3
"""
Test script for Morgana Command Polling System

Demonstrates the lightweight command polling for pause/resume control.
"""

import time
import sys
from pathlib import Path

# Add current directory to Python path
sys.path.insert(0, str(Path(__file__).parent))

from morgana_command_poll import (
    MorganaCommandPoller, 
    integration_check_point,
    pause_execution, 
    resume_execution, 
    stop_execution,
    status_request
)


def test_basic_polling():
    """Test basic command polling functionality."""
    print("=== Testing Basic Command Polling ===")
    
    poller = MorganaCommandPoller()
    print(f"Poller initialized with commands file: {poller.commands_file}")
    print(f"Initial state: {poller.get_state().value}")
    
    # Test polling once with no commands
    print("\n1. Testing poll with no commands...")
    result = poller.poll_once()
    print(f"Poll result: {result} (should be True)")
    print(f"State after poll: {poller.get_state().value}")
    
    # Test sending a pause command
    print("\n2. Sending pause command...")
    pause_success = pause_execution()
    print(f"Pause command sent: {pause_success}")
    
    # Poll to process the pause command
    result = poller.poll_once()
    print(f"Poll result after pause: {result}")
    print(f"State after pause: {poller.get_state().value}")
    
    # Test resume command
    print("\n3. Sending resume command...")
    resume_success = resume_execution()
    print(f"Resume command sent: {resume_success}")
    
    # Poll to process the resume command
    result = poller.poll_once()
    print(f"Poll result after resume: {result}")
    print(f"State after resume: {poller.get_state().value}")
    
    # Test stop command
    print("\n4. Sending stop command...")
    stop_success = stop_execution()
    print(f"Stop command sent: {stop_success}")
    
    # Poll to process the stop command
    result = poller.poll_once()
    print(f"Poll result after stop: {result} (should be False)")
    print(f"State after stop: {poller.get_state().value}")


def test_integration_checkpoint():
    """Test the integration checkpoint function."""
    print("\n=== Testing Integration Checkpoint ===")
    
    # Reset state by creating new poller
    from morgana_command_poll import _global_poller
    if _global_poller:
        _global_poller.state = _global_poller.state.RUNNING
    
    print("1. Testing checkpoint with running state...")
    result = integration_check_point()
    print(f"Checkpoint result: {result} (should be True)")
    
    # Send pause command and test checkpoint
    print("\n2. Sending pause command and testing checkpoint...")
    pause_execution()
    
    # This will poll and then wait while paused
    print("Checkpoint will poll for pause command and then simulate waiting...")
    print("(This test will timeout after a few seconds since no resume is sent)")
    
    start_time = time.time()
    result = integration_check_point()
    elapsed = time.time() - start_time
    
    print(f"Checkpoint result: {result}")
    print(f"Time elapsed: {elapsed:.2f}s")


def test_simulated_work():
    """Test with simulated work that can be paused/resumed."""
    print("\n=== Testing Simulated Work with Command Polling ===")
    
    # Reset state
    from morgana_command_poll import _global_poller
    if _global_poller:
        _global_poller.state = _global_poller.state.RUNNING
    
    print("Starting simulated work loop...")
    print("You can test commands by running in another terminal:")
    print("  python morgana-command.py pause")
    print("  python morgana-command.py resume")
    print("  python morgana-command.py stop")
    print("  python morgana-command.py status")
    print()
    
    for i in range(20):
        print(f"Work iteration {i+1}/20")
        
        # Integration checkpoint - handles pause/resume/stop
        if not integration_check_point():
            print("Work stopped by command!")
            break
        
        # Simulate some work
        time.sleep(1)
    
    print("Work completed or stopped.")


def main():
    """Main test function."""
    print("Morgana Command Polling System Test")
    print("====================================")
    
    if len(sys.argv) > 1:
        test_mode = sys.argv[1]
    else:
        test_mode = "basic"
    
    if test_mode == "basic":
        test_basic_polling()
    elif test_mode == "checkpoint":
        test_integration_checkpoint()
    elif test_mode == "work":
        test_simulated_work()
    else:
        print(f"Unknown test mode: {test_mode}")
        print("Available modes: basic, checkpoint, work")
        sys.exit(1)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nTest interrupted by user")
    except Exception as e:
        print(f"\nTest failed with error: {e}")
        raise