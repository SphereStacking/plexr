# Configuration Schema

Complete reference for Plexr YAML configuration files.

## Schema Overview

```yaml
# Required fields
name: string
version: string
executors: map<string, ExecutorConfig>
steps: array<Step>

# Optional fields
description: string
platforms: map<string, map<string, string>>
```

## Root Fields

### name

**Type:** `string` (required)  
**Description:** Name of the execution plan

```yaml
name: "Development Environment Setup"
```

### version

**Type:** `string` (required)  
**Description:** Semantic version of the plan

```yaml
version: "1.0.0"
```

### description

**Type:** `string` (optional)  
**Description:** Detailed description of the plan

```yaml
description: |
  This plan sets up a complete development environment including:
  - Required tools and dependencies
  - Database initialization
  - Configuration files
```

### executors

**Type:** `map<string, ExecutorConfig>` (required)  
**Description:** Map of executor configurations

```yaml
executors:
  shell:
    type: shell
    config:
      shell: /bin/bash
```

### steps

**Type:** `array<Step>` (required)  
**Description:** List of execution steps

```yaml
steps:
  - id: install_tools
    description: "Install required tools"
    executor: shell
    files:
      - path: "scripts/install.sh"
```

### platforms

**Type:** `map<string, map<string, string>>` (optional)  
**Description:** Platform-specific variables

```yaml
platforms:
  linux:
    package_manager: apt
    install_prefix: /usr/local
  darwin:
    package_manager: brew
    install_prefix: /opt/homebrew
```

## ExecutorConfig

Configuration for an executor.

### type

**Type:** `string` (required)  
**Description:** Type of executor  
**Values:** `shell`, `sql` (future), `http` (future)

### config

**Type:** `map<string, any>` (optional)  
**Description:** Executor-specific configuration

#### Shell Executor Config

```yaml
executors:
  shell:
    type: shell
    config:
      shell: string          # Shell to use (default: /bin/bash or powershell.exe)
      timeout: integer       # Default timeout in seconds (default: 300)
      env: map<string, string>  # Environment variables
      working_dir: string    # Working directory
```

Example:
```yaml
executors:
  shell:
    type: shell
    config:
      shell: /bin/zsh
      timeout: 600
      env:
        NODE_ENV: development
        DEBUG: "true"
      working_dir: /tmp
```

## Step

Definition of a single execution step.

### id

**Type:** `string` (required)  
**Description:** Unique identifier for the step  
**Pattern:** `^[a-zA-Z][a-zA-Z0-9_-]*$`

```yaml
id: install_dependencies
```

### description

**Type:** `string` (required)  
**Description:** Human-readable description

```yaml
description: "Install Node.js dependencies"
```

### executor

**Type:** `string` (required)  
**Description:** Name of executor to use  
**Must match:** Key in `executors` map

```yaml
executor: shell
```

### files

**Type:** `array<FileConfig>` (required)  
**Description:** List of files to execute

```yaml
files:
  - path: "scripts/install.sh"
    timeout: 300
```

### depends_on

**Type:** `array<string>` (optional)  
**Description:** Step IDs that must complete before this step

```yaml
depends_on: [install_tools, create_directories]
```

### skip_if

**Type:** `string` (optional)  
**Description:** Shell command to determine if step should be skipped

```yaml
skip_if: "test -f /usr/local/bin/node"
```

### check_command

**Type:** `string` (optional)  
**Description:** Command to verify if step is already completed

```yaml
check_command: "docker --version"
```

### transaction_mode

**Type:** `string` (optional)  
**Description:** Transaction handling mode  
**Values:** `all`, `none`, `per-file`  
**Default:** `none`

```yaml
transaction_mode: all
```

## FileConfig

Configuration for a file to be executed.

### path

**Type:** `string` (required)  
**Description:** Path to the file relative to plan location

```yaml
path: "scripts/install.sh"
```

### platform

**Type:** `string` (optional)  
**Description:** Platform restriction  
**Values:** `linux`, `darwin`, `windows`

```yaml
platform: linux
```

### timeout

**Type:** `integer` (optional)  
**Description:** Execution timeout in seconds  
**Default:** Executor's default timeout

```yaml
timeout: 600
```

### retry

**Type:** `integer` (optional)  
**Description:** Number of retry attempts on failure  
**Default:** 0

```yaml
retry: 3
```

### skip_if

**Type:** `string` (optional)  
**Description:** Shell command to determine if file should be skipped

```yaml
skip_if: "test -d node_modules"
```

## Complete Example

```yaml
name: "Full Stack Development Environment"
version: "2.1.0"
description: |
  Complete setup for full-stack development including:
  - Node.js and npm packages
  - PostgreSQL database
  - Redis cache
  - Docker containers

platforms:
  linux:
    package_manager: apt
    postgres_package: postgresql-14
  darwin:
    package_manager: brew
    postgres_package: postgresql@14

executors:
  shell:
    type: shell
    config:
      timeout: 300
      env:
        NODE_ENV: development
  
  sql:
    type: sql
    config:
      driver: postgres
      database: myapp_dev

steps:
  - id: install_system_deps
    description: "Install system dependencies"
    executor: shell
    files:
      - path: "scripts/install_deps.sh"
        platform: linux
      - path: "scripts/install_deps_mac.sh"
        platform: darwin

  - id: install_node
    description: "Install Node.js"
    executor: shell
    depends_on: [install_system_deps]
    check_command: "node --version"
    files:
      - path: "scripts/install_node.sh"
        timeout: 600

  - id: install_npm_packages
    description: "Install npm packages"
    executor: shell
    depends_on: [install_node]
    skip_if: "test -d node_modules"
    files:
      - path: "scripts/npm_install.sh"
        retry: 3

  - id: setup_database
    description: "Initialize PostgreSQL database"
    executor: sql
    depends_on: [install_system_deps]
    transaction_mode: all
    files:
      - path: "sql/create_database.sql"
      - path: "sql/create_tables.sql"
      - path: "sql/seed_data.sql"

  - id: configure_environment
    description: "Set up environment configuration"
    executor: shell
    depends_on: [install_npm_packages, setup_database]
    files:
      - path: "scripts/setup_env.sh"
```

## Validation Rules

1. **Unique Step IDs:** All step IDs must be unique
2. **Valid Dependencies:** Steps in `depends_on` must exist
3. **No Circular Dependencies:** Dependency graph must be acyclic
4. **Existing Executors:** Step executors must be defined
5. **File Paths:** Paths should be relative to plan file
6. **Platform Values:** Must be `linux`, `darwin`, or `windows`

## Environment Variables in Scripts

Plexr provides these variables to executed scripts:

| Variable | Description |
|----------|-------------|
| `PLEXR_STEP_ID` | Current step ID |
| `PLEXR_STEP_INDEX` | Current step number (1-based) |
| `PLEXR_TOTAL_STEPS` | Total number of steps |
| `PLEXR_PLATFORM` | Current platform |
| `PLEXR_STATE_FILE` | Path to state file |
| `PLEXR_DRY_RUN` | "true" if in dry-run mode |
| `PLEXR_PLATFORM_*` | Platform-specific variables |

Example usage in scripts:
```bash
echo "Running step $PLEXR_STEP_INDEX of $PLEXR_TOTAL_STEPS"
echo "Package manager: $PLEXR_PLATFORM_package_manager"
```