# Agent Model Optimization Strategy

## Recommended Model Defaults with Token-Efficient Mode

### Agent-Specific Recommendations

#### sprint-planner

- **Keep Current**: gemini-2.5-pro or o3
- **Rationale**: Complex reasoning, dependency analysis, and strategic planning
  require deeper models
- **Prompt Optimization**: Use structured templates but maintain model choice

#### code-implementer

- **Current**: gpt-4.1 or gemini-2.5-flash
- **Recommendation**: Add claude-3-7-sonnet as option for simple implementations
- **Decision Tree**:
  ```
  IF task_complexity == "simple" AND token_efficient_enabled:
    -> claude-3-7-sonnet with structured prompt
  ELIF task_complexity == "moderate":
    -> gemini-2.5-flash
  ELSE:
    -> gpt-4.1
  ```

#### validation-expert

- **Current**: Already uses claude-3-7-sonnet
- **Recommendation**: Make claude-3-7-sonnet PRIMARY with token-efficient
  prompts
- **Structured Prompt Pattern**:
  ```
  Validate: [component]
  Criteria: Security, Performance, Quality
  Output: YAML findings with severity/location
  ```

#### test-specialist

- **Current**: Already uses claude-3-7-sonnet as primary
- **Recommendation**: Optimize with token-efficient prompts
- **Structured Prompt Pattern**:
  ```
  Generate tests for: [component]
  Coverage: Happy/Edge/Error
  Framework: [detected]
  Pattern: AAA
  ```

## Implementation Strategy

### 1. Agent Prompt Templates

Each agent should receive optimized prompts when token-efficient mode is
enabled:

```python
def get_agent_prompt(agent_type, task, token_efficient_enabled):
    if token_efficient_enabled:
        return structured_prompts[agent_type].format(task=task)
    else:
        return standard_prompts[agent_type].format(task=task)
```

### 2. Model Selection Logic

```yaml
model_selection:
  sprint-planner:
    default: gemini-2.5-pro
    token_efficient: gemini-2.5-pro # No change - needs reasoning

  code-implementer:
    simple_task:
      default: gemini-2.5-flash
      token_efficient: claude-3-7-sonnet
    complex_task:
      default: gpt-4.1
      token_efficient: gpt-4.1 # No change - needs capability

  validation-expert:
    default: claude-3-7-sonnet
    token_efficient: claude-3-7-sonnet # Already optimal

  test-specialist:
    default: claude-3-7-sonnet
    token_efficient: claude-3-7-sonnet # Already optimal
```

### 3. Structured Prompt Templates per Agent

#### Sprint Planner Template

```xml
<task>Plan sprint for: {feature}</task>
<output>QDIRECTOR YAML</output>
<constraints>
- 2-4 hour tasks
- Clear dependencies
- Exit criteria required
</constraints>
```

#### Code Implementer Template

```xml
<task>Implement: {feature}</task>
<context>
- Patterns: {project_patterns}
- Style: {code_style}
</context>
<constraints>
- No comments
- Production ready
- Match conventions
</constraints>
```

#### Validation Expert Template

```xml
<task>Validate: {component}</task>
<criteria>
- Security: injections, auth
- Performance: O(n) complexity
- Quality: SOLID, DRY
</criteria>
<output>
severity: issue at location
</output>
```

#### Test Specialist Template

```xml
<task>Test: {component}</task>
<coverage>
- Happy paths
- Edge cases
- Error scenarios
</coverage>
<pattern>AAA</pattern>
<framework>{detected}</framework>
```

## Expected Benefits

### With Token-Efficient Mode + Structured Prompts:

1. **Cost Reduction**: 40-60% average across all agents
2. **Speed Improvement**: 30-50% faster responses
3. **Quality Maintenance**: Same accuracy with clearer outputs
4. **Consistency**: More predictable agent responses

### Measurement Strategy:

```yaml
metrics:
  token_usage:
    baseline: [current average per agent]
    optimized: [with token-efficient mode]

  response_time:
    baseline: [current average]
    optimized: [with structured prompts]

  quality_score:
    validation_pass_rate: [maintain >90%]
    retry_rate: [target <20%]
```

## Rollout Plan

### Phase 1: Test with Non-Critical Agents

- validation-expert (already uses claude-3-7-sonnet)
- test-specialist (already uses claude-3-7-sonnet)

### Phase 2: Expand to Simple Tasks

- code-implementer for simple tasks only
- documentation generation

### Phase 3: Full Integration

- Update QDIRECTOR routing logic
- Add complexity detection
- Implement automatic prompt structuring

## Configuration Updates Needed

1. Update QDIRECTOR model selection logic
2. Add structured prompt templates to each agent
3. Implement complexity detection for code-implementer
4. Add metrics tracking for optimization validation
