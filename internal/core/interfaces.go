package core

import (
	"context"
)

// Executor defines the interface for all executors
type Executor interface {
	Name() string
	Execute(ctx context.Context, file ExecutionFile) (*ExecutionResult, error)
	Validate(config map[string]interface{}) error
}

// ExecutionFile represents a file to be executed
type ExecutionFile struct {
	Path     string
	Timeout  int
	Retry    int
	Platform string
}

// ExecutionResult represents the result of executing a file
type ExecutionResult struct {
	Success bool
	Output  string
	Error   error
	Duration int64 // in milliseconds
}
