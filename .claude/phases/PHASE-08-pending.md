# Phase 08 — Polish & v2.0.0 Release
Status: PENDING

## Goal
Final quality pass, documentation sweep, version bump to 2.0.0.

## Prerequisites
- All prior phases complete (or consciously deferred)

## Tasks

### 8.1 Test Coverage
- [ ] Target >70% overall coverage
- [ ] Add tests for all new Phase 3-7 features
- [ ] Verify: `go test ./... -coverprofile=coverage.out`

### 8.2 Documentation Sweep
- [ ] Update `.claude/architecture.md` with any new packages/components
- [ ] Update `docs/COMMANDS.md` with all new commands
- [ ] Update `README.md` features table
- [ ] Update `docs/ARCHITECTURE.md` with new layers (cache, template, etc.)
- [ ] Update `CHANGELOG.md` — organize all changes by version

### 8.3 Version Bump
- [ ] Update `pkg/version/version.go` → `2.0.0`
- [ ] Create git tag: `git tag -a v2.0.0 -m "Release version 2.0.0"`
- [ ] Create GitHub release with changelog

### 8.4 Final Verification
- [ ] `make build` — binary compiles
- [ ] `make test` — all tests pass, >70% coverage
- [ ] `make lint` — no lint errors
- [ ] `dtiam version` shows 2.0.0
- [ ] Manual smoke test all major commands
- [ ] Verify backward compatibility

## Acceptance Criteria
- [ ] `go test ./... -coverprofile` shows >70%
- [ ] All documentation updated and complete
- [ ] `make build && make test && make lint` all pass
- [ ] `dtiam version` shows 2.0.0
- [ ] All prior phase acceptance criteria still met (no regressions)

## MANDATORY: Follow .claude/rules/command-standards.md for all new code
