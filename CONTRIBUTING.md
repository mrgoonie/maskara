# Contributing

## Development

```bash
go test ./...
go build ./cmd/maskara
```

Keep scanner behavior deterministic and offline. Do not add runtime network
calls for provider validation unless the feature is explicit and opt-in.

## Pull Requests

- Add tests for new detector rules, report behavior, redaction behavior, or
  guardrail installers.
- Avoid committing real secrets, dotenv files, tokens, private keys, or private
  logs.
- Keep CLI output free of raw secret values.

## Branches

- `main` is stable.
- `dev` is beta.
- Feature branches should target `dev` unless a maintainer says otherwise.
