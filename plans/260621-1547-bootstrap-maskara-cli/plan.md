---
title: Bootstrap Maskara CLI
description: ''
status: completed
priority: P1
branch: feat/bootstrap-maskara-cli
tags:
  - cli
  - security
  - guardrails
  - release
blockedBy: []
blocks: []
created: '2026-06-21T08:47:18.153Z'
createdBy: 'ck:plan'
source: skill
---

# Bootstrap Maskara CLI

## Overview

Bootstrap `maskara`, an offline open-source Go CLI that scans local coding-agent
session logs for leaked secrets, reports findings, redacts safely, and installs
agent guardrails to reduce future exposure.

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Scaffold CLI](./phase-01-scaffold-cli.md) | Completed |
| 2 | [Detection and Redaction](./phase-02-detection-and-redaction.md) | Completed |
| 3 | [Guardrails](./phase-03-guardrails.md) | Completed |
| 4 | [Docs and Release](./phase-04-docs-and-release.md) | Completed |
| 5 | [Verification](./phase-05-verification.md) | Completed |

## Dependencies

None.

## Acceptance Criteria

- [x] `maskara`, `maskara scan`, `maskara report`, and `maskara guardrails` work on Windows, macOS, and Linux paths.
- [x] Scanner detects common API keys, tokens, private keys, JWTs, database URLs, and env assignments without network calls.
- [x] Reports default to Markdown and support JSON with `--json`; `--output`/`-o` writes to a file or directory.
- [x] Redaction is explicit, creates backups, validates rewritten files, and preserves readable session logs.
- [x] Guardrails can install Claude/Codex/OpenCode-style instructions and hooks without overwriting user files silently.
- [x] README includes generated 5:2 watercolor technical sketch banner.
- [x] CI supports stable releases from `main` and beta/prerelease builds from `dev`.
- [x] Unit tests and local build pass.

## Verification Evidence

- `go test ./...` passed.
- `go build ./cmd/maskara` passed.
- Smoke test verified scan/report/full-redact/guardrails dry-run behavior.
- `go run ./cmd/maskara scan --root .` found 0 findings in the repo.
- `goreleaser check` passed after GitHub remote setup.

## Scope Boundary

MVP supports filesystem-backed text/session logs first. SQLite-backed agent
stores, live provider validation, and package-manager distribution are deferred.

## References

- https://github.com/Ishannaik/agent-sweep
- https://github.com/simonw/scan-for-secrets
- https://github.com/GitGuardian/ggshield
