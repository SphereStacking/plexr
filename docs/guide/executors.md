# Executors Guide

This guide explains how to use and configure executors in Plexr, including built-in executors and creating custom ones for your specific needs.

## What are Executors?

Executors are the components that actually run your steps. They provide the bridge between your plan definitions and the execution environment. Each executor specializes in running specific types of commands or operations.

## Built-in Executors

### Shell Executor

The shell executor is the default and most commonly used executor. It runs shell commands in your system's default shell.

#### Basic Usage

```yaml
steps:
  - name: "Simple command"
    command: "echo Hello, World!"
    # executor: shell is implicit
```

#### Configuration Options

```yaml
steps:
  - name: "Configured shell command"
    command: "npm install"
    executor: shell
    config:
      shell: "/bin/bash"  # Specify shell (default: /bin/sh)
      timeout: 300        # Timeout in seconds
      workdir: "./app"    # Working directory
      env:               # Environment variables
        NODE_ENV: "production"
        CI: "true"
```

#### Shell Features

**Command Chaining**:
```yaml
steps:
  - name: "Multiple commands"
    command: |
      echo "Starting build"
      npm install
      npm run build
      echo "Build complete"
```

**Error Handling**:
```yaml
steps:
  - name: "Safe command"
    command: |
      set -e  # Exit on error
      set -u  # Exit on undefined variable
      set -o pipefail  # Fail on pipe errors
      
      command1 | command2 | command3
```

**Conditional Execution**:
```yaml
steps:
  - name: "Conditional logic"
    command: |
      if [ -f "package.json" ]; then
        npm install
      else
        echo "No package.json found"
      fi
```

### Script Executor (Future)

Execute scripts in various languages:

```yaml
steps:
  - name: "Python script"
    executor: script
    config:
      language: python
      version: "3.9"
    command: |
      import json
      data = {"status": "success"}
      print(json.dumps(data))
```

### HTTP Executor (Future)

Make HTTP requests:

```yaml
steps:
  - name: "Health check"
    executor: http
    config:
      method: GET
      url: "https://api.example.com/health"
      timeout: 30
      expected_status: 200
```

## Executor Selection

### Automatic Selection

Plexr automatically selects the appropriate executor based on the step configuration:

```yaml
steps:
  - name: "Shell command"
    command: "ls -la"  # Uses shell executor
    
  - name: "HTTP request"
    http:  # Would use http executor
      url: "https://example.com"
```

### Explicit Selection

Specify an executor explicitly:

```yaml
steps:
  - name: "Specific executor"
    executor: shell
    command: "echo 'Using shell executor'"
```

## Executor Configuration

### Global Configuration

Set default configurations for all steps using an executor:

```yaml
executors:
  shell:
    timeout: 120
    env:
      LOG_LEVEL: "info"
      
steps:
  - name: "Uses global config"
    command: "build.sh"
    
  - name: "Override global config"
    command: "long-task.sh"
    config:
      timeout: 600  # Override global timeout
```

### Per-Step Configuration

Configure executors for individual steps:

```yaml
steps:
  - name: "Custom configuration"
    command: "deploy.sh"
    executor: shell
    config:
      workdir: "/opt/app"
      timeout: 300
      env:
        DEPLOY_ENV: "production"
        DEPLOY_KEY: "`{{.deploy_key}}`"
```

## Working with Output

### Capturing Output

```yaml
steps:
  - name: "Capture output"
    command: "generate-report.sh"
    outputs:
      - name: report_path
        from: stdout
      - name: report_size
        from: stderr
        regex: "Size: ([0-9]+) bytes"
```

### Output Formats

**Plain Text**:
```yaml
outputs:
  - name: message
    from: stdout
```

**JSON Parsing**:
```yaml
outputs:
  - name: data
    from: stdout
    json_parse: true
  - name: specific_field
    from: stdout
    json_path: "$.result.id"
```

**Line Selection**:
```yaml
outputs:
  - name: first_line
    from: stdout
    line: 1
  - name: last_line
    from: stdout
    line: -1
```

## Error Handling

### Exit Codes

```yaml
steps:
  - name: "Handle exit codes"
    command: "test-command.sh"
    success_codes: [0, 1]  # Consider 0 and 1 as success
    ignore_failure: false   # Fail the plan on error
```

### Retry Logic

```yaml
steps:
  - name: "Retry on failure"
    command: "flaky-service.sh"
    retry:
      attempts: 3
      delay: 5s
      backoff: exponential  # or linear
      on_codes: [1, 2]     # Retry only on specific exit codes
```

### Error Patterns

```yaml
steps:
  - name: "Pattern-based retry"
    command: "connect-to-service.sh"
    retry:
      attempts: 5
      delay: 10s
      on_output_contains: ["connection refused", "timeout"]
```

## Parallel Execution

### Parallel Steps

```yaml
steps:
  - name: "Parallel execution"
    parallel:
      - name: "Frontend tests"
        command: "npm test"
        workdir: "./frontend"
        
      - name: "Backend tests"
        command: "go test ./..."
        workdir: "./backend"
        
      - name: "Integration tests"
        command: "pytest"
        workdir: "./tests"
```

### Parallel Configuration

```yaml
steps:
  - name: "Controlled parallelism"
    parallel:
      max_concurrent: 2  # Limit concurrent executions
      fail_fast: true    # Stop on first failure
      steps:
        - name: "Task 1"
          command: "task1.sh"
        - name: "Task 2"
          command: "task2.sh"
        - name: "Task 3"
          command: "task3.sh"
```

## Custom Executors

### Creating a Custom Executor

Implement the Executor interface:

```go
package myexecutor

import (
    "context"
    "github.com/plexr/plexr/internal/core"
)

type MyExecutor struct {
    // Executor-specific fields
}

func New() *MyExecutor {
    return &MyExecutor{}
}

func (e *MyExecutor) Execute(ctx context.Context, step core.Step, state core.State) error {
    // Implementation
    return nil
}

func (e *MyExecutor) Validate(step core.Step) error {
    // Validate step configuration
    return nil
}
```

### Registering Custom Executors

```go
// In your main.go or plugin
import (
    "github.com/plexr/plexr/internal/core"
    "myproject/executors/myexecutor"
)

func init() {
    core.RegisterExecutor("myexecutor", myexecutor.New())
}
```

### Using Custom Executors

```yaml
steps:
  - name: "Custom executor step"
    executor: myexecutor
    config:
      custom_option: "value"
    command: "custom command"
```

## Best Practices

### 1. Choose the Right Executor

- Use shell for system commands and scripts
- Use specialized executors for specific tasks (HTTP, database, etc.)
- Create custom executors for complex, reusable logic

### 2. Handle Errors Gracefully

```yaml
steps:
  - name: "Robust execution"
    command: |
      set -euo pipefail
      trap 'echo "Error on line $LINENO"' ERR
      
      # Your commands here
    retry:
      attempts: 3
      delay: 5s
```

### 3. Use Timeouts

```yaml
steps:
  - name: "Time-bound operation"
    command: "long-running-task.sh"
    config:
      timeout: 300  # 5 minutes
    on_timeout: 
      command: "cleanup.sh"  # Run on timeout
```

### 4. Secure Sensitive Data

```yaml
steps:
  - name: "Secure execution"
    command: "deploy.sh"
    env:
      API_KEY: "${API_KEY}"  # From environment
    config:
      mask_output: ["password", "secret", "key"]
```

### 5. Log Appropriately

```yaml
steps:
  - name: "Logged execution"
    command: "process.sh"
    config:
      log_level: "debug"  # debug, info, warn, error
      log_output: true    # Log command output
      log_file: "process.log"
```

## Troubleshooting Executors

### Debug Mode

```yaml
steps:
  - name: "Debug execution"
    command: "problematic-command.sh"
    debug: true  # Enable debug output
    config:
      verbose: true
      dry_run: true  # Show what would be executed
```

### Execution Context

```yaml
steps:
  - name: "Show context"
    command: |
      echo "Working dir: $(pwd)"
      echo "User: $(whoami)"
      echo "Shell: $SHELL"
      echo "PATH: $PATH"
      env | sort
```

### Common Issues

**Command not found**:
```yaml
steps:
  - name: "Explicit path"
    command: "/usr/local/bin/custom-tool"
    # Or update PATH
    env:
      PATH: "/usr/local/bin:${PATH}"
```

**Permission denied**:
```yaml
steps:
  - name: "Ensure permissions"
    command: |
      chmod +x script.sh
      ./script.sh
```

**Working directory issues**:
```yaml
steps:
  - name: "Absolute paths"
    command: "build.sh"
    config:
      workdir: "${PWD}/subdir"  # Use absolute path
```

## Performance Optimization

### Output Buffering

```yaml
steps:
  - name: "Large output"
    command: "generate-large-output.sh"
    config:
      buffer_size: 1048576  # 1MB buffer
      stream_output: true   # Stream instead of buffering
```

### Resource Limits

```yaml
steps:
  - name: "Resource constrained"
    command: "memory-intensive.sh"
    config:
      memory_limit: "2G"
      cpu_limit: 2
      nice: 10  # Lower priority
```

### Caching

```yaml
steps:
  - name: "Cached execution"
    command: "expensive-operation.sh"
    cache:
      key: "`{{.cache_key}}`"
      ttl: 3600  # 1 hour
      paths:
        - "./output"
        - "./artifacts"
```