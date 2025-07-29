# QDIRECTOR Sub-Agents

This directory contains specialized sub-agents for the QDIRECTOR enhanced
orchestration system. Each agent is designed with specific expertise and tool
access to handle different aspects of the software development lifecycle.

## Available Agents

### 🎯 sprint-planner

**Purpose**: Decomposes requirements into structured sprint plans **Expertise**:
Requirements analysis, task breakdown, dependency mapping **Tools**: Read,
Write, TodoWrite, Grep, Glob **Best For**: Initial planning, sprint setup, task
prioritization

### 💻 code-implementer

**Purpose**: Implements clean, secure, performant code **Expertise**: Code
patterns, security practices, performance optimization **Tools**: Read, Write,
Edit, MultiEdit, Bash, Grep, Glob, LS **Best For**: Feature implementation, bug
fixes, refactoring

### 🔍 validation-expert

**Purpose**: Comprehensive code quality and security validation **Expertise**:
Security audits, code review, performance analysis **Tools**: Read, Grep, Glob,
Bash, mcp**zen**codereview, mcp**zen**secaudit, mcp**zen**analyze **Best For**:
Pre-commit validation, security reviews, quality gates

### 🧪 test-specialist

**Purpose**: Creates comprehensive test suites with edge case coverage
**Expertise**: Test patterns, coverage strategies, edge case identification
**Tools**: Read, Write, Edit, MultiEdit, Bash, Grep, Glob, mcp**zen**testgen
**Best For**: Test creation, coverage improvement, regression prevention

## Usage with QDIRECTOR

These agents are designed to work seamlessly with `/qdirector-enhanced` command:

```bash
# Parallel execution example
/qdirector-enhanced
- Task(subagent_type="sprint-planner", prompt="Plan authentication feature")
- Task(subagent_type="validation-expert", prompt="Audit existing auth code")
```

## Agent Communication Protocol

Agents communicate using structured YAML formats:

```yaml
request:
  task_id: "AUTH_001"
  context: "minimal_relevant"
  constraints: ["security_critical"]

response:
  status: "completed"
  outputs: ["sprint-plan.yaml"]
  next_agents: ["code-implementer"]
```

## Creating Custom Agents

To create a new agent:

1. Create a new `.md` file in this directory
2. Add front matter with name, description, and tools
3. Write detailed instructions for the agent's behavior
4. Ensure output formats are QDIRECTOR-compatible

Example structure:

```markdown
---
name: your-agent-name
description: Brief description of agent's purpose
tools: Tool1, Tool2, Tool3
---

Agent instructions here...
```

## Best Practices

1. **Single Responsibility**: Each agent should excel at one thing
2. **Clear Output**: Always produce structured, parseable output
3. **Tool Restrictions**: Only grant necessary tools
4. **Context Efficiency**: Request minimal context to preserve tokens
5. **Error Handling**: Include clear error states and messages

## Integration with Other Commands

These agents enhance other q-commands:

- `/qnew-enhanced` → Uses sprint-planner
- `/qcode` → Uses code-implementer
- `/qtest` → Uses test-specialist
- `/qvalidate-framework` → Uses validation-expert

## Performance Tips

- Run multiple agents in parallel when tasks are independent
- Use appropriate models per agent (configured in QDIRECTOR)
- Cache agent outputs for reuse across retry attempts
- Monitor token usage and optimize context passing

## Token-Efficient Mode

When token-efficient mode is enabled
(`~/.claude/scripts/token-efficient-config.sh enable`):

- Agents automatically benefit from 14-70% token reduction
- Works with Claude 3.7 Sonnet models
- No changes needed to agent configurations
- QDIRECTOR handles model routing automatically

Ideal for:

- High-volume agent workflows
- Cost-sensitive operations
- Faster response times needed
