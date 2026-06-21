---
phase: 2
title: Detection and Redaction
status: completed
priority: P1
dependencies:
  - 1
effort: large
---

# Phase 2: Detection and Redaction

## Overview

Implement secret detection, report models, file scanning, and guarded redaction
with backups.

## Requirements

- Functional: detect common secret categories in text, JSON, JSONL, markdown,
  and env-like logs.
- Non-functional: do not print raw secret values by default; avoid corrupting
  logs; preserve enough context for rotation decisions.

## Architecture

Rules are centralized in `internal/detect`. Scanner streams files and returns
findings with masked previews and byte/line ranges. Redactor rewrites matched
values, creates `.maskara.bak`, and validates JSON/JSONL when applicable.

## Implementation Steps

1. Define finding/report structs with agent, file, line, rule, severity, and masked preview.
2. Implement regex rules for API keys, OAuth tokens, JWTs, private keys, DB URLs, cloud tokens, and env assignments.
3. Add filesystem scanner with binary-file skipping and extension allowlist.
4. Add report renderers for Markdown and JSON.
5. Add redaction engine with backup, containment check, symlink rejection, and validation.
6. Cover scanner and redactor with fixtures.

## Success Criteria

- [x] Scan returns exit code `1` when findings exist, `0` when clean, `2` on errors.
- [x] Markdown and JSON reports include masked previews, not full secrets.
- [x] Redaction replaces secret bytes with `[MASKARA_REDACTED:<rule>]` and creates backup.
- [x] JSONL fixtures remain parseable after redaction.
