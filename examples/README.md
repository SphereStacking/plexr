# Plexr Examples

All examples in this directory are designed to be **side-effect free** and safe to run for demonstration purposes.

## Available Examples

### 1. echo-demo
Simple demonstration of plexr with basic echo commands.
- **Side effects**: None
- **Purpose**: Show basic plexr functionality

### 2. env-variables
Demonstrates environment variable handling.
- **Side effects**: None
- **Purpose**: Show how to read and check environment variables

### 3. file-operations
Demonstrates file operations using temporary files.
- **Side effects**: Creates temporary files (automatically cleaned up)
- **Purpose**: Show file creation and manipulation

### 4. work-directory
Demonstrates the work_directory feature.
- **Side effects**: Creates temporary file in /tmp (automatically cleaned up)
- **Purpose**: Show how to execute scripts in different directories

### 5. global-workdir
Demonstrates global work_directory settings.
- **Side effects**: Creates temporary directory in /tmp (automatically cleaned up)
- **Purpose**: Show global vs step-level work_directory configuration

### 6. basic-setup
Simulates a development environment setup without making actual changes.
- **Side effects**: None (simulation only)
- **Purpose**: Demonstrate what a real setup would do without executing it

## Running Examples

From the plexr root directory:
```bash
# Run an example
./build/plexr execute examples/echo-demo/plan.yml -a

# Or from the example directory
cd examples/echo-demo
../../build/plexr execute plan.yml
```

## Safety Features

- No system modifications
- No package installations
- No configuration changes
- Temporary files are created only in /tmp and cleaned up
- All operations are read-only or simulated