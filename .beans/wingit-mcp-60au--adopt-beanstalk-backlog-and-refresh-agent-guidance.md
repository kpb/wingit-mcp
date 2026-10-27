---
# wingit-mcp-60au
title: Adopt Beanstalk backlog and refresh agent guidance
status: completed
type: task
priority: high
tags:
    - workflow
created_at: 2026-10-27T19:58:08Z
updated_at: "2026-10-27T20:01:05Z"
---
## Goal
Make Beanstalk the execution backlog and document the project working agreements.

## Acceptance criteria
- [x] Create a prioritized v0.2.0 backlog with acceptance criteria, dependency links, and references to existing GitHub issues.
- [x] Keep the recent-sightings resource optional and outside the release-critical milestone.
- [x] Update AGENTS.md with project layout, commands, Beanstalk workflow, MCP stdio safety, data privacy, ADRs, and Conventional Commits.
- [x] Preserve the existing Beanstalk configuration; leave GitHub issues and application code unchanged.
- [x] Verify the backlog structure and agent instructions.

## Outcome and verification
Created milestone wingit-mcp-7zf5 with 11 prioritized release work items, plus optional feature wingit-mcp-f5dg outside the milestone. Migrated the intent of GitHub #2–#7 with original URLs; did not duplicate completed live-client work. Dependencies are clickable Markdown links in bean bodies because this CLI has no dependency flags.

Expanded AGENTS.md with project structure, task workflow, commands, MCP/data safety, numbered ADR guidance under docs/adr/, and Conventional Commits including breaking-change notation. Preserved .beans.yml and .beans/.gitignore; left application code and GitHub unchanged.

Validated all 14 beans, milestone membership, todo statuses for future work, issue references, existing dependency-link targets, acyclic prerequisite relationships, required AGENTS.md sections, and whitespace. Beanstalk reports 11 todo milestone children. No application tests were rerun for this metadata/documentation-only change; race-enabled tests and go vet passed during the preceding stocktake. No commit or push was made.
