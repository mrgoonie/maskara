# Agent Support Reference

## Overview

Maskara keeps supported agent metadata in `internal/agents`. The catalog drives
CLI help, scan target resolution, and guardrail auto-detection.

## Canonical Names

| Canonical name | User-facing agent |
|---|---|
| `claude` | Claude Code |
| `codex` | Codex |
| `cursor` | Cursor |
| `opencode` | OpenCode |
| `antigravity` | Antigravity / Antigravity CLI |
| `kimi` | Kimi Code CLI |
| `droid` | Droid |
| `gemini` | Gemini CLI |
| `github-copilot` | GitHub Copilot |
| `hermes` | Hermes Agent |
| `openclaw` | OpenClaw |
| `kilo` | Kilo Code |
| `kiro` | Kiro CLI |
| `pi` | Pi CLI |
| `qoder` | Qoder |
| `qwen` | Qwen Code |
| `trae` | Trae |

## Alias Rules

Input names are lowercased, trimmed, and normalize spaces or underscores to
hyphens. These aliases map to canonical names:

| Alias examples | Canonical |
|---|---|
| `claude-code`, `claudecode` | `claude` |
| `open-code` | `opencode` |
| `antigravity-cli`, `antigravity code` | `antigravity` |
| `kimi-code`, `kimi-code-cli`, `kimi cli` | `kimi` |
| `gemini-cli`, `gemini cli` | `gemini` |
| `github copilot`, `github-copilot-cli`, `copilot` | `github-copilot` |
| `hermes-agent`, `hermes agent` | `hermes` |
| `open-claw`, `openclaw-cli` | `openclaw` |
| `kilo-code`, `kilo code` | `kilo` |
| `kiro-cli`, `kiro cli` | `kiro` |
| `pi-cli`, `pi cli` | `pi` |
| `qoder-cli`, `qoder cli` | `qoder` |
| `qwen-code`, `qwen code` | `qwen` |
| `trae-cli`, `trae cli` | `trae` |

## Scan Roots

| Agent group | Default roots |
|---|---|
| Claude | `~/.claude/projects` |
| Codex | `~/.codex/sessions` |
| Other catalog agents | home dotdir plus OS config/data dirs |

For catalog agents other than Claude/Codex, Maskara checks a home dotdir first,
then platform locations such as:

- Windows: `%APPDATA%/<AppName>` and `%LOCALAPPDATA%/<AppName>`
- macOS: `~/Library/Application Support/<AppName>`
- Linux: `~/.config/<name>` and `~/.local/share/<name>`

`--root <path>` overrides discovery and scans exactly that directory. With
`--agent all` and no `--root`, Maskara expands to every built-in default target.
With `--root`, Maskara scans the custom directory once.

## Guardrails

Native guardrail paths:

| Agent | Native files |
|---|---|
| Claude | `.claude/CLAUDE.md`, `.claude/skills/maskara-privacy/SKILL.md`, hook |
| Codex | `.codex/AGENTS.md`, `.codex/skills/maskara-privacy/SKILL.md`, hook |
| OpenCode | `.config/opencode/maskara-guardrails.md`, hook |
| Antigravity | `.antigravity/maskara-guardrails.md`, hook |

Other supported agents receive a generic `maskara-guardrails.md` and hook under
their primary catalog root. This gives users a local reference and hook asset,
but it does not guarantee the agent will load the markdown automatically.

## Maintenance Rules

- Add new agents in `internal/agents` first.
- Keep CLI help sourced from `agents.HelpList()`.
- Reuse `agents.RootMarkers()` in guardrail auto-detection.
- Add tests for aliases, default targets, custom-root behavior, and guardrail
  dry-run plans.
- Update this file, README, and roadmap when the catalog changes.

## Unresolved Questions

None.
