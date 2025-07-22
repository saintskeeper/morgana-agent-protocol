# Claude Commands 101: Creating Custom Commands for Claude Code

## Overview

Claude Code allows you to create custom slash commands using markdown files in a `.claude/commands/` directory. These commands become reusable prompts that you can invoke with `/command:name` syntax.

## Basic Structure

### Directory Layout
```
your-project/
├── .claude/
│   ├── commands/
│   │   ├── your-command.md
│   │   ├── subfolder/
│   │   │   └── nested-command.md
│   │   └── ...
│   └── settings.json (optional)
└── ...
```

### Command File Format
Each command is a simple markdown file containing the prompt you want Claude to execute:

```markdown
# Example Command File: .claude/commands/review.md

Review this code for:
1. Security vulnerabilities
2. Performance issues
3. Code style and best practices
4. Potential bugs

Provide specific suggestions for improvement.
```

## Command Types

### 1. Simple Static Commands
Basic prompts that don't change:

```markdown
# .claude/commands/explain.md
Explain this code in simple terms, including:
- What it does
- How it works
- Any notable patterns or techniques used
```

### 2. Commands with Arguments
Use `$ARGUMENTS` to accept user input:

```markdown
# .claude/commands/fix-issue.md
Find and fix issue #$ARGUMENTS. Follow these steps:
1. Understand the issue described in the ticket
2. Locate the relevant code in our codebase
3. Implement a solution that addresses the root cause
4. Add appropriate tests
5. Prepare a concise PR description
```

Usage: `/fix-issue 123` (replaces `$ARGUMENTS` with "123")

### 3. Complex Orchestration Commands
Multi-phase commands like the infinite agentic loop:

```markdown
# .claude/commands/infinite.md

Think deeply about this infinite generation task. You are about to embark on a sophisticated iterative creation process.

## PHASE 1: SPECIFICATION ANALYSIS
Read and deeply understand the specification file at spec_file. This file defines:
- Content requirements and format
- Quality standards
- Naming conventions
- Success criteria

Think carefully about the spec's intent and how each iteration should build upon previous work.

## PHASE 2: OUTPUT DIRECTORY RECONNAISSANCE
Thoroughly analyze the output_dir to understand the current state:
- Count existing iterations
- Analyze patterns and themes used
- Identify gaps or opportunities
- Determine next logical progression

## PHASE 3: ITERATION PLANNING
Based on your analysis, plan the generation approach:
- For count=1: Generate single high-quality iteration
- For count>1: Plan parallel agent coordination
- For count="infinite": Prepare wave-based generation strategy

## PHASE 4: PARALLEL AGENT COORDINATION
Deploy multiple Sub Agents to generate iterations in parallel for maximum efficiency and creative diversity:

For each Sub Agent, provide this context:
```
TASK: Generate iteration [NUMBER] for [SPEC_FILE] in [OUTPUT_DIR]
You are Sub Agent [X] generating iteration [NUMBER].
CONTEXT:
- Specification: [Full spec analysis]
- Existing iterations: [Summary of current output_dir contents]
- Your iteration number: [NUMBER]
- Assigned creative direction: [Specific innovation dimension to explore]

REQUIREMENTS:
1. Read and understand the specification completely
2. Analyze existing iterations to ensure your output is unique
3. Generate content following the spec format exactly
4. Focus on [assigned innovation dimension] while maintaining spec compliance
5. Create file with exact name pattern specified
6. Ensure your iteration adds genuine value and novelty

DELIVERABLE: Single file as specified, with unique innovative content
```

## PHASE 5: INFINITE MODE ORCHESTRATION
For infinite generation mode, orchestrate continuous parallel waves:

```
WHILE context_capacity > threshold:
  1. Assess current output_dir state
  2. Plan next wave of agents (size based on remaining context)
  3. Assign increasingly sophisticated creative directions
  4. Launch parallel Sub Agent wave
  5. Monitor wave completion
  6. Update directory state snapshot
  7. Evaluate context capacity remaining
  8. If sufficient capacity: Continue to next wave
  9. If approaching limits: Complete final wave and summarize
```

Execute this infinite generation system now.
```

## Invoking Commands

### Basic Usage
- `/command` - Runs .claude/commands/command.md
- `/namespace:command` - Runs .claude/commands/namespace/command.md
- `/command arguments` - Passes "arguments" as $ARGUMENTS

### Examples from the Repository
```bash
# Basic infinite command
/project:infinite specs/invent_new_ui_v3.md src 1

# Batch generation
/project:infinite specs/invent_new_ui_v3.md src_new 5

# Infinite mode
/project:infinite specs/invent_new_ui_v3.md infinite_src_new/ infinite
```

## Global vs Project Commands

### Project Commands (.claude/commands/)
- Available only in the specific project
- Shared with team when repository is cloned
- Project-specific prompts and workflows

### Global Commands (~/.claude/commands/)
- Available across all projects
- Personal productivity commands
- Copy useful commands here for reuse

**To make a command global:** Copy the .md file to `~/.claude/commands/`

## Advanced Patterns

### 1. Multi-Agent Coordination
Commands can orchestrate multiple Claude instances:
- Deploy "Sub Agents" for parallel processing
- Assign different creative directions to each agent
- Coordinate outputs and manage context limits

### 2. Context Management
For complex workflows:
- Progressive summarization of state
- Wave-based generation to manage context limits
- State snapshots between agent deployments

### 3. Argument Handling
```markdown
# .claude/commands/test-component.md
Generate comprehensive tests for the $ARGUMENTS component, including:
- Unit tests for all public methods
- Integration tests for component interactions
- Edge case validation
- Performance benchmarks
```

### 4. Conditional Logic
```markdown
# .claude/commands/deploy.md
Deploy the application using these steps:

If $ARGUMENTS contains "production":
- Run full test suite
- Create production build
- Deploy to production servers
- Send deployment notifications

If $ARGUMENTS contains "staging":
- Run basic tests
- Create staging build
- Deploy to staging environment

Otherwise:
- Deploy to development environment
```

## Best Practices

### 1. Command Design
- **Be specific**: Clear, detailed instructions work better
- **Include context**: Explain the purpose and expected outcomes
- **Structure well**: Use phases, steps, or numbered lists
- **Handle arguments**: Make commands flexible with $ARGUMENTS

### 2. Organization
- **Use folders**: Group related commands in subdirectories
- **Descriptive names**: Command names should be self-explanatory
- **Documentation**: Include comments explaining complex commands

### 3. Testing Commands
- Start simple and iterate
- Test with different argument patterns
- Verify commands work in different project contexts
- Get team feedback on shared commands

## Example Command Library

### Development Commands
```markdown
# .claude/commands/refactor.md
Refactor the selected code to improve:
- Readability and maintainability
- Performance where applicable
- Adherence to project conventions
- Removal of code smells

Maintain existing functionality and add tests if needed.
```

### Documentation Commands
```markdown
# .claude/commands/document.md
Create comprehensive documentation for $ARGUMENTS including:
- Purpose and functionality overview
- API reference with examples
- Usage patterns and best practices
- Configuration options
- Troubleshooting guide
```

### Analysis Commands
```markdown
# .claude/commands/security-audit.md
Perform a security audit focusing on:
- Input validation and sanitization
- Authentication and authorization
- Data exposure risks
- Dependency vulnerabilities
- Configuration security

Provide prioritized recommendations with severity levels.
```

## Settings Configuration

The `.claude/settings.json` file can configure permissions and behavior:

```json
{
  "permissions": ["Write", "MultiEdit", "Edit", "Bash"],
  "maxTokens": 4000,
  "temperature": 0.7
}
```

## Getting Started

1. **Create the directory structure**:
   ```bash
   mkdir -p .claude/commands
   ```

2. **Write your first command**:
   ```markdown
   # .claude/commands/hello.md
   Say hello and introduce yourself as a helpful coding assistant.
   Briefly explain what kinds of tasks you can help with.
   ```

3. **Test the command**:
   ```bash
   claude
   > /hello
   ```

4. **Iterate and improve** based on results

5. **Share useful commands** by copying to `~/.claude/commands/` for global access

This system transforms Claude Code into a powerful, customizable coding assistant tailored to your specific workflows and requirements.