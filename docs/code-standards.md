# Code Standards

## Principles

- Keep runtime offline by default.
- Do not log raw secret values.
- Prefer standard library until a dependency removes real complexity.
- Keep file writes guarded by containment checks, backups, and tests.
- Use table-driven tests for detector and CLI behavior.

## Go

- Run `gofmt`.
- Keep packages narrow: `agents`, `detect`, `scanner`, `redact`, `report`,
  `guardrails`, `cli`.
- Keep the agent catalog in `internal/agents`; do not duplicate supported-agent
  lists in CLI, guardrails, scanner, or docs.
- Add table-driven tests when adding aliases, default roots, or guardrail
  support for an agent.
- Return structured errors. Do not panic outside embedded asset loading.
- Avoid global mutable state except build-time version variables.

## Security

- Test credentials must be generated dynamically or clearly fake.
- Reports must use masked previews and fingerprints only.
- Redaction must reject symlinks and files outside scan targets.
- Guardrail writes must stay under the configured root, use backups before
  replacement, and keep dry-run output side-effect-free.

## Unresolved Questions

None.
