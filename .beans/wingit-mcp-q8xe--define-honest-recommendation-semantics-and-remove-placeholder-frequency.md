---
# wingit-mcp-q8xe
title: Define honest recommendation semantics and remove placeholder frequency
status: todo
type: task
priority: high
tags:
    - correctness
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:22Z
updated_at: "2026-10-27T20:00:37Z"
---
## Problem
Every candidate currently receives RecentFrequency = 0.20. Ranking, MinFrequency, and planned frequency bands imply statistics the server does not compute. Raw recent-species observations alone must not be assumed to provide a checklist-frequency denominator.

## Code references
- internal/tools/target_engine.go: constant frequency and ranking/filter behavior.
- internal/prompts/field_checklist.go: frequency grouping instructions.

## Acceptance criteria
- [ ] Document the supported statistical/ranking meaning and data requirements in an ADR under docs/adr/.
- [ ] Either compute a defensible measure from sufficient data or represent frequency as unavailable and rank by supported signals.
- [ ] Remove the fabricated 0.20 and define MinFrequency behavior when data is unavailable.
- [ ] Define whether repeated species observations are aggregated so targets and species limits have clear semantics.
- [ ] Align JSON/schema, prompt, documentation, and tests with the decision; call out compatibility effects.
- [ ] Tests cover ranking/filtering and unavailable-data cases without external network.

## Dependencies
None. Coordinate with the timestamp fix ([wingit-mcp-v92g](wingit-mcp-v92g--parse-live-observation-timestamps-for-recency-ranking.md)); output frequency bands ([wingit-mcp-rf67](wingit-mcp-rf67--add-honest-frequency-bands-and-align-field-checklist-output.md)) depend on this decision.

## Related issue
https://github.com/kpb/wingit-mcp/issues/5 (downstream output polish, not a duplicate of this semantics task).