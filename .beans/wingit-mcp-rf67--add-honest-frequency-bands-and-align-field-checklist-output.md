---
# wingit-mcp-rf67
title: Add honest frequency bands and align field-checklist output
status: todo
type: feature
priority: normal
tags:
    - output
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:53Z
updated_at: "2026-10-27T20:00:37Z"
---
## Source
https://github.com/kpb/wingit-mcp/issues/5

## Scope
Improve field output after recommendation semantics are settled. Do not turn the placeholder constant into a universal high-frequency label.

## Code references
- internal/tools/target_engine.go: TargetRow and recommendation result.
- internal/prompts/field_checklist.go: grouping and field assumptions.

## Acceptance criteria
- [ ] Follow the frequency/ranking ADR from [wingit-mcp-q8xe](wingit-mcp-q8xe--define-honest-recommendation-semantics-and-remove-placeholder-frequency.md).
- [ ] For supported measured frequency, expose high/medium/low bands using documented thresholds and boundary tests.
- [ ] If measured frequency is unavailable, explicitly represent that state instead of inventing a band; document any resulting scope change to GitHub #5 in this bean.
- [ ] Prompt uses the structured band/state rather than independently deriving thresholds.
- [ ] JSON field names, schema, prompt, and examples agree; only promise location details actually present in output.
- [ ] Test measured boundary cases where applicable and the unavailable-data/empty-target cases.

## Dependencies
- Blocked by [wingit-mcp-q8xe](wingit-mcp-q8xe--define-honest-recommendation-semantics-and-remove-placeholder-frequency.md): honest recommendation/frequency semantics.