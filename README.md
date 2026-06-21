![Maskara banner](assets/branding/maskara-banner.png)

# Maskara

[![CI](https://github.com/mrgoonie/maskara/actions/workflows/ci.yml/badge.svg)](https://github.com/mrgoonie/maskara/actions/workflows/ci.yml)

Maskara is an offline CLI for scanning coding-agent session logs, finding
sensitive values, redacting local copies, and producing a rotation report.

Name meaning: `maskara` mixes `mascara` and `mask`. If you do not want to cry
your mascara away after leaking secrets into agent conversations, mask them.

## Install

```bash
go install github.com/mrgoonie/maskara/cmd/maskara@latest
```

Or download a release binary from GitHub Releases.

## Usage

```bash
# full workflow: scan, redact with backups, and write a Markdown report
maskara

# print version
maskara --version
maskara -v

# scan only
maskara scan

# scan only Codex logs
maskara scan --agent codex
maskara scan -a codex

# scan and report
maskara report

# JSON report
maskara report --json

# write report to a file or directory
maskara report --output /path/to/file.md
maskara report -o /path/to/reports/

# install coding-agent guardrails
maskara guardrails
maskara guardrails -a claude
maskara guardrails -a codex
maskara guardrails --dry-run
```

`maskara` with no subcommand runs the full workflow: scan, redact, and report.
Redaction creates `*.maskara.bak` backups before rewriting files.

## Agents

Current target support:

| Agent | Default scan location |
|---|---|
| Claude Code | `~/.claude/projects` |
| Codex | `~/.codex/sessions` |
| OpenCode | platform config/data dirs |
| Antigravity | platform config/data dirs |

Use `--root <path>` to scan a custom log directory.

## Detection

Maskara detects common high-risk patterns:

- OpenAI, Anthropic, GitHub, Slack, Stripe, Google, and AWS token shapes
- JWTs
- private key blocks
- database URLs
- secret-like env assignments

Reports contain masked previews and SHA-256 fingerprints. They do not include
full secret values.

## Guardrails

`maskara guardrails` installs local agent instructions, a small privacy skill,
and hook scripts that discourage commands likely to print secrets. Existing
files are backed up before modification.

Guardrails reduce future leaks. They do not replace secret rotation, provider
revocation, or review of already-shared transcripts.

## Exit Codes

| Code | Meaning |
|---:|---|
| 0 | no findings or guardrail install succeeded |
| 1 | findings detected |
| 2 | runtime error |

## Safety Model

Maskara runs locally and does not validate secrets with remote providers. If a
credential appears in an agent log, assume it may have been exposed and rotate
it. Local redaction removes copies from files on disk only.

## Release Channels

- `main`: stable release branch
- `dev`: beta branch
- stable tags: `vX.Y.Z`
- beta tags: `vX.Y.Z-beta.N`

## References

- [agent-sweep](https://github.com/Ishannaik/agent-sweep)
- [scan-for-secrets](https://github.com/simonw/scan-for-secrets)
- [ggshield](https://github.com/GitGuardian/ggshield)
