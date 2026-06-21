# Code Standards

## Principles

- Keep runtime offline by default.
- Do not log raw secret values.
- Prefer standard library until a dependency removes real complexity.
- Keep file writes guarded by containment checks, backups, and tests.
- Use table-driven tests for detector and CLI behavior.

## Go

- Run `gofmt`.
- Keep packages narrow: `detect`, `scanner`, `redact`, `report`, `guardrails`,
  `cli`.
- Return structured errors. Do not panic outside embedded asset loading.
- Avoid global mutable state except build-time version variables.

## Security

- Test credentials must be generated dynamically or clearly fake.
- Reports must use masked previews and fingerprints only.
- Redaction must reject symlinks and files outside scan targets.

## Unresolved Questions

None.
