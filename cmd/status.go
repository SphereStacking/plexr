/*
Copyright © 2025 Plexr Authors
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/SphereStacking/plexr/internal/core"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <plan.yml>",
	Short: "Show execution status",
	Long: `Display the current execution status of a setup plan.

This command shows:
- Setup information (name, version, platform)
- Execution timeline (start time, last update)
- Current step being executed
- Completed steps
- Failed files (if any)
- Installed tools and versions`,
	Example: `  # Show status of a plan
  plexr status plan.yml

  # Show status with detailed JSON output
  plexr status plan.yml -v`,
	Args: cobra.ExactArgs(1),
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// Check if color output is disabled
func isColorDisabled() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("PLEXR_NO_COLOR") != ""
}

// Add color to string if color is enabled
func colorize(color, text string) string {
	if isColorDisabled() {
		return text
	}
	return color + text + colorReset
}

// Create a progress bar
func createProgressBar(completed, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(completed) / float64(total) * 100
	filled := int(float64(width) * float64(completed) / float64(total))

	bar := "["
	bar += strings.Repeat("=", filled)
	if filled < width {
		bar += ">"
		bar += strings.Repeat(" ", width-filled-1)
	}
	bar += "]"

	return fmt.Sprintf("%s %.1f%%", bar, percentage)
}

func runStatus(cmd *cobra.Command, args []string) error {
	planFile := args[0]
	dir := filepath.Dir(planFile)
	stateFile := filepath.Join(dir, ".plexr_state.json")

	// Check if state file exists
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		fmt.Println(colorize(colorRed, "❌ No execution state found. The plan has not been executed yet."))
		return nil
	}

	// Load the plan configuration
	plan, err := config.LoadExecutionPlan(planFile)
	if err != nil {
		return fmt.Errorf("failed to load plan: %w", err)
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

	// Display status header
	fmt.Printf("📋 Setup: %s %s\n", colorize(colorCyan, state.SetupName), colorize(colorGray, fmt.Sprintf("(v%s)", state.SetupVersion)))
	fmt.Printf("🖥️  Platform: %s\n", colorize(colorBlue, state.Platform))
	fmt.Printf("🕒 Started: %s\n", colorize(colorGray, state.StartedAt.Format(time.RFC3339)))
	fmt.Printf("🔄 Last Updated: %s\n", colorize(colorGray, state.UpdatedAt.Format(time.RFC3339)))

	// Calculate progress
	totalSteps := len(plan.Steps)
	completedCount := len(state.CompletedSteps)

	// Display progress bar
	fmt.Printf("\n📊 Progress: %s %s\n",
		createProgressBar(completedCount, totalSteps, 30),
		colorize(colorCyan, fmt.Sprintf("(%d/%d)", completedCount, totalSteps)))

	// Create a map of completed steps for quick lookup
	completedMap := make(map[string]bool)
	for _, stepID := range state.CompletedSteps {
		completedMap[stepID] = true
	}

	// Display steps with their status
	fmt.Println("\n📝 Steps:")
	for i, step := range plan.Steps {
		var statusIcon, statusColor string

		switch {
		case completedMap[step.ID]:
			statusIcon = "✅"
			statusColor = colorGreen
		case state.CurrentStep == step.ID:
			statusIcon = "⏳"
			statusColor = colorYellow
		default:
			statusIcon = "⏸️ "
			statusColor = colorGray
		}

		// Format step line
		stepLine := fmt.Sprintf("   %s %s", statusIcon, colorize(statusColor, step.ID))

		// Add description if available
		if step.Description != "" {
			stepLine += colorize(colorGray, fmt.Sprintf(" - %s", step.Description))
		}

		// Add dependencies if any
		if len(step.DependsOn) > 0 {
			deps := strings.Join(step.DependsOn, ", ")
			stepLine += colorize(colorGray, fmt.Sprintf(" [depends on: %s]", deps))
		}

		fmt.Printf("%3d. %s\n", i+1, stepLine)
	}

	// Display failed files if any
	if len(state.FailedFiles) > 0 {
		fmt.Printf("\n%s Failed Files:\n", colorize(colorRed, "❌"))
		for _, file := range state.FailedFiles {
			fmt.Printf("   - %s\n", colorize(colorRed, file))
		}
	}

	// Display installed tools if any
	if len(state.InstalledTools) > 0 {
		fmt.Printf("\n🛠️  Installed Tools:\n")
		for tool, version := range state.InstalledTools {
			fmt.Printf("   - %s: %s\n", colorize(colorBlue, tool), colorize(colorGreen, version))
		}
	}

	// Show raw JSON if verbose
	if IsVerbose() {
		fmt.Println("\n📄 Raw State:")
		data, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling state: %v\n", err)
		} else {
			fmt.Println(colorize(colorGray, string(data)))
		}
	}

	return nil
}
