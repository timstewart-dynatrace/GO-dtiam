# Phase 07 — Zones (Legacy)
Status: PENDING

## Goal
Management zone support via Dynatrace entities API. Legacy feature for Python CLI parity.

**Note**: Management zones are being superseded by Grail-based access control. This phase is low priority — implement only if there is user demand.

## Prerequisites
- Phase 2 complete (done)

## Tasks

### 7.1 Zone Handler
- [ ] Create `internal/resources/zones.go`:
  - `ZoneHandler` with List, Get
  - API: `{environment_url}/api/v2/entities?entitySelector=type("MANAGEMENT_ZONE")`
  - Requires `--environment` flag or config environment-url
- [ ] Add `ZoneColumns()` to output/columns.go

### 7.2 Zone Commands
- [ ] Add `get zones [NAME]` to get/get.go:
  - List management zones with --name filter, --environment flag
- [ ] Add `zones export` — export zones to YAML/JSON (--output-dir, --format)
- [ ] Add `zones compare-groups` — compare zone names with IAM group names
  - Useful for auditing: which zones have matching IAM groups?
- [ ] Mark all zone commands as legacy/deprecated in help text

## Acceptance Criteria
- [ ] `dtiam get zones --environment abc12345` lists management zones
- [ ] `dtiam get zones --name Production --environment abc12345` filters by name
- [ ] `dtiam zones compare-groups --environment abc12345` shows zone/group comparison
- [ ] All commands marked as legacy in help text
- [ ] Tests for zone handler and commands

## MANDATORY: Follow .claude/rules/command-standards.md for all new code
