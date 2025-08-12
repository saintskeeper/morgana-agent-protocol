# Getting Started with Morgana Agent Protocol

## Prerequisites

Before installing Morgana, ensure you have:

- ✅ **Git** installed (`git --version`)
- ✅ **Claude Code** CLI installed and authenticated
- ✅ **Go** 1.21+ (`go version`)
- ✅ **Python** 3.8+ (`python3 --version`)
- ✅ **macOS** (current version optimized for macOS)

## Installation

### Step 1: Clone the Repository

```bash
git clone https://github.com/saintskeeper/morgana-agent-protocol.git ~/.claude
```

### Step 2: Build and Install

```bash
cd ~/.claude

# Option A: Install to ~/.claude/bin (no sudo required)
make install-user

# Option B: Install system-wide (requires sudo)
make install
```

This will:

- Build the Morgana binaries
- Install morgana and morgana-monitor
- Set up executable permissions
- Verify installation

### Step 3: Add to PATH and Source Functions

Add to your shell profile (`~/.zshrc` or `~/.bashrc`):

```bash
# Add Morgana to PATH
export PATH=$HOME/.claude/bin:$PATH

# Source agent adapter functions
source ~/.claude/scripts/agent-adapter-wrapper.sh
```

Then reload your shell:

```bash
source ~/.zshrc  # or source ~/.bashrc
```

### Step 4: Verify Installation

```bash
# Run system check
make check

# Test with a sample agent
make test

# Test parallel execution
make test-parallel
```

You should see:

- ✅ "Hello from Morgana!" message
- ✅ "Hook test complete!" confirmation
- ✅ Agent response with greeting

## Your First Agent Execution

### Simple Code Implementation

```bash
morgana -- --agent code-implementer --prompt "Create a fibonacci function in Python"
```

### Planning a Feature

```bash
morgana -- --agent sprint-planner --prompt "Plan a user registration feature"
```

### Generating Tests

```bash
morgana -- --agent test-specialist --prompt "Create tests for a login function"
```

### Code Review

```bash
morgana -- --agent validation-expert --prompt "Review this code for security issues: [paste code]"
```

## Understanding the Output

When you run an agent, you'll see:

```
[INFO] Starting agent: code-implementer
[INFO] Model: claude-3-7-sonnet (token-efficient)
[INFO] Timeout: 5m
...
[Agent Output]
...
[INFO] Execution complete in 2m34s
[INFO] Tokens used: 1,823 (saved 42% with token-efficient mode)
```

## Enabling Token Savings

Reduce API costs by 14-70%:

```bash
# Enable token-efficient mode
~/.claude/scripts/token-efficient-config.sh enable

# Verify status
~/.claude/scripts/token-efficient-config.sh status
```

## Basic Configuration

Create `~/.claude/morgana.yaml`:

```yaml
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m
  timeouts:
    code-implementer: 5m
    sprint-planner: 3m
    test-specialist: 4m
    validation-expert: 2m

execution:
  max_concurrency: 5
  default_mode: sequential

telemetry:
  enabled: true
  exporter: stdout

token_efficient:
  enabled: true
  models:
    - claude-3-7-sonnet
```

## Common First Tasks

### 1. Plan and Implement a Feature

```bash
# Step 1: Plan
morgana -- --agent sprint-planner --prompt "Plan a todo list API"

# Step 2: Implement
morgana -- --agent code-implementer --prompt "Implement the first task from the sprint plan"

# Step 3: Test
morgana -- --agent test-specialist --prompt "Create tests for the todo API"

# Step 4: Validate
morgana -- --agent validation-expert --prompt "Review the todo API implementation"
```

### 2. Parallel Execution

Run multiple agents at once:

```bash
echo '[
  {"agent_type": "code-implementer", "prompt": "Create user model"},
  {"agent_type": "test-specialist", "prompt": "Create user model tests"}
]' | morgana --parallel
```

### 3. Using Slash Commands

In Claude Code, you can use slash commands:

```bash
/morgana-plan Create an authentication system
/morgana-code Implement JWT token generation
/morgana-test Generate tests for JWT functions
/morgana-check Validate the implementation
```

## Monitoring Your Agents

### Start the TUI Monitor

```bash
# Start monitor daemon
make up

# Check status
make status

# In another terminal, run agents
morgana -- --agent code-implementer --prompt "Create API"

# Attach to see live execution
make attach

# View logs
make logs

# Stop monitor when done
make down
```

### View Traces (Optional)

For detailed tracing with Jaeger:

```bash
cd ~/.claude/monitoring
docker-compose up -d

# Access at http://localhost:16686
```

## Troubleshooting Common Issues

### "Python not found"

```bash
# Set Python path explicitly
export PYTHON_PATH=$(which python3)
```

### "morgana: command not found"

```bash
# Add to PATH
export PATH="$HOME/.claude/bin:$PATH"

# Or use full path
~/.claude/bin/morgana -- --agent code-implementer --prompt "Test"
```

### "Agent timeout"

```bash
# Increase timeout for complex tasks
morgana -- --agent code-implementer --prompt "Complex task" --timeout 10m
```

### "Token limit exceeded"

```bash
# Enable token-efficient mode
~/.claude/scripts/token-efficient-config.sh enable

# Or break into smaller tasks
morgana -- --agent sprint-planner --prompt "Break down: [large task]"
```

## Next Steps

Now that you have Morgana running:

1. 📖 Read [Understanding Agents](AGENTS.md) to learn about each specialist
2. 🚀 Explore [Parallel Execution](PARALLEL.md) for complex workflows
3. 💰 Configure [Token Optimization](TOKEN_OPTIMIZATION.md) to reduce costs
4. 📊 Set up [Monitoring](MONITORING.md) for observability
5. 🎯 Try [Example Workflows](../examples/README.md) for real-world patterns

## Getting Help

- Check [Troubleshooting Guide](TROUBLESHOOTING.md) for common issues
- Review [FAQ](FAQ.md) for frequently asked questions
- Report issues on
  [GitHub](https://github.com/saintskeeper/morgana-agent-protocol/issues)
- Join discussions in
  [GitHub Discussions](https://github.com/saintskeeper/morgana-agent-protocol/discussions)

## Quick Command Reference

### Make Commands (Recommended)

| Command              | Purpose                            |
| -------------------- | ---------------------------------- |
| `make help`          | Show all available commands        |
| `make install-user`  | Install to ~/.claude/bin (no sudo) |
| `make install`       | Install system-wide                |
| `make build`         | Build binaries only                |
| `make up`            | Start monitor daemon               |
| `make down`          | Stop monitor daemon                |
| `make status`        | Check monitor status               |
| `make attach`        | View live TUI                      |
| `make logs`          | Show monitor logs                  |
| `make test`          | Run test agent                     |
| `make test-parallel` | Test parallel execution            |
| `make check`         | Full system check                  |
| `make clean`         | Clean up everything                |
| `make restart`       | Restart monitor                    |

### Morgana Commands

| Command                                       | Purpose                          |
| --------------------------------------------- | -------------------------------- |
| `morgana -- --agent [type] --prompt "[task]"` | Run single agent                 |
| `morgana --parallel`                          | Run multiple agents (JSON input) |
| `morgana --list-models`                       | Show available models            |
| `morgana --config [file]`                     | Use custom configuration         |

---

**Ready for more?** Continue to [Understanding Agents](AGENTS.md) →
