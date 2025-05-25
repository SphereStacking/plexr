---
layout: home

hero:
  name: "Plexr"
  text: "Plan + Executor"
  tagline: "Developer-friendly CLI tool for automating local development environment setup"
  image:
    src: https://api.iconify.design/noto:rocket.svg
    alt: Plexr
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/SphereStacking/plexr

features:
  - icon: 📝
    title: Executable Documentation
    details: Turn your README setup instructions into executable YAML configurations
  - icon: 🔄
    title: Stateful Execution
    details: Resume from where you left off if something fails - no need to start over
  - icon: 🖥️
    title: Cross-Platform
    details: Works seamlessly on macOS, Linux, and Windows with platform-specific support
  - icon: 🛡️
    title: Safe by Default
    details: Dry-run mode, skip conditions, and rollback support ensure safe operations
---

## 🎉 Latest Release: v0.1.0

Download the latest release from [GitHub Releases](https://github.com/SphereStacking/plexr/releases/tag/v0.1.0) or install with:

```bash
go install github.com/SphereStacking/plexr/cmd/plexr@v0.1.0
```

## Quick Start

Install Plexr and get started in minutes:

```bash
# Install latest release
go install github.com/SphereStacking/plexr/cmd/plexr@latest

# Or download pre-built binary
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_$(uname -s)_$(uname -m | sed 's/x86_64/x86_64/;s/aarch64/arm64/').tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# Run your first plan
plexr execute setup.yml
```

## Why Plexr?

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

## Example Configuration

```yaml
name: "My Project Setup"
version: "1.0.0"

steps:
  - id: install_deps
    description: "Install dependencies"
    executor: shell
    files:
      - path: "scripts/install.sh"
        platform: linux
      - path: "scripts/install.ps1"
        platform: windows

  - id: setup_database
    description: "Initialize database"
    executor: shell
    depends_on: [install_deps]
    files:
      - path: "scripts/db_setup.sh"
```

## Features

### v0.1.0 (Current Release)
- ✅ Core execution engine with dependency resolution
- ✅ Shell executor for running scripts and commands
- ✅ SQL executor with PostgreSQL support
- ✅ State management with resume capability
- ✅ CLI commands (execute, validate, status, reset)
- ✅ Environment variable expansion
- ✅ Platform-specific file selection
- ✅ Error handling and rollback support

### Coming Soon
- 🚧 Additional database support (MySQL, SQLite)
- 🚧 HTTP executor for API calls
- 🚧 Docker executor
- 🚧 Parallel execution
- 🚧 Advanced conditional logic

## Learn More

- [Installation Guide](/guide/installation) - Get Plexr installed on your system
- [Configuration Reference](/guide/configuration) - Learn about YAML configuration
- [Examples](/examples/) - See real-world usage patterns
- [API Documentation](/api/) - Detailed technical reference