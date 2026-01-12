# Command Reference

> **DISCLAIMER:** This tool is provided "as-is" without warranty. Use at your own risk. This is an independent, community-developed tool and is **NOT produced, endorsed, or supported by Dynatrace**.

Complete reference for all dtiam commands and their options.

## Table of Contents

- [Global Options](#global-options)
- [Environment Variables](#environment-variables)
- [config](#config) - Configuration management
- [get](#get) - List/retrieve resources
- [describe](#describe) - Detailed resource information
- [create](#create) - Create resources
- [delete](#delete) - Delete resources
- [user](#user) - User management
- [service-user](#service-user) - Service user (OAuth client) management
- [group](#group) - Advanced group operations
- [boundary](#boundary) - Boundary management
- [account](#account) - Account limits and subscriptions
- [cache](#cache) - Cache management

---

## Global Options

These options apply to all commands:

```bash
dtiam [OPTIONS] COMMAND [ARGS]
```

| Option            | Short | Description                                 |
| ----------------- | ----- | ------------------------------------------- |
| `--context TEXT`  | `-c`  | Override the current context                |
| `--output FORMAT` | `-o`  | Output format: table, json, yaml, csv, wide |
| `--verbose`       | `-v`  | Enable verbose/debug output                 |
| `--plain`         |       | Plain output mode (no colors, no prompts)   |
| `--dry-run`       |       | Preview changes without applying them       |
| `--version`       | `-V`  | Show version and exit                       |
| `--help`          |       | Show help message                           |

---

## Environment Variables

dtiam supports authentication and configuration via environment variables:

### Authentication Variables

| Variable              | Description            | Use Case                 |
| --------------------- | ---------------------- | ------------------------ |
| `DTIAM_BEARER_TOKEN`  | Static bearer token    | Quick testing, debugging |
| `DTIAM_CLIENT_ID`     | OAuth2 client ID       | Automation (recommended) |
| `DTIAM_CLIENT_SECRET` | OAuth2 client secret   | Automation (recommended) |
| `DTIAM_ACCOUNT_UUID`  | Dynatrace account UUID | Required for all methods |

### Configuration Variables

| Variable        | Description                   |
| --------------- | ----------------------------- |
| `DTIAM_CONTEXT` | Override current context name |
| `DTIAM_OUTPUT`  | Default output format         |
| `DTIAM_VERBOSE` | Enable verbose mode           |

### Authentication Priority

When multiple authentication methods are configured:

1. **Bearer Token** - `DTIAM_BEARER_TOKEN` + `DTIAM_ACCOUNT_UUID`
2. **OAuth2 (env)** - `DTIAM_CLIENT_ID` + `DTIAM_CLIENT_SECRET` + `DTIAM_ACCOUNT_UUID`
3. **Config file** - Context with OAuth2 credentials

### OAuth2 vs Bearer Token

| Feature          | OAuth2 (Recommended) | Bearer Token       |
| ---------------- | -------------------- | ------------------ |
| Auto-refresh     | ✅ Yes               | ❌ No              |
| Long-running     | ✅ Suitable          | ❌ Not recommended |
| Automation       | ✅ Recommended       | ❌ Not recommended |
| Quick testing    | ✅ Works             | ✅ Ideal           |
| Setup complexity | Medium               | Low                |

**Example: OAuth2 Authentication**

```bash
export DTIAM_CLIENT_ID="dt0s01.XXXXX"
export DTIAM_CLIENT_SECRET="dt0s01.XXXXX.YYYYY"
export DTIAM_ACCOUNT_UUID="abc-123-def"
dtiam get groups
```

**Example: Bearer Token Authentication**

```bash
# WARNING: Token will NOT auto-refresh!
export DTIAM_BEARER_TOKEN="dt0c01.XXXXX.YYYYY..."
export DTIAM_ACCOUNT_UUID="abc-123-def"
dtiam get groups
```

---

## config

Manage configuration contexts and credentials.

### config view

Display the current configuration.

```bash
dtiam config view
```

### config get-contexts

List all configured contexts.

```bash
dtiam config get-contexts
```

### config use-context

Switch to a different context.

```bash
dtiam config use-context NAME
```

| Argument | Description               |
| -------- | ------------------------- |
| `NAME`   | Context name to switch to |

### config set-context

Create or update a context.

```bash
dtiam config set-context NAME [OPTIONS]
```

| Argument/Option     | Short | Description                     |
| ------------------- | ----- | ------------------------------- |
| `NAME`              |       | Context name                    |
| `--account-uuid`    | `-a`  | Dynatrace account UUID          |
| `--credentials-ref` | `-c`  | Reference to a named credential |

**Example:**

```bash
dtiam config set-context prod --account-uuid abc-123 --credentials-ref prod-creds
```

### config delete-context

Delete a context.

```bash
dtiam config delete-context NAME
```

### config set-credentials

Store OAuth2 credentials.

```bash
dtiam config set-credentials NAME [OPTIONS]
```

| Argument/Option   | Short | Description          |
| ----------------- | ----- | -------------------- |
| `NAME`            |       | Credential name      |
| `--client-id`     | `-i`  | OAuth2 client ID     |
| `--client-secret` | `-s`  | OAuth2 client secret |

**Example:**

```bash
dtiam config set-credentials prod-creds --client-id dt0s01.XXX --client-secret dt0s01.XXX.YYY
```

### config delete-credentials

Delete stored credentials.

```bash
dtiam config delete-credentials NAME
```

### config get-credentials

List all stored credentials.

```bash
dtiam config get-credentials
```

---

## get

List or retrieve IAM resources.

### get groups

List or get IAM groups.

```bash
dtiam get groups [IDENTIFIER] [OPTIONS]
```

| Argument/Option | Short | Description                   |
| --------------- | ----- | ----------------------------- |
| `IDENTIFIER`    |       | Group UUID or name (optional) |
| `--output`      | `-o`  | Output format                 |

### get users

List or get IAM users.

```bash
dtiam get users [IDENTIFIER] [OPTIONS]
```

| Argument/Option | Short | Description                  |
| --------------- | ----- | ---------------------------- |
| `IDENTIFIER`    |       | User UID or email (optional) |
| `--output`      | `-o`  | Output format                |

### get policies

List or get IAM policies.

```bash
dtiam get policies [IDENTIFIER] [OPTIONS]
```

| Argument/Option | Short | Description                             |
| --------------- | ----- | --------------------------------------- |
| `IDENTIFIER`    |       | Policy UUID or name (optional)          |
| `--level`       | `-l`  | Policy level: account (default), global |
| `--output`      | `-o`  | Output format                           |

### get bindings

List IAM policy bindings.

```bash
dtiam get bindings [OPTIONS]
```

| Option     | Short | Description          |
| ---------- | ----- | -------------------- |
| `--group`  | `-g`  | Filter by group UUID |
| `--output` | `-o`  | Output format        |

### get environments

List or get Dynatrace environments.

```bash
dtiam get environments [IDENTIFIER] [OPTIONS]
```

| Argument/Option | Short | Description                       |
| --------------- | ----- | --------------------------------- |
| `IDENTIFIER`    |       | Environment ID or name (optional) |
| `--output`      | `-o`  | Output format                     |

### get boundaries

List or get IAM policy boundaries.

```bash
dtiam get boundaries [IDENTIFIER] [OPTIONS]
```

| Argument/Option | Short | Description                      |
| --------------- | ----- | -------------------------------- |
| `IDENTIFIER`    |       | Boundary UUID or name (optional) |
| `--output`      | `-o`  | Output format                    |

---

## describe

Show detailed resource information.

### describe group

Show detailed information about an IAM group.

```bash
dtiam describe group IDENTIFIER [--output FORMAT]
```

Displays: UUID, name, description, member count, members list, policy bindings.

### describe user

Show detailed information about an IAM user.

```bash
dtiam describe user IDENTIFIER [--output FORMAT]
```

Displays: UID, email, status, creation date, group memberships.

### describe policy

Show detailed information about an IAM policy.

```bash
dtiam describe policy IDENTIFIER [OPTIONS]
```

| Option     | Short | Description                             |
| ---------- | ----- | --------------------------------------- |
| `--level`  | `-l`  | Policy level: account (default), global |
| `--output` | `-o`  | Output format                           |

Displays: UUID, name, description, statement query, parsed permissions.

### describe environment

Show detailed information about a Dynatrace environment.

```bash
dtiam describe environment IDENTIFIER [--output FORMAT]
```

### describe boundary

Show detailed information about an IAM policy boundary.

```bash
dtiam describe boundary IDENTIFIER [--output FORMAT]
```

Displays: UUID, name, description, boundary query, attached policies count.

---

## create

Create IAM resources.

### create group

Create a new IAM group.

```bash
dtiam create group [OPTIONS]
```

| Option          | Short | Description           |
| --------------- | ----- | --------------------- |
| `--name`        | `-n`  | Group name (required) |
| `--description` | `-d`  | Group description     |
| `--output`      | `-o`  | Output format         |

**Example:**

```bash
dtiam create group --name "DevOps Team" --description "Platform engineering"
```

### create policy

Create a new IAM policy.

```bash
dtiam create policy [OPTIONS]
```

| Option          | Short | Description                       |
| --------------- | ----- | --------------------------------- |
| `--name`        | `-n`  | Policy name (required)            |
| `--statement`   | `-s`  | Policy statement query (required) |
| `--description` | `-d`  | Policy description                |
| `--output`      | `-o`  | Output format                     |

**Example:**

```bash
dtiam create policy --name "viewer" --statement "ALLOW settings:objects:read;"
```

### create binding

Create a policy binding (bind a policy to a group).

```bash
dtiam create binding [OPTIONS]
```

| Option       | Short | Description                      |
| ------------ | ----- | -------------------------------- |
| `--group`    | `-g`  | Group UUID or name (required)    |
| `--policy`   | `-p`  | Policy UUID or name (required)   |
| `--boundary` | `-b`  | Boundary UUID or name (optional) |
| `--output`   | `-o`  | Output format                    |

**Example:**

```bash
dtiam create binding --group "DevOps Team" --policy "admin-policy"
```

### create boundary

Create a new IAM policy boundary.

```bash
dtiam create boundary [OPTIONS]
```

| Option          | Short | Description                        |
| --------------- | ----- | ---------------------------------- |
| `--name`        | `-n`  | Boundary name (required)           |
| `--zones`       | `-z`  | Management zones (comma-separated) |
| `--query`       | `-q`  | Custom boundary query              |
| `--description` | `-d`  | Boundary description               |
| `--output`      | `-o`  | Output format                      |

Either `--zones` or `--query` must be provided.

**Example:**

```bash
dtiam create boundary --name "prod-only" --zones "Production,Staging"
```

---

## delete

Delete IAM resources.

### delete group

Delete an IAM group.

```bash
dtiam delete group IDENTIFIER [--force]
```

| Option    | Short | Description       |
| --------- | ----- | ----------------- |
| `--force` | `-f`  | Skip confirmation |

### delete policy

Delete an IAM policy.

```bash
dtiam delete policy IDENTIFIER [OPTIONS]
```

| Option    | Short | Description       |
| --------- | ----- | ----------------- |
| `--force` | `-f`  | Skip confirmation |

### delete binding

Delete a policy binding.

```bash
dtiam delete binding [OPTIONS]
```

| Option     | Short | Description            |
| ---------- | ----- | ---------------------- |
| `--group`  | `-g`  | Group UUID (required)  |
| `--policy` | `-p`  | Policy UUID (required) |
| `--force`  | `-f`  | Skip confirmation      |

### delete boundary

Delete an IAM policy boundary.

```bash
dtiam delete boundary IDENTIFIER [--force]
```

### delete user

Delete an IAM user.

```bash
dtiam delete user IDENTIFIER [--force]
```

### delete service-user

Delete a service user.

```bash
dtiam delete service-user IDENTIFIER [--force]
```

---

## user

User management operations.

### user create

Create a new user in the account.

```bash
dtiam user create EMAIL [OPTIONS]
```

| Argument/Option | Short | Description                          |
| --------------- | ----- | ------------------------------------ |
| `EMAIL`         |       | User email address (required)        |
| `--first-name`  |       | User's first name                    |
| `--last-name`   |       | User's last name                     |
| `--groups`      | `-g`  | Comma-separated group UUIDs or names |
| `--output`      | `-o`  | Output format                        |

**Examples:**

```bash
dtiam user create user@example.com
dtiam user create user@example.com --first-name John --last-name Doe
dtiam user create user@example.com --groups "DevOps,Platform"
```

### user add-to-groups

Add a user to multiple groups.

```bash
dtiam user add-to-groups EMAIL [OPTIONS]
```

| Argument/Option | Short | Description                                     |
| --------------- | ----- | ----------------------------------------------- |
| `EMAIL`         |       | User email address                              |
| `--groups`      | `-g`  | Comma-separated group UUIDs or names (required) |

**Example:**

```bash
dtiam user add-to-groups user@example.com --groups "DevOps,Platform"
```

### user remove-from-groups

Remove a user from multiple groups.

```bash
dtiam user remove-from-groups EMAIL [OPTIONS]
```

| Argument/Option | Short | Description                                     |
| --------------- | ----- | ----------------------------------------------- |
| `EMAIL`         |       | User email address                              |
| `--groups`      | `-g`  | Comma-separated group UUIDs or names (required) |

### user replace-groups

Replace all group memberships for a user.

```bash
dtiam user replace-groups EMAIL [OPTIONS]
```

| Argument/Option | Short | Description                          |
| --------------- | ----- | ------------------------------------ |
| `EMAIL`         |       | User email address                   |
| `--groups`      | `-g`  | Comma-separated group UUIDs or names |

**Example:**

```bash
dtiam user replace-groups user@example.com --groups "DevOps,Platform"
```

### user list-groups

List all groups a user belongs to.

```bash
dtiam user list-groups IDENTIFIER [--output FORMAT]
```

---

## service-user

Service user (OAuth client) management.

Service users are used for programmatic API access. When you create a service user, you receive OAuth client credentials that can be used to authenticate API requests.

### service-user list

List all service users in the account.

```bash
dtiam service-user list [--output FORMAT]
```

### service-user get

Get details of a service user.

```bash
dtiam service-user get USER [--output FORMAT]
```

| Argument | Description               |
| -------- | ------------------------- |
| `USER`   | Service user UUID or name |

### service-user create

Create a new service user (OAuth client).

**IMPORTANT:** Save the client secret immediately - it cannot be retrieved later!

```bash
dtiam service-user create [OPTIONS]
```

| Option          | Short | Description                          |
| --------------- | ----- | ------------------------------------ |
| `--name`        | `-n`  | Service user name (required)         |
| `--description` | `-d`  | Description                          |
| `--groups`      | `-g`  | Comma-separated group UUIDs or names |
| `--output`      | `-o`  | Output format                        |

**Examples:**

```bash
dtiam service-user create --name "CI Pipeline"
dtiam service-user create --name "CI Pipeline" --groups "DevOps,Automation"
```

### service-user update

Update a service user.

```bash
dtiam service-user update USER [OPTIONS]
```

| Argument/Option | Short | Description               |
| --------------- | ----- | ------------------------- |
| `USER`          |       | Service user UUID or name |
| `--name`        | `-n`  | New name                  |
| `--description` | `-d`  | New description           |
| `--output`      | `-o`  | Output format             |

### service-user delete

Delete a service user.

```bash
dtiam service-user delete USER [--force]
```

**Warning:** Deleting a service user will invalidate any OAuth tokens issued to it.

### service-user add-to-group

Add a service user to a group.

```bash
dtiam service-user add-to-group USER [OPTIONS]
```

| Argument/Option | Short | Description                   |
| --------------- | ----- | ----------------------------- |
| `USER`          |       | Service user UUID or name     |
| `--group`       | `-g`  | Group UUID or name (required) |

### service-user remove-from-group

Remove a service user from a group.

```bash
dtiam service-user remove-from-group USER [OPTIONS]
```

| Argument/Option | Short | Description                   |
| --------------- | ----- | ----------------------------- |
| `USER`          |       | Service user UUID or name     |
| `--group`       | `-g`  | Group UUID or name (required) |

### service-user list-groups

List all groups a service user belongs to.

```bash
dtiam service-user list-groups USER [--output FORMAT]
```

---

## group

Advanced group operations.

### group members

List all members of a group.

```bash
dtiam group members IDENTIFIER [--output FORMAT]
```

### group add-member

Add a user to a group.

```bash
dtiam group add-member IDENTIFIER [OPTIONS]
```

| Argument/Option | Short | Description                        |
| --------------- | ----- | ---------------------------------- |
| `IDENTIFIER`    |       | Group UUID or name                 |
| `--email`       | `-e`  | User email address to add (required) |

### group remove-member

Remove a user from a group.

```bash
dtiam group remove-member IDENTIFIER [OPTIONS]
```

| Argument/Option | Short | Description                  |
| --------------- | ----- | ---------------------------- |
| `IDENTIFIER`    |       | Group UUID or name           |
| `--user`        | `-u`  | User UID to remove (required) |

### group bindings

List all policy bindings for a group.

```bash
dtiam group bindings IDENTIFIER [--output FORMAT]
```

---

## boundary

Boundary attach/detach operations.

### boundary attach

Attach a boundary to an existing binding.

```bash
dtiam boundary attach [OPTIONS]
```

| Option       | Short | Description                      |
| ------------ | ----- | -------------------------------- |
| `--group`    | `-g`  | Group UUID or name (required)    |
| `--policy`   | `-p`  | Policy UUID or name (required)   |
| `--boundary` | `-b`  | Boundary UUID or name (required) |

**Example:**

```bash
dtiam boundary attach --group "DevOps" --policy "admin-policy" --boundary "prod-boundary"
```

### boundary detach

Detach a boundary from a binding.

```bash
dtiam boundary detach [OPTIONS]
```

| Option       | Short | Description                      |
| ------------ | ----- | -------------------------------- |
| `--group`    | `-g`  | Group UUID or name (required)    |
| `--policy`   | `-p`  | Policy UUID or name (required)   |
| `--boundary` | `-b`  | Boundary UUID or name (required) |

### boundary list-attached

List all bindings that use a boundary.

```bash
dtiam boundary list-attached BOUNDARY [--output FORMAT]
```

---

## account

Account limits and subscription information.

### account limits

List account limits and quotas.

```bash
dtiam account limits [--output FORMAT]
```

Shows current usage and maximum allowed values for account resources like users, groups, and environments.

### account subscriptions

List account subscriptions.

```bash
dtiam account subscriptions [--output FORMAT]
```

Shows all subscriptions including type, status, and time period.

### account forecast

Get usage forecast for subscriptions.

```bash
dtiam account forecast [--output FORMAT]
```

---

## cache

Cache management.

### cache stats

Show cache statistics.

```bash
dtiam cache stats
```

### cache clear

Clear cache entries.

```bash
dtiam cache clear
```

---

## Exit Codes

| Code | Description                                         |
| ---- | --------------------------------------------------- |
| 0    | Success                                             |
| 1    | Error (resource not found, permission denied, etc.) |

## See Also

- [Quick Start Guide](QUICK_START.md)
- [Architecture](ARCHITECTURE.md)
- [API Reference](API_REFERENCE.md)
