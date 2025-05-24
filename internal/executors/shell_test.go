package executors

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/SphereStacking/plexr/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellExecutor(t *testing.T) {
	t.Run("NewShellExecutor", func(t *testing.T) {
		executor := NewShellExecutor()
		assert.NotNil(t, executor)
		assert.Equal(t, "shell", executor.Name())
	})

	t.Run("Validate", func(t *testing.T) {
		executor := NewShellExecutor()

		// Valid config
		config := map[string]interface{}{
			"shell": "/bin/bash",
		}
		err := executor.Validate(config)
		assert.NoError(t, err)

		// Empty config (should use default)
		err = executor.Validate(map[string]interface{}{})
		assert.NoError(t, err)

		// Invalid shell type
		config = map[string]interface{}{
			"shell": 123,
		}
		err = executor.Validate(config)
		assert.Error(t, err)
	})

	t.Run("Execute simple command", func(t *testing.T) {
		tmpDir := t.TempDir()
		scriptPath := filepath.Join(tmpDir, "test.sh")

		// Create a simple test script
		script := `#!/bin/bash
echo "Hello, Plexr!"
exit 0`
		err := os.WriteFile(scriptPath, []byte(script), 0755)
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		file := core.ExecutionFile{
			Path: scriptPath,
		}

		result, err := executor.Execute(ctx, file)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Contains(t, result.Output, "Hello, Plexr!")
		assert.Greater(t, result.Duration, int64(0))
	})

	t.Run("Execute with error", func(t *testing.T) {
		tmpDir := t.TempDir()
		scriptPath := filepath.Join(tmpDir, "error.sh")

		// Create a script that fails
		script := `#!/bin/bash
echo "This will fail"
exit 1`
		err := os.WriteFile(scriptPath, []byte(script), 0755)
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		file := core.ExecutionFile{
			Path: scriptPath,
		}

		result, err := executor.Execute(ctx, file)
		assert.Error(t, err)
		assert.False(t, result.Success)
		assert.Contains(t, result.Output, "This will fail")
	})

	t.Run("Execute with timeout", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping timeout test on Windows")
		}

		tmpDir := t.TempDir()
		scriptPath := filepath.Join(tmpDir, "timeout.sh")

		// Create a script that takes too long
		script := `#!/bin/bash
sleep 5
echo "Should not see this"`
		err := os.WriteFile(scriptPath, []byte(script), 0755)
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		file := core.ExecutionFile{
			Path:    scriptPath,
			Timeout: 1, // 1 second timeout
		}

		result, err := executor.Execute(ctx, file)
		assert.Error(t, err)
		assert.False(t, result.Success)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("Execute non-existent file", func(t *testing.T) {
		executor := NewShellExecutor()
		ctx := context.Background()

		file := core.ExecutionFile{
			Path: "/non/existent/script.sh",
		}

		result, err := executor.Execute(ctx, file)
		assert.Error(t, err)
		assert.False(t, result.Success)
	})

	t.Run("Platform-specific execution", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		// Create platform-specific scripts
		darwinScript := filepath.Join(tmpDir, "darwin.sh")
		linuxScript := filepath.Join(tmpDir, "linux.sh")

		err := os.WriteFile(darwinScript, []byte(`#!/bin/bash
echo "Running on macOS"`), 0755)
		require.NoError(t, err)

		err = os.WriteFile(linuxScript, []byte(`#!/bin/bash
echo "Running on Linux"`), 0755)
		require.NoError(t, err)

		executor := NewShellExecutor()
		ctx := context.Background()

		// Test with platform filter
		file := core.ExecutionFile{
			Path:     darwinScript,
			Platform: "darwin",
		}

		result, err := executor.Execute(ctx, file)
		
		if runtime.GOOS == "darwin" {
			require.NoError(t, err)
			assert.True(t, result.Success)
			assert.Contains(t, result.Output, "Running on macOS")
		} else {
			// Should skip on non-darwin platforms
			require.NoError(t, err)
			assert.True(t, result.Success)
			assert.Contains(t, result.Output, "Skipping")
		}
	})
}
