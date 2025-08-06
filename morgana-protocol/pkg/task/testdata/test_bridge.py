#!/usr/bin/env python3
"""Test bridge for integration testing"""
import json
import sys
import os
import time

def main():
    # Read test mode from environment
    test_mode = os.getenv("TEST_MODE", "success")
    
    try:
        # Read JSON input from stdin
        input_data = sys.stdin.read()
        data = json.loads(input_data)
        
        # Simulate different test scenarios
        if test_mode == "success":
            output = {
                "success": True,
                "output": f"Test output for {data['agent_type']}",
                "error": None
            }
        elif test_mode == "error":
            output = {
                "success": False,
                "output": None,
                "error": "Test error scenario"
            }
        elif test_mode == "timeout":
            # Sleep longer than typical timeout
            time.sleep(5)
            output = {
                "success": True,
                "output": "Should not reach here",
                "error": None
            }
        elif test_mode == "invalid_json":
            # Return invalid JSON
            print("This is not valid JSON")
            return
        elif test_mode == "crash":
            # Simulate a crash
            raise RuntimeError("Test crash scenario")
        else:
            output = {
                "success": True,
                "output": f"Unknown test mode: {test_mode}",
                "error": None
            }
    
    except Exception as e:
        output = {
            "success": False,
            "output": None,
            "error": str(e)
        }
    
    # Write JSON response to stdout
    print(json.dumps(output))

if __name__ == "__main__":
    main()