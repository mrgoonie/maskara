# Codebase Summary

## Layout

| Path | Purpose |
|---|---|
| `cmd/maskara` | CLI entrypoint |
| `internal/cli` | command parsing and workflow routing |
| `internal/agents` | agent catalog, name normalization, and default path discovery |
| `internal/detect` | detector rules and masked finding creation |
| `internal/scanner` | filesystem traversal and file scanning |
| `internal/redact` | backup and redaction writes |
| `internal/report` | Markdown and JSON report rendering |
| `internal/guardrails` | guardrail install planner and embedded assets |
| `docs/agent-support-reference.md` | supported agent names, aliases, roots, and guardrail behavior |
| `docs` | maintainer documentation |
| `plans` | implementation plan and phase records |

## Runtime Flow

`cli` resolves targets through `agents`, scans via `scanner`, detects through
`detect`, optionally redacts through `redact`, then writes output through
`report`.

## Agent Catalog

`internal/agents` owns the canonical supported-agent list. CLI help reads from
that catalog, scanner target resolution uses it for `auto` and `all`, and
guardrails reuse its root markers so support does not drift across packages.

## Unresolved Questions

None.
