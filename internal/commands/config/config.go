// Package config provides configuration management commands.
package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/config"
	"github.com/jtimothystewart/dtiam/internal/output"
)

// Cmd is the config command.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage dtiam configuration",
	Long:  "Commands for managing dtiam contexts and credentials.",
}

func init() {
	Cmd.AddCommand(viewCmd)
	Cmd.AddCommand(pathCmd)
	Cmd.AddCommand(getContextsCmd)
	Cmd.AddCommand(currentContextCmd)
	Cmd.AddCommand(useContextCmd)
	Cmd.AddCommand(setContextCmd)
	Cmd.AddCommand(deleteContextCmd)
	Cmd.AddCommand(setCredentialsCmd)
	Cmd.AddCommand(deleteCredentialsCmd)
	Cmd.AddCommand(getCredentialsCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display the current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Mask secrets unless --show-secrets is provided
		showSecrets, _ := cmd.Flags().GetBool("show-secrets")
		if !showSecrets {
			for i := range cfg.Credentials {
				cfg.Credentials[i].Credential.ClientSecret = config.MaskSecret(cfg.Credentials[i].Credential.ClientSecret)
			}
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		fmt.Print(string(data))
		return nil
	},
}

func init() {
	viewCmd.Flags().Bool("show-secrets", false, "Show unmasked secrets")
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Display the configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.GetConfigPath()
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

var getContextsCmd = &cobra.Command{
	Use:     "get-contexts",
	Aliases: []string{"contexts"},
	Short:   "List all contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		printer := cli.GlobalState.NewPrinter()

		// Build data for output
		data := make([]map[string]any, len(cfg.Contexts))
		for i, ctx := range cfg.Contexts {
			current := ""
			if ctx.Name == cfg.CurrentContext {
				current = "*"
			}
			data[i] = map[string]any{
				"name":            ctx.Name,
				"account_uuid":    ctx.Context.AccountUUID,
				"credentials_ref": ctx.Context.CredentialsRef,
				"current":         current,
			}
		}

		return printer.Print(data, output.ContextColumns())
	},
}

var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display the current context name",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.CurrentContext == "" {
			fmt.Println("No current context set")
		} else {
			fmt.Println(cfg.CurrentContext)
		}
		return nil
	},
}

var useContextCmd = &cobra.Command{
	Use:   "use-context NAME",
	Short: "Set the current context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.UseContext(name); err != nil {
			return err
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Switched to context %q\n", name)
		return nil
	},
}

var setContextCmd = &cobra.Command{
	Use:   "set-context NAME",
	Short: "Set or create a context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		accountUUID, _ := cmd.Flags().GetString("account-uuid")
		credentialsRef, _ := cmd.Flags().GetString("credentials-ref")

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var accountPtr, credPtr *string
		if accountUUID != "" {
			accountPtr = &accountUUID
		}
		if credentialsRef != "" {
			credPtr = &credentialsRef
		}

		if err := cfg.SetContext(name, accountPtr, credPtr); err != nil {
			return err
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Context %q updated\n", name)
		return nil
	},
}

func init() {
	setContextCmd.Flags().String("account-uuid", "", "Account UUID")
	setContextCmd.Flags().String("credentials-ref", "", "Credentials reference name")
}

var deleteContextCmd = &cobra.Command{
	Use:   "delete-context NAME",
	Short: "Delete a context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if !cfg.DeleteContext(name) {
			return fmt.Errorf("context %q not found", name)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Context %q deleted\n", name)
		return nil
	},
}

var setCredentialsCmd = &cobra.Command{
	Use:   "set-credentials NAME",
	Short: "Set or create credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")

		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("both --client-id and --client-secret are required")
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg.SetCredential(name, clientID, clientSecret)

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Credentials %q updated\n", name)
		return nil
	},
}

func init() {
	setCredentialsCmd.Flags().String("client-id", "", "OAuth2 client ID")
	setCredentialsCmd.Flags().String("client-secret", "", "OAuth2 client secret")
}

var deleteCredentialsCmd = &cobra.Command{
	Use:   "delete-credentials NAME",
	Short: "Delete credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if !cfg.DeleteCredential(name) {
			return fmt.Errorf("credentials %q not found", name)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Credentials %q deleted\n", name)
		return nil
	},
}

var getCredentialsCmd = &cobra.Command{
	Use:     "get-credentials",
	Aliases: []string{"credentials"},
	Short:   "List all credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		printer := cli.GlobalState.NewPrinter()

		data := make([]map[string]any, len(cfg.Credentials))
		for i, cred := range cfg.Credentials {
			data[i] = map[string]any{
				"name":      cred.Name,
				"client_id": cred.Credential.ClientID,
			}
		}

		return printer.Print(data, output.CredentialColumns())
	},
}
