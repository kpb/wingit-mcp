---
# wingit-mcp-iw61
title: Reuse the live eBird client across tool calls
status: completed
type: bug
priority: high
tags:
    - correctness
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:22Z
updated_at: "2026-10-27T20:56:54Z"
---
## Problem
The MCP handler constructs a fresh client per call, resetting lastRequest and bypassing pacing across calls even though client-level pacing is tested.

## Code references
- cmd/wingit-mcp/main.go: NewClientFromEnv inside target_checklist handler.
- internal/ebird/client.go: Client.waitForRateLimit.

## Acceptance criteria
- [x] Reuse one appropriately scoped client/pacing state across live tool calls.
- [x] Offline mode still works without a token and preserves source precedence.
- [x] Sequential and concurrent handler calls respect the configured request interval.
- [x] Canceled waiting requests return promptly; state remains concurrency-safe.
- [x] Add handler-level regression coverage using httptest and pass the race detector.

## Dependencies
None. Preserve existing public behavior except the pacing fix.

## Outcome and verification
Added a lazily initialized, concurrency-safe live eBird client provider shared by target-checklist calls. It is only initialized for live requests, preserving offline mode and avoiding a token requirement when WINGIT_RECENT_JSON is configured. Offline mode keeps source precedence over the live provider.

Added request-pacing coverage through the target-checklist loading path for sequential and concurrent calls, plus coverage that offline mode succeeds without a token and a canceled waiting live call returns without sending an additional request.

Verification: `go test -count=3 -race ./...`, `go vet ./...`, and `git diff --check` passed. No external network or credentials used. Changes are uncommitted and not pushed.
