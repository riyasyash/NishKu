package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags during build)
var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version, build information, and release details for nishku.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		
		if verbose {
			// Detailed version information
			fmt.Printf("nishku version %s\n", Version)
			fmt.Printf("  Git Commit:  %s\n", GitCommit)
			fmt.Printf("  Build Date:  %s\n", BuildDate)
			fmt.Printf("  Go Version:  %s\n", GoVersion)
			fmt.Printf("  Platform:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		} else {
			// Simple version output
			fmt.Printf("nishku version %s\n", Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("verbose", "v", false, "Show detailed version information")
}
