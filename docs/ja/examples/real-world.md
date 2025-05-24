# 実世界の例

ソフトウェア開発、運用、自動化における一般的なシナリオのための本番環境対応のPlexr設定。

## フルスタックアプリケーションセットアップ

フロントエンド、バックエンド、データベースを含むモダンなフルスタックアプリケーションの完全なセットアップ：

```yaml
name: "フルスタックアプリケーションセットアップ"
version: "1.0"
description: "Reactフロントエンド、Node.jsバックエンド、PostgreSQLデータベースをセットアップ"

vars:
  app_name: "myapp"
  node_version: "18"
  postgres_version: "15"
  redis_version: "7"

steps:
  # 環境チェック
  - name: "前提条件を確認"
    command: |
      echo "必要なツールを確認中..."
      for cmd in git node npm docker docker-compose psql redis-cli; do
        if ! command -v $cmd &> /dev/null; then
          echo "エラー: $cmd がインストールされていません"
          exit 1
        fi
      done
      echo "すべての前提条件が満たされています"

  # データベースセットアップ
  - name: "PostgreSQLを開始"
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

  - name: "PostgreSQLを待機"
    command: |
      for i in {1..30}; do
        if docker exec `{{.app_name}}`-postgres pg_isready -U postgres; then
          echo "PostgreSQLの準備ができました"
          exit 0
        fi
        echo "PostgreSQLを待機中... ($i/30)"
        sleep 2
      done
      exit 1
    depends_on: ["PostgreSQLを開始"]

  # Redisセットアップ
  - name: "Redisを開始"
    command: |
      docker run -d \
        --name `{{.app_name}}`-redis \
        -p 6379:6379 \
        redis:`{{.redis_version}}` \
        redis-server --appendonly yes
    outputs:
      - name: redis_container_id
        from: stdout

  # バックエンドセットアップ
  - name: "バックエンドをセットアップ"
    parallel:
      - name: "バックエンド依存関係をインストール"
        command: "npm install"
        workdir: "./backend"
        
      - name: "環境をセットアップ"
        command: |
          cat > .env << EOF
          NODE_ENV=development
          PORT=3001
          DATABASE_URL=postgresql://postgres:postgres@localhost:5432/`{{.app_name}}`
          REDIS_URL=redis://localhost:6379
          JWT_SECRET=$(openssl rand -hex 32)
          EOF
        workdir: "./backend"

  - name: "データベースマイグレーションを実行"
    command: "npm run migrate"
    workdir: "./backend"
    depends_on: ["PostgreSQLを待機", "バックエンドをセットアップ"]
    env:
      DATABASE_URL: "postgresql://postgres:postgres@localhost:5432/`{{.app_name}}`"

  - name: "データベースをシード"
    command: "npm run seed"
    workdir: "./backend"
    depends_on: ["データベースマイグレーションを実行"]
    condition: '`{{.seed_data | default "true"}}` == "true"'

  # フロントエンドセットアップ
  - name: "フロントエンドをセットアップ"
    parallel:
      - name: "フロントエンド依存関係をインストール"
        command: "npm install"
        workdir: "./frontend"
        
      - name: "フロントエンド環境をセットアップ"
        command: |
          cat > .env.local << EOF
          REACT_APP_API_URL=http://localhost:3001
          REACT_APP_WS_URL=ws://localhost:3001
          EOF
        workdir: "./frontend"

  # サービスを開始
  - name: "バックエンドを開始"
    command: "npm run dev"
    workdir: "./backend"
    background: true
    depends_on: ["データベースマイグレーションを実行", "Redisを開始"]
    outputs:
      - name: backend_pid
        from: stdout

  - name: "フロントエンドを開始"
    command: "npm start"
    workdir: "./frontend"
    background: true
    depends_on: ["フロントエンドをセットアップ"]
    outputs:
      - name: frontend_pid
        from: stdout

  - name: "サービスを待機"
    command: |
      echo "サービスの開始を待機中..."
      sleep 5
      
      # バックエンドを確認
      if curl -f http://localhost:3001/health; then
        echo "バックエンドが実行中"
      else
        echo "バックエンドの開始に失敗"
        exit 1
      fi
      
      # フロントエンドを確認
      if curl -f http://localhost:3000; then
        echo "フロントエンドが実行中"
      else
        echo "フロントエンドの開始に失敗"
        exit 1
      fi
    depends_on: ["バックエンドを開始", "フロントエンドを開始"]

  - name: "アクセスURLを表示"
    command: |
      echo "🚀 アプリケーションの準備ができました！"
      echo "フロントエンド: http://localhost:3000"
      echo "バックエンドAPI: http://localhost:3001"
      echo "データベース: postgresql://localhost:5432/`{{.app_name}}`"
      echo "Redis: redis://localhost:6379"
    depends_on: ["サービスを待機"]
```

## データベースマイグレーションシステム

ロールバックサポート付きの本番グレードのデータベースマイグレーションワークフロー：

```yaml
name: "データベースマイグレーションシステム"
version: "1.0"
description: "ロールバックサポート付きの安全なデータベースマイグレーション"

vars:
  db_host: "${DB_HOST:-localhost}"
  db_name: "${DB_NAME:-myapp}"
  db_user: "${DB_USER:-postgres}"
  migration_table: "schema_migrations"
  migrations_dir: "./migrations"

steps:
  # マイグレーション前チェック
  - name: "データベース接続を確認"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` -c "SELECT 1"
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    
  - name: "マイグレーションテーブルを作成"
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
    depends_on: ["データベース接続を確認"]

  - name: "適用されたマイグレーションを取得"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` -t -c \
        "SELECT version FROM `{{.migration_table}}` WHERE success = true ORDER BY version"
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    outputs:
      - name: applied_migrations
        from: stdout
        as_array: true
    depends_on: ["マイグレーションテーブルを作成"]

  - name: "保留中のマイグレーションを取得"
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
    depends_on: ["適用されたマイグレーションを取得"]

  - name: "マイグレーション計画を表示"
    command: |
      if [ -z "`{{.pending_migrations}}`" ]; then
        echo "保留中のマイグレーションはありません"
      else
        echo "保留中のマイグレーション:"
        echo "`{{.pending_migrations}}`" | tr ' ' '\n'
      fi
    depends_on: ["保留中のマイグレーションを取得"]

  # マイグレーション前のバックアップ
  - name: "バックアップを作成"
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

  # マイグレーションを適用
  - name: "マイグレーションを適用"
    command: |
      migration_file="`{{.migration}}`"
      version=$(basename $migration_file .up.sql)
      checksum=$(sha256sum $migration_file | cut -d' ' -f1)
      
      echo "マイグレーションを適用中: $version"
      start_time=$(date +%s)
      
      # トランザクションを開始
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` << EOF
      BEGIN;
      
      -- マイグレーションを適用
      \i $migration_file
      
      -- マイグレーションを記録
      INSERT INTO `{{.migration_table}}` (version, checksum, execution_time)
      VALUES ('$version', '$checksum', $(date +%s) - $start_time);
      
      COMMIT;
      EOF
      
      echo "マイグレーション $version が正常に適用されました"
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    for_each:
      items: "`{{.pending_migrations}}`"
      var: migration
    depends_on: ["バックアップを作成"]
    on_failure:
      - name: "マイグレーションをロールバック"
        command: |
          version=$(basename `{{.migration}}` .up.sql)
          rollback_file="`{{.migrations_dir}}`/$version.down.sql"
          
          if [ -f "$rollback_file" ]; then
            echo "マイグレーションをロールバック中: $version"
            psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` < $rollback_file
            
            # 失敗としてマーク
            psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` -c \
              "UPDATE `{{.migration_table}}` SET success = false WHERE version = '$version'"
          else
            echo "$version のロールバックファイルが見つかりません"
          fi
        env:
          PGPASSWORD: "${DB_PASSWORD}"

  # マイグレーション後の検証
  - name: "スキーマを検証"
    command: |
      # スキーマ検証テストを実行
      npm run test:schema
    depends_on: ["マイグレーションを適用"]
    condition: '`{{.skip_validation | default "false"}}` != "true"'

  - name: "マイグレーションサマリー"
    command: |
      psql -h `{{.db_host}}` -U `{{.db_user}}` -d `{{.db_name}}` << EOF
      SELECT 
        version,
        applied_at,
        execution_time || '秒' as 期間,
        CASE WHEN success THEN '✓' ELSE '✗' END as ステータス
      FROM `{{.migration_table}}`
      ORDER BY applied_at DESC
      LIMIT 10;
      EOF
    env:
      PGPASSWORD: "${DB_PASSWORD}"
    depends_on: ["マイグレーションを適用"]
```

## CI/CDパイプライン

テスト、ビルド、デプロイを含む完全なCI/CDパイプライン：

```yaml
name: "CI/CDパイプライン"
version: "1.0"
description: "並列テストとステージングデプロイメントを含む完全なCI/CDパイプライン"

vars:
  app_name: "${APP_NAME}"
  git_sha: "${GITHUB_SHA:-$(git rev-parse HEAD)}"
  git_branch: "${GITHUB_REF_NAME:-$(git branch --show-current)}"
  docker_registry: "${DOCKER_REGISTRY:-ghcr.io}"
  docker_repo: "`{{.docker_registry}}`/`{{.github_repository | default .app_name}}`"

steps:
  # コード品質チェック
  - name: "コード品質"
    parallel:
      - name: "リント"
        command: "npm run lint"
        
      - name: "型チェック"
        command: "npm run type-check"
        
      - name: "セキュリティ監査"
        command: |
          npm audit --audit-level=high
          trivy fs --severity HIGH,CRITICAL .
          
      - name: "ライセンスチェック"
        command: "license-checker --onlyAllow 'MIT;Apache-2.0;BSD-3-Clause;ISC'"

  # テスト
  - name: "テストスイート"
    parallel:
      - name: "ユニットテスト"
        command: "npm run test:unit -- --coverage"
        outputs:
          - name: unit_coverage
            from: stdout
            regex: "Statements\\s+:\\s+([0-9.]+)%"
            
      - name: "統合テスト"
        command: |
          docker-compose -f docker-compose.test.yml up -d
          npm run test:integration
          docker-compose -f docker-compose.test.yml down
          
      - name: "E2Eテスト"
        command: |
          npm run build
          npm run test:e2e
        env:
          HEADLESS: "true"
    depends_on: ["コード品質"]

  - name: "カバレッジを確認"
    command: |
      if (( $(echo "`{{.unit_coverage}}` < 80" | bc -l) )); then
        echo "カバレッジ `{{.unit_coverage}}`% が閾値80%を下回っています"
        exit 1
      fi
      echo "カバレッジ `{{.unit_coverage}}`% が閾値を満たしています"
    depends_on: ["テストスイート"]

  # ビルド
  - name: "アプリケーションをビルド"
    command: |
      npm run build
      echo "`{{.git_sha}}`" > dist/VERSION
      echo "`{{.git_branch}}`" > dist/BRANCH
      date -u +"%Y-%m-%dT%H:%M:%SZ" > dist/BUILD_TIME
    outputs:
      - name: build_time
        from: stdout
    depends_on: ["テストスイート"]

  - name: "Dockerイメージをビルド"
    command: |
      docker build \
        --build-arg VERSION=`{{.git_sha}}` \
        --build-arg BUILD_TIME=`{{.build_time}}` \
        -t `{{.docker_repo}}`:`{{.git_sha}}` \
        -t `{{.docker_repo}}`:`{{.git_branch}}`-latest \
        .
    depends_on: ["アプリケーションをビルド"]

  # セキュリティスキャン
  - name: "Dockerイメージをスキャン"
    command: |
      trivy image --severity HIGH,CRITICAL `{{.docker_repo}}`:`{{.git_sha}}`
      grype `{{.docker_repo}}`:`{{.git_sha}}` -f critical
    depends_on: ["Dockerイメージをビルド"]

  # レジストリにプッシュ
  - name: "Dockerイメージをプッシュ"
    command: |
      echo "$DOCKER_PASSWORD" | docker login `{{.docker_registry}}` -u "$DOCKER_USERNAME" --password-stdin
      docker push `{{.docker_repo}}`:`{{.git_sha}}`
      docker push `{{.docker_repo}}`:`{{.git_branch}}`-latest
    depends_on: ["Dockerイメージをスキャン"]
    condition: '`{{.git_branch}}` == "main" || `{{.git_branch}}` == "develop"'

  # 環境にデプロイ
  - name: "ステージングにデプロイ"
    command: |
      kubectl set image deployment/`{{.app_name}}` \
        `{{.app_name}}`=`{{.docker_repo}}`:`{{.git_sha}}` \
        -n staging
        
      kubectl rollout status deployment/`{{.app_name}}` -n staging
    depends_on: ["Dockerイメージをプッシュ"]
    condition: '`{{.git_branch}}` == "develop"'

  - name: "スモークテスト - ステージング"
    command: |
      ./scripts/smoke-tests.sh https://staging.example.com
    depends_on: ["ステージングにデプロイ"]
    retry:
      attempts: 3
      delay: 30s

  - name: "本番環境にデプロイ"
    command: |
      # ブルーグリーンデプロイメント
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
    depends_on: ["スモークテスト - ステージング"]
    condition: '`{{.git_branch}}` == "main" && `{{.manual_approval}}` == "true"'

  - name: "ヘルスチェック - 本番環境"
    command: |
      for i in {1..10}; do
        if curl -f https://example.com/health; then
          echo "本番デプロイメントが健全です"
          exit 0
        fi
        echo "ヘルスチェック試行 $i/10 が失敗しました"
        sleep 6
      done
      exit 1
    depends_on: ["本番環境にデプロイ"]

  - name: "トラフィックを切り替え"
    command: |
      # メインサービスを新しいデプロイメントに向ける
      kubectl patch service `{{.app_name}}` -n production -p \
        '{"spec":{"selector":{"version":"`{{.git_sha}}`"}}`}'
        
      echo "トラフィックがバージョン `{{.git_sha}}` に切り替わりました"
    depends_on: ["ヘルスチェック - 本番環境"]

  - name: "古いデプロイメントをクリーンアップ"
    command: |
      # 最新の3つのデプロイメントのみを保持
      kubectl get deployments -n production -l app=`{{.app_name}}` \
        --sort-by=.metadata.creationTimestamp -o name | \
        head -n -3 | xargs -r kubectl delete -n production
    depends_on: ["トラフィックを切り替え"]

  # 通知
  - name: "通知を送信"
    parallel:
      - name: "Slack通知"
        command: |
          curl -X POST $SLACK_WEBHOOK -H 'Content-type: application/json' -d '{
            "text": "デプロイメント成功",
            "attachments": [{
              "color": "good",
              "fields": [
                {"title": "アプリケーション", "value": "`{{.app_name}}`", "short": true},
                {"title": "バージョン", "value": "`{{.git_sha}}`", "short": true},
                {"title": "環境", "value": "production", "short": true},
                {"title": "ブランチ", "value": "`{{.git_branch}}`", "short": true}
              ]
            }]
          }'
          
      - name: "デプロイメントトラッカーを更新"
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
    depends_on: ["トラフィックを切り替え"]
    always_run: true
```

## チームコラボレーションワークフロー

マルチユーザー開発環境セットアップ：

```yaml
name: "チーム開発環境"
version: "1.0"
description: "チームコラボレーション用の標準化された開発環境"

vars:
  project_name: "${PROJECT_NAME}"
  team_members: "${TEAM_MEMBERS}"  # カンマ区切りリスト
  base_port: 3000

steps:
  # 共有サービスをセットアップ
  - name: "共有ネットワークを作成"
    command: "docker network create `{{.project_name}}`-dev || true"

  - name: "共有データベースを開始"
    command: |
      docker run -d \
        --name `{{.project_name}}`-db \
        --network `{{.project_name}}`-dev \
        -e POSTGRES_PASSWORD=dev \
        -e POSTGRES_DB=`{{.project_name}}` \
        -v `{{.project_name}}`-db-data:/var/lib/postgresql/data \
        -p 5432:5432 \
        postgres:15
        
  - name: "共有Redisを開始"
    command: |
      docker run -d \
        --name `{{.project_name}}`-redis \
        --network `{{.project_name}}`-dev \
        -v `{{.project_name}}`-redis-data:/data \
        -p 6379:6379 \
        redis:7 redis-server --appendonly yes

  # 各チームメンバー用のセットアップ
  - name: "開発者環境をセットアップ"
    command: |
      member="`{{.member}}`"
      port_offset=`{{.index}}`
      
      # 開発者固有の設定を作成
      mkdir -p .dev/$member
      cat > .dev/$member/config.env << EOF
      DEVELOPER=$member
      API_PORT=$((`{{.base_port}}` + port_offset * 10))
      WEB_PORT=$((`{{.base_port}}` + port_offset * 10 + 1))
      DB_NAME=`{{.project_name}}`_${member}
      REDIS_DB=$port_offset
      EOF
      
      # 開発者データベースを作成
      docker exec `{{.project_name}}`-db psql -U postgres -c \
        "CREATE DATABASE `{{.project_name}}`_${member};"
        
      echo "$member の環境が作成されました"
      echo "APIポート: $((`{{.base_port}}` + port_offset * 10))"
      echo "Webポート: $((`{{.base_port}}` + port_offset * 10 + 1))"
    for_each:
      items: "`{{.team_members | split \",\"}}`"
      var: member
      index: index
    depends_on: ["共有データベースを開始", "共有Redisを開始"]

  # 開発ツールをセットアップ
  - name: "開発プロキシを開始"
    command: |
      # チーム用のnginx設定を作成
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

  # コード共有をセットアップ
  - name: "Gitフックを初期化"
    command: |
      # コード標準のためのpre-commitフック
      cat > .git/hooks/pre-commit << 'EOF'
      #!/bin/bash
      echo "pre-commitチェックを実行中..."
      
      # リンティングを実行
      npm run lint || exit 1
      
      # 変更されたファイルのテストを実行
      git diff --cached --name-only | grep -E '\.(js|ts)$' | xargs npm test -- --findRelatedTests
      
      # console.log文をチェック
      if git diff --cached | grep -E '^\+.*console\.log'; then
        echo "エラー: console.log文が見つかりました"
        exit 1
      fi
      EOF
      
      chmod +x .git/hooks/pre-commit

  - name: "コードレビューツールをセットアップ"
    command: |
      # コードレビューツールをインストールして設定
      npm install -D prettier eslint husky lint-staged
      
      # prettierを設定
      cat > .prettierrc << EOF
      {
        "semi": true,
        "singleQuote": true,
        "tabWidth": 2,
        "trailingComma": "es5"
      }
      EOF
      
      # huskyをセットアップ
      npx husky install
      npx husky add .husky/pre-commit "npx lint-staged"
      
      # lint-stagedを設定
      cat > .lintstagedrc.json << EOF
      {
        "*.{js,ts,jsx,tsx}": ["eslint --fix", "prettier --write"],
        "*.{json,md,yml,yaml}": ["prettier --write"]
      }
      EOF

  # ドキュメントとオンボーディング
  - name: "開発者ドキュメントを生成"
    command: |
      cat > .dev/README.md << EOF
      # `{{.project_name}}` 開発環境
      
      ## クイックスタート
      
      1. \`plexr execute\` を実行して環境をセットアップ
      2. 環境をソース: \`source .dev/$USER/config.env\`
      3. サービスを開始: \`npm run dev\`
      
      ## チーム環境
      
      | 開発者 | APIポート | Webポート | URL |
      |-----------|----------|----------|-----|
      `{{range $i, $member := .team_members | split ","}}`| `{{$member}}` | `{{add $.base_port (mul $i 10)}}` | `{{add $.base_port (mul $i 10) 1}}` | http://localhost/`{{$member}}` |
      `{{end}}`
      
      ## 共有サービス
      
      - データベース: postgresql://localhost:5432/`{{.project_name}}`_$USER
      - Redis: redis://localhost:6379 (割り当てられたDB番号を使用)
      - プロキシ: http://localhost (すべての開発者環境にルーティング)
      
      ## コマンド
      
      - \`npm run dev\` - 開発サーバーを開始
      - \`npm test\` - テストを実行
      - \`npm run lint\` - コードスタイルをチェック
      - \`npm run db:migrate\` - データベースマイグレーションを実行
      
      ## Gitワークフロー
      
      1. フィーチャーブランチを作成: \`git checkout -b feature/your-feature\`
      2. 変更を加えてコミット
      3. プッシュしてPRを作成
      4. チームからレビューをリクエスト
      EOF
      
      echo "ドキュメントが .dev/README.md に生成されました"

  - name: "サマリーを表示"
    command: |
      echo "🚀 チーム開発環境の準備ができました！"
      echo ""
      echo "共有サービス:"
      echo "- データベース: postgresql://localhost:5432"
      echo "- Redis: redis://localhost:6379"
      echo "- プロキシ: http://localhost"
      echo ""
      echo "以下の開発者環境が作成されました:"
      echo "`{{.team_members}}`" | tr ',' '\n' | sed 's/^/- /'
      echo ""
      echo "次のステップ:"
      echo "1. 環境をソース: source .dev/$USER/config.env"
      echo "2. マイグレーションを実行: npm run db:migrate"
      echo "3. 開発を開始: npm run dev"
      echo ""
      echo "詳細は .dev/README.md を参照してください"
```

## モニタリングとアラートのセットアップ

本番モニタリングスタックのデプロイメント：

```yaml
name: "モニタリングスタックデプロイメント"
version: "1.0"
description: "Prometheus、Grafana、AlertManagerをデプロイ"

vars:
  stack_name: "monitoring"
  prometheus_retention: "30d"
  grafana_admin_password: "${GRAFANA_ADMIN_PASSWORD}"
  slack_webhook: "${SLACK_WEBHOOK}"
  pagerduty_key: "${PAGERDUTY_KEY}"

steps:
  - name: "モニタリング名前空間を作成"
    command: "kubectl create namespace `{{.stack_name}}` --dry-run=client -o yaml | kubectl apply -f -"

  - name: "Prometheusをデプロイ"
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

  - name: "アラートルールをデプロイ"
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
                    summary: "高いエラー率が検出されました"
                    description: "`{{ $labels.instance }}` のエラー率が `{{ $value }}` です"
                
                - alert: HighMemoryUsage
                  expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) > 0.9
                  for: 5m
                  labels:
                    severity: warning
                  annotations:
                    summary: "高いメモリ使用率"
                    description: "`{{ $labels.instance }}` のメモリ使用率が90%を超えています"
                
                - alert: PodCrashLooping
                  expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
                  for: 5m
                  labels:
                    severity: critical
                  annotations:
                    summary: "Podがクラッシュループしています"
                    description: "Pod `{{ $labels.namespace }}`/`{{ $labels.pod }}` がクラッシュループしています"
      EOF
    depends_on: ["Prometheusをデプロイ"]

  - name: "AlertManagerをデプロイ"
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
                  title: 'アラート: `{{ .GroupLabels.alertname }}`'
                  text: '`{{ range .Alerts }}``{{ .Annotations.description }}``{{ end }}`'
            
            - name: 'critical'
              slack_configs:
                - channel: '#alerts-critical'
                  title: '🚨 クリティカル: `{{ .GroupLabels.alertname }}`'
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

  - name: "Grafanaをデプロイ"
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

  - name: "デプロイメントを待機"
    command: |
      for deployment in prometheus alertmanager grafana; do
        kubectl rollout status deployment/$deployment -n `{{.stack_name}}`
      done
    depends_on: ["Prometheusをデプロイ", "AlertManagerをデプロイ", "Grafanaをデプロイ"]

  - name: "Grafanaダッシュボードをインポート"
    command: |
      # Grafanaの準備を待つ
      sleep 30
      
      # GrafanaのURLを取得
      GRAFANA_URL=$(kubectl get svc grafana -n `{{.stack_name}}` -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
      
      # ダッシュボードをインポート
      for dashboard in dashboards/*.json; do
        curl -X POST \
          -H "Content-Type: application/json" \
          -u "admin:`{{.grafana_admin_password}}`" \
          -d @$dashboard \
          http://$GRAFANA_URL/api/dashboards/db
      done
    depends_on: ["デプロイメントを待機"]

  - name: "アクセス情報を表示"
    command: |
      GRAFANA_URL=$(kubectl get svc grafana -n `{{.stack_name}}` -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
      
      echo "✅ モニタリングスタックが正常にデプロイされました！"
      echo ""
      echo "アクセスURL:"
      echo "- Grafana: http://$GRAFANA_URL (admin/`{{.grafana_admin_password}}`)"
      echo "- Prometheus: kubectl port-forward -n `{{.stack_name}}` svc/prometheus 9090:9090"
      echo "- AlertManager: kubectl port-forward -n `{{.stack_name}}` svc/alertmanager 9093:9093"
      echo ""
      echo "次のステップ:"
      echo "1. アプリケーションがメトリクスを公開するよう設定"
      echo "2. 追加のGrafanaダッシュボードをインポート"
      echo "3. Prometheusのアラートルールをカスタマイズ"
      echo "4. テストアラートをトリガーしてアラートをテスト"
    depends_on: ["Grafanaダッシュボードをインポート"]
```