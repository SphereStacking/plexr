#!/bin/bash

# 開発用ディレクトリの作成
echo "Creating development directory structure..."

# ホームディレクトリに開発用フォルダを作成
DEV_DIR="$HOME/development"
mkdir -p "$DEV_DIR"

# プロジェクト用のサブディレクトリを作成
mkdir -p "$DEV_DIR/projects"
mkdir -p "$DEV_DIR/tools"
mkdir -p "$DEV_DIR/workspaces"

# 各プロジェクトタイプ用のディレクトリ
mkdir -p "$DEV_DIR/projects/go"
mkdir -p "$DEV_DIR/projects/python"
mkdir -p "$DEV_DIR/projects/node"
mkdir -p "$DEV_DIR/projects/docker"

# 設定ファイル用のディレクトリ
mkdir -p "$DEV_DIR/config"
mkdir -p "$DEV_DIR/config/vscode"
mkdir -p "$DEV_DIR/config/shell"

echo "Directory structure created at $DEV_DIR" 
