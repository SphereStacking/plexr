# 高度なパターン

このガイドでは、複雑なワークフロー、条件付きロジック、高度な自動化シナリオのための高度なPlexrパターンを示します。

## 複雑な依存関係チェーン

### ダイヤモンド依存関係

複数のパスが収束する複雑な依存関係グラフを処理：

```yaml
name: "ダイヤモンド依存関係パターン"
version: "1.0"

steps:
  - name: "開始"
    command: "echo 'ワークフローを開始'"
    
  - name: "パスA1"
    command: "process-a1.sh"
    depends_on: ["開始"]
    
  - name: "パスA2"
    command: "process-a2.sh"
    depends_on: ["パスA1"]
    
  - name: "パスB1"
    command: "process-b1.sh"
    depends_on: ["開始"]
    
  - name: "パスB2"
    command: "process-b2.sh"
    depends_on: ["パスB1"]
    
  - name: "収束"
    command: "merge-results.sh"
    depends_on: ["パスA2", "パスB2"]
```

### 条件付き依存関係

前のステップの成功に基づいてステップを実行：

```yaml
steps:
  - name: "環境を確認"
    command: "check-env.sh"
    outputs:
      - name: env_ready
        from: stdout
        
  - name: "環境をセットアップ"
    command: "setup-env.sh"
    condition: '`{{.env_ready}}` != "true"'
    outputs:
      - name: env_ready
        value: "true"
        
  - name: "デプロイ"
    command: "deploy.sh"
    depends_on: ["環境を確認", "環境をセットアップ"]
    condition: '`{{.env_ready}}` == "true"'
```

### 動的依存関係

発見されたリソースに基づいて依存関係を生成：

```yaml
steps:
  - name: "サービスを発見"
    command: "kubectl get services -o json | jq -r '.items[].metadata.name'"
    outputs:
      - name: services
        from: stdout
        as_array: true
        
  - name: "ヘルスチェックを生成"
    command: |
      for service in `{{range .services}}``{{.}}` `{{end}}`; do
        echo "- name: \"$service を確認\""
        echo "  command: \"check-health.sh $service\""
        echo "  tags: [\"health-check\"]"
      done > dynamic-steps.yml
      
  - name: "ヘルスチェックを実行"
    include: "dynamic-steps.yml"
    depends_on: ["ヘルスチェックを生成"]
```

## 条件付き実行

### 複数条件ロジック

複数の条件を組み合わせる：

```yaml
steps:
  - name: "本番デプロイ"
    command: "deploy-prod.sh"
    condition: |
      `{{.environment}}` == "production" &&
      `{{.tests_passed}}` == "true" &&
      `{{.approval_granted}}` == "true"
      
  - name: "ステージングデプロイ"
    command: "deploy-staging.sh"
    condition: |
      `{{.environment}}` == "staging" &&
      `{{.tests_passed}}` == "true"
```

### Switch-Caseパターン

スイッチのような動作を実装：

```yaml
vars:
  deployment_commands:
    dev: "deploy-dev.sh"
    staging: "deploy-staging.sh"
    prod: "deploy-prod.sh"
    
steps:
  - name: "環境にデプロイ"
    command: "`{{index .deployment_commands .environment}}`"
    condition: "`{{hasKey .deployment_commands .environment}}`"
    
  - name: "不明な環境"
    command: "echo '不明な環境: `{{.environment}}`' && exit 1"
    condition: "`{{not (hasKey .deployment_commands .environment)}}`"
```

### 条件付きインクルード

条件に基づいて異なるプランをインクルード：

```yaml
steps:
  - name: "OSを検出"
    command: "uname -s"
    outputs:
      - name: os
        from: stdout
        
  - name: "Linuxセットアップ"
    include: "setup-linux.yml"
    condition: '`{{.os}}` == "Linux"'
    
  - name: "macOSセットアップ"
    include: "setup-macos.yml"
    condition: '`{{.os}}` == "Darwin"'
    
  - name: "Windowsセットアップ"
    include: "setup-windows.yml"
    condition: '`{{contains .os "MINGW"}}`'
```

## エラー処理戦略

### グレースフルデグラデーション

機能を減らして実行を継続：

```yaml
steps:
  - name: "プライマリデータベースを試行"
    command: "connect-primary-db.sh"
    ignore_failure: true
    outputs:
      - name: primary_db_available
        from: exit_code
        transform: '`{{if eq . "0"}}`true`{{else}}`false`{{end}}`'
        
  - name: "フォールバックデータベースを使用"
    command: "connect-fallback-db.sh"
    condition: '`{{.primary_db_available}}` != "true"'
    
  - name: "利用可能なデータベースで実行"
    command: "run-app.sh"
    env:
      USE_FALLBACK: '`{{.primary_db_available}}` != "true"'
```

### バックオフ付きリトライ

高度なリトライ戦略を実装：

```yaml
steps:
  - name: "リトライ付きAPIコール"
    command: "call-api.sh"
    retry:
      attempts: 5
      delay: 2s
      backoff: exponential
      max_delay: 60s
      on_retry: |
        echo "リトライ試行 `{{.retry_count}}` / `{{.retry_max}}`"
        echo "次のリトライまで `{{.retry_delay}}` 秒"
```

### サーキットブレーカーパターン

カスケード障害を防ぐ：

```yaml
vars:
  circuit_breaker:
    failure_threshold: 3
    timeout: 300  # 5分
    
steps:
  - name: "サーキット状態を確認"
    command: |
      if [ -f .circuit_open ] && [ $(( $(date +%s) - $(stat -f %m .circuit_open) )) -lt `{{.circuit_breaker.timeout}}` ]; then
        echo "open"
      else
        echo "closed"
      fi
    outputs:
      - name: circuit_state
        from: stdout
        
  - name: "サービスコール"
    command: "call-service.sh"
    condition: '`{{.circuit_state}}` != "open"'
    on_failure:
      - name: "失敗カウントを増加"
        command: |
          count=$(cat .failure_count 2>/dev/null || echo 0)
          echo $((count + 1)) > .failure_count
          if [ $((count + 1)) -ge `{{.circuit_breaker.failure_threshold}}` ]; then
            touch .circuit_open
          fi
```

## 状態管理パターン

### ステートマシン

ワークフローステートマシンを実装：

```yaml
vars:
  states:
    init: "初期化中"
    build: "ビルド中"
    test: "テスト中"
    deploy: "デプロイ中"
    complete: "完了"
    failed: "失敗"
    
steps:
  - name: "初期状態を設定"
    command: "echo '`{{.states.init}}`'"
    outputs:
      - name: workflow_state
        from: stdout
        
  - name: "ビルドフェーズ"
    command: "build.sh"
    condition: '`{{.workflow_state}}` == "`{{.states.init}}`"'
    outputs:
      - name: workflow_state
        value: "`{{.states.build}}`"
    on_failure:
      outputs:
        - name: workflow_state
          value: "`{{.states.failed}}`"
          
  - name: "テストフェーズ"
    command: "test.sh"
    condition: '`{{.workflow_state}}` == "`{{.states.build}}`"'
    outputs:
      - name: workflow_state
        value: "`{{.states.test}}`"
```

### 実行間の永続状態

実行間で状態を維持：

```yaml
steps:
  - name: "永続状態をロード"
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
        
  - name: "状態を更新"
    command: |
      echo '`{{.persistent_state | toJson}}`' | \
      jq '.last_run = "`{{now | date "2006-01-02T15:04:05Z07:00"}}`"' | \
      jq '.run_count = ((.run_count // 0) + 1)'
    outputs:
      - name: updated_state
        from: stdout
        json_parse: true
        
  - name: "状態を保存"
    command: "echo '`{{.updated_state | toJson}}`' > state.json"
    always_run: true
```

### 分散状態

複数のプラン実行間で状態を共有：

```yaml
steps:
  - name: "ロックを取得"
    command: |
      while ! mkdir .lock 2>/dev/null; do
        echo "ロックを待機中..."
        sleep 1
      done
    timeout: 60
    
  - name: "共有状態を読み取り"
    command: "cat shared-state.json 2>/dev/null || echo '{}'"
    outputs:
      - name: shared_state
        from: stdout
        json_parse: true
        
  - name: "共有状態を更新"
    command: |
      echo '`{{.shared_state | toJson}}`' | \
      jq '.workers["`{{.worker_id}}`"] = {
        "last_seen": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
        "status": "active"
      }' > shared-state.json
      
  - name: "ロックを解放"
    command: "rmdir .lock"
    always_run: true
```

## 動的ワークフロー生成

### テンプレートベースの生成

テンプレートからステップを生成：

```yaml
vars:
  environments: ["dev", "staging", "prod"]
  regions: ["us-east-1", "eu-west-1", "ap-southeast-1"]
  
steps:
  - name: "デプロイステップを生成"
    command: |
      cat > generated-steps.yml << EOF
      steps:
      `{{range $env := .environments}}`
      `{{range $region := $.regions}}`
        - name: "`{{$env}}`-`{{$region}}`にデプロイ"
          command: "deploy.sh"
          env:
            ENVIRONMENT: "`{{$env}}`"
            REGION: "`{{$region}}`"
          tags: ["deploy", "`{{$env}}`", "`{{$region}}`"]
      `{{end}}`
      `{{end}}`
      EOF
      
  - name: "デプロイを実行"
    include: "generated-steps.yml"
    depends_on: ["デプロイステップを生成"]
```

### 発見ベースのワークフロー

発見されたリソースに基づいてワークフローを構築：

```yaml
steps:
  - name: "マイクロサービスを発見"
    command: "find services -name 'Dockerfile' -type f | xargs dirname"
    outputs:
      - name: services
        from: stdout
        as_array: true
        
  - name: "サービスイメージをビルド"
    command: |
      docker build -t `{{.service | base}}`:`{{.version}}` `{{.service}}`
    for_each:
      items: "`{{.services}}`"
      var: service
    parallel:
      max_concurrent: 3
```

## 統合パターン

### イベント駆動ワークフロー

外部イベントに反応：

```yaml
steps:
  - name: "Webhookを待機"
    command: |
      # シンプルなWebhookリスナー
      while true; do
        nc -l 8080 | grep -q "deploy-signal" && break
      done
    timeout: 3600  # 1時間
    
  - name: "Webhookデータを処理"
    command: "parse-webhook.sh"
    outputs:
      - name: deployment_target
        from: stdout
        
  - name: "Webhookに基づいてデプロイ"
    command: "deploy.sh `{{.deployment_target}}`"
```

### ポーリングパターン

定期的に条件をチェック：

```yaml
steps:
  - name: "サービスの準備を待機"
    command: |
      for i in {1..60}; do
        if curl -s http://service:8080/health | grep -q "ok"; then
          echo "ready"
          exit 0
        fi
        echo "試行 $i/60: サービスが準備できていません、待機中..."
        sleep 5
      done
      echo "timeout"
      exit 1
    outputs:
      - name: service_status
        from: stdout
```

### メッセージキュー統合

キューからアイテムを処理：

```yaml
steps:
  - name: "キューメッセージを取得"
    command: "aws sqs receive-message --queue-url `{{.queue_url}}` --max-number-of-messages 10"
    outputs:
      - name: messages
        from: stdout
        json_path: "$.Messages"
        
  - name: "メッセージを処理"
    command: |
      echo '`{{.message | toJson}}`' | process-message.sh
    for_each:
      items: "`{{.messages}}`"
      var: message
    on_success:
      - name: "メッセージを削除"
        command: |
          aws sqs delete-message \
            --queue-url `{{.queue_url}}` \
            --receipt-handle `{{.message.ReceiptHandle}}`
```

## パフォーマンス最適化

### 並列マトリックス実行

組み合わせを並列で実行：

```yaml
vars:
  python_versions: ["3.8", "3.9", "3.10", "3.11"]
  test_suites: ["unit", "integration", "e2e"]
  
steps:
  - name: "テストマトリックス"
    parallel:
      - name: "Python `{{.version}}` - `{{.suite}}` をテスト"
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

### キャッシングパターン

インテリジェントキャッシングを実装：

```yaml
steps:
  - name: "キャッシュキーを計算"
    command: |
      echo "$(cat package-lock.json | sha256sum | cut -d' ' -f1)-$(date +%Y%m%d)"
    outputs:
      - name: cache_key
        from: stdout
        
  - name: "キャッシュを確認"
    command: |
      if [ -d ".cache/`{{.cache_key}}`" ]; then
        echo "hit"
      else
        echo "miss"
      fi
    outputs:
      - name: cache_status
        from: stdout
        
  - name: "キャッシュから復元"
    command: "cp -r .cache/`{{.cache_key}}`/node_modules ."
    condition: '`{{.cache_status}}` == "hit"'
    
  - name: "依存関係をインストール"
    command: "npm ci"
    condition: '`{{.cache_status}}` == "miss"'
    
  - name: "キャッシュに保存"
    command: |
      mkdir -p .cache/`{{.cache_key}}`
      cp -r node_modules .cache/`{{.cache_key}}`/
    condition: '`{{.cache_status}}` == "miss"'
```

### リソースプーリング

制限されたリソースを管理：

```yaml
vars:
  max_db_connections: 5
  
steps:
  - name: "接続プールを初期化"
    command: "echo 0 > .connection_count"
    
  - name: "接続制限でアイテムを処理"
    command: |
      # 利用可能な接続を待つ
      while [ $(cat .connection_count) -ge `{{.max_db_connections}}` ]; do
        sleep 1
      done
      
      # 接続を取得
      count=$(cat .connection_count)
      echo $((count + 1)) > .connection_count
      
      # 接続で処理
      process-with-db.sh `{{.item}}`
      
      # 接続を解放
      count=$(cat .connection_count)
      echo $((count - 1)) > .connection_count
    for_each:
      items: "`{{.items_to_process}}`"
      var: item
    parallel:
      max_concurrent: 10  # キューイングを示すために接続より多い
```

## セキュリティパターン

### シークレットローテーション

シークレットを自動的にローテート：

```yaml
steps:
  - name: "シークレットの経過時間を確認"
    command: |
      secret_date=$(vault read -format=json secret/api-key | jq -r '.data.created')
      age_days=$(( ($(date +%s) - $(date -d "$secret_date" +%s)) / 86400 ))
      echo $age_days
    outputs:
      - name: secret_age_days
        from: stdout
        
  - name: "必要に応じてローテート"
    condition: "`{{.secret_age_days}}` > 30"
    steps:
      - name: "新しいシークレットを生成"
        command: "openssl rand -hex 32"
        outputs:
          - name: new_secret
            from: stdout
            
      - name: "vaultを更新"
        command: |
          vault write secret/api-key \
            value=`{{.new_secret}}` \
            created=$(date -u +%Y-%m-%dT%H:%M:%SZ)
            
      - name: "アプリケーションを更新"
        command: "update-secret.sh `{{.new_secret}}`"
```

### 監査証跡

実行監査ログを維持：

```yaml
steps:
  - name: "実行開始をログ"
    command: |
      jq -n '{
        "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
        "action": "execution_start",
        "user": "`{{.user | default env.USER}}`",
        "plan": "`{{.plan_name}}`",
        "environment": "`{{.environment}}`"
      }' >> audit.log
      
  - name: "監査付きで実行"
    command: "sensitive-operation.sh"
    on_success:
      - name: "成功をログ"
        command: |
          jq -n '{
            "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
            "action": "execution_complete",
            "status": "success"
          }' >> audit.log
    on_failure:
      - name: "失敗をログ"
        command: |
          jq -n '{
            "timestamp": "`{{now | date "2006-01-02T15:04:05Z07:00"}}`",
            "action": "execution_complete",
            "status": "failure",
            "error": "`{{.error_message}}`"
          }' >> audit.log
```