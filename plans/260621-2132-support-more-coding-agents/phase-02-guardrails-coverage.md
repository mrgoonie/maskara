---
phase: 2
title: Guardrails Coverage
status: completed
priority: P1
dependencies:
  - 1
---

# Phase 2: Guardrails Coverage

## Overview

Expand `maskara guardrails` so every supported agent can receive a safe local
guardrail plan. Keep native instruction-file handling for known agents and use
a generic fallback for newer agents where native instruction loading is not
established in this codebase.

## Requirements

- Functional: `guardrails -a <requested-agent> --root <tmp> --dry-run` returns planned files instead of unsupported-agent errors.
- Functional: auto/all detection checks the expanded catalog through reusable agent root markers.
- Non-functional: retain containment checks, backups, and dry-run behavior.

## Architecture

`internal/guardrails` consumes agent catalog metadata instead of duplicating
supported-agent decisions. Known native targets keep existing file paths:
Claude `CLAUDE.md`, Codex `AGENTS.md`, OpenCode config, and Antigravity config.
Other agents install generic guardrail markdown and hook files under their
primary local config root.

## Related Code Files

- Modify: `internal/guardrails/guardrails.go`
- Modify: `internal/guardrails/files.go`
- Modify: `internal/guardrails/guardrails_test.go`

## Implementation Steps

1. Replace duplicated guardrail auto-detection markers with agent catalog roots.
2. Add fallback install plans for supported agents that lack native paths.
3. Preserve existing exact paths and backup semantics for current agents.
4. Add table-driven tests covering dry-run plans for the expanded catalog.

## Success Criteria

- [ ] Existing Codex guardrail backup test still passes.
- [ ] Dry-run install plans every requested agent.
- [ ] Unsupported guardrail agent still errors clearly.
- [ ] File writes remain contained under the configured root.

## Risk Assessment

Generic files may not be automatically read by every agent. Mitigation: docs
call this a local guardrail/reference install where native integration is not
known, not guaranteed runtime enforcement.
