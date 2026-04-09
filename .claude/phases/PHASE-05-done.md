# Phase 05 — Templates, Apply & v2.0.0 Release (FINAL)
Status: DONE

## Goal
Template engine, template commands, declarative `apply` command, bulk create-groups-with-policies, then v2.0.0 release.

## Tasks

### 5.1 Template Engine ✅
- [x] `internal/template/engine.go` — Go text/template renderer with `default` function
- [x] `internal/template/store.go` — XDG filesystem store for custom templates
- [x] `internal/template/builtin.go` — 5 embedded templates via `//go:embed`
- [x] Tests for rendering, variable parsing, extraction, builtin listing

### 5.2 Template Commands ✅
- [x] `template list/show/render/apply/save/delete/path` — all 7 subcommands
- [x] Registered in `cmd/dtiam/main.go`
- [x] Tests for subcommands, flags, getTemplate, columns

### 5.3 Apply Command ✅
- [x] `dtiam apply -f resource.yaml` with `--set` variables and `--dry-run`
- [x] Auto-detect kind (Group, Policy, Boundary, Binding)
- [x] Multi-document YAML support (`---` separators)
- [x] Tests for flags, YAML splitting

### 5.4 Export as Template ✅
- [x] `export policy --as-template` updated to Go template syntax (`{{.name}}`)
- [x] Compatible with `dtiam template apply`

### 5.5 Bulk Groups+Policies ✅
- [x] `bulk create-groups-with-policies --file FILE`
- [x] Creates/reuses groups, resolves policies, creates bindings
- [x] Supports CSV/YAML/JSON, `--continue-on-error`, `--dry-run`

### 5.6 Documentation Sweep ✅
- [x] All docs updated for Phase 5 features

### 5.7 Version Bump ✅
- [x] `pkg/version/version.go` → 2.0.0
- [x] All version references updated across docs

## Acceptance Criteria
- [x] `dtiam template list` shows 5 built-in templates
- [x] `dtiam apply --help` shows correct help
- [x] `dtiam bulk create-groups-with-policies --help` shows correct help
- [x] All 29 packages pass tests
- [x] Version 2.0.0 in version.go
