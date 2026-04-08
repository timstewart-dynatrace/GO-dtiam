# Architecture

## Project Structure

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
│   │   ├── cache/                    # Cache management
│   │   ├── bulk/                     # Bulk operations from files
│   │   ├── export/                   # Export resources for backup
│   │   └── analyze/                  # Permission analysis commands
│   ├── config/
│   │   ├── config.go                 # Config structs
│   │   └── loader.go                 # Config load/save, XDG paths
│   ├── client/
│   │   ├── client.go                 # HTTP client with retry
│   │   ├── errors.go                 # APIError type
│   │   └── urls.go                   # Centralized API URL constants
│   ├── auth/
│   │   ├── auth.go                   # TokenProvider interface
│   │   ├── oauth.go                  # OAuth2 token manager
│   │   └── bearer.go                 # Static bearer token
│   ├── resources/
│   │   ├── handler.go                # Handler interfaces + BaseHandler
│   │   ├── types.go                  # Typed response structs with table tags
│   │   ├── groups.go                 # GroupHandler
│   │   ├── users.go                  # UserHandler
│   │   ├── policies.go               # PolicyHandler
│   │   ├── bindings.go               # BindingHandler
│   │   ├── boundaries.go             # BoundaryHandler
│   │   ├── environments.go           # EnvironmentHandler
│   │   ├── serviceusers.go           # ServiceUserHandler
│   │   ├── limits.go                 # LimitsHandler
│   │   ├── subscriptions.go          # SubscriptionHandler
│   │   ├── tokens.go                 # TokenHandler (platform tokens)
│   │   ├── apps.go                   # AppHandler (App Engine Registry)
│   │   └── schemas.go                # SchemaHandler (Settings API)
│   ├── output/
│   │   ├── format.go                 # Format enum
│   │   ├── printer.go                # Unified Printer
│   │   ├── structprinter.go          # Struct-tag based printer
│   │   ├── table.go                  # Table formatter
│   │   └── columns.go                # Column definitions
│   ├── prompt/
│   │   └── confirm.go                # Confirmation prompts (Confirm, ConfirmDelete)
│   ├── diagnostic/
│   │   └── error.go                  # Enhanced errors with exit codes and suggestions
│   ├── logging/
│   │   └── logger.go                 # Structured logging with logrus
│   ├── suggest/
│   │   └── suggest.go                # Levenshtein command/flag suggestions
│   └── utils/
│       ├── permissions.go            # Permissions calculator, matrix, effective API
│       └── safemap.go                # Safe type assertion helpers
├── pkg/version/version.go            # Version info
├── go.mod
├── Makefile
└── .goreleaser.yaml
```

## Key Components

### CLI Layer (`internal/cli/`)
- Root command with global flags (`--context`, `--output`, `--verbose`, `--plain`, `--dry-run`)
- GlobalState singleton accessed by all commands
- Printer factory method on GlobalState

### Command Layer (`internal/commands/`)
- Verb-noun pattern: `get groups`, `create group`, `delete policy`
- Each verb is a package with a single exported `Cmd`
- Commands use `common.CreateClient()` for API access
- All follow `command-standards.md`

### Resource Layer (`internal/resources/`)
- `BaseHandler` provides generic CRUD via HTTP methods
- Concrete handlers embed BaseHandler and override as needed
- Handler interface: `List()`, `Get()`, `Create()`, `Update()`, `Delete()`

### Output Layer (`internal/output/`)
- Unified Printer supports table, wide, JSON, YAML, CSV
- `--plain` mode forces JSON for machine consumption
- Column definitions per resource type

### Auth Layer (`internal/auth/`)
- `TokenProvider` interface with OAuth2 and Bearer implementations
- OAuth2 auto-refreshes expired tokens
- Bearer is static (no refresh)

## Data Flow

```
CLI Command → common.CreateClient() → Auth (OAuth2/Bearer)
    → Resource Handler → HTTP Client (Resty) → Dynatrace API
    → Response → Printer (table/json/yaml/csv) → stdout
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

**Resolution API** (effective permissions):
`https://api.dynatrace.com/iam/v1/resolution/{level_type}/{level_id}/effectivepermissions`

**App Engine Registry API**:
`https://{environment-id}.apps.dynatrace.com/platform/app-engine/registry/v1/apps`

Level types: `account`, `environment`, `global`

## API Coverage

### Implemented

| Endpoint | Operation | Handler Method |
|----------|-----------|----------------|
| `GET /groups` | List groups | `GroupHandler.List()` |
| `GET /groups/{uuid}` | Get group | `GroupHandler.Get()` |
| `POST /groups` | Create group | `GroupHandler.Create()` |
| `PUT /groups/{uuid}` | Update group | `GroupHandler.Update()` |
| `DELETE /groups/{uuid}` | Delete group | `GroupHandler.Delete()` |
| `GET /users` | List users | `UserHandler.List()` |
| `GET /users/{uid}` | Get user | `UserHandler.Get()` |
| `POST /users` | Create user | `UserHandler.Create()` |
| `DELETE /users/{uid}` | Delete user | `UserHandler.Delete()` |
| `PUT /users/{email}/groups` | Replace user's groups | `UserHandler.ReplaceGroups()` |
| `DELETE /users/{email}/groups` | Remove from groups | `UserHandler.RemoveFromGroups()` |
| `POST /users/{email}` | Add to multiple groups | `UserHandler.AddToGroups()` |
| `GET /service-users` | List service users | `ServiceUserHandler.List()` |
| `POST /service-users` | Create service user | `ServiceUserHandler.Create()` |
| `DELETE /service-users/{uid}` | Delete service user | `ServiceUserHandler.Delete()` |
| `GET /policies` | List policies | `PolicyHandler.List()` |
| `POST /policies` | Create policy | `PolicyHandler.Create()` |
| `DELETE /policies/{uuid}` | Delete policy | `PolicyHandler.Delete()` |
| `GET /bindings` | List bindings | `BindingHandler.List()` |
| `POST /bindings` | Create binding | `BindingHandler.Create()` |
| `DELETE /bindings` | Delete binding | `BindingHandler.Delete()` |
| `GET /boundaries` | List boundaries | `BoundaryHandler.List()` |
| `POST /boundaries` | Create boundary | `BoundaryHandler.Create()` |
| `DELETE /boundaries/{uuid}` | Delete boundary | `BoundaryHandler.Delete()` |
| `GET /limits` | List limits | `LimitsHandler.List()` |
| `GET /subscriptions` | List subscriptions | `SubscriptionHandler.List()` |
| `GET /environments` | List environments | `EnvironmentHandler.List()` |

### Bulk Operations

| Command | Description |
|---------|-------------|
| `bulk add-users-to-group` | Add users from file |
| `bulk remove-users-from-group` | Remove users from file |
| `bulk create-groups` | Create groups from file |
| `bulk create-bindings` | Create bindings from file |
| `bulk export-group-members` | Export group members |

### Export Operations

| Command | Description |
|---------|-------------|
| `export all` | Export all resources |
| `export group` | Export single group |
| `export policy` | Export single policy |

### Analyze Operations

| Command | Description |
|---------|-------------|
| `analyze user-permissions` | Calculate user permissions |
| `analyze group-permissions` | Calculate group permissions |
| `analyze permissions-matrix` | Generate permissions matrix |
| `analyze policy` | Analyze policy permissions |
| `analyze least-privilege` | Least privilege compliance |
| `analyze effective-user` | Get user permissions via API |
| `analyze effective-group` | Get group permissions via API |

### Not Yet Implemented

| Feature | Description | Priority |
|---------|-------------|----------|
| `template` commands | Template-based resource creation | Medium |
| Caching | In-memory caching with TTL | Low |
