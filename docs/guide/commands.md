# Commands

Plexr provides several commands to manage and execute your plans. This guide covers all available commands and their options.

## Overview

```bash
plexr [command] [flags]
```

Available commands:
- `execute` - Run an execution plan
- `validate` - Validate a plan without executing
- `status` - Show current execution status
- `reset` - Reset execution state
- `completion` - Generate shell completions
- `help` - Get help on any command
- `version` - Show version information

## execute

The main command to run execution plans.

### Basic Usage

```bash
plexr execute [plan-file] [flags]
```

### Examples

```bash
# Execute a plan
plexr execute setup.yml

# Dry run to see what would happen
plexr execute setup.yml --dry-run

# Auto-confirm all prompts
plexr execute setup.yml --auto

# Use specific platform
plexr execute setup.yml --platform=linux
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dry-run` | `-n` | Show what would be executed without running | `false` |
| `--auto` | `-y` | Auto-confirm all prompts | `false` |
| `--platform` | `-p` | Override platform detection | auto-detect |
| `--state-file` | `-s` | Custom state file location | `.plexr_state.json` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--force` | `-f` | Force re-execution of completed steps | `false` |

### Execution Flow

1. Load and validate the plan
2. Check current state
3. Resolve dependencies
4. Execute steps in order
5. Update state after each step
6. Handle failures gracefully

### Resume on Failure

If execution fails, simply run the command again to resume:

```bash
# First run fails at step 3
plexr execute setup.yml
# Error: Step 3 failed

# Fix the issue, then resume
plexr execute setup.yml
# Resuming from step 3...
```

## validate

Validate a plan file without executing it.

### Usage

```bash
plexr validate [plan-file] [flags]
```

### Examples

```bash
# Validate syntax and structure
plexr validate setup.yml

# Validate with verbose output
plexr validate setup.yml --verbose
```

### What It Checks

- YAML syntax
- Required fields presence
- Field types and values
- Circular dependencies
- File existence (with `--check-files`)
- Executor availability

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--check-files` | `-c` | Verify referenced files exist | `false` |
| `--verbose` | `-v` | Show detailed validation info | `false` |

## status

Show the current execution state of a plan.

### Usage

```bash
plexr status [plan-file] [flags]
```

### Examples

```bash
# Show current status
plexr status setup.yml

# Show detailed status
plexr status setup.yml --verbose
```

### Output

```
Plan: Development Environment Setup (v1.0.0)
State: In Progress

Progress: 3/5 steps completed (60%)

✓ install_tools      - Install development tools
✓ create_directories - Create project directories  
✓ setup_database     - Initialize database
→ configure_app      - Configure application (current)
○ run_tests          - Run verification tests

Last updated: 2023-12-15 10:30:45
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--state-file` | `-s` | Custom state file location | `.plexr_state.json` |
| `--verbose` | `-v` | Show detailed status | `false` |
| `--json` | `-j` | Output in JSON format | `false` |

## reset

Reset the execution state, allowing you to start fresh.

### Usage

```bash
plexr reset [plan-file] [flags]
```

### Examples

```bash
# Reset all state
plexr reset setup.yml

# Reset without confirmation
plexr reset setup.yml --force

# Reset specific steps only
plexr reset setup.yml --steps install_tools,setup_database
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--state-file` | `-s` | Custom state file location | `.plexr_state.json` |
| `--force` | `-f` | Skip confirmation prompt | `false` |
| `--steps` | | Reset specific steps only | all |

## completion

Generate shell completion scripts.

### Usage

```bash
plexr completion [shell]
```

### Supported Shells

- bash
- zsh
- fish
- powershell

### Examples

```bash
# Bash
plexr completion bash > /etc/bash_completion.d/plexr

# Zsh
plexr completion zsh > "${fpath[1]}/_plexr"

# Fish
plexr completion fish > ~/.config/fish/completions/plexr.fish

# PowerShell
plexr completion powershell > plexr.ps1
```

## version

Display version information.

### Usage

```bash
plexr version [flags]
```

### Examples

```bash
# Simple version
plexr version
# Output: plexr version 1.0.0

# Detailed version info
plexr version --verbose
# Output:
# plexr version 1.0.0
# Go version: go1.21.5
# Built: 2023-12-15T10:30:00Z
# Commit: abc123def
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Show detailed version info | `false` |
| `--json` | `-j` | Output in JSON format | `false` |

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Config file location | `$HOME/.plexr/config.yml` |
| `--log-level` | `-l` | Log level (debug, info, warn, error) | `info` |
| `--no-color` | | Disable colored output | `false` |
| `--help` | `-h` | Show help for command | |

## Environment Variables

Plexr respects these environment variables:

```bash
# Override state file location
export PLEXR_STATE_FILE=/tmp/my-state.json

# Set log level
export PLEXR_LOG_LEVEL=debug

# Disable colors
export PLEXR_NO_COLOR=true

# Override platform
export PLEXR_PLATFORM=linux
```

## Exit Codes

Plexr uses standard exit codes:

- `0`: Success
- `1`: General error
- `2`: Invalid arguments
- `3`: Plan validation failed
- `4`: Execution failed
- `5`: State corruption
- `130`: Interrupted (Ctrl+C)

## Advanced Usage

### Chaining Commands

```bash
# Validate, then execute if successful
plexr validate setup.yml && plexr execute setup.yml

# Reset and execute in one line
plexr reset setup.yml --force && plexr execute setup.yml --auto
```

### Using with CI/CD

```bash
# CI-friendly execution
plexr execute setup.yml \
  --auto \
  --platform=linux \
  --log-level=debug \
  --no-color
```

### Debugging

```bash
# Maximum verbosity
PLEXR_LOG_LEVEL=debug plexr execute setup.yml --verbose

# Dry run with verbose output
plexr execute setup.yml --dry-run --verbose
```

## Next Steps

- Learn about [State Management](/guide/state-management)
- See [Examples](/examples/) for real-world usage
- Read about [Executors](/guide/executors) for extending Plexr