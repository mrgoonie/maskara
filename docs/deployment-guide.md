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

## Automated Releases

Releases are branch-push driven. Do not create release tags manually for normal
shipping.

| Branch | Output |
|---|---|
| `dev` | beta prerelease tag and GitHub release, `vX.Y.Z-beta.N` |
| `main` | stable official tag and GitHub release, `vX.Y.Z` |

The release workflow:

1. Checks out full history and tags.
2. Runs `go test ./...`.
3. Computes the next semantic version from conventional commits since the
   latest stable tag.
4. Generates release notes from conventional commit messages.
5. Creates and pushes the computed tag if it does not already exist.
6. Runs GoReleaser to publish archives and checksums.

Version bump rules:

| Commit signal | Bump |
|---|---|
| `type!:` or `BREAKING CHANGE:` | major |
| `feat:` | minor |
| anything else releasable | patch |

Beta releases use the same base version as the eventual stable release and
increment the beta suffix from existing `vX.Y.Z-beta.N` tags.

## Manual Recovery

If a release job fails after creating a tag, inspect the failed GitHub Actions
run first. Re-run the workflow only after confirming the tag and GitHub release
state. The workflow skips publishing when the computed tag already exists.

## Unresolved Questions

None.
