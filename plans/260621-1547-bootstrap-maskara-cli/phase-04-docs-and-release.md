---
phase: 4
title: Docs and Release
status: completed
priority: P2
dependencies:
  - 1
  - 2
  - 3
effort: medium
---

# Phase 4: Docs and Release

## Overview

Document the project and configure CI/release for an open-source repository.

## Requirements

- Functional: README covers install, usage, safety model, agent support, and
  examples from the prompt.
- Non-functional: release process separates stable `main` and beta `dev`.

## Architecture

GitHub Actions run Go tests on PR/push. GoReleaser builds archives for
Windows/macOS/Linux. `main` tags publish stable releases; `dev` produces beta
pre-releases.

## Implementation Steps

1. Generate and save 5:2 watercolor technical sketch banner under `assets/branding/`.
2. Write README with banner, badges, quick start, commands, threat model, and safety notes.
3. Add docs: project overview, code standards, architecture, roadmap, release guide.
4. Add `.github/workflows/ci.yml`, `.github/workflows/release.yml`, and `.goreleaser.yaml`.
5. Add contribution and security docs.

## Success Criteria

- [x] README embeds local banner path.
- [x] CI covers Windows, macOS, Linux Go test/build.
- [x] Release workflow differentiates stable and beta.
