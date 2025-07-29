## 📋 Migration Guide for Existing Users

This guide helps existing users navigate the restructured README and find
familiar content in its new location.

### What's Changed

#### 🎉 Major Improvements

- **50% more concise** - Focused on what users need most
- **Task-oriented structure** - Commands grouped by workflow
- **Better examples** - Real-world scenarios throughout
- **Clearer navigation** - Find anything in < 30 seconds
- **Progressive disclosure** - Simple for beginners, detailed docs linked

#### 📍 Content Relocation Map

| Old Section              | New Location                                   | Notes                                |
| ------------------------ | ---------------------------------------------- | ------------------------------------ |
| 🏗️ Repository Structure  | [Additional Resources](#-additional-resources) | Moved to reduce initial complexity   |
| 📜 Scripts               | [Commands Reference](#-commands-reference)     | Only user-facing scripts documented  |
| 🪝 Hooks                 | [Configuration](#️-configuration)              | Simplified, technical details linked |
| ⚡ Commands (verbose)    | [Commands Reference](#-commands-reference)     | Reorganized by workflow              |
| 🤖 Specialized Agents    | [Tips & Tricks](#-tips--tricks)                | Agent details in advanced section    |
| 📚 USER GUIDE            | Distributed across sections                    | Content integrated where relevant    |
| 🧪 Experimental Features | [Configuration](#️-configuration)              | Token-efficient mode highlighted     |
| 🔧 Templates             | [Additional Resources](#-additional-resources) | Advanced users only                  |
| ⚙️ Configuration         | [Configuration](#️-configuration)              | Essentials only, details linked      |
| 🚀 Setup                 | [Quickstart](#-quickstart)                     | Streamlined 5-minute setup           |

### Quick Reference for Power Users

#### Your Favorite Commands - New Locations

**Sprint & Planning:**

- `/qnew-enhanced` →
  [Planning & Sprint Management](#-planning--sprint-management)
- `/qplan-enhanced` →
  [Planning & Sprint Management](#-planning--sprint-management)
- `/qdirector-enhanced` →
  [Planning & Sprint Management](#-planning--sprint-management)

**Development:**

- `/qcode` → [Development](#-development)
- `/qtest` → [Development](#-development)

**Validation:**

- All `/qcheck*` commands → [Validation & Quality](#-validation--quality)
- `/qvalidate-framework` → [Validation & Quality](#-validation--quality)

**Utilities:**

- `/qgit` → [Utilities](#-utilities)
- Token efficiency → [Utilities](#-utilities)

#### Advanced Features Still Available

1. **Agent Architecture** - See
   [Tips & Tricks > Advanced Agent Usage](#advanced-agent-usage)
2. **Model Routing** - See
   [Commands Reference > Model Selection Strategy](#-model-selection-strategy)
3. **Parallel Execution** - See
   [Common Workflows > Parallel Task Execution](#parallel-task-execution)
4. **Hook Details** - Full details at `hooks/README.md`
5. **Template Files** - Still in `templates/` directory

### New Features to Explore

- **Real-World Examples** - 10 comprehensive scenarios in
  [Examples](#-real-world-examples)
- **Troubleshooting Guide** - Common issues and solutions
- **Performance Tips** - Token optimization strategies
- **Quick Workflows** - Copy-paste command sequences

### Preserving Your Workflow

Your existing workflows still work! The restructure only affects documentation,
not functionality:

```bash
# These commands work exactly the same:
/qnew-enhanced Create new feature
/qdirector-enhanced Execute sprint plan
/qvalidate-framework --mode deep

# Hooks still run automatically
# Scripts still in ~/.claude/scripts/
# Templates still in ~/.claude/templates/
```

### Getting Help

- **Can't find something?** Use browser search (Cmd+F / Ctrl+F)
- **Need old format?** Previous README backed up at `README.md.backup`
- **Missing details?** Check subdirectory README files for deep dives

### Feedback Welcome

This restructure prioritizes new user experience while maintaining power user
capabilities. If something critical is missing or harder to find, please
[open an issue](https://github.com/saintskeeper/claude-code-configs/issues).
