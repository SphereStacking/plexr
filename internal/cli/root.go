package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	version = "dev"
)

// NewRootCommand creates the root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "plexr",
		Short: "A developer-friendly CLI tool for automating local development environment setup",
		Long: `Plexr helps developers set up and maintain their local development
environments through simple YAML configuration files.

No more "works on my machine" issues or spending hours following
outdated setup documentation.`,
		SilenceUsage: true,
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add subcommands
	rootCmd.AddCommand(NewExecuteCommand())
	rootCmd.AddCommand(NewValidateCommand())
	rootCmd.AddCommand(NewStatusCommand())
	rootCmd.AddCommand(NewResetCommand())
	rootCmd.AddCommand(NewVersionCommand())

	return rootCmd
}

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("plexr version %s\n", version)
		},
	}
}

// Execute runs the CLI
func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		// Error is already printed by cobra
		// Just exit with non-zero status
	}
}
