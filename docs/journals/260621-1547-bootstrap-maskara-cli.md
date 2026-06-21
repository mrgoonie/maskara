# Bootstrap Maskara CLI Journal

## Context

New repo scaffold for `mrgoonie/maskara`, an offline coding-agent session log
secret scanner, redactor, reporter, and guardrails installer.

## What Happened

- Chose Go for cross-platform single-binary distribution.
- Implemented standard-library CLI, detector rules, filesystem scanner, Markdown/JSON reports, backup redaction, and guardrails installer.
- Generated 5:2 watercolor technical sketch banner and embedded it in README.
- Added CI for `main` and `dev`, plus GoReleaser config for stable and beta tags.
- Created public GitHub repo `mrgoonie/maskara`.
- Verified with unit tests, build, smoke commands, self-scan, and GoReleaser config check.

## Decisions

- Runtime stays offline by default.
- Reports include masked previews and SHA-256 fingerprints, never raw secrets.
- Redaction rejects symlinks and writes backups before replacing files.
- Provider validation and automatic rotation are deferred.

## Unresolved Questions

None.
