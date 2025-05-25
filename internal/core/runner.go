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
			fmt.Printf("Skipping completed step: %s\n", stepID)
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
				fmt.Printf("Skipping step %s due to skip_if condition\n", stepID)
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
			return fmt.Errorf("failed to execute step %s: %w", stepID, err)
		}

		// Mark as completed
		err = r.stateManager.MarkStepCompleted(stepID)
		if err != nil {
			return fmt.Errorf("failed to mark step %s as completed: %w", stepID, err)
		}
	}

	return nil
}

// executeStep executes a single step
func (r *Runner) executeStep(ctx context.Context, step *config.Step) error {
	executor, ok := r.executors[step.Executor]
	if !ok {
		return fmt.Errorf("executor not found: %s", step.Executor)
	}

	fmt.Printf("\nExecuting step: %s - %s\n", step.ID, step.Description)
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
			Path:          fileConfig.Path,
			Timeout:       fileConfig.Timeout,
			Retry:         fileConfig.Retry,
			Platform:      fileConfig.Platform,
			WorkDirectory: workDir,
		}

		fmt.Printf("  Executing file: %s\n", file.Path)
		result, err := executor.Execute(ctx, file)
		if err != nil {
			return err
		}

		if !result.Success {
			return fmt.Errorf("execution failed")
		}

		fmt.Printf("  ✓ Success (%dms)\n", result.Duration)

		// Show output if available
		if result.Output != "" {
			fmt.Println("  Output:")
			fmt.Println(result.Output)
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
