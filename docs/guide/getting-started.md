# Getting Started

This guide will help you get started with Plexr in just a few minutes.

## What is Plexr?

Plexr (Plan + Executor) is a CLI tool that automates local development environment setup through YAML-based execution plans. It solves the "works on my machine" problem by ensuring consistent, reproducible environments across your team.

## Prerequisites

- Go 1.21 or higher (for installation from source)
- Basic familiarity with YAML
- Command line terminal

## Installation

### From Source

The easiest way to install Plexr is using Go:

```bash
go install github.com/SphereStacking/plexr@latest
```

### Binary Releases

Coming soon! Binary releases for major platforms will be available.

## Your First Plan

Let's create a simple execution plan to understand how Plexr works.

### 1. Create a Plan File

Create a file named `setup.yml`:

```yaml
name: "Hello Plexr"
version: "1.0.0"
description: "My first Plexr execution plan"

executors:
  shell:
    type: shell

steps:
  - id: welcome
    description: "Welcome message"
    executor: shell
    files:
      - path: "scripts/welcome.sh"
```

### 2. Create the Script

Create a directory for scripts and add the welcome script:

```bash
mkdir scripts
```

Create `scripts/welcome.sh`:

```bash
#!/bin/bash
echo "🚀 Welcome to Plexr!"
echo "Your development environment automation journey begins here."
```

### 3. Run the Plan

Execute your plan:

```bash
plexr execute setup.yml
```

You should see output like:

```
🚀 Starting execution of "Hello Plexr"
✓ Step 1/1: Welcome message
🎉 Execution completed successfully!
```

## Understanding Plans

A Plexr plan consists of:

- **Metadata**: Name, version, and description
- **Executors**: Define how to run different types of files
- **Steps**: The actual tasks to perform

### Step Structure

Each step has:
- `id`: Unique identifier
- `description`: Human-readable description
- `executor`: Which executor to use
- `files`: List of files to execute

## Key Features

### 1. State Management

Plexr tracks execution state, so if something fails, you can fix it and resume:

```bash
# Check current state
plexr status

# Resume from failure
plexr execute setup.yml
```

### 2. Platform Support

Handle different operating systems elegantly:

```yaml
files:
  - path: "scripts/install.sh"
    platform: linux
  - path: "scripts/install.ps1"
    platform: windows
```

### 3. Dependencies

Define execution order with dependencies:

```yaml
steps:
  - id: install_tools
    description: "Install required tools"
    executor: shell
    files:
      - path: "scripts/install.sh"
  
  - id: configure_tools
    description: "Configure tools"
    executor: shell
    depends_on: [install_tools]
    files:
      - path: "scripts/configure.sh"
```

### 4. Skip Conditions

Skip steps based on conditions:

```yaml
steps:
  - id: install_docker
    description: "Install Docker"
    executor: shell
    check_command: "docker --version"
    files:
      - path: "scripts/install_docker.sh"
```

### 5. Working Directory

Control where scripts execute:

```yaml
# Global working directory for all steps
work_directory: "/path/to/project"

steps:
  - id: build
    description: "Build project"
    executor: shell
    # Override for specific step
    work_directory: "/tmp/build"
    files:
      - path: "scripts/build.sh"
```

## Next Steps

Now that you understand the basics:

1. Read the [Configuration Guide](/guide/configuration) for detailed YAML options
2. Explore [Examples](/examples/) for real-world patterns
3. Learn about [Commands](/guide/commands) for advanced usage
4. Understand [State Management](/guide/state-management) for complex workflows

## Getting Help

- Check our [troubleshooting guide](/guide/troubleshooting)
- Open an issue on [GitHub](https://github.com/SphereStacking/plexr/issues)
- Read the [FAQ](/guide/faq)