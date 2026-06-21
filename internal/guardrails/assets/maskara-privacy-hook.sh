#!/usr/bin/env sh
set -eu

payload="${*:-$(cat 2>/dev/null || true)}"
lower="$(printf '%s' "$payload" | tr '[:upper:]' '[:lower:]')"

case "$lower" in
  *".env"*cat*|*cat*".env"*|*printenv*|*env\|*|*"secret"*|*"private key"*|*"authorization:"*|*"bearer "*)
    printf '%s\n' "maskara guardrails: command may expose secrets. Use masked checks or run maskara report." >&2
    exit 2
    ;;
esac

exit 0
