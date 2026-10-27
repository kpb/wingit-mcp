---
# wingit-mcp-7zf5
title: v0.2.0 — trustworthy live data and MCP release readiness
status: todo
type: milestone
priority: high
tags:
    - release
created_at: 2026-10-27T19:58:08Z
updated_at: "2026-10-27T20:00:37Z"
---
## Outcome
Ship a usable birding assistant with reliable offline/live behavior, honest recommendation semantics, tested MCP integration, and reproducible release binaries.

## Release criteria
- [ ] Observation query arguments and reported filters agree.
- [ ] Live pacing, timestamps, and offline failures behave correctly.
- [ ] Ranking and frequency claims are supported by available data, not placeholder values.
- [ ] Source-selection and MCP integration tests run without external network access.
- [ ] Recent-sightings caching has tested hit/miss and expiration behavior.
- [ ] Output polish is consistent with the agreed frequency semantics.
- [ ] Cross-platform builds and setup/version documentation are ready.

## Scope and provenance
Reconciles GitHub issues #2–#6 with the implemented code. Live client issue #1 is already closed and is not duplicated. GitHub #7 (recent-sightings resource) remains optional outside this milestone.

Original milestone: https://github.com/kpb/wingit-mcp/milestone/1
The old March 2026 deadline is not carried forward; choose a new release date after correctness work is understood.

## Execution order
Correctness first, then reliability and caching, then output/release polish. Child beans carry acceptance criteria and explicit dependency links.

Beanstalk is the execution source of truth; GitHub issues remain unchanged as provenance.

## Work index

### Correctness (high priority)
- [wingit-mcp-dta3](wingit-mcp-dta3--normalize-query-arguments-before-fetching-observations.md): Normalize query arguments before fetching observations
- [wingit-mcp-iw61](wingit-mcp-iw61--reuse-the-live-ebird-client-across-tool-calls.md): Reuse the live eBird client across tool calls
- [wingit-mcp-v92g](wingit-mcp-v92g--parse-live-observation-timestamps-for-recency-ranking.md): Parse live observation timestamps for recency ranking
- [wingit-mcp-mqwt](wingit-mcp-mqwt--return-actionable-errors-for-failed-offline-observation-loads.md): Return actionable errors for failed offline observation loads
- [wingit-mcp-q8xe](wingit-mcp-q8xe--define-honest-recommendation-semantics-and-remove-placeholder-frequency.md): Define honest recommendation semantics and remove placeholder frequency

### Reliability
- [wingit-mcp-vxe1](wingit-mcp-vxe1--add-offline-mcp-tools-prompts-and-resources-integration-tests.md): Add offline MCP tools, prompts, and resources integration tests
- [wingit-mcp-18w6](wingit-mcp-18w6--complete-and-test-observation-source-selection.md): Complete and test observation source selection
- [wingit-mcp-duib](wingit-mcp-duib--cache-normalized-recent-sightings-queries-with-ttl.md): Cache normalized recent-sightings queries with TTL

### Output and release readiness
- [wingit-mcp-rf67](wingit-mcp-rf67--add-honest-frequency-bands-and-align-field-checklist-output.md): Add honest frequency bands and align field-checklist output
- [wingit-mcp-t1hi](wingit-mcp-t1hi--add-reproducible-cross-platform-release-builds.md): Add reproducible cross-platform release builds
- [wingit-mcp-3wvj](wingit-mcp-3wvj--refresh-setup-documentation-and-release-version-metadata.md): Refresh setup documentation and release version metadata

### Optional, outside this milestone
- [wingit-mcp-f5dg](wingit-mcp-f5dg--expose-a-read-only-recent-sightings-mcp-resource.md): Expose a read-only recent-sightings MCP resource
