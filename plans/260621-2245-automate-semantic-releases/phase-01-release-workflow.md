---
phase: 1
title: Release Workflow
status: completed
priority: P1
dependencies: []
---

# Phase 1: Release Workflow

## Overview

Update `.github/workflows/release.yml` so branch pushes own release creation.
Use shell logic plus GoReleaser already present in the repo; avoid new package
dependencies.

## Requirements

- Functional: `dev` push produces prerelease beta tags.
- Functional: `main` push produces stable tags.
- Functional: release notes derive from conventional commit subjects/bodies.
- Non-functional: no secrets beyond `GITHUB_TOKEN`; full git history available.

## Architecture

Workflow checks out full history, fetches tags, runs tests, computes release
metadata, creates a local annotated tag, pushes it, and runs GoReleaser with the
generated changelog. The workflow runs only on branch pushes and manual dispatch
to avoid recursive tag-trigger releases.

## Implementation Steps

1. Change release trigger to branch pushes for `dev` and `main`, plus manual dispatch.
2. Add a concurrency group per branch.
3. Compute latest stable tag and next semver bump from conventional commits.
4. For `dev`, compute next beta suffix for that base version.
5. Generate release notes grouped by breaking, features, fixes, and other changes.
6. Create and push the computed tag only if it does not exist.
7. Run GoReleaser with `--release-notes`.

## Success Criteria

- [ ] Workflow syntax is valid YAML.
- [ ] Version script handles no existing tags.
- [ ] GoReleaser gets `GITHUB_TOKEN` and a concrete tag checkout context.
- [ ] Manual dispatch can test either `dev` or `main`.
