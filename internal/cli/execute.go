package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/SphereStacking/plexr/internal/core"
	"github.com/SphereStacking/plexr/internal/utils"
	"github.com/spf13/cobra"
)

var (
	auto     bool
	dryRun   bool
	fromStep string
	platform string
	only     string
)

// NewExecuteCommand creates the execute command
func NewExecuteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "execute <plan.yml>",
		Aliases: []string{"exec", "run"},
		Short:   "Execute a setup plan",
		Long:    `Execute a YAML-based setup plan to configure your local development environment.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runExecute,
	}

	cmd.Flags().BoolVarP(&auto, "auto", "a", false, "Skip confirmation prompts")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be executed without running")
	cmd.Flags().StringVar(&fromStep, "from-step", "", "Start execution from a specific step")
	cmd.Flags().StringVarP(&platform, "platform", "p", "", "Override platform detection")
	cmd.Flags().StringVarP(&only, "only", "o", "", "Execute only specific steps")

	return cmd
}

func runExecute(cmd *cobra.Command, args []string) error {
	// Print logo
	printLogo()

	planFile := args[0]

	// Initialize logger
	err := utils.InitLogger(verbose)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Load execution plan
	fmt.Printf("Loading execution plan: %s\n", planFile)
	plan, err := config.LoadExecutionPlan(planFile)
	if err != nil {
		return fmt.Errorf("failed to load execution plan: %w", err)
	}

	fmt.Printf("\n📋 Execution Plan: %s (v%s)\n", plan.Name, plan.Version)
	if plan.Description != "" {
		fmt.Printf("📝 %s\n", plan.Description)
	}

	if dryRun {
		fmt.Println("\n🔍 DRY RUN MODE - No changes will be made")
		return showExecutionPlan(plan)
	}

	// Create state file path
	dir := filepath.Dir(planFile)
	stateFile := filepath.Join(dir, ".plexr_state.json")

	// Create runner
	runner, err := core.NewRunner(plan, stateFile)
	if err != nil {
		return fmt.Errorf("failed to create runner: %w", err)
	}

	// Confirm execution
	if !auto {
		fmt.Print("\n⚡ Ready to execute. Continue? [y/N]: ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			fmt.Printf("\nFailed to read input: %v\n", err)
			return err
		}
		if response != "y" && response != "Y" {
			fmt.Println("Execution canceled.")
			return nil
		}
	}

	// Execute
	ctx := context.Background()
	fmt.Println("\n🚀 Starting execution...")
	err = runner.Execute(ctx)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	fmt.Println("\n✅ Execution completed successfully!")
	return nil
}

func showExecutionPlan(plan *config.ExecutionPlan) error {
	fmt.Println("\n📋 Steps to be executed:")
	for i, step := range plan.Steps {
		fmt.Printf("\n%d. %s", i+1, step.ID)
		if step.Description != "" {
			fmt.Printf(" - %s", step.Description)
		}
		fmt.Println()

		if len(step.DependsOn) > 0 {
			fmt.Printf("   Dependencies: %v\n", step.DependsOn)
		}

		if step.SkipIf != "" {
			fmt.Printf("   Skip if: %s\n", step.SkipIf)
		}

		fmt.Printf("   Executor: %s\n", step.Executor)
		fmt.Printf("   Files:\n")
		for _, file := range step.Files {
			fmt.Printf("     - %s", file.Path)
			if file.Platform != "" {
				fmt.Printf(" (platform: %s)", file.Platform)
			}
			if file.Timeout > 0 {
				fmt.Printf(" (timeout: %ds)", file.Timeout)
			}
			fmt.Println()
		}
	}
	return nil
}
