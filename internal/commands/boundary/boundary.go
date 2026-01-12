// Package boundary provides boundary management commands.
package boundary

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/common"
	"github.com/jtimothystewart/dtiam/internal/resources"
)

// Cmd is the boundary command.
var Cmd = &cobra.Command{
	Use:   "boundary",
	Short: "Boundary management commands",
	Long:  "Commands for attaching and detaching boundaries from policy bindings.",
}

func init() {
	Cmd.AddCommand(attachCmd)
	Cmd.AddCommand(detachCmd)
	Cmd.AddCommand(listAttachedCmd)
}

var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a boundary to a policy binding",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		policyID, _ := cmd.Flags().GetString("policy")
		boundaryID, _ := cmd.Flags().GetString("boundary")

		if groupID == "" || policyID == "" || boundaryID == "" {
			return fmt.Errorf("--group, --policy, and --boundary are all required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would attach boundary %s to binding (group=%s, policy=%s)\n", boundaryID, groupID, policyID)
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewBindingHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		if err := handler.AddBoundary(ctx, groupID, policyID, boundaryID); err != nil {
			return err
		}

		printer.PrintSuccess("Boundary attached successfully")
		return nil
	},
}

func init() {
	attachCmd.Flags().StringP("group", "g", "", "Group UUID (required)")
	attachCmd.Flags().StringP("policy", "p", "", "Policy UUID (required)")
	attachCmd.Flags().StringP("boundary", "b", "", "Boundary UUID (required)")
}

var detachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a boundary from a policy binding",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		policyID, _ := cmd.Flags().GetString("policy")
		boundaryID, _ := cmd.Flags().GetString("boundary")

		if groupID == "" || policyID == "" || boundaryID == "" {
			return fmt.Errorf("--group, --policy, and --boundary are all required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would detach boundary %s from binding (group=%s, policy=%s)\n", boundaryID, groupID, policyID)
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewBindingHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		if err := handler.RemoveBoundary(ctx, groupID, policyID, boundaryID); err != nil {
			return err
		}

		printer.PrintSuccess("Boundary detached successfully")
		return nil
	},
}

func init() {
	detachCmd.Flags().StringP("group", "g", "", "Group UUID (required)")
	detachCmd.Flags().StringP("policy", "p", "", "Policy UUID (required)")
	detachCmd.Flags().StringP("boundary", "b", "", "Boundary UUID (required)")
}

var listAttachedCmd = &cobra.Command{
	Use:   "list-attached IDENTIFIER",
	Short: "List policies using a boundary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewBoundaryHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		boundary, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if boundary == nil {
			return fmt.Errorf("boundary %q not found", args[0])
		}

		uuid, _ := boundary["uuid"].(string)
		attached, err := handler.GetAttachedPolicies(ctx, uuid)
		if err != nil {
			return err
		}

		if len(attached) == 0 {
			fmt.Println("No policies are using this boundary.")
			return nil
		}

		// Define columns for attached policies
		columns := []struct {
			Key    string
			Header string
		}{
			{"policyUuid", "POLICY_UUID"},
			{"groupUuid", "GROUP_UUID"},
		}

		// Convert to output.Column format
		var cols []struct {
			Key    string
			Header string
		}
		for _, c := range columns {
			cols = append(cols, c)
		}

		return printer.Print(attached, nil)
	},
}
