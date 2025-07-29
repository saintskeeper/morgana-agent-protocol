# Claude Code Configuration

A comprehensive Claude Code enhancement system with automated workflows, quality
assurance, and development commands.

## 🏗️ Repository Structure

```
.
├── commands/           # Claude Code slash commands
├── hooks/             # Event-triggered automation scripts
├── scripts/           # Utility scripts for maintenance
├── templates/         # Reusable configuration templates
├── guidelines/        # Best practices documentation
├── settings.json      # Claude Code hook configuration
├── CLAUDE.md         # Global user instructions
├── CLAUDE_SAFETY.md  # Safety guidelines for configuration editing
├── setup-local.sh    # Local git hooks installer
└── test-hooks.sh     # Hook testing utility
```

## 📜 Scripts (`/scripts/`)

### `validate-claude.sh`

Validates CLAUDE.md integrity after edits:

- Ensures minimum 200 lines
- Validates required sections
- Prevents configuration corruption

### `token-efficient-config.sh`

Manages the token-efficient tool use beta feature:

- Enable/disable token-efficient mode
- Check current status and compatibility
- Reduces output tokens by 14-70% with Claude Sonnet 3.7

### `qsweep.sh`

Organizes AI documentation with conventional commit patterns:

- Moves docs to structured `ai-docs/` hierarchy
- Supports filtering by ticket ID or commit type
- Categories: feat, fix, docs, chore, refactor, test, build, ci, perf, style
- Dry-run mode available

## 🪝 Hooks (`/hooks/`)

### `post-edit.sh`

Main dispatcher that routes to file-specific formatters:

- Supports Go, Markdown, TypeScript/JavaScript, YAML, Rust
- Uses appropriate formatters (Prettier, gofmt, rustfmt)

### `post-branch-create.sh`

Automatically runs documentation cleanup on new branches

### `post-edit-markdown.sh`

Formats Markdown files:

- Fixes end-of-file newlines
- Removes trailing whitespace
- Runs Prettier with prose-wrap

### `post-edit-go.sh`

Auto-formats Go files:

- Runs `gofmt` for formatting
- Runs `goimports` for import organization

## ⚡ Commands (`/commands/`)

### Enhanced Workflow Commands

#### Planning & Sprint Management

- **`/qnew-enhanced`** - Advanced sprint planning generator

  - Creates structured sprint plans with tasks, dependencies, and exit criteria
  - Outputs QDIRECTOR-compatible YAML format
  - Models: `gemini-2.5-pro` or `o3` for comprehensive planning

- **`/qplan-enhanced`** - Technical validation & sprint refinement
  - Validates sprint plans for technical feasibility
  - Enriches tasks with codebase context and patterns
  - Generates dependency graphs and risk assessments

#### Master Orchestration

- **`/qdirector-enhanced`** - Intelligent task orchestrator
  - Manages specialized agents for complex workflows
  - Automatic retry logic with smart model selection
  - Parallel execution for independent tasks
  - State machine: PENDING → READY → IN_PROGRESS → VALIDATION → COMPLETED

### Validation Framework

#### Comprehensive Validations

- **`/qcheck-enhanced`** - Full code validation suite

  - Security, performance, and quality checks
  - Structured YAML output for automated parsing
  - Severity levels: MUST_FIX, SHOULD_FIX, CONSIDER
  - Pass rate thresholds: 90%+ auto-approve, 70-89% retry, <70% human review

- **`/qcheckf-enhanced`** - Function-level analysis

  - Complexity metrics (cyclomatic, cognitive, nesting)
  - Design principles validation
  - Performance characteristics analysis

- **`/qcheckt-enhanced`** - Test quality validation

  - Coverage analysis (line, branch, function)
  - Test effectiveness metrics
  - Anti-pattern detection

- **`/qvalidate-framework`** - Unified validation orchestrator
  - Aggregates all validation results
  - Progressive modes: quick (dev), standard (pre-commit), deep (pre-deploy)
  - Smart retry recommendations based on failure patterns

### Core Workflow Commands

- **`/qcode`** - Implementation with auto-validation

  - Triggers `/qcheckf-enhanced` automatically
  - Models: `gpt-4.1` or `gemini-2.5-flash`

- **`/qgit`** - Git operations with pre-commit validation
  - Runs `/qvalidate-framework --mode standard`
  - Semantic commit messages

### Specialized Tools

- **`/qux`** - UX testing scenario generation
- **`/rules-of-theroad`** - Claude Commands system guide
- **`/important-instruction-reminders`** - Core behavioral constraints
- **`/enhanced-quick-reference`** - Quick guide to enhanced workflow

## 🤖 Specialized Agents (`/agents/`)

### Agent Architecture

The QDIRECTOR system leverages specialized agents for focused tasks:

#### `sprint-planner`

- **Purpose**: Expert requirements decomposition and sprint planning
- **Tools**: Read, Write, TodoWrite, Grep, Glob
- **Outputs**: QDIRECTOR-compatible sprint plans with:
  - Task decomposition (2-4 hour chunks)
  - Clear dependencies and priority levels (P0-P3)
  - Acceptance criteria and exit conditions
  - Risk identification and mitigation

#### `code-implementer`

- **Purpose**: Clean, secure code implementation following conventions
- **Tools**: Read, Write, Edit, MultiEdit, Bash, Grep, Glob, LS
- **Principles**:
  - SOLID principles and design patterns
  - Security by default (input validation, parameterized queries)
  - Performance conscious implementation
  - Convention adherence to project patterns

#### `test-specialist`

- **Purpose**: Comprehensive test suite creation
- **Tools**: Read, Write, Edit, MultiEdit, Bash, Grep, Glob, mcp**zen**testgen
- **Focus**:
  - Behavior-driven testing (not implementation)
  - Edge case coverage
  - AAA pattern (Arrange, Act, Assert)
  - Test pyramid approach (unit > integration > E2E)

#### `validation-expert`

- **Purpose**: Multi-dimensional code validation
- **Tools**: Read, Grep, Glob, Bash, mcp**zen**codereview, mcp**zen**secaudit,
  mcp**zen**analyze
- **Validation Scope**:
  - Code quality (structure, complexity, maintainability)
  - Security (OWASP, injection prevention, auth)
  - Performance (algorithms, resource usage, queries)
  - Best practices compliance

## 📚 USER GUIDE - Intelligent Model Routing

### Understanding Complexity-Based Model Selection

QDIRECTOR automatically analyzes your tasks and routes them to the optimal model
based on complexity. This ensures you get the best balance of quality, speed,
and cost.

#### 🎯 Complexity Keywords & Model Selection

**SIMPLE TASKS → Claude 3.7 Sonnet (Token-Efficient)** Keywords that trigger
simple routing:

- `utility`, `helper`, `function`, `component`
- `convert`, `format`, `validate` (simple)
- `button`, `modal`, `form`, `tooltip`
- `config`, `constant`, `interface`

Examples:

```bash
"Create a date formatting utility"           # → Claude 3.7 Sonnet
"Add a button component"                     # → Claude 3.7 Sonnet
"Write a helper function to validate email"  # → Claude 3.7 Sonnet
```

**MODERATE TASKS → Claude 4 Sonnet** Keywords that trigger moderate routing:

- `api`, `service`, `integration`, `middleware`
- `authentication`, `authorization`, `database schema`
- `cache`, `queue`, `websocket`, `graphql`
- `business logic`, `workflow`, `state management`

Examples:

```bash
"Implement REST API with authentication"     # → Claude 4 Sonnet
"Create caching service with Redis"          # → Claude 4 Sonnet
"Build user authentication middleware"       # → Claude 4 Sonnet
```

**COMPLEX TASKS → Claude 4 Opus** Keywords that trigger complex routing:

- `architect`, `design system`, `refactor entire`
- `migrate`, `distributed`, `concurrent`, `parallel`
- `real-time`, `blockchain`, `machine learning`
- `performance critical`, `security critical`
- `custom algorithm`, `parser`, `compiler`

Examples:

```bash
"Design distributed caching system"          # → Claude 4 Opus
"Refactor entire payment architecture"       # → Claude 4 Opus
"Implement concurrent data processing"       # → Claude 4 Opus
```

#### 🔧 Manual Complexity Analysis

Test task complexity before execution:

```bash
# Analyze complexity
~/.claude/scripts/code-complexity-analyzer.sh analyze "your task description"

# Get model recommendation
~/.claude/scripts/code-complexity-analyzer.sh recommend "your task description" true
```

#### 💡 Pro Tips for Task Description

1. **Be Specific**: More details help accurate routing
2. **Use Keywords**: Include complexity indicators
3. **Mention Scale**: "entire system" vs "single function"
4. **State Requirements**: "high-performance" or "simple utility"

### Token-Efficient Mode Benefits

When enabled, simple tasks automatically use Claude 3.7 Sonnet with:

- **14-70% token reduction**
- **Faster response times**
- **Lower API costs**
- **Maintained quality**

Enable token-efficient mode:

```bash
~/.claude/scripts/token-efficient-config.sh enable
```

### Model Routing Examples

```bash
# Automatic routing based on complexity:
/qdirector-enhanced "Create a simple date formatter"
# → Routes to: Claude 3.7 Sonnet (simple, token-efficient)

/qdirector-enhanced "Build authentication API with JWT"
# → Routes to: Claude 4 Sonnet (moderate complexity)

/qdirector-enhanced "Design microservices architecture"
# → Routes to: Claude 4 Opus (complex, needs deep reasoning)
```

### Understanding Model Capabilities

| Model             | Best For                | Token Efficiency       | Use When                     |
| ----------------- | ----------------------- | ---------------------- | ---------------------------- |
| Claude 3.7 Sonnet | Simple tasks, utilities | Yes (14-70% reduction) | Clear, straightforward tasks |
| Claude 4 Sonnet   | Moderate complexity     | No (works normally)    | APIs, services, integrations |
| Claude 4 Opus     | Complex architecture    | No (works normally)    | System design, algorithms    |
| GPT-4.1           | General coding          | No                     | Fallback for moderate tasks  |
| Gemini 2.5 Flash  | Quick tasks             | No                     | Fast simple tasks            |

### Structured Prompts for Efficiency

Use these templates for optimal results:

**For Simple Tasks:**

```
Task: [specific action]
Input: [data/parameters]
Output: [expected format]
```

**For Complex Tasks:**

```
Implement: [feature]
Requirements:
- [requirement 1]
- [requirement 2]
Constraints:
- [constraint 1]
Architecture: [patterns to follow]
```

## 🧪 Experimental Features

### Token-Efficient Tool Use (Beta)

Reduces output tokens by 14-70% when using Claude Sonnet 3.7:

- **Enable**: `~/.claude/scripts/token-efficient-config.sh enable`
- **Status**: `~/.claude/scripts/token-efficient-config.sh status`
- **Disable**: `~/.claude/scripts/token-efficient-config.sh disable`

Configuration in `settings.json`:

```json
{
  "env": {
    "CLAUDE_TOKEN_EFFICIENT_MODE": "true",
    "CLAUDE_BETA_HEADER": "token-efficient-tools-2025-02-19"
  }
}
```

⚠️ **Limitations**: Only works with Claude Sonnet 3.7, not compatible with
Claude 4 models.

📖 **Learn More**: See [Beta Features Guide](docs/beta-features-guide.md) for
detailed information about what beta features mean and how they work.

## 🔧 Templates (`/templates/`)

### `claude-code-ai-assistant-template.md`

Complete AI assistant architecture template:

- Multi-agent system patterns
- Tool definitions and orchestration
- Performance tracking and security

### `CLAUDE.template.md`

Standard project configuration template:

- Linear integration guidelines
- Implementation best practices
- Documentation management

## ⚙️ Configuration

### `settings.json`

Hook configuration:

- **PostToolUse**: Formatting after edits, audio notifications
- **PreToolUse**: Validation before edits
- **UserPromptSubmit**: Documentation sweep
- **Stop**: Final cleanup before commits

### `CLAUDE.md`

Global user instructions:

- Custom commit message format
- macOS environment specification
- TODO addition guidelines

### `CLAUDE_SAFETY.md`

Safety guidelines for configuration editing:

- MultiEdit usage patterns
- Validation procedures
- Recovery processes

## 🚀 Setup

### Quick Start

```bash
# Clone the repository
git clone git@github.com:saintskeeper/claude-code-configs.git ~/.claude

# Install local git hooks (optional)
chmod +x ~/.claude/setup-local.sh
~/.claude/setup-local.sh
```

### Test Installation

```bash
# Test hooks functionality
chmod +x ~/.claude/test-hooks.sh
~/.claude/test-hooks.sh
```

## 💡 Usage Examples

### Development Workflow

1. Start with `/qplan-enhanced` to validate architectural approach
2. Use `/qnew-enhanced` for new feature development reminders
3. Run `/qcheck-enhanced` before major commits
4. Execute `/qcode` for final validation
5. Commit with `/qgit` for conventional commit format

### Documentation Management

- Documentation is automatically organized with `qsweep.sh`
- Use `--id SPE-249` to filter by ticket
- Use `--type feat` to filter by change type
- Run `--dry-run` to preview changes

### Code Quality

- All edits are automatically formatted via hooks
- CLAUDE.md changes are validated for integrity
- Audio notifications confirm command completion

## 🛡️ Safety Features

- **Pre-edit validation** prevents CLAUDE.md corruption
- **MultiEdit patterns** for safe large file editing
- **Section markers** for reliable configuration updates
- **Backup strategies** for configuration recovery

## 🚀 Enhanced Workflow Example

### Complete Development Cycle with QDIRECTOR

```bash
# 1. Create sprint plan from requirements
/qnew-enhanced Build secure user authentication with JWT and OAuth

# 2. Validate and enrich the plan technically
/qplan-enhanced --sprint sprint-2024-01-auth.md

# 3. Execute with QDIRECTOR orchestration
/qdirector-enhanced Execute sprint plan in sprint-2024-01-auth.md

# QDIRECTOR automatically:
# - Spawns specialized agents in parallel
# - Validates outputs with enhanced commands
# - Retries with focused context on failures
# - Tracks progress through state machine
# - Escalates blockers to human review
```

### Parallel Agent Execution

```yaml
# Example of parallel task execution
parallel_tasks:
  - Task(subagent_type="code-implementer", prompt="Implement JWT service")
  - Task(subagent_type="code-implementer", prompt="Create user model")
  - Task(subagent_type="test-specialist", prompt="Design auth test strategy")
```

## 📊 Validation Pipeline

```
Code Generation → Function Validation → Test Validation → Integration Check → Security Scan
     ↓                    ↓                    ↓                   ↓                ↓
  /qcode         /qcheckf-enhanced    /qcheckt-enhanced      /qcheck-enhanced   Ready
```

## 🎯 Key Benefits

1. **Intelligent Orchestration**: QDIRECTOR manages complex workflows with
   specialized agents
2. **Automated Quality Assurance**: Multi-layer validation with smart retry
   logic
3. **Parallel Execution**: Independent tasks run simultaneously for speed
4. **Structured Workflows**: From sprint planning to validated implementation
5. **Safety Mechanisms**: Configuration protection and automatic validation
6. **Model Optimization**: Right model for each task type
7. **Continuous Learning**: Metrics tracking for process improvement

---

_This enhanced system transforms Claude Code into an intelligent development
orchestrator with specialized agents, comprehensive validation, and automated
workflows while respecting user preferences for commit messages and
attribution._
