---
# wingit-mcp-v92g
title: Parse live observation timestamps for recency ranking
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
The engine parses only YYYY-MM-DD, but the live client preserves timestamps such as 2017-08-23 20:05. These silently become zero time and defeat recency ranking.

## Code references
- internal/tools/target_engine.go: ObsDt parsing and sort.
- internal/ebird/client_test.go: live-format response fixture.

## Acceptance criteria
- [ ] Support the documented eBird timestamp format and existing date-only fixtures.
- [ ] Define date-only/timezone interpretation without depending on the machine local timezone.
- [ ] Missing or malformed dates have deterministic documented ordering.
- [ ] Tests verify newest-first ranking with mixed date formats and stable ties.
- [ ] Preserve the original observation timestamp in output.

## Dependencies
None.
