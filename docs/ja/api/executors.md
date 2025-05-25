# エグゼキューター

エグゼキューターは、Plexrプランのステップを実行する責任を持つコンポーネントです。プラン定義と実際の実行環境の間のインターフェースを提供します。

## 組み込みエグゼキューター

### シェルエグゼキューター

シェルコマンドを実行するためのデフォルトエグゼキューター。

```yaml
steps:
  - name: "依存関係をインストール"
    command: "npm install"
    executor: shell  # オプション、shellがデフォルト
```

#### 設定

```yaml
steps:
  - name: "プロジェクトをビルド"
    command: "make build"
    executor: shell
    env:
      BUILD_ENV: production
    timeout: 300  # 秒
    workdir: ./src
```

### SQLエグゼキューター

PostgreSQLデータベースに対してSQLクエリを実行（MySQLとSQLiteのサポートは計画中）。

#### 設定

`executors`セクションでSQLエグゼキューターを定義：

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

#### 使用方法

```yaml
steps:
  - name: "マイグレーションを実行"
    executor: db
    files:
      - path: sql/001_schema.sql
        timeout: 30
    transaction_mode: all  # オプション: none, each, all
```

#### トランザクションモード

- `none`: トランザクションなし
- `each`: 各SQLステートメントが独自のトランザクション内で実行（デフォルト）
- `all`: すべてのステートメントが単一のトランザクション内で実行

#### 複数のデータベース

異なるデータベース用に複数のSQLエグゼキューターを定義できます：

```yaml
executors:
  main_db:
    type: sql
    driver: postgres
    database: myapp_main
    # ... 接続詳細
    
  analytics_db:
    type: sql
    driver: postgres
    database: myapp_analytics
    # ... 接続詳細
```

## カスタムエグゼキューター

`Executor`インターフェースを実装することで、カスタムエグゼキューターを作成できます：

```go
type Executor interface {
    Execute(ctx context.Context, step Step, state State) error
    Validate(step Step) error
}
```

### カスタムエグゼキューターの例

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
    // Dockerコンテナを実行するための実装
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
        return fmt.Errorf("dockerエグゼキューターには設定に'image'が必要です")
    }
    return nil
}
```

## エグゼキューター設定

### グローバル設定

プランでデフォルトのエグゼキューター設定を設定：

```yaml
config:
  executors:
    shell:
      timeout: 60
      shell: /bin/bash
    docker:
      registry: ghcr.io
      
steps:
  - name: "ビルド"
    command: "make build"
    # グローバルシェル設定を使用
```

### ステップごとの設定

特定のステップのエグゼキューター設定を上書き：

```yaml
steps:
  - name: "長時間実行タスク"
    command: "./long-task.sh"
    executor: shell
    config:
      timeout: 3600  # 1時間
```

## 状態管理

エグゼキューターは共有状態の読み書きができます：

```yaml
steps:
  - name: "バージョンを取得"
    command: "git describe --tags"
    executor: shell
    outputs:
      - name: version
        from: stdout
        
  - name: "バージョンでビルド"
    command: "make build VERSION=`{{.version}}`"
```

## エラー処理

### リトライ設定

```yaml
steps:
  - name: "不安定なテスト"
    command: "npm test"
    retry:
      attempts: 3
      delay: 5s
      backoff: exponential
```

### エラー条件

```yaml
steps:
  - name: "サービスを確認"
    command: "curl http://localhost:8080/health"
    errorConditions:
      - exitCode: [1, 2]
        retry: true
      - output: "connection refused"
        fail: true
```

## 並列実行

複数のステップを同時に実行：

```yaml
steps:
  - name: "並列タスク"
    parallel:
      - name: "フロントエンドをテスト"
        command: "npm test"
      - name: "バックエンドをテスト"
        command: "go test ./..."
      - name: "リント"
        command: "npm run lint"
```

## 環境変数

エグゼキューターは環境変数の展開をサポート：

```yaml
steps:
  - name: "デプロイ"
    command: "deploy.sh"
    env:
      DEPLOY_ENV: "`{{.environment}}`"
      API_KEY: "$DEPLOY_API_KEY"  # ホスト環境から
```

## ロギングと出力

### 出力のキャプチャ

```yaml
steps:
  - name: "レポートを生成"
    command: "./generate-report.sh"
    outputs:
      - name: reportPath
        from: stdout
      - name: reportSize
        from: stderr
        regex: "Size: ([0-9]+) bytes"
```

### ログレベル

```yaml
steps:
  - name: "詳細な操作"
    command: "npm install"
    logLevel: debug  # debug、info、warn、error
```

## セキュリティの考慮事項

- エグゼキューターはPlexrプロセスの権限で実行されます
- 機密データには環境変数を使用
- カスタムエグゼキューターですべての入力を検証
- 信頼できないコードにはサンドボックス環境の使用を検討

## ベストプラクティス

1. **適切なエグゼキューターを使用**: 各タスクに適したエグゼキューターを選択
2. **エラーを適切に処理**: 適切なリトライとエラー処理を設定
3. **状態を慎重に管理**: 複雑なワークフローには出力と条件を使用
4. **実行を監視**: 適切なタイムアウトとロギングを設定
5. **エグゼキューターをテスト**: 本番使用前にカスタムエグゼキューターを徹底的に検証