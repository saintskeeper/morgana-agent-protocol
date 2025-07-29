# Token-Efficient Mode Management

You have been asked to manage the token-efficient tool use feature. This beta
feature from Anthropic can reduce output tokens by 14-70% when using Claude
Sonnet 3.7.

## Current Configuration

Check the current status by running:

```bash
~/.claude/scripts/token-efficient-config.sh status
```

## Available Actions

1. **Enable Token-Efficient Mode**:

   - Run: `~/.claude/scripts/token-efficient-config.sh enable`
   - This will activate the beta header for API calls
   - Only works with Claude Sonnet 3.7 models

2. **Disable Token-Efficient Mode**:

   - Run: `~/.claude/scripts/token-efficient-config.sh disable`
   - Returns to standard API calls

3. **Check Compatibility**:
   - Verify the current model is compatible
   - Check for conflicting settings

## Important Notes

- This is a BETA feature - use with caution
- Not compatible with Claude 4 models (Opus/Sonnet)
- Cannot be used with `disable_parallel_tool_use`
- Test thoroughly before production use

## Implementation Details

The feature works by:

1. Adding a beta header to API requests
2. Using the token-efficient SDK when available
3. Monitoring and reporting token savings

Always inform the user about the current status and any changes made.
