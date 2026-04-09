// Package template provides template management commands.
package template

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/jtimothystewart/dtiam/internal/cli"
	"github.com/jtimothystewart/dtiam/internal/commands/common"
	"github.com/jtimothystewart/dtiam/internal/output"
	"github.com/jtimothystewart/dtiam/internal/prompt"
	"github.com/jtimothystewart/dtiam/internal/resources"
	tmpl "github.com/jtimothystewart/dtiam/internal/template"
	"github.com/jtimothystewart/dtiam/internal/utils"
)

// Cmd is the template command.
var Cmd = &cobra.Command{
	Use:   "template",
	Short: "Manage and use resource templates",
	Long: `Commands for managing IAM resource templates.

Templates are reusable YAML definitions with variable placeholders that can be
rendered and applied to create resources. Built-in templates are provided for
common patterns. Custom templates can be saved and managed.`,
	Example: `  # List all available templates
  dtiam template list

  # Show a template's content and required variables
  dtiam template show policy-readonly

  # Render a template with variables
  dtiam template render policy-readonly --set name=MyPolicy

  # Render and create the resource
  dtiam template apply policy-readonly --set name=MyPolicy

  # Save a custom template
  dtiam template save my-template --file template.yaml

  # Show templates directory
  dtiam template path`,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(renderCmd)
	Cmd.AddCommand(applyCmd)
	Cmd.AddCommand(saveCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(pathCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available templates including built-in and custom templates.`,
	Example: `  # List all templates
  dtiam template list

  # Output as JSON
  dtiam template list -o json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer := cli.GlobalState.NewPrinter()

		var all []map[string]any

		for _, t := range tmpl.ListBuiltin() {
			all = append(all, map[string]any{
				"name":      t.Name,
				"source":    t.Source,
				"variables": strings.Join(t.Vars, ", "),
			})
		}

		store, err := tmpl.NewStore()
		if err != nil {
			return err
		}

		custom, err := store.List()
		if err != nil {
			return err
		}
		for _, t := range custom {
			all = append(all, map[string]any{
				"name":      t.Name,
				"source":    t.Source,
				"variables": strings.Join(t.Vars, ", "),
			})
		}

		if len(all) == 0 {
			fmt.Println("No templates found.")
			return nil
		}

		return printer.Print(all, TemplateColumns())
	},
}

var showCmd = &cobra.Command{
	Use:   "show NAME",
	Short: "Show template content and required variables",
	Long:  `Display the content of a template and list its required variables.`,
	Example: `  # Show a built-in template
  dtiam template show policy-readonly

  # Show a custom template
  dtiam template show my-template`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		content, source, err := getTemplate(name)
		if err != nil {
			return err
		}

		vars := tmpl.ExtractVariables(string(content))

		fmt.Printf("Template: %s (%s)\n", name, source)
		if len(vars) > 0 {
			fmt.Printf("Variables: %s\n", strings.Join(vars, ", "))
		}
		fmt.Println("---")
		fmt.Print(string(content))
		return nil
	},
}

var renderCmd = &cobra.Command{
	Use:   "render NAME",
	Short: "Render a template with variables",
	Long:  `Render a template to stdout with the given variable substitutions.`,
	Example: `  # Render a policy template
  dtiam template render policy-readonly --set name=MyReadOnlyPolicy

  # Render with multiple variables
  dtiam template render group-team --set name=DevOps --set description="DevOps Team"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		setFlags, _ := cmd.Flags().GetStringSlice("set")

		vars, err := tmpl.ParseSetFlags(setFlags)
		if err != nil {
			return err
		}

		content, _, err := getTemplate(args[0])
		if err != nil {
			return err
		}

		rendered, err := tmpl.RenderTemplate(string(content), vars)
		if err != nil {
			return err
		}

		fmt.Print(rendered)
		return nil
	},
}

func init() {
	renderCmd.Flags().StringSlice("set", nil, "Set template variable as key=value (repeatable)")
}

var applyCmd = &cobra.Command{
	Use:   "apply NAME",
	Short: "Render a template and create the resource",
	Long: `Render a template with variables and create the resulting resource.

The template is rendered, the kind field is read, and the resource is created
via the appropriate API handler.`,
	Example: `  # Create a read-only policy from template
  dtiam template apply policy-readonly --set name=MyPolicy

  # Create a group from template
  dtiam template apply group-team --set name=DevOps --set description="DevOps Team"

  # Preview what would be created
  dtiam template apply policy-readonly --set name=MyPolicy --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		setFlags, _ := cmd.Flags().GetStringSlice("set")

		vars, err := tmpl.ParseSetFlags(setFlags)
		if err != nil {
			return err
		}

		content, _, err := getTemplate(args[0])
		if err != nil {
			return err
		}

		rendered, err := tmpl.RenderTemplate(string(content), vars)
		if err != nil {
			return err
		}

		// Parse the rendered YAML
		var doc map[string]any
		if err := yaml.Unmarshal([]byte(rendered), &doc); err != nil {
			return fmt.Errorf("failed to parse rendered template: %w", err)
		}

		kind, _ := doc["kind"].(string)
		spec, _ := doc["spec"].(map[string]any)
		if kind == "" {
			return fmt.Errorf("template missing 'kind' field")
		}
		if spec == nil {
			return fmt.Errorf("template missing 'spec' field")
		}

		printer := cli.GlobalState.NewPrinter()

		if cli.GlobalState.IsDryRun() {
			printer.PrintWarning("Would create %s from template %q:", kind, args[0])
			fmt.Fprint(os.Stderr, rendered)
			return nil
		}

		return createResource(kind, spec, printer)
	},
}

func init() {
	applyCmd.Flags().StringSlice("set", nil, "Set template variable as key=value (repeatable)")
}

var saveCmd = &cobra.Command{
	Use:   "save NAME",
	Short: "Save a custom template from a file",
	Long:  `Save a custom template to the templates directory for reuse.`,
	Example: `  # Save a template from a file
  dtiam template save my-policy --file my-policy-template.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		printer := cli.GlobalState.NewPrinter()

		if cli.GlobalState.IsDryRun() {
			printer.PrintWarning("Would save template %q from %s", args[0], file)
			return nil
		}

		store, err := tmpl.NewStore()
		if err != nil {
			return err
		}

		if err := store.Save(args[0], content); err != nil {
			return err
		}

		printer.PrintSuccess("Template %q saved to %s", args[0], store.Path())
		return nil
	},
}

func init() {
	saveCmd.Flags().StringP("file", "f", "", "Template file to save (required)")
	_ = saveCmd.MarkFlagRequired("file")
}

var deleteCmd = &cobra.Command{
	Use:   "delete NAME",
	Short: "Delete a custom template",
	Long:  `Delete a custom template from the templates directory.`,
	Example: `  # Delete a custom template
  dtiam template delete my-template

  # Delete without confirmation
  dtiam template delete my-template --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		printer := cli.GlobalState.NewPrinter()

		if cli.GlobalState.IsDryRun() {
			printer.PrintWarning("Would delete template %q", args[0])
			return nil
		}

		if !prompt.ConfirmDelete("template", args[0], force || cli.GlobalState.IsPlain()) {
			fmt.Println("Aborted.")
			return nil
		}

		store, err := tmpl.NewStore()
		if err != nil {
			return err
		}

		if err := store.Delete(args[0]); err != nil {
			return err
		}

		printer.PrintSuccess("Template %q deleted", args[0])
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the templates directory path",
	Long:  `Show the filesystem path where custom templates are stored.`,
	Example: `  # Show templates path
  dtiam template path`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := tmpl.NewStore()
		if err != nil {
			return err
		}
		fmt.Println(store.Path())
		return nil
	},
}

// getTemplate returns template content and source, checking custom then builtin.
func getTemplate(name string) ([]byte, string, error) {
	store, err := tmpl.NewStore()
	if err != nil {
		return nil, "", err
	}

	if content, err := store.Get(name); err == nil {
		return content, "custom", nil
	}

	content, err := tmpl.GetBuiltin(name)
	if err != nil {
		return nil, "", fmt.Errorf("template %q not found (checked custom and built-in)", name)
	}
	return content, "builtin", nil
}

// createResource creates a resource based on kind and spec.
func createResource(kind string, spec map[string]any, printer interface {
	PrintSuccess(string, ...any)
	PrintDetail(map[string]any) error
}) error {
	c, err := common.CreateClient()
	if err != nil {
		return err
	}
	defer c.Close()

	ctx := context.Background()

	switch strings.ToLower(kind) {
	case "group":
		handler := resources.NewGroupHandler(c)
		result, err := handler.Create(ctx, spec)
		if err != nil {
			return err
		}
		printer.PrintSuccess("Group created from template")
		return printer.PrintDetail(result)

	case "policy":
		handler := resources.NewPolicyHandler(c)
		result, err := handler.Create(ctx, spec)
		if err != nil {
			return err
		}
		printer.PrintSuccess("Policy created from template")
		return printer.PrintDetail(result)

	case "boundary":
		handler := resources.NewBoundaryHandler(c)
		name := utils.StringFrom(spec, "name")
		var zones []string
		if z, ok := spec["zones"].([]any); ok {
			for _, zone := range z {
				if s, ok := zone.(string); ok {
					zones = append(zones, s)
				}
			}
		}
		var query *string
		if q, ok := spec["boundaryQuery"].(string); ok {
			query = &q
		}
		var desc *string
		if d, ok := spec["description"].(string); ok {
			desc = &d
		}
		result, err := handler.Create(ctx, name, zones, query, desc)
		if err != nil {
			return err
		}
		printer.PrintSuccess("Boundary created from template")
		return printer.PrintDetail(result)

	case "binding":
		handler := resources.NewBindingHandler(c)
		groupUUID := utils.StringFrom(spec, "group")
		policyUUID := utils.StringFrom(spec, "policy")
		if groupUUID == "" || policyUUID == "" {
			return fmt.Errorf("binding spec requires 'group' and 'policy' fields")
		}
		result, err := handler.Create(ctx, groupUUID, policyUUID, nil, nil)
		if err != nil {
			return err
		}
		printer.PrintSuccess("Binding created from template")
		return printer.PrintDetail(result)

	default:
		return fmt.Errorf("unsupported resource kind: %s", kind)
	}
}

// TemplateColumns returns columns for template list output.
func TemplateColumns() []output.Column {
	return []output.Column{
		{Key: "name", Header: "NAME"},
		{Key: "source", Header: "SOURCE"},
		{Key: "variables", Header: "VARIABLES"},
	}
}
