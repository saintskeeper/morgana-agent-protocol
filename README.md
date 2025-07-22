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

### Core Development Workflow

- **`/qcheck`** - Comprehensive code quality analysis for major changes
- **`/qcode`** - Implementation workflow checklist with testing
- **`/qgit`** - Git workflow with conventional commits
- **`/qplan`** - Architectural consistency validation
- **`/qnew`** - Best practices reminder for new development

### Specialized Analysis

- **`/qcheckf`** - Function-specific quality analysis
- **`/qcheckt`** - Test-specific quality analysis
- **`/qux`** - UX testing scenario generation

### System Documentation

- **`/rules-of-theroad`** - Claude Commands system guide
- **`/important-instruction-reminders`** - Core behavioral constraints

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

1. Start with `/qplan` to validate architectural approach
2. Use `/qnew` for new feature development reminders
3. Run `/qcheck` before major commits
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

## 🎯 Key Benefits

1. **Automated Quality Assurance**: Formatting, validation, and testing
2. **Streamlined Workflows**: Commands for all development phases
3. **Documentation Organization**: Automatic conventional commit categorization
4. **Safety Mechanisms**: Configuration protection and recovery
5. **Consistent Development**: Templates and best practices enforcement

---

_This system transforms Claude Code into a highly customized development
environment with automated workflows while respecting user preferences for
commit messages and attribution._
