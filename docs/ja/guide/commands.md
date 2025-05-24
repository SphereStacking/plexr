# コマンド

Plexrは、プランを管理および実行するためのいくつかのコマンドを提供します。このガイドでは、利用可能なすべてのコマンドとそのオプションについて説明します。

## 概要

```bash
plexr [command] [flags]
```

利用可能なコマンド：
- `execute` - 実行プランを実行
- `validate` - 実行せずにプランを検証
- `status` - 現在の実行ステータスを表示
- `reset` - 実行状態をリセット
- `completion` - シェル補完を生成
- `help` - 任意のコマンドのヘルプを取得
- `version` - バージョン情報を表示

## execute

実行プランを実行するメインコマンド。

### 基本的な使用法

```bash
plexr execute [plan-file] [flags]
```

### 例

```bash
# プランを実行
plexr execute setup.yml

# 何が起こるかを確認するドライラン
plexr execute setup.yml --dry-run

# すべてのプロンプトを自動確認
plexr execute setup.yml --auto

# 特定のプラットフォームを使用
plexr execute setup.yml --platform=linux
```

### フラグ

| フラグ | 短縮形 | 説明 | デフォルト |
|------|-------|-------------|---------|
| `--dry-run` | `-n` | 実行せずに何が実行されるかを表示 | `false` |
| `--auto` | `-y` | すべてのプロンプトを自動確認 | `false` |
| `--platform` | `-p` | プラットフォーム検出を上書き | 自動検出 |
| `--state-file` | `-s` | カスタム状態ファイルの場所 | `.plexr_state.json` |
| `--verbose` | `-v` | 詳細な出力を有効化 | `false` |
| `--force` | `-f` | 完了したステップの再実行を強制 | `false` |

### 実行フロー

1. プランをロードして検証
2. 現在の状態を確認
3. 依存関係を解決
4. ステップを順番に実行
5. 各ステップ後に状態を更新
6. エラーを適切に処理

### 失敗時の再開

実行が失敗した場合、コマンドを再度実行するだけで再開できます：

```bash
# 最初の実行がステップ3で失敗
plexr execute setup.yml
# エラー: ステップ3が失敗しました

# 問題を修正してから再開
plexr execute setup.yml
# ステップ3から再開中...
```

## validate

実行せずにプランファイルを検証します。

### 使用法

```bash
plexr validate [plan-file] [flags]
```

### 例

```bash
# 構文と構造を検証
plexr validate setup.yml

# 詳細な出力で検証
plexr validate setup.yml --verbose
```

### チェック内容

- YAML構文
- 必須フィールドの存在
- フィールドの型と値
- 循環依存
- ファイルの存在（`--check-files`を使用）
- エグゼキューターの可用性

### フラグ

| フラグ | 短縮形 | 説明 | デフォルト |
|------|-------|-------------|---------|
| `--check-files` | `-c` | 参照されているファイルが存在するか確認 | `false` |
| `--verbose` | `-v` | 詳細な検証情報を表示 | `false` |

## status

プランの現在の実行状態を表示します。

### 使用法

```bash
plexr status [plan-file] [flags]
```

### 例

```bash
# 現在の状態を表示
plexr status setup.yml

# 詳細な状態を表示
plexr status setup.yml --verbose
```

### 出力

```
プラン: 開発環境セットアップ (v1.0.0)
状態: 進行中

進捗: 3/5 ステップ完了 (60%)

✓ install_tools      - 開発ツールをインストール
✓ create_directories - プロジェクトディレクトリを作成  
✓ setup_database     - データベースを初期化
→ configure_app      - アプリケーションを構成（現在）
○ run_tests          - 検証テストを実行

最終更新: 2023-12-15 10:30:45
```

### フラグ

| フラグ | 短縮形 | 説明 | デフォルト |
|------|-------|-------------|---------|
| `--state-file` | `-s` | カスタム状態ファイルの場所 | `.plexr_state.json` |
| `--verbose` | `-v` | 詳細な状態を表示 | `false` |
| `--json` | `-j` | JSON形式で出力 | `false` |

## reset

実行状態をリセットして、最初からやり直すことができます。

### 使用法

```bash
plexr reset [plan-file] [flags]
```

### 例

```bash
# すべての状態をリセット
plexr reset setup.yml

# 確認なしでリセット
plexr reset setup.yml --force

# 特定のステップのみリセット
plexr reset setup.yml --steps install_tools,setup_database
```

### フラグ

| フラグ | 短縮形 | 説明 | デフォルト |
|------|-------|-------------|---------|
| `--state-file` | `-s` | カスタム状態ファイルの場所 | `.plexr_state.json` |
| `--force` | `-f` | 確認プロンプトをスキップ | `false` |
| `--steps` | | 特定のステップのみリセット | すべて |

## completion

シェル補完スクリプトを生成します。

### 使用法

```bash
plexr completion [shell]
```

### サポートされているシェル

- bash
- zsh
- fish
- powershell

### 例

```bash
# Bash
plexr completion bash > /etc/bash_completion.d/plexr

# Zsh
plexr completion zsh > "${fpath[1]}/_plexr"

# Fish
plexr completion fish > ~/.config/fish/completions/plexr.fish

# PowerShell
plexr completion powershell > plexr.ps1
```

## version

バージョン情報を表示します。

### 使用法

```bash
plexr version [flags]
```

### 例

```bash
# シンプルなバージョン
plexr version
# 出力: plexr version 1.0.0

# 詳細なバージョン情報
plexr version --verbose
# 出力:
# plexr version 1.0.0
# Go version: go1.21.5
# Built: 2023-12-15T10:30:00Z
# Commit: abc123def
```

### フラグ

| フラグ | 短縮形 | 説明 | デフォルト |
|------|-------|-------------|---------|
| `--verbose` | `-v` | 詳細なバージョン情報を表示 | `false` |
| `--json` | `-j` | JSON形式で出力 | `false` |

## グローバルフラグ

これらのフラグはすべてのコマンドで利用可能です：

| フラグ | 短縮形 | 説明 | デフォルト |
|------|-------|-------------|---------|
| `--config` | `-c` | 設定ファイルの場所 | `$HOME/.plexr/config.yml` |
| `--log-level` | `-l` | ログレベル（debug、info、warn、error） | `info` |
| `--no-color` | | カラー出力を無効化 | `false` |
| `--help` | `-h` | コマンドのヘルプを表示 | |

## 環境変数

Plexrは以下の環境変数を使用します：

```bash
# 状態ファイルの場所を上書き
export PLEXR_STATE_FILE=/tmp/my-state.json

# ログレベルを設定
export PLEXR_LOG_LEVEL=debug

# カラーを無効化
export PLEXR_NO_COLOR=true

# プラットフォームを上書き
export PLEXR_PLATFORM=linux
```

## 終了コード

Plexrは標準的な終了コードを使用します：

- `0`: 成功
- `1`: 一般的なエラー
- `2`: 無効な引数
- `3`: プラン検証失敗
- `4`: 実行失敗
- `5`: 状態の破損
- `130`: 中断（Ctrl+C）

## 高度な使用法

### コマンドの連鎖

```bash
# 検証してから、成功した場合に実行
plexr validate setup.yml && plexr execute setup.yml

# リセットして1行で実行
plexr reset setup.yml --force && plexr execute setup.yml --auto
```

### CI/CDでの使用

```bash
# CIフレンドリーな実行
plexr execute setup.yml \
  --auto \
  --platform=linux \
  --log-level=debug \
  --no-color
```

### デバッグ

```bash
# 最大の詳細度
PLEXR_LOG_LEVEL=debug plexr execute setup.yml --verbose

# 詳細な出力でドライラン
plexr execute setup.yml --dry-run --verbose
```

## 次のステップ

- [状態管理](/guide/state-management)について学ぶ
- 実際の使用例の[サンプル](/examples/)を参照
- Plexrを拡張するための[エグゼキューター](/guide/executors)について読む