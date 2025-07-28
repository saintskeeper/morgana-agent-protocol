# QVALIDATE Framework - Unified Validation System

Orchestrates all validation commands (QCHECK, QCHECKF, QCHECKT) for
comprehensive quality assurance.

## Purpose

Provide a unified validation interface that QDIRECTOR can use to validate task
completion across code, functions, and tests.

## Unified Validation Architecture

### 1. Validation Pipeline

```yaml
validation_pipeline:
  stages:
    - name: "syntax_check"
      tools: ["lint", "typecheck"]
      blocking: true

    - name: "function_analysis"
      command: "/qcheckf"
      blocking: true

    - name: "test_validation"
      command: "/qcheckt"
      blocking: true

    - name: "integration_check"
      command: "/qcheck"
      blocking: true

    - name: "security_scan"
      optional: true
      condition: "if security-sensitive"
```

### 2. Aggregate Scoring System

```yaml
validation_score:
  components:
    code_quality:
      weight: 0.3
      sources: ["/qcheck"]

    function_quality:
      weight: 0.25
      sources: ["/qcheckf"]

    test_quality:
      weight: 0.25
      sources: ["/qcheckt"]

    coverage:
      weight: 0.2
      sources: ["coverage reports"]

  final_score: weighted_average(components)
  pass_threshold: 85%
```

### 3. Unified Output Format

```yaml
unified_validation_report:
  task_id: "AUTH_IMPL"
  timestamp: "2024-01-15 10:45:00"
  overall_status: "NEEDS_RETRY"
  overall_score: 78%

  summary:
    files_changed: 12
    functions_added: 8
    tests_added: 15
    coverage_delta: +5%

  blocking_issues:
    - source: "qcheck"
      issue: "SQL injection vulnerability"
      severity: "CRITICAL"
      location: "AuthService.ts:45"

    - source: "qcheckt"
      issue: "Missing error case tests"
      severity: "HIGH"
      impact: "Core auth flow uncovered"

  non_blocking_issues:
    - source: "qcheckf"
      issue: "High complexity"
      severity: "MEDIUM"
      suggestion: "Refactor processLogin()"

  recommendations:
    retry_strategy: "focused"
    focus_areas:
      - "Fix SQL injection in AuthService"
      - "Add timeout tests for login"
    suggested_model: "o3-mini"
    estimated_effort: "30 minutes"
```

### 4. Validation Rules Engine

```typescript
interface ValidationRule {
  id: string;
  name: string;
  checker: (context: ValidationContext) => ValidationResult;
  severity: "CRITICAL" | "HIGH" | "MEDIUM" | "LOW";
  autoFixable: boolean;
}

// Example custom rules
const projectRules: ValidationRule[] = [
  {
    id: "no-console-logs",
    name: "No console.log in production code",
    checker: (ctx) => {
      const hasConsoleLogs = ctx.changes.some((change) =>
        change.content.includes("console.log"),
      );
      return {
        passed: !hasConsoleLogs,
        message: "Remove console.log statements",
        locations: findConsoleLogLocations(ctx),
      };
    },
    severity: "HIGH",
    autoFixable: true,
  },

  {
    id: "api-versioning",
    name: "API endpoints must include version",
    checker: (ctx) => {
      const endpoints = findApiEndpoints(ctx);
      const unversioned = endpoints.filter(
        (e) => !e.path.includes("/v1/") && !e.path.includes("/v2/"),
      );
      return {
        passed: unversioned.length === 0,
        message: "Add version to API endpoints",
        locations: unversioned,
      };
    },
    severity: "HIGH",
    autoFixable: false,
  },
];
```

### 5. Progressive Validation Strategy

```yaml
validation_modes:
  quick:
    description: "Fast validation for rapid iteration"
    includes: ["syntax", "basic tests"]
    duration: "~30 seconds"
    use_when: "development iteration"

  standard:
    description: "Comprehensive validation"
    includes: ["all checks", "coverage", "integration"]
    duration: "~2 minutes"
    use_when: "pre-commit, task completion"

  deep:
    description: "Full analysis with performance testing"
    includes: ["standard", "performance", "security scan"]
    duration: "~5 minutes"
    use_when: "pre-deployment, critical paths"
```

### 6. Smart Retry Recommendations

```yaml
retry_intelligence:
  failure_patterns:
    - pattern: "missing_tests"
      recommendation:
        model: "o3-mini"
        prompt_enhancement: "Focus on edge cases and error scenarios"
        context: ["existing test patterns", "coverage gaps"]

    - pattern: "security_vulnerability"
      recommendation:
        model: "o3"
        prompt_enhancement: "Apply OWASP best practices"
        context: ["security guidelines", "similar fixes"]

    - pattern: "high_complexity"
      recommendation:
        model: "gemini-2.5-pro"
        prompt_enhancement: "Refactor using SOLID principles"
        context: ["design patterns", "clean code examples"]
```

### 7. Validation Orchestration

```typescript
class ValidationOrchestrator {
  async validateTask(
    taskId: string,
    mode: ValidationMode,
  ): Promise<ValidationReport> {
    const pipeline = this.getPipeline(mode);
    const results: StageResult[] = [];

    for (const stage of pipeline) {
      const result = await this.runStage(stage);
      results.push(result);

      // Early exit on critical failures
      if (result.blocking && !result.passed) {
        break;
      }
    }

    return this.aggregateResults(results);
  }

  private async runStage(stage: ValidationStage): Promise<StageResult> {
    switch (stage.type) {
      case "command":
        return this.runCommand(stage.command);
      case "tool":
        return this.runTool(stage.tool);
      case "custom":
        return this.runCustomValidation(stage.rules);
    }
  }

  private aggregateResults(results: StageResult[]): ValidationReport {
    // Combine all validation results
    // Calculate overall score
    // Generate recommendations
    // Format for QDIRECTOR consumption
  }
}
```

### 8. Continuous Improvement

```yaml
validation_metrics:
  tracking:
    - first_pass_success_rate
    - average_retry_count
    - common_failure_patterns
    - time_to_validation_pass

  learning:
    - identify_recurring_issues
    - update_validation_rules
    - adjust_severity_levels
    - optimize_retry_strategies

  reporting:
    weekly_summary:
      total_validations: 156
      pass_rate: 72%
      average_score: 83%

      top_issues:
        - "Missing test coverage": 23%
        - "High complexity": 18%
        - "Security vulnerabilities": 8%

      improvement_actions:
        - "Add pre-commit hooks for common issues"
        - "Create code generation templates"
        - "Update developer guidelines"
```

### 9. Integration with QDIRECTOR

```yaml
qdirector_integration:
  validation_triggers:
    - after_code_generation
    - before_task_completion
    - on_retry_attempt
    - pre_deployment

  decision_logic:
    if score >= 90 and no_critical_issues:
      action: "mark_complete"

    elif score >= 70 and retry_count < 3:
      action: "retry_with_focus"

    elif critical_security_issue:
      action: "escalate_to_human"

    else:
      action: "detailed_human_review"

  context_enhancement:
    on_retry:
      - include_validation_report
      - add_specific_examples
      - reference_best_practices
      - focus_on_critical_issues
```

### 10. Usage Examples

```bash
# Run full validation pipeline
/qvalidate --mode standard

# Quick validation during development
/qvalidate --mode quick --file AuthService.ts

# Deep validation before deployment
/qvalidate --mode deep --include-security

# Validate specific task for QDIRECTOR
/qvalidate --task-id AUTH_IMPL --format qdirector

# Generate validation report
/qvalidate --report --output validation-report.yaml
```

## Benefits

1. **Consistency**: Same validation standards across all code
2. **Efficiency**: Parallel validation where possible
3. **Intelligence**: Smart retry recommendations
4. **Learning**: Continuous improvement from patterns
5. **Integration**: Seamless QDIRECTOR compatibility

This unified framework ensures comprehensive quality validation while
maintaining fast feedback loops for development.
