# Developer Porting Reference

This directory contains documentation for porting dtiam's IAM functionality into [dtctl](https://github.com/dynatrace-oss/dtctl). These documents capture design patterns, algorithms, and empirical API knowledge that would otherwise require reverse-engineering from dtiam's source code.

**This branch is a reference branch and should not be merged to main.**

---

## Documents

| Document | Purpose |
|----------|---------|
| [DTCTL_INTEGRATION_EVALUATION.md](DTCTL_INTEGRATION_EVALUATION.md) | Gap analysis of dtctl's IAM and Account namespace design documents against dtiam's feature set. Identifies what maps cleanly, what's missing, and recommendations. |
| [ANALYSIS_ALGORITHMS.md](ANALYSIS_ALGORITHMS.md) | Permission analysis algorithms used by the 7 `analyze` subcommands. Describes inputs, API call sequences, deduplication logic, and matrix construction. Primary porting target for `dtctl iam analyze`. |
| [API_QUIRKS.md](API_QUIRKS.md) | Empirically discovered Dynatrace Account Management API behaviors: response format inconsistencies, null vs absent fields, pagination patterns, error shapes, OAuth token semantics, and retry behavior. |
| [EXPORT_SCHEMA.md](EXPORT_SCHEMA.md) | Structure of YAML/JSON produced by dtiam's `export` commands. Defines field names, enrichment behavior, and CSV flattening rules. Compatibility target if `dtctl iam export` is implemented. |
| [BOUNDARY_QUERIES.md](BOUNDARY_QUERIES.md) | Boundary query string format and construction rules for app-id and schema-id boundaries. Includes the `buildBoundaryQuery()` function logic and validation workflow. |
| [GROUP_CLONE_ALGORITHM.md](GROUP_CLONE_ALGORITHM.md) | Compound operation sequence for `group clone`: create group, copy members, copy bindings. Documents error handling (best-effort, no rollback) and API call counts. |

## How to Use

These documents are intended for developers implementing `dtctl iam` and `dtctl account`. They complement the design proposals:

- [IAM_INTEGRATION_DESIGN.md](https://github.com/dynatrace-oss/dtctl/blob/docs/iam-integration-design/docs/dev/IAM_INTEGRATION_DESIGN.md) — command structure, config, auth, handlers
- [ACCOUNT_NAMESPACE_DESIGN.md](https://github.com/dynatrace-oss/dtctl/blob/docs/account-namespace-design/docs/dev/ACCOUNT_NAMESPACE_DESIGN.md) — subscriptions, cost, audit, notifications

The design docs describe *what* to build. These docs describe *how dtiam built it* and *what the API actually does* (vs what the docs say).
