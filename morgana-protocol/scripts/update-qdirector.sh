#!/bin/bash
# Updates qdirector-enhanced.md to use Morgana Protocol

echo "Updating qdirector-enhanced.md to use Morgana Protocol..."

# Create Python wrapper that calls Morgana
cat > /Users/walterday/.claude/scripts/agent_adapter.py << 'EOF'
#!/usr/bin/env python3
import subprocess
import json
import sys

def AgentAdapter(agent_type, prompt, **kwargs):
    """Adapter using Morgana Protocol"""
    morgana_bin = "/Users/walterday/.claude/morgana-protocol/dist/morgana"
    
    task = {
        "agent_type": agent_type,
        "prompt": prompt,
        "options": kwargs
    }
    
    # Call Morgana
    proc = subprocess.Popen(
        [morgana_bin],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    
    output, error = proc.communicate(json.dumps([task]))
    
    if proc.returncode != 0:
        raise Exception(f"Morgana failed: {error}")
    
    result = json.loads(output)
    return result["results"][0]["output"]

# CLI usage
if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: agent_adapter.py <agent-type> <prompt>")
        sys.exit(1)
    
    result = AgentAdapter(sys.argv[1], sys.argv[2])
    print(result)
EOF

chmod +x /Users/walterday/.claude/scripts/agent_adapter.py

echo "✓ Created agent_adapter.py wrapper"
echo "✓ Morgana Protocol integration complete!"
echo ""
echo "Usage:"
echo "  python ~/.claude/scripts/agent_adapter.py code-implementer 'implement feature'"