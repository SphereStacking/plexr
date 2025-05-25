package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/SphereStacking/plexr/internal/executors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockExecutor is a mock executor for testing
type MockExecutor struct {
	name            string
	executeFunc     func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error)
	validateFunc    func(config map[string]interface{}) error
	executeCalled   int
	executedFiles   []string
	executionDelays map[string]time.Duration
	mu              sync.Mutex
}

func (m *MockExecutor) Name() string {
	return m.name
}

func (m *MockExecutor) Execute(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
	m.mu.Lock()
	m.executeCalled++
	m.executedFiles = append(m.executedFiles, file.Path)
	m.mu.Unlock()

	// Simulate execution delay if configured
	if delay, ok := m.executionDelays[file.Path]; ok {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return &executors.ExecutionResult{
				Success: false,
				Error:   ctx.Err(),
			}, ctx.Err()
		}
	}

	if m.executeFunc != nil {
		return m.executeFunc(ctx, file)
	}
	return &executors.ExecutionResult{Success: true, Output: "mock output for " + file.Path}, nil
}

func (m *MockExecutor) Validate(config map[string]interface{}) error {
	if m.validateFunc != nil {
		return m.validateFunc(config)
	}
	return nil
}

func (m *MockExecutor) GetExecutedFiles() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	files := make([]string, len(m.executedFiles))
	copy(files, m.executedFiles)
	return files
}

func TestRunner(t *testing.T) {
	t.Run("NewRunner with valid plan", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:        "Test Plan",
			Version:     "1.0.0",
			Description: `Test plan for runner tests`,
			Executors: map[string]config.ExecutorConfig{
				"shell": {Type: "shell"},
			},
			Steps: []config.Step{
				{
					ID:          "test-step",
					Description: "Test step",
					Executor:    "shell",
					Files: []config.FileConfig{
						{Path: "test.sh"},
					},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)
		assert.NotNil(t, runner)
		assert.Equal(t, plan, runner.plan)
		assert.NotNil(t, runner.stateManager)
		assert.NotNil(t, runner.executors)

		// Verify shell executor is registered by default
		assert.Contains(t, runner.executors, "shell")
	})

	t.Run("RegisterExecutor", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		// Register new executor
		mockExec := &MockExecutor{name: "custom"}
		err = runner.RegisterExecutor("custom", mockExec)
		assert.NoError(t, err)
		assert.Contains(t, runner.executors, "custom")

		// Try to register with same name
		err = runner.RegisterExecutor("custom", mockExec)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")

		// Register with empty name
		err = runner.RegisterExecutor("", mockExec)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("Execute with single step", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Single Step Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:          "single-step",
					Description: "A single test step",
					Executor:    "mock",
					Files: []config.FileConfig{
						{Path: "script.sh"},
					},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		executed := false
		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				executed = true
				assert.Equal(t, "script.sh", file.Path)
				return &executors.ExecutionResult{
					Success:  true,
					Output:   "Step executed successfully",
					Duration: 100,
				}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)
		assert.True(t, executed)

		// Verify state
		state, err := runner.stateManager.Load()
		require.NoError(t, err)
		assert.Contains(t, state.CompletedSteps, "single-step")
		assert.Equal(t, "single-step", state.CurrentStep)
	})

	t.Run("Execute with dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Dependency Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "setup",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "setup.sh"}},
				},
				{
					ID:        "build",
					Executor:  "mock",
					DependsOn: []string{"setup"},
					Files:     []config.FileConfig{{Path: "build.sh"}},
				},
				{
					ID:        "test",
					Executor:  "mock",
					DependsOn: []string{"build"},
					Files:     []config.FileConfig{{Path: "test.sh"}},
				},
				{
					ID:        "deploy",
					Executor:  "mock",
					DependsOn: []string{"build", "test"},
					Files:     []config.FileConfig{{Path: "deploy.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{
			name:          "mock",
			executedFiles: []string{},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Verify execution order
		expectedOrder := []string{"setup.sh", "build.sh", "test.sh", "deploy.sh"}
		assert.Equal(t, expectedOrder, mockExec.GetExecutedFiles())
	})

	t.Run("Execute with circular dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Circular Dependency Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:        "step1",
					Executor:  "mock",
					DependsOn: []string{"step3"},
					Files:     []config.FileConfig{{Path: "step1.sh"}},
				},
				{
					ID:        "step2",
					Executor:  "mock",
					DependsOn: []string{"step1"},
					Files:     []config.FileConfig{{Path: "step2.sh"}},
				},
				{
					ID:        "step3",
					Executor:  "mock",
					DependsOn: []string{"step2"},
					Files:     []config.FileConfig{{Path: "step3.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{name: "mock"}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute should fail due to circular dependency
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circular dependency")
	})

	t.Run("Execute with work_directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")
		globalWorkDir := filepath.Join(tmpDir, "global_work")
		stepWorkDir := filepath.Join(tmpDir, "step_work")

		// Create work directories
		err := os.MkdirAll(globalWorkDir, 0755) // #nosec G301 - Test directory
		require.NoError(t, err)
		err = os.MkdirAll(stepWorkDir, 0755) // #nosec G301 - Test directory
		require.NoError(t, err)

		plan := &config.ExecutionPlan{
			Name:          "WorkDir Test",
			Version:       "1.0.0",
			WorkDirectory: globalWorkDir,
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:          "global-workdir-step",
					Description: "Uses global work directory",
					Executor:    "mock",
					Files:       []config.FileConfig{{Path: "global.sh"}},
				},
				{
					ID:            "custom-workdir-step",
					Description:   "Uses custom work directory",
					Executor:      "mock",
					WorkDirectory: stepWorkDir,
					Files:         []config.FileConfig{{Path: "custom.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		workDirs := make(map[string]string)
		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				workDirs[file.Path] = file.WorkDirectory
				return &executors.ExecutionResult{Success: true}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Verify work directories
		assert.Equal(t, globalWorkDir, workDirs["global.sh"])
		assert.Equal(t, stepWorkDir, workDirs["custom.sh"])
	})

	t.Run("Skip completed steps", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		// Create initial state with some steps completed
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		initialState := &ExecutionState{
			SetupName:      "Resume Test",
			SetupVersion:   "1.0.0",
			Platform:       "linux",
			CompletedSteps: []string{"step1", "step2"},
		}
		err = sm.Save(initialState)
		require.NoError(t, err)

		plan := &config.ExecutionPlan{
			Name:    "Resume Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "step1",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "step1.sh"}},
				},
				{
					ID:       "step2",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "step2.sh"}},
				},
				{
					ID:       "step3",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "step3.sh"}},
				},
				{
					ID:       "step4",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "step4.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{
			name:          "mock",
			executedFiles: []string{},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Only step3 and step4 should be executed
		assert.Equal(t, []string{"step3.sh", "step4.sh"}, mockExec.GetExecutedFiles())
	})

	t.Run("Execute with skip_if condition", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		// Create a test file to check in skip_if
		testFile := filepath.Join(tmpDir, "exists.txt")
		err := os.WriteFile(testFile, []byte("test"), 0600) // #nosec G306 - Test file
		require.NoError(t, err)

		plan := &config.ExecutionPlan{
			Name:    "Skip If Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "skip-always",
					Executor: "mock",
					SkipIf:   "true",
					Files:    []config.FileConfig{{Path: "skip-always.sh"}},
				},
				{
					ID:       "skip-never",
					Executor: "mock",
					SkipIf:   "false",
					Files:    []config.FileConfig{{Path: "skip-never.sh"}},
				},
				{
					ID:       "skip-if-file-exists",
					Executor: "mock",
					SkipIf:   "test -f " + testFile,
					Files:    []config.FileConfig{{Path: "skip-exists.sh"}},
				},
				{
					ID:       "skip-if-file-not-exists",
					Executor: "mock",
					SkipIf:   "test -f /non/existent/file",
					Files:    []config.FileConfig{{Path: "skip-not-exists.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{
			name:          "mock",
			executedFiles: []string{},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Since "test -f ..." is treated as a step ID (not a shell command),
		// and no steps with those IDs exist, they won't skip.
		// Only skip-always (with SkipIf: "true") will be skipped.
		expectedFiles := []string{"skip-never.sh", "skip-exists.sh", "skip-not-exists.sh"}
		assert.Equal(t, expectedFiles, mockExec.GetExecutedFiles())

		// Verify all steps are marked as completed (even skipped ones)
		state, err := runner.stateManager.Load()
		require.NoError(t, err)
		assert.Len(t, state.CompletedSteps, 4)
	})

	t.Run("Execute with multiple files per step", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Multi File Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:          "multi-file-step",
					Description: "Step with multiple files",
					Executor:    "mock",
					Files: []config.FileConfig{
						{Path: "file1.sh"},
						{Path: "file2.sh", Timeout: 10},
						{Path: "file3.sh", Platform: "linux"},
					},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		fileDetails := make(map[string]executors.ExecutionFile)
		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				fileDetails[file.Path] = file
				return &executors.ExecutionResult{Success: true}, nil
			},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.NoError(t, err)

		// Verify all files were executed with correct details
		assert.Len(t, fileDetails, 3)
		assert.Equal(t, 0, fileDetails["file1.sh"].Timeout)
		assert.Equal(t, 10, fileDetails["file2.sh"].Timeout)
		assert.Equal(t, "linux", fileDetails["file3.sh"].Platform)
	})

	t.Run("Execute with step failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Failure Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "success-step",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "success.sh"}},
				},
				{
					ID:       "fail-step",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "fail.sh"}},
				},
				{
					ID:       "after-fail-step",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "after-fail.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{
			name: "mock",
			executeFunc: func(ctx context.Context, file executors.ExecutionFile) (*executors.ExecutionResult, error) {
				if file.Path == "fail.sh" {
					return &executors.ExecutionResult{
						Success: false,
						Output:  "Step failed",
						Error:   errors.New("intentional failure"),
					}, errors.New("intentional failure")
				}
				return &executors.ExecutionResult{Success: true}, nil
			},
			executedFiles: []string{},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Execute
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute step fail-step")

		// Verify only steps before failure were executed
		assert.Equal(t, []string{"success.sh", "fail.sh"}, mockExec.GetExecutedFiles())

		// Verify state
		state, err := runner.stateManager.Load()
		require.NoError(t, err)
		assert.Contains(t, state.CompletedSteps, "success-step")
		assert.NotContains(t, state.CompletedSteps, "fail-step")
		assert.NotContains(t, state.CompletedSteps, "after-fail-step")
		assert.Equal(t, "fail-step", state.CurrentStep)
	})

	t.Run("Context cancellation", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Cancellation Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"mock": {Type: "mock"},
			},
			Steps: []config.Step{
				{
					ID:       "quick-step",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "quick.sh"}},
				},
				{
					ID:       "slow-step",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "slow.sh"}},
				},
				{
					ID:       "after-cancel-step",
					Executor: "mock",
					Files:    []config.FileConfig{{Path: "after-cancel.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		mockExec := &MockExecutor{
			name: "mock",
			executionDelays: map[string]time.Duration{
				"slow.sh": 5 * time.Second,
			},
			executedFiles: []string{},
		}
		err = runner.RegisterExecutor("mock", mockExec)
		require.NoError(t, err)

		// Create cancellable context
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		// Execute
		err = runner.Execute(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")

		// Verify only quick step was executed
		executedFiles := mockExec.GetExecutedFiles()
		assert.Contains(t, executedFiles, "quick.sh")
		assert.NotContains(t, executedFiles, "after-cancel.sh")
	})

	t.Run("Missing executor", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		plan := &config.ExecutionPlan{
			Name:    "Missing Executor Test",
			Version: "1.0.0",
			Executors: map[string]config.ExecutorConfig{
				"nonexistent": {Type: "nonexistent"},
			},
			Steps: []config.Step{
				{
					ID:       "test-step",
					Executor: "nonexistent",
					Files:    []config.FileConfig{{Path: "test.sh"}},
				},
			},
		}

		runner, err := NewRunner(plan, stateFile)
		require.NoError(t, err)

		// Execute should fail
		ctx := context.Background()
		err = runner.Execute(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "executor not found: nonexistent")
	})
}
