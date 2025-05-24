package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateManager(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "test_state.json")

	t.Run("NewStateManager", func(t *testing.T) {
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)
		assert.NotNil(t, sm)
		assert.Equal(t, stateFile, sm.filePath)
	})

	t.Run("Save and Load empty state", func(t *testing.T) {
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Create initial state
		state := &ExecutionState{
			SetupName:      "Test Setup",
			SetupVersion:   "1.0.0",
			Platform:       "darwin",
			StartedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			CompletedSteps: []string{},
			CurrentStep:    "",
			FailedFiles:    []string{},
			InstalledTools: map[string]string{},
		}

		// Save state
		err = sm.Save(state)
		require.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(stateFile)
		require.NoError(t, err)

		// Load state
		loaded, err := sm.Load()
		require.NoError(t, err)
		assert.Equal(t, state.SetupName, loaded.SetupName)
		assert.Equal(t, state.SetupVersion, loaded.SetupVersion)
		assert.Equal(t, state.Platform, loaded.Platform)
	})

	t.Run("Update state", func(t *testing.T) {
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Initial state
		state := &ExecutionState{
			SetupName:      "Test Setup",
			SetupVersion:   "1.0.0",
			Platform:       "darwin",
			StartedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			CompletedSteps: []string{"step1"},
			CurrentStep:    "step2",
			FailedFiles:    []string{},
			InstalledTools: map[string]string{
				"docker": "24.0.7",
			},
		}

		err = sm.Save(state)
		require.NoError(t, err)

		// Update state
		state.CompletedSteps = append(state.CompletedSteps, "step2")
		state.CurrentStep = "step3"
		state.InstalledTools["node"] = "18.19.0"
		state.UpdatedAt = time.Now()

		err = sm.Save(state)
		require.NoError(t, err)

		// Load and verify
		loaded, err := sm.Load()
		require.NoError(t, err)
		assert.Len(t, loaded.CompletedSteps, 2)
		assert.Equal(t, "step3", loaded.CurrentStep)
		assert.Equal(t, "18.19.0", loaded.InstalledTools["node"])
	})

	t.Run("Load non-existent state", func(t *testing.T) {
		nonExistentFile := filepath.Join(tmpDir, "non_existent.json")
		sm, err := NewStateManager(nonExistentFile)
		require.NoError(t, err)

		loaded, err := sm.Load()
		assert.Error(t, err)
		assert.Nil(t, loaded)
	})

	t.Run("Reset state", func(t *testing.T) {
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		// Save some state
		state := &ExecutionState{
			SetupName:      "Test Setup",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{"step1", "step2"},
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Reset
		err = sm.Reset()
		require.NoError(t, err)

		// Verify file is gone
		_, err = os.Stat(stateFile)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("IsStepCompleted", func(t *testing.T) {
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		state := &ExecutionState{
			SetupName:      "Test Setup",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{"step1", "step2"},
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Load state internally
		_, err = sm.Load()
		require.NoError(t, err)

		assert.True(t, sm.IsStepCompleted("step1"))
		assert.True(t, sm.IsStepCompleted("step2"))
		assert.False(t, sm.IsStepCompleted("step3"))
	})

	t.Run("MarkStepCompleted", func(t *testing.T) {
		sm, err := NewStateManager(stateFile)
		require.NoError(t, err)

		state := &ExecutionState{
			SetupName:      "Test Setup",
			SetupVersion:   "1.0.0",
			CompletedSteps: []string{},
		}
		err = sm.Save(state)
		require.NoError(t, err)

		// Mark step as completed
		err = sm.MarkStepCompleted("step1")
		require.NoError(t, err)

		// Verify
		loaded, err := sm.Load()
		require.NoError(t, err)
		assert.Contains(t, loaded.CompletedSteps, "step1")
		assert.True(t, sm.IsStepCompleted("step1"))

		// Mark another step
		err = sm.MarkStepCompleted("step2")
		require.NoError(t, err)

		loaded, err = sm.Load()
		require.NoError(t, err)
		assert.Len(t, loaded.CompletedSteps, 2)
	})
}
