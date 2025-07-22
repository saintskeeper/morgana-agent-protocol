# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

<!-- SECTION: Linear Integration -->
## Linear Integration

### Cached Linear Data
To avoid repeated API calls, Linear team and project information is cached in the `.linear/` directory:

<!-- SECTION: Claude Guidelines -->
# Claude Code Guidelines

## Implementation Best Practices

### 0 — Purpose

These rules ensure maintainability, safety, and developer velocity.
**MUST** rules are enforced by CI; **SHOULD** rules are strongly recommended.

---

### 1 — Before Coding

- **BP-1 (MUST)** Ask the user clarifying questions.
- **BP-2 (SHOULD)** Draft and confirm an approach for complex work.
- **BP-3 (SHOULD)** If ≥ 2 approaches exist, list clear pros and cons.

---

### 2 — While Coding

- **C-1 (MUST)** Follow TDD: scaffold stub -> write failing test -> implement.
- **C-2 (MUST)** Name functions with existing domain vocabulary for consistency.
- **C-3 (SHOULD NOT)** Introduce classes when small testable functions suffice.
- **C-4 (SHOULD)** Prefer simple, composable, testable functions.
- **C-5 (MUST)** Prefer branded `type`s for IDs
  ```ts
  type UserId = Brand<string, 'UserId'>   // ✅ Good
  type UserId = string                    // ❌ Bad
  ```
- **C-6 (MUST)** Use `import type { … }` for type-only imports.
- **C-7 (SHOULD NOT)** Add comments except for critical caveats; rely on self‑explanatory code.
- **C-8 (SHOULD)** Default to `type`; use `interface` only when more readable or interface merging is required.
- **C-9 (SHOULD NOT)** Extract a new function unless it will be reused elsewhere, is the only way to unit-test otherwise untestable logic, or drastically improves readability of an opaque block.
- **C-10 (MUST)** When adding TODOs to the codebase, ensure you're not deleting functionality from the codebase; instead, call analyze on your problem if you don't have context before making assumptions.

---

### 3 — After Coding

- **AC-1 (MUST)** Run unit tests for changed code.
- **AC-2 (MUST)** Verify existing tests still pass.
- **AC-3 (SHOULD)** Run integration tests for areas touched.
- **AC-4 (MUST)** Run lint and typecheck commands before committing.
- **AC-5 (SHOULD)** Review code diff before pushing.

---

### 4 — Documentation Management

#### qsweep Function
Organizes completed feature documentation following conventional commit types:

```bash
# Sweep all completed docs to ai-docs/completed/
./qsweep.sh

# Sweep specific ticket (e.g., SPE-249)
./qsweep.sh --id SPE-249

# Sweep by type (feat, fix, docs, chore, refactor, test)
./qsweep.sh --type feat

# Preview what would be moved without actually moving
./qsweep.sh --dry-run
```

#### AI Documentation Structure
```
ai-docs/
├── completed/
│   ├── feat/         # New features (e.g., feat/spe-249/)
│   ├── fix/          # Bug fixes (e.g., fix/spe-234/)
│   ├── docs/         # Documentation improvements
│   ├── chore/        # Maintenance tasks
│   ├── refactor/     # Code refactoring
│   └── test/         # Test-related docs
└── active/           # Currently active work
```

#### Naming Convention
Follows conventional commit types: `{type}/{ticket-id}/`
- Example: `ai-docs/completed/feat/spe-249/`

The qsweep script automatically:
- Scans documentation in: back-end-go/docs/, front-end-next/docs/, operator/docs/
- Extracts ticket IDs (SPE-XXX pattern)
- Determines document type from content
- Preserves file timestamps and structure

#### Claude Hooks Integration
The project includes automated hooks (configured in `.claude/settings.json`):

**Documentation Management:**
- **post-branch-create**: Automatically runs qsweep when creating a new branch
- **pre-commit**: Optionally sweeps docs before committing

**Code Formatting (post-edit):**
- **Go files**: Automatically runs `gofmt` and `goimports`
- **Markdown files**: Fixes EOF newlines and removes trailing whitespace
- **TypeScript/JavaScript**: Runs Prettier (if installed)
- **YAML files**: Runs Prettier (if installed)

To enable git hooks locally:
```bash
./.claude/setup-local.sh
```

---

<!-- SECTION: Important Reminders -->
# important-instruction-reminders
Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.

<!-- END OF FILE -->