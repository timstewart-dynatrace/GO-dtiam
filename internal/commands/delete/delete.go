// Package delete provides commands for deleting resources.
package delete

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/common"
	"github.com/jtimothystewart/dtiam/internal/resources"
)

// Cmd is the delete command.
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource",
	Long:  "Commands for deleting IAM resources.",
}

func init() {
	Cmd.AddCommand(groupCmd)
	Cmd.AddCommand(policyCmd)
	Cmd.AddCommand(bindingCmd)
	Cmd.AddCommand(boundaryCmd)
	Cmd.AddCommand(userCmd)
	Cmd.AddCommand(serviceUserCmd)
}

// confirm asks for user confirmation.
func confirm(message string, force bool) bool {
	if force {
		return true
	}

	fmt.Printf("%s [y/N]: ", message)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}

var groupCmd = &cobra.Command{
	Use:   "group IDENTIFIER",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete group: %s\n", args[0])
			return nil
		}

		if !confirm(fmt.Sprintf("Delete group %q?", args[0]), force) {
			fmt.Println("Aborted.")
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

		// Resolve identifier
		group, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("group %q not found", args[0])
		}

		uuid, _ := group["uuid"].(string)
		if err := handler.Delete(ctx, uuid); err != nil {
			return err
		}

		printer.PrintSuccess("Group %q deleted successfully", args[0])
		return nil
	},
}

func init() {
	groupCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

var policyCmd = &cobra.Command{
	Use:   "policy IDENTIFIER",
	Short: "Delete a policy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete policy: %s\n", args[0])
			return nil
		}

		if !confirm(fmt.Sprintf("Delete policy %q?", args[0]), force) {
			fmt.Println("Aborted.")
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

		policy, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if policy == nil {
			return fmt.Errorf("policy %q not found", args[0])
		}

		uuid, _ := policy["uuid"].(string)
		if err := handler.Delete(ctx, uuid); err != nil {
			return err
		}

		printer.PrintSuccess("Policy %q deleted successfully", args[0])
		return nil
	},
}

func init() {
	policyCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

var bindingCmd = &cobra.Command{
	Use:   "binding",
	Short: "Delete a policy binding",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		policyID, _ := cmd.Flags().GetString("policy")
		force, _ := cmd.Flags().GetBool("force")

		if groupID == "" || policyID == "" {
			return fmt.Errorf("both --group and --policy are required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete binding: group=%s policy=%s\n", groupID, policyID)
			return nil
		}

		if !confirm(fmt.Sprintf("Delete binding for group %q and policy %q?", groupID, policyID), force) {
			fmt.Println("Aborted.")
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

		if err := handler.Delete(ctx, groupID, policyID); err != nil {
			return err
		}

		printer.PrintSuccess("Binding deleted successfully")
		return nil
	},
}

func init() {
	bindingCmd.Flags().StringP("group", "g", "", "Group UUID (required)")
	bindingCmd.Flags().StringP("policy", "p", "", "Policy UUID (required)")
	bindingCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

var boundaryCmd = &cobra.Command{
	Use:   "boundary IDENTIFIER",
	Short: "Delete a boundary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete boundary: %s\n", args[0])
			return nil
		}

		if !confirm(fmt.Sprintf("Delete boundary %q?", args[0]), force) {
			fmt.Println("Aborted.")
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

		boundary, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if boundary == nil {
			return fmt.Errorf("boundary %q not found", args[0])
		}

		uuid, _ := boundary["uuid"].(string)
		if err := handler.Delete(ctx, uuid); err != nil {
			return err
		}

		printer.PrintSuccess("Boundary %q deleted successfully", args[0])
		return nil
	},
}

func init() {
	boundaryCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

var userCmd = &cobra.Command{
	Use:   "user IDENTIFIER",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete user: %s\n", args[0])
			return nil
		}

		if !confirm(fmt.Sprintf("Delete user %q?", args[0]), force) {
			fmt.Println("Aborted.")
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewUserHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		user, err := handler.Get(ctx, args[0])
		if err != nil {
			user, err = handler.GetByEmail(ctx, args[0])
			if err != nil {
				return err
			}
		}
		if user == nil {
			return fmt.Errorf("user %q not found", args[0])
		}

		uid, _ := user["uid"].(string)
		if err := handler.Delete(ctx, uid); err != nil {
			return err
		}

		printer.PrintSuccess("User %q deleted successfully", args[0])
		return nil
	},
}

func init() {
	userCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

var serviceUserCmd = &cobra.Command{
	Use:     "service-user IDENTIFIER",
	Aliases: []string{"serviceuser"},
	Short:   "Delete a service user",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete service user: %s\n", args[0])
			return nil
		}

		if !confirm(fmt.Sprintf("Delete service user %q?", args[0]), force) {
			fmt.Println("Aborted.")
			return nil
		}

		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewServiceUserHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		user, err := handler.Get(ctx, args[0])
		if err != nil {
			user, err = handler.GetByName(ctx, args[0])
			if err != nil {
				return err
			}
		}
		if user == nil {
			return fmt.Errorf("service user %q not found", args[0])
		}

		uid, _ := user["uid"].(string)
		if err := handler.Delete(ctx, uid); err != nil {
			return err
		}

		printer.PrintSuccess("Service user %q deleted successfully", args[0])
		return nil
	},
}

func init() {
	serviceUserCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}
