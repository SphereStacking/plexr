# はじめに

このガイドでは、Plexrを数分で使い始める方法を説明します。

## Plexrとは？

Plexr（Plan + Executor）は、YAMLベースの実行プランを通じてローカル開発環境のセットアップを自動化するCLIツールです。「私のマシンでは動く」問題を解決し、チーム全体で一貫した再現可能な環境を確保します。

## 前提条件

- Go 1.21以上（ソースからのインストールの場合）
- YAMLの基本的な知識
- コマンドラインターミナル

## インストール

### ソースから

Goを使用する最も簡単な方法：

```bash
go install github.com/SphereStacking/plexr@latest
```

### バイナリリリース

近日公開予定！主要プラットフォーム向けのバイナリリリースが利用可能になります。

## 最初のプラン

Plexrの動作を理解するために、シンプルな実行プランを作成しましょう。

### 1. プランファイルの作成

`setup.yml`という名前のファイルを作成：

```yaml
name: "Hello Plexr"
version: "1.0.0"
description: "初めてのPlexr実行プラン"

executors:
  shell:
    type: shell

steps:
  - id: welcome
    description: "ウェルカムメッセージ"
    executor: shell
    files:
      - path: "scripts/welcome.sh"
```

### 2. スクリプトの作成

スクリプト用のディレクトリを作成し、ウェルカムスクリプトを追加：

```bash
mkdir scripts
```

`scripts/welcome.sh`を作成：

```bash
#!/bin/bash
echo "🚀 Plexrへようこそ！"
echo "開発環境自動化の旅がここから始まります。"
```

### 3. プランの実行

プランを実行：

```bash
plexr execute setup.yml
```

以下のような出力が表示されます：

```
🚀 "Hello Plexr"の実行を開始
✓ ステップ 1/1: ウェルカムメッセージ
🎉 実行が正常に完了しました！
```

## プランの理解

Plexrプランは以下で構成されます：

- **メタデータ**: 名前、バージョン、説明
- **エグゼキューター**: 異なるタイプのファイルの実行方法を定義
- **ステップ**: 実行する実際のタスク

### ステップの構造

各ステップには以下があります：
- `id`: 一意の識別子
- `description`: 人間が読める説明
- `executor`: 使用するエグゼキューター
- `files`: 実行するファイルのリスト

## 主要機能

### 1. ステート管理

Plexrは実行状態を追跡するため、何かが失敗した場合は修正して再開できます：

```bash
# 現在の状態を確認
plexr status

# 失敗から再開
plexr execute setup.yml
```

### 2. プラットフォームサポート

異なるオペレーティングシステムを優雅に処理：

```yaml
files:
  - path: "scripts/install.sh"
    platform: linux
  - path: "scripts/install.ps1"
    platform: windows
```

### 3. 依存関係

依存関係で実行順序を定義：

```yaml
steps:
  - id: install_tools
    description: "必要なツールのインストール"
    executor: shell
    files:
      - path: "scripts/install.sh"
  
  - id: configure_tools
    description: "ツールの設定"
    executor: shell
    depends_on: [install_tools]
    files:
      - path: "scripts/configure.sh"
```

### 4. スキップ条件

条件に基づいてステップをスキップ：

```yaml
steps:
  - id: install_docker
    description: "Dockerのインストール"
    executor: shell
    check_command: "docker --version"
    files:
      - path: "scripts/install_docker.sh"
```

## 次のステップ

基本を理解したら：

1. 詳細なYAMLオプションについて[設定ガイド](/ja/guide/configuration)を読む
2. 実世界のパターンについて[サンプル](/ja/examples/)を探索
3. 高度な使用法について[コマンド](/ja/guide/commands)を学ぶ
4. 複雑なワークフローのための[ステート管理](/ja/guide/state-management)を理解

## ヘルプを得る

- [トラブルシューティングガイド](/ja/guide/troubleshooting)を確認
- [GitHub](https://github.com/SphereStacking/plexr/issues)でイシューを開く
- [FAQ](/ja/guide/faq)を読む