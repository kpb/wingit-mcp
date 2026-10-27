---
# wingit-mcp-3wvj
title: Refresh setup documentation and release version metadata
status: todo
type: task
priority: normal
tags:
    - release
    - documentation
parent: wingit-mcp-7zf5
created_at: 2026-10-27T20:00:12Z
updated_at: "2026-10-27T20:00:37Z"
---
## Goal
Make documentation describe the shipped behavior and a reproducible release process, not the old partially completed roadmap. Do not publish or tag a release in this task.

## Code references
- README.md: setup, modes, and outdated roadmap.
- cmd/wingit-mcp/main.go: MCP implementation version currently 0.1.0.
- internal/ebird/client.go: User-Agent currently advertises 0.2.0.
- Makefile: build/test/lint and future release target.

## Acceptance criteria
- [ ] Document build/start commands, the expected normalized personal JSON format with sample data, and a host configuration example.
- [ ] Explain offline precedence, live coordinate requirements, missing-source errors, and cache settings/TTL.
- [ ] Explain actual ranking/frequency semantics and any unavailable-data states without unsupported likelihood claims.
- [ ] Replace stale planned/completed claims with current release scope and Beanstalk tracking guidance.
- [ ] Make MCP and HTTP version metadata consistent with the intended release; avoid a premature tag/publication.
- [ ] Document cross-platform artifact names and release verification commands.
- [ ] Verify examples with synthetic offline fixtures; run the full race-enabled test suite and lint, reporting unavailable checks.

## Dependencies
- Blocked by [wingit-mcp-18w6](wingit-mcp-18w6--complete-and-test-observation-source-selection.md): final observation-source contract.
- Blocked by [wingit-mcp-duib](wingit-mcp-duib--cache-normalized-recent-sightings-queries-with-ttl.md): cache configuration and TTL.
- Blocked by [wingit-mcp-rf67](wingit-mcp-rf67--add-honest-frequency-bands-and-align-field-checklist-output.md): final output/prompt semantics.
- Blocked by [wingit-mcp-t1hi](wingit-mcp-t1hi--add-reproducible-cross-platform-release-builds.md): release build commands.
- Blocked by [wingit-mcp-vxe1](wingit-mcp-vxe1--add-offline-mcp-tools-prompts-and-resources-integration-tests.md): MCP integration coverage.
- Blocked by [wingit-mcp-v92g](wingit-mcp-v92g--parse-live-observation-timestamps-for-recency-ranking.md): live recency ranking fix.
Other correctness fixes are prerequisites through these beans.