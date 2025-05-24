# Frequently Asked Questions

## General Questions

### What is Plexr?

Plexr is a command execution and workflow automation tool that helps you define, manage, and execute complex multi-step processes. It's designed for development workflows, CI/CD pipelines, and infrastructure automation.

### How is Plexr different from Make?

While Make is great for building software, Plexr offers:
- **Better state management**: Share data between steps easily
- **Modern YAML syntax**: More readable and maintainable
- **Built-in parallelism**: Run steps concurrently without complex syntax
- **Conditional execution**: Dynamic workflows based on state
- **Multiple executors**: Not limited to shell commands

### Can I use Plexr in CI/CD?

Yes! Plexr is designed to work well in CI/CD environments:
- Consistent execution across environments
- Clear dependency management
- Detailed logging and error reporting
- State persistence for complex workflows

## Installation

### How do I install Plexr?

```bash
# Using Go
go install github.com/plexr/plexr/cmd/plexr@latest

# From source
git clone https://github.com/plexr/plexr
cd plexr
make install
```

### What are the system requirements?

- Go 1.19 or later (for installation from source)
- Linux, macOS, or Windows
- Shell access (bash, sh, or equivalent)

### How do I update Plexr?

```bash
# Using Go
go install github.com/plexr/plexr/cmd/plexr@latest

# Check version
plexr version
```

## Configuration

### Where should I put my plan file?

By default, Plexr looks for `plan.yml` or `plan.yaml` in the current directory. You can also specify a custom location:

```bash
plexr execute -f path/to/myplan.yml
```

### Can I have multiple plan files?

Yes! You can organize your workflows:

```bash
# Different plans for different purposes
plexr execute -f plans/build.yml
plexr execute -f plans/deploy.yml
plexr execute -f plans/test.yml
```

### How do I pass variables to my plan?

Several ways:

```bash
# Environment variables
export API_KEY=secret
plexr execute

# Command-line variables
plexr execute -v environment=production -v version=1.2.3

# From file
plexr execute --vars-file vars.json
```

## Execution

### Can I run specific steps only?

Yes, use tags or specify steps:

```yaml
steps:
  - name: "Build"
    command: "make build"
    tags: ["build", "ci"]
    
  - name: "Test"
    command: "make test"
    tags: ["test", "ci"]
```

```bash
# Run only tagged steps
plexr execute --tags build

# Run from specific step
plexr execute --from "Test"
```

### How do I debug execution problems?

Enable verbose logging:

```bash
# Basic debug info
plexr execute -v

# Detailed debug info
plexr execute -vv

# Dry run (show what would be executed)
plexr execute --dry-run
```

### Can I resume a failed execution?

Yes:

```bash
# Resume from last successful step
plexr execute --resume

# Retry failed steps
plexr execute --retry-failed
```

## State Management

### Where is state stored?

State is stored in `.plexr/state.json` in your project directory. You can customize this:

```bash
# Use custom state file
plexr execute --state-file /tmp/mystate.json

# Use in-memory state (no persistence)
plexr execute --no-state-file
```

### How do I share data between steps?

Use outputs to capture data:

```yaml
steps:
  - name: "Get version"
    command: "git describe --tags"
    outputs:
      - name: version
        from: stdout
        
  - name: "Use version"
    command: "echo Building version `{{.version}}`"
```

### Can I use external data sources?

Yes, load data from files or commands:

```yaml
vars:
  config: "`{{file \"config.json\" | json}}`"
  
steps:
  - name: "Load secrets"
    command: "vault read -format=json secret/myapp"
    outputs:
      - name: secrets
        from: stdout
        json_parse: true
```

## Security

### How do I handle secrets?

Use environment variables and avoid hardcoding:

```yaml
steps:
  - name: "Deploy"
    command: "deploy.sh"
    env:
      API_KEY: "${API_KEY}"  # From environment
      DB_PASS: "`{{.vault_secret}}`"  # From previous step
```

### Can I mask sensitive output?

Yes:

```yaml
steps:
  - name: "Show config"
    command: "print-config.sh"
    mask_output: ["password", "api_key", "secret"]
```

### Is it safe to commit plan files?

Yes, if you follow best practices:
- Never hardcode secrets
- Use environment variables
- Document required variables
- Use `.gitignore` for state files

## Advanced Usage

### Can I create custom executors?

Yes, implement the Executor interface:

```go
type Executor interface {
    Execute(ctx context.Context, step Step, state State) error
    Validate(step Step) error
}
```

See the [Executors Guide](./executors.md) for details.

### How do I handle complex dependencies?

Use dependency groups and conditions:

```yaml
steps:
  - name: "Prepare"
    command: "prepare.sh"
    
  - name: "Build A"
    command: "build-a.sh"
    depends_on: ["Prepare"]
    
  - name: "Build B"
    command: "build-b.sh"
    depends_on: ["Prepare"]
    
  - name: "Package"
    command: "package.sh"
    depends_on: ["Build A", "Build B"]
```

### Can I use Plexr as a library?

Yes:

```go
import (
    "github.com/plexr/plexr/internal/config"
    "github.com/plexr/plexr/internal/core"
)

plan, err := config.LoadPlan("plan.yml")
runner := core.NewRunner()
err = runner.Execute(ctx, plan)
```

## Integration

### Does Plexr work with Docker?

Yes, you can run commands in containers:

```yaml
steps:
  - name: "Build in Docker"
    command: |
      docker run --rm -v $(pwd):/app -w /app \
        node:16 npm run build
```

### Can I integrate with Kubernetes?

Yes:

```yaml
steps:
  - name: "Deploy to K8s"
    command: |
      kubectl apply -f deployment.yaml
      kubectl wait --for=condition=available --timeout=300s \
        deployment/myapp
```

### How do I use Plexr with Git hooks?

Create a `.git/hooks/pre-commit` file:

```bash
#!/bin/bash
plexr execute -f .plexr/pre-commit.yml
```

## Troubleshooting

### Command not found errors

Ensure commands are in PATH or use absolute paths:

```yaml
steps:
  - name: "Use absolute path"
    command: "/usr/local/bin/mytool"
    
  - name: "Update PATH"
    command: "mytool"
    env:
      PATH: "/usr/local/bin:${PATH}"
```

### State corruption

Reset state if corrupted:

```bash
# Reset state
plexr reset

# Remove all Plexr data
rm -rf .plexr
```

### Performance issues

- Use parallel execution for independent steps
- Limit output capture for verbose commands
- Use appropriate timeouts
- Consider breaking large plans into smaller ones

## Best Practices

### How should I structure large projects?

```
project/
├── plexr/
│   ├── build.yml
│   ├── test.yml
│   ├── deploy.yml
│   └── common/
│       ├── vars.yml
│       └── functions.yml
├── scripts/
│   └── ...
└── plan.yml  # Main orchestration
```

### Should I version control state files?

No, add to `.gitignore`:

```gitignore
.plexr/
*.state.json
```

### How do I make plans reusable?

Use variables and includes:

```yaml
# common/database.yml
steps:
  - name: "Database setup"
    command: "`{{.db_setup_script}}`"
    env:
      DB_NAME: "`{{.database_name}}`"

# main plan.yml
includes:
  - common/database.yml
  
vars:
  database_name: "myapp"
  db_setup_script: "./scripts/setup-db.sh"
```

## Getting Help

### Where can I get support?

- GitHub Issues: https://github.com/plexr/plexr/issues
- Documentation: https://plexr.dev/docs
- Community Discord: https://discord.gg/plexr

### How do I report bugs?

Include:
1. Plexr version (`plexr version`)
2. Plan file (sanitized)
3. Error messages
4. Steps to reproduce
5. Expected vs actual behavior

### Can I contribute?

Yes! We welcome contributions:
- Report bugs
- Suggest features
- Submit pull requests
- Improve documentation
- Share your use cases