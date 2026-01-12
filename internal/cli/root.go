package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jtimothystewart/dtiam/internal/output"
	"github.com/jtimothystewart/dtiam/pkg/version"
)

var (
	// Flags
	contextFlag string
	outputFlag  string
	verboseFlag bool
	plainFlag   bool
	dryRunFlag  bool
)

// RootCmd is the root command for dtiam.
var RootCmd = &cobra.Command{
	Use:   "dtiam",
	Short: "A kubectl-inspired CLI for Dynatrace IAM",
	Long: `dtiam is a command-line tool for managing Dynatrace Identity and Access Management (IAM) resources.

It provides a consistent interface for managing groups, users, policies, bindings,
boundaries, environments, and service users.

DISCLAIMER: This tool is provided "as-is" without warranty. Use at your own risk.
This is an independent, community-developed tool and is NOT produced, endorsed,
or supported by Dynatrace.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Update global state from flags
		GlobalState.Context = contextFlag
		GlobalState.Verbose = verboseFlag
		GlobalState.Plain = plainFlag
		GlobalState.DryRun = dryRunFlag

		if outputFlag != "" {
			format, err := output.ParseFormat(outputFlag)
			if err != nil {
				return err
			}
			GlobalState.Output = format
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVar(&contextFlag, "context", "", "Override current context")
	RootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "Output format: table, wide, json, yaml, csv, plain")
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().BoolVar(&plainFlag, "plain", false, "Disable colors and interactive features")
	RootCmd.PersistentFlags().BoolVar(&dryRunFlag, "dry-run", false, "Preview changes without applying them")

	// Add version command
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dtiam version %s\n", version.Version)
		if version.Commit != "unknown" {
			fmt.Printf("  commit: %s\n", version.Commit)
		}
		if version.Date != "unknown" {
			fmt.Printf("  built:  %s\n", version.Date)
		}
	},
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// AddCommand adds a command to the root command.
func AddCommand(cmd *cobra.Command) {
	RootCmd.AddCommand(cmd)
}
