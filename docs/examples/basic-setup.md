# Basic Setup Example

This example demonstrates a basic development environment setup using Plexr. It shows fundamental concepts like step dependencies, platform-specific scripts, and state management.

## Overview

This plan will:
1. Install essential development tools
2. Configure Git settings
3. Create project directory structure
4. Set up VS Code with extensions

## The Plan File

Create `basic-setup.yml`:

```yaml
name: "Basic Development Environment"
version: "1.0.0"
description: |
  Sets up a basic development environment with:
  - Essential development tools
  - Git configuration
  - Project directories
  - VS Code setup

executors:
  shell:
    type: shell
    config:
      timeout: 300

steps:
  - id: install_tools
    description: "Install essential development tools"
    executor: shell
    files:
      - path: "scripts/install_tools.sh"
        platform: linux
      - path: "scripts/install_tools_mac.sh"
        platform: darwin
      - path: "scripts/install_tools.ps1"
        platform: windows

  - id: setup_git
    description: "Configure Git settings"
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
    description: "Create project directory structure"
    executor: shell
    files:
      - path: "scripts/create_dirs.sh"
        platform: linux
      - path: "scripts/create_dirs.sh"
        platform: darwin
      - path: "scripts/create_dirs.ps1"
        platform: windows

  - id: setup_vscode
    description: "Configure VS Code"
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

## Script Files

### Linux/macOS: `scripts/install_tools.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "🔧 Installing essential development tools..."

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install curl if not present
if ! command_exists curl; then
    echo "Installing curl..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install curl
    else
        sudo apt-get update && sudo apt-get install -y curl
    fi
fi

# Install git if not present
if ! command_exists git; then
    echo "Installing git..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install git
    else
        sudo apt-get install -y git
    fi
fi

# Install Node.js using NodeSource
if ! command_exists node; then
    echo "Installing Node.js..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install node
    else
        curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
        sudo apt-get install -y nodejs
    fi
fi

# Install VS Code
if ! command_exists code; then
    echo "Installing Visual Studio Code..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install --cask visual-studio-code
    else
        # Linux installation
        wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
        sudo install -o root -g root -m 644 packages.microsoft.gpg /etc/apt/trusted.gpg.d/
        sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
        sudo apt-get update
        sudo apt-get install -y code
    fi
fi

echo "✅ Development tools installed successfully!"
```

### Windows: `scripts/install_tools.ps1`

```powershell
# Install development tools on Windows

Write-Host "🔧 Installing essential development tools..." -ForegroundColor Green

# Check if Chocolatey is installed
if (!(Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Chocolatey..."
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
}

# Install tools using Chocolatey
$tools = @('git', 'nodejs', 'vscode', 'curl')

foreach ($tool in $tools) {
    if (!(Get-Command $tool -ErrorAction SilentlyContinue)) {
        Write-Host "Installing $tool..."
        choco install $tool -y
    } else {
        Write-Host "$tool is already installed" -ForegroundColor Yellow
    }
}

Write-Host "✅ Development tools installed successfully!" -ForegroundColor Green
```

### Git Setup: `scripts/setup_git.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "🔧 Configuring Git..."

# Prompt for user information
read -p "Enter your name for Git commits: " git_name
read -p "Enter your email for Git commits: " git_email

# Configure Git
git config --global user.name "$git_name"
git config --global user.email "$git_email"

# Set up useful aliases
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global alias.lg "log --oneline --graph --decorate"

# Configure editor
if command -v code >/dev/null 2>&1; then
    git config --global core.editor "code --wait"
elif command -v vim >/dev/null 2>&1; then
    git config --global core.editor vim
fi

# Set up default branch name
git config --global init.defaultBranch main

echo "✅ Git configured successfully!"
```

### Create Directories: `scripts/create_dirs.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "📁 Creating project directory structure..."

# Base directories
WORKSPACE="$HOME/workspace"
PROJECTS="$WORKSPACE/projects"
SCRIPTS="$WORKSPACE/scripts"
CONFIGS="$WORKSPACE/configs"

# Create directory structure
mkdir -p "$PROJECTS"/{personal,work,learning}
mkdir -p "$SCRIPTS"
mkdir -p "$CONFIGS"

# Create a README in workspace
cat > "$WORKSPACE/README.md" << EOF
# Workspace

This is your development workspace organized as follows:

- **projects/**: Your coding projects
  - personal/: Personal projects
  - work/: Work-related projects
  - learning/: Learning and tutorial projects
- **scripts/**: Useful scripts and tools
- **configs/**: Configuration files and dotfiles

Created by Plexr on $(date)
EOF

echo "✅ Directory structure created at $WORKSPACE"
```

### VS Code Setup: `scripts/setup_vscode.sh`

```bash
#!/bin/bash
set -euo pipefail

echo "🎨 Setting up VS Code..."

# Install useful extensions
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
    echo "Installing extension: $ext"
    code --install-extension "$ext" || true
done

# Create VS Code settings
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

echo "✅ VS Code configured successfully!"
```

## Running the Example

1. Save all files in the correct structure:
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

2. Make scripts executable (Linux/macOS):
   ```bash
   chmod +x scripts/*.sh
   ```

3. Run the plan:
   ```bash
   plexr execute basic-setup.yml
   ```

## What Happens

1. **First Run**: All steps execute in order
2. **Subsequent Runs**: 
   - Steps with `check_command` skip if already done
   - Steps with `skip_if` evaluate their condition
   - Only necessary steps run

## Checking Status

```bash
# Check execution status
plexr status basic-setup.yml

# See what would run without executing
plexr execute basic-setup.yml --dry-run
```

## Customization Ideas

1. **Add More Tools**: Extend the install scripts with your favorite tools
2. **Personal Git Config**: Add your preferred Git aliases and settings
3. **Different Directories**: Modify the directory structure to match your workflow
4. **More VS Code Extensions**: Add extensions for your tech stack

## Troubleshooting

### Script Permissions
```bash
chmod +x scripts/*.sh
```

### Resume After Failure
```bash
# If a step fails, fix the issue and run again
plexr execute basic-setup.yml
```

### Reset and Start Over
```bash
plexr reset basic-setup.yml --force
plexr execute basic-setup.yml
```

## Next Steps

- Try the [Advanced Patterns](/examples/advanced-patterns) example
- Learn about [State Management](/guide/state-management)
- Create your own custom plan