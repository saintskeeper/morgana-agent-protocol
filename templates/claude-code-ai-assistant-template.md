# Claude Code AI Assistant Master Template

## Overview
This template provides patterns for building AI assistants with Claude Code, based on speech-to-speech capabilities and multi-agent architectures.

## Core Architecture

### 1. Main Assistant Controller
The top-level orchestrator that handles user requests and delegates to specialized agents.

```python
# Instead of OpenAI's real-time API, use Claude's conversational flow
async def process_user_request(user_input):
    """Main entry point for processing user requests"""
    # Parse intent and extract parameters
    intent, params = parse_user_intent(user_input)
    
    # Route to appropriate agent
    result = await route_to_agent(intent, params)
    
    return format_response(result)
```

### 2. Tool Definitions

#### Time & Utilities
```python
def get_current_time():
    """Get current time with timezone awareness"""
    return datetime.now().strftime("%I:%M %p on %B %d, %Y")

def generate_random_number(min_val=1, max_val=100):
    """Generate random number in specified range"""
    return random.randint(min_val, max_val)
```

#### Browser Control Agent
```python
async def open_browser_urls(urls):
    """Open multiple URLs in browser tabs"""
    results = []
    for url in urls:
        # Claude Code equivalent: Use WebFetch for URL validation
        # Then use system commands to open browser
        result = await validate_and_open_url(url)
        results.append(result)
    return results

# Personalized URL mapping
BROWSER_SHORTCUTS = {
    "chatgpt": "https://chat.openai.com",
    "claude": "https://claude.ai",
    "gemini": "https://gemini.google.com",
    "hackernews": "https://news.ycombinator.com",
    "simonw": "https://simonwillison.net"
}
```

### 3. File Manipulation Agents

#### Create File Agent
```python
async def create_file_agent(filename, content_type, specifications):
    """Create files with AI-generated content"""
    # Use Claude to generate content based on specifications
    prompt = f"""
    Create a {content_type} file named {filename} with the following specifications:
    {specifications}
    
    Format the output appropriately for the file type.
    """
    
    # Generate content
    content = await generate_content(prompt)
    
    # Write file
    await write_file(filename, content)
    return f"Created {filename}"
```

#### Update File Agent
```python
async def update_file_agent(filename, modifications):
    """Update existing files with specified changes"""
    # Read current content
    current_content = await read_file(filename)
    
    # Generate update prompt
    prompt = f"""
    Update the following file content with these modifications:
    {modifications}
    
    Current content:
    {current_content}
    
    Provide the complete updated file content.
    """
    
    # Generate updated content
    updated_content = await generate_content(prompt)
    
    # Write updated file
    await write_file(filename, updated_content)
    return f"Updated {filename}"
```

#### Delete File Agent
```python
async def delete_file_agent(filename, force=False):
    """Delete files with optional confirmation"""
    if not force:
        # Request confirmation
        return f"Are you sure you want to delete {filename}? Confirm with force=True"
    
    # Perform deletion
    os.remove(filename)
    return f"Deleted {filename}"
```

### 4. Advanced Patterns

#### CSV Manipulation Agent
```python
async def csv_agent(action, filename, params):
    """Handle CSV file operations"""
    if action == "create":
        # Generate mock data
        data = generate_mock_csv_data(params.get("rows", 10))
        df = pd.DataFrame(data)
        df.to_csv(filename, index=False)
        
    elif action == "update":
        df = pd.read_csv(filename)
        
        # Handle various update operations
        if "delete_rows" in params:
            df = df.drop(params["delete_rows"])
        
        if "add_column" in params:
            df[params["add_column"]["name"]] = params["add_column"]["value"]
        
        if "add_rows" in params:
            # Use reasoning model for complex data generation
            new_data = await generate_complex_data(params["add_rows"])
            df = pd.concat([df, pd.DataFrame(new_data)])
        
        df.to_csv(filename, index=False)
    
    return f"CSV operation completed on {filename}"
```

#### Multi-Language Code Generation
```python
async def code_generation_agent(topic, languages):
    """Generate equivalent code examples in multiple languages"""
    results = {}
    
    for lang in languages:
        prompt = f"""
        Create a comprehensive {lang} file explaining {topic}.
        Include:
        - Detailed comments
        - Multiple examples
        - Best practices
        - Common patterns
        
        Format as a complete, runnable {lang} file.
        """
        
        filename = f"{topic.replace(' ', '_')}.{get_extension(lang)}"
        content = await generate_content(prompt)
        await write_file(filename, content)
        results[lang] = filename
    
    return results
```

### 5. Personalization System

```yaml
# personalization.yaml
user_preferences:
  name: "Dan"
  assistant_name: "Ada"
  
  browser_shortcuts:
    - name: "chatgpt"
      url: "https://chat.openai.com"
    - name: "claude"
      url: "https://claude.ai"
    - name: "gemini"
      url: "https://gemini.google.com"
    
  common_tasks:
    - "file_manipulation"
    - "code_generation"
    - "browser_automation"
    
  default_settings:
    confirm_deletions: true
    use_reasoning_model: false
    response_style: "concise"
```

### 6. Agent Orchestration

```python
class AgentOrchestrator:
    """Main orchestrator for multi-agent system"""
    
    def __init__(self):
        self.agents = {
            "browser": BrowserAgent(),
            "file": FileManipulationAgent(),
            "code": CodeGenerationAgent(),
            "csv": CSVAgent(),
            "system": SystemAgent()
        }
        self.personalization = load_personalization()
    
    async def process_request(self, request):
        """Process user request through appropriate agents"""
        # Parse request
        intent, entities = self.parse_request(request)
        
        # Route to agent
        agent = self.select_agent(intent)
        
        # Execute with timing
        start_time = time.time()
        result = await agent.execute(entities)
        execution_time = time.time() - start_time
        
        # Log performance
        self.log_execution(intent, execution_time)
        
        return result
```

### 7. Performance Tracking

```python
class PerformanceTracker:
    """Track agent execution times and optimize"""
    
    def __init__(self):
        self.execution_times = defaultdict(list)
    
    def log_execution(self, agent_name, task_type, duration):
        """Log execution time for analysis"""
        self.execution_times[agent_name].append({
            "task": task_type,
            "duration": duration,
            "timestamp": datetime.now()
        })
    
    def get_stats(self):
        """Get performance statistics"""
        stats = {}
        for agent, times in self.execution_times.items():
            durations = [t["duration"] for t in times]
            stats[agent] = {
                "avg_time": np.mean(durations),
                "min_time": min(durations),
                "max_time": max(durations),
                "total_calls": len(durations)
            }
        return stats
```

### 8. Error Handling & Recovery

```python
class ErrorHandler:
    """Graceful error handling for agent operations"""
    
    async def safe_execute(self, func, *args, **kwargs):
        """Execute function with error handling"""
        try:
            return await func(*args, **kwargs)
        except FileNotFoundError:
            return "File not found. Please check the filename."
        except PermissionError:
            return "Permission denied. Cannot perform this operation."
        except Exception as e:
            # Log error for debugging
            logging.error(f"Agent error: {str(e)}")
            return f"An error occurred: {str(e)}"
```

## Implementation Guidelines

### 1. Start Simple
Begin with basic file operations and gradually add complexity.

### 2. Use Claude's Strengths
- Natural language understanding for intent parsing
- Code generation for multi-language support
- Complex reasoning for data generation

### 3. Avoid Vendor Lock-in
- Use standard Python libraries
- Keep agent interfaces generic
- Store configuration in standard formats

### 4. Performance Optimization
- Implement async operations
- Cache frequently used data
- Monitor execution times

### 5. Security Considerations
- Validate all file paths
- Implement permission checks
- Sanitize user inputs
- Never expose sensitive data

## Usage Example

```python
# Initialize the assistant
assistant = AgentOrchestrator()

# Process user requests
await assistant.process_request("Create a CSV file with user analytics data")
await assistant.process_request("Update the file and add a premium column")
await assistant.process_request("Generate Python examples for loops and comprehensions")
await assistant.process_request("Open my favorite websites")
```

## Next Steps

1. Implement core agents
2. Add personalization features
3. Create performance monitoring
4. Build error recovery systems
5. Extend with custom agents for specific use cases

This template provides a foundation for building powerful AI assistants with Claude Code, focusing on practical file manipulation, code generation, and task automation without the vendor lock-in of proprietary APIs.