# Enhanced Commands Quick Reference

## Workflow Overview

```
1. Plan Sprint    → /qnew-enhanced + /qplan-enhanced
2. Execute Tasks  → /qdirector-enhanced (orchestrates all below)
3. Validate Code  → /qcheck-enhanced, /qcheckf-enhanced, /qcheckt-enhanced
4. Unified Check  → /qvalidate-framework
5. Commit Changes → /qgit (with pre-commit validation)
```

## Command Summary

### Planning Commands

- **`/qnew-enhanced`** - Generate sprint plans with tasks, dependencies, exit
  criteria
- **`/qplan-enhanced`** - Validate technical feasibility, enrich with codebase
  context

### Execution (via QDIRECTOR)

- **`/qdirector-enhanced`** - Master orchestrator with retry logic and
  validation
- **`/qcode`** - Implementation (auto-validates with qcheckf-enhanced)
- **`/qtest`** - Test generation (auto-validates with qcheckt-enhanced)

### Validation Commands

- **`/qcheck-enhanced`** - Comprehensive code validation (security, performance,
  quality)
- **`/qcheckf-enhanced`** - Function-level analysis (complexity, design,
  efficiency)
- **`/qcheckt-enhanced`** - Test quality validation (coverage, patterns,
  effectiveness)
- **`/qvalidate-framework`** - Unified validation orchestration

### Completion

- **`/qgit`** - Semantic commits (runs validation pipeline pre-commit)

## Key Features

### Structured Output

All enhanced commands output YAML for QDIRECTOR parsing:

```yaml
validation_report:
  score: 85%
  must_fix: [critical issues]
  should_fix: [improvements]
  ready_for_merge: true|false
  retry_recommendation:
    model: "o3-mini"
    focus_areas: ["specific fixes"]
```

### Automatic Retry

- Up to 3 attempts before human escalation
- Smart model selection based on issue type
- Focused context on specific problems

### Validation Pipeline

```
Syntax → Functions → Tests → Integration → Security
```

### Model Selection

- Planning: `gemini-2.5-pro`, `o3`
- Implementation: `gpt-4.1`, `gemini-2.5-flash`
- Testing: `o3-mini`, `gemini-2.5-flash`
- Quick fixes: `gemini-2.5-flash`, `o3-mini`

## Example Workflow

```bash
# 1. Create sprint plan
/qnew-enhanced Build secure user authentication with JWT

# 2. Validate and enrich plan
/qplan-enhanced --sprint sprint-2024-01-auth.md

# 3. Execute with director
/qdirector-enhanced Execute sprint plan in sprint-2024-01-auth.md

# Director automatically:
# - Spawns agents with enhanced commands
# - Validates outputs
# - Retries with focused context
# - Escalates when needed
```

## Quick Tips

1. **Always use enhanced versions** for consistency
2. **Check validation scores** before marking tasks complete
3. **Use retry recommendations** for efficient fixes
4. **Monitor validation metrics** for process improvement
5. **Keep sprint plans updated** with actual progress

## Validation Thresholds

- **90%+** → Auto-approve
- **70-89%** → Retry if < 3 attempts
- **<70%** → Human review recommended
- **Critical issues** → Always block completion

## Files Created

- Sprint plans: `sprint-[date]-[feature].md`
- Validation logs: `.qdirector/logs/`
- Metrics: `.qdirector/metrics.yaml`
- Status: `.qdirector/status.md`
