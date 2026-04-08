# Phase 04 — Advanced Operations
Status: PENDING

## Goal
Group cloning, app/schema boundary helpers, group setup, parameterized policies.

## Prerequisites
- Phase 2 complete (done)
- Apps and schemas handlers available (done — Phase 3 original tasks 3.1-3.3)

## Tasks

### 4.1 Group Clone
- [ ] Add `group clone SOURCE --name NEW_NAME` subcommand to group/group.go
- Flags:
  - `--name` (required) — name for new group
  - `--description` — optional description
  - `--include-members` — copy group members to new group
  - `--include-policies` — copy policy bindings to new group
- Logic:
  1. Resolve source group by name/UUID
  2. Create new group with --name and --description
  3. If --include-members: iterate source members, add each to new group
  4. If --include-policies: iterate source bindings, create each for new group
- Supports --dry-run and --force

### 4.2 App Boundaries
- [ ] Add `boundary create-app-boundary NAME` subcommand to boundary/boundary.go
- Flags:
  - `--app-ids` (required, comma-separated)
  - `--not-in` — use NOT IN instead of IN (exclude apps)
  - `--environment` — for validation against App Engine Registry
  - `--description` — optional
  - `--skip-validation` — skip app ID validation
- Logic:
  1. If not --skip-validation: validate each app ID via AppHandler.Get
  2. Generate boundary query: `shared:app-id IN ("app1", "app2")`
  3. Create boundary via BoundaryHandler.Create

### 4.3 Schema Boundaries
- [ ] Add `boundary create-schema-boundary NAME` subcommand
- Same pattern as app boundaries:
  - `--schema-ids` (required), `--not-in`, `--environment`, `--description`, `--skip-validation`
  - Query: `settings:schemaId IN ("builtin:alerting.profile", ...)`
  - Validate via SchemaHandler.Get

### 4.4 Group Setup
- [ ] Add `group setup --name NAME --policies-file FILE` subcommand
  - Reads YAML/JSON file with policies, bindings, and optional boundaries
  - Creates group, then applies all policy bindings from file
  - One-step provisioning for new teams
  - Flags: `--name`, `--description`, `--policies-file`, `--dry-run`

### 4.5 Parameterized Policies
- [ ] Add `--param key=value` repeatable flag to `create binding` (create/create.go)
- [ ] Pass parameters in binding creation payload
- [ ] Update `describe policy` to show bind parameters when present
- [ ] Update BindingHandler.Create to accept optional parameters map

## Acceptance Criteria
- [ ] `dtiam group clone "Source" --name "Copy" --include-members --include-policies` works
- [ ] `dtiam boundary create-app-boundary "My Apps" --app-ids dynatrace.dashboards,dynatrace.logs` validates and creates
- [ ] `dtiam boundary create-schema-boundary "My Schemas" --schema-ids builtin:alerting.profile` validates and creates
- [ ] `dtiam group setup --name "Team" --policies-file setup.yaml` provisions group end-to-end
- [ ] `dtiam create binding --group X --policy Y --param env=production` sends parameters
- [ ] All new commands support --dry-run and --force where applicable
- [ ] Tests for clone logic, boundary query generation, parameter passing

## MANDATORY: Follow .claude/rules/command-standards.md for all new code
