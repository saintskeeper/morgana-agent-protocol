# OpenAI SDK Compatibility

**Source URL:** https://docs.anthropic.com/en/api/openai-sdk **Fetch
Timestamp:** 2025-01-11T11:34:00-05:00

## Overview

Anthropic provides a compatibility layer that allows you to use the OpenAI SDK
to test Anthropic's API. This is primarily intended for model capability testing
rather than long-term production use.

## Getting Started

To use Anthropic's API with the OpenAI SDK:

1. Use an official OpenAI SDK
2. Update the base URL to: `https://api.anthropic.com/v1/`
3. Replace your OpenAI API key with your Anthropic API key
4. Use Claude model names instead of OpenAI model names

## Code Examples

### Python

```python
from openai import OpenAI

client = OpenAI(
    api_key="ANTHROPIC_API_KEY",
    base_url="https://api.anthropic.com/v1/"
)

response = client.chat.completions.create(
    model="claude-opus-4-1-20250805",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": "Who are you?"}
    ]
)

print(response.choices[0].message.content)
```

### TypeScript/JavaScript

```typescript
import OpenAI from "openai";

const openai = new OpenAI({
  apiKey: "ANTHROPIC_API_KEY",
  baseURL: "https://api.anthropic.com/v1/",
});

const completion = await openai.chat.completions.create({
  model: "claude-opus-4-1-20250805",
  messages: [
    { role: "system", content: "You are a helpful assistant." },
    { role: "user", content: "Who are you?" },
  ],
});

console.log(completion.choices[0].message.content);
```

## Supported Features

- Basic chat completion parameters
- Tool/function calling
- Streaming responses
- Token usage tracking
- Message history

## Limitations

- **System Messages**: System and developer messages are automatically hoisted
  to the beginning of the conversation
- **Audio Input**: Not supported
- **Prompt Caching**: Not available through this compatibility layer
- **OpenAI-Specific Parameters**: Many OpenAI-specific parameters are ignored
- **Tool Use JSON**: JSON output from tool use may not strictly follow the
  supplied schema

## Rate Limits

The compatibility layer follows Anthropic's standard rate limits for the
`/v1/messages` endpoint. These limits apply regardless of which SDK you're
using.

## Supported Parameters

### Required Parameters

- `model`: The Claude model to use
- `messages`: Array of message objects with role and content

### Optional Parameters

- `max_tokens`: Maximum tokens for the response
- `temperature`: Controls randomness (0.0-1.0)
- `top_p`: Nucleus sampling threshold
- `stop`: Stop sequences
- `stream`: Enable streaming responses
- `tools`: Define available tools/functions

## Best Practices

1. **Use for Testing**: This compatibility layer is ideal for quickly testing
   and comparing model capabilities
2. **Native API for Production**: For full feature access and production use,
   switch to the native Anthropic API
3. **Monitor Token Usage**: Keep track of token consumption through the usage
   object in responses
4. **Handle Errors**: Implement proper error handling for rate limits and API
   errors

## Migration Path

When ready to move from the compatibility layer to the native Anthropic API:

1. Install the Anthropic SDK for your language
2. Update your code to use Anthropic's message format
3. Take advantage of Anthropic-specific features like prompt caching
4. Adjust any tool/function definitions to match Anthropic's format

## Recommendation

While the OpenAI SDK compatibility layer provides a quick way to test
Anthropic's models, we recommend using the native Anthropic API and SDKs for
production applications to access the full range of features and optimizations.
