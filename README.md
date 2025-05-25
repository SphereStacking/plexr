# Plexr 🚀

A developer-friendly CLI tool for automating local development environment setup with YAML-based execution plans.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> ⚠️ **Note**: This project was created through vibe coding sessions. While the core concepts and architecture are in place, some features described in this documentation may not be fully implemented yet. Please refer to the [Project Status](#-project-status) section for details on what's currently available.

## 🎯 What is Plexr?

Plexr (Plan + Executor) helps developers set up and maintain their local development environments through simple YAML configuration files. No more "works on my machine" issues or spending hours following outdated setup documentation.

### Key Features

- 📝 **Executable Documentation**: Turn your README setup instructions into executable YAML
- 🔄 **Stateful Execution**: Resume from where you left off if something fails
- 💬 **Interactive Error Handling**: Get helpful suggestions when things go wrong
- 🖥️ **Cross-Platform**: Works on macOS, Linux, and Windows
- 🛡️ **Safe by Default**: Dry-run mode, skip conditions, and rollback support
- 🚦 **Dependency Management**: Define execution order and prerequisites

## 🚀 Quick Start

### Installation

```bash
# macOS
brew install plexr

# Linux/Windows
curl -sSL https://raw.githubusercontent.com/SphereStacking/plexr/main/install.sh | bash

# From source
go install github.com/SphereStacking/plexr@latest
```

### Basic Usage

1. Create a `setup.yml` file:

```yaml
name: "My Project Setup"
version: "1.0.0"

executors:
  shell:
    type: shell
  sql:
    type: sql
    driver: postgres
    database: myapp_dev

steps:
  - id: check_tools
    description: "Checking required tools"
    executor: shell
    files:
      - path: "setup/check_dependencies.sh"

  - id: setup_database
    description: "Setting up database"
    executor: sql
    depends_on: [check_tools]
    skip_if: "SELECT 1 FROM pg_database WHERE datname = 'myapp_dev'"
    files:
      - path: "db/schema.sql"
      - path: "db/seed_data.sql"
```

2. Run the setup:

```bash
plexr execute setup.yml
```

### 🚀 Shell Completion

Plexr provides git-style command completion with descriptions:

```bash
# Install completion for your shell
plexr completion bash > ~/.local/share/bash-completion/completions/plexr  # Bash
plexr completion zsh > "${fpath[1]}/_plexr"                              # Zsh
plexr completion fish > ~/.config/fish/completions/plexr.fish            # Fish

# After installation, you'll see completions like:
❯❯❯ plexr [TAB]
completion  -- Generate the autocompletion script for the specified shell
execute     -- Execute a setup plan
exec        -- alias for 'execute'
run         -- alias for 'execute'
help        -- Help about any command
reset       -- Reset execution state
status      -- Show execution status
validate    -- Validate an execution plan
val         -- alias for 'validate'
check       -- alias for 'validate'
version     -- Print version information
```

## 📖 Documentation

- [Getting Started Guide](docs/getting-started.md)
- [Configuration Reference](docs/configuration.md)
- [Examples](examples/)
- [FAQ](docs/faq.md)

## 🤝 Why Plexr?

### The Problem

- 😫 "I followed the README but it doesn't work"
- ⏰ "Setting up the dev environment took all day"
- 🤷 "It works on my machine"
- 🔧 "Everyone's environment is slightly different"

### The Solution

Plexr makes environment setup:
- **Reproducible**: Same result every time
- **Debuggable**: Clear error messages with solutions
- **Maintainable**: Version controlled setup procedures
- **Team-friendly**: Everyone uses the same configuration

## 🛠️ Use Cases

- **New Developer Onboarding**: Get new team members productive in hours, not days
- **Project Setup**: Initialize databases, install dependencies, configure tools
- **Environment Updates**: Apply schema changes, update configurations safely
- **Cross-Platform Development**: Handle OS-specific setup automatically

## 🏗️ Project Status

Plexr is currently in active development. This project was created through vibe coding sessions, and we're working towards a stable v1.0 release.

### Currently Implemented
- [x] Core execution engine
- [x] Shell executor
- [x] State management
- [x] Basic CLI commands (execute, validate, status, reset)
- [x] YAML configuration parsing
- [x] Dependency resolution

### In Progress / Not Yet Implemented
- [ ] SQL executor (documented but not implemented)
- [ ] Platform detection (partial implementation)
- [ ] Interactive help system
- [ ] Plugin system
- [ ] Transaction mode
- [ ] Some advanced features described in documentation

**Important**: While the documentation describes the complete vision for Plexr, not all features are currently implemented. The core functionality for basic execution plans is working, but advanced features like SQL executors, transaction modes, and some platform-specific handling may not be available yet.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

```bash
# Clone the repository
git clone https://github.com/SphereStacking/plexr.git
cd plexr

# Install dependencies
make deps

# Run tests
make test

# Build
make build
```

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Inspired by tools like Ansible, Make, and various database migration tools, but focused specifically on developer environment setup.

---

Made with ❤️ for developers who just want their environment to work.

## 📚 Documentation

Full documentation is available at: https://spherestacking.github.io/plexr/
