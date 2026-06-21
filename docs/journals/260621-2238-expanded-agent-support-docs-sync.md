# Expanded Agent Support Docs Sync

## Context

Expanded agent support shipped to `dev`, but docs sync was incomplete after the
beta flow because beta shipping skipped the full docs update step.

## What Happened

- Added a dedicated agent support reference.
- Added project changelog entry for the expanded catalog.
- Refreshed code standards, architecture, deployment, roadmap, PDR, codebase
  summary, and README links.
- Clarified `--agent all` behavior with and without `--root`.

## Decisions

- Keep agent metadata source of truth in `internal/agents`.
- Document generic guardrails honestly: local files are installed, but automatic
  loading is not guaranteed for every agent.

## Next

- For future agent additions, update code, tests, README, agent support
  reference, roadmap, and changelog in the same PR.

## Unresolved Questions

None.
