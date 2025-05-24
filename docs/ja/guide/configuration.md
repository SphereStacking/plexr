# 設定

PlexrはYAMLファイルを使用して実行プランを定義します。このガイドでは、すべての設定オプションについて詳しく説明します。

## ファイル構造

Plexr設定ファイルには、以下のトップレベル構造があります：

```yaml
name: string          # 必須: プラン名
version: string       # 必須: プランのバージョン
description: string   # オプション: プランの説明
platforms: map        # オプション: プラットフォーム固有の設定
executors: map        # 必須: エグゼキューター設定
steps: array          # 必須: 実行ステップ
```

## 基本的な例

```yaml
name: "開発環境セットアップ"
version: "1.0.0"
description: |
  以下を含む完全な開発環境をセットアップします：
  - 開発ツール
  - データベース
  - 設定ファイル

executors:
  shell:
    type: shell
    config:
      shell: /bin/bash

steps:
  - id: install_tools
    description: "開発ツールをインストール"
    executor: shell
    files:
      - path: "scripts/install.sh"
```

## メタデータフィールド

### name（必須）

実行プランの名前：

```yaml
name: "マイプロジェクトセットアップ"
```

### version（必須）

プランのセマンティックバージョン：

```yaml
version: "1.0.0"
```

### description（オプション）

複数行テキストをサポートする詳細な説明：

```yaml
description: |
  このプランは以下をセットアップします：
  - Node.js環境
  - PostgreSQLデータベース
  - Redisキャッシュ
```

## エグゼキューター

エグゼキューターは、異なるタイプのファイルがどのように実行されるかを定義します。

### シェルエグゼキューター

組み込みのシェルエグゼキューターはシェルスクリプトを実行します：

```yaml
executors:
  shell:
    type: shell
    config:
      shell: /bin/bash      # Linux/macOS
      # shell: powershell.exe  # Windows
      timeout: 300          # デフォルトのタイムアウト（秒）
      env:                  # 環境変数
        NODE_ENV: development
        DEBUG: "true"
```

### カスタムエグゼキューター

将来のバージョンではカスタムエグゼキューターをサポートします：

```yaml
executors:
  sql:
    type: sql
    config:
      driver: postgres
      host: localhost
      port: 5432
      database: myapp
      user: postgres
```

## ステップ

ステップは順番に実行するタスクを定義します。

### 基本的なステップ

```yaml
steps:
  - id: create_directories
    description: "プロジェクトディレクトリを作成"
    executor: shell
    files:
      - path: "scripts/create_dirs.sh"
```

### ステップフィールド

#### id（必須）

ステップの一意の識別子：

```yaml
id: install_dependencies
```

#### description（必須）

人間が読める説明：

```yaml
description: "Node.js依存関係をインストール"
```

#### executor（必須）

使用するエグゼキューター：

```yaml
executor: shell
```

#### files（必須）

実行するファイルのリスト：

```yaml
files:
  - path: "scripts/install.sh"
    timeout: 600
    retry: 3
```

#### depends_on（オプション）

先に完了する必要がある依存関係：

```yaml
depends_on: [install_tools, create_directories]
```

#### skip_if（オプション）

ステップをスキップするかどうかを確認するシェルコマンド：

```yaml
skip_if: "test -f /usr/local/bin/node"
```

#### check_command（オプション）

ステップがすでに完了しているかを確認するコマンド：

```yaml
check_command: "docker --version"
```

## ファイル設定

ステップ内の各ファイルには追加の設定が可能です：

```yaml
files:
  - path: "scripts/install.sh"
    platform: linux         # プラットフォーム固有のファイル
    timeout: 300           # タイムアウトを上書き（秒）
    retry: 3              # 失敗時のリトライ回数
    skip_if: "test -f /usr/local/bin/tool"
```

### プラットフォーム固有のファイル

異なるオペレーティングシステムに対応：

```yaml
files:
  - path: "scripts/install.sh"
    platform: linux
  - path: "scripts/install.sh"
    platform: darwin    # macOS
  - path: "scripts/install.ps1"
    platform: windows
```

プラットフォーム値：
- `linux`: Linuxシステム
- `darwin`: macOS
- `windows`: Windows

## プラットフォーム設定

プラットフォーム固有の変数を定義：

```yaml
platforms:
  linux:
    install_prefix: /usr/local
    package_manager: apt
  darwin:
    install_prefix: /opt/homebrew
    package_manager: brew
  windows:
    install_prefix: C:\Program Files
    package_manager: choco
```

スクリプトでのアクセス：
```bash
echo "インストール先: ${PLEXR_PLATFORM_install_prefix}"
```

## 高度な機能

### トランザクションモード

アトミック実行のためのステップのグループ化（近日公開）：

```yaml
steps:
  - id: database_migration
    description: "データベースマイグレーションを実行"
    executor: sql
    transaction_mode: all  # all、none、またはper-file
    files:
      - path: "migrations/001_create_tables.sql"
      - path: "migrations/002_add_indexes.sql"
```

### 条件付き実行

条件に基づいてステップをスキップ：

```yaml
steps:
  - id: install_docker
    description: "Dockerが存在しない場合はインストール"
    executor: shell
    skip_if: "command -v docker >/dev/null 2>&1"
    files:
      - path: "scripts/install_docker.sh"
```

### 依存関係チェーン

複雑なワークフローを作成：

```yaml
steps:
  - id: install_node
    description: "Node.jsをインストール"
    executor: shell
    files:
      - path: "scripts/install_node.sh"
  
  - id: install_npm_packages
    description: "NPMパッケージをインストール"
    executor: shell
    depends_on: [install_node]
    files:
      - path: "scripts/npm_install.sh"
  
  - id: build_project
    description: "プロジェクトをビルド"
    executor: shell
    depends_on: [install_npm_packages]
    files:
      - path: "scripts/build.sh"
```

## 環境変数

Plexrはスクリプトにいくつかの環境変数を提供します：

- `PLEXR_STEP_ID`: 現在のステップID
- `PLEXR_STEP_INDEX`: 現在のステップ番号
- `PLEXR_TOTAL_STEPS`: ステップの総数
- `PLEXR_PLATFORM`: 現在のプラットフォーム（linux、darwin、windows）
- `PLEXR_STATE_FILE`: 状態ファイルへのパス
- `PLEXR_DRY_RUN`: ドライランモードの場合は"true"

## ベストプラクティス

### 1. 説明的なIDを使用

```yaml
# 良い例
id: install_postgresql_14

# 悪い例
id: step1
```

### 2. 前提条件を確認

```yaml
steps:
  - id: configure_git
    description: "Git設定を構成"
    check_command: "git config --get user.name"
    files:
      - path: "scripts/git_config.sh"
```

### 3. エラーを適切に処理

```bash
#!/bin/bash
set -euo pipefail  # エラー時に終了

# 前提条件を確認
if ! command -v node &> /dev/null; then
    echo "エラー: Node.jsが必要ですがインストールされていません"
    exit 1
fi
```

### 4. スクリプトをべき等にする

```bash
# ディレクトリが存在しない場合のみ作成
mkdir -p "$HOME/projects"

# インストールされていない場合のみインストール
if ! command -v tool &> /dev/null; then
    install_tool
fi
```

### 5. プラットフォーム検出を使用

```yaml
files:
  - path: "scripts/install_common.sh"
  - path: "scripts/install_mac.sh"
    platform: darwin
  - path: "scripts/install_linux.sh"
    platform: linux
```

## 検証

実行前に設定を検証：

```bash
plexr validate setup.yml
```

これは以下をチェックします：
- YAML構文エラー
- 必須フィールド
- 循環依存
- ファイルの存在
- エグゼキューターの可用性

## 次のステップ

- 実際の設定の[サンプル](/examples/)を参照
- これらの設定を使用する[コマンド](/guide/commands)について学ぶ
- 複雑なワークフローのための[状態管理](/guide/state-management)を理解する