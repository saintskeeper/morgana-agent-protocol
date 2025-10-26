# Messages API

**Source URL:** https://docs.anthropic.com/en/api/messages **Fetch Timestamp:**
2025-01-11T11:34:00-05:00

## Endpoint Details

- **URL:** `https://api.anthropic.com/v1/messages`
- **Method:** POST
- **Authentication:** API key required in `x-api-key` header
- **Version Header:** `anthropic-version: 2023-06-01`

## Request Format

### Required Parameters

#### `model` (string, required)

Specifies which Claude model to use. Example values:

- `claude-sonnet-4-20250514`
- `claude-opus-4-1-20250805`
- `claude-haiku-3-20240307`

#### `messages` (array, required)

Array of message objects representing the conversation history. Each message
must have:

- `role`: Either "user" or "assistant"
- `content`: The message content (text or structured content)

#### `max_tokens` (integer, required)

Maximum number of tokens to generate in the response.

- Minimum value: 1
- Helps control response length and prevent excessive generation

### Optional Parameters

#### `temperature` (float, optional)

Controls response randomness.

- Range: 0.0 to 1.0
- Lower values = more deterministic
- Higher values = more creative/random

#### `system` (string, optional)

Provides context, instructions, or guidelines for the model's behavior.

#### `tools` (array, optional)

Defines functions or tools the model can use during the conversation.

#### `stop_sequences` (array, optional)

Custom text sequences that will halt generation when encountered.

#### `top_p` (float, optional)

Nucleus sampling threshold for controlling response diversity.

#### `top_k` (integer, optional)

Limits token selection to the top K most likely tokens.

## Response Format

```json
{
  "id": "msg_01XFDUDYJgAACzvnptvVoYEL",
  "type": "message",
  "role": "assistant",
  "content": [
    {
      "type": "text",
      "text": "Hello! How can I help you today?"
    }
  ],
  "model": "claude-sonnet-4-20250514",
  "stop_reason": "end_turn",
  "stop_sequence": null,
  "usage": {
    "input_tokens": 12,
    "output_tokens": 8
  }
}
```

### Response Fields

- `id`: Unique message identifier
- `type`: Always "message" for standard responses
- `role`: Always "assistant" for model responses
- `content`: Array of content blocks (text, tool use, etc.)
- `model`: The model used for generation
- `stop_reason`: Why generation stopped
  - `"end_turn"`: Natural completion
  - `"max_tokens"`: Hit token limit
  - `"stop_sequence"`: Encountered stop sequence
  - `"tool_use"`: Model invoked a tool
- `usage`: Token consumption details
  - `input_tokens`: Tokens in the request
  - `output_tokens`: Tokens in the response

## Code Examples

### cURL

```bash
curl https://api.anthropic.com/v1/messages \
     --header "x-api-key: $ANTHROPIC_API_KEY" \
     --header "anthropic-version: 2023-06-01" \
     --header "content-type: application/json" \
     --data '{
         "model": "claude-sonnet-4-20250514",
         "max_tokens": 1024,
         "messages": [
             {"role": "user", "content": "Hello, world"}
         ]
     }'
```

### Python

```python
import anthropic

client = anthropic.Anthropic(
    api_key="your-api-key"
)

message = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=1024,
    temperature=0.7,
    system="You are a helpful assistant.",
    messages=[
        {
            "role": "user",
            "content": "Explain quantum computing in simple terms."
        }
    ]
)

print(message.content[0].text)
```

### JavaScript/TypeScript

```javascript
import Anthropic from "@anthropic-ai/sdk";

const anthropic = new Anthropic({
  apiKey: "your-api-key",
});

const message = await anthropic.messages.create({
  model: "claude-sonnet-4-20250514",
  max_tokens: 1024,
  messages: [{ role: "user", content: "What is the capital of France?" }],
});

console.log(message.content[0].text);
```

## Multi-turn Conversations

To maintain context across multiple exchanges:

```python
messages = [
    {"role": "user", "content": "Hi, my name is Alice"},
    {"role": "assistant", "content": "Hello Alice! Nice to meet you."},
    {"role": "user", "content": "What's my name?"}
]

response = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=100,
    messages=messages
)
```

## Content Types

Messages support multiple content types:

### Text Content

```json
{
  "role": "user",
  "content": "Simple text message"
}
```

### Structured Content

```json
{
  "role": "user",
  "content": [
    {
      "type": "text",
      "text": "Describe this image:"
    },
    {
      "type": "image",
      "source": {
        "type": "base64",
        "media_type": "image/jpeg",
        "data": "base64_encoded_image_data"
      }
    }
  ]
}
```

## Error Handling

Common error responses:

### Rate Limit Error (429)

```json
{
  "error": {
    "type": "rate_limit_error",
    "message": "Rate limit exceeded"
  }
}
```

### Invalid Request (400)

```json
{
  "error": {
    "type": "invalid_request_error",
    "message": "max_tokens is required"
  }
}
```

### Authentication Error (401)

```json
{
  "error": {
    "type": "authentication_error",
    "message": "Invalid API key"
  }
}
```

## Best Practices

1. **Clear Prompts**: Use specific, well-structured prompts for best results
2. **System Context**: Leverage the system parameter to set consistent behavior
3. **Token Management**: Monitor token usage to optimize costs
4. **Error Handling**: Implement retry logic for rate limits and transient
   errors
5. **Message History**: Maintain conversation context for multi-turn
   interactions
6. **Stop Sequences**: Use custom stop sequences for controlled generation
7. **Temperature Tuning**: Adjust temperature based on use case (lower for
   factual, higher for creative)

## Rate Limits

Rate limits vary by model and API tier. Monitor the following headers in
responses:

- `x-ratelimit-limit`: Maximum requests allowed
- `x-ratelimit-remaining`: Requests remaining
- `x-ratelimit-reset`: When the limit resets

## Advanced Features

### Tool Use

Define and use tools for enhanced functionality:

```python
tools = [
    {
        "name": "get_weather",
        "description": "Get the current weather",
        "input_schema": {
            "type": "object",
            "properties": {
                "location": {"type": "string"}
            }
        }
    }
]

response = client.messages.create(
    model="claude-sonnet-4-20250514",
    max_tokens=1024,
    tools=tools,
    messages=[
        {"role": "user", "content": "What's the weather in Paris?"}
    ]
)
```

### Streaming Responses

Enable real-time streaming for long responses (see streaming documentation for
details).

## Recommendations

- Use the Messages API for all new integrations
- Leverage system prompts for consistent behavior
- Implement proper error handling and retries
- Monitor token usage for cost optimization
- Consider streaming for real-time applications
