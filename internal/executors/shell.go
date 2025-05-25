package executors

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// ExecutionFile represents a file to be executed
type ExecutionFile struct {
	Path          string
	Timeout       int
	Retry         int
	Platform      string
	WorkDirectory string
}

// ExecutionResult represents the result of executing a file
type ExecutionResult struct {
	Success  bool
	Output   string
	Error    error
	Duration int64 // in milliseconds
}

// ShellExecutor executes shell scripts
type ShellExecutor struct {
	shell string
}

// NewShellExecutor creates a new shell executor
func NewShellExecutor() *ShellExecutor {
	shell := "/bin/bash"
	if runtime.GOOS == "windows" {
		shell = "powershell.exe"
	}
	return &ShellExecutor{
		shell: shell,
	}
}

// Name returns the executor name
func (e *ShellExecutor) Name() string {
	return "shell"
}

// Validate validates the executor configuration
func (e *ShellExecutor) Validate(config map[string]interface{}) error {
	if shellPath, ok := config["shell"]; ok {
		if _, ok := shellPath.(string); !ok {
			return fmt.Errorf("shell must be a string")
		}
	}
	return nil
}

// Execute executes a shell script
func (e *ShellExecutor) Execute(ctx context.Context, file ExecutionFile) (*ExecutionResult, error) {
	start := time.Now()

	// Check platform compatibility
	if file.Platform != "" && file.Platform != runtime.GOOS {
		return &ExecutionResult{
			Success:  true,
			Output:   fmt.Sprintf("Skipping file %s (platform: %s, current: %s)", file.Path, file.Platform, runtime.GOOS),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	// Check if file exists
	if _, err := os.Stat(file.Path); err != nil {
		return &ExecutionResult{
			Success:  false,
			Error:    err,
			Duration: time.Since(start).Milliseconds(),
		}, fmt.Errorf("file not found: %s", file.Path)
	}

	// Create context with timeout if specified
	execCtx := ctx
	if file.Timeout > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, time.Duration(file.Timeout)*time.Second)
		defer cancel()
	}

	// Prepare command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(execCtx, e.shell, "-File", file.Path) // #nosec G204 - file.Path is validated and comes from user configuration
	} else {
		cmd = exec.CommandContext(execCtx, e.shell, file.Path) // #nosec G204 - file.Path is validated and comes from user configuration
	}

	// Set working directory if specified
	if file.WorkDirectory != "" {
		cmd.Dir = file.WorkDirectory
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute
	err := cmd.Run()

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}

	result := &ExecutionResult{
		Success:  err == nil,
		Output:   output,
		Error:    err,
		Duration: time.Since(start).Milliseconds(),
	}

	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("execution timeout after %d seconds", file.Timeout)
		}
		return result, fmt.Errorf("execution failed: %w", err)
	}

	return result, nil
}
