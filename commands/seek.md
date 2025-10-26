---
description: Find relevant files and relationships with line numbers
allowed-tools:
  - Glob
  - Grep
  - Read
---

# Seek

Intelligently find relevant files and their relationships across the Cromie codebase, returning results with file paths and specific line number ranges. Uses OS-agnostic search tools that work on Windows, macOS, and Linux.

## Purpose

This command analyzes the codebase to find files related to a specific topic, feature, or component. It identifies code relationships, dependencies, and provides precise line number references for relevant sections.

## Usage

```
/seek <topic|feature|component>
```

Examples:
- `/seek slack commands`
- `/seek time entry`
- `/seek authentication`
- `/seek database models`

## OS-Agnostic Search Strategy

Using Claude's built-in Glob and Grep tools for cross-platform compatibility:

### 1. Initial Discovery

Use Glob to find relevant files by pattern:
```
Glob pattern: **/*time*.go
Glob pattern: app/server/internal/**/*.go
Glob pattern: app/client/**/*.{ts,tsx}
```

Use Grep for content search:
```
Grep pattern: "time entry" (with output_mode: "files_with_matches")
Grep pattern: "type.*TimeEntry|struct.*TimeEntry" (with output_mode: "content", -n: true)
Grep pattern: "import.*slack" (with -n: true, -B: 2, -A: 2)
```

### 2. Dependency Mapping

Find imports and dependencies:
```
Grep pattern: "import.*{SEARCH_TERM}" (type: "go")
Grep pattern: "from.*{SEARCH_TERM}" (type: "ts")
Grep pattern: "type.*{SEARCH_TERM}|interface.*{SEARCH_TERM}" (glob: "**/*.{go,ts}")
```

### 3. Relationship Analysis

Analyze the following relationship patterns:

#### Backend Services (Go/Gin)
```
# Handler → Service → Model pattern
Glob: app/server/internal/api/handlers/*.go
Glob: app/server/internal/models/*.go
Glob: app/server/internal/middleware/*.go
Glob: app/server/internal/slack/*.go

# Then use Grep with -n flag to get line numbers:
Grep: "func|type|struct" (with -n: true, output_mode: "content")
```

#### Frontend Components (Next.js/React)
```
# Page → Component → Service pattern
Glob: app/client/app/**/*.{ts,tsx}
Glob: app/client/components/**/*.{ts,tsx}
Glob: app/client/lib/*.ts

# Get specific line numbers for key sections:
Grep: "export function|export const|export default" (with -n: true)
```

#### Database Models (GORM)
```
# Model definitions and migrations
Glob: app/server/internal/models/*.go
Glob: migrations/**/*.sql

# Find model definitions and relationships:
Grep: "type.*struct|gorm:\"" (with -n: true)
```

## Search Execution Plan

### Step 1: Find Primary Files
```
# Use Glob to find files by name pattern
Glob pattern: **/*{search_term}*.{go,ts,tsx}

# Use Grep to find files containing the search term
Grep pattern: "{search_term}"
  output_mode: "files_with_matches"
  head_limit: 50
```

### Step 2: Extract Line Numbers for Key Sections
```
# For each relevant file, get specific line numbers
Grep pattern: "^(type|func|const|var).*{search_term}"
  -n: true
  output_mode: "content"

# Find method definitions
Grep pattern: "func.*{search_term}|{search_term}.*func"
  -n: true
  output_mode: "content"
  -i: true
```

### Step 3: Find Dependencies
```
# Find import statements
Grep pattern: "import.*{search_term}|from.*{search_term}"
  -n: true
  -B: 1
  -A: 1
  output_mode: "content"
```

### Step 4: Find Related Tests
```
# Find test files
Glob pattern: **/*_test.go
Glob pattern: **/*.spec.{ts,tsx}
Glob pattern: **/*.test.{ts,tsx}

# Then search within test files
Grep pattern: "{search_term}"
  glob: "*_test.go"
  output_mode: "files_with_matches"
```

## Output Format

Return results organized by service/layer with specific line references:

```markdown
## Search Results for: [TOPIC]

### Primary Implementation Files

**Backend API (Go/Gin)**
- `app/server/internal/api/handlers/slack_events.go` [45-120, 200-250] - Event handlers
- `app/server/internal/models/time_entry.go` [1-100] - TimeEntry model definition
- `app/server/internal/slack/client.go` [50-200] - Slack API client

**Frontend (Next.js/React)**
- `app/client/app/(dashboard)/reports/page.tsx` [25-150] - Reports page
- `app/client/components/ui/table.tsx` [1-200] - Table component
- `app/client/lib/utils.ts` [10-50] - Utility functions

**Database Models**
- `app/server/internal/models/user.go` [1-75] - User model
- `app/server/internal/models/organization.go` [1-60] - Organization model
- `app/server/pkg/database/database.go` [20-100] - Database setup

### Middleware & Authentication
- `app/server/internal/middleware/clerk_auth.go` [1-120] - Clerk JWT validation
- `app/server/internal/middleware/slack_verify.go` [1-80] - Slack signature verification
- `app/client/middleware.ts` [1-50] - Clerk route protection

### Related Configuration Files
- `Dockerfile` [1-30] - Container configuration
- `infrastructure/scripts/deploy-prod.sh` [1-100] - Production deployment
- `.env.example` [1-50] - Environment variables

### Test Coverage
- `app/server/internal/api/handlers/admin_test.go` [1-200] - Handler unit tests
- `app/server/internal/models/time_entry_test.go` [1-150] - Model tests

### Documentation
- `README.md` [15-45] - Feature overview
- `roadmap.md` [1-200] - Development roadmap
- `ai_docs/slack_golang_sdk.md` [1-500] - Slack SDK reference
```

## Domain-Specific Search Patterns

### Time Entry Management
```
Glob: app/server/internal/models/time_entry.go
Glob: app/server/internal/api/handlers/**/*time*.go
Glob: app/client/app/(dashboard)/reports/*.{ts,tsx}

Grep: "TimeEntry|time_entry|timeEntry" (with -i: true, -n: true)
```

### Slack Integration
```
Glob: app/server/internal/slack/*.go
Glob: app/server/internal/api/handlers/slack_*.go

Grep: "slack|modal|command|event|interaction" (with -i: true, -n: true)
```

### User & Organization Management
```
Glob: app/server/internal/models/user*.go
Glob: app/server/internal/models/organization.go
Glob: app/client/app/(dashboard)/users/*.{ts,tsx}

Grep: "User|Organization|user|organization" (with -n: true)
```

### Authentication & Security
```
Glob: app/server/internal/middleware/*auth*.go
Glob: app/server/internal/middleware/*verify*.go
Glob: app/client/middleware.ts
Glob: app/client/app/(auth)/**/*.tsx

Grep: "auth|jwt|token|clerk|signature" (with -i: true, -n: true)
```

### Client Management
```
Glob: app/server/internal/models/client.go
Glob: app/server/internal/models/user_client.go
Glob: app/client/app/(dashboard)/clients/*.{ts,tsx}

Grep: "Client|client_id|UserClient" (with -n: true)
```

### Scheduled Messages
```
Glob: app/server/internal/models/scheduled_message.go
Glob: app/server/internal/scheduler/**/*.go

Grep: "schedule|cron|goroutine|ticker" (with -i: true, -n: true)
```

## Smart Search Functions

### Find by Architectural Layer

**Handlers (REST endpoints)**
```
Glob: app/server/internal/api/handlers/*.go
Grep: "{search_term}" (glob: "*.go", path: app/server/internal/api/handlers, -n: true)
```

**Models (Database schemas)**
```
Glob: app/server/internal/models/*.go
Grep: "{search_term}" (glob: "*.go", path: app/server/internal/models, -n: true)
```

**Middleware (Request processing)**
```
Glob: app/server/internal/middleware/*.go
Grep: "{search_term}" (glob: "*.go", path: app/server/internal/middleware, -n: true)
```

**Pages (UI routes)**
```
Glob: app/client/app/**/*.{ts,tsx}
Grep: "{search_term}" (glob: "*.{ts,tsx}", path: app/client/app, -n: true)
```

**Components (Reusable UI)**
```
Glob: app/client/components/**/*.{ts,tsx}
Grep: "{search_term}" (glob: "*.{ts,tsx}", path: app/client/components, -n: true)
```

### Find Cross-Service Dependencies

**API routes definition**
```
Grep: "router\\.(GET|POST|PUT|DELETE)|r\\.Group"
  path: app/server/internal/api
  type: "go"
  -n: true
```

**Frontend API calls**
```
Grep: "fetch\\(|axios\\.|useQuery|useMutation"
  path: app/client
  type: "ts"
  -n: true
```

**Database queries**
```
Grep: "db\\.|gorm\\.|First|Find|Create|Update|Delete|Where"
  path: app/server
  type: "go"
  -n: true
```

### Find Configuration and Deployment

**Environment variables**
```
Grep: "os\\.Getenv|process\\.env\\.|NEXT_PUBLIC_"
  -n: true
  output_mode: "content"
```

**Docker configurations**
```
Glob: Dockerfile*
Glob: docker-compose*.{yml,yaml}
Grep: "{search_term}" (glob: "Dockerfile*", -n: true)
```

**Railway deployment**
```
Glob: infrastructure/scripts/*.sh
Grep: "railway|deploy" (glob: "*.sh", path: infrastructure/scripts, -n: true)
```

## Integration Points

### Clerk Authentication
```
Glob: app/server/internal/middleware/clerk_auth.go
Glob: app/client/app/layout.tsx
Glob: app/client/middleware.ts
Grep: "clerk|ClerkProvider|clerkMiddleware" (with -i: true, -n: true)
```

### Slack API Integration
```
Glob: app/server/internal/slack/client.go
Glob: app/server/internal/api/handlers/slack_*.go
Grep: "slack.*api|PostMessage|OpenView|UpdateView" (with -n: true)
```

### Database Layer (GORM)
```
Glob: app/server/pkg/database/database.go
Glob: app/server/internal/models/*.go
Grep: "gorm\\.Model|AutoMigrate|db\\..*" (with -n: true)
```

### shadcn/ui Components
```
Glob: app/client/components/ui/*.tsx
Grep: "cn\\(|cva\\(|class.*Variance" (with -n: true)
```

## Advanced Usage

### Trace Request Flow
For a complete request flow trace, use these searches in sequence:

1. **Find REST endpoint**:
   ```
   Glob: app/server/internal/api/handlers/*.go
   Grep: "r\\.(GET|POST|PUT|DELETE).*{endpoint}" (-n: true)
   ```

2. **Trace to handler function**:
   ```
   Grep: "func.*{handlerName}" (path: app/server/internal/api/handlers, -n: true)
   ```

3. **Follow to model operations**:
   ```
   Grep: "db\\..*{modelName}" (path: app/server, -n: true)
   ```

4. **Check frontend call**:
   ```
   Grep: "{endpoint}|fetch.*{endpoint}" (path: app/client, -n: true)
   ```

5. **Find UI component**:
   ```
   Glob: app/client/app/**/*.{ts,tsx}
   Grep: "useEffect|useState|{endpointName}" (-n: true)
   ```

### Find Feature Implementation
To understand a complete feature:

1. **Check specifications**:
   ```
   Glob: issues/feature/*{feature}*.md
   ```

2. **Find backend implementation**:
   ```
   Glob: app/server/**/*{feature}*.go
   ```

3. **Locate frontend components**:
   ```
   Glob: app/client/**/*{feature}*.{ts,tsx}
   ```

4. **Review tests**:
   ```
   Glob: **/*{feature}*_test.go
   Glob: **/*{feature}*.spec.{ts,tsx}
   ```

5. **Check documentation**:
   ```
   Glob: ai_docs/*{feature}*.md
   Glob: **/*README*.md
   ```

## Example Searches

### Example 1: Find Time Entry Implementation
```
# Find all time entry related files
Glob: **/*time*entry*.*

# Find TimeEntry struct/type definitions with line numbers
Grep: "type.*TimeEntry.*struct" (-n: true, output_mode: "content")

# Find time entry handler methods
Grep: "func.*TimeEntry|TimeEntry.*func" (-i: true, -n: true)

# Find time entry API endpoints
Grep: "router\\.(GET|POST|PUT|DELETE).*time" (-i: true, -n: true)
```

### Example 2: Trace Slack Command Flow
```
# Find slash command handlers
Glob: app/server/internal/api/handlers/slack_*.go

# Find command processing logic
Grep: "command|slash|/cromie" (-i: true, -n: true)

# Find modal interactions
Grep: "modal|view|OpenView|UpdateView" (-n: true)

# Find event handlers
Grep: "HandleSlackEvent|socketmode|EventsAPI" (-n: true)
```

### Example 3: Find Authentication Implementation
```
# Find Clerk integration
Grep: "clerk|jwt|validateToken" (-i: true, -n: true)

# Find middleware usage
Grep: "Use\\(.*Auth|AuthMiddleware" (-n: true)

# Find protected routes
Grep: "clerkMiddleware|authMiddleware|requireAuth" (path: app/client, -n: true)

# Find token validation
Grep: "jwks|ValidateToken|VerifyToken" (-n: true)
```

## Instructions for Execution

When using `/seek <search_term>`:

1. **Initial File Discovery**:
   - Run Glob patterns to find files by name
   - Run Grep with `output_mode: "files_with_matches"` for content search
   - Combine results and deduplicate

2. **Extract Detailed Information**:
   - For each relevant file, run Grep with `-n: true` to get line numbers
   - Search for function definitions, type definitions, imports
   - Use `-B` and `-A` flags for context

3. **Analyze Relationships**:
   - Find imports and dependencies between files
   - Trace data flow from handlers to models to database
   - Map frontend to backend API connections

4. **Organize Results**:
   - Group by architectural layer (Backend, Frontend, Database, etc.)
   - Include file paths with line number ranges
   - Provide brief descriptions of each file's purpose

5. **Provide Context**:
   - Explain how files relate to each other
   - Identify entry points and data flows
   - Highlight configuration and deployment files

## Notes

- Line numbers are extracted using the `-n` flag in Grep
- Use `-B` and `-A` flags to get context around matches (before/after lines)
- The `head_limit` parameter helps manage large result sets
- Use `glob` parameter to filter by file type
- Use `type` parameter for language-specific searches (e.g., "go", "ts", "js")
- Combine multiple search patterns with regex OR operator: `pattern1|pattern2|pattern3`
- For case-insensitive searches, use the `-i` flag
- Always provide line number ranges in the format `[start-end]` or `[line]`

## Integration with /bug and /plan Commands

After running `/seek`, use the results to populate:

**For `/bug` command**:
- **Files Involved** section with line numbers from seek results
- **Root Cause Analysis** informed by code relationships found
- **Files to Modify** with specific line ranges

**For `/plan` command**:
- **Files to Change** section with precise paths and line ranges
- **Architecture Changes** based on dependency analysis
- **Implementation Plan** phases organized by architectural layer

Example workflow:
```bash
# 1. Find all relevant files
/seek slack commands

# 2. Create bug report with specific file references
/bug slack-modal-timeout The modal doesn't open when user runs /cromie log. Based on seek results, issue is in app/server/internal/api/handlers/slack_events.go:150-180

# OR create feature plan with file changes
/plan Add edit functionality to /cromie command. Modify files found in seek: app/server/internal/slack/client.go, app/server/internal/api/handlers/slack_events.go
```
