# CLAUDE.md

This file provides guidance for AI agents working with the dtiam codebase.

> **DISCLAIMER:** This tool is provided "as-is" without warranty. Use at your own risk. This is an independent, community-developed tool and is **NOT produced, endorsed, or supported by Dynatrace**.

## Project Overview

**dtiam** is a kubectl-inspired CLI for managing Dynatrace Identity and Access Management (IAM) resources. It provides a consistent interface for managing groups, users, policies, bindings, boundaries, environments, and service users.

**Language:** Go 1.22+

## Quick Reference

### Build & Run

```bash
# Build the CLI
make build

# Install locally
make install

# Run CLI
./bin/dtiam --help
dtiam --help  # if installed

# Run with verbose output
dtiam -v get groups

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

### Project Structure

```
dtiam/
├── cmd/dtiam/main.go                 # Entry point
├── internal/
│   ├── cli/
│   │   ├── root.go                   # Root command, global flags
│   │   └── state.go                  # Global state (context, output, verbose)
│   ├── commands/
│   │   ├── common/                   # Shared command utilities
│   │   ├── config/                   # Config management commands
│   │   ├── get/                      # List/retrieve resources
│   │   ├── describe/                 # Detailed resource views
│   │   ├── create/                   # Create resources
│   │   ├── delete/                   # Delete resources
│   │   ├── user/                     # User lifecycle commands
│   │   ├── serviceuser/              # OAuth client management
│   │   ├── group/                    # Advanced group ops
│   │   ├── boundary/                 # Boundary attach/detach
│   │   ├── account/                  # Limits and subscriptions
│   │   └── cache/                    # Cache management
│   ├── config/
│   │   ├── config.go                 # Config structs
│   │   └── loader.go                 # Config load/save, XDG paths
│   ├── client/
│   │   ├── client.go                 # HTTP client with retry
│   │   └── errors.go                 # APIError type
│   ├── auth/
│   │   ├── auth.go                   # TokenProvider interface
│   │   ├── oauth.go                  # OAuth2 token manager
│   │   └── bearer.go                 # Static bearer token
│   ├── resources/
│   │   ├── handler.go                # Handler interfaces
│   │   ├── groups.go                 # GroupHandler
│   │   ├── users.go                  # UserHandler
│   │   ├── policies.go               # PolicyHandler
│   │   ├── bindings.go               # BindingHandler
│   │   ├── boundaries.go             # BoundaryHandler
│   │   ├── environments.go           # EnvironmentHandler
│   │   ├── serviceusers.go           # ServiceUserHandler
│   │   ├── limits.go                 # LimitsHandler
│   │   └── subscriptions.go          # SubscriptionHandler
│   ├── output/
│   │   ├── format.go                 # Format enum
│   │   ├── printer.go                # Unified Printer
│   │   ├── table.go                  # Table formatter
│   │   └── columns.go                # Column definitions
│   └── utils/
│       └── resolver.go               # Name-to-UUID resolution
├── pkg/version/version.go            # Version info
├── go.mod
├── Makefile
└── .goreleaser.yaml
```

## Authentication

dtiam supports two authentication methods:

### OAuth2 (Recommended)
- Auto-refreshes tokens when expired
- Best for automation, CI/CD, long-running processes
- Requires `DTIAM_CLIENT_ID`, `DTIAM_CLIENT_SECRET`, `DTIAM_ACCOUNT_UUID`

### Bearer Token (Static)
- Does NOT auto-refresh (fails when token expires)
- Best for quick testing, debugging, one-off operations
- Requires `DTIAM_BEARER_TOKEN`, `DTIAM_ACCOUNT_UUID`

### Environment Variables

| Variable | Description |
|----------|-------------|
| `DTIAM_BEARER_TOKEN` | Static bearer token (alternative to OAuth2) |
| `DTIAM_CLIENT_ID` | OAuth2 client ID |
| `DTIAM_CLIENT_SECRET` | OAuth2 client secret |
| `DTIAM_ACCOUNT_UUID` | Dynatrace account UUID |
| `DTIAM_CONTEXT` | Override current context |

## Key Patterns

### Adding a New Command

1. Create command file in `internal/commands/<name>/<name>.go`:
```go
package newfeature

import (
    "context"
    "github.com/spf13/cobra"
    "github.com/jtimothystewart/dtiam/internal/cli"
    "github.com/jtimothystewart/dtiam/internal/commands/common"
)

var Cmd = &cobra.Command{
    Use:   "new-feature",
    Short: "New feature operations",
}

func init() {
    Cmd.AddCommand(doSomethingCmd)
}

var doSomethingCmd = &cobra.Command{
    Use:   "do-something NAME",
    Short: "Do something useful",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        c, err := common.CreateClient()
        if err != nil {
            return err
        }
        defer c.Close()

        printer := cli.GlobalState.NewPrinter()
        ctx := context.Background()

        // Implementation here
        return nil
    },
}
```

2. Register in `cmd/dtiam/main.go`:
```go
import "github.com/jtimothystewart/dtiam/internal/commands/newfeature"

func main() {
    cli.AddCommand(newfeature.Cmd)
    // ...
}
```

### Adding a New Resource Handler

1. Create handler in `internal/resources/<name>.go`:
```go
package resources

import "github.com/jtimothystewart/dtiam/internal/client"

type NewResourceHandler struct {
    BaseHandler
}

func NewNewResourceHandler(c *client.Client) *NewResourceHandler {
    return &NewResourceHandler{
        BaseHandler: BaseHandler{
            Client:    c,
            Name:      "new-resource",
            Path:      "/new-resources",
            ListKey:   "items",
            IDField:   "uuid",
            NameField: "name",
        },
    }
}
```

2. Add columns in `internal/output/columns.go`:
```go
func NewResourceColumns() []Column {
    return []Column{
        {Key: "uuid", Header: "UUID"},
        {Key: "name", Header: "NAME"},
        {Key: "description", Header: "DESCRIPTION"},
    }
}
```

### Global State Access

Commands access global CLI state through the `cli` package:
```go
import "github.com/jtimothystewart/dtiam/internal/cli"

// Available:
cli.GlobalState.Context   // string - context override
cli.GlobalState.Output    // output.Format - output format
cli.GlobalState.Verbose   // bool - verbose mode
cli.GlobalState.Plain     // bool - plain mode (no colors)
cli.GlobalState.DryRun    // bool - dry-run mode

// Create printer with current settings
printer := cli.GlobalState.NewPrinter()
```

### HTTP Client Usage

Always close the client when done:
```go
c, err := common.CreateClient()
if err != nil {
    return err
}
defer c.Close()

ctx := context.Background()
body, err := c.Get(ctx, "/groups", nil)
```

### Output Formatting

Use the Printer for consistent output:
```go
printer := cli.GlobalState.NewPrinter()
printer.Print(data, output.GroupColumns())      // List
printer.PrintSingle(data, output.GroupColumns()) // Single item
printer.PrintDetail(data)                        // Key-value pairs
printer.PrintSuccess("Operation completed")      // Success message
printer.PrintWarning("Warning message")          // Warning message
```

## API Endpoints

Base URL: `https://api.dynatrace.com/iam/v1/accounts/{account_uuid}`

| Resource | Path |
|----------|------|
| Groups | `/groups` |
| Users | `/users` |
| Service Users | `/service-users` |
| Limits | `/limits` |
| Policies | `/repo/{level_type}/{level_id}/policies` |
| Bindings | `/repo/{level_type}/{level_id}/bindings` |
| Boundaries | `/repo/account/{uuid}/boundaries` |

**Environment API**: `https://api.dynatrace.com/env/v2/accounts/{uuid}/environments`

**Subscription API**: `https://api.dynatrace.com/sub/v2/accounts/{uuid}/subscriptions`

Level types: `account`, `environment`, `global`

## Configuration

Config file: `~/.config/dtiam/config` (YAML)

```yaml
api-version: v1
kind: Config
current-context: production
contexts:
  - name: production
    context:
      account-uuid: abc-123
      credentials-ref: prod-creds
credentials:
  - name: prod-creds
    credential:
      client-id: dt0s01.XXX
      client-secret: dt0s01.XXX.YYY
```

Environment variable overrides:
- `DTIAM_CONTEXT` - context name
- `DTIAM_OUTPUT` - output format
- `DTIAM_CLIENT_ID` - OAuth2 client ID
- `DTIAM_CLIENT_SECRET` - OAuth2 client secret
- `DTIAM_ACCOUNT_UUID` - account UUID

## Common Tasks

### Build

```bash
make build           # Build binary to bin/dtiam
make build-all       # Build for all platforms
make install         # Install to $GOPATH/bin
```

### Test

```bash
make test            # Run all tests
make test-coverage   # Run tests with coverage report
```

### Lint

```bash
make lint           # Run golangci-lint
make fmt            # Format code
make vet            # Run go vet
```

### Debug Authentication

```bash
# Verbose mode shows HTTP requests
dtiam -v get groups
```

Check `~/.config/dtiam/config` for credential configuration.

## Code Style

- Use Go idioms and conventions
- Error handling: return errors, don't panic
- Use `context.Context` for cancellation
- Close resources with defer
- Use interfaces for testability
- Keep packages focused and minimal

## Dependencies

```
github.com/spf13/cobra      # CLI framework
github.com/olekukonko/tablewriter  # Table output
golang.org/x/oauth2         # OAuth2 support
gopkg.in/yaml.v3            # YAML parsing
```

## Troubleshooting

### "No context configured"

Run:
```bash
dtiam config set-credentials NAME --client-id XXX --client-secret YYY
dtiam config set-context NAME --account-uuid UUID --credentials-ref NAME
dtiam config use-context NAME
```

### "Permission denied"

OAuth2 client needs appropriate scopes:
- `account-idm-read` / `account-idm-write`
- `iam-policies-management`
- `account-env-read`
- `iam:effective-permissions:read` (for effective permissions API)

### Build Errors

Ensure Go 1.22+ is installed:
```bash
go version
go mod tidy
```
