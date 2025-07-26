# QNEW Command - Sprint Planning Generator

Understand all BEST PRACTICES listed in CLAUDE.md. Your code SHOULD ALWAYS
follow these best practices.

## Purpose

Generate structured sprint plans that can be executed by QDIRECTOR with clear
tasks, dependencies, and exit criteria.

## Workflow

### 1. Gather Requirements

When user provides project requirements:

- Identify key features and goals
- Break down into logical components
- Consider technical constraints
- Assess complexity and risks

### 2. Generate Sprint Plan

Create a sprint plan following the template at
`/templates/sprint-plan-template.md`:

```markdown
# Sprint Plan: [Descriptive Name]

## Sprint Overview

- Sprint ID: sprint-YYYY-MM-DD
- Duration: [Estimated based on complexity]
- Goal: [Clear, measurable objective]
- Success Criteria: [Overall deliverables]

## Task Definitions

For each identified task:

- Assign unique TASK_ID (e.g., AUTH_DESIGN, API_IMPL)
- Set appropriate priority (P0-P3)
- Define task type (architecture/implementation/testing/etc)
- List dependencies clearly
- Create specific, measurable exit criteria
- Include relevant context files
```

### 3. Task Decomposition Strategy

**Architecture Tasks** (use suffix \_DESIGN):

- System design documents
- API contract definitions
- Database schema design
- Security architecture
- Exit criteria: Documents created, reviewed, approved

**Implementation Tasks** (use suffix \_IMPL):

- Core feature development
- Service implementations
- API endpoints
- Business logic
- Exit criteria: Code complete, unit tests pass, follows conventions

**Testing Tasks** (use suffix \_TEST):

- Unit test suites
- Integration tests
- E2E test scenarios
- Performance tests
- Exit criteria: Coverage targets met, all tests pass

**Infrastructure Tasks** (use suffix \_INFRA):

- CI/CD setup
- Database migrations
- Deployment configs
- Monitoring setup
- Exit criteria: Scripts work, environments ready

### 4. Dependency Rules

**Sequential Dependencies**:

- Design → Implementation → Testing → Deployment
- Database schema → Data models → Services
- API contracts → Frontend & Backend implementation

**Parallel Opportunities**:

- Frontend and backend after API contracts
- Unit tests alongside implementation
- Documentation in parallel with development
- Independent features can run concurrently

### 5. Priority Guidelines

**P0-Critical**: Core functionality, blockers **P1-High**: Important features,
must-have for sprint **P2-Medium**: Nice-to-have, enhances functionality
**P3-Low**: Polish, optimizations, tech debt

### 6. Complexity Estimation

**Simple**: < 4 hours, well-defined, minimal dependencies **Medium**: 4-16
hours, some complexity, few dependencies **Complex**: 16-40 hours, significant
logic, multiple dependencies **Critical**: 40+ hours, architectural decisions,
high risk

### 7. Exit Criteria Best Practices

Make them SMART (Specific, Measurable, Achievable, Relevant, Time-bound):

- ❌ "Tests written"
- ✅ "Unit test coverage ≥ 80% for all new code"

- ❌ "API works"
- ✅ "All 5 API endpoints return correct data with < 200ms response time"

- ❌ "Secure implementation"
- ✅ "Passes security scan, implements JWT with 1hr expiry, rate limiting
  enabled"

### 8. Output Format

Save the sprint plan as `sprint-[date]-[feature].md` and inform user:

1. Total tasks created
2. Critical path identified
3. Estimated sprint duration
4. Key risks and dependencies
5. Suggested team allocation

### Example Usage

```
User: Create a sprint plan for user authentication with OAuth and JWT

QNEW:
1. Analyzes requirements
2. Creates 8 tasks:
   - AUTH_DESIGN (architecture)
   - OAUTH_IMPL (OAuth integration)
   - JWT_IMPL (JWT service)
   - AUTH_API (API endpoints)
   - AUTH_TEST (test suite)
   - AUTH_DOCS (documentation)
   - SECURITY_REVIEW (security audit)
   - AUTH_DEPLOY (deployment)
3. Defines dependencies and exit criteria
4. Saves as sprint-2024-01-15-authentication.md
5. Ready for QDIRECTOR execution
```

## Integration with QPLAN

After generating initial sprint plan:

- Use `/qplan` to validate technical approach
- Refine task definitions based on codebase analysis
- Adjust dependencies if conflicts found
- Update complexity estimates with actual code context
