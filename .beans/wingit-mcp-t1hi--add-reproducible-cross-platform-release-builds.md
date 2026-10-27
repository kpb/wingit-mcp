---
# wingit-mcp-t1hi
title: Add reproducible cross-platform release builds
status: todo
type: task
priority: normal
tags:
    - release
parent: wingit-mcp-7zf5
created_at: 2026-10-27T19:59:22Z
updated_at: 2026-10-27T19:59:22Z
---
## Source
https://github.com/kpb/wingit-mcp/issues/6

## Current state and scope
Makefile only builds the local binary. Add a release-build target for linux/amd64 and darwin/arm64; leave extra platforms and automated publishing out of scope.

## Code references
- Makefile: build, test, lint, and clean targets.
- .gitignore: dist/ is already ignored.

## Acceptance criteria
- [ ] make release-build produces dist/wingit-mcp_linux_amd64 and dist/wingit-mcp_darwin_arm64.
- [ ] Use explicit GOOS/GOARCH and CGO_ENABLED=0 with reproducible build settings.
- [ ] Existing local build, test, and lint targets are unchanged.
- [ ] Verify both binaries build and document the invocation.
- [ ] No binaries are committed and no release is published as part of this task.

## Dependencies
None; independent of recommendation correctness, but not a substitute for release verification.
