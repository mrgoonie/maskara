---
phase: 1
title: Scaffold CLI
status: completed
priority: P1
dependencies: []
effort: medium
---

# Phase 1: Scaffold CLI

## Overview

Create the Go module, CLI command structure, repository metadata, and minimal
offline data flow for `maskara`.

## Requirements

- Functional: root command runs the full workflow; subcommands exist for
  `scan`, `report`, and `guardrails`.
- Non-functional: no network calls at runtime, cross-platform path handling,
  deterministic tests, simple module boundaries.

## Architecture

`cmd/maskara` owns process entry. `internal/cli` parses commands and routes to
domain packages. `internal/agents` discovers log roots. `internal/scanner`,
`internal/report`, `internal/redact`, and `internal/guardrails` own behavior.

## Implementation Steps

1. Initialize `go.mod` for `github.com/mrgoonie/maskara`.
2. Add `cmd/maskara/main.go` and internal package layout.
3. Add option parsing with standard `flag` package to keep dependencies lean.
4. Add root workflow equivalent to scan, redact, report.
5. Add `.gitignore`, license, and initial test fixtures.

## Success Criteria

- [x] `go test ./...` can discover all packages.
- [x] `go run ./cmd/maskara --help` prints usable command help.
- [x] `go run ./cmd/maskara scan --agent codex --root <fixture>` runs read-only.
