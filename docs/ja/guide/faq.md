# よくある質問

## 一般的な質問

### Plexrとは何ですか？

Plexrは、複雑な多段階プロセスを定義、管理、実行するためのコマンド実行およびワークフロー自動化ツールです。開発ワークフロー、CI/CDパイプライン、インフラストラクチャの自動化のために設計されています。

### PlexrはMakeとどう違いますか？

Makeはソフトウェアのビルドに最適ですが、Plexrは以下を提供します：
- **より良い状態管理**: ステップ間でデータを簡単に共有
- **モダンなYAML構文**: より読みやすくメンテナンスしやすい
- **組み込みの並列処理**: 複雑な構文なしでステップを同時実行
- **条件付き実行**: 状態に基づく動的なワークフロー
- **複数のエグゼキューター**: シェルコマンドに限定されない

### CI/CDでPlexrを使用できますか？

はい！PlexrはCI/CD環境でうまく動作するように設計されています：
- 環境間での一貫した実行
- 明確な依存関係管理
- 詳細なロギングとエラーレポート
- 複雑なワークフローのための状態の永続性

## インストール

### Plexrをインストールするにはどうすればよいですか？

```bash
# Goを使用
go install github.com/plexr/plexr/cmd/plexr@latest

# ソースから
git clone https://github.com/plexr/plexr
cd plexr
make install
```

### システム要件は何ですか？

- Go 1.19以降（ソースからのインストールの場合）
- Linux、macOS、またはWindows
- シェルアクセス（bash、sh、または同等のもの）

### Plexrをアップデートするにはどうすればよいですか？

```bash
# Goを使用
go install github.com/plexr/plexr/cmd/plexr@latest

# バージョンを確認
plexr version
```

## 設定

### プランファイルはどこに置くべきですか？

デフォルトでは、Plexrは現在のディレクトリで`plan.yml`または`plan.yaml`を探します。カスタムの場所を指定することもできます：

```bash
plexr execute -f path/to/myplan.yml
```

### 複数のプランファイルを持つことはできますか？

はい！ワークフローを整理できます：

```bash
# 異なる目的のための異なるプラン
plexr execute -f plans/build.yml
plexr execute -f plans/deploy.yml
plexr execute -f plans/test.yml
```

### プランに変数を渡すにはどうすればよいですか？

いくつかの方法があります：

```bash
# 環境変数
export API_KEY=secret
plexr execute

# コマンドライン変数
plexr execute -v environment=production -v version=1.2.3

# ファイルから
plexr execute --vars-file vars.json
```

## 実行

### 特定のステップのみを実行できますか？

はい、タグを使用するか、ステップを指定します：

```yaml
steps:
  - name: "ビルド"
    command: "make build"
    tags: ["build", "ci"]
    
  - name: "テスト"
    command: "make test"
    tags: ["test", "ci"]
```

```bash
# タグ付きステップのみを実行
plexr execute --tags build

# 特定のステップから実行
plexr execute --from "テスト"
```

### 実行の問題をデバッグするにはどうすればよいですか？

詳細なロギングを有効にします：

```bash
# 基本的なデバッグ情報
plexr execute -v

# 詳細なデバッグ情報
plexr execute -vv

# ドライラン（実行される内容を表示）
plexr execute --dry-run
```

### 失敗した実行を再開できますか？

はい：

```bash
# 最後に成功したステップから再開
plexr execute --resume

# 失敗したステップを再試行
plexr execute --retry-failed
```

## 状態管理

### 状態はどこに保存されますか？

状態はプロジェクトディレクトリの`.plexr/state.json`に保存されます。これをカスタマイズできます：

```bash
# カスタム状態ファイルを使用
plexr execute --state-file /tmp/mystate.json

# インメモリ状態を使用（永続性なし）
plexr execute --no-state-file
```

### ステップ間でデータを共有するにはどうすればよいですか？

出力を使用してデータをキャプチャ：

```yaml
steps:
  - name: "バージョンを取得"
    command: "git describe --tags"
    outputs:
      - name: version
        from: stdout
        
  - name: "バージョンを使用"
    command: "echo バージョン `{{.version}}` をビルド中"
```

### 外部データソースを使用できますか？

はい、ファイルやコマンドからデータをロードできます：

```yaml
vars:
  config: "`{{file \"config.json\" | json}}`"
  
steps:
  - name: "シークレットをロード"
    command: "vault read -format=json secret/myapp"
    outputs:
      - name: secrets
        from: stdout
        json_parse: true
```

## セキュリティ

### シークレットをどのように扱いますか？

環境変数を使用し、ハードコーディングを避けます：

```yaml
steps:
  - name: "デプロイ"
    command: "deploy.sh"
    env:
      API_KEY: "${API_KEY}"  # 環境から
      DB_PASS: "`{{.vault_secret}}`"  # 前のステップから
```

### 機密出力をマスクできますか？

はい：

```yaml
steps:
  - name: "設定を表示"
    command: "print-config.sh"
    mask_output: ["password", "api_key", "secret"]
```

### プランファイルをコミットしても安全ですか？

ベストプラクティスに従えば安全です：
- シークレットをハードコードしない
- 環境変数を使用
- 必要な変数をドキュメント化
- 状態ファイルに`.gitignore`を使用

## 高度な使用法

### カスタムエグゼキューターを作成できますか？

はい、Executorインターフェースを実装します：

```go
type Executor interface {
    Execute(ctx context.Context, step Step, state State) error
    Validate(step Step) error
}
```

詳細は[エグゼキューターガイド](./executors.md)を参照してください。

### 複雑な依存関係をどのように処理しますか？

依存関係グループと条件を使用：

```yaml
steps:
  - name: "準備"
    command: "prepare.sh"
    
  - name: "ビルドA"
    command: "build-a.sh"
    depends_on: ["準備"]
    
  - name: "ビルドB"
    command: "build-b.sh"
    depends_on: ["準備"]
    
  - name: "パッケージ"
    command: "package.sh"
    depends_on: ["ビルドA", "ビルドB"]
```

### Plexrをライブラリとして使用できますか？

はい：

```go
import (
    "github.com/plexr/plexr/internal/config"
    "github.com/plexr/plexr/internal/core"
)

plan, err := config.LoadPlan("plan.yml")
runner := core.NewRunner()
err = runner.Execute(ctx, plan)
```

## 統合

### PlexrはDockerと連携しますか？

はい、コンテナでコマンドを実行できます：

```yaml
steps:
  - name: "Dockerでビルド"
    command: |
      docker run --rm -v $(pwd):/app -w /app \
        node:16 npm run build
```

### Kubernetesと統合できますか？

はい：

```yaml
steps:
  - name: "K8sにデプロイ"
    command: |
      kubectl apply -f deployment.yaml
      kubectl wait --for=condition=available --timeout=300s \
        deployment/myapp
```

### PlexrをGitフックで使用するにはどうすればよいですか？

`.git/hooks/pre-commit`ファイルを作成：

```bash
#!/bin/bash
plexr execute -f .plexr/pre-commit.yml
```

## トラブルシューティング

### コマンドが見つからないエラー

コマンドがPATHにあることを確認するか、絶対パスを使用：

```yaml
steps:
  - name: "絶対パスを使用"
    command: "/usr/local/bin/mytool"
    
  - name: "PATHを更新"
    command: "mytool"
    env:
      PATH: "/usr/local/bin:${PATH}"
```

### 状態の破損

状態が破損した場合はリセット：

```bash
# 状態をリセット
plexr reset

# すべてのPlexrデータを削除
rm -rf .plexr
```

### パフォーマンスの問題

- 独立したステップには並列実行を使用
- 詳細なコマンドの出力キャプチャを制限
- 適切なタイムアウトを使用
- 大きなプランを小さなものに分割することを検討

## ベストプラクティス

### 大規模プロジェクトをどのように構造化すべきですか？

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
└── plan.yml  # メインオーケストレーション
```

### 状態ファイルをバージョン管理すべきですか？

いいえ、`.gitignore`に追加してください：

```gitignore
.plexr/
*.state.json
```

### プランを再利用可能にするにはどうすればよいですか？

変数とインクルードを使用：

```yaml
# common/database.yml
steps:
  - name: "データベースセットアップ"
    command: "`{{.db_setup_script}}`"
    env:
      DB_NAME: "`{{.database_name}}`"

# メインのplan.yml
includes:
  - common/database.yml
  
vars:
  database_name: "myapp"
  db_setup_script: "./scripts/setup-db.sh"
```

## ヘルプを得る

### どこでサポートを受けられますか？

- GitHubイシュー: https://github.com/plexr/plexr/issues
- ドキュメント: https://plexr.dev/docs
- コミュニティDiscord: https://discord.gg/plexr

### バグを報告するにはどうすればよいですか？

以下を含めてください：
1. Plexrバージョン（`plexr version`）
2. プランファイル（サニタイズ済み）
3. エラーメッセージ
4. 再現手順
5. 期待される動作と実際の動作

### 貢献できますか？

はい！貢献を歓迎します：
- バグを報告
- 機能を提案
- プルリクエストを送信
- ドキュメントを改善
- 使用例を共有