#!/bin/bash

# VSCodeの設定
echo "Setting up VSCode..."

# VSCodeのインストール（まだの場合）
if ! command -v code &> /dev/null; then
    echo "Installing VSCode..."
    curl -fsSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
    sudo install -o root -g root -m 644 packages.microsoft.gpg /usr/share/keyrings/
    sudo sh -c 'echo "deb [arch=amd64 signed-by=/usr/share/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/vscode stable main" > /etc/apt/sources.list.d/vscode.list'
    sudo apt-get update
    sudo apt-get install -y code
    rm packages.microsoft.gpg
fi

# 推奨拡張機能のインストール
echo "Installing recommended extensions..."
code --install-extension ms-vscode.go
code --install-extension ms-python.python
code --install-extension dbaeumer.vscode-eslint
code --install-extension esbenp.prettier-vscode
code --install-extension ms-azuretools.vscode-docker
code --install-extension eamodio.gitlens
code --install-extension ms-vscode-remote.remote-containers

# 設定ファイルのコピー
echo "Copying VSCode settings..."
VSCODE_CONFIG_DIR="$HOME/.config/Code/User"
mkdir -p "$VSCODE_CONFIG_DIR"

# 基本的な設定
cat > "$VSCODE_CONFIG_DIR/settings.json" << EOF
{
    "editor.formatOnSave": true,
    "editor.rulers": [80, 100],
    "editor.tabSize": 4,
    "files.trimTrailingWhitespace": true,
    "files.insertFinalNewline": true,
    "go.formatTool": "goimports",
    "python.formatting.provider": "black",
    "python.linting.enabled": true,
    "python.linting.pylintEnabled": true
}
EOF

echo "VSCode setup completed!" 
