# CLI Commands Reference

This is a comprehensive reference for all Plexr CLI commands.

## Command Structure

```
plexr [global-flags] <command> [command-flags] [arguments]
```

## Commands

### plexr execute

Execute a plan file.

```bash
plexr execute <plan-file> [flags]
```

#### Arguments

- `<plan-file>` - Path to the YAML plan file (required)

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dry-run, -n` | bool | false | Preview execution without making changes |
| `--auto, -y` | bool | false | Automatically confirm all prompts |
| `--platform, -p` | string | auto | Override platform detection (linux, darwin, windows) |
| `--state-file, -s` | string | .plexr_state.json | Path to state file |
| `--force, -f` | bool | false | Force re-execution of completed steps |
| `--verbose, -v` | bool | false | Enable verbose output |

#### Examples

```bash
# Basic execution
plexr execute setup.yml

# Dry run
plexr execute setup.yml --dry-run

# Force re-execution
plexr execute setup.yml --force

# Custom state file
plexr execute setup.yml --state-file=/tmp/state.json
```

#### Exit Codes

- `0` - Successful execution
- `1` - General error
- `4` - Execution failed
- `130` - User interruption (Ctrl+C)

---

### plexr validate

Validate a plan file without executing it.

```bash
plexr validate <plan-file> [flags]
```

#### Arguments

- `<plan-file>` - Path to the YAML plan file (required)

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--check-files, -c` | bool | false | Verify that referenced files exist |
| `--verbose, -v` | bool | false | Show detailed validation information |

#### Examples

```bash
# Basic validation
plexr validate setup.yml

# Check file existence
plexr validate setup.yml --check-files

# Verbose output
plexr validate setup.yml --verbose
```

#### Exit Codes

- `0` - Valid plan
- `3` - Validation failed
- `2` - Invalid arguments

---

### plexr status

Display the current execution status.

```bash
plexr status <plan-file> [flags]
```

#### Arguments

- `<plan-file>` - Path to the YAML plan file (required)

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state-file, -s` | string | .plexr_state.json | Path to state file |
| `--verbose, -v` | bool | false | Show detailed status information |
| `--json, -j` | bool | false | Output in JSON format |

#### Examples

```bash
# Show status
plexr status setup.yml

# JSON output
plexr status setup.yml --json

# Custom state file
plexr status setup.yml --state-file=/tmp/state.json
```

#### Output Format (JSON)

```json
{
  "plan": {
    "name": "Development Environment Setup",
    "version": "1.0.0"
  },
  "state": "in_progress",
  "progress": {
    "completed": 3,
    "total": 5,
    "percentage": 60
  },
  "current_step": "configure_app",
  "steps": [
    {
      "id": "install_tools",
      "status": "completed",
      "description": "Install development tools"
    }
  ],
  "last_updated": "2023-12-15T10:30:45Z"
}
```

---

### plexr reset

Reset the execution state.

```bash
plexr reset <plan-file> [flags]
```

#### Arguments

- `<plan-file>` - Path to the YAML plan file (required)

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state-file, -s` | string | .plexr_state.json | Path to state file |
| `--force, -f` | bool | false | Skip confirmation prompt |
| `--steps` | string | | Comma-separated list of step IDs to reset |

#### Examples

```bash
# Reset with confirmation
plexr reset setup.yml

# Force reset
plexr reset setup.yml --force

# Reset specific steps
plexr reset setup.yml --steps=install_tools,setup_database
```

---

### plexr completion

Generate shell completion scripts.

```bash
plexr completion <shell>
```

#### Arguments

- `<shell>` - Target shell: bash, zsh, fish, or powershell (required)

#### Examples

```bash
# Bash
plexr completion bash > /etc/bash_completion.d/plexr

# Zsh
plexr completion zsh > "${fpath[1]}/_plexr"

# Fish
plexr completion fish > ~/.config/fish/completions/plexr.fish

# PowerShell
plexr completion powershell | Out-String | Invoke-Expression
```

---

### plexr version

Display version information.

```bash
plexr version [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--verbose, -v` | bool | false | Show detailed version information |
| `--json, -j` | bool | false | Output in JSON format |

#### Examples

```bash
# Simple version
plexr version

# Detailed version
plexr version --verbose

# JSON output
plexr version --json
```

#### Output Format (JSON)

```json
{
  "version": "1.0.0",
  "go_version": "go1.21.5",
  "build_time": "2023-12-15T10:30:00Z",
  "git_commit": "abc123def",
  "platform": "linux/amd64"
}
```

---

### plexr help

Display help information.

```bash
plexr help [command]
```

#### Arguments

- `[command]` - Optional command to get help for

#### Examples

```bash
# General help
plexr help

# Command-specific help
plexr help execute

# Also works with --help flag
plexr execute --help
```

## Global Flags

These flags are available for all commands:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config, -c` | string | $HOME/.plexr/config.yml | Config file path |
| `--log-level, -l` | string | info | Log level: debug, info, warn, error |
| `--no-color` | bool | false | Disable colored output |
| `--help, -h` | bool | false | Show help information |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PLEXR_STATE_FILE` | Override default state file location | .plexr_state.json |
| `PLEXR_LOG_LEVEL` | Set log level | info |
| `PLEXR_NO_COLOR` | Disable colored output | false |
| `PLEXR_PLATFORM` | Override platform detection | auto |
| `PLEXR_CONFIG` | Override config file location | $HOME/.plexr/config.yml |

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments or usage |
| 3 | Plan validation failed |
| 4 | Execution failed |
| 5 | State file corrupted |
| 130 | Interrupted by user (Ctrl+C) |

## Command Aliases

Some commands support shorter aliases:

- `exec` → `execute`
- `val` → `validate`
- `stat` → `status`

Example:
```bash
plexr exec setup.yml
```