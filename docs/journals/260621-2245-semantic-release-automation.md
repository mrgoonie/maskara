# Semantic Release Automation

## Context

Release workflow previously built beta snapshots from `dev` and official
releases from manually pushed tags. The desired behavior is branch-push release
automation.

## What Changed

- `dev` pushes compute and publish beta tags like `vX.Y.Z-beta.N`.
- `main` pushes compute and publish stable tags like `vX.Y.Z`.
- Version bump uses conventional commit signals.
- Release notes are generated from commit messages and passed to GoReleaser.
- README now includes a Mermaid `How it works` diagram.

## Safety

- The workflow fetches full tags and skips publishing if the computed tag
  already exists.
- Go tests run before tag creation.
- GoReleaser remains the release publisher for archives/checksums.

## Unresolved Questions

None.
