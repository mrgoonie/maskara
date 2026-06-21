# Maskara Privacy Guardrails

Never print, quote, summarize, or store raw secrets from `.env` files, shell
history, cloud CLIs, keychains, password managers, session logs, or agent logs.

When a task needs to verify a credential exists, access it by variable name and
report only presence, prefix class, provider, or a masked preview. Prefer
commands that avoid echoing values. Do not paste full API keys, private keys,
database URLs, bearer tokens, cookies, or session tokens into conversation text.

Before sharing logs, reports, repro bundles, or terminal output, run
`maskara scan` or `maskara report` and redact findings. If a secret may have
been exposed to an agent or remote provider, tell the user to rotate it.

If a tool blocks access to a sensitive file, ask the user for explicit approval
instead of trying alternate commands to bypass the block.
