# エグゼキューターガイド

このガイドでは、Plexrのエグゼキューターの使用方法と設定方法について説明します。ビルトインエグゼキューターと、特定のニーズに応じたカスタムエグゼキューターの作成方法を含みます。

## エグゼキューターとは？

エグゼキューターは、実際にステップを実行するコンポーネントです。プラン定義と実行環境の間の橋渡しをします。各エグゼキューターは、特定のタイプのコマンドや操作の実行に特化しています。

## ビルトインエグゼキューター

### シェルエグゼキューター

シェルエグゼキューターは、デフォルトで最も一般的に使用されるエグゼキューターです。システムのデフォルトシェルでシェルコマンドを実行します。

#### 基本的な使用方法

```yaml
steps:
  - name: "シンプルなコマンド"
    command: "echo Hello, World!"
    # executor: shell は暗黙的
```

#### 設定オプション

```yaml
steps:
  - name: "設定されたシェルコマンド"
    command: "npm install"
    executor: shell
    config:
      shell: "/bin/bash"  # シェルを指定（デフォルト: /bin/sh）
      timeout: 300        # タイムアウト（秒）
      workdir: "./app"    # 作業ディレクトリ
      env:               # 環境変数
        NODE_ENV: "production"
        CI: "true"
```

#### シェル機能

**コマンドチェーン**:
```yaml
steps:
  - name: "複数のコマンド"
    command: |
      echo "ビルド開始"
      npm install
      npm run build
      echo "ビルド完了"
```

**エラーハンドリング**:
```yaml
steps:
  - name: "安全なコマンド"
    command: |
      set -e  # エラー時に終了
      set -u  # 未定義変数で終了
      set -o pipefail  # パイプエラーで失敗
      
      command1 | command2 | command3
```

**条件付き実行**:
```yaml
steps:
  - name: "条件付きロジック"
    command: |
      if [ -f "package.json" ]; then
        npm install
      else
        echo "package.jsonが見つかりません"
      fi
```

### スクリプトエグゼキューター（将来）

様々な言語でスクリプトを実行:

```yaml
steps:
  - name: "Pythonスクリプト"
    executor: script
    config:
      language: python
      version: "3.9"
    command: |
      import json
      data = {"status": "success"}
      print(json.dumps(data))
```

### HTTPエグゼキューター（将来）

HTTPリクエストを実行:

```yaml
steps:
  - name: "ヘルスチェック"
    executor: http
    config:
      method: GET
      url: "https://api.example.com/health"
      timeout: 30
      expected_status: 200
```

## エグゼキューターの選択

### 自動選択

Plexrは、ステップの設定に基づいて適切なエグゼキューターを自動的に選択します:

```yaml
steps:
  - name: "シェルコマンド"
    command: "ls -la"  # シェルエグゼキューターを使用
    
  - name: "HTTPリクエスト"
    http:  # httpエグゼキューターを使用
      url: "https://example.com"
```

### 明示的な選択

エグゼキューターを明示的に指定:

```yaml
steps:
  - name: "特定のエグゼキューター"
    executor: shell
    command: "echo 'シェルエグゼキューターを使用'"
```

## エグゼキューターの設定

### グローバル設定

エグゼキューターを使用するすべてのステップのデフォルト設定を設定:

```yaml
executors:
  shell:
    timeout: 120
    env:
      LOG_LEVEL: "info"
      
steps:
  - name: "グローバル設定を使用"
    command: "build.sh"
    
  - name: "グローバル設定をオーバーライド"
    command: "long-task.sh"
    config:
      timeout: 600  # グローバルタイムアウトをオーバーライド
```

### ステップごとの設定

個々のステップのエグゼキューターを設定:

```yaml
steps:
  - name: "カスタム設定"
    command: "deploy.sh"
    executor: shell
    config:
      workdir: "/opt/app"
      timeout: 300
      env:
        DEPLOY_ENV: "production"
        DEPLOY_KEY: "<code>&#123;&#123;.deploy_key&#125;&#125;</code>"
```

## 出力の処理

### 出力のキャプチャ

```yaml
steps:
  - name: "出力をキャプチャ"
    command: "generate-report.sh"
    outputs:
      - name: report_path
        from: stdout
      - name: report_size
        from: stderr
        regex: "Size: ([0-9]+) bytes"
```

### 出力フォーマット

**プレーンテキスト**:
```yaml
outputs:
  - name: message
    from: stdout
```

**JSONパース**:
```yaml
outputs:
  - name: data
    from: stdout
    json_parse: true
  - name: specific_field
    from: stdout
    json_path: "$.result.id"
```

**行選択**:
```yaml
outputs:
  - name: first_line
    from: stdout
    line: 1
  - name: last_line
    from: stdout
    line: -1
```

## エラーハンドリング

### 終了コード

```yaml
steps:
  - name: "終了コードを処理"
    command: "test-command.sh"
    success_codes: [0, 1]  # 0と1を成功とみなす
    ignore_failure: false   # エラー時にプランを失敗させる
```

### リトライロジック

```yaml
steps:
  - name: "失敗時にリトライ"
    command: "flaky-service.sh"
    retry:
      attempts: 3
      delay: 5s
      backoff: exponential  # または linear
      on_codes: [1, 2]     # 特定の終了コードでのみリトライ
```

### エラーパターン

```yaml
steps:
  - name: "パターンベースのリトライ"
    command: "connect-to-service.sh"
    retry:
      attempts: 5
      delay: 10s
      on_output_contains: ["connection refused", "timeout"]
```

## 並列実行

### 並列ステップ

```yaml
steps:
  - name: "並列実行"
    parallel:
      - name: "フロントエンドテスト"
        command: "npm test"
        workdir: "./frontend"
        
      - name: "バックエンドテスト"
        command: "go test ./..."
        workdir: "./backend"
        
      - name: "統合テスト"
        command: "pytest"
        workdir: "./tests"
```

### 並列設定

```yaml
steps:
  - name: "制御された並列性"
    parallel:
      max_concurrent: 2  # 同時実行数を制限
      fail_fast: true    # 最初の失敗で停止
      steps:
        - name: "タスク1"
          command: "task1.sh"
        - name: "タスク2"
          command: "task2.sh"
        - name: "タスク3"
          command: "task3.sh"
```

## カスタムエグゼキューター

### カスタムエグゼキューターの作成

Executorインターフェースを実装:

```go
package myexecutor

import (
    "context"
    "github.com/plexr/plexr/internal/core"
)

type MyExecutor struct {
    // エグゼキューター固有のフィールド
}

func New() *MyExecutor {
    return &MyExecutor{}
}

func (e *MyExecutor) Execute(ctx context.Context, step core.Step, state core.State) error {
    // 実装
    return nil
}

func (e *MyExecutor) Validate(step core.Step) error {
    // ステップ設定を検証
    return nil
}
```

### カスタムエグゼキューターの登録

```go
// main.goまたはプラグイン内
import (
    "github.com/plexr/plexr/internal/core"
    "myproject/executors/myexecutor"
)

func init() {
    core.RegisterExecutor("myexecutor", myexecutor.New())
}
```

### カスタムエグゼキューターの使用

```yaml
steps:
  - name: "カスタムエグゼキューターステップ"
    executor: myexecutor
    config:
      custom_option: "value"
    command: "custom command"
```

## ベストプラクティス

### 1. 適切なエグゼキューターを選択

- システムコマンドとスクリプトにはシェルを使用
- 特定のタスクには専門のエグゼキューターを使用（HTTP、データベースなど）
- 複雑で再利用可能なロジックにはカスタムエグゼキューターを作成

### 2. エラーを適切に処理

```yaml
steps:
  - name: "堅牢な実行"
    command: |
      set -euo pipefail
      trap 'echo "エラーが行 $LINENO で発生"' ERR
      
      # コマンドをここに記述
    retry:
      attempts: 3
      delay: 5s
```

### 3. タイムアウトを使用

```yaml
steps:
  - name: "時間制限付き操作"
    command: "long-running-task.sh"
    config:
      timeout: 300  # 5分
    on_timeout: 
      command: "cleanup.sh"  # タイムアウト時に実行
```

### 4. 機密データを保護

```yaml
steps:
  - name: "安全な実行"
    command: "deploy.sh"
    env:
      API_KEY: "${API_KEY}"  # 環境から
    config:
      mask_output: ["password", "secret", "key"]
```

### 5. 適切にログを記録

```yaml
steps:
  - name: "ログ付き実行"
    command: "process.sh"
    config:
      log_level: "debug"  # debug, info, warn, error
      log_output: true    # コマンド出力をログ
      log_file: "process.log"
```

## エグゼキューターのトラブルシューティング

### デバッグモード

```yaml
steps:
  - name: "デバッグ実行"
    command: "problematic-command.sh"
    debug: true  # デバッグ出力を有効化
    config:
      verbose: true
      dry_run: true  # 実行内容を表示
```

### 実行コンテキスト

```yaml
steps:
  - name: "コンテキストを表示"
    command: |
      echo "作業ディレクトリ: $(pwd)"
      echo "ユーザー: $(whoami)"
      echo "シェル: $SHELL"
      echo "PATH: $PATH"
      env | sort
```

### 一般的な問題

**コマンドが見つからない**:
```yaml
steps:
  - name: "明示的なパス"
    command: "/usr/local/bin/custom-tool"
    # またはPATHを更新
    env:
      PATH: "/usr/local/bin:${PATH}"
```

**権限拒否**:
```yaml
steps:
  - name: "権限を確保"
    command: |
      chmod +x script.sh
      ./script.sh
```

**作業ディレクトリの問題**:
```yaml
steps:
  - name: "絶対パス"
    command: "build.sh"
    config:
      workdir: "${PWD}/subdir"  # 絶対パスを使用
```

## パフォーマンス最適化

### 出力バッファリング

```yaml
steps:
  - name: "大きな出力"
    command: "generate-large-output.sh"
    config:
      buffer_size: 1048576  # 1MBバッファ
      stream_output: true   # バッファリングの代わりにストリーム
```

### リソース制限

```yaml
steps:
  - name: "リソース制約"
    command: "memory-intensive.sh"
    config:
      memory_limit: "2G"
      cpu_limit: 2
      nice: 10  # 低優先度
```

### キャッシング

```yaml
steps:
  - name: "キャッシュされた実行"
    command: "expensive-operation.sh"
    cache:
      key: "<code>&#123;&#123;.cache_key&#125;&#125;</code>"
      ttl: 3600  # 1時間
      paths:
        - "./output"
        - "./artifacts"
```