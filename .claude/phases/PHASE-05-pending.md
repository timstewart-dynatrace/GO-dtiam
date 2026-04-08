# Phase 05 — Templates & Apply
Status: PENDING

## Goal
Template engine, template commands, declarative `apply` command, bulk create-groups-with-policies.

## Prerequisites
- Phase 4 complete (parameterized policies for template variable support)

## Reference
- go-dtctl-main: pkg/util/template/template.go, pkg/apply/apply.go
- Python-IAM-CLI: src/dtiam/commands/template.py, src/dtiam/utils/templates.py

## Tasks

### 5.1 Template Engine
- [ ] Create `internal/template/engine.go`:
  - Go `text/template` renderer with `--set key=value`
  - `RenderTemplate(content string, vars map[string]string) (string, error)`
  - `ParseSetFlags(flags []string) (map[string]string, error)`
  - Custom functions: `default` (provide default values)
- [ ] Create `internal/template/store.go`:
  - Template storage at `$XDG_DATA_HOME/dtiam/templates/`
  - List, Get, Save, Delete operations
- [ ] Create `internal/template/builtin.go`:
  - Embed built-in templates via Go `embed`:
    - `group-team`, `policy-readonly`, `policy-admin`, `binding-simple`, `boundary-mz`

### 5.2 Template Commands
- [ ] Create `internal/commands/template/template.go`:
  - `template list` — list available templates (built-in + custom)
  - `template show NAME` — display content and required variables
  - `template render NAME --set key=value` — render to stdout
  - `template apply NAME --set key=value` — render and create resource
  - `template save NAME --file FILE` — save custom template
  - `template delete NAME` — delete custom template (with --force)
  - `template path` — show templates directory
- [ ] Register in cmd/dtiam/main.go

### 5.3 Apply Command
- [ ] Create `internal/commands/apply/apply.go`:
  - `dtiam apply -f resource.yaml`
  - Auto-detect resource type from `kind` field
  - Route to handler: Group → GroupHandler.Create, Policy → PolicyHandler.Create, etc.
  - Support `--set key=value` for template variables
  - Handle create vs update (check existence by name/UUID)
  - Support `--dry-run`
- [ ] Register in cmd/dtiam/main.go

### 5.4 Export as Template
- [ ] Enhance `export policy --as-template` with full Go template syntax
  - Convert to `{{.name}}`, `{{.statement}}` placeholders

### 5.5 Bulk Create-Groups-With-Policies
- [ ] Add `bulk create-groups-with-policies --file FILE` to bulk/bulk.go
- CSV columns: group_name, description, policy_name, level, level_id, management_zones, boundary_name, parameters
- Per row: create group → resolve policy → create binding → optional boundary
- Flags: `--file`, `--continue-on-error`, `--dry-run`

## Acceptance Criteria
- [ ] `dtiam template list` shows built-in templates
- [ ] `dtiam template render policy-readonly --set name=MyPolicy` renders YAML
- [ ] `dtiam template apply policy-readonly --set name=MyPolicy` creates the policy
- [ ] `dtiam apply -f group.yaml` creates resource from file
- [ ] `dtiam apply -f group.yaml --dry-run` previews without creating
- [ ] `dtiam bulk create-groups-with-policies -f groups.csv` works end-to-end
- [ ] Tests for template rendering, variable substitution, apply routing

## MANDATORY: Follow .claude/rules/command-standards.md for all new code
