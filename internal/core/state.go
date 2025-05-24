package core

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// ExecutionState represents the current state of an execution
type ExecutionState struct {
	SetupName      string            `json:"setup_name"`
	SetupVersion   string            `json:"setup_version"`
	Platform       string            `json:"platform"`
	StartedAt      time.Time         `json:"started_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	CompletedSteps []string          `json:"completed_steps"`
	CurrentStep    string            `json:"current_step"`
	FailedFiles    []string          `json:"failed_files"`
	InstalledTools map[string]string `json:"installed_tools"`
}

// StateManager manages the execution state
type StateManager struct {
	filePath string
	state    *ExecutionState
	mu       sync.RWMutex
}

// NewStateManager creates a new state manager
func NewStateManager(filePath string) (*StateManager, error) {
	return &StateManager{
		filePath: filePath,
	}, nil
}

// Load loads the state from file
func (sm *StateManager) Load() (*ExecutionState, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state ExecutionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	sm.state = &state
	return &state, nil
}

// Save saves the state to file
func (sm *StateManager) Save(state *ExecutionState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	state.UpdatedAt = time.Now()
	sm.state = state

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(sm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// Reset removes the state file
func (sm *StateManager) Reset() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if err := os.Remove(sm.filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove state file: %w", err)
	}

	sm.state = nil
	return nil
}

// IsStepCompleted checks if a step has been completed
func (sm *StateManager) IsStepCompleted(stepID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.state == nil {
		return false
	}

	for _, completed := range sm.state.CompletedSteps {
		if completed == stepID {
			return true
		}
	}
	return false
}

// MarkStepCompleted marks a step as completed
func (sm *StateManager) MarkStepCompleted(stepID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state == nil {
		return fmt.Errorf("state not loaded")
	}

	// Check if already completed
	for _, completed := range sm.state.CompletedSteps {
		if completed == stepID {
			return nil
		}
	}

	sm.state.CompletedSteps = append(sm.state.CompletedSteps, stepID)
	sm.state.UpdatedAt = time.Now()

	// Save immediately
	data, err := json.MarshalIndent(sm.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(sm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// SetCurrentStep sets the current step being executed
func (sm *StateManager) SetCurrentStep(stepID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state == nil {
		return fmt.Errorf("state not loaded")
	}

	sm.state.CurrentStep = stepID
	sm.state.UpdatedAt = time.Now()

	// Save immediately
	data, err := json.MarshalIndent(sm.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(sm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// AddInstalledTool records an installed tool and its version
func (sm *StateManager) AddInstalledTool(name, version string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state == nil {
		return fmt.Errorf("state not loaded")
	}

	if sm.state.InstalledTools == nil {
		sm.state.InstalledTools = make(map[string]string)
	}

	sm.state.InstalledTools[name] = version
	sm.state.UpdatedAt = time.Now()

	// Save immediately
	data, err := json.MarshalIndent(sm.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(sm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}
