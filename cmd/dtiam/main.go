// Package main is the entry point for the dtiam CLI.
package main

import (
	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/account"
	"github.com/jtimothystewart/dtiam/internal/commands/boundary"
	"github.com/jtimothystewart/dtiam/internal/commands/cache"
	configcmd "github.com/jtimothystewart/dtiam/internal/commands/config"
	"github.com/jtimothystewart/dtiam/internal/commands/create"
	deletecmd "github.com/jtimothystewart/dtiam/internal/commands/delete"
	"github.com/jtimothystewart/dtiam/internal/commands/describe"
	"github.com/jtimothystewart/dtiam/internal/commands/get"
	"github.com/jtimothystewart/dtiam/internal/commands/group"
	"github.com/jtimothystewart/dtiam/internal/commands/serviceuser"
	"github.com/jtimothystewart/dtiam/internal/commands/user"
)

func main() {
	// Register commands
	cli.AddCommand(configcmd.Cmd)
	cli.AddCommand(get.Cmd)
	cli.AddCommand(describe.Cmd)
	cli.AddCommand(create.Cmd)
	cli.AddCommand(deletecmd.Cmd)
	cli.AddCommand(user.Cmd)
	cli.AddCommand(serviceuser.Cmd)
	cli.AddCommand(group.Cmd)
	cli.AddCommand(boundary.Cmd)
	cli.AddCommand(account.Cmd)
	cli.AddCommand(cache.Cmd)

	// Execute
	cli.Execute()
}
