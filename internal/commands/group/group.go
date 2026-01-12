// Package group provides advanced group management commands.
package group

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/common"
	"github.com/jtimothystewart/dtiam/internal/output"
	"github.com/jtimothystewart/dtiam/internal/resources"
)

// Cmd is the group command.
var Cmd = &cobra.Command{
	Use:   "group",
	Short: "Advanced group management commands",
	Long:  "Commands for advanced group operations like cloning and member management.",
}

func init() {
	Cmd.AddCommand(membersCmd)
	Cmd.AddCommand(addMemberCmd)
	Cmd.AddCommand(removeMemberCmd)
	Cmd.AddCommand(bindingsCmd)
}

var membersCmd = &cobra.Command{
	Use:   "members IDENTIFIER",
	Short: "List members of a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewGroupHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		group, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("group %q not found", args[0])
		}

		uuid, _ := group["uuid"].(string)
		members, err := handler.GetMembers(ctx, uuid)
		if err != nil {
			return err
		}

		return printer.Print(members, output.UserColumns())
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add-member IDENTIFIER",
	Short: "Add a user to a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		if email == "" {
			return fmt.Errorf("--email is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would add user %s to group %s\n", email, args[0])
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

		group, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("group %q not found", args[0])
		}

		uuid, _ := group["uuid"].(string)
		if err := handler.AddMember(ctx, uuid, email); err != nil {
			return err
		}

		printer.PrintSuccess("User %s added to group", email)
		return nil
	},
}

func init() {
	addMemberCmd.Flags().StringP("email", "e", "", "User email to add (required)")
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove-member IDENTIFIER",
	Short: "Remove a user from a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user")
		if userID == "" {
			return fmt.Errorf("--user is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would remove user %s from group %s\n", userID, args[0])
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

		group, err := resources.GetOrResolve(ctx, handler, args[0])
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("group %q not found", args[0])
		}

		uuid, _ := group["uuid"].(string)
		if err := handler.RemoveMember(ctx, uuid, userID); err != nil {
			return err
		}

		printer.PrintSuccess("User removed from group")
		return nil
	},
}

func init() {
	removeMemberCmd.Flags().StringP("user", "u", "", "User UID to remove (required)")
}

var bindingsCmd = &cobra.Command{
	Use:   "bindings IDENTIFIER",
	Short: "List policy bindings for a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		groupHandler := resources.NewGroupHandler(c)
		bindingHandler := resources.NewBindingHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		group, err := resources.GetOrResolve(ctx, groupHandler, args[0])
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("group %q not found", args[0])
		}

		uuid, _ := group["uuid"].(string)
		bindings, err := bindingHandler.GetForGroup(ctx, uuid)
		if err != nil {
			return err
		}

		return printer.Print(bindings, output.BindingColumns())
	},
}
