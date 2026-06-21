# Project Overview PDR

## Problem

Coding agents persist conversations, tool output, and terminal logs. Secrets can
land in those stores when users or agents print `.env` files, API keys, database
URLs, private keys, cookies, or provider tokens.

## Product

Maskara is a local CLI binary that scans agent session logs, reports detected
secret exposure, redacts local files with backups, and installs guardrails to
reduce future leaks.

## Users

- Developers using Claude Code, Codex, Cursor, OpenCode, Antigravity, Gemini
  CLI, Qwen Code, and similar local coding agents
- Open-source maintainers cleaning local agent transcripts before sharing repros
- Security-conscious teams adding local privacy hygiene to coding workflows

## MVP Requirements

- Cross-platform Go binary
- Offline scanning
- Markdown and JSON reports
- Safe local redaction with backups
- Guardrails installer for common coding agents, with generic local guardrail
  files where native instruction paths are not known
- CI and release workflows for stable and beta branches

## Out Of Scope

- Provider-side secret validation
- Automatic credential rotation
- Scanning remote SaaS logs
- SQLite/browser profile parsing
- Organization policy management

## Unresolved Questions

None.
