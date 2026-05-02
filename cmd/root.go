package cmd

import (
	"fmt"
	"os"

	"devkit/config"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	cfgFile string
	
	// Global config manager
	cfg *config.Manager
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devkit",
	Short: "DevOps toolkit for health checking, log analysis, and notifications",
	Long: `devkit is a production-ready command-line tool that provides essential DevOps utilities.

Features:
  - Health checking: Concurrent HTTP health checks with configurable timeouts
  - Continuous monitoring: Real-time URL monitoring with Slack alerting
  - Log analysis: Parse and summarize Apache/Nginx, JSON, and plain text logs
  - Slack notifications: Send formatted notifications to Slack channels
  - Flexible output: Support for table, JSON, and plain text formats

Configuration:
  devkit uses .devkit.yaml for configuration. The file is searched in:
    1. Current directory (./.devkit.yaml)
    2. Home directory (~/.devkit.yaml)
  
  You can also specify a custom config file with the --config flag.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration manager
		cfg = config.NewManager()
		
		// Load configuration
		var err error
		if cfgFile != "" {
			// Load from specified config file
			err = cfg.LoadFrom(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config file: %w", err)
			}
		} else {
			// Load from default locations
			err = cfg.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
		}
		
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.devkit.yaml or ~/.devkit.yaml)")
}

// GetConfig returns the global configuration manager
// This is used by subcommands to access configuration
func GetConfig() *config.Manager {
	return cfg
}
