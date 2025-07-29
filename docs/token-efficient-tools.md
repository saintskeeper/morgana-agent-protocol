# Token-Efficient Tool Use Configuration

This feature enables Anthropic's beta token-efficient tool use, which can reduce
output tokens by an average of 14% (up to 70%) and reduce latency when using
Claude Sonnet 3.7.

## Quick Start

1. **Enable token-efficient mode:**

   ```bash
   ~/.claude/scripts/token-efficient-config.sh enable
   ```

2. **Check status:**

   ```bash
   ~/.claude/scripts/token-efficient-config.sh status
   ```

3. **Disable when needed:**
   ```bash
   ~/.claude/scripts/token-efficient-config.sh disable
   ```

## Configuration

The feature is configured in `~/.claude/settings.json` using environment
variables:

```json
{
  "env": {
    "CLAUDE_TOKEN_EFFICIENT_MODE": "true",
    "CLAUDE_BETA_HEADER": "token-efficient-tools-2025-02-19"
  }
}
```

### Configuration Options

- **CLAUDE_TOKEN_EFFICIENT_MODE**: Set to "true" to enable, "false" to disable
- **CLAUDE_BETA_HEADER**: The beta header value (don't change unless instructed
  by Anthropic)

## Environment Variables

You can also enable token-efficient mode using environment variables:

```bash
export CLAUDE_TOKEN_EFFICIENT_MODE=true
export CLAUDE_BETA_HEADER="token-efficient-tools-2025-02-19"
```

## Compatibility

⚠️ **Important Limitations:**

- Only works with Claude Sonnet 3.7 (`claude-3-7-sonnet-20250219`)
- Not compatible with Claude 4 models (Opus and Sonnet)
- Cannot be used with `disable_parallel_tool_use`
- This is a beta feature - test thoroughly before production use

## Monitoring Token Savings

When `showTokenSavings` is enabled and running in verbose mode:

```bash
CLAUDE_VERBOSE=1 claude "your prompt here"
```

You'll see token usage statistics comparing standard vs token-efficient modes.

## Troubleshooting

1. **Feature not working?**

   - Ensure you're using a compatible model
   - Check that the settings file is valid JSON
   - Verify the hook is being executed

2. **Unexpected behavior?**
   - Disable the feature and test again
   - This is a beta feature - report issues to Anthropic

## Best Practices

1. **Test First**: Always test with the beta feature in a development
   environment
2. **Monitor Performance**: Keep track of token savings and response quality
3. **Gradual Rollout**: Enable for specific use cases before broad adoption
4. **Feedback**: Provide feedback to Anthropic about your experience

## Integration with CI/CD

For automated environments, set the environment variable:

```yaml
# GitHub Actions example
env:
  CLAUDE_TOKEN_EFFICIENT_MODE: true
```

## References

- [Anthropic Documentation](https://docs.anthropic.com/en/docs/agents-and-tools/tool-use/token-efficient-tool-use)
- [Beta Features Guide](https://docs.anthropic.com/en/docs/beta-features)
