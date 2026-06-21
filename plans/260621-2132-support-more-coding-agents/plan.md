---
title: Support More Coding Agents
description: >-
  Expand Maskara scan and guardrail support across the requested coding-agent
  catalog.
status: completed
priority: P2
branch: codex/feat/support-more-coding-agents
tags:
  - feature
  - cli
  - security
blockedBy: []
blocks: []
created: '2026-06-21T14:32:54.387Z'
createdBy: 'ck:plan'
source: skill
---

# Support More Coding Agents

## Overview

Expand Maskara's supported agent catalog beyond Claude/Codex/OpenCode/Antigravity.
The CLI should normalize common names and aliases, resolve likely local config
or session roots, install generic guardrail files where native instruction paths
are not known, and document the broadened support without claiming remote or
provider-side coverage.

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Agent Registry](./phase-01-agent-registry.md) | Completed |
| 2 | [Guardrails Coverage](./phase-02-guardrails-coverage.md) | Completed |
| 3 | [Docs and Verification](./phase-03-docs-and-verification.md) | Completed |

## Dependencies

None.

## Acceptance Criteria

- [x] `maskara scan --agent <name> --root <tmp>` accepts each requested agent and labels findings with the normalized agent name.
- [x] `maskara scan --agent all` and `auto` include the expanded supported catalog when discovering existing roots.
- [x] `maskara guardrails -a <name> --root <tmp> --dry-run` plans files for each supported agent without unsupported-agent errors.
- [x] CLI help and README list the expanded agent catalog concisely.
- [x] Maintainer docs note expanded scan/guardrail scope and current text-file limitation.
- [x] `gofmt`, `go test ./...`, and focused CLI smoke checks pass.

## Verification Evidence

- `go test ./...` passed.
- `go run ./cmd/maskara scan --agent gemini --root <tmp>` passed with 0 findings.
- `go run ./cmd/maskara guardrails -a qwen --root <tmp> --dry-run` planned generic guardrail files.
- `go run ./cmd/maskara --help` listed the expanded canonical agent catalog.

## Validation Log

### Verification Results

- Tier: Standard
- Claims checked: 8
- Verified: 8
- Failed: 0
- Unverified: 0
- Evidence: `internal/agents/agents.go`, `internal/guardrails/guardrails.go`, `internal/cli/cli.go`, `internal/cli/help.go`, `README.md`, `docs/project-roadmap.md`.

### Red Team Review

- Security adversary: accepted. Do not add network validation or provider API calls; keep all behavior local and file-backed.
- Assumption destroyer: accepted. For agents without verified native instruction files, install a generic `maskara-guardrails.md` under likely config roots instead of claiming native enforcement.
- Failure mode analyst: accepted. Keep `--root` behavior backward compatible and add tests so custom roots still bypass host-specific path discovery.

### Whole-Plan Consistency Sweep

- Files reread: plan.md, phase-01-agent-registry.md, phase-02-guardrails-coverage.md, phase-03-docs-and-verification.md.
- Decision deltas checked: expanded catalog, generic guardrail fallback, docs scope wording.
- Reconciled stale references: 0.
- Unresolved contradictions: 0.
