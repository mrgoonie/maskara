---
title: Automate Semantic Releases
description: >-
  Automate beta and stable GitHub releases from branch pushes using conventional
  commit messages.
status: completed
priority: P1
branch: codex/feat/semantic-release-automation
tags:
  - ci
  - release
  - automation
blockedBy: []
blocks: []
created: '2026-06-21T15:45:45.281Z'
createdBy: 'ck:plan'
source: skill
---

# Automate Semantic Releases

## Overview

Replace the current manual/tag-driven release flow with branch push automation.
Pushes to `dev` create beta prerelease tags and GitHub releases. Pushes to
`main` create stable official tags and GitHub releases. Version bump and release
notes come from conventional commit messages since the latest relevant tag.

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Release Workflow](./phase-01-release-workflow.md) | Completed |
| 2 | [Documentation](./phase-02-documentation.md) | Completed |
| 3 | [Verification](./phase-03-verification.md) | Completed |

## Dependencies

None.

## Acceptance Criteria

- [x] `dev` push runs tests, computes the next beta semver, creates a `vX.Y.Z-beta.N` tag, generates changelog from conventional commits, and publishes a prerelease.
- [x] `main` push runs tests, computes the next stable semver, creates a `vX.Y.Z` tag, generates changelog from conventional commits, and publishes an official release.
- [x] Version bump rules: breaking changes or `!`/`BREAKING CHANGE` -> major, `feat` -> minor, otherwise patch.
- [x] Beta numbering increments from existing beta tags for the next base version.
- [x] Release job is idempotent when the computed tag already exists.
- [x] Docs explain the automated release flow and no longer instruct manual tag creation as the default.
- [x] YAML and Go tests pass locally.

## Verification Evidence

- `.github/workflows/release.yml` parsed with Python `yaml.safe_load`.
- Release compute script block passed `bash -n`.
- `git diff --check` passed.
- `go test ./...` passed.
- `node %USERPROFILE%\.claude\scripts\validate-docs.cjs docs/` passed.
- `goreleaser check` not run locally because `goreleaser` is not installed in PATH; workflow installs it through `goreleaser/goreleaser-action@v6`.
