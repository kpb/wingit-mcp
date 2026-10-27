---
# wingit-mcp-vxe1
title: Add offline MCP tools, prompts, and resources integration tests
status: todo
type: task
priority: high
tags:
    - reliability
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:53Z
updated_at: 2026-10-27T19:59:53Z
---
## Source
https://github.com/kpb/wingit-mcp/issues/4

## Current state and scope
The existing fixture integration test calls the pure engine, not an MCP session. Add coverage of the actual registered server surface.

## Code references
- cmd/wingit-mcp/main.go: registration and server construction.
- internal/prompts/register.go
- internal/mcp/resources.go
- internal/tools/target_integration_test.go: existing engine-only integration.

## Acceptance criteria
- [ ] Construct the production registrations in-process through a testable server factory or equivalent shared wiring.
- [ ] tools/list includes target_checklist; prompts/list includes field_checklist; resources/list includes wingit://personal-checklist.
- [ ] Invoke the tool, retrieve the prompt, and read the resource using fixture-backed MCP sessions.
- [ ] Check structured tool output/schema and the prompt assumptions agree, including field names.
- [ ] Assert meaningful resource JSON and exercise an error response.
- [ ] Tests need neither external network nor real personal data and pass with the race detector.

## Dependencies
None. Establish regression coverage early; update expectations as correctness and schema fixes land.
