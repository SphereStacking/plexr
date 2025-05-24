# Troubleshooting

This guide helps you diagnose and fix common issues when using Plexr.

## Common Issues

### Installation Problems

#### Go module errors

**Problem**: `go: module github.com/plexr/plexr: git ls-remote -q origin` fails

**Solution**:
```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download
```

#### Permission denied during installation

**Problem**: `permission denied` when running `go install`

**Solution**:
```bash
# Install to user's bin directory
go install github.com/plexr/plexr/cmd/plexr@latest

# Or use sudo (not recommended)
sudo go install github.com/plexr/plexr/cmd/plexr@latest
```

### Configuration Issues

#### Plan file not found

**Problem**: `Error: plan file not found: plan.yml`

**Solution**:
1. Ensure you're in the correct directory
2. Check file name spelling (plan.yml vs plan.yaml)
3. Use the `-f` flag to specify a custom path:
   ```bash
   plexr execute -f path/to/myplan.yml
   ```

#### Invalid YAML syntax

**Problem**: `Error: failed to parse plan: yaml: line X`

**Solution**:
1. Validate your YAML syntax using an online validator
2. Check for common issues:
   - Incorrect indentation (use spaces, not tabs)
   - Missing colons after keys
   - Unquoted special characters

Example of correct syntax:
```yaml
version: "1.0"
name: "My Plan"
steps:
  - name: "Step 1"  # Note the quotes and proper indentation
    command: "echo hello"
```

### Execution Problems

#### Step fails with "command not found"

**Problem**: `bash: command: command not found`

**Solution**:
1. Ensure the command is installed and in PATH
2. Use absolute paths for custom scripts:
   ```yaml
   steps:
     - name: "Run script"
       command: "./scripts/my-script.sh"  # Relative path
       # or
       command: "/home/user/project/scripts/my-script.sh"  # Absolute path
   ```

#### Timeout errors

**Problem**: `Error: step timed out after 60s`

**Solution**:
```yaml
steps:
  - name: "Long running task"
    command: "./long-task.sh"
    timeout: 300  # Increase timeout to 5 minutes
```

#### Working directory issues

**Problem**: `Error: no such file or directory`

**Solution**:
```yaml
steps:
  - name: "Build in subdirectory"
    command: "make build"
    workdir: "./src"  # Specify working directory
```

### State Management Issues

#### Variable not found

**Problem**: `Error: variable 'myvar' not found in state`

**Solution**:
1. Ensure the variable is set before use:
   ```yaml
   steps:
     - name: "Set variable"
       command: "echo 'value'"
       outputs:
         - name: myvar
           from: stdout
     
     - name: "Use variable"
       command: "echo `{{.myvar}}`"  # Now it exists
   ```

#### State file corruption

**Problem**: `Error: failed to load state: invalid character`

**Solution**:
```bash
# Reset state
plexr reset

# Or manually remove state file
rm .plexr/state.json
```

### Dependency Issues

#### Circular dependencies

**Problem**: `Error: circular dependency detected`

**Solution**:
Review your step dependencies and remove cycles:
```yaml
# Bad - circular dependency
steps:
  - name: "A"
    depends_on: ["B"]
  - name: "B"
    depends_on: ["A"]

# Good - no cycles
steps:
  - name: "A"
  - name: "B"
    depends_on: ["A"]
```

#### Step skipped unexpectedly

**Problem**: Step doesn't run even though dependencies are met

**Solution**:
Check conditions:
```yaml
steps:
  - name: "Conditional step"
    command: "deploy.sh"
    condition: "`{{.environment}}` == 'production'"
    # This won't run unless environment is set to "production"
```

## Debugging Tips

### Enable verbose logging

```bash
# Show detailed execution information
plexr execute -v

# Show even more details
plexr execute -vv
```

### Check execution plan

```bash
# Validate without executing
plexr validate

# Show execution order
plexr status --show-order
```

### Inspect state

```bash
# View current state
plexr status

# View state file directly
cat .plexr/state.json | jq
```

### Test individual steps

```yaml
# Add a test mode to your steps
steps:
  - name: "Deploy"
    command: |
      if [ "$TEST_MODE" = "true" ]; then
        echo "Would deploy to production"
      else
        ./deploy.sh
      fi
    env:
      TEST_MODE: "`{{.test_mode | default \"false\"}}`"
```

## Environment-Specific Issues

### Docker container issues

**Problem**: Steps fail when running in Docker

**Solution**:
1. Mount necessary volumes:
   ```dockerfile
   docker run -v $(pwd):/workspace plexr execute
   ```

2. Set working directory:
   ```yaml
   steps:
     - name: "Build in container"
       command: "make build"
       workdir: "/workspace"
   ```

### CI/CD pipeline failures

**Problem**: Plexr works locally but fails in CI

**Solution**:
1. Check environment variables are set
2. Ensure all dependencies are installed
3. Use explicit paths and versions:
   ```yaml
   steps:
     - name: "CI Build"
       command: "/usr/local/go/bin/go build"
       env:
         GOPATH: "/home/runner/go"
   ```

## Performance Issues

### Slow execution

**Solutions**:
1. Run independent steps in parallel:
   ```yaml
   steps:
     - name: "Parallel tasks"
       parallel:
         - name: "Test 1"
           command: "test1.sh"
         - name: "Test 2"
           command: "test2.sh"
   ```

2. Cache dependencies:
   ```yaml
   steps:
     - name: "Install deps"
       command: "npm ci"
       condition: "`{{.deps_cached}}` != 'true'"
   ```

### High memory usage

**Solution**:
Limit output capture for large commands:
```yaml
steps:
  - name: "Large output"
    command: "find / -name '*.log'"
    capture_output: false  # Don't store in state
```

## Getting Help

### Collect diagnostic information

```bash
# System information
plexr version
go version
uname -a

# Plan validation
plexr validate -v

# Execution trace
plexr execute -vv 2>&1 | tee plexr-debug.log
```

### Report issues

When reporting issues, include:
1. Plexr version (`plexr version`)
2. Plan file (sanitized of sensitive data)
3. Error messages and logs
4. Steps to reproduce

Report issues at: https://github.com/plexr/plexr/issues

## FAQ

**Q: Can I use environment variables in plan files?**
A: Yes, use `${VAR_NAME}` for environment variables and <code>&#123;&#123;.var_name&#125;&#125;</code> for state variables.

**Q: How do I handle secrets?**
A: Use environment variables and never commit secrets to your plan files.

**Q: Can I run Plexr in the background?**
A: Yes: `nohup plexr execute > plexr.log 2>&1 &`

**Q: How do I update Plexr?**
A: Run `go install github.com/plexr/plexr/cmd/plexr@latest`