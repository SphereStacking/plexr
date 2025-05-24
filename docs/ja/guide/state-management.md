# 状態管理

Plexrは、ステップ間でデータを共有し、実行の進捗を追跡し、ワークフローのコンテキストを維持するための強力な状態管理機能を提供します。

## 概要

Plexrの状態は、プランの実行全体を通じて永続化するキー・バリューストアです。これにより、ステップは以下のことができます：
- 後続のステップとデータを共有
- 実行結果を保存
- 条件付きの決定を行う
- ワークフローの進捗を追跡

## 状態ファイル

Plexrは、プロジェクトディレクトリの`.plexr/state.json`に状態を保存します：

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

## 状態変数の設定

### コマンド出力から

コマンド出力を状態変数にキャプチャ：

```yaml
steps:
  - name: "バージョンを取得"
    command: "git describe --tags --always"
    outputs:
      - name: version
        from: stdout
```

### JSON出力から

JSON出力を解析して特定のフィールドを抽出：

```yaml
steps:
  - name: "ビルド情報を取得"
    command: "npm run build --json"
    outputs:
      - name: build_id
        from: stdout
        json_path: "$.buildId"
      - name: build_time
        from: stdout
        json_path: "$.timestamp"
```

### 正規表現を使用

正規表現パターンを使用してデータを抽出：

```yaml
steps:
  - name: "デプロイURLを解析"
    command: "deploy.sh"
    outputs:
      - name: app_url
        from: stdout
        regex: "Deployed to: (https://.*)"
        regex_group: 1
```

### 複数の出力

単一のステップから複数の値をキャプチャ：

```yaml
steps:
  - name: "システム情報"
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

## 状態変数の使用

### 変数の置換

Goテンプレート構文を使用して状態変数を参照：

```yaml
steps:
  - name: "バージョンでビルド"
    command: "docker build -t myapp:`{{.version}}` ."
    
  - name: "デプロイ"
    command: "kubectl set image deployment/app app=myapp:`{{.version}}`"
```

### デフォルト値

欠落している変数のデフォルトを提供：

```yaml
steps:
  - name: "環境を設定"
    command: "setup.sh"
    env:
      ENVIRONMENT: "`{{.environment | default \"development\"}}`"
      DEBUG: "`{{.debug | default \"false\"}}`"
```

### 条件付きロジック

条件で状態変数を使用：

```yaml
steps:
  - name: "本番デプロイ"
    command: "deploy-prod.sh"
    condition: '`{{.environment}}` == "production" && `{{.tests_passed}}` == "true"'
```

## 高度な状態管理

### ネストされた変数

複雑なデータ構造を扱う：

```yaml
steps:
  - name: "設定を取得"
    command: "cat config.json"
    outputs:
      - name: config
        from: stdout
        json_parse: true
        
  - name: "ネストされた設定を使用"
    command: "connect.sh"
    env:
      DB_HOST: "`{{.config.database.host}}`"
      DB_PORT: "`{{.config.database.port}}`"
```

### 配列とループ

配列データを処理：

```yaml
steps:
  - name: "サービスを取得"
    command: "kubectl get services -o json"
    outputs:
      - name: services
        from: stdout
        json_path: "$.items[*].metadata.name"
        
  - name: "サービスを処理"
    command: "check-service.sh `{{.service}}`"
    for_each:
      items: "`{{.services}}`"
      var: service
```

### 状態の変更

状態変数を変換：

```yaml
steps:
  - name: "バージョンをインクリメント"
    command: |
      current=`{{.build_number | default "0"}}`
      echo $((current + 1))
    outputs:
      - name: build_number
        from: stdout
```

## 状態の永続化

### 状態の保存

状態は各ステップ後に自動的に保存されます。手動でチェックポイントを作成することもできます：

```yaml
steps:
  - name: "重要な操作"
    command: "important-task.sh"
    save_state: true  # 即座に状態を保存
```

### 以前の状態のロード

前回の実行から続行：

```bash
# 最後の状態から再開
plexr execute --resume

# 特定の状態ファイルをロード
plexr execute --state-file backup-state.json
```

### 状態のリセット

すべての状態データをクリア：

```bash
# 現在のプランの状態をリセット
plexr reset

# 状態ファイルをリセットして削除
plexr reset --clean
```

## 状態のスコープ

### グローバル変数

すべてのステップで利用可能な変数を設定：

```yaml
vars:
  project_name: "myapp"
  region: "us-east-1"
  
steps:
  - name: "デプロイ"
    command: "deploy.sh `{{.project_name}}` `{{.region}}`"
```

### ステップローカル変数

特定のステップにスコープされた変数：

```yaml
steps:
  - name: "ビルドバリアント"
    parallel:
      - name: "デバッグビルド"
        command: "build.sh"
        env:
          BUILD_TYPE: "debug"
      - name: "リリースビルド"
        command: "build.sh"
        env:
          BUILD_TYPE: "release"
```

### 環境変数の統合

環境変数と状態変数を混在：

```yaml
steps:
  - name: "デプロイ"
    command: "deploy.sh"
    env:
      APP_VERSION: "`{{.version}}`"  # 状態から
      API_KEY: "${DEPLOY_API_KEY}"  # 環境から
```

## 状態テンプレート

### 文字列操作

```yaml
steps:
  - name: "メッセージをフォーマット"
    command: |
      echo "`{{.message | upper}}`"
      echo "`{{.path | base}}`"
      echo "`{{.text | replace \" \" \"_\"}}`"
```

### 算術演算

```yaml
steps:
  - name: "計算"
    command: |
      echo "合計: `{{add .value1 .value2}}`"
      echo "差: `{{sub .value1 .value2}}`"
```

### 日付と時刻

```yaml
steps:
  - name: "タイムスタンプ"
    command: "echo `{{now | date \"2006-01-02 15:04:05\"}}`"
    outputs:
      - name: build_time
        from: stdout
```

## ベストプラクティス

### 1. 変数の命名

明確で説明的な名前を使用：
```yaml
# 良い例
outputs:
  - name: docker_image_tag
  - name: deployment_url
  - name: test_coverage_percent

# 避けるべき例
outputs:
  - name: tag
  - name: url
  - name: coverage
```

### 2. 状態の検証

使用前に状態を検証：
```yaml
steps:
  - name: "前提条件を確認"
    command: |
      if [ -z "`{{.api_key}}`" ]; then
        echo "エラー: api_keyが設定されていません"
        exit 1
      fi
```

### 3. 状態のドキュメント化

期待される状態変数をドキュメント化：
```yaml
# このプランは以下の変数を期待します：
# - environment: ターゲット環境（dev、staging、prod）
# - version: デプロイするアプリケーションバージョン
# - region: デプロイ用のAWSリージョン

requires:
  - environment
  - version
  - region
```

### 4. 状態のクリーンアップ

機密データをクリーン：
```yaml
steps:
  - name: "シークレットをクリーンアップ"
    command: "echo ''"
    outputs:
      - name: api_key
        value: ""  # 機密データをクリア
    always_run: true
```

## 状態のデバッグ

### 現在の状態を表示

```bash
# すべての状態を表示
plexr status

# 特定の変数を表示
plexr status --var version

# 状態をJSONとしてエクスポート
plexr status --json > state-backup.json
```

### 状態履歴

状態の変更を追跡：
```yaml
steps:
  - name: "状態変更をログ"
    command: |
      echo "以前: `{{.version | default \"none\"}}`"
      echo "新規: `{{.new_version}}`"
    debug: true
```

### インタラクティブな状態更新

開発中に状態を変更：
```bash
# 変数を設定
plexr state set version "1.2.3"

# 変数を削除
plexr state unset debug_mode

# 状態をインポート
plexr state import < custom-state.json
```

## 一般的なパターン

### フィーチャーフラグ

```yaml
vars:
  features:
    new_ui: true
    beta_api: false
    
steps:
  - name: "機能付きでデプロイ"
    command: "deploy.sh"
    env:
      ENABLE_NEW_UI: "`{{.features.new_ui}}`"
      ENABLE_BETA_API: "`{{.features.beta_api}}`"
```

### ビルドマトリックス

```yaml
vars:
  platforms: ["linux", "darwin", "windows"]
  architectures: ["amd64", "arm64"]
  
steps:
  - name: "ビルドマトリックス"
    command: "build.sh -os `{{.platform}}` -arch `{{.arch}}`"
    for_each:
      platforms: "`{{.platforms}}`"
      architectures: "`{{.architectures}}`"
      as:
        platform: platform
        arch: arch
```

### ワークフロー状態マシン

```yaml
steps:
  - name: "状態を確認"
    command: "get-workflow-state.sh"
    outputs:
      - name: workflow_state
        
  - name: "保留中を処理"
    command: "process.sh"
    condition: '`{{.workflow_state}}` == "pending"'
    
  - name: "エラーを処理"
    command: "error-handler.sh"
    condition: '`{{.workflow_state}}` == "error"'
```