# Configuration

Plexr uses YAML files to define execution plans. This guide covers all configuration options in detail.

## File Structure

A Plexr configuration file has the following top-level structure:

```yaml
name: string          # Required: Plan name
version: string       # Required: Plan version
description: string   # Optional: Plan description
platforms: map        # Optional: Platform-specific settings
executors: map        # Required: Executor configurations
steps: array          # Required: Execution steps
```

## Basic Example

```yaml
name: "Development Environment Setup"
version: "1.0.0"
description: |
  Sets up a complete development environment including:
  - Development tools
  - Database
  - Configuration files

executors:
  shell:
    type: shell
    config:
      shell: /bin/bash

steps:
  - id: install_tools
    description: "Install development tools"
    executor: shell
    files:
      - path: "scripts/install.sh"
```

## Metadata Fields

### name (Required)

The name of your execution plan:

```yaml
name: "My Project Setup"
```

### version (Required)

Semantic version of your plan:

```yaml
version: "1.0.0"
```

### description (Optional)

Detailed description supporting multiline text:

```yaml
description: |
  This plan sets up:
  - Node.js environment
  - PostgreSQL database
  - Redis cache
```

## Executors

Executors define how different types of files are executed.

### Shell Executor

The built-in shell executor runs shell scripts:

```yaml
executors:
  shell:
    type: shell
    config:
      shell: /bin/bash      # Linux/macOS
      # shell: powershell.exe  # Windows
      timeout: 300          # Default timeout in seconds
      env:                  # Environment variables
        NODE_ENV: development
        DEBUG: "true"
```

### Custom Executors

Future versions will support custom executors:

```yaml
executors:
  sql:
    type: sql
    config:
      driver: postgres
      host: localhost
      port: 5432
      database: myapp
      user: postgres
```

## Steps

Steps define the tasks to execute in order.

### Basic Step

```yaml
steps:
  - id: create_directories
    description: "Create project directories"
    executor: shell
    files:
      - path: "scripts/create_dirs.sh"
```

### Step Fields

#### id (Required)

Unique identifier for the step:

```yaml
id: install_dependencies
```

#### description (Required)

Human-readable description:

```yaml
description: "Install Node.js dependencies"
```

#### executor (Required)

Which executor to use:

```yaml
executor: shell
```

#### files (Required)

List of files to execute:

```yaml
files:
  - path: "scripts/install.sh"
    timeout: 600
    retry: 3
```

#### depends_on (Optional)

Dependencies that must complete first:

```yaml
depends_on: [install_tools, create_directories]
```

#### skip_if (Optional)

Shell command to check if step should be skipped:

```yaml
skip_if: "test -f /usr/local/bin/node"
```

#### check_command (Optional)

Command to verify if step is already completed:

```yaml
check_command: "docker --version"
```

## File Configuration

Each file in a step can have additional configuration:

```yaml
files:
  - path: "scripts/install.sh"
    platform: linux         # Platform-specific file
    timeout: 300           # Override timeout (seconds)
    retry: 3              # Number of retries on failure
    skip_if: "test -f /usr/local/bin/tool"
```

### Platform-Specific Files

Handle different operating systems:

```yaml
files:
  - path: "scripts/install.sh"
    platform: linux
  - path: "scripts/install.sh"
    platform: darwin    # macOS
  - path: "scripts/install.ps1"
    platform: windows
```

Platform values:
- `linux`: Linux systems
- `darwin`: macOS
- `windows`: Windows

## Platform Configuration

Define platform-specific variables:

```yaml
platforms:
  linux:
    install_prefix: /usr/local
    package_manager: apt
  darwin:
    install_prefix: /opt/homebrew
    package_manager: brew
  windows:
    install_prefix: C:\Program Files
    package_manager: choco
```

Access in scripts:
```bash
echo "Installing to: ${PLEXR_PLATFORM_install_prefix}"
```

## Advanced Features

### Transaction Mode

Group steps for atomic execution (coming soon):

```yaml
steps:
  - id: database_migration
    description: "Run database migrations"
    executor: sql
    transaction_mode: all  # all, none, or per-file
    files:
      - path: "migrations/001_create_tables.sql"
      - path: "migrations/002_add_indexes.sql"
```

### Conditional Execution

Skip steps based on conditions:

```yaml
steps:
  - id: install_docker
    description: "Install Docker if not present"
    executor: shell
    skip_if: "command -v docker >/dev/null 2>&1"
    files:
      - path: "scripts/install_docker.sh"
```

### Dependency Chains

Create complex workflows:

```yaml
steps:
  - id: install_node
    description: "Install Node.js"
    executor: shell
    files:
      - path: "scripts/install_node.sh"
  
  - id: install_npm_packages
    description: "Install NPM packages"
    executor: shell
    depends_on: [install_node]
    files:
      - path: "scripts/npm_install.sh"
  
  - id: build_project
    description: "Build the project"
    executor: shell
    depends_on: [install_npm_packages]
    files:
      - path: "scripts/build.sh"
```

## Environment Variables

Plexr provides several environment variables to scripts:

- `PLEXR_STEP_ID`: Current step ID
- `PLEXR_STEP_INDEX`: Current step number
- `PLEXR_TOTAL_STEPS`: Total number of steps
- `PLEXR_PLATFORM`: Current platform (linux, darwin, windows)
- `PLEXR_STATE_FILE`: Path to state file
- `PLEXR_DRY_RUN`: "true" if in dry-run mode

## Best Practices

### 1. Use Descriptive IDs

```yaml
# Good
id: install_postgresql_14

# Bad
id: step1
```

### 2. Check Prerequisites

```yaml
steps:
  - id: configure_git
    description: "Configure Git settings"
    check_command: "git config --get user.name"
    files:
      - path: "scripts/git_config.sh"
```

### 3. Handle Errors Gracefully

```bash
#!/bin/bash
set -euo pipefail  # Exit on error

# Check prerequisites
if ! command -v node &> /dev/null; then
    echo "Error: Node.js is required but not installed"
    exit 1
fi
```

### 4. Make Scripts Idempotent

```bash
# Create directory only if it doesn't exist
mkdir -p "$HOME/projects"

# Install only if not present
if ! command -v tool &> /dev/null; then
    install_tool
fi
```

### 5. Use Platform Detection

```yaml
files:
  - path: "scripts/install_common.sh"
  - path: "scripts/install_mac.sh"
    platform: darwin
  - path: "scripts/install_linux.sh"
    platform: linux
```

## Validation

Validate your configuration before running:

```bash
plexr validate setup.yml
```

This checks for:
- YAML syntax errors
- Required fields
- Circular dependencies
- File existence
- Executor availability

## Next Steps

- See [Examples](/examples/) for real-world configurations
- Learn about [Commands](/guide/commands) to use these configurations
- Understand [State Management](/guide/state-management) for complex workflows