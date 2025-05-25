package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SphereStacking/plexr/internal/core"
	"github.com/spf13/cobra"
)

// NewResetCommand creates the reset command
func NewResetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset <plan.yml>",
		Short: "Reset execution state",
		Long:  `Reset the execution state of a setup plan, allowing it to be run from the beginning.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runReset,
	}

	cmd.Flags().BoolVarP(&auto, "auto", "a", false, "Skip confirmation prompt")

	return cmd
}

func runReset(cmd *cobra.Command, args []string) error {
	planFile := args[0]
	dir := filepath.Dir(planFile)
	stateFile := filepath.Join(dir, ".plexr_state.json")

	// Check if state file exists
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		fmt.Println("❌ No execution state found. Nothing to reset.")
		return nil
	}

	// Confirm reset
	if !auto {
		fmt.Print("⚠️  This will reset all execution state. Continue? [y/N]: ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			fmt.Printf("\nFailed to read input: %v\n", err)
			return err
		}
		if response != "y" && response != "Y" {
			fmt.Println("Reset canceled.")
			return nil
		}
	}

	// Reset state
	sm, err := core.NewStateManager(stateFile)
	if err != nil {
		return fmt.Errorf("failed to create state manager: %w", err)
	}

	err = sm.Reset()
	if err != nil {
		return fmt.Errorf("failed to reset state: %w", err)
	}

	fmt.Println("✅ Execution state has been reset.")
	return nil
}
