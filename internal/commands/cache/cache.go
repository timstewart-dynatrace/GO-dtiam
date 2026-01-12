// Package cache provides cache management commands.
package cache

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd is the cache command.
var Cmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management commands",
	Long:  "Commands for managing the in-memory cache.",
}

func init() {
	Cmd.AddCommand(clearCmd)
	Cmd.AddCommand(statsCmd)
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		// In-memory cache is cleared on each CLI invocation
		// This is a no-op in the Go implementation
		fmt.Println("Cache cleared")
		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show cache statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		// In-memory cache stats
		// The Go implementation doesn't maintain persistent cache stats
		fmt.Println("Cache is not persistent in this implementation.")
		fmt.Println("Each CLI invocation starts with a fresh cache.")
		return nil
	},
}
