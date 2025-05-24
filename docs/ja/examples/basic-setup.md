# 基本セットアップの例

この例では、Plexrを使用した基本的な開発環境のセットアップを示します。ステップの依存関係、プラットフォーム固有のスクリプト、状態管理などの基本的な概念を示しています。

## 概要

このプランは以下を行います：
1. 必須の開発ツールをインストール
2. Git設定を構成
3. プロジェクトディレクトリ構造を作成
4. 拡張機能を含むVS Codeをセットアップ

## プランファイル

`basic-setup.yml`を作成：

```yaml
name: "基本開発環境"
version: "1.0.0"
description: |
  以下を含む基本的な開発環境をセットアップします：
  - 必須の開発ツール
  - Git設定
  - プロジェクトディレクトリ
  - VS Codeセットアップ

executors:
  shell:
    type: shell
    config:
      timeout: 300

steps:
  - id: install_tools
    description: "必須の開発ツールをインストール"
    executor: shell
    files:
      - path: "scripts/install_tools.sh"
        platform: linux
      - path: "scripts/install_tools_mac.sh"
        platform: darwin
      - path: "scripts/install_tools.ps1"
        platform: windows

  - id: setup_git
    description: "Git設定を構成"
    executor: shell
    depends_on: [install_tools]
    check_command: "git config --get user.name"
    files:
      - path: "scripts/setup_git.sh"
        platform: linux
      - path: "scripts/setup_git.sh"
        platform: darwin
      - path: "scripts/setup_git.ps1"
        platform: windows

  - id: create_directories
    description: "プロジェクトディレクトリ構造を作成"
    executor: shell
    files:
      - path: "scripts/create_dirs.sh"
        platform: linux
      - path: "scripts/create_dirs.sh"
        platform: darwin
      - path: "scripts/create_dirs.ps1"
        platform: windows

  - id: setup_vscode
    description: "VS Codeを設定"
    executor: shell
    depends_on: [install_tools]
    skip_if: "! command -v code >/dev/null 2>&1"
    files:
      - path: "scripts/setup_vscode.sh"
        platform: linux
      - path: "scripts/setup_vscode.sh"
        platform: darwin
      - path: "scripts/setup_vscode.ps1"
        platform: windows
```

## スクリプトファイル

### Linux/macOS: `scripts/install_tools.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "🔧 必須の開発ツールをインストール中..."

# コマンドが存在するかチェックする関数
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# curlがなければインストール
if ! command_exists curl; then
    echo "curlをインストール中..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install curl
    else
        sudo apt-get update && sudo apt-get install -y curl
    fi
fi

# gitがなければインストール
if ! command_exists git; then
    echo "gitをインストール中..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install git
    else
        sudo apt-get install -y git
    fi
fi

# NodeSourceを使用してNode.jsをインストール
if ! command_exists node; then
    echo "Node.jsをインストール中..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install node
    else
        curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
        sudo apt-get install -y nodejs
    fi
fi

# VS Codeをインストール
if ! command_exists code; then
    echo "Visual Studio Codeをインストール中..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install --cask visual-studio-code
    else
        # Linuxインストール
        wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
        sudo install -o root -g root -m 644 packages.microsoft.gpg /etc/apt/trusted.gpg.d/
        sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
        sudo apt-get update
        sudo apt-get install -y code
    fi
fi

echo "✅ 開発ツールのインストールが成功しました！"
```

### Windows: `scripts/install_tools.ps1`

```powershell
# Windowsで開発ツールをインストール

Write-Host "🔧 必須の開発ツールをインストール中..." -ForegroundColor Green

# Chocolateyがインストールされているか確認
if (!(Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Chocolateyをインストール中..."
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
}

# Chocolateyを使用してツールをインストール
$tools = @('git', 'nodejs', 'vscode', 'curl')

foreach ($tool in $tools) {
    if (!(Get-Command $tool -ErrorAction SilentlyContinue)) {
        Write-Host "$tool をインストール中..."
        choco install $tool -y
    } else {
        Write-Host "$tool はすでにインストールされています" -ForegroundColor Yellow
    }
}

Write-Host "✅ 開発ツールのインストールが成功しました！" -ForegroundColor Green
```

### Gitセットアップ: `scripts/setup_git.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "🔧 Gitを設定中..."

# ユーザー情報を入力
read -p "Gitコミット用の名前を入力してください: " git_name
read -p "Gitコミット用のメールアドレスを入力してください: " git_email

# Gitを設定
git config --global user.name "$git_name"
git config --global user.email "$git_email"

# 便利なエイリアスを設定
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global alias.lg "log --oneline --graph --decorate"

# エディタを設定
if command -v code >/dev/null 2>&1; then
    git config --global core.editor "code --wait"
elif command -v vim >/dev/null 2>&1; then
    git config --global core.editor vim
fi

# デフォルトのブランチ名を設定
git config --global init.defaultBranch main

echo "✅ Gitの設定が成功しました！"
```

### ディレクトリ作成: `scripts/create_dirs.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "📁 プロジェクトディレクトリ構造を作成中..."

# ベースディレクトリ
WORKSPACE="$HOME/workspace"
PROJECTS="$WORKSPACE/projects"
SCRIPTS="$WORKSPACE/scripts"
CONFIGS="$WORKSPACE/configs"

# ディレクトリ構造を作成
mkdir -p "$PROJECTS"/{personal,work,learning}
mkdir -p "$SCRIPTS"
mkdir -p "$CONFIGS"

# ワークスペースにREADMEを作成
cat > "$WORKSPACE/README.md" << EOF
# ワークスペース

これはあなたの開発ワークスペースで、以下のように整理されています：

- **projects/**: コーディングプロジェクト
  - personal/: 個人プロジェクト
  - work/: 仕事関連のプロジェクト
  - learning/: 学習とチュートリアルプロジェクト
- **scripts/**: 便利なスクリプトとツール
- **configs/**: 設定ファイルとドットファイル

Plexrによって$(date)に作成されました
EOF

echo "✅ ディレクトリ構造が$WORKSPACEに作成されました"
```

### VS Codeセットアップ: `scripts/setup_vscode.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "🎨 VS Codeをセットアップ中..."

# 便利な拡張機能をインストール
extensions=(
    "ms-vscode.vscode-typescript-next"
    "dbaeumer.vscode-eslint"
    "esbenp.prettier-vscode"
    "eamodio.gitlens"
    "PKief.material-icon-theme"
    "GitHub.copilot"
    "ms-python.python"
    "golang.go"
    "rust-lang.rust-analyzer"
)

for ext in "${extensions[@]}"; do
    echo "拡張機能をインストール中: $ext"
    code --install-extension "$ext" || true
done

# VS Code設定を作成
VSCODE_CONFIG="$HOME/.config/Code/User"
mkdir -p "$VSCODE_CONFIG"

cat > "$VSCODE_CONFIG/settings.json" << 'EOF'
{
    "editor.fontSize": 14,
    "editor.tabSize": 2,
    "editor.formatOnSave": true,
    "editor.minimap.enabled": false,
    "workbench.iconTheme": "material-icon-theme",
    "terminal.integrated.fontSize": 14,
    "files.autoSave": "afterDelay",
    "files.autoSaveDelay": 1000,
    "git.autofetch": true,
    "git.confirmSync": false
}
EOF

echo "✅ VS Codeの設定が成功しました！"
```

## 例を実行する

1. すべてのファイルを正しい構造で保存：
   ```
   basic-setup/
   ├── basic-setup.yml
   └── scripts/
       ├── install_tools.sh
       ├── install_tools_mac.sh
       ├── install_tools.ps1
       ├── setup_git.sh
       ├── setup_git.ps1
       ├── create_dirs.sh
       ├── create_dirs.ps1
       ├── setup_vscode.sh
       └── setup_vscode.ps1
   ```

2. スクリプトを実行可能にする（Linux/macOS）：
   ```bash
   chmod +x scripts/*.sh
   ```

3. プランを実行：
   ```bash
   plexr execute basic-setup.yml
   ```

## 何が起こるか

1. **初回実行**: すべてのステップが順番に実行されます
2. **後続の実行**: 
   - `check_command`を持つステップは、すでに完了している場合はスキップされます
   - `skip_if`を持つステップは条件を評価します
   - 必要なステップのみが実行されます

## ステータスの確認

```bash
# 実行ステータスを確認
plexr status basic-setup.yml

# 実行せずに何が実行されるかを確認
plexr execute basic-setup.yml --dry-run
```

## カスタマイズのアイデア

1. **さらにツールを追加**: お気に入りのツールでインストールスクリプトを拡張
2. **個人のGit設定**: 好みのGitエイリアスと設定を追加
3. **異なるディレクトリ**: ワークフローに合わせてディレクトリ構造を変更
4. **さらにVS Code拡張機能**: 技術スタック用の拡張機能を追加

## トラブルシューティング

### スクリプトの権限
```bash
chmod +x scripts/*.sh
```

### 失敗後の再開
```bash
# ステップが失敗した場合、問題を修正して再度実行
plexr execute basic-setup.yml
```

### リセットして最初からやり直す
```bash
plexr reset basic-setup.yml --force
plexr execute basic-setup.yml
```

## 次のステップ

- [高度なパターン](/examples/advanced-patterns)の例を試す
- [状態管理](/guide/state-management)について学ぶ
- 独自のカスタムプランを作成