# Decisions

## 2025-01-01 — kubectl-style CLI Design
**Chosen:** Verb-noun command structure modeled after kubectl and go-dtctl-main
**Alternatives:** POSIX-style flags only, interactive TUI, REST-wrapper approach
**Why:** Target users are Dynatrace admins familiar with kubectl. Verb-noun is discoverable (`get groups`, `delete policy`) and scriptable. go-dtctl-main provided a proven reference implementation.
**Trade-offs:** More complex command registration than flat flag-based CLI
**Revisit if:** Target audience shifts away from kubectl-familiar users

---

## 2025-01-01 — BaseHandler Pattern for Resources
**Chosen:** Generic `BaseHandler` struct that all resource handlers embed, providing default CRUD via HTTP
**Alternatives:** Individual handler implementations, code generation from OpenAPI spec
**Why:** Dynatrace IAM API follows consistent patterns (list/get/create/delete). BaseHandler eliminates 80% of boilerplate while allowing per-resource overrides.
**Trade-offs:** Less type safety than generated clients; handlers return `map[string]any` instead of typed structs
**Revisit if:** API diverges significantly between resources, or typed response structs are needed

---

## 2025-01-01 — Resty over stdlib HTTP
**Chosen:** `go-resty/resty/v2` for HTTP client
**Alternatives:** Standard `net/http`, `hashicorp/go-retryablehttp`
**Why:** Built-in retry with backoff, request/response middleware, cleaner API for REST operations. Reduces boilerplate for auth header injection, error handling, and JSON marshaling.
**Trade-offs:** External dependency for something stdlib can do
**Revisit if:** Resty maintenance stalls or dependency footprint becomes a concern

---

## 2025-01-01 — Unified Printer Abstraction
**Chosen:** Single `Printer` type that handles table, wide, JSON, YAML, CSV output with `--plain` mode for machine consumption
**Alternatives:** Per-format output functions, template-based rendering
**Why:** Every command needs consistent output formatting. `--plain` mode forces JSON and strips colors so AI agents and scripts get clean structured data. Centralizing this prevents format drift across 30+ commands.
**Trade-offs:** Printer is a relatively large interface; all commands coupled to it
**Revisit if:** Output requirements diverge significantly between command categories

---

## 2025-01-01 — OAuth2 with stdlib (no external dep)
**Chosen:** Custom OAuth2 implementation using `net/http` and `net/url`
**Alternatives:** `golang.org/x/oauth2`, third-party OAuth libraries
**Why:** Dynatrace OAuth2 flow is straightforward client_credentials grant. Custom implementation avoids pulling in `golang.org/x/oauth2` and its transitive dependencies for a simple token exchange.
**Trade-offs:** Must maintain token refresh logic manually
**Revisit if:** Need to support additional OAuth2 flows (authorization code, PKCE)

---

## 2025-01-01 — Client-side Filtering
**Chosen:** Fetch full resource list from API, filter client-side with `--name`/`--email`
**Alternatives:** Server-side filtering (API query params), GraphQL
**Why:** Dynatrace IAM API does not support server-side filtering for most resources. Client-side substring matching provides consistent UX across all resource types.
**Trade-offs:** Fetches more data than needed; won't scale to very large accounts
**Revisit if:** API adds server-side filtering, or accounts exceed 10k resources

---

## 2026-04-08 — Modular .claude/ Rule Structure
**Chosen:** Break monolithic root CLAUDE.md into `.claude/CLAUDE.md` + modular rule files under `.claude/rules/`
**Alternatives:** Keep single CLAUDE.md, use only `@` includes from root
**Why:** 500+ line CLAUDE.md mixed workflow rules, code patterns, API docs, and architecture. Modular files align with PROJECT-TEMPLATES standard, improve maintainability, and allow rules to be updated independently.
**Trade-offs:** More files to maintain; must keep root CLAUDE.md in sync as pointer
**Revisit if:** Claude Code changes how it loads instructions and modular files become unnecessary
