# Phase 03 — Quick Wins
Status: DONE

## Goal
Complete partial implementations and fill obvious gaps with minimal effort.

## Tasks

### 3.1 Account Capabilities ✅
- [x] Added `account capabilities [SUBSCRIPTION]` subcommand
- [x] Uses existing `SubscriptionHandler.GetCapabilities()`
- [x] Added `CapabilityColumns()` to output/columns.go
- [x] Tests for command and handler

### 3.2 Per-Resource Exports ✅
- [x] Added `export environments`, `export users`, `export bindings`, `export boundaries`, `export service-users`
- [x] All support `--output`, `--format`, `--prefix`, `--detailed` flags
- [x] Extracted `exportResourceToDir()` helper to reduce duplication
- [x] Tests for subcommand registration, flags, aliases, examples

### 3.3 User Info Alias ✅
- [x] Added `user info IDENTIFIER` subcommand
- [x] Equivalent to `describe user` — resolves by UID or email, shows expanded view
- [x] Tests for args and examples

## Acceptance Criteria
- [x] `dtiam account capabilities` shows capabilities table
- [x] `dtiam export users --format yaml` exports users
- [x] `dtiam user info EMAIL` shows user details
- [x] All new commands have Example help text
- [x] All new commands support --plain and -o json/yaml/table
- [x] Tests for each new command
- [x] All 26 packages pass tests
