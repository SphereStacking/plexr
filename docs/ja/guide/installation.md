# インストール

Plexrは、お使いのプラットフォームや好みに応じて、いくつかの方法でインストールできます。

## システム要件

- **オペレーティングシステム**: macOS、Linux、またはWindows
- **アーキテクチャ**: amd64またはarm64
- **Go**: バージョン1.21以上（ソースからビルドする場合のみ）

## インストール方法

### Go Installを使用する（推奨）

Goがインストールされている場合、これが最も簡単な方法です：

```bash
# 最新バージョンをインストール
go install github.com/SphereStacking/plexr/cmd/plexr@latest

# または特定のバージョンをインストール（例：v0.1.0）
go install github.com/SphereStacking/plexr/cmd/plexr@v0.1.0
```

これにより、Plexrが`$GOPATH/bin`ディレクトリにインストールされます。

### ソースからビルドする

最新の開発バージョンを使用したい場合や、貢献したい場合：

```bash
# リポジトリをクローン
git clone https://github.com/SphereStacking/plexr.git
cd plexr

# 依存関係をインストール
make deps

# バイナリをビルド
make build

# PATHにインストール
make install
```

### パッケージマネージャー

#### Homebrew (macOS/Linux)

近日公開予定：
```bash
brew install plexr
```

#### Scoop (Windows)

近日公開予定：
```bash
scoop install plexr
```

### バイナリリリース

複数のプラットフォーム向けに事前ビルドされたバイナリが利用可能です。[リリースページ](https://github.com/SphereStacking/plexr/releases)からダウンロードしてください。

#### Linux (x86_64)
```bash
# 最新リリースをダウンロード
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Linux_x86_64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# または特定のバージョンをダウンロード（例：v0.1.0）
curl -sSL https://github.com/SphereStacking/plexr/releases/download/v0.1.0/plexr_Linux_x86_64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/
```

#### Linux (arm64)
```bash
# 最新リリースをダウンロード
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Linux_arm64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# または特定のバージョンをダウンロード（例：v0.1.0）
curl -sSL https://github.com/SphereStacking/plexr/releases/download/v0.1.0/plexr_Linux_arm64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/
```

#### macOS (Intel)
```bash
# 最新リリースをダウンロード
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Darwin_x86_64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# または特定のバージョンをダウンロード（例：v0.1.0）
curl -sSL https://github.com/SphereStacking/plexr/releases/download/v0.1.0/plexr_Darwin_x86_64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/
```

#### macOS (Apple Silicon)
```bash
# 最新リリースをダウンロード
curl -sSL https://github.com/SphereStacking/plexr/releases/latest/download/plexr_Darwin_arm64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/

# または特定のバージョンをダウンロード（例：v0.1.0）
curl -sSL https://github.com/SphereStacking/plexr/releases/download/v0.1.0/plexr_Darwin_arm64.tar.gz | tar xz
sudo mv plexr /usr/local/bin/
```

#### Windows
1. [リリースページ](https://github.com/SphereStacking/plexr/releases)から適切なファイルをダウンロード：
   - 最新版：`plexr_Windows_x86_64.zip`
   - v0.1.0：v0.1.0のリリースページからダウンロード
2. zipファイルを展開
3. ディレクトリをPATHに追加するか、`plexr.exe`をPATH内のディレクトリに移動

#### インストールの確認
```bash
plexr version
```

## インストールの確認

インストール後、Plexrが正しくインストールされていることを確認します：

```bash
plexr --version
```

次のような出力が表示されるはずです：
```
plexr version 0.1.0
```

## シェル補完

Plexrは、bash、zsh、fish、PowerShellのシェル補完をサポートしています。

### Bash

```bash
# ~/.bashrcに追加
echo 'source <(plexr completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```bash
# ~/.zshrcに追加
echo 'source <(plexr completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish

```bash
plexr completion fish | source
# 永続化するには：
plexr completion fish > ~/.config/fish/completions/plexr.fish
```

### PowerShell

```powershell
# PowerShellプロファイルに追加
plexr completion powershell | Out-String | Invoke-Expression
```

## 環境変数

Plexrは以下の環境変数を使用します：

- `PLEXR_STATE_FILE`: デフォルトの状態ファイルの場所を上書き
- `PLEXR_LOG_LEVEL`: ログレベルを設定（debug、info、warn、error）
- `PLEXR_NO_COLOR`: カラー出力を無効化

例：
```bash
export PLEXR_LOG_LEVEL=debug
export PLEXR_STATE_FILE=/tmp/plexr_state.json
```

## アップグレード

### Goを使用

```bash
# 最新バージョンにアップグレード
go install github.com/SphereStacking/plexr/cmd/plexr@latest

# または特定のバージョンにアップグレード
go install github.com/SphereStacking/plexr/cmd/plexr@v0.1.0
```

### ソースから

```bash
cd plexr
git pull
make clean build install
```

## アンインストール

### Goでインストールした場合

```bash
rm $(go env GOPATH)/bin/plexr
```

### 手動インストール

```bash
rm /usr/local/bin/plexr
```

### 状態ファイルのクリーンアップ

Plexrはプロジェクトディレクトリに状態ファイルを作成します：

```bash
# 状態ファイルを削除
find . -name ".plexr_state.json" -delete
```

## トラブルシューティング

### コマンドが見つからない

インストール後に「コマンドが見つかりません」というエラーが表示される場合：

1. バイナリがPATHにあるか確認：
   ```bash
   which plexr
   ```

2. Goインストールの場合、`$GOPATH/bin`がPATHにあることを確認：
   ```bash
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

### 権限が拒否されました

権限エラーが発生する場合：

```bash
chmod +x /path/to/plexr
```

### バージョンの競合

複数のバージョンがインストールされている場合：

```bash
# すべてのplexrインストールを検索
which -a plexr

# 特定のバージョンを使用
/usr/local/bin/plexr --version
```

## 次のステップ

- [はじめに](/guide/getting-started)を読む
- [設定](/guide/configuration)について学ぶ
- 実際の使用例の[サンプル](/examples/)を参照