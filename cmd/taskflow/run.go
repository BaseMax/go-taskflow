package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BaseMax/go-taskflow/pkg/parser"
	"github.com/BaseMax/go-taskflow/pkg/workflow"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [workflow-file]",
	Short: "Run a workflow from a YAML file",
	Long:  `Execute a workflow defined in a YAML file with support for dependencies, retries, and parallel execution.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workflowFile := args[0]

		// Parse workflow
		fmt.Printf("📋 Loading workflow from: %s\n", workflowFile)
		wf, err := parser.ParseWorkflowFile(workflowFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error parsing workflow: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("📝 Workflow: %s\n", wf.Name)
		if wf.Description != "" {
			fmt.Printf("   %s\n", wf.Description)
		}
		fmt.Printf("   Tasks: %d\n\n", len(wf.Tasks))

		// Create engine and run workflow
		engine := workflow.NewEngine(wf)
		ctx := context.Background()

		startTime := time.Now()
		fmt.Println("🚀 Starting workflow execution...")
		fmt.Println()

		results, err := engine.Run(ctx)
		
		elapsed := time.Since(startTime)
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════")
		fmt.Println("📊 Execution Summary")
		fmt.Println("═══════════════════════════════════════════")

		successCount := 0
		failureCount := 0

		for _, result := range results {
			duration := result.EndTime.Sub(result.StartTime)
			if result.Success {
				fmt.Printf("✅ %s (%.2fs)\n", result.TaskName, duration.Seconds())
				if result.Output != "" && len(result.Output) < 200 {
					fmt.Printf("   Output: %s\n", result.Output)
				}
				successCount++
			} else {
				fmt.Printf("❌ %s (%.2fs)\n", result.TaskName, duration.Seconds())
				if result.Error != nil {
					fmt.Printf("   Error: %v\n", result.Error)
				}
				failureCount++
			}
		}

		fmt.Println("═══════════════════════════════════════════")
		fmt.Printf("Total: %d tasks | ✅ Success: %d | ❌ Failed: %d\n", len(results), successCount, failureCount)
		fmt.Printf("Total time: %.2fs\n", elapsed.Seconds())
		fmt.Println("═══════════════════════════════════════════")

		if err != nil {
			fmt.Fprintf(os.Stderr, "\n⚠️  Workflow completed with errors: %v\n", err)
			os.Exit(1)
		}

		if failureCount > 0 {
			os.Exit(1)
		}

		fmt.Println("\n✨ Workflow completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
