# Real-World Examples

Production-ready Plexr configurations for common scenarios in software development, operations, and automation.

## Full-Stack Application Setup

Complete setup for a modern full-stack application with frontend, backend, and database:

```yaml
name: "Full-stack application setup"
version: "1.0"
description: "Setup React frontend, Node.js backend, and PostgreSQL database"

vars:
  app_name: "myapp"
  node_version: "18"
  postgres_version: "15"
  redis_version: "7"

steps:
  # Environment check
  - name: "Check prerequisites"
    command: |
      echo "Checking required tools..."
      for cmd in git node npm docker docker-compose psql redis-cli; do
        if ! command -v $cmd &> /dev/null; then
          echo "ERROR: $cmd is not installed"
          exit 1
        fi
      done
      echo "All prerequisites met"

  # Database setup
  - name: "Start PostgreSQL"
    command: |
      docker run -d \
        --name `{{.app_name}}`-postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB=`{{.app_name}}` \
        -p 5432:5432 \
        postgres:`{{.postgres_version}}`
    outputs:
      - name: db_container_id
        from: stdout

  - name: "Wait for PostgreSQL"
    command: |
      for i in {1..30}; do
        if docker exec `{{.app_name}}`-postgres pg_isready -U postgres; then
          echo "PostgreSQL is ready"
          exit 0
        fi
        echo "Waiting for PostgreSQL... ($i/30)"
        sleep 2
      done
      exit 1
    depends_on: ["Start PostgreSQL"]

  # Redis setup
  - name: "Start Redis"
    command: |
      docker run -d \
        --name `{{.app_name}}`-redis \
        -p 6379:6379 \
        redis:`{{.redis_version}}` \
        redis-server --appendonly yes
    outputs:
      - name: redis_container_id
        from: stdout

  # Backend setup
  - name: "Setup backend"
    parallel:
      - name: "Install backend dependencies"
        command: "npm install"
        workdir: "./backend"
        
      - name: "Setup environment"
        command: |
          cat > .env << EOF
          NODE_ENV=development
          PORT=3001
          DATABASE_URL=postgresql://postgres:postgres@localhost:5432/`{{.app_name}}`
          REDIS_URL=redis://localhost:6379
          JWT_SECRET=$(openssl rand -hex 32)
          EOF
        workdir: "./backend"

  - name: "Run database migrations"
    command: "npm run migrate"
    workdir: "./backend"
    depends_on: ["Wait for PostgreSQL", "Setup backend"]
    env:
      DATABASE_URL: "postgresql://postgres:postgres@localhost:5432/`{{.app_name}}`"

  - name: "Seed database"
    command: "npm run seed"
    workdir: "./backend"
    depends_on: ["Run database migrations"]
    condition: '`{{.seed_data | default "true"}}` == "true"'

  # Frontend setup
  - name: "Setup frontend"
    parallel:
      - name: "Install frontend dependencies"
        command: "npm install"
        workdir: "./frontend"
        
      - name: "Setup frontend environment"
        command: |
          cat > .env.local << EOF
          REACT_APP_API_URL=http://localhost:3001
          REACT_APP_WS_URL=ws://localhost:3001
          EOF
        workdir: "./frontend"

  # Start services
  - name: "Start backend"
    command: "npm run dev"
    workdir: "./backend"
    background: true
    depends_on: ["Run database migrations", "Start Redis"]
    outputs:
      - name: backend_pid
        from: stdout

  - name: "Start frontend"
    command: "npm start"
    workdir: "./frontend"
    background: true
    depends_on: ["Setup frontend"]
    outputs:
      - name: frontend_pid
        from: stdout

  - name: "Wait for services"
    command: |
      echo "Waiting for services to start..."
      sleep 5
      
      # Check backend
      if curl -f http://localhost:3001/health; then
        echo "Backend is running"
      else
        echo "Backend failed to start"
        exit 1
      fi
      
      # Check frontend
      if curl -f http://localhost:3000; then
        echo "Frontend is running"
      else
        echo "Frontend failed to start"
        exit 1
      fi
    depends_on: ["Start backend", "Start frontend"]

  - name: "Show access URLs"
    command: |
      echo "🚀 Application is ready!"
      echo "Frontend: http://localhost:3000"
      echo "Backend API: http://localhost:3001"
      echo "Database: postgresql://localhost:5432/`{{.app_name}}`"
      echo "Redis: redis://localhost:6379"
    depends_on: ["Wait for services"]
```

## Database Migration System

Production-grade database migration workflow:

```yaml
name: "Database migration system"
version: "1.0"
description: "Safe database migration with rollback support"

vars:
  db_host: "${DB_HOST:-localhost}"
  db_name: "${DB_NAME:-myapp}"
  db_user: "${DB_USER:-postgres}"
  migration_table: "schema_migrations"
  migrations_dir: "./migrations"

steps:
  # Pre-migration checks
  - name: "Verify database connection"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` -c "SELECT 1"
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    
  - name: "Create migration table"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` << EOF
      CREATE TABLE IF NOT EXISTS `{{.migration_table}}` (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        checksum VARCHAR(64),
        execution_time INTEGER,
        success BOOLEAN DEFAULT TRUE
      );
      EOF
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    depends_on: ["Verify database connection"]

  - name: "Get applied migrations"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` -t -c \
        "SELECT version FROM `{{.migration_table}}` WHERE success = true ORDER BY version"
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    outputs:
      - name: applied_migrations
        from: stdout
        as_array: true
    depends_on: ["Create migration table"]

  - name: "Get pending migrations"
    command: |
      ls `{{.migrations_dir}}`/*.up.sql | sort | while read file; do
        version=$(basename $file .up.sql)
        if ! echo "`{{.applied_migrations}}`" | grep -q "$version"; then
          echo "$file"
        fi
      done
    outputs:
      - name: pending_migrations
        from: stdout
        as_array: true
    depends_on: ["Get applied migrations"]

  - name: "Show migration plan"
    command: |
      if [ -z "`{{.pending_migrations}}`" ]; then
        echo "No pending migrations"
      else
        echo "Pending migrations:"
        echo "`{{.pending_migrations}}`" | tr ' ' '\n'
      fi
    depends_on: ["Get pending migrations"]

  # Backup before migration
  - name: "Create backup"
    command: |
      backup_file="backup-$(date +%Y%m%d-%H%M%S).sql"
      pg_dump -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` > $backup_file
      echo $backup_file
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    outputs:
      - name: backup_file
        from: stdout
    condition: '`{{.skip_backup | default "false"}}` != "true"'

  # Apply migrations
  - name: "Apply migration"
    command: |
      migration_file="`{{.migration}}`"
      version=$(basename $migration_file .up.sql)
      checksum=$(sha256sum $migration_file | cut -d' ' -f1)
      
      echo "Applying migration: $version"
      start_time=$(date +%s)
      
      # Begin transaction
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` << EOF
      BEGIN;
      
      -- Apply migration
      \i $migration_file
      
      -- Record migration
      INSERT INTO `{{.migration_table}}` (version, checksum, execution_time)
      VALUES ('$version', '$checksum', $(date +%s) - $start_time);
      
      COMMIT;
      EOF
      
      echo "Migration $version applied successfully"
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    for_each:
      items: "`{{.pending_migrations}}`"
      var: migration
    depends_on: ["Create backup"]
    on_failure:
      - name: "Rollback migration"
        command: |
          version=$(basename `{{.migration}}` .up.sql)
          rollback_file="`{{.migrations_dir}}`/$version.down.sql"
          
          if [ -f "$rollback_file" ]; then
            echo "Rolling back migration: $version"
            psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` < $rollback_file
            
            # Mark as failed
            psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` -c \
              "UPDATE `{{.migration_table}}` SET success = false WHERE version = '$version'"
          else
            echo "No rollback file found for $version"
          fi
        env:
          PGPASSWORD: "${DB_PASSWORD}"

  # Post-migration validation
  - name: "Validate schema"
    command: |
      # Run schema validation tests
      npm run test:schema
    depends_on: ["Apply migration"]
    condition: '`{{.skip_validation | default "false"}}` != "true"'

  - name: "Migration summary"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` << EOF
      SELECT 
        version,
        applied_at,
        execution_time || 's' as duration,
        CASE WHEN success THEN '✓' ELSE '✗' END as status
      FROM `{{.migration_table}}`
      ORDER BY applied_at DESC
      LIMIT 10;
      EOF
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    depends_on: ["Apply migration"]
```

## CI/CD Pipeline

Complete CI/CD pipeline with testing, building, and deployment:

```yaml
name: "CI/CD Pipeline"
version: "1.0"
description: "Full CI/CD pipeline with parallel testing and staged deployments"

vars:
  app_name: "${APP_NAME}"
  git_sha: "${GITHUB_SHA:-$(git rev-parse HEAD)}"
  git_branch: "${GITHUB_REF_NAME:-$(git branch --show-current)}"
  docker_registry: "${DOCKER_REGISTRY:-ghcr.io}"
  docker_repo: "`{{.docker_registry}}`/`{{.github_repository | default .app_name}}`"

steps:
  # Code quality checks
  - name: "Code quality"
    parallel:
      - name: "Lint"
        command: "npm run lint"
        
      - name: "Type check"
        command: "npm run type-check"
        
      - name: "Security audit"
        command: |
          npm audit --audit-level=high
          trivy fs --severity HIGH,CRITICAL .
          
      - name: "License check"
        command: "license-checker --onlyAllow 'MIT;Apache-2.0;BSD-3-Clause;ISC'"

  # Testing
  - name: "Test suite"
    parallel:
      - name: "Unit tests"
        command: "npm run test:unit -- --coverage"
        outputs:
          - name: unit_coverage
            from: stdout
            regex: "Statements\\s+:\\s+([0-9.]+)%"
            
      - name: "Integration tests"
        command: |
          docker-compose -f docker-compose.test.yml up -d
          npm run test:integration
          docker-compose -f docker-compose.test.yml down
          
      - name: "E2E tests"
        command: |
          npm run build
          npm run test:e2e
        env:
          HEADLESS: "true"
    depends_on: ["Code quality"]

  - name: "Check coverage"
    command: |
      if (( $(echo "`{{.unit_coverage}}` < 80" | bc -l) )); then
        echo "Coverage `{{.unit_coverage}}`% is below threshold of 80%"
        exit 1
      fi
      echo "Coverage `{{.unit_coverage}}`% meets threshold"
    depends_on: ["Test suite"]

  # Build
  - name: "Build application"
    command: |
      npm run build
      echo "`{{.git_sha}}`" > dist/VERSION
      echo "`{{.git_branch}}`" > dist/BRANCH
      date -u +"%Y-%m-%dT%H:%M:%SZ" > dist/BUILD_TIME
    outputs:
      - name: build_time
        from: stdout
    depends_on: ["Test suite"]

  - name: "Build Docker image"
    command: |
      docker build \
        --build-arg VERSION=`{{.git_sha}}` \
        --build-arg BUILD_TIME=`{{.build_time}}` \
        -t `{{.docker_repo}}`:`{{.git_sha}}` \
        -t `{{.docker_repo}}`:`{{.git_branch}}`-latest \
        .
    depends_on: ["Build application"]

  # Security scanning
  - name: "Scan Docker image"
    command: |
      trivy image --severity HIGH,CRITICAL `{{.docker_repo}}`:`{{.git_sha}}`
      grype `{{.docker_repo}}`:`{{.git_sha}}` -f critical
    depends_on: ["Build Docker image"]

  # Push to registry
  - name: "Push Docker image"
    command: |
      echo "$DOCKER_PASSWORD" | docker login `{{.docker_registry}}` -u "$DOCKER_USERNAME" --password-stdin
      docker push `{{.docker_repo}}`:`{{.git_sha}}`
      docker push `{{.docker_repo}}`:`{{.git_branch}}`-latest
    depends_on: ["Scan Docker image"]
    condition: '`{{.git_branch}}` == "main" || `{{.git_branch}}` == "develop"'

  # Deploy to environments
  - name: "Deploy to staging"
    command: |
      kubectl set image deployment/`{{.app_name}}` \
        `{{.app_name}}`=`{{.docker_repo}}`:`{{.git_sha}}` \
        -n staging
        
      kubectl rollout status deployment/`{{.app_name}}` -n staging
    depends_on: ["Push Docker image"]
    condition: '`{{.git_branch}}` == "develop"'

  - name: "Smoke tests - staging"
    command: |
      ./scripts/smoke-tests.sh https://staging.example.com
    depends_on: ["Deploy to staging"]
    retry:
      attempts: 3
      delay: 30s

  - name: "Deploy to production"
    command: |
      # Blue-green deployment
      kubectl apply -f - << EOF
      apiVersion: v1
      kind: Service
      metadata:
        name: `{{.app_name}}`-green
        namespace: production
      spec:
        selector:
          app: `{{.app_name}}`
          version: `{{.git_sha}}`
        ports:
        - port: 80
          targetPort: 8080
      ---
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: `{{.app_name}}`-`{{.git_sha}}`
        namespace: production
      spec:
        replicas: 3
        selector:
          matchLabels:
            app: `{{.app_name}}`
            version: `{{.git_sha}}`
        template:
          metadata:
            labels:
              app: `{{.app_name}}`
              version: `{{.git_sha}}`
          spec:
            containers:
            - name: `{{.app_name}}`
              image: `{{.docker_repo}}`:`{{.git_sha}}`
              ports:
              - containerPort: 8080
      EOF
      
      kubectl rollout status deployment/`{{.app_name}}`-`{{.git_sha}}` -n production
    depends_on: ["Smoke tests - staging"]
    condition: '`{{.git_branch}}` == "main" && `{{.manual_approval}}` == "true"'

  - name: "Health check - production"
    command: |
      for i in {1..10}; do
        if curl -f https://example.com/health; then
          echo "Production deployment healthy"
          exit 0
        fi
        echo "Health check attempt $i/10 failed"
        sleep 6
      done
      exit 1
    depends_on: ["Deploy to production"]

  - name: "Switch traffic"
    command: |
      # Update main service to point to new deployment
      kubectl patch service `{{.app_name}}` -n production -p \
        '{"spec":{"selector":{"version":"`{{.git_sha}}`"}}`}'
        
      echo "Traffic switched to version `{{.git_sha}}`"
    depends_on: ["Health check - production"]

  - name: "Cleanup old deployments"
    command: |
      # Keep only last 3 deployments
      kubectl get deployments -n production -l app=`{{.app_name}}` \
        --sort-by=.metadata.creationTimestamp -o name | \
        head -n -3 | xargs -r kubectl delete -n production
    depends_on: ["Switch traffic"]

  # Notifications
  - name: "Send notifications"
    parallel:
      - name: "Slack notification"
        command: |
          curl -X POST $SLACK_WEBHOOK -H 'Content-type: application/json' -d '{
            "text": "Deployment successful",
            "attachments": [{
              "color": "good",
              "fields": [
                {"title": "Application", "value": "`{{.app_name}}`", "short": true},
                {"title": "Version", "value": "`{{.git_sha}}`", "short": true},
                {"title": "Environment", "value": "production", "short": true},
                {"title": "Branch", "value": "`{{.git_branch}}`", "short": true}
              ]
            }]
          }'
          
      - name: "Update deployment tracker"
        command: |
          curl -X POST https://deployments.example.com/api/deployments \
            -H "Authorization: Bearer $DEPLOY_TOKEN" \
            -H "Content-Type: application/json" \
            -d '{
              "app": "`{{.app_name}}`",
              "version": "`{{.git_sha}}`",
              "environment": "production",
              "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
              "status": "success"
            }'
    depends_on: ["Switch traffic"]
    always_run: true
```

## Team Collaboration Workflow

Multi-user development environment setup:

```yaml
name: "Team development environment"
version: "1.0"
description: "Standardized development environment for team collaboration"

vars:
  project_name: "${PROJECT_NAME}"
  team_members: "${TEAM_MEMBERS}"  # Comma-separated list
  base_port: 3000

steps:
  # Setup shared services
  - name: "Create shared network"
    command: "docker network create `{{.project_name}}`-dev || true"

  - name: "Start shared database"
    command: |
      docker run -d \
        --name `{{.project_name}}`-db \
        --network `{{.project_name}}`-dev \
        -e POSTGRES_PASSWORD=dev \
        -e POSTGRES_DB=`{{.project_name}}` \
        -v `{{.project_name}}`-db-data:/var/lib/postgresql/data \
        -p 5432:5432 \
        postgres:15
        
  - name: "Start shared Redis"
    command: |
      docker run -d \
        --name `{{.project_name}}`-redis \
        --network `{{.project_name}}`-dev \
        -v `{{.project_name}}`-redis-data:/data \
        -p 6379:6379 \
        redis:7 redis-server --appendonly yes

  # Setup for each team member
  - name: "Setup developer environment"
    command: |
      member="`{{.member}}`"
      port_offset=`{{.index}}`
      
      # Create developer-specific config
      mkdir -p .dev/$member
      cat > .dev/$member/config.env << EOF
      DEVELOPER=$member
      API_PORT=$((`{{.base_port}}` + port_offset * 10))
      WEB_PORT=$((`{{.base_port}}` + port_offset * 10 + 1))
      DB_NAME=`{{.project_name}}`_${member}
      REDIS_DB=$port_offset
      EOF
      
      # Create developer database
      docker exec `{{.project_name}}`-db psql -U postgres -c \
        "CREATE DATABASE `{{.project_name}}`_${member};"
        
      echo "Environment created for $member"
      echo "API Port: $((`{{.base_port}}` + port_offset * 10))"
      echo "Web Port: $((`{{.base_port}}` + port_offset * 10 + 1))"
    for_each:
      items: "`{{.team_members | split \",\"}}`"
      var: member
      index: index
    depends_on: ["Start shared database", "Start shared Redis"]

  # Setup development tools
  - name: "Start development proxy"
    command: |
      # Create nginx config for team
      cat > .dev/nginx.conf << 'EOF'
      events { worker_connections 1024; }
      http {
        `{{range $i, $member := .team_members | split ","}}`
        upstream `{{$member}}`_api {
          server localhost:`{{add $.base_port (mul $i 10)}}`;
        }
        upstream `{{$member}}`_web {
          server localhost:`{{add $.base_port (mul $i 10) 1}}`;
        }
        `{{end}}`
        
        server {
          listen 80;
          
          `{{range $i, $member := .team_members | split ","}}`
          location /`{{$member}}`/api {
            rewrite ^/`{{$member}}`/api(.*)$ $1 break;
            proxy_pass http://`{{$member}}`_api;
          }
          location /`{{$member}}` {
            proxy_pass http://`{{$member}}`_web;
          }
          `{{end}}`
        }
      }
      EOF
      
      docker run -d \
        --name `{{.project_name}}`-proxy \
        -p 80:80 \
        -v $(pwd)/.dev/nginx.conf:/etc/nginx/nginx.conf:ro \
        nginx:alpine

  # Setup code sharing
  - name: "Initialize Git hooks"
    command: |
      # Pre-commit hook for code standards
      cat > .git/hooks/pre-commit << 'EOF'
      #!/bin/bash
      echo "Running pre-commit checks..."
      
      # Run linting
      npm run lint || exit 1
      
      # Run tests for changed files
      git diff --cached --name-only | grep -E '\.(js|ts)$' | xargs npm test -- --findRelatedTests
      
      # Check for console.log statements
      if git diff --cached | grep -E '^\+.*console\.log'; then
        echo "Error: console.log statements found"
        exit 1
      fi
      EOF
      
      chmod +x .git/hooks/pre-commit

  - name: "Setup code review tools"
    command: |
      # Install and configure code review tools
      npm install -D prettier eslint husky lint-staged
      
      # Configure prettier
      cat > .prettierrc << EOF
      {
        "semi": true,
        "singleQuote": true,
        "tabWidth": 2,
        "trailingComma": "es5"
      }
      EOF
      
      # Setup husky
      npx husky install
      npx husky add .husky/pre-commit "npx lint-staged"
      
      # Configure lint-staged
      cat > .lintstagedrc.json << EOF
      {
        "*.{js,ts,jsx,tsx}": ["eslint --fix", "prettier --write"],
        "*.{json,md,yml,yaml}": ["prettier --write"]
      }
      EOF

  # Documentation and onboarding
  - name: "Generate developer docs"
    command: |
      cat > .dev/README.md << EOF
      # `{{.project_name}}` Development Environment
      
      ## Quick Start
      
      1. Run \`plexr execute\` to setup your environment
      2. Source your environment: \`source .dev/$USER/config.env\`
      3. Start your services: \`npm run dev\`
      
      ## Team Environments
      
      | Developer | API Port | Web Port | URL |
      |-----------|----------|----------|-----|
      `{{range $i, $member := .team_members | split ","}}`| `{{$member}}` | `{{add $.base_port (mul $i 10)}}` | `{{add $.base_port (mul $i 10) 1}}` | http://localhost/`{{$member}}` |
      `{{end}}`
      
      ## Shared Services
      
      - Database: postgresql://localhost:5432/`{{.project_name}}`_$USER
      - Redis: redis://localhost:6379 (use your assigned DB number)
      - Proxy: http://localhost (routes to all developer environments)
      
      ## Commands
      
      - \`npm run dev\` - Start your development server
      - \`npm test\` - Run tests
      - \`npm run lint\` - Check code style
      - \`npm run db:migrate\` - Run database migrations
      
      ## Git Workflow
      
      1. Create feature branch: \`git checkout -b feature/your-feature\`
      2. Make changes and commit
      3. Push and create PR
      4. Request review from team
      EOF
      
      echo "Documentation generated at .dev/README.md"

  - name: "Show summary"
    command: |
      echo "🚀 Team development environment ready!"
      echo ""
      echo "Shared services:"
      echo "- Database: postgresql://localhost:5432"
      echo "- Redis: redis://localhost:6379"
      echo "- Proxy: http://localhost"
      echo ""
      echo "Developer environments created for:"
      echo "`{{.team_members}}`" | tr ',' '\n' | sed 's/^/- /'
      echo ""
      echo "Next steps:"
      echo "1. Source your environment: source .dev/$USER/config.env"
      echo "2. Run migrations: npm run db:migrate"
      echo "3. Start development: npm run dev"
      echo ""
      echo "See .dev/README.md for more information"
```

## Monitoring and Alerting Setup

Production monitoring stack deployment:

```yaml
name: "Monitoring stack deployment"
version: "1.0"
description: "Deploy Prometheus, Grafana, and AlertManager"

vars:
  stack_name: "monitoring"
  prometheus_retention: "30d"
  grafana_admin_password: "${GRAFANA_ADMIN_PASSWORD}"
  slack_webhook: "${SLACK_WEBHOOK}"
  pagerduty_key: "${PAGERDUTY_KEY}"

steps:
  - name: "Create monitoring namespace"
    command: "kubectl create namespace `{{.stack_name}}` --dry-run=client -o yaml | kubectl apply -f -"

  - name: "Deploy Prometheus"
    command: |
      cat << EOF | kubectl apply -f -
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: prometheus-config
        namespace: `{{.stack_name}}`
      data:
        prometheus.yml: |
          global:
            scrape_interval: 30s
            evaluation_interval: 30s
          
          rule_files:
            - /etc/prometheus/rules/*.yml
          
          alerting:
            alertmanagers:
              - static_configs:
                  - targets: ['alertmanager:9093']
          
          scrape_configs:
            - job_name: 'kubernetes-pods'
              kubernetes_sd_configs:
                - role: pod
              relabel_configs:
                - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
                  action: keep
                  regex: true
                - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
                  action: replace
                  target_label: __metrics_path__
                  regex: (.+)
                - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
                  action: replace
                  regex: ([^:]+)(?::\d+)?;(\d+)
                  replacement: \$1:\$2
                  target_label: __address__
      ---
      apiVersion: apps/v1
      kind: StatefulSet
      metadata:
        name: prometheus
        namespace: `{{.stack_name}}`
      spec:
        serviceName: prometheus
        replicas: 1
        selector:
          matchLabels:
            app: prometheus
        template:
          metadata:
            labels:
              app: prometheus
          spec:
            containers:
            - name: prometheus
              image: prom/prometheus:latest
              args:
                - '--config.file=/etc/prometheus/prometheus.yml'
                - '--storage.tsdb.path=/prometheus'
                - '--storage.tsdb.retention.time=`{{.prometheus_retention}}`'
              ports:
              - containerPort: 9090
              volumeMounts:
              - name: config
                mountPath: /etc/prometheus
              - name: storage
                mountPath: /prometheus
            volumes:
            - name: config
              configMap:
                name: prometheus-config
        volumeClaimTemplates:
        - metadata:
            name: storage
          spec:
            accessModes: ["ReadWriteOnce"]
            resources:
              requests:
                storage: 50Gi
      ---
      apiVersion: v1
      kind: Service
      metadata:
        name: prometheus
        namespace: `{{.stack_name}}`
      spec:
        ports:
        - port: 9090
        selector:
          app: prometheus
      EOF

  - name: "Deploy alert rules"
    command: |
      cat << 'EOF' | kubectl apply -f -
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: prometheus-rules
        namespace: `{{.stack_name}}`
      data:
        alerts.yml: |
          groups:
            - name: application
              interval: 30s
              rules:
                - alert: HighErrorRate
                  expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
                  for: 5m
                  labels:
                    severity: critical
                  annotations:
                    summary: "High error rate detected"
                    description: "Error rate is `{{ $value }}` for `{{ $labels.instance }}`"
                
                - alert: HighMemoryUsage
                  expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) > 0.9
                  for: 5m
                  labels:
                    severity: warning
                  annotations:
                    summary: "High memory usage"
                    description: "Memory usage is above 90% on `{{ $labels.instance }}`"
                
                - alert: PodCrashLooping
                  expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
                  for: 5m
                  labels:
                    severity: critical
                  annotations:
                    summary: "Pod is crash looping"
                    description: "Pod `{{ $labels.namespace }}`/`{{ $labels.pod }}` is crash looping"
      EOF
    depends_on: ["Deploy Prometheus"]

  - name: "Deploy AlertManager"
    command: |
      cat << EOF | kubectl apply -f -
      apiVersion: v1
      kind: Secret
      metadata:
        name: alertmanager-config
        namespace: `{{.stack_name}}`
      stringData:
        alertmanager.yml: |
          global:
            resolve_timeout: 5m
            slack_api_url: '`{{.slack_webhook}}`'
          
          route:
            group_by: ['alertname', 'cluster', 'service']
            group_wait: 10s
            group_interval: 10s
            repeat_interval: 12h
            receiver: 'default'
            routes:
              - match:
                  severity: critical
                receiver: 'critical'
          
          receivers:
            - name: 'default'
              slack_configs:
                - channel: '#alerts'
                  title: 'Alert: `{{ .GroupLabels.alertname }}`'
                  text: '`{{ range .Alerts }}``{{ .Annotations.description }}``{{ end }}`'
            
            - name: 'critical'
              slack_configs:
                - channel: '#alerts-critical'
                  title: '🚨 CRITICAL: `{{ .GroupLabels.alertname }}`'
                  text: '`{{ range .Alerts }}``{{ .Annotations.description }}``{{ end }}`'
              pagerduty_configs:
                - service_key: '`{{.pagerduty_key}}`'
      ---
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: alertmanager
        namespace: `{{.stack_name}}`
      spec:
        replicas: 1
        selector:
          matchLabels:
            app: alertmanager
        template:
          metadata:
            labels:
              app: alertmanager
          spec:
            containers:
            - name: alertmanager
              image: prom/alertmanager:latest
              args:
                - '--config.file=/etc/alertmanager/alertmanager.yml'
              ports:
              - containerPort: 9093
              volumeMounts:
              - name: config
                mountPath: /etc/alertmanager
            volumes:
            - name: config
              secret:
                secretName: alertmanager-config
      ---
      apiVersion: v1
      kind: Service
      metadata:
        name: alertmanager
        namespace: `{{.stack_name}}`
      spec:
        ports:
        - port: 9093
        selector:
          app: alertmanager
      EOF

  - name: "Deploy Grafana"
    command: |
      cat << EOF | kubectl apply -f -
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: grafana-datasources
        namespace: `{{.stack_name}}`
      data:
        prometheus.yaml: |
          apiVersion: 1
          datasources:
            - name: Prometheus
              type: prometheus
              access: proxy
              url: http://prometheus:9090
              isDefault: true
      ---
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: grafana
        namespace: `{{.stack_name}}`
      spec:
        replicas: 1
        selector:
          matchLabels:
            app: grafana
        template:
          metadata:
            labels:
              app: grafana
          spec:
            containers:
            - name: grafana
              image: grafana/grafana:latest
              env:
              - name: GF_SECURITY_ADMIN_PASSWORD
                value: "`{{.grafana_admin_password}}`"
              - name: GF_INSTALL_PLUGINS
                value: "grafana-clock-panel,grafana-simple-json-datasource"
              ports:
              - containerPort: 3000
              volumeMounts:
              - name: datasources
                mountPath: /etc/grafana/provisioning/datasources
            volumes:
            - name: datasources
              configMap:
                name: grafana-datasources
      ---
      apiVersion: v1
      kind: Service
      metadata:
        name: grafana
        namespace: `{{.stack_name}}`
      spec:
        type: LoadBalancer
        ports:
        - port: 80
          targetPort: 3000
        selector:
          app: grafana
      EOF

  - name: "Wait for deployments"
    command: |
      for deployment in prometheus alertmanager grafana; do
        kubectl rollout status deployment/$deployment -n `{{.stack_name}}`
      done
    depends_on: ["Deploy Prometheus", "Deploy AlertManager", "Deploy Grafana"]

  - name: "Import Grafana dashboards"
    command: |
      # Wait for Grafana to be ready
      sleep 30
      
      # Get Grafana URL
      GRAFANA_URL=$(kubectl get svc grafana -n `{{.stack_name}}` -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
      
      # Import dashboards
      for dashboard in dashboards/*.json; do
        curl -X POST \
          -H "Content-Type: application/json" \
          -u "admin:`{{.grafana_admin_password}}`" \
          -d @$dashboard \
          http://$GRAFANA_URL/api/dashboards/db
      done
    depends_on: ["Wait for deployments"]

  - name: "Show access information"
    command: |
      GRAFANA_URL=$(kubectl get svc grafana -n `{{.stack_name}}` -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
      
      echo "✅ Monitoring stack deployed successfully!"
      echo ""
      echo "Access URLs:"
      echo "- Grafana: http://$GRAFANA_URL (admin/`{{.grafana_admin_password}}`)"
      echo "- Prometheus: kubectl port-forward -n `{{.stack_name}}` svc/prometheus 9090:9090"
      echo "- AlertManager: kubectl port-forward -n `{{.stack_name}}` svc/alertmanager 9093:9093"
      echo ""
      echo "Next steps:"
      echo "1. Configure your applications to expose metrics"
      echo "2. Import additional Grafana dashboards"
      echo "3. Customize alert rules in Prometheus"
      echo "4. Test alerting by triggering a test alert"
    depends_on: ["Import Grafana dashboards"]
```