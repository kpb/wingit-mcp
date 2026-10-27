---
# wingit-mcp-dta3
title: Normalize query arguments before fetching observations
status: completed
type: bug
priority: high
tags:
    - correctness
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:22Z
updated_at: "2026-10-27T20:19:51Z"
---
## Problem
The live fetch receives raw arguments before the engine normalizes them. Omitted values currently query 0 km / 1 day but report 20 km / 7 days; upper bounds can also disagree.

## Code references
- cmd/wingit-mcp/main.go: live RecentNearby call precedes BuildTargetChecklist.
- internal/tools/target_engine.go: normalizeArgs.
- internal/ebird/client.go: independent API clamping.

## Acceptance criteria
- [x] Establish one validated, normalized query before observation loading.
- [x] Omitted/default, explicit, negative, and out-of-range inputs have documented behavior consistent with eBird limits.
- [x] Reject non-finite coordinates/numeric values where applicable.
- [x] HTTP query parameters and returned Filters agree for default and boundary cases.
- [x] Add regression tests using mocked HTTP, with no external network.

## Dependencies
None. Start here; source-selection tests and caching depend on this contract.

## Outcome and verification
Added exported `tools.NormalizeTargetArgs` and apply it before the observation loader; the engine uses the same normalization for reported filters. Non-positive radius/days use defaults (20 km / 7 days), values above eBird limits clamp to 50 km / 30 days, frequency stays clamped to [0, 1], and non-finite radius/frequency/coordinates are rejected. Documented query normalization in README.md.

Added a mocked-HTTP regression test that verifies requested `dist`/`back` parameters match returned filters for defaults, explicit values, non-positive inputs, and upper-bound clamping. Added normalization and invalid-number tests. No external network or credentials used.

Verification: `go test -count=1 -race ./...`, `go vet ./...`, and `git diff --check` passed. `make lint` ran `go vet` but could not complete because `staticcheck` is not installed. Changes are uncommitted and not pushed.
