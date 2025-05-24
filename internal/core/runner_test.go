package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/SphereStacking/plexr/internal/executors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockExecutor is a mock executor for testing
type MockExecutor struct {
	name          string
	executeFunc   func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error)
	validateFunc  func(config map[string]interface{}) error
	executeCalled int
}

func (m *MockExecutor) Name() string {
	return m.name
}

func (m *MockExecutor) Execute(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
	m.executeCalled++
	if m.executeFunc != nil {
		return m.executeFunc(ctx, file)
	}
	return &executors.ExecutionResult{Success: true, Output: "mock output"}, nil
}

func (m *MockExecutor) Validate(config map[string]interface{}) error {
	if m.validateFunc != nil {
		return m.validateFunc(config)
	}
	return nil
}

func TestRunner(t *testing.T) {
	t.Run("NewRunner", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state_newrunner.json")

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)
		assert.NotNil(t, runner)
		assert.Equal(t, plan, runner.plan)
		assert.NotNil(t, runner.stateManager)
		assert.NotNil(t, runner.executors)
	})

	t.Run("RegisterExecutor", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state_registerexecutor.json")

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{name: "mock"}
		err = runner.RegisterExecutor("mock", mockExec)
		assert.NoError(t, err)

		// Try to register with same name
		err = runner.RegisterExecutor("mock", mockExec)
		assert.Error(t, err)
	})

	t.Run("Execute simple plan", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state_simple.json")

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "step1",
					Executor: "mock",
					Files: []config.FileConfig{
						{Path: "file1.sh"},
					},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				return &executors.ExecutionResult{
					Success: true,
					Output:  "executed " + file.Path,
				}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, mockExec.executeCalled)

		// Check state
		state, err := runner.stateManager.Load()
		require.NoError(t, err)
		assert.Contains(t, state.CompletedSteps, "step1")
	})

	t.Run("Execute with dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state_deps.json")

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "step1",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "file1.sh"}},
				},
				{
					ID:        "step2",
					Executor:  "mock",
					DependsOn: []string{"step1"},
					Files:     []config.FileConfig{{Path: "file2.sh"}},
				},
				{
					ID:        "step3",
					Executor:  "mock",
					DependsOn: []string{"step1", "step2"},
					Files:     []config.FileConfig{{Path: "file3.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		executionOrder := []string{}
		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				executionOrder = append(executionOrder, file.Path)
				return &executors.ExecutionResult{Success: true}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Verify execution order
		assert.Equal(t, []string{"file1.sh", "file2.sh", "file3.sh"}, executionOrder)
	})

	t.Run("Skip completed steps", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state_skip.json")

		// Create initial state with step1 completed
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		initialState := &ExecutionState{
			SetupName:      "Test",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{"step1"},
		}
		err = sm.Save(initialState)
		require.NoError(t, err)

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "step1",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "file1.sh"}},
				},
				{
					ID:       "step2",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "file2.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		executedFiles := []string{}
		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				executedFiles = append(executedFiles, file.Path)
				return &executors.ExecutionResult{Success: true}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Only step2 should be executed
		assert.Equal(t, []string{"file2.sh"}, executedFiles)
	})

	t.Run("Execute with skip_if condition", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state_skipif.json")

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "step1",
					Executor: "mock",
					SkipIf:   "true", // Always skip
					Files:    []config.FileConfig{{Path: "file1.sh"}},
				},
				{
					ID:       "step2",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "file2.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		executedFiles := []string{}
		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				executedFiles = append(executedFiles, file.Path)
				return &executors.ExecutionResult{Success: true}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Only step2 should be executed
		assert.Equal(t, []string{"file2.sh"}, executedFiles)
	})
}
