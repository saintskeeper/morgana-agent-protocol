# Enhanced Commands Quick Reference

## Workflow Overview

```
1. Plan Sprint    → /morgana-plan + /morgana-validate
2. Execute Tasks  → /morgana-director (orchestrates all below)
3. Validate Code  → /morgana-check, /morgana-check-function, /morgana-check-tests
4. Unified Check  → /morgana-validate-all
5. Commit Changes → /morgana-commit (with pre-commit validation)
```

## Command Summary

### Planning Commands

- **`/morgana-plan`** - Generate sprint plans with tasks, dependencies, exit
  criteria
- **`/morgana-validate`** - Validate technical feasibility, enrich with codebase
  context

### Execution (via MORGANA-DIRECTOR)

- **`/morgana-director`** - Master orchestrator with retry logic and validation
- **`/morgana-code`** - Implementation (auto-validates with
  morgana-check-function)
- **`/morgana-test`** - Test generation (auto-validates with
  morgana-check-tests)

### Validation Commands

- **`/morgana-check`** - Comprehensive code validation (security, performance,
  quality)
- **`/morgana-check-function`** - Function-level analysis (complexity, design,
  efficiency)
- **`/morgana-check-tests`** - Test quality validation (coverage, patterns,
  effectiveness)
- **`/morgana-validate-all`** - Unified validation orchestration

### Completion

- **`/morgana-commit`** - Semantic commits (runs validation pipeline pre-commit)

## Key Features

### Structured Output

All enhanced commands output YAML for MORGANA-DIRECTOR parsing:

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
/morgana-plan Build secure user authentication with JWT

# 2. Validate and enrich plan
/morgana-validate --sprint sprint-2024-01-auth.md

# 3. Execute with director
/morgana-director Execute sprint plan in sprint-2024-01-auth.md

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
