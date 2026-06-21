# Deployment Guide

## Local Build

```bash
go test ./...
go build -o bin/maskara ./cmd/maskara
```

Before beta or stable shipping, run focused CLI smoke checks for changed public
behavior. For agent-catalog changes:

```bash
go run ./cmd/maskara --help
go run ./cmd/maskara scan --agent gemini --root /tmp/maskara-empty
go run ./cmd/maskara guardrails -a qwen --root /tmp/maskara-home --dry-run
```

## Stable Release

Stable releases come from `main`.

```bash
git checkout main
git tag vX.Y.Z
git push origin vX.Y.Z
```

The release workflow runs GoReleaser and publishes archives/checksums.

## Beta Release

Beta releases come from `dev`.

```bash
git checkout dev
git tag vX.Y.Z-beta.N
git push origin vX.Y.Z-beta.N
```

Tags with prerelease suffixes are marked prerelease by GoReleaser.

For beta branch validation, CI must pass across Ubuntu, macOS, and Windows. The
release workflow also builds a beta snapshot on `dev`.

## Unresolved Questions

None.
