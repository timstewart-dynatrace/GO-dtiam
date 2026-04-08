# Phase 06 — Analysis & Caching
Status: PENDING

## Goal
Enhanced analysis with permission diff and risk scoring, functional in-memory cache.

## Prerequisites
- Phase 2 complete (done)

## Tasks

### 6.1 Permission Diff
- [ ] Add `analyze diff-permissions ENTITY1 ENTITY2` subcommand
  - Compare effective permissions between two users or two groups
  - Show: permissions only in ENTITY1, only in ENTITY2, shared
  - Support `--format summary` (group by service)
  - Table output: permission | entity1 | entity2 | status (added/removed/shared)

### 6.2 Permission Gaps
- [ ] Add `analyze permission-gaps` subcommand
  - Find policies bound to 0 groups (unused)
  - Find groups with 0 policy bindings (no permissions)
  - Output as structured report
  - Support `--export FILE` for JSON/CSV

### 6.3 Enhanced Least-Privilege
- [ ] Enhance existing `analyze least-privilege`:
  - Add numeric risk scoring (1-10 per policy)
  - Add `--min-severity` filter flag (low, medium, high)
  - Detect: unused policies, duplicate policies, overlapping permissions
  - Support `--export FILE` for JSON/CSV report

### 6.4 TTL Cache
- [ ] Create `internal/cache/cache.go`:
  - Generic TTL cache with `sync.RWMutex`
  - Default TTL: 300 seconds (5 minutes)
  - Stats tracking: hits, misses, hit rate
- [ ] Create `internal/cache/middleware.go`:
  - HTTP middleware that caches GET responses by URL+params
  - Integrate into `internal/client/client.go`
  - Active within single CLI invocation (useful for analyze, export --detailed)

### 6.5 Cache Commands (Make Functional)
- [ ] Rewrite `internal/commands/cache/cache.go`:
  - `cache stats` — display hit rate, entry count, TTL
  - `cache clear` — clear all entries (with --force, --expired-only)
  - `cache keys` — list cached keys
  - `cache set-ttl SECONDS` — set default TTL
  - `cache reset-stats` — reset counters

## Acceptance Criteria
- [ ] `dtiam analyze diff-permissions user1@example.com user2@example.com` shows diff
- [ ] `dtiam analyze permission-gaps` finds unused policies and empty groups
- [ ] `dtiam analyze least-privilege --min-severity medium` filters findings
- [ ] `dtiam cache stats` shows real hit/miss statistics
- [ ] Running `dtiam export all --detailed` uses cache (fewer API calls)
- [ ] Tests for cache TTL, concurrent access, stats, analysis logic

## MANDATORY: Follow .claude/rules/command-standards.md for all new code
