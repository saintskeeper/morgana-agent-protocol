# Client SDKs

**Source URL:** https://docs.anthropic.com/en/api/client-sdks **Fetch
Timestamp:** 2025-01-11T11:34:00-05:00

## Overview

Anthropic provides official client SDKs for multiple programming languages to
simplify API integration. These SDKs offer a consistent interface across
different languages while handling authentication, request formatting, and
response parsing.

## Supported Languages

### Production-Ready SDKs

- **Python** - Full feature support
- **TypeScript/JavaScript** - Node.js and browser compatible
- **Java** - Enterprise-grade implementation
- **Go** - High-performance client
- **Ruby** - Rails-friendly integration

### Beta SDKs

- **PHP** (beta) - Community-driven development

## Installation

### Python

```bash
pip install anthropic
```

### TypeScript/JavaScript

```bash
npm install @anthropic-ai/sdk
# or
yarn add @anthropic-ai/sdk
```

### Java

```xml
<dependency>
    <groupId>com.anthropic</groupId>
    <artifactId>anthropic-java</artifactId>
    <version>latest</version>
</dependency>
```

### Go

```bash
go get github.com/anthropic/anthropic-go
```

### Ruby

```ruby
gem 'anthropic'
```

### PHP (Beta)

```bash
composer require anthropic/anthropic-php
```

## Authentication

All SDKs support two authentication methods:

### 1. Environment Variable (Recommended)

Set the `ANTHROPIC_API_KEY` environment variable:

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

### 2. Direct Configuration

Pass the API key directly when initializing the client.

## Basic Usage Examples

### Python

```python
import anthropic

# Using environment variable
client = anthropic.Anthropic()

# Or with explicit API key
client = anthropic.Anthropic(api_key="sk-ant-...")

message = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "Hello, Claude"}
    ]
)

print(message.content[0].text)
```

### TypeScript/JavaScript

```typescript
import Anthropic from "@anthropic-ai/sdk";

// Using environment variable
const anthropic = new Anthropic();

// Or with explicit API key
const anthropic = new Anthropic({
  apiKey: "sk-ant-...",
});

const message = await anthropic.messages.create({
  model: "claude-sonnet-4-20250514",
  max_tokens: 1024,
  messages: [{ role: "user", content: "Hello, Claude" }],
});

console.log(message.content[0].text);
```

### Java

```java
import com.anthropic.AnthropicClient;
import com.anthropic.models.*;

// Using environment variable
AnthropicClient client = AnthropicClient.builder().build();

// Or with explicit API key
AnthropicClient client = AnthropicClient.builder()
    .apiKey("sk-ant-...")
    .build();

MessageResponse response = client.messages().create(
    CreateMessageRequest.builder()
        .model("claude-sonnet-4-20250514")
        .maxTokens(1024)
        .messages(List.of(
            Message.builder()
                .role("user")
                .content("Hello, Claude")
                .build()
        ))
        .build()
);

System.out.println(response.getContent().get(0).getText());
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "github.com/anthropic/anthropic-go"
)

func main() {
    // Using environment variable
    client := anthropic.NewClient()

    // Or with explicit API key
    client := anthropic.NewClient(
        anthropic.WithAPIKey("sk-ant-..."),
    )

    response, err := client.Messages.Create(context.Background(), &anthropic.MessageRequest{
        Model:     "claude-sonnet-4-20250514",
        MaxTokens: 1024,
        Messages: []anthropic.Message{
            {
                Role:    "user",
                Content: "Hello, Claude",
            },
        },
    })

    if err != nil {
        panic(err)
    }

    fmt.Println(response.Content[0].Text)
}
```

### Ruby

```ruby
require 'anthropic'

# Using environment variable
client = Anthropic::Client.new

# Or with explicit API key
client = Anthropic::Client.new(api_key: 'sk-ant-...')

response = client.messages.create(
  model: 'claude-sonnet-4-20250514',
  max_tokens: 1024,
  messages: [
    { role: 'user', content: 'Hello, Claude' }
  ]
)

puts response.content[0].text
```

## Available Models

All SDKs support the full range of Claude models:

### Claude 4 Series

- `claude-opus-4-1-20250805` - Most capable model
- `claude-sonnet-4-20250514` - Balanced performance
- `claude-4-latest` - Latest Claude 4 model alias

### Claude 3.7 Series

- `claude-3.7-opus-latest` - Top-tier Claude 3.7
- `claude-3.7-sonnet-latest` - Balanced Claude 3.7

### Claude 3.5 Series

- `claude-3-5-opus-latest` - Enhanced Claude 3.5
- `claude-3-5-sonnet-latest` - Fast Claude 3.5
- `claude-3-5-haiku-latest` - Efficient Claude 3.5

### Claude 3 Series

- `claude-3-opus-latest` - Powerful Claude 3
- `claude-3-sonnet-20240229` - Balanced Claude 3
- `claude-3-haiku-20240307` - Fast Claude 3

### Deprecated Models

Models marked as deprecated should be migrated to newer versions.

## Advanced Features

### Streaming Responses

All SDKs support streaming for real-time response generation:

```python
# Python streaming example
with client.messages.stream(
    model="claude-sonnet-4-20250514",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Tell me a story"}]
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)
```

### Tool Use

Define and use tools for enhanced functionality:

```typescript
// TypeScript tool use example
const response = await anthropic.messages.create({
  model: "claude-sonnet-4-20250514",
  max_tokens: 1024,
  tools: [
    {
      name: "get_weather",
      description: "Get current weather",
      input_schema: {
        type: "object",
        properties: {
          location: { type: "string" },
        },
      },
    },
  ],
  messages: [{ role: "user", content: "What's the weather in Paris?" }],
});
```

### Beta Features

Each SDK includes a `beta` namespace for experimental features:

```python
# Python beta features
from anthropic import Anthropic

client = Anthropic()

# Access beta features
response = client.beta.messages.create(
    # Beta-specific parameters
)
```

## Error Handling

All SDKs provide consistent error handling:

```python
# Python error handling
import anthropic
from anthropic import APIError, RateLimitError

client = anthropic.Anthropic()

try:
    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=1024,
        messages=[{"role": "user", "content": "Hello"}]
    )
except RateLimitError as e:
    print(f"Rate limit exceeded: {e}")
    # Implement retry logic
except APIError as e:
    print(f"API error: {e}")
    # Handle other API errors
```

## Partner Platforms

### Amazon Bedrock

Additional configuration required for AWS integration:

```python
# Python with Bedrock
import boto3
from anthropic_bedrock import AnthropicBedrock

bedrock_client = boto3.client('bedrock-runtime')
client = AnthropicBedrock(bedrock_client=bedrock_client)
```

### Google Cloud Vertex AI

Configuration for GCP integration:

```python
# Python with Vertex AI
from anthropic import AnthropicVertex

client = AnthropicVertex(
    project_id="your-project-id",
    region="us-central1"
)
```

## Best Practices

1. **Use Environment Variables**: Store API keys securely in environment
   variables
2. **Implement Retry Logic**: Handle rate limits and transient errors gracefully
3. **Monitor Usage**: Track token consumption using the usage object in
   responses
4. **Keep SDKs Updated**: Regularly update to get new features and bug fixes
5. **Use Type Hints**: Leverage TypeScript/Python type hints for better IDE
   support
6. **Handle Errors**: Implement comprehensive error handling for production apps
7. **Use Streaming**: For long responses, use streaming to improve user
   experience

## Migration Guide

When migrating between SDK versions:

1. Check the changelog for breaking changes
2. Update model names if deprecated
3. Adjust parameter names if changed
4. Test thoroughly in development environment
5. Monitor for any behavior changes

## SDK-Specific Features

### Python

- Async support with `AsyncAnthropic`
- Type hints for all methods
- Comprehensive docstrings

### TypeScript/JavaScript

- Full TypeScript definitions
- Browser and Node.js support
- Promise-based and async/await patterns

### Java

- Builder pattern for request construction
- Thread-safe client implementation
- Comprehensive JavaDoc

### Go

- Context support for cancellation
- Idiomatic Go error handling
- Efficient memory usage

## Support and Resources

- **Documentation**: Full API reference for each SDK
- **GitHub**: Source code and issue tracking
- **Discord**: Community support and discussions
- **Enterprise Support**: Available for production deployments

## Recommendations

1. Start with the SDK for your primary language
2. Use the latest stable version
3. Follow language-specific conventions
4. Implement proper error handling
5. Use streaming for interactive applications
6. Monitor API usage and costs
7. Keep API keys secure and rotate regularly
