# 設定スキーマ

Plexr YAML設定ファイルの完全なリファレンス。

## スキーマ概要

```yaml
# 必須フィールド
name: string
version: string
executors: map<string, ExecutorConfig>
steps: array<Step>

# オプションフィールド
description: string
platforms: map<string, map<string, string>>
```

## ルートフィールド

### name

**型:** `string`（必須）  
**説明:** 実行プランの名前

```yaml
name: "開発環境セットアップ"
```

### version

**型:** `string`（必須）  
**説明:** プランのセマンティックバージョン

```yaml
version: "1.0.0"
```

### description

**型:** `string`（オプション）  
**説明:** プランの詳細な説明

```yaml
description: |
  このプランは以下を含む完全な開発環境をセットアップします：
  - 必要なツールと依存関係
  - データベースの初期化
  - 設定ファイル
```

### executors

**型:** `map<string, ExecutorConfig>`（必須）  
**説明:** エグゼキューター設定のマップ

```yaml
executors:
  shell:
    type: shell
    config:
      shell: /bin/bash
```

### steps

**型:** `array<Step>`（必須）  
**説明:** 実行ステップのリスト

```yaml
steps:
  - id: install_tools
    description: "必要なツールをインストール"
    executor: shell
    files:
      - path: "scripts/install.sh"
```

### platforms

**型:** `map<string, map<string, string>>`（オプション）  
**説明:** プラットフォーム固有の変数

```yaml
platforms:
  linux:
    package_manager: apt
    install_prefix: /usr/local
  darwin:
    package_manager: brew
    install_prefix: /opt/homebrew
```

## ExecutorConfig

エグゼキューターの設定。

### type

**型:** `string`（必須）  
**説明:** エグゼキューターのタイプ  
**値:** `shell`、`sql`（将来）、`http`（将来）

### config

**型:** `map<string, any>`（オプション）  
**説明:** エグゼキューター固有の設定

#### シェルエグゼキューター設定

```yaml
executors:
  shell:
    type: shell
    config:
      shell: string          # 使用するシェル（デフォルト: /bin/bashまたはpowershell.exe）
      timeout: integer       # デフォルトのタイムアウト（秒）（デフォルト: 300）
      env: map<string, string>  # 環境変数
      working_dir: string    # 作業ディレクトリ
```

例：
```yaml
executors:
  shell:
    type: shell
    config:
      shell: /bin/zsh
      timeout: 600
      env:
        NODE_ENV: development
        DEBUG: "true"
      working_dir: /tmp
```

## Step

単一の実行ステップの定義。

### id

**型:** `string`（必須）  
**説明:** ステップの一意の識別子  
**パターン:** `^[a-zA-Z][a-zA-Z0-9_-]*$`

```yaml
id: install_dependencies
```

### description

**型:** `string`（必須）  
**説明:** 人間が読める説明

```yaml
description: "Node.js依存関係をインストール"
```

### executor

**型:** `string`（必須）  
**説明:** 使用するエグゼキューターの名前  
**一致する必要:** `executors`マップのキー

```yaml
executor: shell
```

### files

**型:** `array<FileConfig>`（必須）  
**説明:** 実行するファイルのリスト

```yaml
files:
  - path: "scripts/install.sh"
    timeout: 300
```

### depends_on

**型:** `array<string>`（オプション）  
**説明:** このステップの前に完了する必要があるステップID

```yaml
depends_on: [install_tools, create_directories]
```

### skip_if

**型:** `string`（オプション）  
**説明:** ステップをスキップするかどうかを判断するシェルコマンド

```yaml
skip_if: "test -f /usr/local/bin/node"
```

### check_command

**型:** `string`（オプション）  
**説明:** ステップがすでに完了しているかを確認するコマンド

```yaml
check_command: "docker --version"
```

### transaction_mode

**型:** `string`（オプション）  
**説明:** トランザクション処理モード  
**値:** `all`、`none`、`per-file`  
**デフォルト:** `none`

```yaml
transaction_mode: all
```

## FileConfig

実行するファイルの設定。

### path

**型:** `string`（必須）  
**説明:** プランの場所を基準としたファイルへのパス

```yaml
path: "scripts/install.sh"
```

### platform

**型:** `string`（オプション）  
**説明:** プラットフォーム制限  
**値:** `linux`、`darwin`、`windows`

```yaml
platform: linux
```

### timeout

**型:** `integer`（オプション）  
**説明:** 実行タイムアウト（秒）  
**デフォルト:** エグゼキューターのデフォルトタイムアウト

```yaml
timeout: 600
```

### retry

**型:** `integer`（オプション）  
**説明:** 失敗時のリトライ試行回数  
**デフォルト:** 0

```yaml
retry: 3
```

### skip_if

**型:** `string`（オプション）  
**説明:** ファイルをスキップするかどうかを判断するシェルコマンド

```yaml
skip_if: "test -d node_modules"
```

## 完全な例

```yaml
name: "フルスタック開発環境"
version: "2.1.0"
description: |
  以下を含む完全なフルスタック開発のセットアップ：
  - Node.jsとnpmパッケージ
  - PostgreSQLデータベース
  - Redisキャッシュ
  - Dockerコンテナ

platforms:
  linux:
    package_manager: apt
    postgres_package: postgresql-14
  darwin:
    package_manager: brew
    postgres_package: postgresql@14

executors:
  shell:
    type: shell
    config:
      timeout: 300
      env:
        NODE_ENV: development
  
  sql:
    type: sql
    config:
      driver: postgres
      database: myapp_dev

steps:
  - id: install_system_deps
    description: "システム依存関係をインストール"
    executor: shell
    files:
      - path: "scripts/install_deps.sh"
        platform: linux
      - path: "scripts/install_deps_mac.sh"
        platform: darwin

  - id: install_node
    description: "Node.jsをインストール"
    executor: shell
    depends_on: [install_system_deps]
    check_command: "node --version"
    files:
      - path: "scripts/install_node.sh"
        timeout: 600

  - id: install_npm_packages
    description: "npmパッケージをインストール"
    executor: shell
    depends_on: [install_node]
    skip_if: "test -d node_modules"
    files:
      - path: "scripts/npm_install.sh"
        retry: 3

  - id: setup_database
    description: "PostgreSQLデータベースを初期化"
    executor: sql
    depends_on: [install_system_deps]
    transaction_mode: all
    files:
      - path: "sql/create_database.sql"
      - path: "sql/create_tables.sql"
      - path: "sql/seed_data.sql"

  - id: configure_environment
    description: "環境設定をセットアップ"
    executor: shell
    depends_on: [install_npm_packages, setup_database]
    files:
      - path: "scripts/setup_env.sh"
```

## 検証ルール

1. **一意のステップID:** すべてのステップIDは一意でなければならない
2. **有効な依存関係:** `depends_on`のステップは存在する必要がある
3. **循環依存なし:** 依存関係グラフは非環式でなければならない
4. **既存のエグゼキューター:** ステップのエグゼキューターは定義されている必要がある
5. **ファイルパス:** パスはプランファイルに対して相対的であるべき
6. **プラットフォーム値:** `linux`、`darwin`、または`windows`でなければならない

## スクリプト内の環境変数

Plexrは実行されるスクリプトに以下の変数を提供します：

| 変数 | 説明 |
|----------|-------------|
| `PLEXR_STEP_ID` | 現在のステップID |
| `PLEXR_STEP_INDEX` | 現在のステップ番号（1ベース） |
| `PLEXR_TOTAL_STEPS` | ステップの総数 |
| `PLEXR_PLATFORM` | 現在のプラットフォーム |
| `PLEXR_STATE_FILE` | 状態ファイルへのパス |
| `PLEXR_DRY_RUN` | ドライランモードの場合は"true" |
| `PLEXR_PLATFORM_*` | プラットフォーム固有の変数 |

スクリプトでの使用例：
```bash
echo "ステップ $PLEXR_STEP_INDEX / $PLEXR_TOTAL_STEPS を実行中"
echo "パッケージマネージャー: $PLEXR_PLATFORM_package_manager"
```