# CLAUDE.md Safety Guidelines

## Global Safety Rules for All Projects

### 1. Always Use MultiEdit for CLAUDE.md
```python
# NEVER use Edit for CLAUDE.md changes
# ALWAYS use MultiEdit with specific, unique strings

MultiEdit(
    file_path="CLAUDE.md",
    edits=[{
        "old_string": "<!-- SECTION: Specific Section -->\n## Section Title",
        "new_string": "<!-- SECTION: Specific Section -->\n## Section Title\n\n### New Content"
    }]
)
```

### 2. Section Markers Are Your Friends
When creating a new CLAUDE.md, always include:
```markdown
<!-- SECTION: Section Name -->
## Section Title
...content...

<!-- END OF FILE -->
```

### 3. Validate After Edits
Run validation after any CLAUDE.md changes:
```bash
~/.claude/scripts/validate-claude.sh
```

### 4. Template Location
Global template stored at: `~/.claude/templates/CLAUDE.template.md`

### 5. Recovery Process
If CLAUDE.md gets truncated:
```bash
# Check git history
git log --oneline CLAUDE.md

# Restore from template
cp ~/.claude/templates/CLAUDE.template.md CLAUDE.md

# Or restore from git
git checkout HEAD~1 CLAUDE.md
```

## Best Practices

1. **Small, Targeted Edits**: Make one section change at a time
2. **Unique Match Strings**: Include 2-3 lines of context
3. **Test First**: Use grep to verify your match string is unique
4. **Backup Before Major Changes**: `cp CLAUDE.md CLAUDE.md.bak`
5. **Use Section Comments**: Target edits to specific sections

## Common Patterns

### Adding a New Section
```python
MultiEdit(
    file_path="CLAUDE.md",
    edits=[{
        "old_string": "<!-- SECTION: Important Reminders -->",
        "new_string": "<!-- SECTION: New Section -->\n## New Section\n\nContent here...\n\n<!-- SECTION: Important Reminders -->"
    }]
)
```

### Updating Commands
```python
MultiEdit(
    file_path="CLAUDE.md",
    edits=[{
        "old_string": "```bash\ncd front-end-next\nnpm run dev",
        "new_string": "```bash\ncd front-end-next\nnpm run dev\nnpm run test  # Run tests"
    }]
)
```

### Adding to End of File
```python
MultiEdit(
    file_path="CLAUDE.md",
    edits=[{
        "old_string": "<!-- END OF FILE -->",
        "new_string": "## New Final Section\n\nContent...\n\n<!-- END OF FILE -->"
    }]
)
```

Remember: CLAUDE.md is critical for Claude Code's understanding of your project. Treat it with care!