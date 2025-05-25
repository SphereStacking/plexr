# Plexr 🚀

A developer-friendly CLI tool for automating local development environment setup with YAML-based execution plans.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/SphereStacking/plexr?include_prereleases)](https://github.com/SphereStacking/plexr/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

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

#### Using go install (Recommended)
```bash
go install github.com/SphereStacking/plexr/cmd/plexr@latest
```

#### Download Pre-built Binary
Download the latest binary from [Releases](https://github.com/SphereStacking/plexr/releases) page.

```bash
# Linux (x86_64)
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Linux_x86_64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# macOS (Intel)
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Darwin_x86_64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# macOS (Apple Silicon)
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Darwin_arm64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# Windows
# Download plexr_Windows_x86_64.zip from releases page and extract
```

#### Build from Source
```bash
git clone https://github.com/SphereStacking/plexr.git
cd plexr
make build
sudo make install
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

## 📦 Current Release

**Latest Version**: [v0.1.0](https://github.com/SphereStacking/plexr/releases/tag/v0.1.0)

### Features Available in v0.1.0

✅ **Core Features**
- Shell executor for running scripts
- SQL executor for PostgreSQL databases
- State management for resumable execution
- YAML-based execution plans
- Dependency resolution between steps
- Environment variable expansion
- Transaction mode support (none, each, all)

✅ **CLI Commands**
- `execute` - Run execution plans
- `validate` - Validate plan syntax
- `status` - Check execution status
- `reset` - Reset execution state
- `version` - Show version information
- `completion` - Generate shell completions

✅ **Platform Support**
- Linux (x86_64, arm64)
- macOS (Intel, Apple Silicon)
- Windows (x86_64)

### Roadmap

🚧 **Planned for Future Releases**
- SQL executor for MySQL and SQLite
- Advanced platform detection
- Interactive error recovery
- Plugin system
- Web UI for monitoring
- Remote state storage

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
