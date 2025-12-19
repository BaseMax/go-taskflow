package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "taskflow",
	Short: "TaskFlow - A declarative task automation tool",
	Long: `TaskFlow is a powerful automation engine that executes workflows defined in YAML.
It supports shell commands, HTTP requests, file operations, task dependencies,
parallel execution, retries, and conditional logic.`,
	Version: version,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`TaskFlow {{.Version}}
A declarative task automation tool for developers and ops teams.
`)
}
