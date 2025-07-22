#!/bin/bash
# Test the post-edit hooks manually

echo "🧪 Testing Claude hooks manually..."

# Test with a Go file
if [ -f "back-end-go/main.go" ]; then
    echo "Testing Go formatter..."
    ./.claude/hooks/post-edit.sh "back-end-go/main.go"
fi

# Test with this README
echo -e "\nTesting Markdown formatter..."
./.claude/hooks/post-edit.sh ".claude/README.md"

echo -e "\n✅ Hook test complete!"
