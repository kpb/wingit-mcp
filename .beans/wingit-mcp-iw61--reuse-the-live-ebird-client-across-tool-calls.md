---
# wingit-mcp-iw61
title: Reuse the live eBird client across tool calls
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
The MCP handler constructs a fresh client per call, resetting lastRequest and bypassing pacing across calls even though client-level pacing is tested.

## Code references
- cmd/wingit-mcp/main.go: NewClientFromEnv inside target_checklist handler.
- internal/ebird/client.go: Client.waitForRateLimit.

## Acceptance criteria
- [ ] Reuse one appropriately scoped client/pacing state across live tool calls.
- [ ] Offline mode still works without a token and preserves source precedence.
- [ ] Sequential and concurrent handler calls respect the configured request interval.
- [ ] Canceled waiting requests return promptly; state remains concurrency-safe.
- [ ] Add handler-level regression coverage using httptest and pass the race detector.

## Dependencies
None. Preserve existing public behavior except the pacing fix.
