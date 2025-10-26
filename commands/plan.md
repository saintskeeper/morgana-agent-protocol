# Plan
Create a detailed implementation plan for a feature and save it to the issues directory.

---
description: Create feature implementation plan with file changes
allowed-tools:
  - Read
  - Write
  - Grep
  - Glob
  - Bash(git:*)
---

Create a detailed implementation plan for: $ARGUMENTS

Follow these steps:

1. **Analyze the Feature Request**
   - Understand the requirements and scope
   - Identify affected components and modules

2. **Search the Codebase**
   - Use Grep and Glob to find relevant files
   - Read key files to understand current implementation
   - Identify patterns and conventions

3. **Identify Files to Change**
   For each file that needs modification or creation:
   - Specify the exact file path
   - Indicate if it's NEW or EDIT
   - Specify line ranges if editing (e.g., lines 25-27, 100-150)
   - Brief description of changes needed

4. **Create Plan Document**
   Generate a markdown file at `issues/feature/<feature-name>.md` with:

   ```markdown
   # Feature: <Feature Name>

   ## Overview
   Brief description of the feature and its purpose.

   ## Requirements
   - Functional requirements
   - Technical requirements
   - Dependencies

   ## Architecture Changes
   - Component interactions
   - Data flow
   - API changes

   ## Implementation Plan

   ### Phase 1: [Phase Name]
   **Files to Modify:**

   #### 📝 EDIT: `path/to/file.go` (lines 25-27, 100-150)
   - Description of changes
   - Why these changes are needed

   #### ✨ NEW: `path/to/new_file.go`
   - Purpose of new file
   - Key functionality

   ### Phase 2: [Phase Name]
   ...

   ## Testing Strategy
   - Unit tests needed
   - Integration tests
   - Manual testing steps

   ## Rollout Plan
   - Deployment steps
   - Monitoring requirements
   - Rollback procedure

   ## Timeline Estimate
   - Phase breakdown with estimates

   ## Risk Assessment
   - Potential issues
   - Mitigation strategies
   ```

5. **Validate Plan**
   - Ensure all necessary files are identified
   - Check for missing dependencies
   - Verify line ranges are accurate

6. **Report**
   Confirm plan creation and provide path to the generated file.

Format the feature name as lowercase-with-hyphens for the filename.
