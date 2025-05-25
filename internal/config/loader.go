package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadExecutionPlan loads and parses an execution plan from a YAML file
func LoadExecutionPlan(path string) (*ExecutionPlan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var plan ExecutionPlan
	if err := yaml.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := ValidateExecutionPlan(&plan); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &plan, nil
}

// ValidateExecutionPlan validates the execution plan
func ValidateExecutionPlan(plan *ExecutionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan is nil")
	}
	
	if plan.Name == "" {
		return fmt.Errorf("name is required")
	}

	if plan.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(plan.Steps) == 0 {
		return fmt.Errorf("at least one step is required")
	}

	// Check for duplicate step IDs
	stepIDs := make(map[string]bool)
	for _, step := range plan.Steps {
		if step.ID == "" {
			return fmt.Errorf("step ID is required")
		}
		if stepIDs[step.ID] {
			return fmt.Errorf("duplicate step id: %s", step.ID)
		}
		stepIDs[step.ID] = true
	}

	// Check that all steps reference defined executors
	for _, step := range plan.Steps {

		if step.Executor == "" {
			return fmt.Errorf("executor is required for step %s", step.ID)
		}

		if _, ok := plan.Executors[step.Executor]; !ok {
			return fmt.Errorf("undefined executor '%s' in step '%s'", step.Executor, step.ID)
		}

		if len(step.Files) == 0 {
			return fmt.Errorf("at least one file is required for step %s", step.ID)
		}

		// Validate file paths and platform values
		validPlatforms := map[string]bool{
			"":        true, // empty is valid (means all platforms)
			"linux":   true,
			"darwin":  true,
			"windows": true,
		}
		for _, file := range step.Files {
			// Validate file path
			if file.Path == "" {
				return fmt.Errorf("file path cannot be empty in step '%s'", step.ID)
			}
			if strings.HasPrefix(file.Path, "/") {
				return fmt.Errorf("file path cannot be absolute in step '%s': %s", step.ID, file.Path)
			}
			if strings.Contains(file.Path, "..") {
				return fmt.Errorf("file path cannot contain '..' in step '%s': %s", step.ID, file.Path)
			}
			
			// Validate platform
			if !validPlatforms[file.Platform] {
				return fmt.Errorf("invalid platform '%s' in step '%s'", file.Platform, step.ID)
			}
		}

		// Validate transaction mode
		validTransactionModes := map[string]bool{
			"":       true, // empty is valid (no transaction mode)
			"none":   true,
			"each":   true,
			"all":    true,
		}
		if !validTransactionModes[step.TransactionMode] {
			return fmt.Errorf("invalid transaction_mode '%s' in step '%s'", step.TransactionMode, step.ID)
		}
	}

	// Check for circular dependencies
	if err := checkCircularDependencies(plan.Steps); err != nil {
		return err
	}

	return nil
}

// checkCircularDependencies checks for circular dependencies in steps
func checkCircularDependencies(steps []Step) error {
	// Build a map of step IDs to their dependencies
	deps := make(map[string][]string)
	stepExists := make(map[string]bool)

	for _, step := range steps {
		deps[step.ID] = step.DependsOn
		stepExists[step.ID] = true
	}

	// Check that all dependencies exist
	for stepID, stepDeps := range deps {
		for _, dep := range stepDeps {
			if !stepExists[dep] {
				return fmt.Errorf("step '%s' depends on undefined step '%s'", stepID, dep)
			}
		}
	}

	// Check for circular dependencies using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(string) bool
	hasCycle = func(stepID string) bool {
		visited[stepID] = true
		recStack[stepID] = true

		for _, dep := range deps[stepID] {
			if !visited[dep] {
				if hasCycle(dep) {
					return true
				}
			} else if recStack[dep] {
				return true
			}
		}

		recStack[stepID] = false
		return false
	}

	for stepID := range deps {
		if !visited[stepID] {
			if hasCycle(stepID) {
				return fmt.Errorf("circular dependency detected")
			}
		}
	}

	return nil
}
