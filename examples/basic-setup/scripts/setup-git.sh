#!/bin/bash

# Gitのグローバル設定
echo "Setting up Git configuration..."

# ユーザー名とメールアドレスの設定
read -p "Enter your Git username: " git_username
read -p "Enter your Git email: " git_email

git config --global user.name "$git_username"
git config --global user.email "$git_email"

# デフォルトブランチ名の設定
git config --global init.defaultBranch main

# 便利なエイリアスの設定
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit

# カラーの設定
git config --global color.ui true
git config --global color.status auto
git config --global color.branch auto

echo "Git configuration completed!" 
