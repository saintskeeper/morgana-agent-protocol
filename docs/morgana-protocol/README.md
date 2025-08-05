# Morgana Protocol: Agent Orchestration System

## Overview

The Morgana Protocol is a Go-based agent orchestration system designed to
overcome the Task tool limitation in Claude Code where only "general-purpose"
subagent types are recognized. It provides true parallel execution through
goroutines and cross-platform deployment capabilities.

## Architecture

```
morgana-protocol/
├── cmd/
│   └── morgana/
│       └── main.go          # Entry point
├── internal/
│   ├── adapter/
│   │   ├── adapter.go       # Core adapter logic
│   │   └── adapter_test.go
│   ├── orchestrator/
│   │   ├── parallel.go      # Goroutine-based orchestration
│   │   └── parallel_test.go
│   └── prompt/
│       ├── loader.go        # Agent prompt loader
│       └── loader_test.go
├── pkg/
│   └── task/
│       └── client.go        # Task tool interface
├── scripts/
│   ├── build.sh
│   └── install.sh
├── go.mod
├── go.sum
└── Makefile
```

## Key Features

1. **Agent Type Translation**: Maps specialized agent types to general-purpose
2. **True Parallelism**: Goroutines enable concurrent agent execution
3. **Cross-Platform**: Compiles to native binaries for macOS, Linux, Windows
4. **Zero Dependencies**: Single binary deployment
5. **JSON API**: Compatible with existing Claude Code workflows

## Quick Start

```bash
# Install Morgana Protocol
~/.claude/docs/morgana-protocol/scripts/install.sh

# Run parallel agents
morgana run --parallel \
  --agent code-implementer --prompt "implement auth service" \
  --agent test-specialist --prompt "create auth tests"
```

## Implementation Status

- [ ] Core adapter infrastructure
- [ ] Parallel orchestration engine
- [ ] Agent prompt loading system
- [ ] Multi-platform build pipeline
- [ ] Integration with Claude Code

## Why "Morgana"?

Named after Morgana le Fay, the powerful enchantress from Arthurian legend who
could transform and adapt - much like how this protocol transforms specialized
agent requests into a format the Task tool can understand.
