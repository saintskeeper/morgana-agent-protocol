# Morgana Agent Protocol 🚀

**Supercharge Claude Code with parallel agent execution, 70% token savings, and
enterprise observability.**

Morgana transforms Claude Code into a multi-agent orchestration platform,
enabling parallel task execution, intelligent retry mechanisms, and
comprehensive monitoring - all while reducing API costs by up to 70%.

## Why Morgana?

| Standard Claude Code   | With Morgana Protocol            |
| ---------------------- | -------------------------------- |
| Sequential execution   | Parallel agent orchestration     |
| Single model usage     | Intelligent model routing        |
| Basic error handling   | Automatic retry with escalation  |
| Standard token usage   | 14-70% token reduction           |
| Limited observability  | Full OpenTelemetry tracing       |
| Manual task management | Automated workflow orchestration |

## 🎯 Quick Start (60 seconds)

```bash
# 1. Clone and setup
git clone https://github.com/saintskeeper/morgana-agent-protocol.git ~/.claude
cd ~/.claude

# 2. Build and install
make install-user  # No sudo required, installs to ~/.claude/bin

# 3. Test installation
make test  # Runs a test agent

# 4. Your first agent execution
morgana -- --agent code-implementer --prompt "Create a REST API endpoint"
```

## 📚 Essential Documentation

### Getting Started

- **[Installation & Setup](docs/GETTING_STARTED.md)** - Complete setup guide
  with prerequisites
- **[Quick Reference](docs/QUICK_REFERENCE.md)** - All make commands and common
  workflows
- **[Agent Selection Guide](docs/AGENTS.md)** - When and how to use each
  specialized agent

### Quick Links

- **[Examples](examples/)** - Real-world usage patterns
- **[Commands Reference](commands/)** - All slash commands
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions

## 🏗️ Architecture Overview

```
User Input → Morgana CLI → Go Orchestrator → Python Bridge → Claude Code
                ↓               ↓                              ↓
           Parallel Exec    Model Routing               Task Execution
                ↓               ↓                              ↓
           Monitoring      Token Optimization          Result Processing
```

## 🤖 Specialized Agents

| Agent                 | Purpose                      | Best For                       | Token Efficient         |
| --------------------- | ---------------------------- | ------------------------------ | ----------------------- |
| **code-implementer**  | Write production code        | Features, APIs, services       | ✅ Yes (14-70% savings) |
| **sprint-planner**    | Plan development sprints     | Task breakdown, prioritization | ❌ No                   |
| **test-specialist**   | Generate comprehensive tests | Unit, integration, E2E tests   | ✅ Yes                  |
| **validation-expert** | Review and validate code     | Security, quality, standards   | ✅ Yes                  |

[→ Full Agent Guide](docs/AGENTS.md)

## 💰 Token Savings

Enable token-efficient mode for 14-70% cost reduction:

```bash
~/.claude/scripts/token-efficient-config.sh enable
```

## 🔥 Key Features

- **Parallel Execution**: Run multiple agents simultaneously
- **Smart Retry**: Automatic escalation to more capable models
- **Real-time Monitoring**: Built-in TUI for live execution tracking
- **Full Observability**: OpenTelemetry tracing and metrics
- **Workflow Automation**: Chain commands for complex tasks

## 🚦 Common Workflows

### Feature Development

```bash
/morgana-plan "User authentication system"
/morgana-director  # Orchestrates implementation
/morgana-validate-all  # Comprehensive validation
```

### Quick Code Review

```bash
/morgana-check --file src/auth.service.ts
```

### Test Generation

```bash
/morgana-test generate unit --file src/user.model.ts
```

### Parallel Agent Execution

```bash
echo '[
  {"agent_type":"code-implementer","prompt":"Build feature"},
  {"agent_type":"test-specialist","prompt":"Write tests"}
]' | morgana --parallel
```

## 📊 Monitoring

```bash
# Start monitor daemon
make up

# View live execution
make attach

# Check system status
make status

# View logs
make logs

# Access web dashboards (optional)
# Jaeger: http://localhost:16686
# Grafana: http://localhost:3000
```

## 🛠️ Quick Command Reference

### Essential Make Commands

| Command             | Purpose                            |
| ------------------- | ---------------------------------- |
| `make help`         | Show all available commands        |
| `make install-user` | Install to ~/.claude/bin (no sudo) |
| `make up`           | Start monitor daemon               |
| `make attach`       | View live TUI                      |
| `make status`       | Check system status                |
| `make test`         | Run test agent                     |
| `make check`        | Full system health check           |
| `make clean`        | Clean and reset everything         |

[→ Full Command Reference](docs/QUICK_REFERENCE.md)

## Configuration

Basic configuration in `~/.claude/morgana.yaml`:

```yaml
agents:
  default_timeout: 2m
execution:
  max_concurrency: 5
telemetry:
  enabled: true
```

## 🎯 Best Practices

1. **Start with planning**: Use sprint-planner before implementation
2. **Parallel where possible**: Run independent agents simultaneously
3. **Enable token savings**: Use token-efficient mode for 14-70% reduction
4. **Monitor execution**: Use `make attach` to watch live progress
5. **Chain workflows**: Plan → Implement → Test → Validate

## 🔧 Troubleshooting

### Common Issues

**"morgana: command not found"**

```bash
export PATH=$HOME/.claude/bin:$PATH
```

**"Monitor not running"**

```bash
make up         # Start it
make status     # Verify
make attach     # Connect to view
```

**"Token limit exceeded"**

```bash
# Enable token-efficient mode
~/.claude/scripts/token-efficient-config.sh enable
```

[→ Full Troubleshooting Guide](docs/TROUBLESHOOTING.md)

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📝 License

MIT License - See [LICENSE](LICENSE) for details.

## 🔗 Resources

- [GitHub Repository](https://github.com/saintskeeper/morgana-agent-protocol)
- [Issue Tracker](https://github.com/saintskeeper/morgana-agent-protocol/issues)
- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)

---

<p align="center">
  <strong>Ready to supercharge your development?</strong><br>
  Start with <a href="docs/GETTING_STARTED.md">Getting Started</a> •
  <a href="docs/QUICK_REFERENCE.md">Quick Reference</a> •
  <a href="docs/AGENTS.md">Agent Guide</a>
</p>
