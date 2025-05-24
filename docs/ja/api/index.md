# APIリファレンス

Plexr APIリファレンスドキュメントへようこそ。このセクションでは、Plexrのすべての側面に関する詳細な技術情報を提供します。

## 利用可能なリファレンス

### [CLIコマンド](/api/cli-commands)

すべてのコマンドラインインターフェースコマンドの完全なリファレンス：
- `execute` - 実行プランを実行
- `validate` - プラン構文を検証
- `status` - 実行ステータスを確認
- `reset` - 実行状態をリセット
- その他...

### [設定スキーマ](/api/configuration-schema)

詳細なYAML設定スキーマドキュメント：
- プラン構造とフィールド
- ステップ定義
- エグゼキューター設定
- プラットフォーム固有の設定
- 検証ルール

### [エグゼキューターAPI](/api/executors)

組み込みおよびカスタムエグゼキューターに関する情報：
- シェルエグゼキューター
- SQLエグゼキューター（近日公開）
- カスタムエグゼキューターの作成
- エグゼキューターインターフェース

## クイックリンク

### コマンドライン

- [executeコマンド](/api/cli-commands#plexr-execute) - メイン実行コマンド
- [グローバルフラグ](/api/cli-commands#global-flags) - すべてのコマンドで利用可能なフラグ
- [終了コード](/api/cli-commands#exit-codes) - 戻り値の理解

### 設定

- [ルートフィールド](/api/configuration-schema#root-fields) - トップレベルの設定
- [ステップスキーマ](/api/configuration-schema#step) - ステップ設定の詳細
- [ファイル設定](/api/configuration-schema#fileconfig) - ファイル実行オプション

### 開発

- [エグゼキューターインターフェース](/api/executors#executor-interface) - エグゼキューターの実装
- [状態管理](/api/executors#state-management) - 状態の操作
- [エラー処理](/api/executors#error-handling) - ベストプラクティス

## 環境変数

Plexrは設定のためにいくつかの環境変数を使用します：

| 変数 | 説明 | デフォルト |
|----------|-------------|---------|
| `PLEXR_STATE_FILE` | 状態ファイルの場所 | `.plexr_state.json` |
| `PLEXR_LOG_LEVEL` | ログレベル | `info` |
| `PLEXR_NO_COLOR` | カラーを無効化 | `false` |
| `PLEXR_PLATFORM` | プラットフォームを上書き | 自動検出 |

## ファイル形式

### 状態ファイル形式

状態ファイル（`.plexr_state.json`）は実行の進捗を追跡します：

```json
{
  "version": "1.0",
  "plan_name": "開発環境セットアップ",
  "plan_version": "1.0.0",
  "started_at": "2023-12-15T10:00:00Z",
  "updated_at": "2023-12-15T10:30:00Z",
  "current_step": "configure_app",
  "completed_steps": [
    {
      "id": "install_tools",
      "completed_at": "2023-12-15T10:10:00Z"
    }
  ],
  "failed_steps": [],
  "installed_tools": {
    "node": "20.10.0",
    "docker": "24.0.7"
  }
}
```

### 設定ファイル形式

完全なYAML形式のドキュメントについては、[設定スキーマ](/api/configuration-schema)を参照してください。

## エラーコード

Plexrはすべての操作で一貫したエラーコードを使用します：

| コード | カテゴリ | 説明 |
|------|----------|-------------|
| 0 | 成功 | 操作が正常に完了 |
| 1-99 | 一般 | 一般的なエラー |
| 100-199 | 検証 | 設定検証エラー |
| 200-299 | 実行 | ランタイム実行エラー |
| 300-399 | 状態 | 状態管理エラー |
| 400-499 | プラットフォーム | プラットフォーム固有のエラー |

## バージョニング

Plexrはセマンティックバージョニングに従います：

- **メジャー:** CLIまたは設定形式の破壊的変更
- **マイナー:** 新機能、後方互換性あり
- **パッチ:** バグ修正と小さな改善

## サポート

追加のヘルプについて：
- [GitHubイシュー](https://github.com/SphereStacking/plexr/issues)
- [サンプル](/examples/)
- [トラブルシューティングガイド](/guide/troubleshooting)