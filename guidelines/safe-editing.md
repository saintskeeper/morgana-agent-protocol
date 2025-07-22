# Safe Editing Guidelines for Claude Code

## Preventing File Truncation

### 1. Use MultiEdit for Large Files
When editing CLAUDE.md or other critical files, use MultiEdit instead of Edit:

```python
# GOOD - Multiple targeted edits
MultiEdit(
    file_path="CLAUDE.md",
    edits=[
        {
            "old_string": "### 4 — Documentation Management",
            "new_string": "### 4 — Documentation Management\n\n#### New Section"
        },
        {
            "old_string": "# important-instruction-reminders",
            "new_string": "### 5 — Another Section\n\n# important-instruction-reminders"
        }
    ]
)
```

### 2. Use Unique Markers
Add unique section markers to CLAUDE.md:

```markdown
<!-- SECTION: Linear Integration -->
## Linear Integration
...

<!-- SECTION: Architecture Overview -->
## Architecture Overview
...

<!-- SECTION: Claude Guidelines -->
# Claude Code Guidelines
...

<!-- SECTION: Important Reminders -->
# important-instruction-reminders
...
<!-- END OF FILE -->
```

### 3. Create a Template
Keep a template version for reference:

```bash
cp CLAUDE.md .claude/templates/CLAUDE.template.md
```

### 4. Use Append Pattern
Instead of editing, append new sections:

```python
# Read current content
current_content = Read("CLAUDE.md")

# Find the insertion point
insertion_point = "<!-- SECTION: Important Reminders -->"

# Split and insert
parts = current_content.split(insertion_point)
new_content = parts[0] + new_section + "\n\n" + insertion_point + parts[1]

# Write back
Write("CLAUDE.md", new_content)
```

### 5. Create Edit Helpers
Add helper scripts for common edits:

```bash
.claude/scripts/
├── add-section.sh      # Add new section to CLAUDE.md
├── update-commands.sh  # Update command references
└── validate-claude.sh  # Check CLAUDE.md integrity
```

## Best Practices

1. **Always use specific, unique strings** for old_string
2. **Include context** - match at least 2-3 lines
3. **Test with grep first** to ensure unique match
4. **Use MultiEdit** for multiple changes
5. **Keep backups** before major edits

## Example Safe Edit

```python
# First, verify the string exists and is unique
Bash("grep -n 'unique string to match' CLAUDE.md")

# Then use MultiEdit with context
MultiEdit(
    file_path="CLAUDE.md",
    edits=[{
        "old_string": "### Specific Section Title\n\nSome content here",
        "new_string": "### Specific Section Title\n\nSome content here\n\n#### New Subsection"
    }]
)
```