# JSONL Redaction Escape Fix

## Context

Maskara could stop full redaction with `redaction would break JSONL` on valid
Codex session logs. The failing shape was a detected bounded value, such as a
database URL, ending immediately before an escaped quote inside a JSON string.

## What Changed

- Added escaped quote delimiter trimming after detector range selection.
- Preserved the broad detector regexes so plain-text secrets containing
  backslashes still match.
- Added detector tests for escaped quote boundaries and backslash-containing env
  values.
- Added a redaction integration test proving JSONL remains parseable after the
  escaped delimiter case is redacted.

## Decisions

- Keep structured JSON/JSONL validation as the final safety gate.
- Fix the detector range source instead of weakening redaction validation.
- Prefer targeted delimiter trimming over excluding backslashes from detector
  character classes.

## Verification

- `go test ./internal/detect ./internal/redact`
- `go test ./...`
- Manual repro with a valid JSONL Codex-style line containing an escaped
  database URL quote.

## Unresolved Questions

None.
