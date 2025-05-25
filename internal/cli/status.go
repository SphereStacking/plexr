package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SphereStacking/plexr/internal/core"
	"github.com/spf13/cobra"
)

// NewStatusCommand creates the status command
func NewStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status <plan.yml>",
		Short: "Show execution status",
		Long:  `Display the current execution status of a setup plan.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runStatus,
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	planFile := args[0]
	dir := filepath.Dir(planFile)
	stateFile := filepath.Join(dir, ".plexr_state.json")

	// Check if state file exists
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		fmt.Println("❌ No execution state found. The plan has not been executed yet.")
		return nil
	}

	// Load state
	sm, err := core.NewStateManager(stateFile)
	if err != nil {
		return fmt.Errorf("failed to create state manager: %w", err)
	}

	state, err := sm.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Display status
	fmt.Printf("📋 Setup: %s (v%s)\n", state.SetupName, state.SetupVersion)
	fmt.Printf("🖥️  Platform: %s\n", state.Platform)
	fmt.Printf("🕒 Started: %s\n", state.StartedAt.Format(time.RFC3339))
	fmt.Printf("🔄 Last Updated: %s\n", state.UpdatedAt.Format(time.RFC3339))

	if state.CurrentStep != "" {
		fmt.Printf("⏳ Current Step: %s\n", state.CurrentStep)
	}

	if len(state.CompletedSteps) > 0 {
		fmt.Printf("\n✅ Completed Steps (%d):\n", len(state.CompletedSteps))
		for _, step := range state.CompletedSteps {
			fmt.Printf("   - %s\n", step)
		}
	}

	if len(state.FailedFiles) > 0 {
		fmt.Printf("\n❌ Failed Files:\n")
		for _, file := range state.FailedFiles {
			fmt.Printf("   - %s\n", file)
		}
	}

	if len(state.InstalledTools) > 0 {
		fmt.Printf("\n🛠️  Installed Tools:\n")
		for tool, version := range state.InstalledTools {
			fmt.Printf("   - %s: %s\n", tool, version)
		}
	}

	// Show raw JSON if verbose
	if verbose {
		fmt.Println("\n📄 Raw State:")
		data, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling state: %v\n", err)
		} else {
			fmt.Println(string(data))
		}
	}

	return nil
}
