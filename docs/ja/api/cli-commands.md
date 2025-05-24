# CLIコマンドリファレンス

これは、すべてのPlexr CLIコマンドの包括的なリファレンスです。

## コマンド構造

```
plexr [global-flags] <command> [command-flags] [arguments]
```

## コマンド

### plexr execute

プランファイルを実行します。

```bash
plexr execute <plan-file> [flags]
```

#### 引数

- `<plan-file>` - YAMLプランファイルへのパス（必須）

#### フラグ

| フラグ | 型 | デフォルト | 説明 |
|------|------|---------|-------------|
| `--dry-run, -n` | bool | false | 変更を加えずに実行をプレビュー |
| `--auto, -y` | bool | false | すべてのプロンプトを自動的に確認 |
| `--platform, -p` | string | auto | プラットフォーム検出を上書き（linux、darwin、windows） |
| `--state-file, -s` | string | .plexr_state.json | 状態ファイルへのパス |
| `--force, -f` | bool | false | 完了したステップの再実行を強制 |
| `--verbose, -v` | bool | false | 詳細な出力を有効化 |

#### 例

```bash
# 基本的な実行
plexr execute setup.yml

# ドライラン
plexr execute setup.yml --dry-run

# 強制再実行
plexr execute setup.yml --force

# カスタム状態ファイル
plexr execute setup.yml --state-file=/tmp/state.json
```

#### 終了コード

- `0` - 実行成功
- `1` - 一般的なエラー
- `4` - 実行失敗
- `130` - ユーザーによる中断（Ctrl+C）

---

### plexr validate

実行せずにプランファイルを検証します。

```bash
plexr validate <plan-file> [flags]
```

#### 引数

- `<plan-file>` - YAMLプランファイルへのパス（必須）

#### フラグ

| フラグ | 型 | デフォルト | 説明 |
|------|------|---------|-------------|
| `--check-files, -c` | bool | false | 参照されているファイルが存在することを確認 |
| `--verbose, -v` | bool | false | 詳細な検証情報を表示 |

#### 例

```bash
# 基本的な検証
plexr validate setup.yml

# ファイルの存在を確認
plexr validate setup.yml --check-files

# 詳細な出力
plexr validate setup.yml --verbose
```

#### 終了コード

- `0` - 有効なプラン
- `3` - 検証失敗
- `2` - 無効な引数

---

### plexr status

現在の実行ステータスを表示します。

```bash
plexr status <plan-file> [flags]
```

#### 引数

- `<plan-file>` - YAMLプランファイルへのパス（必須）

#### フラグ

| フラグ | 型 | デフォルト | 説明 |
|------|------|---------|-------------|
| `--state-file, -s` | string | .plexr_state.json | 状態ファイルへのパス |
| `--verbose, -v` | bool | false | 詳細なステータス情報を表示 |
| `--json, -j` | bool | false | JSON形式で出力 |

#### 例

```bash
# ステータスを表示
plexr status setup.yml

# JSON出力
plexr status setup.yml --json

# カスタム状態ファイル
plexr status setup.yml --state-file=/tmp/state.json
```

#### 出力形式（JSON）

```json
{
  "plan": {
    "name": "開発環境セットアップ",
    "version": "1.0.0"
  },
  "state": "in_progress",
  "progress": {
    "completed": 3,
    "total": 5,
    "percentage": 60
  },
  "current_step": "configure_app",
  "steps": [
    {
      "id": "install_tools",
      "status": "completed",
      "description": "開発ツールをインストール"
    }
  ],
  "last_updated": "2023-12-15T10:30:45Z"
}
```

---

### plexr reset

実行状態をリセットします。

```bash
plexr reset <plan-file> [flags]
```

#### 引数

- `<plan-file>` - YAMLプランファイルへのパス（必須）

#### フラグ

| フラグ | 型 | デフォルト | 説明 |
|------|------|---------|-------------|
| `--state-file, -s` | string | .plexr_state.json | 状態ファイルへのパス |
| `--force, -f` | bool | false | 確認プロンプトをスキップ |
| `--steps` | string | | リセットするステップIDのカンマ区切りリスト |

#### 例

```bash
# 確認付きでリセット
plexr reset setup.yml

# 強制リセット
plexr reset setup.yml --force

# 特定のステップをリセット
plexr reset setup.yml --steps=install_tools,setup_database
```

---

### plexr completion

シェル補完スクリプトを生成します。

```bash
plexr completion <shell>
```

#### 引数

- `<shell>` - ターゲットシェル: bash、zsh、fish、またはpowershell（必須）

#### 例

```bash
# Bash
plexr completion bash > /etc/bash_completion.d/plexr

# Zsh
plexr completion zsh > "${fpath[1]}/_plexr"

# Fish
plexr completion fish > ~/.config/fish/completions/plexr.fish

# PowerShell
plexr completion powershell | Out-String | Invoke-Expression
```

---

### plexr version

バージョン情報を表示します。

```bash
plexr version [flags]
```

#### フラグ

| フラグ | 型 | デフォルト | 説明 |
|------|------|---------|-------------|
| `--verbose, -v` | bool | false | 詳細なバージョン情報を表示 |
| `--json, -j` | bool | false | JSON形式で出力 |

#### 例

```bash
# シンプルなバージョン
plexr version

# 詳細なバージョン
plexr version --verbose

# JSON出力
plexr version --json
```

#### 出力形式（JSON）

```json
{
  "version": "1.0.0",
  "go_version": "go1.21.5",
  "build_time": "2023-12-15T10:30:00Z",
  "git_commit": "abc123def",
  "platform": "linux/amd64"
}
```

---

### plexr help

ヘルプ情報を表示します。

```bash
plexr help [command]
```

#### 引数

- `[command]` - ヘルプを取得するオプションのコマンド

#### 例

```bash
# 一般的なヘルプ
plexr help

# コマンド固有のヘルプ
plexr help execute

# --helpフラグでも動作
plexr execute --help
```

## グローバルフラグ

これらのフラグはすべてのコマンドで利用可能です：

| フラグ | 型 | デフォルト | 説明 |
|------|------|---------|-------------|
| `--config, -c` | string | $HOME/.plexr/config.yml | 設定ファイルパス |
| `--log-level, -l` | string | info | ログレベル: debug、info、warn、error |
| `--no-color` | bool | false | カラー出力を無効化 |
| `--help, -h` | bool | false | ヘルプ情報を表示 |

## 環境変数

| 変数 | 説明 | デフォルト |
|----------|-------------|---------|
| `PLEXR_STATE_FILE` | デフォルトの状態ファイルの場所を上書き | .plexr_state.json |
| `PLEXR_LOG_LEVEL` | ログレベルを設定 | info |
| `PLEXR_NO_COLOR` | カラー出力を無効化 | false |
| `PLEXR_PLATFORM` | プラットフォーム検出を上書き | auto |
| `PLEXR_CONFIG` | 設定ファイルの場所を上書き | $HOME/.plexr/config.yml |

## 終了コード

| コード | 説明 |
|------|-------------|
| 0 | 成功 |
| 1 | 一般的なエラー |
| 2 | 無効な引数または使用法 |
| 3 | プラン検証失敗 |
| 4 | 実行失敗 |
| 5 | 状態ファイル破損 |
| 130 | ユーザーによる中断（Ctrl+C） |

## コマンドエイリアス

一部のコマンドは短いエイリアスをサポートします：

- `exec` → `execute`
- `val` → `validate`
- `stat` → `status`

例：
```bash
plexr exec setup.yml
```