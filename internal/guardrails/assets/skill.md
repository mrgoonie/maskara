---
name: maskara-privacy
description: Prevent accidental disclosure of secrets in coding-agent sessions and route cleanup through maskara.
---

# Maskara Privacy Skill

Use this skill whenever a task touches secrets, credentials, `.env` files,
agent session logs, shell history, or generated reports that may contain private
values.

Rules:
- Do not print raw secrets.
- Use masked previews only.
- Prefer checking variable names or key presence over reading values.
- Run `maskara scan` before sharing agent logs.
- Run `maskara report` when the user needs rotation guidance.
- Run `maskara` only when the user wants scan, redaction, and report together.
- If redaction finds credentials, advise rotation. Local redaction does not
  revoke secrets already shared with providers.
