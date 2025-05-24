# Advanced Patterns

This guide demonstrates advanced Plexr patterns for complex workflows, conditional logic, and sophisticated automation scenarios.

## Complex Dependency Chains

### Diamond Dependencies

Handle complex dependency graphs where multiple paths converge:

```yaml
name: "Diamond dependency pattern"
version: "1.0"

steps:
  - name: "Start"
    command: "echo 'Starting workflow'"
    
  - name: "Path A1"
    command: "process-a1.sh"
    depends_on: ["Start"]
    
  - name: "Path A2"
    command: "process-a2.sh"
    depends_on: ["Path A1"]
    
  - name: "Path B1"
    command: "process-b1.sh"
    depends_on: ["Start"]
    
  - name: "Path B2"
    command: "process-b2.sh"
    depends_on: ["Path B1"]
    
  - name: "Converge"
    command: "merge-results.sh"
    depends_on: ["Path A2", "Path B2"]
```

### Conditional Dependencies

Execute steps based on the success of previous steps:

```yaml
steps:
  - name: "Check environment"
    command: "check-env.sh"
    outputs:
      - name: env_ready
        from: stdout
        
  - name: "Setup environment"
    command: "setup-env.sh"
    condition: '`{{.env_ready}}` != "true"'
    outputs:
      - name: env_ready
        value: "true"
        
  - name: "Deploy"
    command: "deploy.sh"
    depends_on: ["Check environment", "Setup environment"]
    condition: '`{{.env_ready}}` == "true"'
```

### Dynamic Dependencies

Generate dependencies based on discovered resources:

```yaml
steps:
  - name: "Discover services"
    command: "kubectl get services -o json | jq -r '.items[].metadata.name'"
    outputs:
      - name: services
        from: stdout
        as_array: true
        
  - name: "Generate health checks"
    command: |
      for service in `{{range .services}}``{{.}}` `{{end}}`; do
        echo "- name: \"Check $service\""
        echo "  command: \"check-health.sh $service\""
        echo "  tags: [\"health-check\"]"
      done > dynamic-steps.yml
      
  - name: "Run health checks"
    include: "dynamic-steps.yml"
    depends_on: ["Generate health checks"]
```

## Conditional Execution

### Multi-Condition Logic

Combine multiple conditions:

```yaml
steps:
  - name: "Production deployment"
    command: "deploy-prod.sh"
    condition: |
      `{{.environment}}` == "production" &&
      `{{.tests_passed}}` == "true" &&
      `{{.approval_granted}}` == "true"
      
  - name: "Staging deployment"
    command: "deploy-staging.sh"
    condition: |
      `{{.environment}}` == "staging" &&
      `{{.tests_passed}}` == "true"
```

### Switch-Case Pattern

Implement switch-like behavior:

```yaml
vars:
  deployment_commands:
    dev: "deploy-dev.sh"
    staging: "deploy-staging.sh"
    prod: "deploy-prod.sh"
    
steps:
  - name: "Deploy to environment"
    command: "`{{index .deployment_commands .environment}}`"
    condition: "`{{hasKey .deployment_commands .environment}}`"
    
  - name: "Unknown environment"
    command: "echo 'Unknown environment: `{{.environment}}`' && exit 1"
    condition: "`{{not (hasKey .deployment_commands .environment)}}`"
```

### Conditional Includes

Include different plans based on conditions:

```yaml
steps:
  - name: "Detect OS"
    command: "uname -s"
    outputs:
      - name: os
        from: stdout
        
  - name: "Linux setup"
    include: "setup-linux.yml"
    condition: '`{{.os}}` == "Linux"'
    
  - name: "macOS setup"
    include: "setup-macos.yml"
    condition: '`{{.os}}` == "Darwin"'
    
  - name: "Windows setup"
    include: "setup-windows.yml"
    condition: '`{{contains .os "MINGW"}}`'
```

## Error Handling Strategies

### Graceful Degradation

Continue execution with reduced functionality:

```yaml
steps:
  - name: "Try primary database"
    command: "connect-primary-db.sh"
    ignore_failure: true
    outputs:
      - name: primary_db_available
        from: exit_code
        transform: '`{{if eq . "0"}}`true`{{else}}`false`{{end}}`'
        
  - name: "Use fallback database"
    command: "connect-fallback-db.sh"
    condition: '`{{.primary_db_available}}` != "true"'
    
  - name: "Run with available database"
    command: "run-app.sh"
    env:
      USE_FALLBACK: '`{{.primary_db_available}}` != "true"'
```

### Retry with Backoff

Implement sophisticated retry strategies:

```yaml
steps:
  - name: "API call with retry"
    command: "call-api.sh"
    retry:
      attempts: 5
      delay: 2s
      backoff: exponential
      max_delay: 60s
      on_retry: |
        echo "Retry attempt `{{.retry_count}}` of `{{.retry_max}}`"
        echo "Next retry in `{{.retry_delay}}` seconds"
```

### Circuit Breaker Pattern

Prevent cascading failures:

```yaml
vars:
  circuit_breaker:
    failure_threshold: 3
    timeout: 300  # 5 minutes
    
steps:
  - name: "Check circuit state"
    command: |
      if [ -f .circuit_open ] && [ $(( $(date +%s) - $(stat -f %m .circuit_open) )) -lt `{{.circuit_breaker.timeout}}` ]; then
        echo "open"
      else
        echo "closed"
      fi
    outputs:
      - name: circuit_state
        from: stdout
        
  - name: "Service call"
    command: "call-service.sh"
    condition: '`{{.circuit_state}}` != "open"'
    on_failure:
      - name: "Increment failure count"
        command: |
          count=$(cat .failure_count 2>/dev/null || echo 0)
          echo $((count + 1)) > .failure_count
          if [ $((count + 1)) -ge `{{.circuit_breaker.failure_threshold}}` ]; then
            touch .circuit_open
          fi
```

## State Management Patterns

### State Machines

Implement workflow state machines:

```yaml
vars:
  states:
    init: "initializing"
    build: "building"
    test: "testing"
    deploy: "deploying"
    complete: "completed"
    failed: "failed"
    
steps:
  - name: "Set initial state"
    command: "echo '`{{.states.init}}`'"
    outputs:
      - name: workflow_state
        from: stdout
        
  - name: "Build phase"
    command: "build.sh"
    condition: '`{{.workflow_state}}` == "`{{.states.init}}`"'
    outputs:
      - name: workflow_state
        value: "`{{.states.build}}`"
    on_failure:
      outputs:
        - name: workflow_state
          value: "`{{.states.failed}}`"
          
  - name: "Test phase"
    command: "test.sh"
    condition: '`{{.workflow_state}}` == "`{{.states.build}}`"'
    outputs:
      - name: workflow_state
        value: "`{{.states.test}}`"
```

### Persistent State Across Runs

Maintain state between executions:

```yaml
steps:
  - name: "Load persistent state"
    command: |
      if [ -f state.json ]; then
        cat state.json
      else
        echo '{}'
      fi
    outputs:
      - name: persistent_state
        from: stdout
        json_parse: true
        
  - name: "Update state"
    command: |
      echo '`{{.persistent_state | toJson}}`' | \
      jq '.last_run = "`{{now | date "2006-01-02T15:04:05Z07:00"}}`"' | \
      jq '.run_count = ((.run_count // 0) + 1)'
    outputs:
      - name: updated_state
        from: stdout
        json_parse: true
        
  - name: "Save state"
    command: "echo '`{{.updated_state | toJson}}`' > state.json"
    always_run: true
```

### Distributed State

Share state across multiple plan executions:

```yaml
steps:
  - name: "Acquire lock"
    command: |
      while ! mkdir .lock 2>/dev/null; do
        echo "Waiting for lock..."
        sleep 1
      done
    timeout: 60
    
  - name: "Read shared state"
    command: "cat shared-state.json 2>/dev/null || echo '{}'"
    outputs:
      - name: shared_state
        from: stdout
        json_parse: true
        
  - name: "Update shared state"
    command: |
      echo '`{{.shared_state | toJson}}`' | \
      jq '.workers["`{{.worker_id}}`"] = {
        "last_seen": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
        "status": "active"
      }' > shared-state.json
      
  - name: "Release lock"
    command: "rmdir .lock"
    always_run: true
```

## Dynamic Workflow Generation

### Template-Based Generation

Generate steps from templates:

```yaml
vars:
  environments: ["dev", "staging", "prod"]
  regions: ["us-east-1", "eu-west-1", "ap-southeast-1"]
  
steps:
  - name: "Generate deployment steps"
    command: |
      cat > generated-steps.yml << EOF
      steps:
      `{{range $env := .environments}}`
      `{{range $region := $.regions}}`
        - name: "Deploy to `{{$env}}`-`{{$region}}`"
          command: "deploy.sh"
          env:
            ENVIRONMENT: "`{{$env}}`"
            REGION: "`{{$region}}`"
          tags: ["deploy", "`{{$env}}`", "`{{$region}}`"]
      `{{end}}`
      `{{end}}`
      EOF
      
  - name: "Execute deployments"
    include: "generated-steps.yml"
    depends_on: ["Generate deployment steps"]
```

### Discovery-Based Workflow

Build workflows based on discovered resources:

```yaml
steps:
  - name: "Discover microservices"
    command: "find services -name 'Dockerfile' -type f | xargs dirname"
    outputs:
      - name: services
        from: stdout
        as_array: true
        
  - name: "Build service images"
    command: |
      docker build -t `{{.service | base}}`:`{{.version}}` `{{.service}}`
    for_each:
      items: "`{{.services}}`"
      var: service
    parallel:
      max_concurrent: 3
```

## Integration Patterns

### Event-Driven Workflows

React to external events:

```yaml
steps:
  - name: "Wait for webhook"
    command: |
      # Simple webhook listener
      while true; do
        nc -l 8080 | grep -q "deploy-signal" && break
      done
    timeout: 3600  # 1 hour
    
  - name: "Process webhook data"
    command: "parse-webhook.sh"
    outputs:
      - name: deployment_target
        from: stdout
        
  - name: "Deploy based on webhook"
    command: "deploy.sh `{{.deployment_target}}`"
```

### Polling Pattern

Check for conditions periodically:

```yaml
steps:
  - name: "Wait for service ready"
    command: |
      for i in {1..60}; do
        if curl -s http://service:8080/health | grep -q "ok"; then
          echo "ready"
          exit 0
        fi
        echo "Attempt $i/60: Service not ready, waiting..."
        sleep 5
      done
      echo "timeout"
      exit 1
    outputs:
      - name: service_status
        from: stdout
```

### Message Queue Integration

Process items from a queue:

```yaml
steps:
  - name: "Get queue messages"
    command: "aws sqs receive-message --queue-url `{{.queue_url}}` --max-number-of-messages 10"
    outputs:
      - name: messages
        from: stdout
        json_path: "$.Messages"
        
  - name: "Process messages"
    command: |
      echo '`{{.message | toJson}}`' | process-message.sh
    for_each:
      items: "`{{.messages}}`"
      var: message
    on_success:
      - name: "Delete message"
        command: |
          aws sqs delete-message \
            --queue-url `{{.queue_url}}` \
            --receipt-handle `{{.message.ReceiptHandle}}`
```

## Performance Optimization

### Parallel Matrix Execution

Run combinations in parallel:

```yaml
vars:
  python_versions: ["3.8", "3.9", "3.10", "3.11"]
  test_suites: ["unit", "integration", "e2e"]
  
steps:
  - name: "Test matrix"
    parallel:
      - name: "Test Python `{{.version}}` - `{{.suite}}`"
        command: |
          pyenv local `{{.version}}`
          pytest tests/`{{.suite}}`
        for_each:
          python_version: "`{{.python_versions}}`"
          test_suite: "`{{.test_suites}}`"
          as:
            version: python_version
            suite: test_suite
      max_concurrent: 4
```

### Caching Pattern

Implement intelligent caching:

```yaml
steps:
  - name: "Calculate cache key"
    command: |
      echo "$(cat package-lock.json | sha256sum | cut -d' ' -f1)-$(date +%Y%m%d)"
    outputs:
      - name: cache_key
        from: stdout
        
  - name: "Check cache"
    command: |
      if [ -d ".cache/`{{.cache_key}}`" ]; then
        echo "hit"
      else
        echo "miss"
      fi
    outputs:
      - name: cache_status
        from: stdout
        
  - name: "Restore from cache"
    command: "cp -r .cache/`{{.cache_key}}`/node_modules ."
    condition: '`{{.cache_status}}` == "hit"'
    
  - name: "Install dependencies"
    command: "npm ci"
    condition: '`{{.cache_status}}` == "miss"'
    
  - name: "Save to cache"
    command: |
      mkdir -p .cache/`{{.cache_key}}`
      cp -r node_modules .cache/`{{.cache_key}}`/
    condition: '`{{.cache_status}}` == "miss"'
```

### Resource Pooling

Manage limited resources:

```yaml
vars:
  max_db_connections: 5
  
steps:
  - name: "Initialize connection pool"
    command: "echo 0 > .connection_count"
    
  - name: "Process items with connection limit"
    command: |
      # Wait for available connection
      while [ $(cat .connection_count) -ge `{{.max_db_connections}}` ]; do
        sleep 1
      done
      
      # Acquire connection
      count=$(cat .connection_count)
      echo $((count + 1)) > .connection_count
      
      # Process with connection
      process-with-db.sh `{{.item}}`
      
      # Release connection
      count=$(cat .connection_count)
      echo $((count - 1)) > .connection_count
    for_each:
      items: "`{{.items_to_process}}`"
      var: item
    parallel:
      max_concurrent: 10  # More than connections to show queueing
```

## Security Patterns

### Secret Rotation

Automatically rotate secrets:

```yaml
steps:
  - name: "Check secret age"
    command: |
      secret_date=$(vault read -format=json secret/api-key | jq -r '.data.created')
      age_days=$(( ($(date +%s) - $(date -d "$secret_date" +%s)) / 86400 ))
      echo $age_days
    outputs:
      - name: secret_age_days
        from: stdout
        
  - name: "Rotate if needed"
    condition: "`{{.secret_age_days}}` > 30"
    steps:
      - name: "Generate new secret"
        command: "openssl rand -hex 32"
        outputs:
          - name: new_secret
            from: stdout
            
      - name: "Update vault"
        command: |
          vault write secret/api-key \
            value=`{{.new_secret}}` \
            created=$(date -u +%Y-%m-%dT%H:%M:%SZ)
            
      - name: "Update applications"
        command: "update-secret.sh `{{.new_secret}}`"
```

### Audit Trail

Maintain execution audit logs:

```yaml
steps:
  - name: "Log execution start"
    command: |
      jq -n '{
        "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
        "action": "execution_start",
        "user": "`{{.user | default env.USER}}`",
        "plan": "`{{.plan_name}}`",
        "environment": "`{{.environment}}`"
      }' >> audit.log
      
  - name: "Execute with audit"
    command: "sensitive-operation.sh"
    on_success:
      - name: "Log success"
        command: |
          jq -n '{
            "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
            "action": "execution_complete",
            "status": "success"
          }' >> audit.log
    on_failure:
      - name: "Log failure"
        command: |
          jq -n '{
            "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
            "action": "execution_complete",
            "status": "failure",
            "error": "`{{.error_message}}`"
          }' >> audit.log
```