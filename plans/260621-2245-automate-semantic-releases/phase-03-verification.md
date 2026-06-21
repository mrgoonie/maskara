---
phase: 3
title: Verification
status: completed
priority: P1
dependencies:
  - 1
  - 2
---

# Phase 3: Verification

## Overview

Validate workflow syntax, local tests, docs consistency, and generated release
metadata logic as much as possible without pushing test tags.

## Implementation Steps

1. Run YAML parsing check for workflow files.
2. Run `go test ./...`.
3. Run docs validator.
4. Use shell dry-run snippets or manual inspection to verify version/changelog logic.

## Success Criteria

- [ ] `.github/workflows/release.yml` parses as YAML.
- [ ] `go test ./...` passes.
- [ ] Docs validator passes.
- [ ] No generated binary or release artifacts are staged.
