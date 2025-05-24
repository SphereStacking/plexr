package cli

import (
	"fmt"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/spf13/cobra"
)

// NewValidateCommand creates the validate command
func NewValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <plan.yml>",
		Short: "Validate an execution plan",
		Long:  `Validate a YAML execution plan for syntax errors and configuration issues.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runValidate,
	}
}

func runValidate(cmd *cobra.Command, args []string) error {
	planFile := args[0]

	fmt.Printf("Validating execution plan: %s\n", planFile)

	// Load and validate
	plan, err := config.LoadExecutionPlan(planFile)
	if err != nil {
		return fmt.Errorf("❌ Validation failed: %w", err)
	}

	fmt.Printf("\n✅ Execution plan is valid!\n")
	fmt.Printf("📋 Name: %s\n", plan.Name)
	fmt.Printf("📌 Version: %s\n", plan.Version)
	fmt.Printf("📊 Steps: %d\n", len(plan.Steps))
	fmt.Printf("🔧 Executors: %d\n", len(plan.Executors))

	// Check for warnings
	warnings := checkWarnings(plan)
	if len(warnings) > 0 {
		fmt.Println("\n⚠️  Warnings:")
		for _, warning := range warnings {
			fmt.Printf("   - %s\n", warning)
		}
	}

	return nil
}

func checkWarnings(plan *config.ExecutionPlan) []string {
	var warnings []string

	// Check for unused executors
	usedExecutors := make(map[string]bool)
	for _, step := range plan.Steps {
		usedExecutors[step.Executor] = true
	}

	for name := range plan.Executors {
		if !usedExecutors[name] {
			warnings = append(warnings, fmt.Sprintf("Executor '%s' is defined but not used", name))
		}
	}

	// Check for steps with no description
	for _, step := range plan.Steps {
		if step.Description == "" {
			warnings = append(warnings, fmt.Sprintf("Step '%s' has no description", step.ID))
		}
	}

	return warnings
}
