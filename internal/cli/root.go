package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	version = "0.0.1"
)

// printLogo prints the ASCII art logo
func printLogo() {
	fmt.Print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf(`   ____  __    ____  _  _  ____ ` + "\n")
	fmt.Printf(`  (  _ \(  )  (  __)( \/ )(  _ \` + "\n")
	fmt.Printf(`   ) __// (_/\ ) _)  )  (  )   /` + "\n")
	fmt.Printf(`  (___) \____/(____)(_/\_)(__\_)` + "\n")
	fmt.Println()
	fmt.Printf("  Plan + Execute v%s \n", version)
	fmt.Println()
	fmt.Print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

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
		Run: func(cmd *cobra.Command, args []string) {
			printLogo()
			_ = cmd.Help()
		},
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
			printLogo()
		},
	}
}

// Execute runs the CLI
func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		// Error is already printed by cobra
		os.Exit(1)
	}
}
