# State Management

Plexr provides powerful state management capabilities to share data between steps, track execution progress, and maintain workflow context.

## Overview

State in Plexr is a key-value store that persists throughout the execution of a plan. It allows steps to:
- Share data with subsequent steps
- Store execution results
- Make conditional decisions
- Track workflow progress

## State File

Plexr stores state in `.plexr/state.json` in your project directory:

```json
{
  "variables": {
    "version": "1.2.3",
    "build_id": "abc123",
    "environment": "production"
  },
  "steps": {
    "build": {
      "status": "completed",
      "start_time": "2024-01-20T10:00:00Z",
      "end_time": "2024-01-20T10:05:00Z",
      "outputs": {
        "artifact_path": "/dist/app.tar.gz"
      }
    }
  }
}
```

## Setting State Variables

### From Command Output

Capture command output into state variables:

```yaml
steps:
  - name: "Get version"
    command: "git describe --tags --always"
    outputs:
      - name: version
        from: stdout
```

### From JSON Output

Parse JSON output and extract specific fields:

```yaml
steps:
  - name: "Get build info"
    command: "npm run build --json"
    outputs:
      - name: build_id
        from: stdout
        json_path: "$.buildId"
      - name: build_time
        from: stdout
        json_path: "$.timestamp"
```

### Using Regular Expressions

Extract data using regex patterns:

```yaml
steps:
  - name: "Parse deployment URL"
    command: "deploy.sh"
    outputs:
      - name: app_url
        from: stdout
        regex: "Deployed to: (https://.*)"
        regex_group: 1
```

### Multiple Outputs

Capture multiple values from a single step:

```yaml
steps:
  - name: "System info"
    command: |
      echo "OS: $(uname -s)"
      echo "Arch: $(uname -m)"
      echo "Host: $(hostname)"
    outputs:
      - name: os
        from: stdout
        regex: "OS: (.*)"
      - name: arch
        from: stdout
        regex: "Arch: (.*)"
      - name: hostname
        from: stdout
        regex: "Host: (.*)"
```

## Using State Variables

### Variable Substitution

Reference state variables using Go template syntax:

```yaml
steps:
  - name: "Build with version"
    command: "docker build -t myapp:`{{.version}}` ."
    
  - name: "Deploy"
    command: "kubectl set image deployment/app app=myapp:`{{.version}}`"
```

### Default Values

Provide defaults for missing variables:

```yaml
steps:
  - name: "Configure environment"
    command: "setup.sh"
    env:
      ENVIRONMENT: "`{{.environment | default \"development\"}}`"
      DEBUG: "`{{.debug | default \"false\"}}`"
```

### Conditional Logic

Use state variables in conditions:

```yaml
steps:
  - name: "Production deployment"
    command: "deploy-prod.sh"
    condition: '`{{.environment}}` == "production" && `{{.tests_passed}}` == "true"'
```

## Advanced State Management

### Nested Variables

Work with complex data structures:

```yaml
steps:
  - name: "Get config"
    command: "cat config.json"
    outputs:
      - name: config
        from: stdout
        json_parse: true
        
  - name: "Use nested config"
    command: "connect.sh"
    env:
      DB_HOST: "`{{.config.database.host}}`"
      DB_PORT: "`{{.config.database.port}}`"
```

### Arrays and Loops

Handle array data:

```yaml
steps:
  - name: "Get services"
    command: "kubectl get services -o json"
    outputs:
      - name: services
        from: stdout
        json_path: "$.items[*].metadata.name"
        
  - name: "Process services"
    command: "check-service.sh `{{.service}}`"
    for_each:
      items: "`{{.services}}`"
      var: service
```

### State Mutations

Transform state variables:

```yaml
steps:
  - name: "Increment version"
    command: |
      current=`{{.build_number | default "0"}}`
      echo $((current + 1))
    outputs:
      - name: build_number
        from: stdout
```

## State Persistence

### Saving State

State is automatically saved after each step. You can also manually checkpoint:

```yaml
steps:
  - name: "Critical operation"
    command: "important-task.sh"
    save_state: true  # Force immediate state save
```

### Loading Previous State

Continue from a previous execution:

```bash
# Resume from last state
plexr execute --resume

# Load specific state file
plexr execute --state-file backup-state.json
```

### Resetting State

Clear all state data:

```bash
# Reset state for current plan
plexr reset

# Reset and remove state file
plexr reset --clean
```

## State Scoping

### Global Variables

Set variables available to all steps:

```yaml
vars:
  project_name: "myapp"
  region: "us-east-1"
  
steps:
  - name: "Deploy"
    command: "deploy.sh `{{.project_name}}` `{{.region}}`"
```

### Step-Local Variables

Variables scoped to specific steps:

```yaml
steps:
  - name: "Build variants"
    parallel:
      - name: "Build debug"
        command: "build.sh"
        env:
          BUILD_TYPE: "debug"
      - name: "Build release"
        command: "build.sh"
        env:
          BUILD_TYPE: "release"
```

### Environment Variable Integration

Mix environment and state variables:

```yaml
steps:
  - name: "Deploy"
    command: "deploy.sh"
    env:
      APP_VERSION: "`{{.version}}`"  # From state
      API_KEY: "${DEPLOY_API_KEY}"  # From environment
```

## State Templates

### String Manipulation

```yaml
steps:
  - name: "Format message"
    command: |
      echo "`{{.message | upper}}`"
      echo "`{{.path | base}}`"
      echo "`{{.text | replace \" \" \"_\"}}`"
```

### Arithmetic Operations

```yaml
steps:
  - name: "Calculate"
    command: |
      echo "Total: `{{add .value1 .value2}}`"
      echo "Difference: `{{sub .value1 .value2}}`"
```

### Date and Time

```yaml
steps:
  - name: "Timestamp"
    command: "echo `{{now | date \"2006-01-02 15:04:05\"}}`"
    outputs:
      - name: build_time
        from: stdout
```

## Best Practices

### 1. Variable Naming

Use clear, descriptive names:
```yaml
# Good
outputs:
  - name: docker_image_tag
  - name: deployment_url
  - name: test_coverage_percent

# Avoid
outputs:
  - name: tag
  - name: url
  - name: coverage
```

### 2. State Validation

Validate state before using:
```yaml
steps:
  - name: "Check prerequisites"
    command: |
      if [ -z "`{{.api_key}}`" ]; then
        echo "Error: api_key not set"
        exit 1
      fi
```

### 3. State Documentation

Document expected state variables:
```yaml
# This plan expects the following variables:
# - environment: Target environment (dev, staging, prod)
# - version: Application version to deploy
# - region: AWS region for deployment

requires:
  - environment
  - version
  - region
```

### 4. State Cleanup

Clean sensitive data:
```yaml
steps:
  - name: "Cleanup secrets"
    command: "echo ''"
    outputs:
      - name: api_key
        value: ""  # Clear sensitive data
    always_run: true
```

## Debugging State

### View Current State

```bash
# Show all state
plexr status

# Show specific variable
plexr status --var version

# Export state as JSON
plexr status --json > state-backup.json
```

### State History

Track state changes:
```yaml
steps:
  - name: "Log state change"
    command: |
      echo "Previous: `{{.version | default \"none\"}}`"
      echo "New: `{{.new_version}}`"
    debug: true
```

### Interactive State Updates

Modify state during development:
```bash
# Set variable
plexr state set version "1.2.3"

# Remove variable
plexr state unset debug_mode

# Import state
plexr state import < custom-state.json
```

## Common Patterns

### Feature Flags

```yaml
vars:
  features:
    new_ui: true
    beta_api: false
    
steps:
  - name: "Deploy with features"
    command: "deploy.sh"
    env:
      ENABLE_NEW_UI: "`{{.features.new_ui}}`"
      ENABLE_BETA_API: "`{{.features.beta_api}}`"
```

### Build Matrix

```yaml
vars:
  platforms: ["linux", "darwin", "windows"]
  architectures: ["amd64", "arm64"]
  
steps:
  - name: "Build matrix"
    command: "build.sh -os `{{.platform}}` -arch `{{.arch}}`"
    for_each:
      platforms: "`{{.platforms}}`"
      architectures: "`{{.architectures}}`"
      as:
        platform: platform
        arch: arch
```

### Workflow State Machine

```yaml
steps:
  - name: "Check state"
    command: "get-workflow-state.sh"
    outputs:
      - name: workflow_state
        
  - name: "Process pending"
    command: "process.sh"
    condition: '`{{.workflow_state}}` == "pending"'
    
  - name: "Handle error"
    command: "error-handler.sh"
    condition: '`{{.workflow_state}}` == "error"'
```