package core

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/SphereStacking/plexr/internal/executors"
)

// Runner manages the execution of an execution plan
type Runner struct {
	plan         *config.ExecutionPlan
	stateManager *StateManager
	executors    map[string]Executor
	platform     string
	// Progress tracking
	progressCallback func(stepID string, event string, data interface{})
}

// NewRunner creates a new runner
func NewRunner(plan *config.ExecutionPlan, stateFile string) (*Runner, error) {
	sm, err := NewStateManager(stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create state manager: %w", err)
	}

	r := &Runner{
		plan:         plan,
		stateManager: sm,
		executors:    make(map[string]Executor),
		platform:     runtime.GOOS,
	}

	// Register built-in executors
	r.executors["shell"] = executors.NewShellExecutor()
	r.executors["sql"] = executors.NewSQLExecutor()

	// Initialize and validate executors with their configurations
	if plan.Executors != nil {
		for name, config := range plan.Executors {
			executorType, ok := config["type"].(string)
			if !ok {
				return nil, fmt.Errorf("executor %s missing type field", name)
			}

			// Skip if it's already registered (e.g., from tests)
			if _, exists := r.executors[name]; exists {
				continue
			}

			// Get the executor by type
			baseExecutor, ok := r.executors[executorType]
			if !ok {
				// If not a built-in type, skip - it might be registered later in tests
				continue
			}

			// For SQL executor, create a new instance for each configuration
			var executor Executor
			if executorType == "sql" {
				sqlExecutor := executors.NewSQLExecutor()
				// Validate the configuration
				if err := sqlExecutor.Validate(config); err != nil {
					return nil, fmt.Errorf("invalid configuration for executor %s: %w", name, err)
				}
				executor = sqlExecutor
			} else {
				// For other executors, validate with the base executor
				if err := baseExecutor.Validate(config); err != nil {
					return nil, fmt.Errorf("invalid configuration for executor %s: %w", name, err)
				}
				executor = baseExecutor
			}

			// Store the configured executor with the custom name
			r.executors[name] = executor
		}
	}

	return r, nil
}

// RegisterExecutor registers a custom executor
func (r *Runner) RegisterExecutor(name string, executor Executor) error {
	if name == "" {
		return fmt.Errorf("executor name cannot be empty")
	}
	if _, exists := r.executors[name]; exists {
		return fmt.Errorf("executor %s already registered", name)
	}
	r.executors[name] = executor
	return nil
}

// SetProgressCallback sets a callback function for progress notifications
func (r *Runner) SetProgressCallback(callback func(stepID string, event string, data interface{})) {
	r.progressCallback = callback
}

// notifyProgress sends a progress notification if a callback is registered
func (r *Runner) notifyProgress(stepID string, event string, data interface{}) {
	if r.progressCallback != nil {
		r.progressCallback(stepID, event, data)
	}
}

// State returns the current execution state
func (r *Runner) State() *ExecutionState {
	state, _ := r.stateManager.Load()
	return state
}

// Execute runs the execution plan
func (r *Runner) Execute(ctx context.Context) error {
	// Load or create state
	_, err := r.stateManager.Load()
	if err != nil {
		// Create new state
		state := &ExecutionState{
			SetupName:      r.plan.Name,
			SetupVersion:   r.plan.Version,
			Platform:       r.platform,
			StartedAt:      time.Now(),
			CompletedSteps: []string{},
			InstalledTools: make(map[string]string),
		}
		err = r.stateManager.Save(state)
		if err != nil {
			return fmt.Errorf("failed to save initial state: %w", err)
		}
	}

	// Build execution order based on dependencies
	order, err := r.buildExecutionOrder()
	if err != nil {
		return fmt.Errorf("failed to build execution order: %w", err)
	}

	// Execute steps in order
	for _, stepID := range order {
		step := r.findStep(stepID)
		if step == nil {
			return fmt.Errorf("step not found: %s", stepID)
		}

		// Check if already completed
		if r.stateManager.IsStepCompleted(stepID) {
			r.notifyProgress(stepID, "skipped", map[string]interface{}{"reason": "already_completed"})
			continue
		}

		// Check skip condition
		if step.SkipIf != "" {
			// Evaluate skip condition
			shouldSkip := false
			switch step.SkipIf {
			case "true":
				shouldSkip = true
			case "false":
				shouldSkip = false
			default:
				// Check if the condition is a step ID that has been completed
				shouldSkip = r.stateManager.IsStepCompleted(step.SkipIf)
			}

			if shouldSkip {
				r.notifyProgress(stepID, "skipped", map[string]interface{}{"reason": "skip_if_condition", "condition": step.SkipIf})
				// Mark as completed even if skipped
				err = r.stateManager.MarkStepCompleted(stepID)
				if err != nil {
					return fmt.Errorf("failed to mark skipped step %s as completed: %w", stepID, err)
				}
				continue
			}
		}

		// Execute step
		err = r.executeStep(ctx, step)
		if err != nil {
			r.notifyProgress(stepID, "failed", map[string]interface{}{"error": err.Error()})
			return fmt.Errorf("failed to execute step %s: %w", stepID, err)
		}

		// Mark as completed
		err = r.stateManager.MarkStepCompleted(stepID)
		if err != nil {
			return fmt.Errorf("failed to mark step %s as completed: %w", stepID, err)
		}
		r.notifyProgress(stepID, "completed", nil)
	}

	return nil
}

// executeStep executes a single step
func (r *Runner) executeStep(ctx context.Context, step *config.Step) error {
	executor, ok := r.executors[step.Executor]
	if !ok {
		return fmt.Errorf("executor not found: %s", step.Executor)
	}

	r.notifyProgress(step.ID, "started", map[string]interface{}{"description": step.Description})
	err := r.stateManager.SetCurrentStep(step.ID)
	if err != nil {
		return err
	}

	for _, fileConfig := range step.Files {
		// Use step work_directory if specified, otherwise use global work_directory
		workDir := step.WorkDirectory
		if workDir == "" {
			workDir = r.plan.WorkDirectory
		}

		file := executors.ExecutionFile{
			Path:            fileConfig.Path,
			Timeout:         fileConfig.Timeout,
			Retry:           fileConfig.Retry,
			Platform:        fileConfig.Platform,
			WorkDirectory:   workDir,
			TransactionMode: step.TransactionMode,
		}

		r.notifyProgress(step.ID, "executing_file", map[string]interface{}{"file": file.Path})
		result, err := executor.Execute(ctx, file)
		if err != nil {
			return err
		}

		if !result.Success {
			return fmt.Errorf("execution failed")
		}

		// Show output if available
		if result.Output != "" {
			r.notifyProgress(step.ID, "output", map[string]interface{}{"output": result.Output})
		}
	}

	return nil
}

// buildExecutionOrder builds the execution order based on dependencies
func (r *Runner) buildExecutionOrder() ([]string, error) {
	var order []string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var visit func(stepID string) error
	visit = func(stepID string) error {
		if visited[stepID] {
			return nil
		}

		if recStack[stepID] {
			return fmt.Errorf("circular dependency detected")
		}

		recStack[stepID] = true

		step := r.findStep(stepID)
		if step == nil {
			return fmt.Errorf("step not found: %s", stepID)
		}

		for _, dep := range step.DependsOn {
			err := visit(dep)
			if err != nil {
				return err
			}
		}

		recStack[stepID] = false
		visited[stepID] = true
		order = append(order, stepID)
		return nil
	}

	for _, step := range r.plan.Steps {
		if !visited[step.ID] {
			err := visit(step.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return order, nil
}

// findStep finds a step by ID
func (r *Runner) findStep(id string) *config.Step {
	for i := range r.plan.Steps {
		if r.plan.Steps[i].ID == id {
			return &r.plan.Steps[i]
		}
	}
	return nil
}
