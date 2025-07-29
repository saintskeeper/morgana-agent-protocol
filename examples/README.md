## 📘 Real-World Examples

### Example 1: Building a REST API from Scratch

```bash
# 1. Create sprint plan for the API
/qnew-enhanced Build RESTful API for user management with CRUD operations, JWT auth, and PostgreSQL

# 2. QDIRECTOR generates tasks like:
# - API_DESIGN: OpenAPI specification
# - DB_SCHEMA: User table design
# - AUTH_IMPL: JWT middleware
# - CRUD_IMPL: User endpoints
# - TEST_IMPL: Integration tests

# 3. Execute the sprint
/qdirector-enhanced sprint-2025-01-api.md

# 4. Monitor progress
# QDIRECTOR shows:
# ✅ API_DESIGN (completed)
# 🔄 AUTH_IMPL (in progress - implementing JWT)
# ⏳ CRUD_IMPL (pending - waiting for auth)
```

### Example 2: Debugging Production Issue

```bash
# 1. Describe the issue
/qdirector-enhanced Debug memory leak in payment processing service causing OOM after 48 hours

# 2. QDIRECTOR automatically:
# - Spawns validation-expert to analyze memory patterns
# - Uses code-implementer to trace allocations
# - Generates memory profiling tests
# - Identifies leak in connection pooling

# 3. Fix implementation
/qcode Fix connection pool leak by implementing proper cleanup in payment.service.ts

# 4. Validate fix
/qvalidate-framework --mode deep --focus performance
```

### Example 3: Refactoring Legacy Code

```bash
# 1. Analyze existing code
/qcheckf-enhanced --file src/legacy/OrderProcessor.js
# Output: Complexity: 47, Lines: 892, Recommendation: urgent refactor

# 2. Plan refactoring sprint
/qnew-enhanced Refactor OrderProcessor into microservices with clean architecture

# 3. Execute with parallel agents
/qdirector-enhanced
# Parallel execution:
# - Agent 1: Extract payment logic
# - Agent 2: Extract inventory logic
# - Agent 3: Create integration tests

# 4. Validate improvements
/qcheckf-enhanced --dir src/services/
# Output: Average complexity: 8, All functions < 50 lines
```

### Example 4: Security Audit Workflow

```bash
# 1. Run comprehensive security audit
/qcheck-enhanced --focus security --path ./src/

# 2. Critical issue found: SQL injection in UserService
validation_report:
  must_fix:
    - issue: "SQL injection vulnerability"
      location: "UserService.ts:45"
      severity: "critical"

# 3. Fix with focused agent
/qdirector-enhanced Fix SQL injection in UserService using parameterized queries

# 4. Verify fix
/qvalidate-framework --mode deep --security-only
# Output: All security checks passed
```

### Example 5: Test Coverage Improvement

```bash
# 1. Check current coverage
/qcheckt-enhanced --dir src/
# Output: Overall coverage: 62%, Critical paths: 45%

# 2. Generate comprehensive tests
/qtest generate comprehensive --focus critical-paths

# 3. Validate test quality
/qcheckt-enhanced --dir src/
# Output: Overall coverage: 91%, Critical paths: 95%

# 4. Commit with confidence
/qgit feat: improve test coverage to 91% with critical path focus
```

### Example 6: Feature Development with Token Efficiency

```bash
# 1. Enable token-efficient mode
~/.claude/scripts/token-efficient-config.sh enable

# 2. Create simple utility (uses Claude 3.7 Sonnet automatically)
/qcode Create date formatting utility with timezone support
# Token usage: 2,847 (vs 6,712 without efficiency)

# 3. Complex feature (automatically uses Claude 4 Opus)
/qcode Design distributed caching system with Redis cluster
# Uses powerful model for complex architecture

# 4. Check token savings
~/.claude/scripts/token-efficient-config.sh status
# Savings this session: 67%
```

### Example 7: Continuous Integration Setup

```bash
# 1. Create CI/CD sprint
/qnew-enhanced Setup GitHub Actions CI/CD pipeline with testing, linting, and deployment

# 2. QDIRECTOR creates tasks:
# - CI_DESIGN: Pipeline architecture
# - TEST_SETUP: Test automation
# - LINT_CONFIG: Code quality rules
# - DEPLOY_IMPL: Deployment scripts

# 3. Execute with validation
/qdirector-enhanced --validate-each-step

# 4. Test the pipeline
/qvalidate-framework --ci-mode
# Simulates CI environment locally
```

### Example 8: Documentation Generation

```bash
# 1. Analyze codebase for documentation needs
/qdirector-enhanced Generate comprehensive API documentation

# 2. Agents work in parallel:
# - Extract JSDoc comments
# - Generate OpenAPI specs
# - Create usage examples
# - Build interactive docs

# 3. Validate documentation
/qcheck-enhanced --docs-only
# Ensures all public APIs documented
```

### Example 9: Performance Optimization Sprint

```bash
# 1. Profile application
/qvalidate-framework --mode performance --baseline

# 2. Create optimization plan
/qnew-enhanced Optimize API response time from 800ms to under 200ms

# 3. Execute optimizations
/qdirector-enhanced
# - Implements caching layer
# - Optimizes database queries
# - Adds request pooling
# - Implements lazy loading

# 4. Verify improvements
/qvalidate-framework --mode performance --compare baseline
# Output: Average response time: 187ms (76% improvement)
```

### Example 10: Multi-Repository Workflow

```bash
# 1. Coordinate changes across repos
/qnew-enhanced Implement shared authentication across frontend, backend, and mobile apps

# 2. QDIRECTOR manages:
# - Backend: JWT service implementation
# - Frontend: Auth context and hooks
# - Mobile: Native auth integration
# - Shared: Auth types package

# 3. Validate integration
/qvalidate-framework --cross-repo --repos frontend,backend,mobile

# 4. Coordinate deployment
/qgit --all-repos feat: implement unified authentication system
```

Each example demonstrates real-world usage patterns and the power of the
QDIRECTOR orchestration system with intelligent agent routing and automatic
validation.
