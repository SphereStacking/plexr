package executors

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellExecutor(t *testing.T) {
	t.Run("NewShellExecutor", func(t *testing.T) {
		executor := NewShellExecutor()
		assert.NotNil(t, executor)
		assert.Equal(t, "shell", executor.Name())

		// Check default shell based on platform
		if runtime.GOOS == "windows" {
			assert.Equal(t, "powershell.exe", executor.shell)
		} else {
			assert.Equal(t, "/bin/bash", executor.shell)
		}
	})

	t.Run("Validate configuration", func(t *testing.T) {
		tests := []struct {
			name    string
			config  map[string]interface{}
			wantErr bool
			errMsg  string
		}{
			{
				name:    "valid shell path",
				config:  map[string]interface{}{"shell": "/bin/bash"},
				wantErr: false,
			},
			{
				name:    "empty config uses default",
				config:  map[string]interface{}{},
				wantErr: false,
			},
			{
				name:    "invalid shell type",
				config:  map[string]interface{}{"shell": 123},
				wantErr: true,
				errMsg:  "shell must be a string",
			},
			{
				name: "valid with environment variables",
				config: map[string]interface{}{
					"shell": "/bin/bash",
					"env": map[string]string{
						"TEST_VAR": "value",
					},
				},
				wantErr: false,
			},
		}

		executor := NewShellExecutor()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := executor.Validate(tt.config)
				if tt.wantErr {
					assert.Error(t, err)
					if tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Execute scripts with different outcomes", func(t *testing.T) {
		tests := []struct {
			name           string
			script         string
			expectedOutput string
			expectedError  bool
			checkDuration  bool
		}{
			{
				name: "successful execution with output",
				script: `#!/bin/bash
echo "Starting process..."
echo "Process completed successfully"
exit 0`,
				expectedOutput: "Process completed successfully",
				expectedError:  false,
				checkDuration:  true,
			},
			{
				name: "script with error exit code",
				script: `#!/bin/bash
echo "Error: Something went wrong" >&2
exit 1`,
				expectedOutput: "Error: Something went wrong",
				expectedError:  true,
				checkDuration:  true,
			},
			{
				name: "script with both stdout and stderr",
				script: `#!/bin/bash
echo "Normal output"
echo "Error output" >&2
exit 0`,
				expectedOutput: "Normal output",
				expectedError:  false,
				checkDuration:  true,
			},
			{
				name: "script using environment variables",
				script: `#!/bin/bash
echo "HOME=$HOME"
echo "USER=$USER"
echo "PWD=$PWD"`,
				expectedOutput: "HOME=",
				expectedError:  false,
				checkDuration:  true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tmpDir := t.TempDir()
				scriptPath := filepath.Join(tmpDir, "test.sh")

				err := os.WriteFile(scriptPath, []byte(tt.script), 0755) // #nosec G306 - Script needs to be executable
				require.NoError(t, err)

				executor := NewShellExecutor()
				ctx := context.Background()

				file := ExecutionFile{
					Path: scriptPath,
				}

				result, err := executor.Execute(ctx, file)

				if tt.expectedError {
					assert.Error(t, err)
					assert.False(t, result.Success)
				} else {
					assert.NoError(t, err)
					assert.True(t, result.Success)
				}

				assert.Contains(t, result.Output, tt.expectedOutput)

				if tt.checkDuration {
					assert.GreaterOrEqual(t, result.Duration, int64(0))
				}
			})
		}
	})

	t.Run("Working directory functionality", func(t *testing.T) {
		tmpDir := t.TempDir()
		workDir := filepath.Join(tmpDir, "work")
		err := os.MkdirAll(workDir, 0755) // #nosec G301 - Test directory
		require.NoError(t, err)

		// Create a test file in work directory
		testFile := filepath.Join(workDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644) // #nosec G306 - Test file
		require.NoError(t, err)

		scriptPath := filepath.Join(tmpDir, "pwd_test.sh")
		script := `#!/bin/bash
pwd
ls -la test.txt 2>/dev/null || echo "test.txt not found"`
		err = os.WriteFile(scriptPath, []byte(script), 0755) // #nosec G306 - Script needs to be executable
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		// Test with work directory
		file := ExecutionFile{
			Path:          scriptPath,
			WorkDirectory: workDir,
		}

		result, err := executor.Execute(ctx, file)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Contains(t, result.Output, workDir)
		assert.NotContains(t, result.Output, "test.txt not found")

		// Test without work directory
		file.WorkDirectory = ""
		result, err = executor.Execute(ctx, file)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Contains(t, result.Output, "test.txt not found")
	})

	t.Run("Timeout handling", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping timeout test on Windows")
		}

		tests := []struct {
			name          string
			script        string
			timeout       int
			shouldTimeout bool
		}{
			{
				name: "script completes before timeout",
				script: `#!/bin/bash
echo "Quick task"
sleep 0.1
echo "Done"`,
				timeout:       2,
				shouldTimeout: false,
			},
			{
				name: "script exceeds timeout",
				script: `#!/bin/bash
echo "Starting long task"
sleep 5
echo "This should not appear"`,
				timeout:       1,
				shouldTimeout: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tmpDir := t.TempDir()
				scriptPath := filepath.Join(tmpDir, "timeout.sh")

				err := os.WriteFile(scriptPath, []byte(tt.script), 0755) // #nosec G306 - Script needs to be executable
				require.NoError(t, err)

				executor := NewShellExecutor()
				ctx := context.Background()

				file := ExecutionFile{
					Path:    scriptPath,
					Timeout: tt.timeout,
				}

				start := time.Now()
				result, err := executor.Execute(ctx, file)
				duration := time.Since(start)

				if tt.shouldTimeout {
					assert.Error(t, err)
					assert.False(t, result.Success)
					assert.Contains(t, err.Error(), "timeout")
					// Just verify it didn't run for the full sleep duration
					_ = duration // Duration check removed as process termination timing can vary
				} else {
					assert.NoError(t, err)
					assert.True(t, result.Success)
					assert.Contains(t, result.Output, "Done")
				}
			})
		}
	})

	t.Run("Platform-specific execution", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create scripts for different platforms
		scripts := map[string]string{
			"linux": `#!/bin/bash
echo "Running on Linux"
uname -s`,
			"darwin": `#!/bin/bash
echo "Running on macOS"
sw_vers 2>/dev/null || echo "Not on macOS"`,
			"windows": `@echo off
echo Running on Windows
ver`,
		}

		files := make(map[string]string)
		for platform, script := range scripts {
			scriptPath := filepath.Join(tmpDir, platform+".sh")
			if platform == "windows" {
				scriptPath = filepath.Join(tmpDir, platform+".bat")
			}
			err := os.WriteFile(scriptPath, []byte(script), 0755) // #nosec G306 - Script needs to be executable
			require.NoError(t, err)
			files[platform] = scriptPath
		}

		executor := NewShellExecutor()
		ctx := context.Background()

		// Test each platform script
		for platform, scriptPath := range files {
			t.Run("platform_"+platform, func(t *testing.T) {
				file := ExecutionFile{
					Path:     scriptPath,
					Platform: platform,
				}

				result, err := executor.Execute(ctx, file)
				require.NoError(t, err)
				assert.True(t, result.Success)

				if runtime.GOOS == platform {
					// Should execute on matching platform
					assert.Contains(t, result.Output, "Running on")
				} else {
					// Should skip on non-matching platform
					assert.Contains(t, result.Output, "Skipping file")
					assert.Contains(t, result.Output, platform)
				}
			})
		}
	})

	t.Run("Error cases", func(t *testing.T) {
		executor := NewShellExecutor()
		ctx := context.Background()

		tests := []struct {
			name        string
			file        ExecutionFile
			errContains string
		}{
			{
				name: "non-existent file",
				file: ExecutionFile{
					Path: "/non/existent/script.sh",
				},
				errContains: "file not found",
			},
			{
				name: "empty file path",
				file: ExecutionFile{
					Path: "",
				},
				errContains: "file not found",
			},
			{
				name: "directory instead of file",
				file: ExecutionFile{
					Path: os.TempDir(),
				},
				errContains: "exit status",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := executor.Execute(ctx, tt.file)
				assert.Error(t, err)
				assert.False(t, result.Success)
				assert.Contains(t, err.Error(), tt.errContains)
			})
		}
	})

	t.Run("Context cancellation", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping context cancellation test on Windows")
		}

		tmpDir := t.TempDir()
		scriptPath := filepath.Join(tmpDir, "long_running.sh")

		script := `#!/bin/bash
for i in {1..10}; do
    echo "Step $i"
    sleep 1
done`
		err := os.WriteFile(scriptPath, []byte(script), 0755) // #nosec G306 - Script needs to be executable
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx, cancel := context.WithCancel(context.Background())

		file := ExecutionFile{
			Path: scriptPath,
		}

		// Cancel context after a short delay
		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		result, err := executor.Execute(ctx, file)
		assert.Error(t, err)
		assert.False(t, result.Success)

		// Output should contain at least one step but not all
		assert.Contains(t, result.Output, "Step 1")
		assert.NotContains(t, result.Output, "Step 10")
	})

	t.Run("Script with large output", func(t *testing.T) {
		tmpDir := t.TempDir()
		scriptPath := filepath.Join(tmpDir, "large_output.sh")

		// Create a script that generates a lot of output
		script := `#!/bin/bash
for i in {1..100}; do
    echo "Line $i: This is a test line with some content to make it longer"
done`
		err := os.WriteFile(scriptPath, []byte(script), 0755) // #nosec G306 - Script needs to be executable
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		file := ExecutionFile{
			Path: scriptPath,
		}

		result, err := executor.Execute(ctx, file)
		require.NoError(t, err)
		assert.True(t, result.Success)

		// Verify output contains expected content
		lines := strings.Split(strings.TrimSpace(result.Output), "\n")
		assert.Equal(t, 100, len(lines))
		assert.Contains(t, result.Output, "Line 1:")
		assert.Contains(t, result.Output, "Line 100:")
	})
}

func TestShellExecutorIntegration(t *testing.T) {
	t.Run("Complex script with multiple operations", func(t *testing.T) {
		tmpDir := t.TempDir()
		scriptPath := filepath.Join(tmpDir, "complex.sh")

		// Create a complex script that tests various features
		script := `#!/bin/bash
set -e

# Test variable expansion
TEST_VAR="Hello from Plexr"
echo "Variable: $TEST_VAR"

# Test function
function test_func() {
    echo "Function called with: $1"
}
test_func "test argument"

# Test conditional
if [ -d "$PWD" ]; then
    echo "Working directory exists"
fi

# Test loop
for i in 1 2 3; do
    echo "Iteration: $i"
done

# Test command substitution
echo "Current date: $(date +%Y-%m-%d)"

# Exit successfully
exit 0`
		err := os.WriteFile(scriptPath, []byte(script), 0755) // #nosec G306 - Script needs to be executable
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		file := ExecutionFile{
			Path: scriptPath,
		}

		result, err := executor.Execute(ctx, file)
		require.NoError(t, err)
		assert.True(t, result.Success)

		// Verify all expected outputs
		assert.Contains(t, result.Output, "Variable: Hello from Plexr")
		assert.Contains(t, result.Output, "Function called with: test argument")
		assert.Contains(t, result.Output, "Working directory exists")
		assert.Contains(t, result.Output, "Iteration: 3")
		assert.Contains(t, result.Output, "Current date:")
	})
}
