# Executors

Executors are the components responsible for running steps in your Plexr plans. They provide the interface between your plan definitions and the actual execution environment.

## Built-in Executors

### Shell Executor

The default executor for running shell commands.

```yaml
steps:
  - name: "Install dependencies"
    command: "npm install"
    executor: shell  # Optional, shell is the default
```

#### Configuration

```yaml
steps:
  - name: "Build project"
    command: "make build"
    executor: shell
    env:
      BUILD_ENV: production
    timeout: 300  # seconds
    workdir: ./src
```

### SQL Executor

Execute SQL queries against PostgreSQL databases (MySQL and SQLite support planned).

#### Configuration

Define SQL executors in the `executors` section:

```yaml
executors:
  db:
    type: sql
    driver: postgres
    host: ${DB_HOST:-localhost}
    port: ${DB_PORT:-5432}
    database: ${DB_NAME:-myapp}
    username: ${DB_USER:-postgres}
    password: ${DB_PASSWORD}
    sslmode: ${DB_SSLMODE:-disable}
```

#### Usage

```yaml
steps:
  - name: "Run migrations"
    executor: db
    files:
      - path: sql/001_schema.sql
        timeout: 30
    transaction_mode: all  # Options: none, each, all
```

#### Transaction Modes

- `none`: No transaction wrapping
- `each`: Each SQL statement in its own transaction (default)
- `all`: All statements in a single transaction

#### Multiple Databases

You can define multiple SQL executors for different databases:

```yaml
executors:
  main_db:
    type: sql
    driver: postgres
    database: myapp_main
    # ... connection details
    
  analytics_db:
    type: sql
    driver: postgres
    database: myapp_analytics
    # ... connection details
```

## Custom Executors

You can create custom executors by implementing the `Executor` interface:

```go
type Executor interface {
    Execute(ctx context.Context, step Step, state State) error
    Validate(step Step) error
}
```

### Example Custom Executor

```go
package executors

import (
    "context"
    "github.com/plexr/plexr/internal/core"
)

type DockerExecutor struct {
    client *docker.Client
}

func (e *DockerExecutor) Execute(ctx context.Context, step core.Step, state core.State) error {
    // Implementation for running Docker containers
    container, err := e.client.CreateContainer(docker.CreateContainerOptions{
        Config: &docker.Config{
            Image: step.Config["image"],
            Cmd:   []string{step.Command},
        },
    })
    if err != nil {
        return err
    }
    
    return e.client.StartContainer(container.ID, nil)
}

func (e *DockerExecutor) Validate(step core.Step) error {
    if step.Config["image"] == "" {
        return fmt.Errorf("docker executor requires 'image' in config")
    }
    return nil
}
```

## Executor Configuration

### Global Configuration

Set default executor configurations in your plan:

```yaml
config:
  executors:
    shell:
      timeout: 60
      shell: /bin/bash
    docker:
      registry: ghcr.io
      
steps:
  - name: "Build"
    command: "make build"
    # Uses global shell config
```

### Per-Step Configuration

Override executor settings for specific steps:

```yaml
steps:
  - name: "Long running task"
    command: "./long-task.sh"
    executor: shell
    config:
      timeout: 3600  # 1 hour
```

## State Management

Executors can read and write to the shared state:

```yaml
steps:
  - name: "Get version"
    command: "git describe --tags"
    executor: shell
    outputs:
      - name: version
        from: stdout
        
  - name: "Build with version"
    command: "make build VERSION=`{{.version}}`"
```

## Error Handling

### Retry Configuration

```yaml
steps:
  - name: "Flaky test"
    command: "npm test"
    retry:
      attempts: 3
      delay: 5s
      backoff: exponential
```

### Error Conditions

```yaml
steps:
  - name: "Check service"
    command: "curl http://localhost:8080/health"
    errorConditions:
      - exitCode: [1, 2]
        retry: true
      - output: "connection refused"
        fail: true
```

## Parallel Execution

Run multiple steps concurrently:

```yaml
steps:
  - name: "Parallel tasks"
    parallel:
      - name: "Test frontend"
        command: "npm test"
      - name: "Test backend"
        command: "go test ./..."
      - name: "Lint"
        command: "npm run lint"
```

## Environment Variables

Executors support environment variable expansion:

```yaml
steps:
  - name: "Deploy"
    command: "deploy.sh"
    env:
      DEPLOY_ENV: "`{{.environment}}`"
      API_KEY: "$DEPLOY_API_KEY"  # From host environment
```

## Logging and Output

### Output Capture

```yaml
steps:
  - name: "Generate report"
    command: "./generate-report.sh"
    outputs:
      - name: reportPath
        from: stdout
        regex: "Report saved to: (.*)"
```

### Log Levels

```yaml
steps:
  - name: "Verbose operation"
    command: "npm install"
    logLevel: debug  # debug, info, warn, error
```

## Security Considerations

- Executors run with the permissions of the Plexr process
- Use environment variables for sensitive data
- Validate all inputs in custom executors
- Consider using sandboxed environments for untrusted code

## Best Practices

1. **Use specific executors**: Choose the right executor for each task
2. **Handle errors gracefully**: Configure appropriate retry and error handling
3. **Manage state carefully**: Use outputs and conditions for complex workflows
4. **Monitor execution**: Set appropriate timeouts and logging
5. **Test executors**: Validate custom executors thoroughly before production use