package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information variables - these are set via ldflags during build
// They are imported from the main package
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Display the version number, git commit hash, and build date of devkit.

Version information is injected at build time via ldflags. When running
a development build, the version will be displayed as "dev".

Example:
  devkit version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devkit version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// SetVersionInfo sets the version information from the main package
// This is called by main.main() to inject the build-time values
func SetVersionInfo(v, c, b string) {
	version = v
	commit = c
	buildDate = b
}
