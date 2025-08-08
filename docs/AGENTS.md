# Agent Selection Guide

## Quick Decision Tree

```
What do you need to do?
│
├─ Plan work? → sprint-planner
├─ Write code? → code-implementer
├─ Create tests? → test-specialist
└─ Review/validate? → validation-expert
```

## Agent Comparison Matrix

| Aspect               | code-implementer         | sprint-planner                  | test-specialist            | validation-expert           |
| -------------------- | ------------------------ | ------------------------------- | -------------------------- | --------------------------- |
| **Primary Focus**    | Implementation           | Planning                        | Testing                    | Quality Assurance           |
| **Default Model**    | Claude 3.7 Sonnet        | Gemini 2.5 Pro                  | O3-mini                    | Claude 3.7 Sonnet           |
| **Token Efficient**  | ✅ Yes (14-70% savings)  | ❌ No                           | ✅ Yes                     | ✅ Yes                      |
| **Typical Duration** | 2-5 minutes              | 1-3 minutes                     | 2-4 minutes                | 1-2 minutes                 |
| **Retry Escalation** | Claude 4 → Opus          | O3 → Gemini Pro                 | Gemini Flash → O3          | Claude 4 → Opus             |
| **Best Use Cases**   | Features, APIs, Services | Sprint planning, Task breakdown | Unit/Integration/E2E tests | Code review, Security audit |

## When to Use Each Agent

### 🛠️ code-implementer

**Use when you need to:**

- Implement new features or functionality
- Create API endpoints or services
- Refactor existing code
- Fix bugs or issues
- Build UI components

**Strengths:**

- Follows project conventions automatically
- Implements security best practices
- Optimizes for performance
- Produces production-ready code

**Example:**

```bash
morgana -- --agent code-implementer --prompt "Create user authentication service with JWT"
```

### 📋 sprint-planner

**Use when you need to:**

- Break down large projects into tasks
- Create development roadmaps
- Prioritize work items
- Define task dependencies
- Estimate implementation effort

**Strengths:**

- Creates structured sprint plans
- Identifies dependencies and blockers
- Provides clear exit criteria
- Generates task IDs for tracking

**Example:**

```bash
morgana -- --agent sprint-planner --prompt "Plan 2-week sprint for payment processing system"
```

### 🧪 test-specialist

**Use when you need to:**

- Generate comprehensive test suites
- Create edge case scenarios
- Build integration tests
- Design E2E test flows
- Improve test coverage

**Strengths:**

- Identifies edge cases automatically
- Creates realistic test data
- Follows testing best practices
- Achieves high coverage targets

**Example:**

```bash
morgana -- --agent test-specialist --prompt "Create unit tests for UserService with edge cases"
```

### ✅ validation-expert

**Use when you need to:**

- Review code quality
- Perform security audits
- Check compliance standards
- Validate architecture decisions
- Assess performance implications

**Strengths:**

- Identifies security vulnerabilities
- Catches performance issues
- Ensures code maintainability
- Validates against best practices

**Example:**

```bash
morgana -- --agent validation-expert --prompt "Review authentication module for security issues"
```

## Parallel Agent Execution

Run multiple agents simultaneously for complex tasks:

```bash
echo '[
  {
    "agent_type": "sprint-planner",
    "prompt": "Plan user management features"
  },
  {
    "agent_type": "validation-expert",
    "prompt": "Audit existing auth code"
  },
  {
    "agent_type": "test-specialist",
    "prompt": "Design test strategy for auth"
  }
]' | morgana --parallel
```

## Model Escalation Strategy

Each agent uses intelligent model escalation:

1. **Initial Attempt**: Default model (optimized for speed/cost)
2. **Retry 1**: Enhanced model (better reasoning)
3. **Retry 2+**: Maximum capability model
4. **Validation Failure**: Specialized recovery model

### Token Optimization

Agents marked with token-efficient mode can save 14-70% on API costs:

```yaml
# Enable globally
~/.claude/scripts/token-efficient-config.sh enable

# Agents that support token-efficient mode:
- code-implementer
- test-specialist
- validation-expert
```

## Custom Agent Creation

Create your own specialized agents:

1. Create agent definition in `~/.claude/agents/my-agent.md`:

```markdown
---
name: my-custom-agent
description: Specialized agent for specific tasks
tools: Read, Write, Edit, Bash
model_selection:
  default: gemini-2.5-flash
  escalation:
    retry_1: o3-mini
    retry_2: gemini-2.5-pro
---

You are a specialized agent for [specific purpose]...
```

2. Use your custom agent:

```bash
morgana -- --agent my-custom-agent --prompt "Perform specialized task"
```

## Performance Benchmarks

| Agent             | Avg Execution Time | Token Usage | Success Rate |
| ----------------- | ------------------ | ----------- | ------------ |
| code-implementer  | 3.2 min            | 2.8k tokens | 94%          |
| sprint-planner    | 1.8 min            | 1.2k tokens | 98%          |
| test-specialist   | 2.5 min            | 2.1k tokens | 96%          |
| validation-expert | 1.3 min            | 950 tokens  | 97%          |

_With token-efficient mode enabled, token usage reduced by average 42%_

## Best Practices

1. **Start with planning**: Use sprint-planner before implementation
2. **Parallel where possible**: Run independent agents simultaneously
3. **Validate early**: Use validation-expert during development, not just at the
   end
4. **Test comprehensively**: test-specialist for critical functionality
5. **Chain agents**: Plan → Implement → Test → Validate workflow

## Troubleshooting

### Agent Timeout

```bash
# Increase timeout for complex tasks
morgana -- --agent code-implementer --prompt "Complex task" --timeout 10m
```

### Model Not Available

```bash
# Check available models
morgana --list-models

# Use fallback model
morgana -- --agent code-implementer --prompt "Task" --model claude-3-7-sonnet
```

### Token Limit Exceeded

```bash
# Enable token-efficient mode
~/.claude/scripts/token-efficient-config.sh enable

# Split into smaller tasks
morgana -- --agent sprint-planner --prompt "Break down: [large task]"
```
