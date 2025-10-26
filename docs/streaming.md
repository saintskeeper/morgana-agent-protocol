# Streaming Messages

**Source URL:** https://docs.anthropic.com/en/api/streaming **Fetch Timestamp:**
2025-01-11T11:34:00-05:00

## Overview

Anthropic's streaming API allows you to receive responses incrementally using
server-sent events (SSE). This enables real-time, chunk-by-chunk delivery of
responses, improving user experience for long-form content generation.

## Key Benefits

- **Real-time feedback**: Users see responses as they're generated
- **Better UX**: No waiting for complete response before displaying content
- **Granular control**: Process response chunks as they arrive
- **Efficient handling**: Parse and display content progressively

## How Streaming Works

When streaming is enabled, the API sends responses as a series of server-sent
events. Each event contains a portion of the response, allowing you to update
your UI progressively.

### Event Flow Sequence

1. **`message_start`**: Initializes an empty message object
2. **Content block events**: Multiple events for content generation
   - `content_block_start`: Begins a new content block
   - `content_block_delta`: Sends incremental content updates
   - `content_block_stop`: Marks the end of a content block
3. **`message_delta`**: Updates to top-level message fields (usage, stop_reason)
4. **`message_stop`**: Signals the complete end of the response

## Supported Content Types

### Text Streaming

Standard text responses are streamed character by character or in small chunks.

### Tool Use Streaming

When using tools, the API streams:

- Tool name and ID
- Input parameters as they're generated
- Fine-grained parameter updates

### Extended Thinking

For models with extended thinking capabilities:

- Step-by-step reasoning process
- Intermediate thoughts and analysis
- Final conclusions

## Implementation Examples

### Python SDK

#### Basic Streaming

```python
import anthropic

client = anthropic.Anthropic()

# Using context manager for automatic cleanup
with client.messages.stream(
    model="claude-opus-4-1-20250805",
    messages=[{"role": "user", "content": "Write a short story"}],
    max_tokens=1024
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)
```

#### Handling Different Event Types

```python
with client.messages.stream(
    model="claude-opus-4-1-20250805",
    messages=[{"role": "user", "content": "Analyze this data"}],
    max_tokens=1024
) as stream:
    for event in stream:
        if event.type == "content_block_start":
            print("\n[New content block]")
        elif event.type == "content_block_delta":
            print(event.delta.text, end="", flush=True)
        elif event.type == "message_stop":
            print("\n[Response complete]")
```

#### Async Streaming

```python
import asyncio
from anthropic import AsyncAnthropic

async def stream_response():
    client = AsyncAnthropic()

    async with client.messages.stream(
        model="claude-opus-4-1-20250805",
        messages=[{"role": "user", "content": "Explain quantum computing"}],
        max_tokens=1024
    ) as stream:
        async for text in stream.text_stream:
            print(text, end="", flush=True)

asyncio.run(stream_response())
```

### TypeScript/JavaScript SDK

#### Basic Streaming

```typescript
import Anthropic from "@anthropic-ai/sdk";

const anthropic = new Anthropic();

const stream = await anthropic.messages.stream({
  model: "claude-opus-4-1-20250805",
  messages: [{ role: "user", content: "Tell me a joke" }],
  max_tokens: 1024,
});

stream.on("text", (text) => {
  process.stdout.write(text);
});

stream.on("end", () => {
  console.log("\n[Stream complete]");
});
```

#### Handling Events

```typescript
const stream = await anthropic.messages.stream({
  model: "claude-opus-4-1-20250805",
  messages: [{ role: "user", content: "Analyze this problem" }],
  max_tokens: 1024,
});

stream.on("message_start", (message) => {
  console.log("Starting new message:", message.id);
});

stream.on("content_block_delta", (delta) => {
  if (delta.type === "text_delta") {
    process.stdout.write(delta.text);
  }
});

stream.on("message_stop", () => {
  console.log("\nMessage complete");
});
```

#### Async Iterator Pattern

```typescript
const stream = await anthropic.messages.stream({
  model: "claude-opus-4-1-20250805",
  messages: [{ role: "user", content: "Write a poem" }],
  max_tokens: 1024,
});

for await (const chunk of stream) {
  if (chunk.type === "content_block_delta") {
    process.stdout.write(chunk.delta.text);
  }
}
```

### Raw SSE Format

If not using an SDK, you can consume the raw SSE stream:

```bash
curl https://api.anthropic.com/v1/messages \
  --header "x-api-key: $ANTHROPIC_API_KEY" \
  --header "anthropic-version: 2023-06-01" \
  --header "content-type: application/json" \
  --data '{
    "model": "claude-opus-4-1-20250805",
    "messages": [{"role": "user", "content": "Hello"}],
    "max_tokens": 1024,
    "stream": true
  }'
```

Example SSE events:

```
event: message_start
data: {"type":"message_start","message":{"id":"msg_123","type":"message","role":"assistant","content":[],"model":"claude-opus-4-1-20250805","usage":{"input_tokens":10,"output_tokens":0}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"! How"}}

event: content_block_stop
data: {"type":"content_block_stop","index":0}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn","usage":{"output_tokens":10}}}

event: message_stop
data: {"type":"message_stop"}
```

## Error Recovery Best Practices

When streaming connections are interrupted, you can resume from where you left
off:

### 1. Capture Partial Response

```python
partial_response = ""
try:
    with client.messages.stream(...) as stream:
        for text in stream.text_stream:
            partial_response += text
            print(text, end="", flush=True)
except Exception as e:
    print(f"\nStream interrupted: {e}")
    # partial_response contains what was received
```

### 2. Construct Continuation Request

```python
# Resume from interruption
continuation_messages = [
    {"role": "user", "content": original_prompt},
    {"role": "assistant", "content": partial_response},
    {"role": "user", "content": "Please continue from where you left off"}
]

# Start new stream with context
with client.messages.stream(
    model="claude-opus-4-1-20250805",
    messages=continuation_messages,
    max_tokens=1024
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)
```

## Tool Use with Streaming

When using tools, streaming provides granular updates:

```python
tools = [{
    "name": "calculator",
    "description": "Performs basic math operations",
    "input_schema": {
        "type": "object",
        "properties": {
            "operation": {"type": "string"},
            "a": {"type": "number"},
            "b": {"type": "number"}
        }
    }
}]

with client.messages.stream(
    model="claude-opus-4-1-20250805",
    messages=[{"role": "user", "content": "What is 15 * 23?"}],
    tools=tools,
    max_tokens=1024
) as stream:
    for event in stream:
        if event.type == "content_block_start" and hasattr(event.content_block, 'name'):
            print(f"Using tool: {event.content_block.name}")
        elif event.type == "content_block_delta":
            if hasattr(event.delta, 'partial_json'):
                print(f"Tool input: {event.delta.partial_json}")
```

## Performance Considerations

### Buffer Management

- Process chunks immediately to avoid memory buildup
- Use appropriate buffer sizes for your application
- Consider implementing backpressure for slow consumers

### Connection Handling

- Implement timeout handling for stalled streams
- Use keep-alive mechanisms for long-running streams
- Handle network interruptions gracefully

### UI Updates

- Batch UI updates to avoid excessive re-renders
- Use debouncing for rapid chunk arrivals
- Implement smooth scrolling for content display

## Best Practices

1. **Use Official SDKs**: They handle SSE parsing, reconnection, and error
   recovery
2. **Handle All Event Types**: Don't assume only text events will arrive
3. **Implement Error Recovery**: Network issues can interrupt streams
4. **Monitor Token Usage**: Track usage in `message_delta` events
5. **Clean Up Resources**: Ensure streams are properly closed
6. **Test Edge Cases**: Handle empty responses, early termination
7. **Optimize UI Updates**: Balance responsiveness with performance

## Common Patterns

### Progress Indication

```typescript
let totalTokens = 0;

stream.on("message_delta", (delta) => {
  if (delta.usage) {
    totalTokens = delta.usage.output_tokens;
    updateProgressBar(totalTokens / maxTokens);
  }
});
```

### Content Accumulation

```python
full_response = ""
with client.messages.stream(...) as stream:
    for text in stream.text_stream:
        full_response += text
        # Process incrementally
        update_ui(text)

    # Process complete response
    save_to_database(full_response)
```

### Multi-Modal Streaming

```python
with client.messages.stream(...) as stream:
    for event in stream:
        if event.type == "content_block_start":
            if event.content_block.type == "text":
                prepare_text_display()
            elif event.content_block.type == "tool_use":
                prepare_tool_display()
```

## Troubleshooting

### Stream Not Starting

- Verify `stream: true` is set in the request
- Check API key and permissions
- Ensure model supports streaming

### Incomplete Responses

- Check for `message_stop` event
- Monitor `stop_reason` in message_delta
- Verify max_tokens setting

### Connection Drops

- Implement exponential backoff for retries
- Save partial responses for recovery
- Use connection pooling for stability

## Recommendations

1. **Always use streaming for long responses** - Improves perceived performance
2. **Implement proper error handling** - Networks are unreliable
3. **Test with various response lengths** - Ensure UI handles all cases
4. **Monitor streaming performance** - Track latency and throughput
5. **Use streaming-specific features** - Like partial JSON for tool use
6. **Consider user experience** - Show loading states and progress
7. **Optimize for your use case** - Balance latency vs. throughput
