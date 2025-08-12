#!/usr/bin/env python3
"""
Morgana Command Utility

Simple command-line tool for sending pause/resume/stop commands
to running Morgana Protocol agents.
"""

import sys
import argparse
from pathlib import Path

# Import the command writing functions
try:
    from morgana_command_poll import (
        pause_execution, resume_execution, stop_execution, status_request,
        write_command, get_poller
    )
except ImportError:
    print("Error: morgana_command_poll module not found", file=sys.stderr)
    sys.exit(1)


def main():
    """Main command-line interface."""
    parser = argparse.ArgumentParser(
        description="Send control commands to Morgana Protocol agents",
        epilog="""
Examples:
  morgana-command.py pause         # Pause execution
  morgana-command.py resume        # Resume execution  
  morgana-command.py stop          # Stop execution
  morgana-command.py status        # Request status
  morgana-command.py check         # Check current state
        """,
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument(
        "command",
        choices=["pause", "resume", "stop", "status", "check"],
        help="Command to send or action to perform"
    )
    
    parser.add_argument(
        "--commands-file",
        default="/tmp/morgana/commands.txt",
        help="Path to commands file (default: /tmp/morgana/commands.txt)"
    )
    
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Verbose output"
    )
    
    args = parser.parse_args()
    
    if args.command == "check":
        # Check current state without sending command
        try:
            poller = get_poller()
            state = poller.get_state()
            print(f"Current state: {state.value}")
            
            commands_file = Path(args.commands_file)
            if commands_file.exists():
                print(f"Pending commands file exists: {commands_file}")
            else:
                print("No pending commands")
                
        except Exception as e:
            print(f"Error checking state: {e}", file=sys.stderr)
            sys.exit(1)
        return
    
    # Send command
    command_functions = {
        "pause": pause_execution,
        "resume": resume_execution,
        "stop": stop_execution,
        "status": status_request
    }
    
    func = command_functions[args.command]
    
    if args.verbose:
        print(f"Sending '{args.command}' command...")
    
    try:
        success = func()
        if success:
            if args.verbose:
                print(f"Command '{args.command}' sent successfully")
                print(f"Commands file: {args.commands_file}")
            else:
                print(f"✓ {args.command}")
        else:
            print(f"✗ Failed to send '{args.command}' command", file=sys.stderr)
            sys.exit(1)
    except Exception as e:
        print(f"Error sending command: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()