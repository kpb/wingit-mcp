---
# wingit-mcp-18w6
title: Complete and test observation source selection
status: todo
type: task
priority: normal
tags:
    - reliability
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:53Z
updated_at: "2026-10-27T20:00:37Z"
---
## Source
https://github.com/kpb/wingit-mcp/issues/3

## Current state and scope
Offline-first/live-second selection already exists and README describes both. Finish configuration errors and mode coverage instead of reimplementing completed behavior. Caching is a separate bean ([wingit-mcp-duib](wingit-mcp-duib--cache-normalized-recent-sightings-queries-with-ttl.md)); the original caching requirement is delegated there.

## Code references
- cmd/wingit-mcp/main.go: environment configuration and handler.
- README.md: required environment variables and recent observations.

## Acceptance criteria
- [ ] WINGIT_RECENT_JSON wins when both recent-source settings are present.
- [ ] Offline mode works without an API token.
- [ ] Otherwise WINGIT_EBIRD_TOKEN selects live access with normalized coordinates/query settings.
- [ ] Neither setting produces a friendly error explaining both supported source configurations.
- [ ] WINGIT_PERSONAL_JSON remains required.
- [ ] Tests cover each mode, precedence, and failure behavior without external network.
- [ ] README matches behavior; stdout remains protocol-only.

## Dependencies
- Blocked by [wingit-mcp-dta3](wingit-mcp-dta3--normalize-query-arguments-before-fetching-observations.md): normalize arguments before fetch.
- Blocked by [wingit-mcp-iw61](wingit-mcp-iw61--reuse-the-live-ebird-client-across-tool-calls.md): reuse live client/pacing state.
- Blocked by [wingit-mcp-mqwt](wingit-mcp-mqwt--return-actionable-errors-for-failed-offline-observation-loads.md): actionable offline-load errors.