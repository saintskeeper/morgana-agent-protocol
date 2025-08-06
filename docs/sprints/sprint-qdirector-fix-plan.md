# Technical Implementation Plan: Fix QDIRECTOR Agent Type Issue

## Sprint Overview

**Goal**: Fix the Task tool limitation that only recognizes "general-purpose"
agent type **Approach**: Implement Agent Adapter pattern to maintain
functionality **Priority**: P0-Critical (System is non-functional without this
fix) **Estimated Effort**: 4-8 hours

## Task Breakdown

### TASK_001: Create Agent Adapter Infrastructure

**Title**: Implement AgentAdapter wrapper function **Priority**: P0-Critical
**Dependencies**: None **Complexity**: Medium **Exit Criteria**:

- [ ] AgentAdapter function created that wraps Task() calls
- [ ] Agent prompts loaded from .claude/agents/ files
- [ ] Error handling for missing agent types
- [ ] Maintains parallel execution capability
- [ ] Preserves all Task() parameters

**Technical Approach**:

```python
# Add to beginning of qdirector-enhanced.md
def AgentAdapter(agent_type, prompt, **kwargs):
    """
    Adapter to bridge custom agent types with general-purpose Task tool
    """
    agent_prompts = {
        "code-implementer": load_agent_prompt("code-implementer.md"),
        "sprint-planner": load_agent_prompt("sprint-planner.md"),
        "test-specialist": load_agent_prompt("test-specialist.md"),
        "validation-expert": load_agent_prompt("validation-expert.md")
    }

    if agent_type not in agent_prompts:
        raise ValueError(f"Unknown agent type: {agent_type}")

    # Combine agent system prompt with task prompt
    full_prompt = f"{agent_prompts[agent_type]}\n\nTask: {prompt}"

    # Call Task with general-purpose
    return Task(subagent_type="general-purpose", prompt=full_prompt, **kwargs)
```

**Implementation Notes**:

- Follow Python function patterns in lines 407-416
- Use structured format similar to agent_request (lines 530-555)
- Include error handling patterns from validation sections

### TASK_002: Implement Agent Prompt Loader

**Title**: Create function to load agent prompts from files **Priority**:
P0-Critical **Dependencies**: TASK_001 **Complexity**: Simple **Exit Criteria**:

- [ ] Function reads agent markdown files
- [ ] Extracts content after YAML frontmatter
- [ ] Caches loaded prompts for performance
- [ ] Handles file not found gracefully

**Technical Approach**:

```python
def load_agent_prompt(agent_file):
    """Load agent system prompt from markdown file"""
    # Cache for performance
    if hasattr(load_agent_prompt, 'cache'):
        if agent_file in load_agent_prompt.cache:
            return load_agent_prompt.cache[agent_file]
    else:
        load_agent_prompt.cache = {}

    path = f"/Users/walterday/.claude/agents/{agent_file}"
    try:
        with open(path, 'r') as f:
            content = f.read()
            # Skip YAML frontmatter
            parts = content.split('---', 2)
            if len(parts) >= 3:
                prompt = parts[2].strip()
                load_agent_prompt.cache[agent_file] = prompt
                return prompt
    except FileNotFoundError:
        return f"Act as a {agent_file.replace('.md', '')} specialist"
```

### TASK_003: Refactor All Task Calls

**Title**: Replace Task() with AgentAdapter() throughout qdirector-enhanced.md
**Priority**: P0-Critical **Dependencies**: TASK_001, TASK_002 **Complexity**:
Simple (but repetitive) **Exit Criteria**:

- [ ] All 15+ Task() calls updated to use AgentAdapter()
- [ ] Parallel execution arrays maintained
- [ ] All agent_type values preserved
- [ ] Command still reads clearly

**Refactoring Pattern**:

```python
# Before:
Task(subagent_type="code-implementer", prompt="implement auth service")

# After:
AgentAdapter("code-implementer", "implement auth service")

# Parallel execution before:
parallel_tasks = [
    Task(subagent_type="code-implementer", prompt="implement auth service"),
    Task(subagent_type="test-specialist", prompt="create auth test suite")
]

# Parallel execution after:
parallel_tasks = [
    AgentAdapter("code-implementer", "implement auth service"),
    AgentAdapter("test-specialist", "create auth test suite")
]
```

**Files to Update**:

- Lines 164, 166-167, 169 (execution pattern examples)
- Lines 174-176 (parallel tasks example)
- Lines 275, 281-283, 289-291, 299-301, 305, 308, 314 (workflow examples)

### TASK_004: Add Logging and Debugging

**Title**: Implement logging for agent adapter usage **Priority**: P1-High
**Dependencies**: TASK_003 **Complexity**: Simple **Exit Criteria**:

- [ ] Log which agent type is requested
- [ ] Log successful adapter translations
- [ ] Log any fallback behaviors
- [ ] Include in monitoring section (lines 425-458)

### TASK_005: Create Fallback Mechanism

**Title**: Add direct agent invocation fallback **Priority**: P2-Medium
**Dependencies**: TASK_003 **Complexity**: Medium **Exit Criteria**:

- [ ] Detect if Task() call fails
- [ ] Fallback to direct agent request
- [ ] Log fallback usage
- [ ] Document in error handling section

## Validation Steps

1. **Unit Testing**:

   - Test AgentAdapter with each agent type
   - Verify prompt combination
   - Test error cases (missing agent)

2. **Integration Testing**:

   - Run simple workflow: `/qdirector-enhanced implement hello world function`
   - Verify parallel execution works
   - Check agent outputs are formatted correctly

3. **Regression Testing**:
   - Ensure color formatting still works (lines 390-416)
   - Verify TodoWrite integration maintained
   - Check validation pipeline triggers

## Risk Assessment

### Technical Risks

1. **Risk**: Agent prompts might be too long for context

   - **Mitigation**: Implement prompt truncation if needed
   - **Severity**: Medium

2. **Risk**: Task tool might reject complex prompts

   - **Mitigation**: Test with actual Task tool early
   - **Severity**: High

3. **Risk**: Performance impact from file I/O
   - **Mitigation**: Cache loaded prompts
   - **Severity**: Low

### Implementation Risks

1. **Risk**: Missing some Task() calls during refactoring
   - **Mitigation**: Use grep to find all instances
   - **Severity**: Medium

## Dependencies

- No external dependencies
- Uses existing Task tool
- Reads existing agent files

## Success Metrics

- [ ] All Task() calls successfully execute
- [ ] Parallel execution maintains performance
- [ ] Agent behaviors match expected specialization
- [ ] No regression in existing functionality

## Next Steps After Implementation

1. Test with real workflows
2. Monitor performance impact
3. Document workaround for team
4. Consider contributing fix to Task tool upstream
