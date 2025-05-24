# トラブルシューティング

このガイドは、Plexrを使用する際の一般的な問題の診断と修正に役立ちます。

## よくある問題

### インストールの問題

#### Goモジュールエラー

**問題**: `go: module github.com/plexr/plexr: git ls-remote -q origin` が失敗する

**解決策**:
```bash
# モジュールキャッシュをクリア
go clean -modcache

# 依存関係を再ダウンロード
go mod download
```

#### インストール中の権限拒否

**問題**: `go install`実行時に`permission denied`

**解決策**:
```bash
# ユーザーのbinディレクトリにインストール
go install github.com/plexr/plexr/cmd/plexr@latest

# またはsudoを使用（推奨されません）
sudo go install github.com/plexr/plexr/cmd/plexr@latest
```

### 設定の問題

#### プランファイルが見つからない

**問題**: `Error: plan file not found: plan.yml`

**解決策**:
1. 正しいディレクトリにいることを確認
2. ファイル名のスペルをチェック（plan.yml vs plan.yaml）
3. `-f`フラグを使用してカスタムパスを指定：
   ```bash
   plexr execute -f path/to/myplan.yml
   ```

#### 無効なYAML構文

**問題**: `Error: failed to parse plan: yaml: line X`

**解決策**:
1. オンラインバリデーターを使用してYAML構文を検証
2. 一般的な問題をチェック：
   - 不正なインデント（タブではなくスペースを使用）
   - キーの後のコロンの欠落
   - 引用符で囲まれていない特殊文字

正しい構文の例：
```yaml
version: "1.0"
name: "マイプラン"
steps:
  - name: "ステップ1"  # 引用符と適切なインデントに注意
    command: "echo hello"
```

### 実行の問題

#### "コマンドが見つかりません"でステップが失敗する

**問題**: `bash: command: command not found`

**解決策**:
1. コマンドがインストールされ、PATHにあることを確認
2. カスタムスクリプトには絶対パスを使用：
   ```yaml
   steps:
     - name: "スクリプトを実行"
       command: "./scripts/my-script.sh"  # 相対パス
       # または
       command: "/home/user/project/scripts/my-script.sh"  # 絶対パス
   ```

#### タイムアウトエラー

**問題**: `Error: step timed out after 60s`

**解決策**:
```yaml
steps:
  - name: "長時間実行タスク"
    command: "./long-task.sh"
    timeout: 300  # タイムアウトを5分に増加
```

#### 作業ディレクトリの問題

**問題**: `Error: no such file or directory`

**解決策**:
```yaml
steps:
  - name: "サブディレクトリでビルド"
    command: "make build"
    workdir: "./src"  # 作業ディレクトリを指定
```

### 状態管理の問題

#### 変数が見つからない

**問題**: `Error: variable 'myvar' not found in state`

**解決策**:
1. 使用前に変数が設定されていることを確認：
   ```yaml
   steps:
     - name: "変数を設定"
       command: "echo 'value'"
       outputs:
         - name: myvar
           from: stdout
     
     - name: "変数を使用"
       command: "echo `{{.myvar}}`"  # これで存在します
   ```

#### 状態ファイルの破損

**問題**: `Error: failed to load state: invalid character`

**解決策**:
```bash
# 状態をリセット
plexr reset

# または手動で状態ファイルを削除
rm .plexr/state.json
```

### 依存関係の問題

#### 循環依存

**問題**: `Error: circular dependency detected`

**解決策**:
ステップの依存関係を確認し、サイクルを削除：
```yaml
# 悪い例 - 循環依存
steps:
  - name: "A"
    depends_on: ["B"]
  - name: "B"
    depends_on: ["A"]

# 良い例 - サイクルなし
steps:
  - name: "A"
  - name: "B"
    depends_on: ["A"]
```

#### ステップが予期せずスキップされる

**問題**: 依存関係が満たされているのにステップが実行されない

**解決策**:
条件を確認：
```yaml
steps:
  - name: "条件付きステップ"
    command: "deploy.sh"
    condition: "`{{.environment}}` == 'production'"
    # environmentが"production"に設定されていない限り実行されません
```

## デバッグのヒント

### 詳細なロギングを有効化

```bash
# 詳細な実行情報を表示
plexr execute -v

# さらに詳細を表示
plexr execute -vv
```

### 実行計画を確認

```bash
# 実行せずに検証
plexr validate

# 実行順序を表示
plexr status --show-order
```

### 状態を検査

```bash
# 現在の状態を表示
plexr status

# 状態ファイルを直接表示
cat .plexr/state.json | jq
```

### 個別のステップをテスト

```yaml
# ステップにテストモードを追加
steps:
  - name: "デプロイ"
    command: |
      if [ "$TEST_MODE" = "true" ]; then
        echo "本番環境にデプロイします"
      else
        ./deploy.sh
      fi
    env:
      TEST_MODE: "`{{.test_mode | default \"false\"}}`"
```

## 環境固有の問題

### Dockerコンテナの問題

**問題**: Dockerで実行するとステップが失敗する

**解決策**:
1. 必要なボリュームをマウント：
   ```dockerfile
   docker run -v $(pwd):/workspace plexr execute
   ```

2. 作業ディレクトリを設定：
   ```yaml
   steps:
     - name: "コンテナでビルド"
       command: "make build"
       workdir: "/workspace"
   ```

### CI/CDパイプラインの失敗

**問題**: Plexrはローカルで動作するがCIで失敗する

**解決策**:
1. 環境変数が設定されていることを確認
2. すべての依存関係がインストールされていることを確認
3. 明示的なパスとバージョンを使用：
   ```yaml
   steps:
     - name: "CIビルド"
       command: "/usr/local/go/bin/go build"
       env:
         GOPATH: "/home/runner/go"
   ```

## パフォーマンスの問題

### 実行が遅い

**解決策**:
1. 独立したステップを並列で実行：
   ```yaml
   steps:
     - name: "並列タスク"
       parallel:
         - name: "テスト1"
           command: "test1.sh"
         - name: "テスト2"
           command: "test2.sh"
   ```

2. 依存関係をキャッシュ：
   ```yaml
   steps:
     - name: "依存関係をインストール"
       command: "npm ci"
       condition: "`{{.deps_cached}}` != 'true'"
   ```

### 高メモリ使用量

**解決策**:
大きなコマンドの出力キャプチャを制限：
```yaml
steps:
  - name: "大きな出力"
    command: "find / -name '*.log'"
    capture_output: false  # 状態に保存しない
```

## ヘルプを得る

### 診断情報を収集

```bash
# システム情報
plexr version
go version
uname -a

# プラン検証
plexr validate -v

# 実行トレース
plexr execute -vv 2>&1 | tee plexr-debug.log
```

### 問題を報告

問題を報告する際は、以下を含めてください：
1. Plexrバージョン（`plexr version`）
2. プランファイル（機密データを削除したもの）
3. エラーメッセージとログ
4. 再現手順

問題を報告する場所: https://github.com/plexr/plexr/issues

## FAQ

**Q: プランファイルで環境変数を使用できますか？**
A: はい、環境変数には`${VAR_NAME}`を、状態変数には<code>&#123;&#123;.var_name&#125;&#125;</code>を使用します。

**Q: シークレットをどのように扱いますか？**
A: 環境変数を使用し、プランファイルにシークレットをコミットしないでください。

**Q: Plexrをバックグラウンドで実行できますか？**
A: はい: `nohup plexr execute > plexr.log 2>&1 &`

**Q: Plexrをアップデートするにはどうすればよいですか？**
A: `go install github.com/plexr/plexr/cmd/plexr@latest`を実行