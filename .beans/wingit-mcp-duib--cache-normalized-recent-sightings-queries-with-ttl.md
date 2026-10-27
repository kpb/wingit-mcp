---
# wingit-mcp-duib
title: Cache normalized recent-sightings queries with TTL
status: todo
type: feature
priority: normal
tags:
    - reliability
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:53Z
updated_at: "2026-10-27T20:00:37Z"
---
## Source
https://github.com/kpb/wingit-mcp/issues/2
Also owns the caching portion of https://github.com/kpb/wingit-mcp/issues/3.

## Scope
Avoid repeated live requests for the same effective lat/lon, radius, and days-back query. Prefer optional file persistence under WINGIT_CACHE_DIR for repeatable demos; decide and document the storage/TTL design before implementation. Offline fixtures remain independent of this live-query cache.

## Code references
- internal/ebird/client.go: request construction and response handling.
- cmd/wingit-mcp/main.go: source/client lifetime.

## Acceptance criteria
- [ ] Record cache storage, TTL, key normalization, concurrency, and data/privacy tradeoffs in an ADR.
- [ ] Keys reflect the effective HTTP query, including normalized coordinates and filters.
- [ ] Cache hits avoid HTTP requests; misses and expired entries fetch fresh data.
- [ ] Define optional WINGIT_CACHE_DIR behavior, including the unset/default case and corrupt/unwritable cache handling.
- [ ] Do not cache failed responses or include API tokens in cache contents/keys.
- [ ] Deterministic tests cover hit, miss, expiration, failure recovery, and concurrent access without external network.
- [ ] Document TTL/configuration and persistence/cleanup behavior.

## Dependencies
- Blocked by [wingit-mcp-dta3](wingit-mcp-dta3--normalize-query-arguments-before-fetching-observations.md): stable normalized query contract.
- Blocked by [wingit-mcp-iw61](wingit-mcp-iw61--reuse-the-live-ebird-client-across-tool-calls.md): shared client lifetime.