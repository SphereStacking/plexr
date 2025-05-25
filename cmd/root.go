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

// printLogo prints the ASCII art logo with colors
func printLogo() {
	// ANSI color codes
	const (
		reset  = "\033[0m"
		bold   = "\033[1m"
		blue   = "\033[34m"
		cyan   = "\033[36m"
		green  = "\033[32m"
		yellow = "\033[33m"
		purple = "\033[35m"
	)

	// Check if we should disable colors (for non-TTY environments)
	noColor := os.Getenv("NO_COLOR") != "" || os.Getenv("PLEXR_NO_COLOR") != ""

	if noColor {
		// Print without colors
		logo := `
PPPPPPPPPPPPPPPPPP   lllllll                                                           
P::::::::::::::::P  l:::::l                                                           
P::::::PPPPPP:::::P l:::::l                                                           
PP:::::P     P:::::Pl:::::l                                                           
  P::::P     P:::::P l::::l     eeeeeeeeeeee  xxxxxxx      xxxxxxxrrrrr   rrrrrrrrr   
  P::::P     P:::::P l::::l   ee::::::::::::ee x:::::x    x:::::x r::::rrr:::::::::r  
  P::::PPPPPP:::::P  l::::l  e::::::eeeee:::::eex:::::x  x:::::x  r:::::::::::::::::r 
  P:::::::::::::PP   l::::l e::::::e     e:::::e x:::::xx:::::x   rr::::::rrrrr::::::r
  P::::PPPPPPPPP     l::::l e:::::::eeeee::::::e  x::::::::::x     r:::::r     r:::::r
  P::::P             l::::l e:::::::::::::::::e    x::::::::x      r:::::r     rrrrrrr
  P::::P             l::::l e::::::eeeeeeeeeee     x::::::::x      r:::::r            
  P::::P             l::::l e:::::::e             x::::::::::x     r:::::r            
PP::::::PP          l::::::le::::::::e           x:::::xx:::::x    r:::::r            
P::::::::P          l::::::l e::::::::eeeeeeee  x:::::x  x:::::x   r:::::r            
P::::::::P          l::::::l  ee:::::::::::::e x:::::x    x:::::x  r:::::r            
PPPPPPPPPP          llllllll    eeeeeeeeeeeeeexxxxxxx      xxxxxxx rrrrrrr            
`
		fmt.Print(logo)
		fmt.Printf("\n  Plan + Execute v%s \n", Version)
		fmt.Println()
		return
	}

	// Print with colors
	fmt.Println()
	fmt.Printf("%s%sPPPPPPPPPPPPPPPPPP   %slllllll                                                           %s\n", bold, blue, cyan, reset)
	fmt.Printf("%s%sP::::::::::::::::P  %sl:::::l                                                           %s\n", bold, blue, cyan, reset)
	fmt.Printf("%s%sP::::::PPPPPP:::::P %sl:::::l                                                           %s\n", bold, blue, cyan, reset)
	fmt.Printf("%s%sPP:::::P     P:::::P%sl:::::l                                                           %s\n", bold, blue, cyan, reset)
	fmt.Printf("%s%s  P::::P     P:::::P %sl::::l     %seeeeeeeeeeee  %sxxxxxxx      xxxxxxx%srrrrr   rrrrrrrrr   %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P::::P     P:::::P %sl::::l   %see::::::::::::ee %sx:::::x    x:::::x %sr::::rrr:::::::::r  %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P::::PPPPPP:::::P  %sl::::l  %se::::::eeeee:::::ee%sx:::::x  x:::::x  %sr:::::::::::::::::r %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P:::::::::::::PP   %sl::::l %se::::::e     e:::::e %sx:::::xx:::::x   %srr::::::rrrrr::::::r%s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P::::PPPPPPPPP     %sl::::l %se:::::::eeeee::::::e  %sx::::::::::x     %sr:::::r     r:::::r%s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P::::P             %sl::::l %se:::::::::::::::::e    %sx::::::::x      %sr:::::r     rrrrrrr%s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P::::P             %sl::::l %se::::::eeeeeeeeeee     %sx::::::::x      %sr:::::r            %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%s  P::::P             %sl::::l %se:::::::e             %sx::::::::::x     %sr:::::r            %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%sPP::::::PP          %sl::::::l%se::::::::e           %sx:::::xx:::::x    %sr:::::r            %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%sP::::::::P          %sl::::::l %se::::::::eeeeeeee  %sx:::::x  x:::::x   %sr:::::r            %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%sP::::::::P          %sl::::::l  %see:::::::::::::e %sx:::::x    x:::::x  %sr:::::r            %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Printf("%s%sPPPPPPPPPP          %sllllllll    %seeeeeeeeeeeeee%sxxxxxxx      xxxxxxx %srrrrrrr            %s\n", bold, blue, cyan, green, yellow, purple, reset)
	fmt.Println()
	fmt.Printf("  %sPlan + Execute %sv%s%s \n", bold, cyan, Version, reset)
	fmt.Println()
}

// IsVerbose returns true if verbose flag is set
func IsVerbose() bool {
	return verbose
}
