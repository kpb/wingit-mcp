# Agent Instructions

IMPORTANT: Before you do anything else, run `beanstalk prime` and follow its output.

## Project context

WingIt-MCP is a Go stdio MCP server that compares personal eBird history with
recent observations to suggest species the user has not seen. It supports
offline JSON fixtures and live eBird data. See `README.md` for configuration.

- `cmd/wingit-mcp/`: server startup, configuration, and tool wiring.
- `internal/ebird/`: checklist loaders and live HTTP client.
- `internal/tools/`: target-checklist engine and argument normalization.
- `internal/prompts/`: field-checklist prompt and registration.
- `internal/mcp/`: MCP resource registration.
- `internal/types/`: shared eBird data types.
- `data/` and package-local `testdata/`: sample data and test fixtures.

Keep domain logic testable independently of MCP transport and HTTP. Match the
existing Go structure and style; avoid unrelated refactors.

## Task tracking

Beanstalk is the execution source of truth. Tasks live in `.beans/`, with project
configuration in `.beans.yml`. Preserve this setup rather than reinitializing it.

- Follow `beanstalk prime` for task creation, claiming, and completion rules.
- Do not create beans for questions, exploration, analysis, or planning alone.
- Before untracked implementation, ask whether the user wants a bean unless
  they have already explicitly requested task creation or tracking.
- Create requested work as `todo`; read and claim its bean before implementation.
- Find work with `beanstalk list --status todo --json` and inspect details with
  `beanstalk show <id> --json`.
- Read dependency links in bean bodies before claiming work. A `todo` status
  alone does not mean a task is unblocked.
- Include scope, testable acceptance criteria, relevant code references, and
  dependencies. Preserve original GitHub issue URLs for migrated work.
- Keep outcomes and verification notes in the bean. Mark work `completed` only
  when its acceptance criteria are satisfied; record why work is `scrapped`.
- GitHub issues are provenance/intake, not a second execution backlog. Do not
  close or edit them without explicit user approval.

## Build and verification

- `make build`: build the stdio server into `bin/`.
- `go test -count=1 -race -cover ./...`: tests without repository coverage files.
- `make test`: race-enabled tests plus `coverage.out` and a coverage summary.
- `make lint`: `go vet` and `staticcheck` (requires an installed `staticcheck`).
- `make fmt`: format Go code. Run it only for intended changes and inspect the diff.

Run relevant tests for changes and the full suite when practical. Report any
checks that could not run. Keep generated artifacts out of commits. Do not run
`go mod tidy` or upgrade dependencies unless needed for the requested work.

## MCP and data safety

- Stdout is exclusively for the MCP protocol; send logs and diagnostics to stderr.
- Keep tool, prompt, resource, and JSON-schema contracts aligned with host-facing
  documentation. Test protocol wiring as well as pure domain behavior.
- `WINGIT_PERSONAL_JSON` supplies personal history. `WINGIT_RECENT_JSON` takes
  precedence over live access via `WINGIT_EBIRD_TOKEN`.
- Keep automated tests offline and deterministic: use fixtures and `httptest`,
  not real eBird credentials or network access. A live check requires user approval.
- Never commit tokens, `.env` files, or actual personal exports. Use synthetic
  sample data and avoid logging secrets or personal checklist contents.
- Do not present placeholder statistics as measured frequency or likelihood.
- Keep configuration, runtime behavior, README instructions, and release/version
  metadata consistent when changing them.

## Architecture decision records (ADRs)

Record significant architectural decisions in `docs/adr/`; create the directory
when the first ADR is needed. Routine fixes do not require an ADR.

- Use sequential filenames such as `0001-observation-source-and-cache.md`.
- Include a title, date, status, context, decision, and consequences/tradeoffs.
- Capture consequential choices about data semantics, persistence/caching, MCP
  contracts, or major dependencies when they arise; do not invent historical ADRs.
- Reference relevant ADRs from beans and link implementation work back to them.
- Preserve decision history: supersede an accepted ADR with a linked new record
  rather than silently rewriting the original decision.

## Conventional Commits

When making commits, use `type(scope): description`, with an optional scope.
Examples: `fix(ebird): reuse client across tool calls` and
`docs(agents): document task and architecture workflows`.

- Use an appropriate type: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`,
  `build`, `ci`, or `chore`.
- Keep the description concise and imperative; explain motivation in the body
  where useful and reference relevant bean IDs and GitHub issues.
- Mark breaking changes with `!` and include a `BREAKING CHANGE:` explanation
  in the commit body/footer.
- Keep commits focused. Inspect the diff and exclude unrelated user changes.
- Do not create commits, push, or publish releases unless the user requests it.
