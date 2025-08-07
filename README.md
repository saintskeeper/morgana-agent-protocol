# Morgana Agent Protocol 🚀

> A comprehensive agent orchestration system with parallel execution,
> intelligent workflows, and quality assurance for modern development teams.

[![Version](https://img.shields.io/badge/version-2.0-blue.svg)]()
[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![Token Efficient](https://img.shields.io/badge/token%20efficient-beta-orange.svg)]()

## 📋 Table of Contents

- [🚀 Quickstart](#-quickstart)
- [📚 Commands Reference](#-commands-reference)
  - [Planning & Sprint Management](#-planning--sprint-management)
  - [Development](#-development)
  - [Validation & Quality](#-validation--quality)
  - [Utilities](#-utilities)
- [🔄 Common Workflows](#-common-workflows)
- [🎯 Best Practices](#-best-practices)
- [💡 Tips & Tricks](#-tips--tricks)
- [🛠️ Configuration](#️-configuration)
- [🔧 Troubleshooting](#-troubleshooting)
- [📖 Additional Resources](#-additional-resources)

## 🚀 Quickstart

Get up and running with Morgana Agent Protocol in under 5 minutes.

### Prerequisites

- **Git** installed on your system
- **Claude Code** CLI installed and authenticated
- **Go** 1.21+ for Morgana Protocol binary
- **Python** 3.8+ for bridge integration
- **macOS** (current configuration is optimized for macOS)
- Optional: `gofmt`, `prettier`, and `pre-commit` for formatting features

### Installation

```bash
# 1. Clone the Morgana Agent Protocol repository to your home directory
git clone git@github.com:saintskeeper/morgana-agent-protocol.git ~/.claude

# 2. Make scripts executable
chmod +x ~/.claude/setup-local.sh ~/.claude/test-hooks.sh

# 3. Install git hooks for automated workflows (recommended)
~/.claude/setup-local.sh
```

### First Command

Test that everything is working:

```bash
# Run a simple validation to ensure Claude commands are accessible
claude /rules-of-theroad

# Test the hooks functionality
~/.claude/test-hooks.sh
```

### Verification

✅ You're ready when you see:

- "✅ Post-checkout hook installed!" after running setup
- "✅ Hook test complete!" after running the test script
- Claude responds to slash commands like `/morgana-check`

### What's Next?

- **Enable token-efficient mode** (saves 14-70% on API costs):

  ```bash
  ~/.claude/scripts/token-efficient-config.sh enable
  ```

- **Try your first enhanced workflow**:

  ```bash
  # Create a simple utility function with auto-validation
  /morgana-code Create a date formatting utility function
  ```

- **Explore the command reference**: See all available commands with
  `/enhanced-quick-reference`

💡 **Tip**: The system automatically routes tasks to the optimal AI model based
on complexity. Simple tasks use efficient models, while complex architecture
work uses more powerful ones.

## 📚 Commands Reference

Morgana Agent Protocol commands are organized by workflow to help you accomplish
your development tasks efficiently. Each command leverages parallel agent
execution for maximum efficiency and reliability.

### 🎯 Planning & Sprint Management

#### `/morgana-plan` - Sprint Planning Generator

**Purpose**: Generate structured sprint plans with clear tasks, dependencies,
and exit criteria **Usage**: `/morgana-plan [project requirements]` **Model**:
`gemini-2.5-pro` or `o3` for comprehensive planning **Example**:

```bash
/morgana-plan Create authentication system with OAuth and JWT
# Generates: sprint-2024-01-15-authentication.md with 8 prioritized tasks
```

#### `/morgana-validate` - Technical Validation & Refinement

**Purpose**: Validate sprint plans against codebase patterns and technical
feasibility **Usage**: `/morgana-validate --sprint [sprint-file]` **Model**:
`pro` or Claude Opus for analysis **Options**:

- `--sprint`: Path to sprint plan file
- Analyzes code patterns, validates dependencies, identifies risks

**Example**:

```bash
/morgana-validate --sprint sprint-2024-01-15-authentication.md
# Output: Enhanced task definitions with codebase context and risk mitigation
```

#### `/morgana-director` - Master Orchestration System

**Purpose**: Orchestrate complex multi-task workflows with intelligent retry and
validation **Usage**: `/morgana-director [task description]` **Features**:

- Parallel task execution
- Automatic validation
- Smart retry with model escalation
- Human-in-the-loop for critical decisions

**Example**:

```bash
/morgana-director build complete authentication system with OAuth, JWT, and 2FA
# Orchestrates: Sprint planning → Implementation → Testing → Validation
```

### 💻 Development

#### `/morgana-code` - Code Implementation

**Purpose**: Implement features following project standards and best practices
**Usage**: Automatically invoked by MORGANA-DIRECTOR or used directly **Model**:
Complexity-based selection (Claude 3.7 → Claude 4 → GPT-4.1) **Features**:

- Pre-commit hook execution
- Race condition testing (Go)
- Type checking (TypeScript)
- Automatic formatting

**Example**:

```bash
/morgana-code implement user profile service with avatar upload
# Runs: pre-commit hooks → implementation → testing → formatting
```

#### `/morgana-test` - Comprehensive Test Generation

**Purpose**: Create thorough test suites with edge case coverage **Usage**:
`/morgana-test generate [type] --file [path]` **Model**: `o3-mini` or
`gemini-2.5-flash` **Options**:

- `unit` - Unit tests for functions/methods
- `integration` - Component interaction tests
- `e2e` - End-to-end user workflows
- `edge-cases` - Boundary and error scenarios

**Example**:

```bash
/morgana-test generate unit --file src/auth/jwt.service.ts
# Creates: Comprehensive unit tests with 90%+ coverage
```

### ✅ Validation & Quality

#### `/morgana-check` - Comprehensive Code Validation

**Purpose**: Validate code against best practices and security standards
**Usage**: Automatically triggered after code generation **Output**: Structured
YAML report with:

- Must-fix issues (blocking)
- Should-fix issues (recommended)
- Consider items (optional)
- Retry recommendations

**Example**:

```yaml
validation_report:
  pass_rate: 75%
  must_fix:
    - issue: "SQL injection vulnerability"
      location: "AuthService.ts:45"
  ready_for_merge: false
```

#### `/morgana-check-function` - Function-Level Validation

**Purpose**: Deep analysis of function quality, complexity, and maintainability
**Usage**: `/morgana-check-function [function-name]` or `--file [path]`
**Metrics**:

- Cyclomatic complexity (target: ≤10)
- Line count (target: ≤50)
- Parameter count (target: ≤3)
- Single responsibility adherence

**Example**:

```bash
/morgana-check-function processPayment
# Output: Complexity score: 12, Recommendation: refactor
```

#### `/morgana-check-tests` - Test Quality Validation

**Purpose**: Ensure tests are comprehensive, maintainable, and effective
**Usage**: `/morgana-check-tests [test-file]` or `--dir [directory]`
**Metrics**:

- Line coverage (target: ≥80%)
- Branch coverage (target: ≥75%)
- Test structure (AAA pattern)
- Mock quality

**Example**:

```bash
/morgana-check-tests PaymentService.test.ts
# Output: Coverage: 85%, Issues: missing timeout tests
```

#### `/morgana-validate-all` - Unified Validation System

**Purpose**: Orchestrate all validation commands for comprehensive quality
assurance **Usage**: `/morgana-validate-all --mode [quick|standard|deep]`
**Modes**:

- `quick`: Fast validation for development (~30s)
- `standard`: Comprehensive pre-commit validation (~2min)
- `deep`: Full analysis with security scanning (~5min)

**Example**:

```bash
/morgana-validate-all --mode standard --task-id AUTH_IMPL
# Runs: syntax → functions → tests → integration → security
```

### 🚀 Agent Orchestration with Morgana Protocol

#### Morgana Protocol - Parallel Agent Execution

**Purpose**: Execute multiple specialized agents in parallel with true Go
concurrency **Features**:

- 🚀 Goroutine-based parallel execution
- 🔍 OpenTelemetry tracing for observability
- ⏱️ Per-agent timeout configuration
- 🧪 Comprehensive integration testing
- 🐍 Python bridge for Claude Code integration

**Usage**:

```bash
# Single agent execution
morgana -- --agent code-implementer --prompt "implement auth service"

# Parallel execution with JSON
echo '[
  {"agent_type":"code-implementer","prompt":"implement feature"},
  {"agent_type":"test-specialist","prompt":"write tests"}
]' | morgana --parallel

# With configuration file
morgana --config morgana.yaml -- --agent sprint-planner --prompt "plan sprint"
```

**Testing Morgana Protocol**:

```bash
# Run all Morgana Protocol tests
cd ~/.claude/tests/morgana
./test-morgana-integration.sh

# Test director workflow
./test-morgana-director-workflow.sh

# Test specific agent execution
./test-qdirector-morgana.sh

# Test the Python bridge integration
cd ../../
morgana -- --agent code-implementer --prompt "test integration"
```

**Example Integration Test**:

```bash
# Set up test environment
export TEST_MODE=success

# Run integration test suite
make test-integration

# Check test coverage
go test -v -tags=integration -cover ./...
```

### 🔧 Utilities

#### `/morgana-commit` - Git Operations

**Purpose**: Add, commit with semantic messages, and push changes **Usage**:
`/morgana-commit [commit message]` **Features**:

- Semantic commit format (feat/fix/chore)
- Pre-commit validation
- No Claude attribution in commits

**Example**:

```bash
/morgana-commit feat: implement JWT authentication service
# Stages all changes → Creates semantic commit → Pushes to remote
```

#### `/qux` - UX Testing Scenarios

**Purpose**: Generate comprehensive user testing scenarios **Usage**: `/qux`
after implementing a feature **Output**: Prioritized list of test scenarios from
user perspective

**Example**:

```bash
/qux
# Output: 15 prioritized test scenarios for authentication flow
```

#### `/qprompt` - Prompt Template Helper

**Purpose**: Structure requests using token-efficient patterns **Usage**:
`/qprompt [task-type] - [description]` **Task Types**:

- `simple` - Direct execution
- `analyze` - Code review/analysis
- `implement` - Feature development
- `test` - Test generation
- `plan` - Sprint planning

**Example**:

```bash
/qprompt analyze - review auth.service.ts for security vulnerabilities
# Formats request with optimal structure for analysis
```

#### `/qtoken-efficient` - Token Optimization Management

**Purpose**: Enable/disable Anthropic's beta token-efficient mode **Usage**:
`/qtoken-efficient [enable|disable|status]` **Benefits**: 14-70% token reduction
with Claude 3.7 Sonnet **Note**: Beta feature - not compatible with Claude 4
models

**Example**:

```bash
~/.claude/scripts/token-efficient-config.sh enable
# Activates token-efficient mode for compatible operations
```

### 🔄 Workflow Integration

Commands work together in an intelligent workflow:

1. **Planning Phase**

   ```
   /morgana-plan → /morgana-validate → /morgana-director
   ```

2. **Implementation Phase**

   ```
   /morgana-code → automatic /morgana-check-function validation
   /morgana-test → automatic /morgana-check-tests validation
   ```

3. **Validation Phase**

   ```
   /morgana-validate-all orchestrates:
   → /morgana-check → /morgana-check-function → /morgana-check-tests
   ```

4. **Completion Phase**
   ```
   /morgana-commit with pre-commit /morgana-validate-all
   ```

### 🤖 Model Selection Strategy

Commands automatically select optimal models based on task complexity:

| Task Type    | Complexity | Primary Model       | Token-Efficient |
| ------------ | ---------- | ------------------- | --------------- |
| Planning     | High       | `gemini-2.5-pro`    | No              |
| Complex Code | High       | `claude-4-opus`     | No              |
| Simple Code  | Low        | `claude-3-7-sonnet` | Yes             |
| Testing      | Any        | `o3-mini`           | Yes             |
| Validation   | Any        | `claude-3-7-sonnet` | Yes             |

### 📊 Validation Severity Levels

- **MUST_FIX**: Blocks completion (security, data corruption, breaking changes)
- **SHOULD_FIX**: Retry recommended (performance, complexity, poor patterns)
- **CONSIDER**: Optional improvements (style, minor optimizations)

Each command is designed to work standalone or as part of the orchestrated
MORGANA-DIRECTOR workflow, providing flexibility for both automated and manual
development processes.

### 🧪 Integration Testing

The Morgana Protocol includes comprehensive integration tests:

| Test Type   | Coverage | Description                                             |
| ----------- | -------- | ------------------------------------------------------- |
| Task Client | ✅       | Tests Python bridge execution, timeouts, error handling |
| Adapter     | ✅       | Tests agent orchestration, concurrent execution         |
| Bridge      | ✅       | Tests multi-language communication (Go ↔ Python)       |
| E2E         | ✅       | Tests complete agent workflow with real Task tool       |

## 🔄 Common Workflows

### Complete Feature Development Flow

```bash
# 1. Plan the sprint
/morgana-plan Create user authentication system with JWT

# 2. Validate and enrich the plan
/morgana-validate --sprint sprint-2025-01-auth.md

# 3. Execute with MORGANA-DIRECTOR orchestration
/morgana-director
- Load sprint plan
- Execute tasks with automatic retry
- Validate outputs at each stage

# 4. Commit changes
/morgana-commit "feat: implement JWT authentication system"
```

### Quick Code Review

```bash
# For focused file review
/morgana-check-function auth_service.go

# For comprehensive validation
/morgana-validate-all --path ./src/auth/

# For security-focused review
/morgana-check --focus security
```

### Test Generation Workflow

```bash
# Generate tests for specific function
/qtest --function authenticate --file auth.go

# Validate test quality
/qcheckt-enhanced auth_test.go

# Run full test suite validation
/qvalidate-framework --tests-only
```

### Documentation Cleanup

```bash
# Organize AI docs by commit type
./scripts/qsweep.sh --filter feat

# Enable token-efficient mode
./scripts/token-efficient-config.sh enable

# Validate configuration
./scripts/validate-claude.sh
```

### Parallel Task Execution

```yaml
# In QDIRECTOR
parallel_tasks:
  - Task(subagent_type="validation-expert", prompt="Audit auth code")
  - Task(subagent_type="code-implementer", prompt="Research best practices")
  - Task(subagent_type="test-specialist", prompt="Plan test strategy")
```

## 🎯 Best Practices

### 1. Safe File Editing

- **Always use MultiEdit for critical files** like CLAUDE.md to prevent
  truncation
- **Test with grep first** to ensure unique string matches before editing
- **Include context** - match at least 2-3 lines for safer edits
- **Keep backups** before major configuration changes

### 2. Efficient Command Usage

- **Start with planning commands** (`/qnew-enhanced`, `/qplan-enhanced`) before
  implementation
- **Use token-efficient mode** for high-volume operations (saves 14-70% tokens)
- **Leverage parallel execution** in QDIRECTOR for independent tasks
- **Choose appropriate models** based on task complexity (see model optimization
  guide)

### 3. Agent Orchestration

- **Single responsibility principle** - each agent excels at one thing
- **Minimal context passing** - request only necessary information
- **Structured outputs** - always use QDIRECTOR-compatible formats
- **Clear error handling** - include retry recommendations

### 4. Hook Configuration

- **Language-specific formatters** automatically run on file edits
- **Validation hooks** ensure CLAUDE.md integrity after changes
- **Branch creation hooks** trigger automatic documentation cleanup
- **Test all hooks** with `test-hooks.sh` before relying on them

### 5. Sprint Management

- **Task sizing** - keep tasks to 2-4 hour implementation windows
- **Clear dependencies** - make task relationships explicit
- **Exit criteria** - define measurable success conditions
- **Priority tagging** - use P0-P3 levels consistently

## 💡 Tips & Tricks

### Performance Optimization

- **Cache agent outputs** between retry attempts to save tokens
- **Run independent tasks in parallel** using QDIRECTOR
- **Use flash models** for simple tasks, reserve pro models for complex analysis
- **Enable token-efficient mode** for 14-70% reduction in token usage

### Advanced Agent Usage

- **Chain agents intelligently**: planning → implementation → testing →
  validation
- **Share minimal context** between agents to preserve tokens
- **Use agent-specific models** configured in QDIRECTOR for optimal performance
- **Create custom agents** by adding markdown files to `/agents/` directory

### Smart Model Selection

```yaml
# Complexity-based routing
simple_task: "gemini-2.5-flash"
code_generation: "gpt-4.1"
deep_analysis: "gemini-2.5-pro"
security_audit: "o3"
comprehensive_planning: "o3" or "gemini-2.5-pro"
```

### Efficient Sprint Planning

- **Break down epics** into 2-week sprints maximum
- **Front-load risky tasks** to identify blockers early
- **Define clear interfaces** between components for parallel work
- **Include buffer time** for validation and iteration

### Code Quality Shortcuts

```bash
# Quick quality check before commit
alias qqa='./scripts/validate-claude.sh && /qvalidate-framework --quick'

# Auto-organize AI docs after session
alias qclean='./scripts/qsweep.sh --auto'

# Full validation pipeline
alias qfull='/qvalidate-framework --comprehensive'
```

### Context Management

- **Use .claudeignore** to exclude irrelevant files from context
- **Reference specific files** in prompts rather than "check the codebase"
- **Keep CLAUDE.md focused** - move project-specific rules to local .claude/
- **Use grep/glob first** to find files, then read only what's needed

### Debugging Commands

- **Add --verbose flag** to see detailed agent reasoning
- **Use --dry-run** to preview what would be executed
- **Check intermediate outputs** in task working directories
- **Enable debug logging** in QDIRECTOR for state machine visibility

### Integration Tips

- **Linear Integration**: Set project context with prep commands
- **CI/CD**: Use validation commands in pre-commit hooks
- **IDE Integration**: Map commands to keyboard shortcuts
- **Team Workflows**: Share sprint plans via version control

## 🛠️ Configuration

### Essential Settings

The system works out-of-the-box, but you can customize behavior through
`settings.json`:

```json
{
  "hooks": {
    "postToolUse": ["./hooks/post-edit.sh"],
    "userPromptSubmit": ["./hooks/qsweep.sh"]
  },
  "env": {
    "CLAUDE_TOKEN_EFFICIENT_MODE": "true",
    "CLAUDE_BETA_HEADER": "token-efficient-tools-2025-02-19"
  }
}
```

### Token-Efficient Mode (Beta)

Save 14-70% on API costs with Claude 3.7 Sonnet:

```bash
# Enable
~/.claude/scripts/token-efficient-config.sh enable

# Check status
~/.claude/scripts/token-efficient-config.sh status

# Disable
~/.claude/scripts/token-efficient-config.sh disable
```

⚠️ **Note**: Only works with Claude 3.7 Sonnet. Other models operate normally.

### CLAUDE.md Customization

Add project-specific instructions to `CLAUDE.md`:

- Commit message formats
- Code style preferences
- Project-specific rules
- Team conventions

### Advanced Configuration

For detailed configuration options:

- **Hooks Documentation**: See `hooks/README.md`
- **Script Options**: See `scripts/README.md`
- **Template Customization**: See `templates/README.md`

## 🔧 Troubleshooting Guide

### CLAUDE.md Gets Truncated

**Problem**: File becomes corrupted or sections disappear after edits

**Solution**:

1. Use MultiEdit instead of Edit for changes
2. Restore from template: `cp templates/CLAUDE.template.md CLAUDE.md`
3. Run validation: `./scripts/validate-claude.sh`
4. Use unique section markers for safer edits

### Hooks Not Running

**Problem**: Post-edit formatting not happening automatically

**Solution**:

1. Check hook installation: `./setup-local.sh`
2. Verify settings.json has correct hook paths
3. Test specific hook: `./test-hooks.sh post-edit`
4. Check file permissions: `chmod +x hooks/*.sh`

### QDIRECTOR Task Failures

**Problem**: Tasks stuck in RETRY state or failing repeatedly

**Solution**:

1. Check task dependencies are properly defined
2. Verify agent has necessary tool access
3. Review validation output for specific issues
4. Use manual retry with different model:
   ```yaml
   retry_with_model: "gemini-2.5-pro"
   additional_context: "Previous attempt failed due to..."
   ```

### Token Limit Exceeded

**Problem**: Commands failing due to context size

**Solution**:

1. Enable token-efficient mode: `./scripts/token-efficient-config.sh enable`
2. Use focused commands instead of comprehensive ones
3. Split large tasks into smaller subtasks
4. Choose appropriate models for task complexity

### Model Not Available

**Problem**: Specified model returns errors

**Solution**:

1. Check available models: `/qvalidate-framework --list-models`
2. Verify API keys are configured correctly
3. Use fallback models in QDIRECTOR configuration
4. Check model-specific context limits

### Git Hooks Conflict

**Problem**: Local git hooks interfere with Claude hooks

**Solution**:

1. Backup existing hooks: `cp .git/hooks/* .git/hooks.backup/`
2. Integrate Claude hooks with existing ones
3. Use hook chaining in .git/hooks scripts
4. Test combined functionality thoroughly

## 📖 Additional Resources

### 📚 Detailed Documentation

- **[Repository Structure](docs/repository-structure.md)** - Full directory
  layout and file purposes
- **[Agent Architecture](agents/README.md)** - Deep dive into specialized agents
- **[Hook System](hooks/README.md)** - Advanced hook configuration
- **[Script Reference](scripts/README.md)** - All utility scripts explained
- **[Template Guide](templates/README.md)** - Customization templates

### 🔗 External Resources

- **[Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)** -
  Official docs
- **[GitHub Repository](https://github.com/saintskeeper/morgana-agent-protocol)** -
  Source code
- **[Issue Tracker](https://github.com/saintskeeper/morgana-agent-protocol/issues)** -
  Report bugs
- **[Claude Code Updates](https://github.com/anthropics/claude-code/releases)** -
  Latest features

### 🎓 Learning Resources

- **[Beta Features Guide](docs/beta-features-guide.md)** - Understanding beta
  features
- **[Model Capabilities](docs/model-comparison.md)** - Detailed model
  comparisons
- **[Security Best Practices](guidelines/security.md)** - Secure coding
  guidelines
- **[Performance Optimization](guidelines/performance.md)** - Speed and
  efficiency tips

### 🤝 Contributing

Want to improve Morgana Agent Protocol?

1. Fork the repository
2. Create a feature branch
3. Follow the contribution guidelines
4. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed instructions.

### 📞 Getting Help

- **Quick questions**: Check [Troubleshooting](#-troubleshooting) first
- **Bug reports**:
  [Open an issue](https://github.com/saintskeeper/morgana-agent-protocol/issues)
- **Feature requests**: Use the issue template
- **Community**: Join discussions in issues

### 🏛️ Architecture Overview

```
Morgana Agent Protocol
├── Commands (User Interface)
│   ├── Planning & Sprint Management
│   ├── Development & Testing
│   └── Validation & Quality
├── Morgana Director (Orchestration Layer)
│   ├── Parallel Task Execution
│   ├── Agent Routing & Load Balancing
│   └── Validation Pipeline
├── Specialized Agents (Execution Layer)
│   ├── sprint-planner
│   ├── code-implementer
│   ├── test-specialist
│   └── validation-expert
└── Infrastructure (Support Layer)
    ├── Go Binary & Python Bridge
    ├── OpenTelemetry Tracing
    └── Scripts & Utilities
```

### 🔄 Execution Flow with Claude REPL

This diagram shows how Morgana Protocol integrates with a running Claude Code
session:

```mermaid
graph TB
    subgraph "Claude Code Environment"
        REPL["Claude REPL<br/>(Interactive Session)"]
        TASK["Task() Function<br/>(Built-in)"]
    end

    subgraph "User Shell"
        USER["User Command<br/>$ morgana --agent code-implementer<br/>--prompt 'implement auth'"]
        WRAPPER["Shell Wrapper<br/>(agent-adapter-wrapper.sh)"]
    end

    subgraph "Morgana Protocol Core"
        CLI["morgana CLI<br/>(Go Binary)"]
        ADAPTER["Adapter<br/>(Agent Type → Task)"]
        ORCH["Orchestrator<br/>(Sequential/Parallel)"]
        CLIENT["Task Client<br/>(Go)"]
    end

    subgraph "Python Bridge Layer"
        BRIDGE["task_bridge.py<br/>(Python Script)"]
    end

    subgraph "Event & Monitoring System"
        EVENTBUS["Event Bus<br/>(5M events/sec)"]
        IPC["IPC Client<br/>(Unix Socket)"]
        MONITOR["morgana-monitor<br/>(Daemon Process)"]
        TUI["TUI Display<br/>(Real-time Updates)"]
    end

    subgraph "Observability Stack"
        OTEL["OpenTelemetry<br/>(Tracing)"]
        JAEGER["Jaeger UI<br/>(:16686)"]
        PROM["Prometheus<br/>(:9090)"]
        GRAFANA["Grafana<br/>(:3000)"]
    end

    %% Main execution flow
    USER -->|1. Execute| CLI
    WRAPPER -.->|Source| USER
    CLI -->|2. Parse & Validate| ADAPTER
    ADAPTER -->|3. Load Prompt<br/>Template| ADAPTER
    ADAPTER -->|4. Select Model<br/>(complexity-based)| ADAPTER
    ADAPTER -->|5. Create Task| ORCH
    ORCH -->|6. Execute| CLIENT
    CLIENT -->|7. Shell Out<br/>(JSON via stdin)| BRIDGE
    BRIDGE -->|8. Call Task()<br/>'general-purpose'| TASK
    TASK -->|9. Execute in<br/>Claude Code| REPL
    REPL -->|10. Return Result| TASK
    TASK -->|11. JSON Response| BRIDGE
    BRIDGE -->|12. stdout| CLIENT
    CLIENT -->|13. Parse Result| ORCH
    ORCH -->|14. Return| ADAPTER
    ADAPTER -->|15. Output| CLI

    %% Event flow
    ADAPTER -.->|Publish Events| EVENTBUS
    ORCH -.->|Task Events| EVENTBUS
    EVENTBUS -.->|Forward| IPC
    IPC -.->|Unix Socket<br/>/tmp/morgana.sock| MONITOR
    MONITOR -.->|Update| TUI

    %% Telemetry flow
    ADAPTER -.->|Spans| OTEL
    ORCH -.->|Traces| OTEL
    OTEL -.->|Export| JAEGER
    OTEL -.->|Metrics| PROM
    PROM -.->|Visualize| GRAFANA

    %% Parallel execution branch
    ORCH -->|Parallel Mode| POOL["Goroutine Pool<br/>(Max 5 concurrent)"]
    POOL -->|Multiple Tasks| CLIENT

    style REPL fill:#e1f5fe
    style TASK fill:#e1f5fe
    style BRIDGE fill:#fff3e0
    style EVENTBUS fill:#f3e5f5
    style TUI fill:#f3e5f5
    style MONITOR fill:#f3e5f5
```

#### Key Flow Points:

1. **Task Execution Pipeline**: User invokes `morgana` → Go validates & loads
   agent template → Shells to Python bridge → Calls Claude's Task() → Returns
   result

2. **Python Bridge (Current Bottleneck)**: Each task spawns new Python process
   with JSON serialization overhead. Acceptable for current use but limits
   high-frequency operations.

3. **Event System**: Real-time event bus forwards to monitor daemon via Unix
   socket, enabling live TUI updates across multiple instances.

4. **Observability**: Full OpenTelemetry tracing with Jaeger UI, Prometheus
   metrics, and Grafana dashboards for comprehensive monitoring.

5. **Parallel Execution**: Goroutine pool (default 5) manages concurrent tasks,
   each still requiring separate Python process.

---

_Built with ❤️ for the Morgana Agent Protocol community_

## 📘 Real-World Examples

For comprehensive examples demonstrating the power of QDIRECTOR orchestration:

- **[View All Examples](examples/README.md)** - 10 detailed scenarios
- **[Migration Guide](MIGRATION-GUIDE.md)** - For existing users
