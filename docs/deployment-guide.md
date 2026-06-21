# Deployment Guide

## Local Build

```bash
go test ./...
go build -o bin/maskara ./cmd/maskara
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

## Unresolved Questions

None.
