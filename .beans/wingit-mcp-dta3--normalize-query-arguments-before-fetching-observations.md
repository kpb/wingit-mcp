---
# wingit-mcp-dta3
title: Normalize query arguments before fetching observations
status: todo
type: bug
priority: high
tags:
    - correctness
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:22Z
updated_at: 2026-10-27T19:59:22Z
---
## Problem
The live fetch receives raw arguments before the engine normalizes them. Omitted values currently query 0 km / 1 day but report 20 km / 7 days; upper bounds can also disagree.

## Code references
- cmd/wingit-mcp/main.go: live RecentNearby call precedes BuildTargetChecklist.
- internal/tools/target_engine.go: normalizeArgs.
- internal/ebird/client.go: independent API clamping.

## Acceptance criteria
- [ ] Establish one validated, normalized query before observation loading.
- [ ] Omitted/default, explicit, negative, and out-of-range inputs have documented behavior consistent with eBird limits.
- [ ] Reject non-finite coordinates/numeric values where applicable.
- [ ] HTTP query parameters and returned Filters agree for default and boundary cases.
- [ ] Add regression tests using mocked HTTP, with no external network.

## Dependencies
None. Start here; source-selection tests and caching depend on this contract.
