// Package serviceuser provides service user (OAuth client) management commands.
package serviceuser

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/common"
	"github.com/jtimothystewart/dtiam/internal/output"
	"github.com/jtimothystewart/dtiam/internal/resources"
)

// Cmd is the service-user command.
var Cmd = &cobra.Command{
	Use:     "service-user",
	Aliases: []string{"serviceuser"},
	Short:   "Service user (OAuth client) management commands",
	Long:    "Commands for managing service users (OAuth clients).",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(addToGroupCmd)
	Cmd.AddCommand(removeFromGroupCmd)
	Cmd.AddCommand(listGroupsCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List service users",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := common.CreateClient()
		if err != nil {
			return err
		}
		defer c.Close()

		handler := resources.NewServiceUserHandler(c)
		printer := cli.GlobalState.NewPrinter()
		ctx := context.Background()

		users, err := handler.List(ctx, nil)
		if err != nil {
			return err
		}

		return printer.Print(users, output.ServiceUserColumns())
	},
}

var getCmd = &cobra.Command{
	Use:   "get IDENTIFIER",
	Short: "Get a service user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		return printer.PrintSingle(user, output.ServiceUserColumns())
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a service user",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		groupsStr, _ := cmd.Flags().GetString("groups")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would create service user: %s\n", name)
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

		var descPtr *string
		if description != "" {
			descPtr = &description
		}

		var groups []string
		if groupsStr != "" {
			groups = strings.Split(groupsStr, ",")
			for i := range groups {
				groups[i] = strings.TrimSpace(groups[i])
			}
		}

		user, err := handler.Create(ctx, name, descPtr, groups)
		if err != nil {
			return err
		}

		printer.PrintSuccess("Service user created successfully")
		printer.PrintWarning("Save the credentials below - they cannot be retrieved later!")
		return printer.PrintDetail(user)
	},
}

func init() {
	createCmd.Flags().StringP("name", "n", "", "Service user name (required)")
	createCmd.Flags().StringP("description", "d", "", "Service user description")
	createCmd.Flags().StringP("groups", "g", "", "Comma-separated list of group UUIDs")
}

var updateCmd = &cobra.Command{
	Use:   "update IDENTIFIER",
	Short: "Update a service user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would update service user: %s\n", args[0])
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

		// Find the user first
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

		var namePtr, descPtr *string
		if name != "" {
			namePtr = &name
		}
		if description != "" {
			descPtr = &description
		}

		updated, err := handler.Update(ctx, uid, namePtr, descPtr, nil)
		if err != nil {
			return err
		}

		printer.PrintSuccess("Service user updated successfully")
		return printer.PrintSingle(updated, output.ServiceUserColumns())
	},
}

func init() {
	updateCmd.Flags().StringP("name", "n", "", "New name")
	updateCmd.Flags().StringP("description", "d", "", "New description")
}

var deleteCmd = &cobra.Command{
	Use:   "delete IDENTIFIER",
	Short: "Delete a service user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would delete service user: %s\n", args[0])
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

		printer.PrintSuccess("Service user deleted successfully")
		return nil
	},
}

var addToGroupCmd = &cobra.Command{
	Use:   "add-to-group IDENTIFIER",
	Short: "Add a service user to a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		if groupID == "" {
			return fmt.Errorf("--group is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would add service user %s to group %s\n", args[0], groupID)
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
		if err := handler.AddToGroup(ctx, uid, groupID); err != nil {
			return err
		}

		printer.PrintSuccess("Service user added to group")
		return nil
	},
}

func init() {
	addToGroupCmd.Flags().StringP("group", "g", "", "Group UUID (required)")
}

var removeFromGroupCmd = &cobra.Command{
	Use:   "remove-from-group IDENTIFIER",
	Short: "Remove a service user from a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		if groupID == "" {
			return fmt.Errorf("--group is required")
		}

		if cli.GlobalState.IsDryRun() {
			fmt.Printf("Would remove service user %s from group %s\n", args[0], groupID)
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
		if err := handler.RemoveFromGroup(ctx, uid, groupID); err != nil {
			return err
		}

		printer.PrintSuccess("Service user removed from group")
		return nil
	},
}

func init() {
	removeFromGroupCmd.Flags().StringP("group", "g", "", "Group UUID (required)")
}

var listGroupsCmd = &cobra.Command{
	Use:   "list-groups IDENTIFIER",
	Short: "List groups a service user belongs to",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
		groups, err := handler.GetGroups(ctx, uid)
		if err != nil {
			return err
		}

		return printer.Print(groups, output.GroupColumns())
	},
}
