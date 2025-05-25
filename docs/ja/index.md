---
layout: home

hero:
  name: "Plexr"
  text: "Plan + Executor"
  tagline: "ローカル開発環境の自動セットアップを実現する開発者フレンドリーなCLIツール"
  image:
    src: https://api.iconify.design/noto:rocket.svg
    alt: Plexr
  actions:
    - theme: brand
      text: はじめる
      link: /ja/guide/getting-started
    - theme: alt
      text: GitHubで見る
      link: https://github.com/SphereStacking/plexr

features:
  - icon: 📝
    title: 実行可能なドキュメント
    details: READMEのセットアップ手順をYAML形式の実行可能な設定に変換
  - icon: 🔄
    title: ステートフルな実行
    details: 失敗しても中断した場所から再開可能 - 最初からやり直す必要なし
  - icon: 🖥️
    title: クロスプラットフォーム
    details: macOS、Linux、Windowsで動作し、プラットフォーム固有の処理にも対応
  - icon: 🛡️
    title: 安全性重視
    details: ドライラン、スキップ条件、ロールバック機能で安全な操作を保証
---

## 🎉 最新リリース: v0.1.0

[GitHub Releases](https://github.com/SphereStacking/plexr/releases/tag/v0.1.0)から最新版をダウンロードするか、以下でインストール：

```bash
go install github.com/SphereStacking/plexr/cmd/plexr@v0.1.0
```

## クイックスタート

Plexrをインストールして、数分で使い始められます：

```bash
# 最新リリースをインストール
go install github.com/SphereStacking/plexr/cmd/plexr@latest

# またはプリビルドバイナリをダウンロード
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_$(uname -s)_$(uname -m | sed 's/x86_64/x86_64/;s/aarch64/arm64/').tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# 最初のプランを実行
plexr execute setup.yml
```

## なぜPlexr？

### よくある問題

- 😫 「READMEの通りにやったのに動かない」
- ⏰ 「開発環境のセットアップに丸一日かかった」
- 🤷 「私のマシンでは動くんだけど...」
- 🔧 「みんなの環境が微妙に違う」

### Plexrが提供する解決策

環境セットアップを以下のように改善：
- **再現性**: 毎回同じ結果を保証
- **デバッグ可能**: 明確なエラーメッセージと解決策の提示
- **保守性**: バージョン管理されたセットアップ手順
- **チーム対応**: 全員が同じ設定を使用

## 設定例

```yaml
name: "プロジェクトセットアップ"
version: "1.0.0"

steps:
  - id: install_deps
    description: "依存関係のインストール"
    executor: shell
    files:
      - path: "scripts/install.sh"
        platform: linux
      - path: "scripts/install.ps1"
        platform: windows

  - id: setup_database
    description: "データベースの初期化"
    executor: shell
    depends_on: [install_deps]
    files:
      - path: "scripts/db_setup.sh"
```

## 機能

### v0.1.0（最新リリース）
- ✅ 依存関係解決を備えたコア実行エンジン
- ✅ スクリプトとコマンド実行用のシェルエグゼキューター
- ✅ PostgreSQLサポートを含むSQLエグゼキューター
- ✅ 再開機能付きステート管理
- ✅ CLIコマンド（execute, validate, status, reset）
- ✅ 環境変数の展開
- ✅ プラットフォーム固有のファイル選択
- ✅ エラーハンドリングとロールバックサポート

### 今後の予定
- 🚧 追加データベースサポート（MySQL、SQLite）
- 🚧 API呼び出し用HTTPエグゼキューター
- 🚧 Dockerエグゼキューター
- 🚧 並列実行
- 🚧 高度な条件分岐ロジック

## さらに詳しく

- [インストールガイド](/ja/guide/installation) - システムへのPlexrのインストール
- [設定リファレンス](/ja/guide/configuration) - YAML設定について
- [サンプル](/ja/examples/) - 実際の使用例
- [APIドキュメント](/ja/api/) - 詳細な技術リファレンス