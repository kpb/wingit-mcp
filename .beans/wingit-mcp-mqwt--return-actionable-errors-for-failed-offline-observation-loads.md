---
# wingit-mcp-mqwt
title: Return actionable errors for failed offline observation loads
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
An unreadable or malformed WINGIT_RECENT_JSON file is logged and converted into an empty successful checklist, making data failure look like no candidate lifers.

## Code references
- cmd/wingit-mcp/main.go: LoadRecentNearby warning-and-continue branch.
- internal/ebird/loader.go: file read and decode errors.

## Acceptance criteria
- [ ] Missing, unreadable, and malformed configured observation files return actionable tool errors.
- [ ] A valid empty observation array remains a successful empty result.
- [ ] Do not silently fall back to live access when an explicit offline source fails.
- [ ] Test failures and empty success through the handler without external network.
- [ ] Keep diagnostics off stdout and avoid leaking file contents or secrets.

## Dependencies
None. Source-selection completion depends on this behavior.
