# Plexr 🚀

A developer-friendly CLI tool for automating local development environment setup with YAML-based execution plans.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
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

Plexr is currently in active development. We're working towards a stable v1.0 release.

- [x] Core execution engine
- [x] Shell executor
- [x] SQL executor
- [x] State management
- [ ] Platform detection
- [ ] Interactive help system
- [ ] Plugin system

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
