---
# wingit-mcp-2iuj
title: Prevent canceled live requests from consuming pacing slots
status: todo
type: bug
priority: normal
tags:
    - correctness
parent: wingit-mcp-7zf5
created_at: 2026-10-27T22:53:46Z
updated_at: 2026-10-27T22:53:46Z
---
## Problem
The shared eBird client reserves its next pacing timestamp before waiting. If a request context is canceled during that wait, the reservation remains in `lastRequest`; later valid requests are delayed by the abandoned slot. Repeated canceled waits can extend the delay further.

## Code references
- internal/ebird/client.go: Client.waitForRateLimit reserves `lastRequest` before selecting on `ctx.Done()`.
- cmd/wingit-mcp/main.go: live client is now shared across tool calls.
- cmd/wingit-mcp/main_test.go: current cancellation test checks prompt return/no HTTP request but not recovery of the next valid call.

## Acceptance criteria
- [ ] A canceled request waiting for pacing does not consume or extend a future request slot.
- [ ] Concurrent requests remain rate-limited and race-safe when one or more waiters cancel.
- [ ] After a canceled wait, the next valid request proceeds according to the configured interval, not an extra abandoned interval.
- [ ] Tests cover canceled-wait recovery, repeated cancellations, concurrent valid calls, and prompt context cancellation.
- [ ] `go test -count=1 -race ./...` passes; tests use httptest and no external network.

## Dependencies
Depends on completed bean [wingit-mcp-iw61](.beans/wingit-mcp-iw61--reuse-the-live-ebird-client-across-tool-calls.md), which introduced shared client pacing.

## Review provenance
Found during review after commit `31186b7` (`fix(ebird): reuse paced client across tool calls`).
