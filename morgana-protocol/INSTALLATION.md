# Morgana Protocol Installation Guide

This guide covers different ways to install the Morgana Protocol CLI with TUI
support.

## Quick Installation (Recommended)

### Option 1: System Installation (`/usr/local/bin`)

Installs system-wide, accessible to all users. **Requires sudo**.

```bash
# Build from source and install to /usr/local/bin
sudo ./scripts/setup.sh
```

### Option 2: User Installation (`~/.local/bin`)

Installs for current user only. **No sudo required**.

```bash
# Build from source and install to ~/.local/bin
./scripts/setup-user.sh
```

### Option 3: Pre-built Binary Installation

If you have pre-built binaries in the `dist/` directory.

```bash
# Install pre-built binary to ~/.claude/bin
./scripts/install.sh
```

## Manual Installation

### Prerequisites

- **Go 1.21+** - Install from [golang.org](https://golang.org/dl/)
- **Git** - For cloning the repository

### Build from Source

1. **Clone the repository:**

   ```bash
   git clone https://github.com/saintskeeper/claude-code-configs.git
   cd claude-code-configs/morgana-protocol
   ```

2. **Build the CLI:**

   ```bash
   go mod tidy
   go build -o morgana ./cmd/morgana
   ```

3. **Install manually:**

   ```bash
   # System-wide (requires sudo)
   sudo cp morgana /usr/local/bin/

   # User installation
   mkdir -p ~/.local/bin
   cp morgana ~/.local/bin/
   ```

4. **Add to PATH if needed:**
   ```bash
   # For ~/.local/bin installation
   echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
   source ~/.bashrc
   ```

## Installation Options Comparison

| Method                    | Location         | Sudo Required | Scope        | Best For                   |
| ------------------------- | ---------------- | ------------- | ------------ | -------------------------- |
| `./scripts/setup.sh`      | `/usr/local/bin` | Yes           | All users    | Production, shared systems |
| `./scripts/setup-user.sh` | `~/.local/bin`   | No            | Current user | Development, personal use  |
| `./scripts/install.sh`    | `~/.claude/bin`  | No            | Current user | Pre-built binaries         |
| Manual                    | Custom           | Varies        | Custom       | Advanced users             |

## Verification

After installation, verify it works:

```bash
# Check version
morgana --version

# Test TUI mode
morgana --tui --version

# Quick TUI demo with built-in tasks
morgana --tui --tui-mode dev
```

Expected output:

```
Morgana Protocol v1.0.0-abc123
🧙‍♂️ Morgana Protocol TUI started in dev mode. Press 'q' in TUI or Ctrl+C to quit.
```

## Configuration

### Default Configuration

The CLI looks for configuration in these locations (in order):

1. `--config` flag: `morgana --config /path/to/config.yaml`
2. `morgana.yaml` in current directory
3. `~/.claude/morgana.yaml`
4. Built-in defaults

### Sample Configuration

Create `morgana.yaml` with TUI settings:

```yaml
# TUI configuration
tui:
  enabled: true
  performance:
    refresh_rate: 16ms # 60 FPS
    max_log_lines: 10000
  visual:
    show_debug_info: false
    show_timestamps: true
  features:
    enable_filtering: true
    enable_search: true

# Agent configuration
agents:
  prompt_dir: ~/.claude/agents
  default_timeout: 2m

# Execution configuration
execution:
  max_concurrency: 5
  default_mode: sequential

# Telemetry configuration
telemetry:
  enabled: true
  exporter: stdout
  service_name: morgana-protocol
```

## Troubleshooting Installation

### Common Issues

#### 1. Go not found

```bash
# Install Go
curl -L https://go.dev/dl/go1.21.5.linux-amd64.tar.gz | sudo tar -C /usr/local -xz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. Permission denied (Linux/macOS)

```bash
# Use user installation instead
./scripts/setup-user.sh

# Or fix permissions for system installation
sudo chown -R $(whoami) /usr/local/bin
```

#### 3. Binary not in PATH

```bash
# Check current PATH
echo $PATH

# Add installation directory to PATH
# For system installation:
export PATH="$PATH:/usr/local/bin"

# For user installation:
export PATH="$PATH:$HOME/.local/bin"

# Make permanent:
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
```

#### 4. TUI not working

```bash
# Check terminal support
echo $TERM
tput colors

# Try different terminal
# Use iTerm2, gnome-terminal, or Terminal.app

# Test basic functionality
morgana --version
morgana --help
```

#### 5. Build failures

```bash
# Clean and retry
go clean -cache
go mod tidy
go build -v ./cmd/morgana
```

### Platform-Specific Notes

#### macOS

- Use Homebrew if available: `brew install go`
- Terminal.app and iTerm2 both support TUI
- May need to allow in System Preferences > Security

#### Linux

- Install via package manager: `sudo apt install golang-go` (Ubuntu)
- Most modern terminals support TUI
- May need to install `build-essential`

#### Windows

- Use Git Bash, PowerShell, or WSL
- Windows Terminal recommended for best TUI support
- May require Visual Studio Build Tools

## Uninstallation

### Remove Binary

```bash
# System installation
sudo rm /usr/local/bin/morgana

# User installation
rm ~/.local/bin/morgana
# or
rm ~/.claude/bin/morgana
```

### Remove Configuration

```bash
# Remove user configuration
rm -rf ~/.claude/morgana.yaml

# Remove project configuration (if not needed)
rm morgana.yaml
```

## Getting Help

- **Documentation**: See `TUI_USER_GUIDE.md` for TUI usage
- **Integration**: See `MONITORING_INTEGRATION.md` for technical details
- **Issues**: Report problems with installation details:
  - Operating system and version
  - Go version (`go version`)
  - Terminal type and version
  - Full error messages

## Next Steps

After installation:

1. **Quick Start**: `morgana --tui`
2. **Configure**: Edit `morgana.yaml`
3. **Learn TUI**: Read `TUI_USER_GUIDE.md`
4. **Set up agents**: Configure your `~/.claude/agents` directory

Happy orchestrating! 🧙‍♂️
