# Installation

Plexr can be installed in several ways depending on your platform and preferences.

## System Requirements

- **Operating System**: macOS, Linux, or Windows
- **Architecture**: amd64 or arm64
- **Go**: Version 1.21+ (only for building from source)

## Installation Methods

### Using Go Install (Recommended)

If you have Go installed, this is the simplest method:

```bash
go install github.com/SphereStacking/plexr@latest
```

This will install the latest version of Plexr to your `$GOPATH/bin` directory.

### Building from Source

For the latest development version or to contribute:

```bash
# Clone the repository
git clone https://github.com/SphereStacking/plexr.git
cd plexr

# Install dependencies
make deps

# Build the binary
make build

# Install to your PATH
make install
```

### Package Managers

#### Homebrew (macOS/Linux)

Coming soon:
```bash
brew install plexr
```

#### Scoop (Windows)

Coming soon:
```bash
scoop install plexr
```

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/SphereStacking/plexr/releases):

1. Download the appropriate archive for your platform
2. Extract the binary
3. Move it to a directory in your PATH

Example for Linux/macOS:
```bash
# Download (replace VERSION and PLATFORM)
curl -L https://github.com/SphereStacking/plexr/releases/download/vVERSION/plexr_PLATFORM.tar.gz -o plexr.tar.gz

# Extract
tar -xzf plexr.tar.gz

# Move to PATH
sudo mv plexr /usr/local/bin/

# Verify installation
plexr --version
```

## Verifying Installation

After installation, verify that Plexr is correctly installed:

```bash
plexr --version
```

You should see output like:
```
plexr version 1.0.0
```

## Shell Completion

Plexr supports shell completion for bash, zsh, fish, and PowerShell.

### Bash

```bash
# Add to ~/.bashrc
echo 'source <(plexr completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```bash
# Add to ~/.zshrc
echo 'source <(plexr completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish

```bash
plexr completion fish | source
# To persist:
plexr completion fish > ~/.config/fish/completions/plexr.fish
```

### PowerShell

```powershell
# Add to your PowerShell profile
plexr completion powershell | Out-String | Invoke-Expression
```

## Environment Variables

Plexr respects the following environment variables:

- `PLEXR_STATE_FILE`: Override the default state file location
- `PLEXR_LOG_LEVEL`: Set logging level (debug, info, warn, error)
- `PLEXR_NO_COLOR`: Disable colored output

Example:
```bash
export PLEXR_LOG_LEVEL=debug
export PLEXR_STATE_FILE=/tmp/plexr_state.json
```

## Upgrading

### Using Go

```bash
go install github.com/SphereStacking/plexr@latest
```

### From Source

```bash
cd plexr
git pull
make clean build install
```

## Uninstalling

### Installed with Go

```bash
rm $(go env GOPATH)/bin/plexr
```

### Manual Installation

```bash
rm /usr/local/bin/plexr
```

### Cleanup State Files

Plexr creates state files in your project directories:

```bash
# Remove state files
find . -name ".plexr_state.json" -delete
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Check if the binary is in your PATH:
   ```bash
   which plexr
   ```

2. For Go installations, ensure `$GOPATH/bin` is in your PATH:
   ```bash
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

### Permission Denied

If you get permission errors:

```bash
chmod +x /path/to/plexr
```

### Version Conflicts

If you have multiple versions installed:

```bash
# Find all plexr installations
which -a plexr

# Use specific version
/usr/local/bin/plexr --version
```

## Next Steps

- Read the [Getting Started Guide](/guide/getting-started)
- Learn about [Configuration](/guide/configuration)
- See [Examples](/examples/) of real-world usage