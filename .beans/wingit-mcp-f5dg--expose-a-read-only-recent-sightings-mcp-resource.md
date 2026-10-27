---
# wingit-mcp-f5dg
title: Expose a read-only recent-sightings MCP resource
status: todo
type: feature
priority: low
tags:
    - optional
    - stretch
created_at: 2026-10-27T19:59:22Z
updated_at: "2026-10-27T20:00:37Z"
---
## Source
https://github.com/kpb/wingit-mcp/issues/7

## Scope
Optional follow-up, intentionally outside the release-critical v0.2.0 milestone. Expose last-known loaded/fetched sightings for host UX and debugging, without triggering hidden fetches.

## Code references
- internal/mcp/resources.go: existing personal-checklist resource pattern.
- cmd/wingit-mcp/main.go: observation source and server wiring.

## Acceptance criteria
- [ ] resources/list advertises wingit://recent-sightings.
- [ ] resources/read returns valid JSON with an explicit empty-before-first-load state.
- [ ] Define last-successful-data behavior after failed queries.
- [ ] Shared observation state is safe under concurrent reads and tool calls.
- [ ] Add offline MCP integration and race-enabled tests.
- [ ] Document the resource and its data visibility/privacy implications.

## Dependencies
- Blocked by [wingit-mcp-18w6](wingit-mcp-18w6--complete-and-test-observation-source-selection.md): completed observation-source selection.
- Blocked by [wingit-mcp-vxe1](wingit-mcp-vxe1--add-offline-mcp-tools-prompts-and-resources-integration-tests.md): MCP integration coverage.