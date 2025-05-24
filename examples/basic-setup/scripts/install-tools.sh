#!/bin/bash

# 必要なパッケージのインストール
echo "Installing required packages..."
sudo apt-get update
sudo apt-get install -y \
    git \
    curl \
    wget \
    build-essential \
    python3 \
    python3-pip \
    nodejs \
    npm

# Goのインストール
echo "Installing Go..."
if ! command -v go &> /dev/null; then
    wget https://go.dev/dl/go1.21.7.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.21.7.linux-amd64.tar.gz
    rm go1.21.7.linux-amd64.tar.gz
fi

# Dockerのインストール
echo "Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh
fi

echo "Tool installation completed!" 
