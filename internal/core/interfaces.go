package core

import (
	"context"

	"github.com/SphereStacking/plexr/internal/executors"
)

// Executor defines the interface for all executors
type Executor interface {
	Name() string
	Execute(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error)
	Validate(config map[string]interface{}) error
}
