# Phase 03 — Quick Wins
Status: PENDING

## Goal
Complete partial implementations and fill obvious gaps with minimal effort.

## Prerequisites
- Phase 2 complete (done)

## Tasks

### 3.1 Account Capabilities
- [ ] Add `account capabilities [SUBSCRIPTION]` subcommand to account/account.go
  - Uses existing `SubscriptionHandler.GetCapabilities()` (internal/resources/subscriptions.go)
  - Optional subscription UUID/name argument to filter by subscription
  - Table output with capability name, enabled status, subscription name
- [ ] Add `CapabilityColumns()` to output/columns.go
- [ ] Tests for command flag parsing and output

### 3.2 Per-Resource Exports
- [ ] Add individual export subcommands to export/export.go:
  - `export environments` — export all environments
  - `export users` — export all users
  - `export bindings` — export all bindings
  - `export boundaries` — export all boundaries
  - `export service-users` — export all service users
- [ ] Same flags as `export all`: `--output`, `--format`, `--detailed`
- [ ] Tests for each new subcommand

### 3.3 User Info Alias
- [ ] Add `user info IDENTIFIER` subcommand to user/user.go
  - Alias for `describe user` — resolves user and shows detailed view
  - For Python CLI parity (users expect `user info`)
- [ ] Tests for command

## Acceptance Criteria
- [ ] `dtiam account capabilities` shows capabilities table
- [ ] `dtiam export users --format yaml` exports users
- [ ] `dtiam user info EMAIL` shows user details
- [ ] All new commands have Example help text
- [ ] All new commands support --plain and -o json/yaml/table
- [ ] Tests for each new command

## MANDATORY: Follow .claude/rules/command-standards.md for all new code
