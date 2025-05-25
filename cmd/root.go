/*
Copyright © 2025 Plexr Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags
	verbose bool

	// Version information (set by build)
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plexr",
	Short: "A developer-friendly CLI tool for automating local development environment setup",
	Long: `Plexr helps developers set up and maintain their local development
environments through simple YAML configuration files.

No more "works on my machine" issues or spending hours following
outdated setup documentation.`,
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   false,
		DisableNoDescFlag:   false,
		DisableDescriptions: false,
		HiddenDefaultCmd:    false,
	},
	Run: func(cmd *cobra.Command, args []string) {
		printLogo()
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

// printLogo prints the ASCII art logo
func printLogo() {
	fmt.Print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf(`   ____  __    ____  _  _  ____ ` + "\n")
	fmt.Printf(`  (  _ \(  )  (  __)( \/ )(  _ \` + "\n")
	fmt.Printf(`   ) __// (_/\ ) _)  )  (  )   /` + "\n")
	fmt.Printf(`  (___) \____/(____)(_/\_)(__\_)` + "\n")
	fmt.Println()
	fmt.Printf("  Plan + Execute v%s \n", Version)
	fmt.Println()
	fmt.Print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// IsVerbose returns true if verbose flag is set
func IsVerbose() bool {
	return verbose
}
