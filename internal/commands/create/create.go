// Package create provides commands for creating resources.
package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/common"
	"github.com/jtimothystewart/dtiam/internal/output"
	"github.com/jtimothystewart/dtiam/internal/resources"
)

// Cmd is the create command.
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
	Long:  "Commands for creating IAM resources.",
}

func init() {
	Cmd.AddCommand(groupCmd)
	Cmd.AddCommand(policyCmd)
	Cmd.AddCommand(bindingCmd)
	Cmd.AddCommand(boundaryCmd)
}

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Create a new group",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would create group: %s\n", name)
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewGroupHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		data := map[string]any{
			"name": name,
		}
		if description != "" {
			data["description"] = description
		}

		group, err := handler.Create(ctx, data)
		if err != nil {
			return err
		}

		printer.PrintSuccess("Group created successfully")
		return printer.Print([]map[string]any{group}, output.GroupColumns())
	},
}

func init() {
	groupCmd.Flags().StringP("name", "n", "", "Group name (required)")
	groupCmd.Flags().StringP("description", "d", "", "Group description")
}

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Create a new policy",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		statement, _ := cmd.Flags().GetString("statement")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if statement == "" {
			return fmt.Errorf("--statement is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would create policy: %s\n", name)
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewPolicyHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		data := map[string]any{
			"name":           name,
			"statementQuery": statement,
		}
		if description != "" {
			data["description"] = description
		}

		policy, err := handler.Create(ctx, data)
		if err != nil {
			return err
		}

		printer.PrintSuccess("Policy created successfully")
		return printer.Print([]map[string]any{policy}, output.PolicyColumns())
	},
}

func init() {
	policyCmd.Flags().StringP("name", "n", "", "Policy name (required)")
	policyCmd.Flags().StringP("statement", "s", "", "Policy statement query (required)")
	policyCmd.Flags().StringP("description", "d", "", "Policy description")
}

var bindingCmd = &cobra.Command{
	Use:   "binding",
	Short: "Create a new policy binding",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		policyID, _ := cmd.Flags().GetString("policy")
		boundaries, _ := cmd.Flags().GetStringSlice("boundary")

		if groupID == "" {
			return fmt.Errorf("--group is required")
		}
		if policyID == "" {
			return fmt.Errorf("--policy is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would create binding: group=%s policy=%s\n", groupID, policyID)
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

		binding, err := handler.Create(ctx, groupID, policyID, boundaries)
		if err != nil {
			return err
		}

		printer.PrintSuccess("Binding created successfully")
		return printer.Print([]map[string]any{binding}, output.BindingColumns())
	},
}

func init() {
	bindingCmd.Flags().StringP("group", "g", "", "Group UUID (required)")
	bindingCmd.Flags().StringP("policy", "p", "", "Policy UUID (required)")
	bindingCmd.Flags().StringSliceP("boundary", "b", nil, "Boundary UUIDs")
}

var boundaryCmd = &cobra.Command{
	Use:   "boundary",
	Short: "Create a new boundary",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		zones, _ := cmd.Flags().GetStringSlice("zone")
		query, _ := cmd.Flags().GetString("query")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if len(zones) == 0 && query == "" {
			return fmt.Errorf("either --zone or --query is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would create boundary: %s\n", name)
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewBoundaryHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		var queryPtr, descPtr *string
		if query != "" {
			queryPtr = &query
		}
		if description != "" {
			descPtr = &description
		}

		boundary, err := handler.Create(ctx, name, zones, queryPtr, descPtr)
		if err != nil {
			return err
		}

		printer.PrintSuccess("Boundary created successfully")
		return printer.Print([]map[string]any{boundary}, output.BoundaryColumns())
	},
}

func init() {
	boundaryCmd.Flags().StringP("name", "n", "", "Boundary name (required)")
	boundaryCmd.Flags().StringSliceP("zone", "z", nil, "Management zone names")
	boundaryCmd.Flags().StringP("query", "q", "", "Boundary query")
	boundaryCmd.Flags().StringP("description", "d", "", "Boundary description")
}
