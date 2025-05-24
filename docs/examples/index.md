# Examples

Learn how to use Plexr through practical examples. These examples cover common use cases and demonstrate best practices.

## Available Examples

### [Basic Setup](/examples/basic-setup)
A simple example showing how to set up a basic development environment with common tools.

**What you'll learn:**
- Creating your first plan
- Using platform-specific scripts
- Basic step dependencies

### [Advanced Patterns](/examples/advanced-patterns)
Advanced techniques for complex setups.

**What you'll learn:**
- Complex dependency chains
- Conditional execution
- Error handling strategies
- State management

### [Real World Examples](/examples/real-world)
Production-ready examples from actual projects.

**What you'll learn:**
- Full-stack application setup
- Database migrations
- CI/CD integration
- Team collaboration patterns

## Quick Start Example

Here's a minimal example to get you started:

```yaml
name: "Quick Start"
version: "1.0.0"

executors:
  shell:
    type: shell

steps:
  - id: hello_world
    description: "Say hello"
    executor: shell
    files:
      - path: "hello.sh"
```

With `hello.sh`:
```bash
#!/bin/bash
echo "Hello from Plexr! 🚀"
```

Run it:
```bash
plexr execute quickstart.yml
```

## Common Patterns

### 1. Tool Installation with Verification

```yaml
steps:
  - id: install_node
    description: "Install Node.js"
    executor: shell
    check_command: "node --version"
    files:
      - path: "install_node.sh"
```

### 2. Platform-Specific Scripts

```yaml
steps:
  - id: install_deps
    description: "Install dependencies"
    executor: shell
    files:
      - path: "install_linux.sh"
        platform: linux
      - path: "install_mac.sh"
        platform: darwin
      - path: "install_windows.ps1"
        platform: windows
```

### 3. Sequential Dependencies

```yaml
steps:
  - id: step1
    description: "First step"
    executor: shell
    files:
      - path: "step1.sh"
  
  - id: step2
    description: "Second step"
    depends_on: [step1]
    executor: shell
    files:
      - path: "step2.sh"
  
  - id: step3
    description: "Third step"
    depends_on: [step2]
    executor: shell
    files:
      - path: "step3.sh"
```

### 4. Parallel Execution

```yaml
steps:
  - id: download_tools
    description: "Download tools"
    executor: shell
    files:
      - path: "download.sh"
  
  - id: create_dirs
    description: "Create directories"
    executor: shell
    files:
      - path: "mkdirs.sh"
  
  - id: setup_configs
    description: "Setup configurations"
    depends_on: [download_tools, create_dirs]
    executor: shell
    files:
      - path: "configure.sh"
```

## Best Practices from Examples

### 1. Always Use Check Commands

```yaml
check_command: "command -v docker >/dev/null 2>&1"
```

### 2. Make Scripts Idempotent

```bash
# Bad
mkdir ~/workspace

# Good
mkdir -p ~/workspace
```

### 3. Handle Errors Gracefully

```bash
set -euo pipefail

if ! command -v node &> /dev/null; then
    echo "Node.js is required but not installed"
    exit 1
fi
```

### 4. Use Descriptive Step IDs

```yaml
# Bad
id: step1

# Good
id: install_postgresql_14
```

## Example Repository Structure

A typical project using Plexr:

```
my-project/
├── setup.yml              # Main plan file
├── scripts/              # Execution scripts
│   ├── install/
│   │   ├── node.sh
│   │   ├── docker.sh
│   │   └── postgres.sh
│   ├── configure/
│   │   ├── git.sh
│   │   └── env.sh
│   └── verify/
│       └── health_check.sh
├── sql/                  # SQL scripts
│   ├── create_db.sql
│   └── migrations/
└── configs/              # Configuration templates
    ├── .env.template
    └── docker-compose.yml
```

## Contributing Examples

Have a great example? We'd love to include it! 

1. Fork the repository
2. Add your example to `examples/`
3. Include a README explaining the use case
4. Submit a pull request

## Next Steps

- Try the [Basic Setup Example](/examples/basic-setup)
- Read the [Configuration Guide](/guide/configuration)
- Check the [API Reference](/api/) for detailed options