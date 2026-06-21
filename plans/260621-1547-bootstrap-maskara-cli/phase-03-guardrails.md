---
phase: 3
title: Guardrails
status: completed
priority: P1
dependencies:
  - 1
  - 2
effort: medium
---

# Phase 3: Guardrails

## Overview

Create `maskara guardrails` assets for coding agents so future sessions avoid
printing or persisting sensitive values.

## Requirements

- Functional: support auto-detect plus `--agent claude|codex|opencode`; write
  agent instructions, skill docs, and hook scripts where applicable.
- Non-functional: never overwrite user files without backup; output exact file
  paths changed; generated hooks avoid secret value exfiltration.

## Architecture

`internal/guardrails` exposes install plans and writes files through a safe
writer. Assets are embedded with Go `embed` from `internal/guardrails/assets`.

## Implementation Steps

1. Add reusable guardrail instruction text for coding agents.
2. Add Claude/Codex target path detection for Windows, macOS, and Linux.
3. Add hook scripts that block obvious secret-printing commands and privacy-sensitive file reads.
4. Add `--dry-run` and backup-on-overwrite behavior.
5. Add tests for generated install plans.

## Success Criteria

- [x] `maskara guardrails --dry-run` shows planned files without writing.
- [x] `maskara guardrails -a codex --root <tmp>` installs Codex guardrails in tests.
- [x] Existing files are backed up with `.maskara.bak` before replacement.
