# Beta Features Guide

## What Are Beta Features?

Beta features are experimental capabilities that Anthropic releases for early
testing and feedback. They provide access to cutting-edge functionality before
it becomes generally available.

### Key Characteristics of Beta Features:

1. **Experimental Nature**

   - May change or be discontinued
   - Not guaranteed for production use
   - Subject to modifications based on feedback

2. **Opt-in Requirement**

   - Must be explicitly enabled
   - Requires special headers or configuration
   - Can be disabled at any time

3. **Limited Compatibility**
   - Only work with specific models
   - May have additional restrictions
   - Some features may conflict with others

## Currently Enabled Beta Features

### Token-Efficient Tool Use

**Status**: BETA **Header**: `token-efficient-tools-2025-02-19` **Compatible
Models**: Claude 3.7 Sonnet only

#### What It Does:

- Reduces output tokens by 14-70% on average
- Improves response latency
- Maintains output quality while being more concise
- Optimizes tool use patterns for efficiency

#### What It Means for You:

- **Lower Costs**: Fewer tokens = reduced API costs
- **Faster Responses**: Less data to generate and transfer
- **Same Quality**: Optimizations don't compromise accuracy
- **Automatic**: Works transparently when enabled

#### Important Limitations:

- **Model Specific**: Only works with Claude 3.7 Sonnet
- **No-op on Claude 4**: Claude 4 models ignore this header (works normally)
- **Incompatible with**: `disable_parallel_tool_use` option
- **Beta Status**: May change or be discontinued

## How Beta Features Work

### 1. Detection and Activation

```bash
# When you enable token-efficient mode:
~/.claude/scripts/token-efficient-config.sh enable

# The system:
1. Sets environment variables in settings.json
2. Hooks intercept API calls
3. Checks model compatibility
4. Applies beta headers if compatible
5. Falls back gracefully if not
```

### 2. Graceful Degradation

- **Compatible Model**: Feature activates, benefits applied
- **Incompatible Model**: Feature ignored, normal operation
- **Conflicts Detected**: Feature disabled, warning logged
- **Errors**: Never breaks functionality, always safe

### 3. Monitoring and Feedback

```bash
# Enable verbose mode to see beta feature activity:
CLAUDE_VERBOSE=1 claude "your prompt"

# You'll see:
[Token-Efficient Mode] Activated with header: token-efficient-tools-2025-02-19
[Token-Efficient Mode] Model: claude-3-7-sonnet-20250219
```

## What This Means for Your Workflow

### When Beta Features Are Active:

1. **Development Work**

   - Continue working normally
   - Enjoy performance benefits automatically
   - No code changes required

2. **Cost Optimization**

   - Monitor token usage reduction
   - Track API cost savings
   - Evaluate efficiency gains

3. **Quality Assurance**
   - Test outputs remain consistent
   - Verify no functionality loss
   - Report any issues to Anthropic

### Best Practices:

1. **Test First**

   - Enable in development before production
   - Compare outputs with/without beta features
   - Ensure compatibility with your use cases

2. **Monitor Changes**

   - Watch Anthropic's announcements
   - Be prepared for feature evolution
   - Have fallback plans if features change

3. **Provide Feedback**
   - Report issues to Anthropic
   - Share performance improvements
   - Suggest enhancements

## Risk Management

### Low Risk Factors:

- ✅ No-op design prevents breaking changes
- ✅ Automatic compatibility checking
- ✅ Easy enable/disable mechanism
- ✅ Graceful fallbacks

### Considerations:

- ⚠️ Features may be removed/changed
- ⚠️ Behavior might evolve
- ⚠️ Not suitable for critical production without testing

## Future Beta Features

Anthropic regularly releases new beta features. Our configuration system is
designed to:

- Easily add new beta features
- Manage multiple features simultaneously
- Provide clear status and compatibility info
- Ensure safe operation regardless of changes

## Quick Reference

### Check Beta Feature Status:

```bash
~/.claude/scripts/token-efficient-config.sh status
```

### Enable/Disable:

```bash
# Enable
~/.claude/scripts/token-efficient-config.sh enable

# Disable
~/.claude/scripts/token-efficient-config.sh disable
```

### Troubleshooting:

- Feature not working? Check model compatibility
- Unexpected behavior? Disable and compare
- Need help? Check verbose logs first

Remember: Beta features are optional enhancements. Your Claude Code experience
remains fully functional with or without them.
