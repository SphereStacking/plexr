package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateManager(t *testing.T) {
	t.Run("NewStateManager", func(t *testing.T) {
		tests := []struct {
			name      string
			setupFunc func() string
			wantErr   bool
			errMsg    string
		}{
			{
				name: "valid file path",
				setupFunc: func() string {
					tmpDir := t.TempDir()
					return filepath.Join(tmpDir, "state.json")
				},
				wantErr: false,
			},
			{
				name: "create missing directory",
				setupFunc: func() string {
					tmpDir := t.TempDir()
					return filepath.Join(tmpDir, "subdir", "state.json")
				},
				wantErr: false,
			},
			{
				name: "invalid path",
				setupFunc: func() string {
					return ""
				},
				wantErr: true,
				errMsg:  "state file path cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filePath := tt.setupFunc()
				sm, err := NewStateManager(filePath)

				if tt.wantErr {
					assert.Error(t, err)
					if tt.errMsg != "" && err != nil {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, sm)
					assert.Equal(t, filePath, sm.filePath)
					// Note: We can't directly check sm.mu due to copylocks warning
					// The mutex is properly initialized by the constructor
				}
			})
		}
	})

	t.Run("Save and Load state", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Test various state configurations
		tests := []struct {
			name  string
			state *ExecutionState
		}{
			{
				name: "minimal state",
				state: &ExecutionState{
					SetupName:    "Minimal",
					SetupVersion: "1.0.0",
				},
			},
			{
				name: "complete state",
				state: &ExecutionState{
					SetupName:      "Complete Setup",
					SetupVersion:   "2.0.0",
					Platform:       "linux",
					StartedAt:      time.Now().Add(-1 * time.Hour),
					UpdatedAt:      time.Now(),
					CompletedSteps: []string{"step1", "step2", "step3"},
					CurrentStep:    "step4",
					FailedFiles:    []string{"error.sh"},
					InstalledTools: map[string]string{
						"docker": "24.0.7",
						"node":   "20.10.0",
						"go":     "1.21.5",
					},
				},
			},
			{
				name: "state with special characters",
				state: &ExecutionState{
					SetupName:      "Test/Setup with \"quotes\"",
					SetupVersion:   "1.0.0-beta",
					CurrentStep:    "step with spaces & symbols!",
					CompletedSteps: []string{"step-1", "step_2", "step.3"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Save state
				err := sm.Save(tt.state)
				require.NoError(t, err)

				// Verify file exists and is valid JSON
				data, err := os.ReadFile(stateFile)
				require.NoError(t, err)

				var jsonState map[string]interface{}
				err = json.Unmarshal(data, &jsonState)
				require.NoError(t, err)

				// Load state
				loaded, err := sm.Load()
				require.NoError(t, err)

				// Verify fields
				assert.Equal(t, tt.state.SetupName, loaded.SetupName)
				assert.Equal(t, tt.state.SetupVersion, loaded.SetupVersion)
				assert.Equal(t, tt.state.Platform, loaded.Platform)
				assert.Equal(t, tt.state.CurrentStep, loaded.CurrentStep)
				assert.ElementsMatch(t, tt.state.CompletedSteps, loaded.CompletedSteps)
				assert.ElementsMatch(t, tt.state.FailedFiles, loaded.FailedFiles)

				// Verify maps
				if tt.state.InstalledTools != nil {
					assert.Equal(t, tt.state.InstalledTools, loaded.InstalledTools)
				}
			})
		}
	})

	t.Run("State persistence across instances", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		// First instance saves state
		sm1, err := NewStateManager(stateFile)
		require.NoError(t, err)

		originalState := &ExecutionState{
			SetupName:      "Persistence Test",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{"step1", "step2"},
			InstalledTools: map[string]string{"tool1": "v1"},
		}
		err = sm1.Save(originalState)
		require.NoError(t, err)

		// Second instance loads state
		sm2, err := NewStateManager(stateFile)
		require.NoError(t, err)

		loadedState, err := sm2.Load()
		require.NoError(t, err)

		assert.Equal(t, originalState.SetupName, loadedState.SetupName)
		assert.Equal(t, originalState.SetupVersion, loadedState.SetupVersion)
		assert.ElementsMatch(t, originalState.CompletedSteps, loadedState.CompletedSteps)
		assert.Equal(t, originalState.InstalledTools, loadedState.InstalledTools)
	})

	t.Run("Concurrent operations", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Initialize state
		initialState := &ExecutionState{
			SetupName:      "Concurrent Test",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{},
		}
		err = sm.Save(initialState)
		require.NoError(t, err)

		// Concurrent updates
		var wg sync.WaitGroup
		stepCount := 10

		for i := 0; i < stepCount; i++ {
			wg.Add(1)
			go func(stepNum int) {
				defer wg.Done()
				stepID := fmt.Sprintf("step%d", stepNum)
				err := sm.MarkStepCompleted(stepID)
				assert.NoError(t, err)
			}(i)
		}

		wg.Wait()

		// Verify all steps were saved
		finalState, err := sm.Load()
		require.NoError(t, err)
		assert.Len(t, finalState.CompletedSteps, stepCount)
	})

	t.Run("Reset functionality", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Save state
		state := &ExecutionState{
			SetupName:      "Reset Test",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{"step1", "step2"},
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(stateFile)
		require.NoError(t, err)

		// Reset
		err = sm.Reset()
		require.NoError(t, err)

		// Verify file is removed
		_, err = os.Stat(stateFile)
		assert.True(t, os.IsNotExist(err))

		// Verify internal state is cleared - state manager doesn't expose internal state

		// Load should fail
		loaded, err := sm.Load()
		assert.Error(t, err)
		assert.Nil(t, loaded)
	})

	t.Run("Step completion tracking", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Initialize
		state := &ExecutionState{
			SetupName:      "Step Tracking Test",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{"existing-step"},
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Test IsStepCompleted
		_, err = sm.Load() // Load state into memory
		require.NoError(t, err)
		assert.True(t, sm.IsStepCompleted("existing-step"))
		assert.False(t, sm.IsStepCompleted("new-step"))
		assert.False(t, sm.IsStepCompleted(""))

		// Mark new step as completed
		err = sm.MarkStepCompleted("new-step")
		require.NoError(t, err)
		assert.True(t, sm.IsStepCompleted("new-step"))

		// Verify persistence
		loaded, err := sm.Load()
		require.NoError(t, err)
		assert.Contains(t, loaded.CompletedSteps, "existing-step")
		assert.Contains(t, loaded.CompletedSteps, "new-step")

		// Mark duplicate step (should not duplicate)
		err = sm.MarkStepCompleted("new-step")
		require.NoError(t, err)

		loaded, err = sm.Load()
		require.NoError(t, err)

		// Count occurrences
		count := 0
		for _, step := range loaded.CompletedSteps {
			if step == "new-step" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})

	t.Run("SetCurrentStep", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Initialize
		state := &ExecutionState{
			SetupName:    "Current Step Test",
			SetupVersion: "1.0.0",
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Set current step
		err = sm.SetCurrentStep("step1")
		require.NoError(t, err)

		// Verify
		loaded, err := sm.Load()
		require.NoError(t, err)
		assert.Equal(t, "step1", loaded.CurrentStep)

		// Update current step
		err = sm.SetCurrentStep("step2")
		require.NoError(t, err)

		loaded, err = sm.Load()
		require.NoError(t, err)
		assert.Equal(t, "step2", loaded.CurrentStep)
	})

	t.Run("Corrupted state file", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		// Write invalid JSON
		err := os.WriteFile(stateFile, []byte("{ invalid json"), 0600) // #nosec G306 - Test file
		require.NoError(t, err)

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Load should fail
		loaded, err := sm.Load()
		assert.Error(t, err)
		assert.Nil(t, loaded)
		assert.Contains(t, err.Error(), "failed to parse state file")
	})

	t.Run("File permissions", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Save initial state
		state := &ExecutionState{
			SetupName:    "Permission Test",
			SetupVersion: "1.0.0",
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Make file read-only
		err = os.Chmod(stateFile, 0444) // #nosec G302 - Test file permissions
		require.NoError(t, err)

		// Save should fail
		state.CompletedSteps = []string{"step1"}
		err = sm.Save(state)
		assert.Error(t, err)

		// Restore permissions for cleanup
		_ = os.Chmod(stateFile, 0644) // #nosec G302 - Test file permissions
	})

	t.Run("UpdatedAt timestamp", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Save initial state
		state := &ExecutionState{
			SetupName:    "Timestamp Test",
			SetupVersion: "1.0.0",
			UpdatedAt:    time.Now().Add(-1 * time.Hour),
		}
		originalTime := state.UpdatedAt
		err = sm.Save(state)
		require.NoError(t, err)

		// Update state
		time.Sleep(10 * time.Millisecond)
		err = sm.MarkStepCompleted("step1")
		require.NoError(t, err)

		// Verify UpdatedAt was updated
		loaded, err := sm.Load()
		require.NoError(t, err)
		assert.True(t, loaded.UpdatedAt.After(originalTime))
	})

	t.Run("Empty state operations", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "state.json")

		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Operations on non-existent state
		assert.False(t, sm.IsStepCompleted("any-step"))

		err = sm.SetCurrentStep("step1")
		assert.Error(t, err)

		err = sm.MarkStepCompleted("step1")
		assert.Error(t, err)
	})
}

func TestExecutionState(t *testing.T) {
	t.Run("JSON marshaling", func(t *testing.T) {
		state := &ExecutionState{
			SetupName:      "JSON Test",
			SetupVersion:   "1.0.0",
			Platform:       "linux",
			StartedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			CompletedSteps: []string{"step1", "step2"},
			CurrentStep:    "step3",
			FailedFiles:    []string{"fail.sh"},
			InstalledTools: map[string]string{
				"tool1": "v1.0",
				"tool2": "v2.0",
			},
		}

		// Marshal
		data, err := json.MarshalIndent(state, "", "  ")
		require.NoError(t, err)

		// Unmarshal
		var loaded ExecutionState
		err = json.Unmarshal(data, &loaded)
		require.NoError(t, err)

		// Verify
		assert.Equal(t, state.SetupName, loaded.SetupName)
		assert.Equal(t, state.SetupVersion, loaded.SetupVersion)
		assert.Equal(t, state.Platform, loaded.Platform)
		assert.Equal(t, state.CurrentStep, loaded.CurrentStep)
		assert.ElementsMatch(t, state.CompletedSteps, loaded.CompletedSteps)
		assert.ElementsMatch(t, state.FailedFiles, loaded.FailedFiles)
		assert.Equal(t, state.InstalledTools, loaded.InstalledTools)

		// Times should be close (accounting for JSON marshaling precision)
		assert.WithinDuration(t, state.StartedAt, loaded.StartedAt, time.Second)
		assert.WithinDuration(t, state.UpdatedAt, loaded.UpdatedAt, time.Second)
	})
}
