---
phase: 3
title: Docs and Verification
status: completed
priority: P2
dependencies:
  - 1
  - 2
---

# Phase 3: Docs and Verification

## Overview

Update user-facing and maintainer documentation, then verify the expanded
catalog through focused tests and CLI smoke checks.

## Requirements

- Functional: README agent table lists all supported agents with concise default location descriptions.
- Functional: help text exposes the full accepted agent set without hiding `all`.
- Non-functional: docs remain honest about text/session-log scanning and do not imply provider validation or SQLite parsing.

## Architecture

Docs stay in existing project locations. No new docs directory structure is
needed beyond this plan.

## Related Code Files

- Modify: `README.md`
- Modify: `docs/project-overview-pdr.md`
- Modify: `docs/project-roadmap.md`
- Modify as needed: `docs/codebase-summary.md`, `docs/system-architecture.md`

## Implementation Steps

1. Update README usage/help examples and agent support table.
2. Update docs that name the supported agent set.
3. Run `gofmt`.
4. Run focused package tests, then `go test ./...`.
5. Smoke-test CLI help and scan/guardrail behavior for representative new agents.

## Success Criteria

- [ ] README and docs match implemented catalog.
- [ ] `go test ./...` passes.
- [ ] `go run ./cmd/maskara scan --agent gemini --root <tmp>` accepts the agent.
- [ ] `go run ./cmd/maskara guardrails -a qwen --root <tmp> --dry-run` returns planned changes.

## Risk Assessment

Help text can become noisy with many names. Mitigation: present a compact agent
name section instead of repeating the full list inline everywhere.
