# サンプル

実用的な例を通してPlexrの使い方を学びます。これらの例は、一般的なユースケースをカバーし、ベストプラクティスを示しています。

## 利用可能なサンプル

### [基本セットアップ](/examples/basic-setup)
一般的なツールを使用した基本的な開発環境のセットアップ方法を示すシンプルな例。

**学習内容:**
- 最初のプランの作成
- プラットフォーム固有のスクリプトの使用
- 基本的なステップの依存関係

### [高度なパターン](/examples/advanced-patterns)
複雑なセットアップのための高度なテクニック。

**学習内容:**
- 複雑な依存関係チェーン
- 条件付き実行
- エラー処理戦略
- 状態管理

### [実世界の例](/examples/real-world)
実際のプロジェクトからの本番環境対応の例。

**学習内容:**
- フルスタックアプリケーションのセットアップ
- データベースマイグレーション
- CI/CD統合
- チームコラボレーションパターン

## クイックスタートの例

始めるための最小限の例：

```yaml
name: "クイックスタート"
version: "1.0.0"

executors:
  shell:
    type: shell

steps:
  - id: hello_world
    description: "挨拶する"
    executor: shell
    files:
      - path: "hello.sh"
```

`hello.sh`：
```bash
#!/bin/bash
echo "Plexrからこんにちは！ 🚀"
```

実行：
```bash
plexr execute quickstart.yml
```

## 一般的なパターン

### 1. 検証付きツールインストール

```yaml
steps:
  - id: install_node
    description: "Node.jsをインストール"
    executor: shell
    check_command: "node --version"
    files:
      - path: "install_node.sh"
```

### 2. プラットフォーム固有のスクリプト

```yaml
steps:
  - id: install_deps
    description: "依存関係をインストール"
    executor: shell
    files:
      - path: "install_linux.sh"
        platform: linux
      - path: "install_mac.sh"
        platform: darwin
      - path: "install_windows.ps1"
        platform: windows
```

### 3. 順次依存関係

```yaml
steps:
  - id: step1
    description: "最初のステップ"
    executor: shell
    files:
      - path: "step1.sh"
  
  - id: step2
    description: "2番目のステップ"
    depends_on: [step1]
    executor: shell
    files:
      - path: "step2.sh"
  
  - id: step3
    description: "3番目のステップ"
    depends_on: [step2]
    executor: shell
    files:
      - path: "step3.sh"
```

### 4. 並列実行

```yaml
steps:
  - id: download_tools
    description: "ツールをダウンロード"
    executor: shell
    files:
      - path: "download.sh"
  
  - id: create_dirs
    description: "ディレクトリを作成"
    executor: shell
    files:
      - path: "mkdirs.sh"
  
  - id: setup_configs
    description: "設定をセットアップ"
    depends_on: [download_tools, create_dirs]
    executor: shell
    files:
      - path: "configure.sh"
```

## サンプルからのベストプラクティス

### 1. 常にチェックコマンドを使用

```yaml
check_command: "command -v docker >/dev/null 2>&1"
```

### 2. スクリプトをべき等にする

```bash
# 悪い例
mkdir ~/workspace

# 良い例
mkdir -p ~/workspace
```

### 3. エラーを適切に処理

```bash
set -euo pipefail

if ! command -v node &> /dev/null; then
    echo "Node.jsが必要ですがインストールされていません"
    exit 1
fi
```

### 4. 説明的なステップIDを使用

```yaml
# 悪い例
id: step1

# 良い例
id: install_postgresql_14
```

## サンプルリポジトリ構造

Plexrを使用する典型的なプロジェクト：

```
my-project/
├── setup.yml              # メインプランファイル
├── scripts/              # 実行スクリプト
│   ├── install/
│   │   ├── node.sh
│   │   ├── docker.sh
│   │   └── postgres.sh
│   ├── configure/
│   │   ├── git.sh
│   │   └── env.sh
│   └── verify/
│       └── health_check.sh
├── sql/                  # SQLスクリプト
│   ├── create_db.sql
│   └── migrations/
└── configs/              # 設定テンプレート
    ├── .env.template
    └── docker-compose.yml
```

## サンプルの貢献

素晴らしい例をお持ちですか？ぜひ含めたいと思います！

1. リポジトリをフォーク
2. `examples/`に例を追加
3. ユースケースを説明するREADMEを含める
4. プルリクエストを送信

## 次のステップ

- [基本セットアップの例](/examples/basic-setup)を試す
- [設定ガイド](/guide/configuration)を読む
- 詳細なオプションについて[APIリファレンス](/api/)を確認