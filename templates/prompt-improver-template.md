# Prompt Improver Template

Based on Anthropic's prompt improver recommendations, this template helps create
structured prompts that work efficiently with token-efficient mode.

## Template Structure

### 1. Task Definition

```xml
<task>
[Clear, specific description of what needs to be accomplished]
</task>
```

### 2. Context and Constraints

```xml
<context>
[Relevant background information]
[Specific constraints or requirements]
[Available resources or tools]
</context>
```

### 3. Step-by-Step Reasoning (For Complex Tasks)

```xml
<reasoning_steps>
1. [First analysis step]
2. [Second analysis step]
3. [Continue as needed]
</reasoning_steps>
```

### 4. Expected Output Format

```xml
<output_format>
[Specific structure or format required]
[Examples if helpful]
</output_format>
```

### 5. Examples (Optional but Recommended)

```xml
<examples>
<example>
Input: [Sample input]
Output: [Expected output]
</example>
</examples>
```

## Token-Efficient Prompt Patterns

### Pattern 1: Direct Task Execution

Best for: Simple, well-defined tasks

```
Task: [Specific action]
Input: [Data/parameters]
Output: [Expected result format]
```

### Pattern 2: Analysis with Structure

Best for: Code review, analysis, validation

```
Analyze [target] for:
- Criterion 1: [specific check]
- Criterion 2: [specific check]
- Criterion 3: [specific check]

Report findings as:
[Structured format]
```

### Pattern 3: Implementation Guide

Best for: Code generation, feature implementation

```
Implement [feature] with:
Requirements:
- [Requirement 1]
- [Requirement 2]

Constraints:
- [Constraint 1]
- [Constraint 2]

Use existing patterns from: [reference]
```

## Tips for Token Efficiency

1. **Be Specific**: Vague prompts lead to verbose responses
2. **Use Structure**: XML tags or clear sections guide concise output
3. **Provide Examples**: One good example prevents lengthy explanations
4. **Set Boundaries**: Specify what NOT to include
5. **Request Format**: Define exact output structure

## Agent-Specific Templates

### For sprint-planner

```yaml
Task: Plan sprint for [feature]
Output: QDIRECTOR YAML format
Focus: Task dependencies and exit criteria
Constraints: 2-4 hour task chunks
```

### For code-implementer

```
Implement: [specific functionality]
Following: [project patterns from X]
Constraints: No comments, match style
Output: Production-ready code
```

### For validation-expert

```
Validate: [code/component]
Criteria: Security, performance, quality
Severity levels: CRITICAL, HIGH, MEDIUM, LOW
Output: Structured findings with locations
```

### For test-specialist

```
Generate tests for: [component/function]
Coverage: Happy path, edge cases, errors
Framework: [detected from codebase]
Pattern: AAA (Arrange, Act, Assert)
```

## Prompt Improvement Checklist

Before using a prompt, verify:

- [ ] Task is clearly defined
- [ ] Output format is specified
- [ ] Constraints are explicit
- [ ] Examples provided for complex tasks
- [ ] Unnecessary verbosity is discouraged
- [ ] Structure guides the response

## Integration with Token-Efficient Mode

When token-efficient mode is enabled:

1. These templates naturally produce more concise outputs
2. The beta header optimizes response generation
3. Structured prompts maximize token savings
4. Clear boundaries prevent over-generation

Remember: The goal is clarity and precision, which naturally leads to
efficiency.
