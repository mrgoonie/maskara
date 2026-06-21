# Project Changelog

## 2026-06-21

### Added

- Automated branch-push release workflow: `dev` publishes beta prereleases and
  `main` publishes stable releases.
- Semantic version bumping and generated GitHub release notes based on
  conventional commit messages.
- README `How it works` diagram for the CLI scan/report/redact/guardrails flow.
- Expanded supported coding-agent catalog across scan target resolution, CLI
  help, guardrail planning, tests, and docs.
- Added shared agent metadata for Cursor, Antigravity CLI, Kimi Code CLI, Droid,
  Gemini CLI, GitHub Copilot, Hermes Agent, OpenClaw, Kilo Code, Kiro CLI, Pi
  CLI, Qoder, Qwen Code, and Trae.
- Added `docs/agent-support-reference.md` to document canonical names, aliases,
  default root heuristics, and guardrail behavior.

### Changed

- Release docs now describe automatic branch-driven release behavior instead of
  manual tag creation as the normal path.
- Guardrails now use the shared agent root markers and generic guardrail files
  for agents without known native instruction paths.
- Documentation now states the `--agent all` and `--root` interaction clearly.

### Fixed

- JSONL redaction now preserves escaped quote delimiters when detected values
  end before JSON-escaped quotes, preventing valid Codex session logs from
  failing structured validation during redaction.

### Verification

- `go test ./...` passed locally.
- JSONL escaped delimiter regression passed locally.
- Manual CI on `dev` passed across Ubuntu, macOS, and Windows.
- Beta release snapshot on `dev` passed.

## Unresolved Questions

None.
