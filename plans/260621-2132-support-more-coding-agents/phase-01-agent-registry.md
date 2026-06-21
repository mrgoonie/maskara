---
phase: 1
title: Agent Registry
status: completed
priority: P1
dependencies: []
---

# Phase 1: Agent Registry

## Overview

Create a maintainable agent catalog for scan target resolution and alias
normalization. Add the requested agents while preserving existing names,
defaults, and custom `--root` behavior.

## Requirements

- Functional: support `cursor`, `antigravity`, `kimi`, `droid`, `gemini`, `github-copilot`, `hermes`, `openclaw`, `kilo`, `kiro`, `pi`, `qoder`, `qwen`, and `trae`.
- Functional: normalize common CLI/name variants such as `gemini-cli`, `qwen-code`, `kimi-code`, `github copilot`, and `antigravity-cli`.
- Functional: preserve existing aliases `claude-code`, `claudecode`, and `open-code`.
- Non-functional: keep scan resolution offline and standard-library only.

## Architecture

`internal/agents` remains the source of truth. Replace one-off switch logic
with a small catalog mapping canonical agent names to aliases and path resolver
functions. Expose helper data for CLI help and guardrail root detection.

## Related Code Files

- Modify: `internal/agents/agents.go`
- Create/modify: `internal/agents/*_test.go`
- Modify: `internal/cli/cli.go`
- Modify: `internal/cli/help.go`

## Implementation Steps

1. Define the canonical supported agent list and alias map.
2. Add default path resolvers using home dotdirs plus platform config/data dirs.
3. Keep `Resolve(agent, root)` behavior: custom root returns one explicit target and uses `custom` for `auto`.
4. Add agent package tests for aliases, support list coverage, and custom-root behavior.
5. Route CLI help/flag descriptions through the catalog helper where practical.

## Success Criteria

- [ ] Every requested agent normalizes to a supported canonical name.
- [ ] Unsupported agents still fail when no custom root is supplied.
- [ ] Custom-root scans remain backward compatible.
- [ ] Agent package tests pass.

## Risk Assessment

Default directories for newer tools can drift. Mitigation: scan multiple likely
local roots per agent and document that `--root` remains the precise override.
